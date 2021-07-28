package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	"fmt"

	inbox_upload "github.com/Files-com/files-sdk-go/inboxupload"
)

var (
	InboxUploads = &cobra.Command{}
)

func InboxUploadsInit() {
	InboxUploads = &cobra.Command{
		Use:  "inbox-uploads [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command inbox-uploads\n\t%v", args[0])
		},
	}
	var fieldsList string
	paramsInboxUploadList := files_sdk.InboxUploadListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsInboxUploadList
			params.MaxPages = MaxPagesList

			client := inbox_upload.Client{Config: *config}
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
	cmdList.Flags().StringVarP(&paramsInboxUploadList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().Int64VarP(&paramsInboxUploadList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().Int64VarP(&paramsInboxUploadList.InboxRegistrationId, "inbox-registration-id", "d", 0, "InboxRegistration ID")
	cmdList.Flags().Int64VarP(&paramsInboxUploadList.InboxId, "inbox-id", "i", 0, "Inbox ID")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	InboxUploads.AddCommand(cmdList)
}
