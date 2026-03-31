package lib

import "github.com/spf13/pflag"

type flagDisplayTypeValue struct {
	pflag.Value
	displayType string
}

func (v flagDisplayTypeValue) Type() string {
	return v.displayType
}

func SetFlagDisplayType(flags *pflag.FlagSet, name string, displayType string) {
	flag := flags.Lookup(name)
	if flag == nil {
		return
	}

	flag.Value = flagDisplayTypeValue{
		Value:       flag.Value,
		displayType: displayType,
	}
}
