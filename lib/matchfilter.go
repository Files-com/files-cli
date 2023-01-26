package lib

import (
	"fmt"

	"github.com/IGLOU-EU/go-wildcard"
)

func MatchFilter(filterBy map[string]string, i interface{}) (bool, error) {
	for field, pattern := range filterBy {
		m, selectedFields, err := OnlyFields([]string{field}, i)
		if err != nil {
			return false, err
		}
		value := selectedFields[0]
		if wildcard.Match(pattern, fmt.Sprintf("%v", m[value])) {
			return true, nil
		}
	}

	return false, nil
}
