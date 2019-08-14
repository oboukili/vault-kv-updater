package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/jeremywohl/flatten"
	"io/ioutil"
	"log"
	"os"
)

func processFile(path string) error {
	var input []byte

	ok, err := isSopsEncryptedYamlFile(path)
	if err != nil {
		return fmt.Errorf("processFile: %s", err)
	}
	if ok {
		input, err = decrypt(path, "yaml")
		if err != nil {
			return fmt.Errorf("processFile: %s", err)
		}
	}

	// unmarshal yaml, because we know we are working with yaml
	var contents interface{}
	if err := yaml.Unmarshal(input, &contents); err != nil {
		log.Fatalln(err)
	}

	// decode yaml contents to map[string]interface{} instance
	decodedContents, err := decode(contents)
	if err != nil {
		return fmt.Errorf("decode: %s", err)
	}

	content, err := json.Marshal(decodedContents)
	if err != nil {
		return fmt.Errorf("jsonMarshal: %s", err)
	}

	unescapedContent, err := unescapeUnicodeCharactersInJSON(json.RawMessage(content))
	if err != nil {
		return fmt.Errorf("unescapeUnicodeCharactersInJSON: %s", err)
	}

	if flattened, err := flatten.FlattenString(string(unescapedContent), "", flatten.DotStyle); err != nil {
		return fmt.Errorf("flatten: %s", err)
	} else {
		fmt.Println(string(flattened))
	}
	return nil
}

func processStdin() error {
	_, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	flag.Parse()

	if len(os.Args) > 1 {
		for _, f := range os.Args[1:] {
			// TODO: use a goroutine + channels for major speedups
			err := processFile(f)
			if err != nil {
				log.Fatalln(fmt.Errorf("main: %s", err))
			}
		}
	} else if stdinStat, _ := os.Stdin.Stat(); (stdinStat.Mode() & os.ModeCharDevice) == 0 {
		if err := processStdin(); err != nil {
			log.Fatalln(err)
		}
	} else {
		log.Fatalln(errors.New("missing input"))
	}
}
