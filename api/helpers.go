package api

import (
	"fmt"
	"github.com/jeremywohl/flatten"
	"os"
	"strconv"
)

func Mapper(i interface{}) interface{} {
	m := make(map[string]interface{})
	switch i.(type) {
	case map[interface{}]interface{}:
		for k, v := range i.(map[interface{}]interface{}) {
			switch k := k.(type) {
			case string:
				switch v := v.(type) {
				case map[interface{}]interface{}:
					m[k] = Mapper(v)
				// Unfortunately, Json Marshaller infers float64 type for every numeric value from the Vault API
				// (of type encoding/json.Number)
				// To ensure strict structure equality to ensure idempotence, we have to convert those to float64
				case int:
					m[k] = float64(v)
				case uint8:
					m[k] = float64(v)
				case uint16:
					m[k] = float64(v)
				case uint32:
					m[k] = float64(v)
				case uint64:
					m[k] = float64(v)
				case int8:
					m[k] = float64(v)
				case int16:
					m[k] = float64(v)
				case int32:
					m[k] = float64(v)
				case int64:
					m[k] = float64(v)
				case float32:
					m[k] = float64(v)
				default:
					m[k] = v
				}
			}
		}
	}
	return m
}

func FlattenMap(data map[string]interface{}) (flattened map[string]interface{}, err error) {
	shouldFlatten := true
	shouldFlattenString, ok := os.LookupEnv(EnvFlatten)
	if ok {
		shouldFlatten, err = strconv.ParseBool(shouldFlattenString)
		if err != nil {
			return nil, fmt.Errorf("ERROR: %s value %s should be boolean compatible: %s", EnvFlatten, shouldFlattenString, err)
		}
	}
	switch shouldFlatten {
	case true:
		flattened, err = flatten.Flatten(data, "", flatten.DotStyle)
		if err != nil {
			return nil, err
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
			r = r[0 : len(r)-2]
		}
	}
	return
}
