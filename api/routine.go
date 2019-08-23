package api

import (
	"github.com/go-yaml/yaml"
	vault "github.com/hashicorp/vault/api"
	"log"
)

func Routine(i interface{}, kvMount string, kvVersion int, kvPath string, c *vault.Client) (err error) {
	ok, err, input := isSopsEncrypted(i)
	if err != nil {
		return
	}
	if ok {
		input, err = decrypt(input, "yaml")
		if err != nil {
			return
		}
	}
	// unmarshal yaml, because we know we are working with yaml
	var contents interface{}
	if err := yaml.Unmarshal(*input, &contents); err != nil {
		return err
	}

	// Convert contents to a manipulable map type
	convertedContents := Mapper(contents)
	flattened, err := FlattenMap(convertedContents.(map[string]interface{}))
	if err != nil {
		return err
	}

	err = VaultKVIdempotentWrite(flattened, kvMount, kvVersion, kvPath, c)
	if err != nil {
		log.Fatalln(err)
	}
	return
}
