package cmd

import (
	"reflect"

	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/clickwrap"
)

var (
	Clickwraps = &cobra.Command{}
)

func ClickwrapsInit() {
	Clickwraps = &cobra.Command{
		Use:  "clickwraps [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command clickwraps\n\t%v", args[0])
		},
	}
	var fieldsList string
	var formatList string
	usePagerList := true
	paramsClickwrapList := files_sdk.ClickwrapListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List Clickwraps",
		Long:  `List Clickwraps`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsClickwrapList
			params.MaxPages = MaxPagesList

			client := clickwrap.Client{Config: *config}
			it, err := client.List(ctx, params)
			it.OnPageError = func(err error) (*[]interface{}, error) {
				overriddenValues, newErr := lib.ErrorWithOriginalResponse(err, formatList, config.Logger())
				values, ok := overriddenValues.([]interface{})
				if ok {
					return &values, newErr
				} else {
					return &[]interface{}{}, newErr
				}
			}
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(ctx, it, formatList, fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
		},
	}

	cmdList.Flags().StringVar(&paramsClickwrapList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsClickwrapList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	Clickwraps.AddCommand(cmdList)
	var fieldsFind string
	var formatFind string
	usePagerFind := true
	paramsClickwrapFind := files_sdk.ClickwrapFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Clickwrap`,
		Long:  `Show Clickwrap`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := clickwrap.Client{Config: *config}

			var clickwrap interface{}
			var err error
			clickwrap, err = client.Find(ctx, paramsClickwrapFind)
			lib.HandleResponse(ctx, clickwrap, err, formatFind, fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdFind.Flags().Int64Var(&paramsClickwrapFind.Id, "id", 0, "Clickwrap ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	Clickwraps.AddCommand(cmdFind)
	var fieldsCreate string
	var formatCreate string
	usePagerCreate := true
	paramsClickwrapCreate := files_sdk.ClickwrapCreateParams{}
	ClickwrapCreateUseWithBundles := ""
	ClickwrapCreateUseWithInboxes := ""
	ClickwrapCreateUseWithUsers := ""

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Clickwrap`,
		Long:  `Create Clickwrap`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := clickwrap.Client{Config: *config}

			var clickwrap interface{}
			var err error
			paramsClickwrapCreate.UseWithBundles = paramsClickwrapCreate.UseWithBundles.Enum()[ClickwrapCreateUseWithBundles]
			paramsClickwrapCreate.UseWithInboxes = paramsClickwrapCreate.UseWithInboxes.Enum()[ClickwrapCreateUseWithInboxes]
			paramsClickwrapCreate.UseWithUsers = paramsClickwrapCreate.UseWithUsers.Enum()[ClickwrapCreateUseWithUsers]
			clickwrap, err = client.Create(ctx, paramsClickwrapCreate)
			lib.HandleResponse(ctx, clickwrap, err, formatCreate, fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdCreate.Flags().StringVar(&paramsClickwrapCreate.Name, "name", "", "Name of the Clickwrap agreement (used when selecting from multiple Clickwrap agreements.)")
	cmdCreate.Flags().StringVar(&paramsClickwrapCreate.Body, "body", "", "Body text of Clickwrap (supports Markdown formatting).")
	cmdCreate.Flags().StringVar(&ClickwrapCreateUseWithBundles, "use-with-bundles", "", fmt.Sprintf("Use this Clickwrap for Bundles? %v", reflect.ValueOf(paramsClickwrapCreate.UseWithBundles.Enum()).MapKeys()))
	cmdCreate.Flags().StringVar(&ClickwrapCreateUseWithInboxes, "use-with-inboxes", "", fmt.Sprintf("Use this Clickwrap for Inboxes? %v", reflect.ValueOf(paramsClickwrapCreate.UseWithInboxes.Enum()).MapKeys()))
	cmdCreate.Flags().StringVar(&ClickwrapCreateUseWithUsers, "use-with-users", "", fmt.Sprintf("Use this Clickwrap for User Registrations?  Note: This only applies to User Registrations where the User is invited to your Files.com site using an E-Mail invitation process where they then set their own password. %v", reflect.ValueOf(paramsClickwrapCreate.UseWithUsers.Enum()).MapKeys()))

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	Clickwraps.AddCommand(cmdCreate)
	var fieldsUpdate string
	var formatUpdate string
	usePagerUpdate := true
	paramsClickwrapUpdate := files_sdk.ClickwrapUpdateParams{}
	ClickwrapUpdateUseWithBundles := ""
	ClickwrapUpdateUseWithInboxes := ""
	ClickwrapUpdateUseWithUsers := ""

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Clickwrap`,
		Long:  `Update Clickwrap`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := clickwrap.Client{Config: *config}

			var clickwrap interface{}
			var err error
			paramsClickwrapUpdate.UseWithBundles = paramsClickwrapUpdate.UseWithBundles.Enum()[ClickwrapUpdateUseWithBundles]
			paramsClickwrapUpdate.UseWithInboxes = paramsClickwrapUpdate.UseWithInboxes.Enum()[ClickwrapUpdateUseWithInboxes]
			paramsClickwrapUpdate.UseWithUsers = paramsClickwrapUpdate.UseWithUsers.Enum()[ClickwrapUpdateUseWithUsers]
			clickwrap, err = client.Update(ctx, paramsClickwrapUpdate)
			lib.HandleResponse(ctx, clickwrap, err, formatUpdate, fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsClickwrapUpdate.Id, "id", 0, "Clickwrap ID.")
	cmdUpdate.Flags().StringVar(&paramsClickwrapUpdate.Name, "name", "", "Name of the Clickwrap agreement (used when selecting from multiple Clickwrap agreements.)")
	cmdUpdate.Flags().StringVar(&paramsClickwrapUpdate.Body, "body", "", "Body text of Clickwrap (supports Markdown formatting).")
	cmdUpdate.Flags().StringVar(&ClickwrapUpdateUseWithBundles, "use-with-bundles", "", fmt.Sprintf("Use this Clickwrap for Bundles? %v", reflect.ValueOf(paramsClickwrapUpdate.UseWithBundles.Enum()).MapKeys()))
	cmdUpdate.Flags().StringVar(&ClickwrapUpdateUseWithInboxes, "use-with-inboxes", "", fmt.Sprintf("Use this Clickwrap for Inboxes? %v", reflect.ValueOf(paramsClickwrapUpdate.UseWithInboxes.Enum()).MapKeys()))
	cmdUpdate.Flags().StringVar(&ClickwrapUpdateUseWithUsers, "use-with-users", "", fmt.Sprintf("Use this Clickwrap for User Registrations?  Note: This only applies to User Registrations where the User is invited to your Files.com site using an E-Mail invitation process where they then set their own password. %v", reflect.ValueOf(paramsClickwrapUpdate.UseWithUsers.Enum()).MapKeys()))

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	cmdUpdate.Flags().StringVarP(&formatUpdate, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	Clickwraps.AddCommand(cmdUpdate)
	var fieldsDelete string
	var formatDelete string
	usePagerDelete := true
	paramsClickwrapDelete := files_sdk.ClickwrapDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Clickwrap`,
		Long:  `Delete Clickwrap`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := clickwrap.Client{Config: *config}

			var err error
			err = client.Delete(ctx, paramsClickwrapDelete)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
		},
	}
	cmdDelete.Flags().Int64Var(&paramsClickwrapDelete.Id, "id", 0, "Clickwrap ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	Clickwraps.AddCommand(cmdDelete)
}
