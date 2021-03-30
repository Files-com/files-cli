package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/message"
)

var (
	Messages = &cobra.Command{
		Use:  "messages [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func MessagesInit() {
	var fieldsList string
	paramsMessageList := files_sdk.MessageListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsMessageList
			params.MaxPages = MaxPagesList
			client := message.Client{Config: files_sdk.GlobalConfig}
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
	cmdList.Flags().Int64VarP(&paramsMessageList.UserId, "user-id", "u", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVarP(&paramsMessageList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().IntVarP(&paramsMessageList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().Int64VarP(&paramsMessageList.ProjectId, "project-id", "r", 0, "Project for which to return messages.")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	Messages.AddCommand(cmdList)
	var fieldsFind string
	paramsMessageFind := files_sdk.MessageFindParams{}
	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			client := message.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Find(paramsMessageFind)
			if err != nil {
				lib.ClientError(err)
			}

			err = lib.JsonMarshal(result, fieldsFind)
			if err != nil {
				lib.ClientError(err)
			}
		},
	}
	cmdFind.Flags().Int64VarP(&paramsMessageFind.Id, "id", "i", 0, "Message ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	Messages.AddCommand(cmdFind)
	var fieldsCreate string
	paramsMessageCreate := files_sdk.MessageCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			client := message.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Create(paramsMessageCreate)
			if err != nil {
				lib.ClientError(err)
			}

			err = lib.JsonMarshal(result, fieldsCreate)
			if err != nil {
				lib.ClientError(err)
			}
		},
	}
	cmdCreate.Flags().Int64VarP(&paramsMessageCreate.UserId, "user-id", "u", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().Int64VarP(&paramsMessageCreate.ProjectId, "project-id", "p", 0, "Project to which the message should be attached.")
	cmdCreate.Flags().StringVarP(&paramsMessageCreate.Subject, "subject", "s", "", "Message subject.")
	cmdCreate.Flags().StringVarP(&paramsMessageCreate.Body, "body", "b", "", "Message body.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	Messages.AddCommand(cmdCreate)
	var fieldsUpdate string
	paramsMessageUpdate := files_sdk.MessageUpdateParams{}
	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			client := message.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Update(paramsMessageUpdate)
			if err != nil {
				lib.ClientError(err)
			}

			err = lib.JsonMarshal(result, fieldsUpdate)
			if err != nil {
				lib.ClientError(err)
			}
		},
	}
	cmdUpdate.Flags().Int64VarP(&paramsMessageUpdate.Id, "id", "i", 0, "Message ID.")
	cmdUpdate.Flags().Int64VarP(&paramsMessageUpdate.ProjectId, "project-id", "p", 0, "Project to which the message should be attached.")
	cmdUpdate.Flags().StringVarP(&paramsMessageUpdate.Subject, "subject", "s", "", "Message subject.")
	cmdUpdate.Flags().StringVarP(&paramsMessageUpdate.Body, "body", "b", "", "Message body.")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	Messages.AddCommand(cmdUpdate)
	var fieldsDelete string
	paramsMessageDelete := files_sdk.MessageDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			client := message.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Delete(paramsMessageDelete)
			if err != nil {
				lib.ClientError(err)
			}

			err = lib.JsonMarshal(result, fieldsDelete)
			if err != nil {
				lib.ClientError(err)
			}
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsMessageDelete.Id, "id", "i", 0, "Message ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	Messages.AddCommand(cmdDelete)
}
