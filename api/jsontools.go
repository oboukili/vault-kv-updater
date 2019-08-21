package api

import (
	"encoding/json"
	"strconv"
	"strings"
)

// TODO: find a way to tell json.Marshal to escape UnicodeErrors instead
func unescapeUnicodeCharactersInJSON(_jsonRaw json.RawMessage) (json.RawMessage, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(_jsonRaw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}
