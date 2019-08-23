package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var validFileExtensions = [...]string{".yml", ".yaml"}

type ExtendedFileInfo struct {
	os.FileInfo
	FilePath    string
	VaultKVPath string
}

func GetVaultKVPath(e *ExtendedFileInfo) (path string, err error) {
	var r string
	var b strings.Builder

	filePrefix, ok := os.LookupEnv(EnvAutoCompleteFilePrefix)
	if !ok {
		filePrefix = ""
	}

	kvPathPrefix, ok := os.LookupEnv(EnvAutoCompleteVaultKVPathPrefix)
	if ok && kvPathPrefix != "" {
		b.WriteString(kvPathPrefix)
		b.WriteString("/")
	}
	if !strings.EqualFold(filePrefix, "") && strings.HasPrefix((*e).Name(), filePrefix) {
		r = strings.Replace(e.Name(), filePrefix, "", 1)
	} else {
		r = e.Name()
	}

	for _, suffixExpression := range validFileExtensions {
		if strings.HasSuffix(r, suffixExpression) {
			r = strings.Replace(r, suffixExpression, "", 1)
		}
	}
	_, err = b.WriteString(r)
	if err != nil {
		return
	}
	return b.String(), nil
}

func AutoCompleteGetFiles(directories []string) (*[]ExtendedFileInfo, error) {

	if directories == nil || len(directories) == 0 {
		err := fmt.Errorf("at least one folder has to be specified as an argument")
		return nil, err
	}

	results := make([]ExtendedFileInfo, 0)

	for _, d := range directories {
		di, err := os.Stat(d)
		//f.FilePath = d
		if err != nil {
			if os.IsNotExist(err) {
				log.Printf("WARN: path does not exist: %s, ignoring...", d)
				continue
			}
			return nil, fmt.Errorf("%s: %s", d, err)
		}
		if !di.IsDir() {
			err = fmt.Errorf("ERROR: %s should be a directory when using autocomplete mode", di.Name())
			return nil, err
		}
		files, err := ioutil.ReadDir(d)
		if err != nil {
			return nil, err
		}
		var b strings.Builder
		for _, f := range files {
			if f.IsDir() {
				continue
			}
			for _, s := range validFileExtensions {
				if strings.HasSuffix(f.Name(), s) {
					b.WriteString(d)
					b.WriteString("/")
					b.WriteString(f.Name())
					results = append(results, ExtendedFileInfo{
						FileInfo: f,
						FilePath: b.String(),
					})
					b.Reset()
					continue
				}
			}
		}
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("ERROR: no file was detected to be processed during autocomplete mode")
	}
	//sort.Slice(results, func(i, j int) bool {
	//	return results[i].FilePath > results[j].FilePath
	//})
	for i, _ := range results {
		path, err := GetVaultKVPath(&results[i])
		if err != nil {
			return nil, err
		}
		results[i].VaultKVPath = path
		if results[i].VaultKVPath == "" {
			err = fmt.Errorf("ERROR: autocomplete computed path was empty for file %s", results[i].Name())
		}
	}
	return &results, nil
}

func AutoCompleteInit() (err error) {

	if len(os.Args) < 2 {
		err = fmt.Errorf("at least one folder should be specified as argument when using autocomplete mode")
		return
	}
	kvMount, ok := os.LookupEnv(EnvVaultKvMount)
	if !ok || strings.EqualFold(kvMount, "") {
		err = fmt.Errorf("%s should be set as a non-empty string environment variable", EnvVaultKvMount)
	}
	return
}
