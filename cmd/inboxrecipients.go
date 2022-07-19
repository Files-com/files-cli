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
	usePagerList := true
	paramsInboxRecipientList := files_sdk.InboxRecipientListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List Inbox Recipients",
		Long:  `List Inbox Recipients`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsInboxRecipientList
			params.MaxPages = MaxPagesList

			client := inbox_recipient.Client{Config: *config}
			it, err := client.List(ctx, params)
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
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(ctx, it, formatList, fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
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
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	InboxRecipients.AddCommand(cmdList)
	var fieldsCreate string
	var formatCreate string
	usePagerCreate := true
	createShareAfterCreate := false
	paramsInboxRecipientCreate := files_sdk.InboxRecipientCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Inbox Recipient`,
		Long:  `Create Inbox Recipient`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := inbox_recipient.Client{Config: *config}

			if createShareAfterCreate {
				paramsInboxRecipientCreate.ShareAfterCreate = flib.Bool(true)
			}

			var inboxRecipient interface{}
			var err error
			inboxRecipient, err = client.Create(ctx, paramsInboxRecipientCreate)
			lib.HandleResponse(ctx, inboxRecipient, err, formatCreate, fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
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
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	InboxRecipients.AddCommand(cmdCreate)
}
