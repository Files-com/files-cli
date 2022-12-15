package lib

import (
	"reflect"

	"github.com/samber/lo"

	"github.com/spf13/cobra"
)

func FlagUpdate[T comparable](cmd *cobra.Command, flag string, value T, m map[string]interface{}) {
	if lo.IsEmpty(value) {
		m[flag] = nil
	} else {
		m[flag] = value
	}
}

func FlagUpdateLen(cmd *cobra.Command, flag string, value interface{}, m map[string]interface{}) {
	if LenIsEmpty(value) {
		m[flag] = nil
	} else {
		m[flag] = value
	}
}

func LenIsEmpty(value interface{}) bool {
	switch reflect.TypeOf(value).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(value)
		return s.Len() == 0
	default:
		return false
	}
}
