package lib

import (
	"github.com/Files-com/files-cli/lib/clierr"
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
			return value, clierr.Errorf(clierr.ErrorCodeFatal,
				`unknown flag "%v" for "%v"

Did you mean this?
	"%v"`, key, entity, knownKey)
		} else {
			return value, clierr.Errorf(clierr.ErrorCodeFatal, `unknown flag "%v" for "%v"`, key, entity)
		}

	}

	return value, nil
}
