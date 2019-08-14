package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/jeremywohl/flatten"
	"log"
	"os"
)

func routine(i interface{}) (err error) {
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
	if err := yaml.Unmarshal(input, &contents); err != nil {
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

	unescapedContent, err := unescapeUnicodeCharactersInJSON(json.RawMessage(content))
	if err != nil {
		return
	}

	if flattened, err := flatten.FlattenString(string(unescapedContent), "", flatten.DotStyle); err != nil {
		return err
	} else {
		fmt.Println(flattened)
	}
	return
}

func main() {
	flag.Parse()

	if len(os.Args) > 1 {
		for _, f := range os.Args[1:] {
			// TODO: use a goroutine + channels for major speedups
			err := routine(f)
			if err != nil {
				log.Fatalln(fmt.Errorf("main: %s", err))
			}
		}
	} else if stdinStat, _ := os.Stdin.Stat(); (stdinStat.Mode() & os.ModeCharDevice) == 0 {
		if err := routine(os.Stdin); err != nil {
			log.Fatalln(err)
		}
	} else {
		// TODO: implement a proper help message with a cli library
		log.Fatalln(errors.New("missing input"))
	}
}
