package lib

import (
	"encoding/json"
	"strings"

	"github.com/Files-com/files-cli/lib/clierr"
)

func parseJSONFlag(flagName string, value string) (interface{}, error) {
	if strings.TrimSpace(value) == "" {
		return nil, clierr.Errorf(clierr.ErrorCodeUsage, "invalid value for --%s: expected valid JSON", flagName)
	}

	var parsed interface{}
	if err := json.Unmarshal([]byte(value), &parsed); err != nil {
		return nil, clierr.Errorf(clierr.ErrorCodeUsage, "invalid value for --%s: expected valid JSON: %v", flagName, err)
	}

	return parsed, nil
}

func ParseJSONObjectFlag(flagName string, value string) (map[string]interface{}, error) {
	parsed, err := parseJSONFlag(flagName, value)
	if err != nil {
		return nil, err
	}
	if parsed == nil {
		return nil, nil
	}

	object, ok := parsed.(map[string]interface{})
	if !ok {
		return nil, clierr.Errorf(clierr.ErrorCodeUsage, "invalid value for --%s: expected JSON object", flagName)
	}

	return object, nil
}

func ParseJSONArrayObjectFlag(flagName string, value string) ([]map[string]interface{}, error) {
	parsed, err := parseJSONFlag(flagName, value)
	if err != nil {
		return nil, err
	}
	if parsed == nil {
		return nil, nil
	}

	items, ok := parsed.([]interface{})
	if !ok {
		return nil, clierr.Errorf(clierr.ErrorCodeUsage, "invalid value for --%s: expected JSON array of objects", flagName)
	}

	objects := make([]map[string]interface{}, 0, len(items))
	for index, item := range items {
		object, ok := item.(map[string]interface{})
		if !ok {
			return nil, clierr.Errorf(clierr.ErrorCodeUsage, "invalid value for --%s: item %d must be a JSON object", flagName, index+1)
		}
		objects = append(objects, object)
	}

	return objects, nil
}
