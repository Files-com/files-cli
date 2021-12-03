package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	"github.com/Files-com/files-sdk-go/v2/behavior"
	flib "github.com/Files-com/files-sdk-go/v2/lib"
)

var (
	Behaviors = &cobra.Command{}
)

func BehaviorsInit() {
	Behaviors = &cobra.Command{
		Use:  "behaviors [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command behaviors\n\t%v", args[0])
		},
	}
	var fieldsList string
	var formatList string
	paramsBehaviorList := files_sdk.BehaviorListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsBehaviorList
			params.MaxPages = MaxPagesList

			client := behavior.Client{Config: *config}
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

	cmdList.Flags().StringVar(&paramsBehaviorList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsBehaviorList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVar(&paramsBehaviorList.Behavior, "behavior", "", "If set, only shows folder behaviors matching this behavior type.")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright - (tables not supported for `list-for --recursive`)")
	Behaviors.AddCommand(cmdList)
	var fieldsFind string
	var formatFind string
	paramsBehaviorFind := files_sdk.BehaviorFindParams{}

	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := behavior.Client{Config: *config}

			result, err := client.Find(ctx, paramsBehaviorFind)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(result, formatFind, fieldsFind, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}
	cmdFind.Flags().Int64Var(&paramsBehaviorFind.Id, "id", 0, "Behavior ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Behaviors.AddCommand(cmdFind)
	var fieldsListFor string
	var formatListFor string
	paramsBehaviorListFor := files_sdk.BehaviorListForParams{}
	var MaxPagesListFor int64

	cmdListFor := &cobra.Command{
		Use:   "list-for [path]",
		Short: "list-for",
		Long:  `list-for`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsBehaviorListFor
			params.MaxPages = MaxPagesListFor
			if len(args) > 0 && args[0] != "" {
				params.Path = args[0]
			}

			client := behavior.Client{Config: *config}
			it, err := client.ListFor(ctx, params)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(it, formatListFor, fieldsListFor, listFilter, cmd.OutOrStdout())
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
		},
	}

	cmdListFor.Flags().StringVar(&paramsBehaviorListFor.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdListFor.Flags().Int64Var(&paramsBehaviorListFor.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdListFor.Flags().StringVar(&paramsBehaviorListFor.Path, "path", "", "Path to operate on.")
	cmdListFor.Flags().StringVar(&paramsBehaviorListFor.Recursive, "recursive", "", "Show behaviors above this path?")
	cmdListFor.Flags().StringVar(&paramsBehaviorListFor.Behavior, "behavior", "", "DEPRECATED: If set only shows folder behaviors matching this behavior type. Use `filter[behavior]` instead.")

	cmdListFor.Flags().Int64VarP(&MaxPagesListFor, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdListFor.Flags().StringVarP(&fieldsListFor, "fields", "", "", "comma separated list of field names to include in response")
	cmdListFor.Flags().StringVarP(&formatListFor, "format", "", "table", "json, csv, table, table-dark, table-bright - (tables not supported for `list-for --recursive`)")
	Behaviors.AddCommand(cmdListFor)
	var fieldsCreate string
	var formatCreate string
	paramsBehaviorCreate := files_sdk.BehaviorCreateParams{}

	cmdCreate := &cobra.Command{
		Use: "create [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := behavior.Client{Config: *config}

			if len(args) > 0 && args[0] != "" {
				paramsBehaviorCreate.Path = args[0]
			}

			result, err := client.Create(ctx, paramsBehaviorCreate)
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
	cmdCreate.Flags().StringVar(&paramsBehaviorCreate.Value, "value", "", "The value of the folder behavior.  Can be a integer, array, or hash depending on the type of folder behavior. See The Behavior Types section for example values for each type of behavior.")
	cmdCreate.Flags().StringVar(&paramsBehaviorCreate.Name, "name", "", "Name for this behavior.")
	cmdCreate.Flags().StringVar(&paramsBehaviorCreate.Description, "description", "", "Description for this behavior.")
	cmdCreate.Flags().StringVar(&paramsBehaviorCreate.Path, "path", "", "Folder behaviors path.")
	cmdCreate.Flags().StringVar(&paramsBehaviorCreate.Behavior, "behavior", "", "Behavior type.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Behaviors.AddCommand(cmdCreate)
	var fieldsWebhookTest string
	var formatWebhookTest string
	paramsBehaviorWebhookTest := files_sdk.BehaviorWebhookTestParams{}

	cmdWebhookTest := &cobra.Command{
		Use: "webhook-test",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := behavior.Client{Config: *config}

			result, err := client.WebhookTest(ctx, paramsBehaviorWebhookTest)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(result, formatWebhookTest, fieldsWebhookTest, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}
	cmdWebhookTest.Flags().StringVar(&paramsBehaviorWebhookTest.Url, "url", "", "URL for testing the webhook.")
	cmdWebhookTest.Flags().StringVar(&paramsBehaviorWebhookTest.Method, "method", "", "HTTP method(GET or POST).")
	cmdWebhookTest.Flags().StringVar(&paramsBehaviorWebhookTest.Encoding, "encoding", "", "HTTP encoding method.  Can be JSON, XML, or RAW (form data).")
	cmdWebhookTest.Flags().StringVar(&paramsBehaviorWebhookTest.Action, "action", "", "action for test body")

	cmdWebhookTest.Flags().StringVarP(&fieldsWebhookTest, "fields", "", "", "comma separated list of field names")
	cmdWebhookTest.Flags().StringVarP(&formatWebhookTest, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Behaviors.AddCommand(cmdWebhookTest)
	var fieldsUpdate string
	var formatUpdate string
	updateAttachmentDelete := false
	paramsBehaviorUpdate := files_sdk.BehaviorUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use: "update [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := behavior.Client{Config: *config}

			if updateAttachmentDelete {
				paramsBehaviorUpdate.AttachmentDelete = flib.Bool(true)
			}

			if len(args) > 0 && args[0] != "" {
				paramsBehaviorUpdate.Path = args[0]
			}

			result, err := client.Update(ctx, paramsBehaviorUpdate)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(result, formatUpdate, fieldsUpdate, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsBehaviorUpdate.Id, "id", 0, "Behavior ID.")
	cmdUpdate.Flags().StringVar(&paramsBehaviorUpdate.Value, "value", "", "The value of the folder behavior.  Can be a integer, array, or hash depending on the type of folder behavior. See The Behavior Types section for example values for each type of behavior.")
	cmdUpdate.Flags().StringVar(&paramsBehaviorUpdate.Name, "name", "", "Name for this behavior.")
	cmdUpdate.Flags().StringVar(&paramsBehaviorUpdate.Description, "description", "", "Description for this behavior.")
	cmdUpdate.Flags().StringVar(&paramsBehaviorUpdate.Behavior, "behavior", "", "Behavior type.")
	cmdUpdate.Flags().StringVar(&paramsBehaviorUpdate.Path, "path", "", "Folder behaviors path.")
	cmdUpdate.Flags().BoolVar(&updateAttachmentDelete, "attachment-delete", updateAttachmentDelete, "If true, will delete the file stored in attachment")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	cmdUpdate.Flags().StringVarP(&formatUpdate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Behaviors.AddCommand(cmdUpdate)
	var fieldsDelete string
	var formatDelete string
	paramsBehaviorDelete := files_sdk.BehaviorDeleteParams{}

	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := behavior.Client{Config: *config}

			result, err := client.Delete(ctx, paramsBehaviorDelete)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(result, formatDelete, fieldsDelete, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}
	cmdDelete.Flags().Int64Var(&paramsBehaviorDelete.Id, "id", 0, "Behavior ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Behaviors.AddCommand(cmdDelete)
}
