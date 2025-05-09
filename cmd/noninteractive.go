package cmd

import (
	"strings"

	"github.com/Files-com/files-cli/lib/clierr"
	"github.com/spf13/cobra"
)

var (
	errNonInteractivePagingUnavailable = clierr.Errorf(clierr.ErrorCodeUsage, "--%s not supported in non-interactive mode", flagNameUsePager)
	errNonInteractiveReauthentication  = clierr.Errorf(clierr.ErrorCodeUsage, "--%s not supported in non-interactive mode", flagNameReauthentication)
)

// checkNonInteractiveMode checks the various combinations of the non-interactive
// flag with other flags and returns an error if the combination is not supported.
func checkNonInteractiveMode(cmd *cobra.Command) error {
	// non-interactive is not set, no need to check anything else.
	if !nonInteractive {
		return nil
	}

	// a valid session, with reauthentication will cause a prompt. If the non-interactive
	// flag is set, return an error that the combination is not supported.
	if Reauthentication && nonInteractive {
		return errNonInteractiveReauthentication
	}

	usePagerFlag := cmd.Flag(flagNameUsePager)
	if usePagerFlag != nil {
		usePager, err := cmd.Flags().GetBool(flagNameUsePager)
		if err != nil {
			return err
		}

		if nonInteractive && usePager && usePagerFlag.Changed {
			// this is the case where the user has set both flags
			// which does not make sense, so return an error.
			return errNonInteractivePagingUnavailable
		} else if nonInteractive {
			// otherwise, if the user explictly set the non-interactive flag
			// but not use-pager, set use-pager to false so that subcommands
			// avoid trying to page the output.
			usePagerFlag.Value.Set("false")
		}
	}

	// check if there's a format flag, and the content of the format flag indicates an interactive
	// table. If so, return an error that the combination is not supported.
	formatFlag := cmd.Flag(flagNameFormat)
	if formatFlag != nil {
		// e.g. --format 'table,interactive,vertical' or --format='table,interactive,vertical'
		format := formatFlag.Value.String()
		if formatFlag != nil && strings.Contains(format, flagNameInteractive) && nonInteractive {
			// [ and ] are added in the string representation of the flag so remove them to get
			// the actual format string.
			format = strings.ReplaceAll(format, "[", "")
			format = strings.ReplaceAll(format, "]", "")
			return clierr.Errorf(clierr.ErrorCodeUsage, "--%s='%s' is not supported in non-interactive mode", flagNameFormat, format)
		}
	}

	return nil
}
