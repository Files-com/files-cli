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
	paramsClickwrapList := files_sdk.ClickwrapListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsClickwrapList
			params.MaxPages = MaxPagesList

			client := clickwrap.Client{Config: *config}
			it, err := client.List(ctx, params)
			if err != nil {
				lib.ClientError(ctx, err)
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(it, formatList, fieldsList, listFilter)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}

	cmdList.Flags().StringVar(&paramsClickwrapList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsClickwrapList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Clickwraps.AddCommand(cmdList)
	var fieldsFind string
	var formatFind string
	paramsClickwrapFind := files_sdk.ClickwrapFindParams{}

	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := clickwrap.Client{Config: *config}

			result, err := client.Find(ctx, paramsClickwrapFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatFind, fieldsFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdFind.Flags().Int64Var(&paramsClickwrapFind.Id, "id", 0, "Clickwrap ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Clickwraps.AddCommand(cmdFind)
	var fieldsCreate string
	var formatCreate string
	paramsClickwrapCreate := files_sdk.ClickwrapCreateParams{}
	ClickwrapCreateUseWithBundles := ""
	ClickwrapCreateUseWithInboxes := ""
	ClickwrapCreateUseWithUsers := ""

	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := clickwrap.Client{Config: *config}

			paramsClickwrapCreate.UseWithBundles = paramsClickwrapCreate.UseWithBundles.Enum()[ClickwrapCreateUseWithBundles]
			paramsClickwrapCreate.UseWithInboxes = paramsClickwrapCreate.UseWithInboxes.Enum()[ClickwrapCreateUseWithInboxes]
			paramsClickwrapCreate.UseWithUsers = paramsClickwrapCreate.UseWithUsers.Enum()[ClickwrapCreateUseWithUsers]

			result, err := client.Create(ctx, paramsClickwrapCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatCreate, fieldsCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdCreate.Flags().StringVar(&paramsClickwrapCreate.Name, "name", "", "Name of the Clickwrap agreement (used when selecting from multiple Clickwrap agreements.)")
	cmdCreate.Flags().StringVar(&paramsClickwrapCreate.Body, "body", "", "Body text of Clickwrap (supports Markdown formatting).")
	cmdCreate.Flags().StringVar(&ClickwrapCreateUseWithBundles, "use-with-bundles", "", fmt.Sprintf("Use this Clickwrap for Bundles? %v", reflect.ValueOf(paramsClickwrapCreate.UseWithBundles.Enum()).MapKeys()))
	cmdCreate.Flags().StringVar(&ClickwrapCreateUseWithInboxes, "use-with-inboxes", "", fmt.Sprintf("Use this Clickwrap for Inboxes? %v", reflect.ValueOf(paramsClickwrapCreate.UseWithInboxes.Enum()).MapKeys()))
	cmdCreate.Flags().StringVar(&ClickwrapCreateUseWithUsers, "use-with-users", "", fmt.Sprintf("Use this Clickwrap for User Registrations?  Note: This only applies to User Registrations where the User is invited to your Files.com site using an E-Mail invitation process where they then set their own password. %v", reflect.ValueOf(paramsClickwrapCreate.UseWithUsers.Enum()).MapKeys()))

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Clickwraps.AddCommand(cmdCreate)
	var fieldsUpdate string
	var formatUpdate string
	paramsClickwrapUpdate := files_sdk.ClickwrapUpdateParams{}
	ClickwrapUpdateUseWithBundles := ""
	ClickwrapUpdateUseWithInboxes := ""
	ClickwrapUpdateUseWithUsers := ""

	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := clickwrap.Client{Config: *config}

			paramsClickwrapUpdate.UseWithBundles = paramsClickwrapUpdate.UseWithBundles.Enum()[ClickwrapUpdateUseWithBundles]
			paramsClickwrapUpdate.UseWithInboxes = paramsClickwrapUpdate.UseWithInboxes.Enum()[ClickwrapUpdateUseWithInboxes]
			paramsClickwrapUpdate.UseWithUsers = paramsClickwrapUpdate.UseWithUsers.Enum()[ClickwrapUpdateUseWithUsers]

			result, err := client.Update(ctx, paramsClickwrapUpdate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatUpdate, fieldsUpdate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsClickwrapUpdate.Id, "id", 0, "Clickwrap ID.")
	cmdUpdate.Flags().StringVar(&paramsClickwrapUpdate.Name, "name", "", "Name of the Clickwrap agreement (used when selecting from multiple Clickwrap agreements.)")
	cmdUpdate.Flags().StringVar(&paramsClickwrapUpdate.Body, "body", "", "Body text of Clickwrap (supports Markdown formatting).")
	cmdUpdate.Flags().StringVar(&ClickwrapUpdateUseWithBundles, "use-with-bundles", "", fmt.Sprintf("Use this Clickwrap for Bundles? %v", reflect.ValueOf(paramsClickwrapUpdate.UseWithBundles.Enum()).MapKeys()))
	cmdUpdate.Flags().StringVar(&ClickwrapUpdateUseWithInboxes, "use-with-inboxes", "", fmt.Sprintf("Use this Clickwrap for Inboxes? %v", reflect.ValueOf(paramsClickwrapUpdate.UseWithInboxes.Enum()).MapKeys()))
	cmdUpdate.Flags().StringVar(&ClickwrapUpdateUseWithUsers, "use-with-users", "", fmt.Sprintf("Use this Clickwrap for User Registrations?  Note: This only applies to User Registrations where the User is invited to your Files.com site using an E-Mail invitation process where they then set their own password. %v", reflect.ValueOf(paramsClickwrapUpdate.UseWithUsers.Enum()).MapKeys()))

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	cmdUpdate.Flags().StringVarP(&formatUpdate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Clickwraps.AddCommand(cmdUpdate)
	var fieldsDelete string
	var formatDelete string
	paramsClickwrapDelete := files_sdk.ClickwrapDeleteParams{}

	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := clickwrap.Client{Config: *config}

			result, err := client.Delete(ctx, paramsClickwrapDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatDelete, fieldsDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdDelete.Flags().Int64Var(&paramsClickwrapDelete.Id, "id", 0, "Clickwrap ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Clickwraps.AddCommand(cmdDelete)
}
