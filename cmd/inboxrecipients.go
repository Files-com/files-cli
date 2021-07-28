package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	flib "github.com/Files-com/files-sdk-go/lib"

	"fmt"

	inbox_recipient "github.com/Files-com/files-sdk-go/inboxrecipient"
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
				lib.ClientError(ctx, err)
			}
			err = lib.JsonMarshalIter(it, fieldsList)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdList.Flags().Int64VarP(&paramsInboxRecipientList.UserId, "user-id", "u", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVarP(&paramsInboxRecipientList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().Int64VarP(&paramsInboxRecipientList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().Int64VarP(&paramsInboxRecipientList.InboxId, "inbox-id", "i", 0, "List recipients for the inbox with this ID.")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	InboxRecipients.AddCommand(cmdList)
	var fieldsCreate string
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
				lib.ClientError(ctx, err)
			}

			err = lib.JsonMarshal(result, fieldsCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdCreate.Flags().Int64VarP(&paramsInboxRecipientCreate.UserId, "user-id", "u", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().Int64VarP(&paramsInboxRecipientCreate.InboxId, "inbox-id", "i", 0, "Inbox to share.")
	cmdCreate.Flags().StringVarP(&paramsInboxRecipientCreate.Recipient, "recipient", "r", "", "Email address to share this inbox with.")
	cmdCreate.Flags().StringVarP(&paramsInboxRecipientCreate.Name, "name", "n", "", "Name of recipient.")
	cmdCreate.Flags().StringVarP(&paramsInboxRecipientCreate.Company, "company", "c", "", "Company of recipient.")
	cmdCreate.Flags().StringVarP(&paramsInboxRecipientCreate.Note, "note", "o", "", "Note to include in email.")
	cmdCreate.Flags().BoolVarP(&createShareAfterCreate, "share-after-create", "s", createShareAfterCreate, "Set to true to share the link with the recipient upon creation.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	InboxRecipients.AddCommand(cmdCreate)
}
