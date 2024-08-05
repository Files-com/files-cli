package cmd

import (
	"fmt"

	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/message"
	"github.com/spf13/cobra"
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
	filterbyList := make(map[string]string)
	paramsMessageList := files_sdk.MessageListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Messages",
		Long:    `List Messages`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsMessageList
			params.MaxPages = MaxPagesList

			client := message.Client{Config: config}
			it, err := client.List(params, files_sdk.WithContext(ctx))
			it.OnPageError = func(err error) (*[]interface{}, error) {
				overriddenValues, newErr := lib.ErrorWithOriginalResponse(err, config.Logger)
				values, ok := overriddenValues.([]interface{})
				if ok {
					return &values, newErr
				} else {
					return &[]interface{}{}, newErr
				}
			}
			if err != nil {
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			if len(filterbyList) > 0 {
				listFilter = func(i interface{}) (interface{}, bool, error) {
					matchOk, err := lib.MatchFilter(filterbyList, i)
					return i, matchOk, err
				}
			}
			err = lib.FormatIter(ctx, it, Profile(cmd).Current().SetResourceFormat(cmd, formatList), fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdList.Flags().StringToStringVar(&filterbyList, "filter-by", filterbyList, `Client side filtering: field-name=*.jpg,field-name=?ello`)

	cmdList.Flags().Int64Var(&paramsMessageList.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVar(&paramsMessageList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsMessageList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().Int64Var(&paramsMessageList.ProjectId, "project-id", 0, "Project for which to return messages.")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
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
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := message.Client{Config: config}

			var message interface{}
			var err error
			message, err = client.Find(paramsMessageFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), message, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsMessageFind.Id, "id", 0, "Message ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
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
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := message.Client{Config: config}

			var message interface{}
			var err error
			message, err = client.Create(paramsMessageCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), message, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().Int64Var(&paramsMessageCreate.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().Int64Var(&paramsMessageCreate.ProjectId, "project-id", 0, "Project to which the message should be attached.")
	cmdCreate.Flags().StringVar(&paramsMessageCreate.Subject, "subject", "", "Message subject.")
	cmdCreate.Flags().StringVar(&paramsMessageCreate.Body, "body", "", "Message body.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
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
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := message.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.MessageUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsMessageUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("project-id") {
				lib.FlagUpdate(cmd, "project_id", paramsMessageUpdate.ProjectId, mapParams)
			}
			if cmd.Flags().Changed("subject") {
				lib.FlagUpdate(cmd, "subject", paramsMessageUpdate.Subject, mapParams)
			}
			if cmd.Flags().Changed("body") {
				lib.FlagUpdate(cmd, "body", paramsMessageUpdate.Body, mapParams)
			}

			var message interface{}
			var err error
			message, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), message, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsMessageUpdate.Id, "id", 0, "Message ID.")
	cmdUpdate.Flags().Int64Var(&paramsMessageUpdate.ProjectId, "project-id", 0, "Project to which the message should be attached.")
	cmdUpdate.Flags().StringVar(&paramsMessageUpdate.Subject, "subject", "", "Message subject.")
	cmdUpdate.Flags().StringVar(&paramsMessageUpdate.Body, "body", "", "Message body.")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
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
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := message.Client{Config: config}

			var err error
			err = client.Delete(paramsMessageDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsMessageDelete.Id, "id", 0, "Message ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	Messages.AddCommand(cmdDelete)
	return Messages
}
