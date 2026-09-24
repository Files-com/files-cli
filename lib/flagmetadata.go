package lib

import (
	"maps"
	"slices"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Flag annotations read by command discovery. Cobra ignores annotation keys
// it does not define, so these never change flag parsing or validation.
const (
	flagAnnotationAPIRequired = "files-cli-api-required"
	flagAnnotationEnum        = "files-cli-enum"
)

// FlagRequired reports whether the command marked the flag with
// MarkFlagRequired, so Cobra rejects the command when it is missing.
func FlagRequired(flag *pflag.Flag) bool {
	required := flag.Annotations[cobra.BashCompOneRequiredFlag]
	return len(required) > 0 && required[0] == "true"
}

// SetFlagAPIRequired records that the Files.com API requires the parameter
// behind this flag. It is informational only; the CLI does not enforce it.
func SetFlagAPIRequired(flags *pflag.FlagSet, name string) {
	_ = flags.SetAnnotation(name, flagAnnotationAPIRequired, []string{"true"})
}

// FlagAPIRequired reports whether SetFlagAPIRequired marked the flag.
func FlagAPIRequired(flag *pflag.Flag) bool {
	return len(flag.Annotations[flagAnnotationAPIRequired]) > 0
}

// SetFlagEnum records the values an enum flag accepts.
func SetFlagEnum[T any](flags *pflag.FlagSet, name string, values map[string]T) {
	_ = flags.SetAnnotation(name, flagAnnotationEnum, slices.Sorted(maps.Keys(values)))
}

// FlagEnum returns the values recorded by SetFlagEnum, sorted.
func FlagEnum(flag *pflag.Flag) []string {
	return flag.Annotations[flagAnnotationEnum]
}

// FlagTypes returns the pflag value type of flag and, when SetFlagDisplayType
// replaced the type shown in help, that display type as the value syntax.
func FlagTypes(flag *pflag.Flag) (valueType string, syntax string) {
	if display, ok := flag.Value.(flagDisplayTypeValue); ok {
		return display.Value.Type(), display.displayType
	}
	return flag.Value.Type(), ""
}
