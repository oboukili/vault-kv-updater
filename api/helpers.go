package api

import (
	"encoding/json"
	"fmt"
	"github.com/jeremywohl/flatten"
	"os"
	"strconv"
	"strings"
)

// TODO: find a way to tell json.Marshal to escape UnicodeErrors instead
func UnescapeUnicodeCharactersInJSON(_jsonRaw json.RawMessage) (json.RawMessage, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(_jsonRaw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}

func Flatten(data string) (flattened string, err error) {
	shouldFlatten := true
	shouldFlattenString, ok := os.LookupEnv(EnvFlatten)
	if ok {
		shouldFlatten, err = strconv.ParseBool(shouldFlattenString)
		if err != nil {
			return "", fmt.Errorf("ERROR: %s value %s should be boolean compatible: %s", EnvFlatten, shouldFlattenString, err)
		}
	}
	switch shouldFlatten {
	case true:
		flattened, err = flatten.FlattenString(data, "", flatten.DotStyle)
		if err != nil {
			return "", err
		}
	case false:
		return data, nil
	}
	return
}

func TrimQuotes(s string) (r string) {
	r = s
	if len(s) > 2 {
		if s[0] == '"' {
			r = s[1:]
		}
		if s[len(s)-1] == '"' {
			r = r[0:len(r)-2]
		}
	}
	return
}