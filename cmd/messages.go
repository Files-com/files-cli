package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/message"
)

var (
	Messages = &cobra.Command{}
)

func MessagesInit() {
	Messages = &cobra.Command{
		Use:  "messages [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command messages\n\t%v", args[0])
		},
	}
	var fieldsList string
	var formatList string
	paramsMessageList := files_sdk.MessageListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsMessageList
			params.MaxPages = MaxPagesList

			client := message.Client{Config: *config}
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
	cmdList.Flags().Int64VarP(&paramsMessageList.UserId, "user-id", "u", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVarP(&paramsMessageList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().Int64VarP(&paramsMessageList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().Int64VarP(&paramsMessageList.ProjectId, "project-id", "r", 0, "Project for which to return messages.")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Messages.AddCommand(cmdList)
	var fieldsFind string
	var formatFind string
	paramsMessageFind := files_sdk.MessageFindParams{}

	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := message.Client{Config: *config}

			result, err := client.Find(ctx, paramsMessageFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatFind, fieldsFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdFind.Flags().Int64VarP(&paramsMessageFind.Id, "id", "i", 0, "Message ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Messages.AddCommand(cmdFind)
	var fieldsCreate string
	var formatCreate string
	paramsMessageCreate := files_sdk.MessageCreateParams{}

	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := message.Client{Config: *config}

			result, err := client.Create(ctx, paramsMessageCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatCreate, fieldsCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdCreate.Flags().Int64VarP(&paramsMessageCreate.UserId, "user-id", "u", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().Int64VarP(&paramsMessageCreate.ProjectId, "project-id", "p", 0, "Project to which the message should be attached.")
	cmdCreate.Flags().StringVarP(&paramsMessageCreate.Subject, "subject", "s", "", "Message subject.")
	cmdCreate.Flags().StringVarP(&paramsMessageCreate.Body, "body", "b", "", "Message body.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Messages.AddCommand(cmdCreate)
	var fieldsUpdate string
	var formatUpdate string
	paramsMessageUpdate := files_sdk.MessageUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := message.Client{Config: *config}

			result, err := client.Update(ctx, paramsMessageUpdate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatUpdate, fieldsUpdate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdUpdate.Flags().Int64VarP(&paramsMessageUpdate.Id, "id", "i", 0, "Message ID.")
	cmdUpdate.Flags().Int64VarP(&paramsMessageUpdate.ProjectId, "project-id", "p", 0, "Project to which the message should be attached.")
	cmdUpdate.Flags().StringVarP(&paramsMessageUpdate.Subject, "subject", "s", "", "Message subject.")
	cmdUpdate.Flags().StringVarP(&paramsMessageUpdate.Body, "body", "b", "", "Message body.")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	cmdUpdate.Flags().StringVarP(&formatUpdate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Messages.AddCommand(cmdUpdate)
	var fieldsDelete string
	var formatDelete string
	paramsMessageDelete := files_sdk.MessageDeleteParams{}

	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := message.Client{Config: *config}

			result, err := client.Delete(ctx, paramsMessageDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatDelete, fieldsDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsMessageDelete.Id, "id", "i", 0, "Message ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Messages.AddCommand(cmdDelete)
}
