package main

import (
	"fmt"
	"github.com/go-yaml/yaml"
	sops "go.mozilla.org/sops/decrypt"
	"io/ioutil"
	"log"
	"reflect"
)

func isSopsEncryptedYamlFile(path string) (bool ,error) {
	input, err := ioutil.ReadFile(path)
	if err != nil {
		return false, fmt.Errorf("isSopsEncryptedYamlFile; %s", err)
	}

	var contents map[interface{}]interface{}
	if err := yaml.Unmarshal(input, &contents); err != nil {
		return false, fmt.Errorf("isSopsEncryptedYamlFile; %s", err)
	}
	
	if isSopsEncryptedData(contents) {
		log.Printf("Detected Sops encrypted data for %s", path)
		return true, nil
	}
	return false, nil
}

func isSopsEncryptedData(input interface{}) bool {
	switch input.(type) {
	//TODO? more readable type assertions?
	case map[interface{}]interface{}:
		if _, ok := input.(map[interface{}]interface{})["sops"]; ok {
			switch input.(map[interface{}]interface{})["sops"].(type) {
			case map[interface{}]interface{}:
				_, ok := input.(map[interface{}]interface{})["sops"].(map[interface{}]interface{})["version"]
				return ok
			}
		}
	}
	return false
}

func decrypt(input interface{}, format string) (contents []byte, err error) {
	switch input.(type) {
	case []byte:
		if contents, err = sops.Data(input.([]byte), format); err != nil {
			return nil, fmt.Errorf("decrypt: cannot sops decrypt data: %s", err)
		}
		return
	case string:
		if contents, err = sops.File(input.(string), format); err != nil {
			return nil, fmt.Errorf("decrypt: cannot sops decrypt file: %s", err)
		}
		return
	default:
		return nil, fmt.Errorf("decrypt: unsupported input type: %s", reflect.TypeOf(input))
	}
}
