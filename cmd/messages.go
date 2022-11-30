package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/message"
)

func init() {
	RootCmd.AddCommand(Messages())
}

func Messages() *cobra.Command {
	Messages := &cobra.Command{
		Use:  "messages [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command messages\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	paramsMessageList := files_sdk.MessageListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List Messages",
		Long:  `List Messages`,
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsMessageList
			params.MaxPages = MaxPagesList

			client := message.Client{Config: *config}
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
				return lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(ctx, it, formatList, fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			if err != nil {
				return lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}

	cmdList.Flags().Int64Var(&paramsMessageList.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVar(&paramsMessageList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsMessageList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().Int64Var(&paramsMessageList.ProjectId, "project-id", 0, "Project for which to return messages.")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
        table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
        json-styles: {raw, pretty}
        `)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	Messages.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsMessageFind := files_sdk.MessageFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Message`,
		Long:  `Show Message`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := message.Client{Config: *config}

			var message interface{}
			var err error
			message, err = client.Find(ctx, paramsMessageFind)
			lib.HandleResponse(ctx, Profile(cmd), message, err, formatFind, fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdFind.Flags().Int64Var(&paramsMessageFind.Id, "id", 0, "Message ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	Messages.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	paramsMessageCreate := files_sdk.MessageCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Message`,
		Long:  `Create Message`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := message.Client{Config: *config}

			var message interface{}
			var err error
			message, err = client.Create(ctx, paramsMessageCreate)
			lib.HandleResponse(ctx, Profile(cmd), message, err, formatCreate, fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdCreate.Flags().Int64Var(&paramsMessageCreate.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().Int64Var(&paramsMessageCreate.ProjectId, "project-id", 0, "Project to which the message should be attached.")
	cmdCreate.Flags().StringVar(&paramsMessageCreate.Subject, "subject", "", "Message subject.")
	cmdCreate.Flags().StringVar(&paramsMessageCreate.Body, "body", "", "Message body.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	Messages.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	paramsMessageUpdate := files_sdk.MessageUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Message`,
		Long:  `Update Message`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := message.Client{Config: *config}

			var message interface{}
			var err error
			message, err = client.Update(ctx, paramsMessageUpdate)
			lib.HandleResponse(ctx, Profile(cmd), message, err, formatUpdate, fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsMessageUpdate.Id, "id", 0, "Message ID.")
	cmdUpdate.Flags().Int64Var(&paramsMessageUpdate.ProjectId, "project-id", 0, "Project to which the message should be attached.")
	cmdUpdate.Flags().StringVar(&paramsMessageUpdate.Subject, "subject", "", "Message subject.")
	cmdUpdate.Flags().StringVar(&paramsMessageUpdate.Body, "body", "", "Message body.")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	Messages.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsMessageDelete := files_sdk.MessageDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Message`,
		Long:  `Delete Message`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := message.Client{Config: *config}

			var err error
			err = client.Delete(ctx, paramsMessageDelete)
			if err != nil {
				return lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsMessageDelete.Id, "id", 0, "Message ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	Messages.AddCommand(cmdDelete)
	return Messages
}
