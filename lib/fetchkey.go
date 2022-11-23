package lib

import (
	"fmt"

	"github.com/sc0vu/didyoumean"
)

func FetchKey[V any](entity string, m map[string]V, key string) (V, error) {
	value, ok := FuzzyKeys(m, key)

	if !ok {
		var keys []string
		for k := range m {
			keys = append(keys, k)
		}
		didyoumean.ThresholdRate = 0.6
		didyoumean.CaseInsensitive = true
		knownKey := didyoumean.FirstMatch(key, keys)
		if knownKey != "" {
			return value, fmt.Errorf(
				`unknown flag "%v" for "%v"

Did you mean this?
	"%v"`, key, entity, knownKey)
		} else {
			return value, fmt.Errorf(`unknown flag "%v" for "%v"`, key, entity)
		}

	}

	return value, nil
}
