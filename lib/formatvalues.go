package lib

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/dustin/go-humanize"
)

func formatValue(value interface{}) interface{} {
	if value == nil {
		value = ""
	}

	if reflect.ValueOf(value).Kind() == reflect.Map {
		jsonBytes, err := json.Marshal(value)
		if err == nil {
			return string(jsonBytes)
		}
	}

	switch value.(type) {
	case float64:
		value = int64(value.(float64))
	}
	return value
}

func formatValuePretty(key string, value interface{}) interface{} {
	if key == "size" {
		return formatSize(value)
	}

	return formatValue(value)
}

// displayCell formats a record value for human table display, escaping any
// terminal control characters carried by the value.
func displayCell(key string, value interface{}) string {
	return escapeTerminalControls(fmt.Sprintf("%v", formatValuePretty(key, value)))
}

func formatSize(value interface{}) interface{} {
	switch value.(type) {
	case float64:
		value = humanize.Bytes(uint64(value.(float64)))
	case string:
		v, err := strconv.ParseInt(value.(string), 10, 64)
		if err == nil {
			value = humanize.Bytes(uint64(v))
		}
	case nil:
		value = ""
	}
	return value
}
