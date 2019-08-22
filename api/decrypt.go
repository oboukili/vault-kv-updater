package api

import (
	"fmt"
	"github.com/go-yaml/yaml"
	sops "go.mozilla.org/sops/decrypt"
	"io/ioutil"
	"os"
	"reflect"
)

// Asserts whether the input (stdin or file) is a sops encrypted stream, also returns a pointer to the read data
func isSopsEncrypted(i interface{}) (ok bool, err error, p *[]byte) {

	var input []byte
	p = &input

	switch i.(type) {
	case string:
		input, err = ioutil.ReadFile(i.(string))
	case *os.File:
		input, err = ioutil.ReadAll(i.(*os.File))
	}
	if err != nil {
		err = fmt.Errorf("ERROR: isSopsEncrypted: %s", err)
		return
	}

	var contents map[interface{}]interface{}
	err = yaml.Unmarshal(input, &contents)
	if err != nil {
		err = fmt.Errorf("ERROR: isSopsEncrypted: %s", err)
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

func decrypt(data interface{}, format string) (p *[]byte, err error) {
	var contents []byte
	p = &contents

	switch data.(type) {
	case *[]byte:
		if contents, err = sops.Data(*data.(*[]byte), format); err != nil {
			return nil, fmt.Errorf("decrypt: cannot sops decrypt data: %s", err)
		}
		return
	default:
		return nil, fmt.Errorf("decrypt: unsupported input type: %s", reflect.TypeOf(data))
	}
}
