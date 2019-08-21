package api

import (
	"fmt"
	"reflect"
	"strconv"
)

func decode(input interface{}) (interface{}, error) {
	var err error
	switch in := input.(type) {
	case map[interface{}]interface{}:
		rec := map[string]interface{}{}
		for k, v := range in {
			switch k.(type) {
			case bool:
				rec[strconv.FormatBool(k.(bool))], err = decode(v)
			case int:
				rec[strconv.Itoa(k.(int))], err = decode(v)
			case float64:
				rec[fmt.Sprintf("%f", k.(float64))], err = decode(v)
			case string:
				rec[k.(string)], err = decode(v)
			default:
				return nil, fmt.Errorf("decode: unsupported key type: %s", reflect.TypeOf(k))
			}
		}
		return rec, err
	case []interface{}:
		for k, v := range in {
			in[k], err = decode(v)
		}
	}
	return input, err
}
