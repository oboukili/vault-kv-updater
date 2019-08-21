package api

import (
	"encoding/json"
	"github.com/go-yaml/yaml"
	"github.com/jeremywohl/flatten"
	"log"
	vault "github.com/hashicorp/vault/api"
)

func Routine(i interface{}, kvPath string, c *vault.Client) (err error) {
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
	
	// decode yaml contents to map[string]interface{} instance
	decodedContents, err := decode(contents)
	if err != nil {
		return
	}
	
	content, err := json.Marshal(decodedContents)
	if err != nil {
		return
	}
	
	// TODO: introduce a boolean for unicode characters json escaping opt-out
	unescapedContent, err := unescapeUnicodeCharactersInJSON(content)
	if err != nil {
		return
	}
	
	// TODO: introduce a boolean for flattening opt-in
	flattened, err := flatten.FlattenString(string(unescapedContent), "", flatten.DotStyle)
	if err != nil {
		return err
	} else {
		err = VaultKVIdempotentWrite(flattened, kvPath, c)
		if err != nil {
			log.Fatalln(err)
		}
	}
	return
}

