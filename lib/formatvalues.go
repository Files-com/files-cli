package lib

import (
	"encoding/json"
	"reflect"
	"strconv"
)

func formatValues(key string, value interface{}) interface{} {
	if value == nil {
		value = ""
	}

	if reflect.ValueOf(value).Kind() == reflect.Map {
		jsonBytes, err := json.Marshal(value)
		if err == nil {
			return string(jsonBytes)
		}
	}

	switch key {
	case "size":
		value = formatSize(value)
	default:
		return value
	}
	return value
}

func formatSize(value interface{}) interface{} {
	switch value.(type) {
	case float64:
		value = ByteCountSIFloat64(value.(float64))
	case string:
		v, err := strconv.ParseInt(value.(string), 10, 64)
		if err == nil {
			value = ByteCountSI(v)
		}
	}
	return value
}
