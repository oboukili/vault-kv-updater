package main

import (
	
	"errors"
	"fmt"
	"gitlab.com/maltcommunity/ops/vault-kv-updater.git/api"
	"log"
	"os"
	"strconv"
)

func main() {
	var kvPath string
	var autoComplete bool

	ac, ok := os.LookupEnv("AUTO_COMPLETE")
	switch ok {
	case true:
		autoComplete, err := strconv.ParseBool(ac)
		if err != nil {
			log.Fatalln(fmt.Errorf("AUTO_COMPLETE environment variable must be boolean compatible: %s", ac))
		}
		if autoComplete {
			if err := AutoCompleteInit(); err != nil {
				log.Fatalln(err)
			}
		}
	case false:
		if len(os.Args) > 2 {
			log.Fatalln("only one YAML document may be specified at a time when not using autocomplete mode")
		}

		kvPath, ok = os.LookupEnv("VAULT_KV_PATH")
		if !ok {
			log.Fatalln("VAULT_KV_PATH must be specified")
		}
	}

	// Vault client initialization
	c, err := VaultClientInit()
	if err != nil {
		log.Fatalln(err)
	}

	// Run main routines
	if len(os.Args) > 1 {
		switch autoComplete {
		case false:
			for _, file := range os.Args[1] {
				err := Routine(file, kvPath, c)
				if err != nil {
					log.Fatalln(err)
				}
			}
		case true:
			r, err := AutoCompleteGetFiles(os.Args[:1])
			if err != nil {
				log.Fatalln(err)
			}
			for _, f := range *r {
				// TODO: use goroutines + channels for major speedups
				err := Routine(f.FilePath, f.VaultKVPath(), c)
				if err != nil {
					log.Fatalln(err)
				}
			}
		}
	} else if stdinStat, _ := os.Stdin.Stat(); (stdinStat.Mode() & os.ModeCharDevice) == 0 {
		if err := Routine(os.Stdin, kvPath, c); err != nil {
			log.Fatalln(err)
		}
	} else {
		// TODO: implement a proper help message with a cli library
		log.Fatalln(errors.New("missing input"))
	}
}
