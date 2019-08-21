package api

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	EnvAutoCompleteFilePrefix        = "AUTO_COMPLETE_FILE_PREFIX"
	EnvAutoCompleteVaultKVMount      = "AUTO_COMPLETE_VAULT_KV_MOUNT"
	EnvAutoCompleteVaultKVPathPrefix = "AUTO_COMPLETE_VAULT_KV_PATH_PREFIX"
)

var validFileExtensions = [...]string{".yml", ".yaml"}

type ExtendedFileInfo struct {
	os.FileInfo
	FilePath string
}

func (e *ExtendedFileInfo) VaultKVPath() string {
	var r string

	filePrefix, ok := os.LookupEnv(EnvAutoCompleteFilePrefix)
	if !ok {
		filePrefix = ""
	}
	b := strings.Builder{}
	b.WriteString(EnvAutoCompleteVaultKVMount)
	b.WriteString("/")
	kvPathPrefix, ok := os.LookupEnv(EnvAutoCompleteVaultKVPathPrefix)
	if ok && kvPathPrefix != "" {
		b.WriteString(kvPathPrefix)
		b.WriteString("/")
	}
	if !strings.EqualFold(filePrefix, "") && strings.HasPrefix((*e).Name(), filePrefix) {
		r = strings.Replace((*e).Name(), filePrefix, "", 1)
	}
	for _, suffixExpression := range validFileExtensions {
		if strings.HasSuffix(r, suffixExpression) {
			r = strings.Replace(r, suffixExpression, "", 1)
		}
	}
	b.WriteString(r)
	return b.String()
}

func AutoCompleteGetFiles(directories []string) (results *[]ExtendedFileInfo, err error) {

	if directories == nil || len(directories) == 0 {
		err = fmt.Errorf("at least one folder has to be specified as an argument")
		return
	}

	*(results) = make([]ExtendedFileInfo, len(directories)-1)

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
			err = fmt.Errorf("%s should be a directory when using autocomplete mode", di.Name())
			return nil, err
		}
		err = filepath.Walk(d, func(path string, info os.FileInfo, err error) (walkFnErr error) {
			for _, s := range validFileExtensions {
				if info.IsDir() {
					return filepath.SkipDir
				}
				if strings.HasSuffix(info.Name(), s) {
					*results = append(*results, ExtendedFileInfo{
						FileInfo: info,
						FilePath: path,
					})
					return
				}
			}
			return
		})
		if err != nil {
			return nil, err
		}
	}
	return
}

func AutoCompleteInit() (err error) {

	if len(os.Args) < 2 {
		err = fmt.Errorf("at least one folder should be specified as argument when using autocomplete mode")
		return
	}
	kvMount, ok := os.LookupEnv(EnvAutoCompleteVaultKVMount)
	if !ok || strings.EqualFold(kvMount, "") {
		err = fmt.Errorf("%s should be set as a non-empty string environment variable", EnvAutoCompleteVaultKVMount)
	}
	return
}
