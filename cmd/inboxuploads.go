package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	inbox_upload "github.com/Files-com/files-sdk-go/inboxupload"
)

var (
	InboxUploads = &cobra.Command{
		Use:  "inbox-uploads [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func InboxUploadsInit() {
	var fieldsList string
	paramsInboxUploadList := files_sdk.InboxUploadListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsInboxUploadList
			params.MaxPages = MaxPagesList
			client := inbox_upload.Client{Config: files_sdk.GlobalConfig}
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
	cmdList.Flags().StringVarP(&paramsInboxUploadList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().IntVarP(&paramsInboxUploadList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().Int64VarP(&paramsInboxUploadList.InboxRegistrationId, "inbox-registration-id", "d", 0, "InboxRegistration ID")
	cmdList.Flags().Int64VarP(&paramsInboxUploadList.InboxId, "inbox-id", "i", 0, "Inbox ID")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	InboxUploads.AddCommand(cmdList)
}
