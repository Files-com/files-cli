package lib

import (
	"strings"

	"github.com/Files-com/files-cli/lib/clierr"
)

func parseAPIListQueryValue(flagName string, value string) (string, string, error) {
	key, fieldValue, hasDelimiter := strings.Cut(value, "=")
	if !hasDelimiter {
		return "", "", clierr.Errorf(clierr.ErrorCodeUsage, "invalid value for --%s: expected key=value", flagName)
	}

	key = strings.TrimSpace(key)
	fieldValue = strings.TrimSpace(fieldValue)
	if key == "" {
		return "", "", clierr.Errorf(clierr.ErrorCodeUsage, "invalid value for --%s: key cannot be empty", flagName)
	}

	return key, fieldValue, nil
}

func ParseAPIListQueryFlag(flagName string, values []string) (map[string]interface{}, error) {
	if len(values) == 0 {
		return nil, nil
	}

	parsed := make(map[string]interface{}, len(values))

	for _, value := range values {
		key, fieldValue, err := parseAPIListQueryValue(flagName, value)
		if err != nil {
			return nil, err
		}

		parsed[key] = fieldValue
	}

	return parsed, nil
}

func ParseAPIListSortFlag(flagName string, value string) (map[string]interface{}, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}

	key, fieldValue, err := parseAPIListQueryValue(flagName, value)
	if err != nil {
		return nil, err
	}
	if fieldValue != "asc" && fieldValue != "desc" {
		return nil, clierr.Errorf(clierr.ErrorCodeUsage, "invalid value for --%s: sort direction must be asc or desc", flagName)
	}

	return map[string]interface{}{key: fieldValue}, nil
}
