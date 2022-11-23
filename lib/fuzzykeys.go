package lib

import "strings"

func FuzzyKeys[V any](m map[string]V, key string) (V, bool) {
	value, ok := m[key]
	if ok {
		return value, ok
	}

	value, ok = m[strings.ToLower(strings.Replace(key, "-", "_", -1))]
	return value, ok
}
