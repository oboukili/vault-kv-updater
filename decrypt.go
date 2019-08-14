package main

import (
	"fmt"
	"github.com/go-yaml/yaml"
	sops "go.mozilla.org/sops/decrypt"
	"io/ioutil"
	"os"
	"reflect"
)

// Asserts whether the input (stdin or file) is a sops encrypted stream, returns
func isSopsEncrypted(i interface{}) (ok bool, err error, input []byte) {

	switch i.(type) {
	case string:
		input, err = ioutil.ReadFile(i.(string))
	case *os.File:
		input, err = ioutil.ReadAll(i.(*os.File))
	}
	if err != nil {
		err = fmt.Errorf("isSopsEncrypted; %s", err)
		return
	}

	var contents map[interface{}]interface{}
	err = yaml.Unmarshal(input, &contents)
	if err != nil {
		err = fmt.Errorf("isSopsEncrypted; %s", err)
		return
	}

	if _, ok = contents["sops"]; ok {
		switch contents["sops"].(type) {
		case map[interface{}]interface{}:
			_, ok = contents["sops"].(map[interface{}]interface{})["version"]
		}
	}
	return
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
