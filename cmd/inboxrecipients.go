package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	inbox_recipient "github.com/Files-com/files-sdk-go/inboxrecipient"
)

var (
	InboxRecipients = &cobra.Command{
		Use:  "inbox-recipients [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func InboxRecipientsInit() {
	var fieldsList string
	paramsInboxRecipientList := files_sdk.InboxRecipientListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsInboxRecipientList
			params.MaxPages = MaxPagesList
			client := inbox_recipient.Client{Config: files_sdk.GlobalConfig}
			it, err := client.List(params)
			if err != nil {
				lib.ClientError(err)
			}
			err = lib.JsonMarshalIter(it, fieldsList)
			if err != nil {
				lib.ClientError(err)
			}
		},
	}
	cmdList.Flags().Int64VarP(&paramsInboxRecipientList.UserId, "user-id", "u", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVarP(&paramsInboxRecipientList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().IntVarP(&paramsInboxRecipientList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().Int64VarP(&paramsInboxRecipientList.InboxId, "inbox-id", "i", 0, "List recipients for the inbox with this ID.")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	InboxRecipients.AddCommand(cmdList)
	var fieldsCreate string
	paramsInboxRecipientCreate := files_sdk.InboxRecipientCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			client := inbox_recipient.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Create(paramsInboxRecipientCreate)
			if err != nil {
				lib.ClientError(err)
			}

			err = lib.JsonMarshal(result, fieldsCreate)
			if err != nil {
				lib.ClientError(err)
			}
		},
	}
	cmdCreate.Flags().Int64VarP(&paramsInboxRecipientCreate.UserId, "user-id", "u", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().Int64VarP(&paramsInboxRecipientCreate.InboxId, "inbox-id", "i", 0, "Inbox to share.")
	cmdCreate.Flags().StringVarP(&paramsInboxRecipientCreate.Recipient, "recipient", "r", "", "Email addresses to share this inbox with.")
	cmdCreate.Flags().StringVarP(&paramsInboxRecipientCreate.Name, "name", "n", "", "Name of recipient.")
	cmdCreate.Flags().StringVarP(&paramsInboxRecipientCreate.Company, "company", "c", "", "Company of recipient.")
	cmdCreate.Flags().StringVarP(&paramsInboxRecipientCreate.Note, "note", "o", "", "Note to include in email.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	InboxRecipients.AddCommand(cmdCreate)
}
