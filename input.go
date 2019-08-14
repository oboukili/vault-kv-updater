package main

import (
	"errors"
	"io/ioutil"
	"os"
)

func readInput() (input []byte, err error) {
	if len(os.Args) > 1 {
		if input, err = ioutil.ReadFile(os.Args[1]); err != nil {
			return nil, err
		}
	} else if stdinStat, _ := os.Stdin.Stat(); (stdinStat.Mode() & os.ModeCharDevice) == 0 {
		if input, err = ioutil.ReadAll(os.Stdin); err != nil {
			return nil, err
		}
	} else {
		err = errors.New("missing input")
	}
	return
}
