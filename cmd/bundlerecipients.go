package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	flib "github.com/Files-com/files-sdk-go/lib"

	"fmt"

	bundle_recipient "github.com/Files-com/files-sdk-go/bundlerecipient"
)

var (
	BundleRecipients = &cobra.Command{}
)

func BundleRecipientsInit() {
	BundleRecipients = &cobra.Command{
		Use:  "bundle-recipients [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command bundle-recipients\n\t%v", args[0])
		},
	}
	var fieldsList string
	var formatList string
	paramsBundleRecipientList := files_sdk.BundleRecipientListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsBundleRecipientList
			params.MaxPages = MaxPagesList

			client := bundle_recipient.Client{Config: *config}
			it, err := client.List(ctx, params)
			if err != nil {
				lib.ClientError(ctx, err)
			}
			err = lib.FormatIter(it, formatList, fieldsList)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdList.Flags().Int64VarP(&paramsBundleRecipientList.UserId, "user-id", "u", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVarP(&paramsBundleRecipientList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().Int64VarP(&paramsBundleRecipientList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().Int64VarP(&paramsBundleRecipientList.BundleId, "bundle-id", "b", 0, "List recipients for the bundle with this ID.")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright")
	BundleRecipients.AddCommand(cmdList)
	var fieldsCreate string
	var formatCreate string
	createShareAfterCreate := false
	paramsBundleRecipientCreate := files_sdk.BundleRecipientCreateParams{}

	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := bundle_recipient.Client{Config: *config}

			if createShareAfterCreate {
				paramsBundleRecipientCreate.ShareAfterCreate = flib.Bool(true)
			}

			result, err := client.Create(ctx, paramsBundleRecipientCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatCreate, fieldsCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdCreate.Flags().Int64VarP(&paramsBundleRecipientCreate.UserId, "user-id", "u", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().Int64VarP(&paramsBundleRecipientCreate.BundleId, "bundle-id", "b", 0, "Bundle to share.")
	cmdCreate.Flags().StringVarP(&paramsBundleRecipientCreate.Recipient, "recipient", "r", "", "Email addresses to share this bundle with.")
	cmdCreate.Flags().StringVarP(&paramsBundleRecipientCreate.Name, "name", "n", "", "Name of recipient.")
	cmdCreate.Flags().StringVarP(&paramsBundleRecipientCreate.Company, "company", "c", "", "Company of recipient.")
	cmdCreate.Flags().StringVarP(&paramsBundleRecipientCreate.Note, "note", "o", "", "Note to include in email.")
	cmdCreate.Flags().BoolVarP(&createShareAfterCreate, "share-after-create", "s", createShareAfterCreate, "Set to true to share the link with the recipient upon creation.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	BundleRecipients.AddCommand(cmdCreate)
}
