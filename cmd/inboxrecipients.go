package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	flib "github.com/Files-com/files-sdk-go/v2/lib"

	"fmt"

	inbox_recipient "github.com/Files-com/files-sdk-go/v2/inboxrecipient"
)

var (
	InboxRecipients = &cobra.Command{}
)

func InboxRecipientsInit() {
	InboxRecipients = &cobra.Command{
		Use:  "inbox-recipients [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command inbox-recipients\n\t%v", args[0])
		},
	}
	var fieldsList string
	var formatList string
	paramsInboxRecipientList := files_sdk.InboxRecipientListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsInboxRecipientList
			params.MaxPages = MaxPagesList

			client := inbox_recipient.Client{Config: *config}
			it, err := client.List(ctx, params)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(it, formatList, fieldsList, listFilter, cmd.OutOrStdout())
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
		},
	}

	cmdList.Flags().Int64Var(&paramsInboxRecipientList.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVar(&paramsInboxRecipientList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsInboxRecipientList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().Int64Var(&paramsInboxRecipientList.InboxId, "inbox-id", 0, "List recipients for the inbox with this ID.")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright - (tables not supported for `list-for --recursive`)")
	InboxRecipients.AddCommand(cmdList)
	var fieldsCreate string
	var formatCreate string
	createShareAfterCreate := false
	paramsInboxRecipientCreate := files_sdk.InboxRecipientCreateParams{}

	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := inbox_recipient.Client{Config: *config}

			if createShareAfterCreate {
				paramsInboxRecipientCreate.ShareAfterCreate = flib.Bool(true)
			}

			result, err := client.Create(ctx, paramsInboxRecipientCreate)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(result, formatCreate, fieldsCreate, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}
	cmdCreate.Flags().Int64Var(&paramsInboxRecipientCreate.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().Int64Var(&paramsInboxRecipientCreate.InboxId, "inbox-id", 0, "Inbox to share.")
	cmdCreate.Flags().StringVar(&paramsInboxRecipientCreate.Recipient, "recipient", "", "Email address to share this inbox with.")
	cmdCreate.Flags().StringVar(&paramsInboxRecipientCreate.Name, "name", "", "Name of recipient.")
	cmdCreate.Flags().StringVar(&paramsInboxRecipientCreate.Company, "company", "", "Company of recipient.")
	cmdCreate.Flags().StringVar(&paramsInboxRecipientCreate.Note, "note", "", "Note to include in email.")
	cmdCreate.Flags().BoolVar(&createShareAfterCreate, "share-after-create", createShareAfterCreate, "Set to true to share the link with the recipient upon creation.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	InboxRecipients.AddCommand(cmdCreate)
}
