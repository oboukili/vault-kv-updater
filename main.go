package main

import (
	"errors"
	"fmt"
	"gitlab.com/maltcommunity/public/vault-kv-updater.git/api"
	"log"
	"os"
	"strconv"
)

func main() {
	var kvPath string
	var autoComplete bool
	var err error

	kvVersionString, ok := os.LookupEnv(api.EnvVaultKvVersion)
	if !ok {
		kvVersionString = "1"
	}

	kvMount, ok := os.LookupEnv(api.EnvVaultKvMount)
	if !ok {
		log.Fatalln(fmt.Errorf("%s environment variable must be specified", api.EnvVaultKvMount))
	}

	kvVersion, err := strconv.Atoi(kvVersionString)
	if err != nil {
		log.Fatalln(fmt.Errorf("%s environment variable must be boolean compatible: %s", api.EnvVaultKvVersion, kvVersionString))
	}

	// Autocomplete mode
	ac, ok := os.LookupEnv(api.EnvAutoComplete)
	if ok {
		autoComplete, err = strconv.ParseBool(ac)
		if err != nil {
			log.Fatalln(fmt.Errorf("%s environment variable must be boolean compatible: %s", api.EnvAutoComplete, ac))
		}
		if autoComplete {
			if err := api.AutoCompleteInit(); err != nil {
				log.Fatalln(err)
			}
		}
	}
	// Single secret mode
	if !autoComplete {
		if len(os.Args) > 2 {
			log.Fatalln("only one YAML document may be specified at a time when not using autocomplete mode")
		}
		kvPath, ok = os.LookupEnv(api.EnvVaultKvPath)
		if !ok {
			log.Fatalln(fmt.Errorf("%s environment variable must be specified", api.EnvVaultKvPath))
		}
	}

	// Vault client initialization
	c, err := api.VaultClientInit()
	if err != nil {
		log.Fatalln(err)
	}

	// Run main routines
	if len(os.Args) > 1 {
		switch autoComplete {
		case false:
			log.Print("Simple mode enabled")
			for _, file := range os.Args[1:] {
				err := api.Routine(file, kvMount, kvVersion, kvPath, c)
				if err != nil {
					log.Fatalln(err)
				}
			}
		case true:
			log.Print("Autocomplete mode enabled")
			r, err := api.AutoCompleteGetFiles(os.Args[1:])
			if err != nil {
				log.Fatalln(err)
			}
			for _, f := range *r {
				// TODO: use goroutines + channels for major speedups

				err := api.Routine(f.FilePath, kvMount, kvVersion, f.VaultKVPath, c)
				if err != nil {
					log.Fatalln(err)
				}
			}
		}
	} else if stdinStat, _ := os.Stdin.Stat(); (stdinStat.Mode() & os.ModeCharDevice) == 0 {
		if err := api.Routine(os.Stdin, kvMount, kvVersion, kvPath, c); err != nil {
			log.Fatalln(err)
		}
	} else {
		// TODO: implement a proper help message with a cli library
		log.Fatalln(errors.New("missing input"))
	}
	log.Print("All good! ;)")
}
