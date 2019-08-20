package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const EnvAutoCompleteFilePrefix = "AUTO_COMPLETE_FILE_PREFIX"
const EnvAutoCompleteVaultKVMount = "AUTO_COMPLETE_VAULT_KV_MOUNT"
const EnvAutoCompleteVaultKVPathPrefix = "AUTO_COMPLETE_VAULT_KV_PATH_PREFIX"

var validFileExtensions = [...]string{".yml", ".yaml"}

type ExtendedFileInfo struct {
	os.FileInfo
	FilePath    string
	VaultKVPath string
}

func (e *ExtendedFileInfo) CompleteVaultKVPath() {
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
	(*e).VaultKVPath = b.String()
}

func AutoCompleteGetFiles(folders *[]ExtendedFileInfo) (results *[]ExtendedFileInfo, err error) {
	if folders == nil || len(*folders) == 0 {
		err = fmt.Errorf("at least one folder has to be specified as an argument")
		return
	}
	for _, f := range *folders {
		if !f.IsDir() {
			err = fmt.Errorf("%s should be a directory when using autocomplete mode", f.Name())
			return
		}
		err = filepath.Walk(f.FilePath, func(path string, info os.FileInfo, err error) (walkFnErr error) {
			for _, s := range validFileExtensions {
				if strings.HasSuffix(info.Name(), s) && !info.IsDir() {
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
			return
		}
		for _, r := range *results {
			r.CompleteVaultKVPath()
		}
	}
	return
}

func AutoCompleteInit(p *[]ExtendedFileInfo) (err error) {

	if len(os.Args) < 2 {
		err = fmt.Errorf("at least one folder should be specified as argument when using autocomplete mode")
		return
	}

	*p = make([]ExtendedFileInfo, len(os.Args)-1)

	kvMount, ok := os.LookupEnv(EnvAutoCompleteVaultKVMount)
	if !ok || strings.EqualFold(kvMount, "") {
		err = fmt.Errorf("%s should be set as a non-empty string environment variable", EnvAutoCompleteVaultKVMount)
	}
	for i, path := range os.Args[:1] {
		f, err := os.Stat(path)
		if err != nil {
			err = fmt.Errorf("%s: %s", path, err)
			return err
		}
		(*p)[i] = ExtendedFileInfo{
			FileInfo: f,
			FilePath: path,
		}
	}
	return
}
