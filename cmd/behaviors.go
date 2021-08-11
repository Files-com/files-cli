package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/behavior"
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
				lib.ClientError(ctx, err)
			}
			err = lib.FormatIter(it, formatList, fieldsList)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdList.Flags().StringVarP(&paramsBehaviorList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().Int64VarP(&paramsBehaviorList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVarP(&paramsBehaviorList.Behavior, "behavior", "b", "", "If set, only shows folder behaviors matching this behavior type.")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-light")
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
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatFind, fieldsFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdFind.Flags().Int64VarP(&paramsBehaviorFind.Id, "id", "i", 0, "Behavior ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-light")
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
				lib.ClientError(ctx, err)
			}
			err = lib.FormatIter(it, formatListFor, fieldsListFor)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdListFor.Flags().StringVarP(&paramsBehaviorListFor.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdListFor.Flags().Int64VarP(&paramsBehaviorListFor.PerPage, "per-page", "r", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdListFor.Flags().StringVarP(&paramsBehaviorListFor.Path, "path", "p", "", "Path to operate on.")
	cmdListFor.Flags().StringVarP(&paramsBehaviorListFor.Recursive, "recursive", "u", "", "Show behaviors above this path?")
	cmdListFor.Flags().StringVarP(&paramsBehaviorListFor.Behavior, "behavior", "b", "", "DEPRECATED: If set only shows folder behaviors matching this behavior type. Use `filter[behavior]` instead.")

	cmdListFor.Flags().Int64VarP(&MaxPagesListFor, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdListFor.Flags().StringVarP(&fieldsListFor, "fields", "", "", "comma separated list of field names to include in response")
	cmdListFor.Flags().StringVarP(&formatListFor, "format", "", "table", "json, csv, table, table-dark, table-light")
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
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatCreate, fieldsCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdCreate.Flags().StringVarP(&paramsBehaviorCreate.Value, "value", "v", "", "The value of the folder behavior.  Can be a integer, array, or hash depending on the type of folder behavior. See The Behavior Types section for example values for each type of behavior.")
	cmdCreate.Flags().StringVarP(&paramsBehaviorCreate.Path, "path", "p", "", "Folder behaviors path.")
	cmdCreate.Flags().StringVarP(&paramsBehaviorCreate.Behavior, "behavior", "b", "", "Behavior type.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-light")
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
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatWebhookTest, fieldsWebhookTest)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdWebhookTest.Flags().StringVarP(&paramsBehaviorWebhookTest.Url, "url", "u", "", "URL for testing the webhook.")
	cmdWebhookTest.Flags().StringVarP(&paramsBehaviorWebhookTest.Method, "method", "t", "", "HTTP method(GET or POST).")
	cmdWebhookTest.Flags().StringVarP(&paramsBehaviorWebhookTest.Encoding, "encoding", "e", "", "HTTP encoding method.  Can be JSON, XML, or RAW (form data).")
	cmdWebhookTest.Flags().StringVarP(&paramsBehaviorWebhookTest.Action, "action", "a", "", "action for test body")

	cmdWebhookTest.Flags().StringVarP(&fieldsWebhookTest, "fields", "", "", "comma separated list of field names")
	cmdWebhookTest.Flags().StringVarP(&formatWebhookTest, "format", "", "table", "json, csv, table, table-dark, table-light")
	Behaviors.AddCommand(cmdWebhookTest)
	var fieldsUpdate string
	var formatUpdate string
	paramsBehaviorUpdate := files_sdk.BehaviorUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use: "update [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := behavior.Client{Config: *config}

			if len(args) > 0 && args[0] != "" {
				paramsBehaviorUpdate.Path = args[0]
			}

			result, err := client.Update(ctx, paramsBehaviorUpdate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatUpdate, fieldsUpdate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdUpdate.Flags().Int64VarP(&paramsBehaviorUpdate.Id, "id", "i", 0, "Behavior ID.")
	cmdUpdate.Flags().StringVarP(&paramsBehaviorUpdate.Value, "value", "v", "", "The value of the folder behavior.  Can be a integer, array, or hash depending on the type of folder behavior. See The Behavior Types section for example values for each type of behavior.")
	cmdUpdate.Flags().StringVarP(&paramsBehaviorUpdate.Behavior, "behavior", "b", "", "Behavior type.")
	cmdUpdate.Flags().StringVarP(&paramsBehaviorUpdate.Path, "path", "p", "", "Folder behaviors path.")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	cmdUpdate.Flags().StringVarP(&formatUpdate, "format", "", "table", "json, csv, table, table-dark, table-light")
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
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatDelete, fieldsDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsBehaviorDelete.Id, "id", "i", 0, "Behavior ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-light")
	Behaviors.AddCommand(cmdDelete)
}
