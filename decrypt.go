package main

import (
	"fmt"
	sops "go.mozilla.org/sops/decrypt"
)

func isSopsEncrypted(input interface{}) bool {
	switch input.(type) {
	case map[string]interface{}:
		if _, ok := input.(map[string]interface{})["sops"]; ok {
			switch input.(map[string]interface{})["sops"].(type) {
			case map[string]interface{}:
				_, ok := input.(map[string]interface{})["sops"].(map[string]interface{})["version"]
				return ok
			}
		}
	}
	return false
}

func decrypt(input interface{}, format string) (contents []byte, err error) {
	switch input.(type) {
	case []byte:
		return sops.Data(input.([]byte), format)
	case string:
		return sops.File(input.(string), format)
	default:
		return nil, fmt.Errorf("unsupported input type")
	}
}
