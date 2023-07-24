package cmd

import (
	"fmt"

	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	bundle_recipient "github.com/Files-com/files-sdk-go/v2/bundlerecipient"
	flib "github.com/Files-com/files-sdk-go/v2/lib"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(BundleRecipients())
}

func BundleRecipients() *cobra.Command {
	BundleRecipients := &cobra.Command{
		Use:  "bundle-recipients [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command bundle-recipients\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsBundleRecipientList := files_sdk.BundleRecipientListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Bundle Recipients",
		Long:    `List Bundle Recipients`,
		Args:    cobra.MinimumNArgs(0),
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsBundleRecipientList
			params.MaxPages = MaxPagesList

			client := bundle_recipient.Client{Config: *config}
			it, err := client.List(params, files_sdk.WithContext(ctx))
			it.OnPageError = func(err error) (*[]interface{}, error) {
				overriddenValues, newErr := lib.ErrorWithOriginalResponse(err, config.Logger())
				values, ok := overriddenValues.([]interface{})
				if ok {
					return &values, newErr
				} else {
					return &[]interface{}{}, newErr
				}
			}
			if err != nil {
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			if len(filterbyList) > 0 {
				listFilter = func(i interface{}) (interface{}, bool, error) {
					matchOk, err := lib.MatchFilter(filterbyList, i)
					return i, matchOk, err
				}
			}
			err = lib.FormatIter(ctx, it, formatList, fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdList.Flags().StringToStringVar(&filterbyList, "filter-by", filterbyList, `Client side filtering: field-name=*.jpg,field-name=?ello`)

	cmdList.Flags().Int64Var(&paramsBundleRecipientList.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVar(&paramsBundleRecipientList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsBundleRecipientList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().Int64Var(&paramsBundleRecipientList.BundleId, "bundle-id", 0, "List recipients for the bundle with this ID.")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
json-styles: {raw, pretty}
        `)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	BundleRecipients.AddCommand(cmdList)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	createShareAfterCreate := true
	paramsBundleRecipientCreate := files_sdk.BundleRecipientCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Bundle Recipient`,
		Long:  `Create Bundle Recipient`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := bundle_recipient.Client{Config: *config}

			if cmd.Flags().Changed("share-after-create") {
				paramsBundleRecipientCreate.ShareAfterCreate = flib.Bool(createShareAfterCreate)
			}

			var bundleRecipient interface{}
			var err error
			bundleRecipient, err = client.Create(paramsBundleRecipientCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), bundleRecipient, err, formatCreate, fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdCreate.Flags().Int64Var(&paramsBundleRecipientCreate.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().Int64Var(&paramsBundleRecipientCreate.BundleId, "bundle-id", 0, "Bundle to share.")
	cmdCreate.Flags().StringVar(&paramsBundleRecipientCreate.Recipient, "recipient", "", "Email addresses to share this bundle with.")
	cmdCreate.Flags().StringVar(&paramsBundleRecipientCreate.Name, "name", "", "Name of recipient.")
	cmdCreate.Flags().StringVar(&paramsBundleRecipientCreate.Company, "company", "", "Company of recipient.")
	cmdCreate.Flags().StringVar(&paramsBundleRecipientCreate.Note, "note", "", "Note to include in email.")
	cmdCreate.Flags().BoolVar(&createShareAfterCreate, "share-after-create", createShareAfterCreate, "Set to true to share the link with the recipient upon creation.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
json-styles: {raw, pretty}`)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	BundleRecipients.AddCommand(cmdCreate)
	return BundleRecipients
}
