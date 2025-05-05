package lib

import "strings"

func FuzzyKeys[V any](m map[string]V, key string) (V, bool) {
	value, ok := m[key]
	if ok {
		return value, ok
	}

	value, ok = m[strings.ToLower(strings.ReplaceAll(key, "-", "_"))]
	return value, ok
}
