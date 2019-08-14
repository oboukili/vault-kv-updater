package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/jeremywohl/flatten"
	"log"
)

func main() {
	flag.Parse()

	input, err := readInput()
	if err != nil {
		log.Fatalln(err)
	}
	
	// unmarshal yaml, because we know we are working yaml
	var contents interface{}
	if err = yaml.Unmarshal(input, &contents); err != nil {
		log.Fatalln(err)
	}

	// decode yaml contents to map[string]interface{} instance
	decodedContents, err := decode(contents)
	if err != nil {
		log.Fatalln(err)
	}
	
	// check whether we are working with sops encrypted data
	if isSopsEncrypted(decodedContents) {
		decodedContents, err = decrypt(decodedContents, "yaml")
		if err != nil {
			log.Fatalln(err)
		}
	}
	
	if content, err := json.Marshal(decodedContents); err != nil {
		log.Fatalln(err)
	} else {
		if unescapedContent, err := _UnescapeUnicodeCharactersInJSON(json.RawMessage(content)); err != nil {
			log.Fatalln(err)
		} else {
			if flattened, err := flatten.FlattenString(string(unescapedContent), "", flatten.DotStyle); err != nil {
				log.Fatalln(err)
			} else {
				fmt.Println(string(flattened))
			}
		}
	}
}
