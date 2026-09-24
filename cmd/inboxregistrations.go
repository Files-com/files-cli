package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	inbox_registration "github.com/Files-com/files-sdk-go/v3/inboxregistration"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(InboxRegistrations())
}

func InboxRegistrations() *cobra.Command {
	InboxRegistrations := &cobra.Command{
		Use:   "inbox-registrations [command]",
		Short: "An InboxRegistration is created when a user fills out the form to access the inbox.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command inbox-registrations\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsInboxRegistrationList := files_sdk.InboxRegistrationListParams{}
	var MaxPagesList int64
	var jsonEnvelopeList bool

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Inbox Registrations",
		Long:    `List Inbox Registrations`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsInboxRegistrationList
			params.MaxPages = MaxPagesList
			var envelopeStyle string
			if jsonEnvelopeList {
				var envelopeErr error
				if envelopeStyle, envelopeErr = lib.PrepareJSONEnvelope(cmd, Profile(cmd).Current().SetResourceFormat(cmd, formatList), &params.MaxPages); envelopeErr != nil {
					return envelopeErr
				}
			}

			client := inbox_registration.Client{Config: config}
			it, err := client.List(params, files_sdk.WithContext(ctx))
			it.OnPageError = func(err error) (*[]interface{}, error) {
				overriddenValues, newErr := lib.ErrorWithOriginalResponse(err, config.Logger)
				values, ok := overriddenValues.([]interface{})
				if ok {
					return &values, newErr
				} else {
					return &[]interface{}{}, newErr
				}
			}
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			if len(filterbyList) > 0 {
				listFilter = func(i interface{}) (interface{}, bool, error) {
					matchOk, err := lib.MatchFilter(filterbyList, i)
					return i, matchOk, err
				}
			}
			if jsonEnvelopeList {
				err = lib.JSONEnvelopeIter(it, fieldsList, listFilter, usePagerList, envelopeStyle, cmd.OutOrStdout())
			} else {
				err = lib.FormatIter(ctx, it, Profile(cmd).Current().SetResourceFormat(cmd, formatList), fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			}
			return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdList.Flags().StringToStringVar(&filterbyList, "filter-by", filterbyList, "Client-side wildcard filtering, for example field-name=*.jpg or field-name=?ello")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter-by", "field=pattern")

	cmdList.Flags().StringVar(&paramsInboxRegistrationList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsInboxRegistrationList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10000, 1,000 or less is recommended).")
	cmdList.Flags().Int64Var(&paramsInboxRegistrationList.FolderBehaviorId, "folder-behavior-id", 0, "ID of the associated Inbox. This is required if the user is not a site admin.")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	cmdList.Flags().BoolVar(&jsonEnvelopeList, "json-envelope", false, lib.JSONEnvelopeHelpText)
	InboxRegistrations.AddCommand(cmdList)
	return InboxRegistrations
}
