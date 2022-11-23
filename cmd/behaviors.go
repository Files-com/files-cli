package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	"github.com/Files-com/files-sdk-go/v2/behavior"
	flib "github.com/Files-com/files-sdk-go/v2/lib"
)

func init() {
	RootCmd.AddCommand(Behaviors())
}

func Behaviors() *cobra.Command {
	Behaviors := &cobra.Command{
		Use:  "behaviors [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command behaviors\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	paramsBehaviorList := files_sdk.BehaviorListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List Behaviors",
		Long:  `List Behaviors`,
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsBehaviorList
			params.MaxPages = MaxPagesList

			client := behavior.Client{Config: *config}
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
				lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(ctx, it, formatList, fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			if err != nil {
				lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}

	cmdList.Flags().StringVar(&paramsBehaviorList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsBehaviorList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVar(&paramsBehaviorList.Behavior, "behavior", "", "If set, only shows folder behaviors matching this behavior type.")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
        table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
        json-styles: {raw, pretty}
        `)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	Behaviors.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsBehaviorFind := files_sdk.BehaviorFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Behavior`,
		Long:  `Show Behavior`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := behavior.Client{Config: *config}

			var behavior interface{}
			var err error
			behavior, err = client.Find(ctx, paramsBehaviorFind)
			lib.HandleResponse(ctx, Profile(cmd), behavior, err, formatFind, fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdFind.Flags().Int64Var(&paramsBehaviorFind.Id, "id", 0, "Behavior ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	Behaviors.AddCommand(cmdFind)
	var fieldsListFor []string
	var formatListFor []string
	usePagerListFor := true
	paramsBehaviorListFor := files_sdk.BehaviorListForParams{}
	var MaxPagesListFor int64

	cmdListFor := &cobra.Command{
		Use:   "list-for [path]",
		Short: "List Behaviors by path",
		Long:  `List Behaviors by path`,
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsBehaviorListFor
			params.MaxPages = MaxPagesListFor
			if len(args) > 0 && args[0] != "" {
				params.Path = args[0]
			}

			client := behavior.Client{Config: *config}
			it, err := client.ListFor(ctx, params)
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
				lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(ctx, it, formatListFor, fieldsListFor, usePagerListFor, listFilter, cmd.OutOrStdout())
			if err != nil {
				lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}

	cmdListFor.Flags().StringVar(&paramsBehaviorListFor.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdListFor.Flags().Int64Var(&paramsBehaviorListFor.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdListFor.Flags().StringVar(&paramsBehaviorListFor.Path, "path", "", "Path to operate on.")
	cmdListFor.Flags().StringVar(&paramsBehaviorListFor.Recursive, "recursive", "", "Show behaviors above this path?")
	cmdListFor.Flags().StringVar(&paramsBehaviorListFor.Behavior, "behavior", "", "DEPRECATED: If set only shows folder behaviors matching this behavior type. Use `filter[behavior]` instead.")

	cmdListFor.Flags().Int64VarP(&MaxPagesListFor, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdListFor.Flags().StringSliceVar(&fieldsListFor, "fields", []string{}, "comma separated list of field names to include in response")
	cmdListFor.Flags().StringSliceVar(&formatListFor, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
        table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
        json-styles: {raw, pretty}
        `)
	cmdListFor.Flags().BoolVar(&usePagerListFor, "use-pager", usePagerListFor, "Use $PAGER (.ie less, more, etc)")
	Behaviors.AddCommand(cmdListFor)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	paramsBehaviorCreate := files_sdk.BehaviorCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create [path]",
		Short: `Create Behavior`,
		Long:  `Create Behavior`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := behavior.Client{Config: *config}

			if len(args) > 0 && args[0] != "" {
				paramsBehaviorCreate.Path = args[0]
			}
			var behavior interface{}
			var err error
			behavior, err = client.Create(ctx, paramsBehaviorCreate)
			lib.HandleResponse(ctx, Profile(cmd), behavior, err, formatCreate, fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdCreate.Flags().StringVar(&paramsBehaviorCreate.Value, "value", "", "The value of the folder behavior.  Can be a integer, array, or hash depending on the type of folder behavior. See The Behavior Types section for example values for each type of behavior.")
	cmdCreate.Flags().StringVar(&paramsBehaviorCreate.Name, "name", "", "Name for this behavior.")
	cmdCreate.Flags().StringVar(&paramsBehaviorCreate.Description, "description", "", "Description for this behavior.")
	cmdCreate.Flags().StringVar(&paramsBehaviorCreate.Path, "path", "", "Folder behaviors path.")
	cmdCreate.Flags().StringVar(&paramsBehaviorCreate.Behavior, "behavior", "", "Behavior type.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	Behaviors.AddCommand(cmdCreate)
	var fieldsWebhookTest []string
	var formatWebhookTest []string
	usePagerWebhookTest := true
	paramsBehaviorWebhookTest := files_sdk.BehaviorWebhookTestParams{}

	cmdWebhookTest := &cobra.Command{
		Use:   "webhook-test",
		Short: `Test webhook`,
		Long:  `Test webhook`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := behavior.Client{Config: *config}

			var err error
			err = client.WebhookTest(ctx, paramsBehaviorWebhookTest)
			if err != nil {
				lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdWebhookTest.Flags().StringVar(&paramsBehaviorWebhookTest.Url, "url", "", "URL for testing the webhook.")
	cmdWebhookTest.Flags().StringVar(&paramsBehaviorWebhookTest.Method, "method", "", "HTTP method(GET or POST).")
	cmdWebhookTest.Flags().StringVar(&paramsBehaviorWebhookTest.Encoding, "encoding", "", "HTTP encoding method.  Can be JSON, XML, or RAW (form data).")
	cmdWebhookTest.Flags().StringVar(&paramsBehaviorWebhookTest.Action, "action", "", "action for test body")

	cmdWebhookTest.Flags().StringSliceVar(&fieldsWebhookTest, "fields", []string{}, "comma separated list of field names")
	cmdWebhookTest.Flags().StringSliceVar(&formatWebhookTest, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdWebhookTest.Flags().BoolVar(&usePagerWebhookTest, "use-pager", usePagerWebhookTest, "Use $PAGER (.ie less, more, etc)")

	Behaviors.AddCommand(cmdWebhookTest)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	updateAttachmentDelete := true
	paramsBehaviorUpdate := files_sdk.BehaviorUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update [path]",
		Short: `Update Behavior`,
		Long:  `Update Behavior`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := behavior.Client{Config: *config}

			if cmd.Flags().Changed("attachment-delete") {
				paramsBehaviorUpdate.AttachmentDelete = flib.Bool(updateAttachmentDelete)
			}

			if len(args) > 0 && args[0] != "" {
				paramsBehaviorUpdate.Path = args[0]
			}
			var behavior interface{}
			var err error
			behavior, err = client.Update(ctx, paramsBehaviorUpdate)
			lib.HandleResponse(ctx, Profile(cmd), behavior, err, formatUpdate, fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsBehaviorUpdate.Id, "id", 0, "Behavior ID.")
	cmdUpdate.Flags().StringVar(&paramsBehaviorUpdate.Value, "value", "", "The value of the folder behavior.  Can be a integer, array, or hash depending on the type of folder behavior. See The Behavior Types section for example values for each type of behavior.")
	cmdUpdate.Flags().StringVar(&paramsBehaviorUpdate.Name, "name", "", "Name for this behavior.")
	cmdUpdate.Flags().StringVar(&paramsBehaviorUpdate.Description, "description", "", "Description for this behavior.")
	cmdUpdate.Flags().StringVar(&paramsBehaviorUpdate.Behavior, "behavior", "", "Behavior type.")
	cmdUpdate.Flags().StringVar(&paramsBehaviorUpdate.Path, "path", "", "Folder behaviors path.")
	cmdUpdate.Flags().BoolVar(&updateAttachmentDelete, "attachment-delete", updateAttachmentDelete, "If true, will delete the file stored in attachment")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	Behaviors.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsBehaviorDelete := files_sdk.BehaviorDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Behavior`,
		Long:  `Delete Behavior`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := behavior.Client{Config: *config}

			var err error
			err = client.Delete(ctx, paramsBehaviorDelete)
			if err != nil {
				lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsBehaviorDelete.Id, "id", 0, "Behavior ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	Behaviors.AddCommand(cmdDelete)
	return Behaviors
}
