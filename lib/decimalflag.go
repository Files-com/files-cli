package lib

import (
	"github.com/Files-com/files-cli/lib/clierr"
	flib "github.com/Files-com/files-sdk-go/v3/lib"
)

// ParseDecimalFlag returns the value of --flagName unchanged when it is a
// decimal number such as 1.5 or 2e-3, so it reaches the API without rounding.
func ParseDecimalFlag(flagName string, value string) (string, error) {
	if !flib.IsDecimalText(value) {
		return "", clierr.Errorf(clierr.ErrorCodeUsage, "invalid value for --%s: expected a decimal number such as 1.5 or 2e-3", flagName)
	}
	return value, nil
}
