package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-yaml/yaml"
	vault "github.com/hashicorp/vault/api"
	"github.com/jeremywohl/flatten"
	"log"
	"os"
)

func routine(i interface{}, kvPath string, c *vault.Client) (err error) {
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
		err = vaultKVIdempotentWrite(flattened, kvPath, c)
		if err != nil {
			log.Fatalln(err)
		}
	}
	return
}

func main() {
	kvPath, ok := os.LookupEnv("VAULT_KV_PATH")
	if !ok {
		log.Fatalln("VAULT_KV_PATH must be specified")
	}

	c := vaultClientInit()
	// Run main routines
	if len(os.Args) > 1 {
		for _, f := range os.Args[1:] {
			// TODO: use a goroutine + channels for major speedups
			err := routine(f, kvPath, c)
			if err != nil {
				log.Fatalln(fmt.Errorf("main: %s", err))
			}
		}
	} else if stdinStat, _ := os.Stdin.Stat(); (stdinStat.Mode() & os.ModeCharDevice) == 0 {
		if err := routine(os.Stdin, kvPath, c); err != nil {
			log.Fatalln(err)
		}
	} else {
		// TODO: implement a proper help message with a cli library
		log.Fatalln(errors.New("missing input"))
	}
}
