package cmd

import (
	"reflect"

	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/clickwrap"
)

func init() {
	RootCmd.AddCommand(Clickwraps())
}

func Clickwraps() *cobra.Command {
	Clickwraps := &cobra.Command{
		Use:  "clickwraps [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command clickwraps\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	paramsClickwrapList := files_sdk.ClickwrapListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List Clickwraps",
		Long:  `List Clickwraps`,
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsClickwrapList
			params.MaxPages = MaxPagesList

			client := clickwrap.Client{Config: *config}
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

	cmdList.Flags().StringVar(&paramsClickwrapList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsClickwrapList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
        table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
        json-styles: {raw, pretty}
        `)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	Clickwraps.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsClickwrapFind := files_sdk.ClickwrapFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Clickwrap`,
		Long:  `Show Clickwrap`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := clickwrap.Client{Config: *config}

			var clickwrap interface{}
			var err error
			clickwrap, err = client.Find(ctx, paramsClickwrapFind)
			lib.HandleResponse(ctx, Profile(cmd), clickwrap, err, formatFind, fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdFind.Flags().Int64Var(&paramsClickwrapFind.Id, "id", 0, "Clickwrap ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	Clickwraps.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	paramsClickwrapCreate := files_sdk.ClickwrapCreateParams{}
	ClickwrapCreateUseWithBundles := ""
	ClickwrapCreateUseWithInboxes := ""
	ClickwrapCreateUseWithUsers := ""

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Clickwrap`,
		Long:  `Create Clickwrap`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := clickwrap.Client{Config: *config}

			var clickwrap interface{}
			var err error
			var ClickwrapCreateUseWithBundlesErr error
			paramsClickwrapCreate.UseWithBundles, ClickwrapCreateUseWithBundlesErr = lib.FetchKey("use-with-bundles", paramsClickwrapCreate.UseWithBundles.Enum(), ClickwrapCreateUseWithBundles)
			if ClickwrapCreateUseWithBundles != "" && ClickwrapCreateUseWithBundlesErr != nil {
				return ClickwrapCreateUseWithBundlesErr
			}
			var ClickwrapCreateUseWithInboxesErr error
			paramsClickwrapCreate.UseWithInboxes, ClickwrapCreateUseWithInboxesErr = lib.FetchKey("use-with-inboxes", paramsClickwrapCreate.UseWithInboxes.Enum(), ClickwrapCreateUseWithInboxes)
			if ClickwrapCreateUseWithInboxes != "" && ClickwrapCreateUseWithInboxesErr != nil {
				return ClickwrapCreateUseWithInboxesErr
			}
			var ClickwrapCreateUseWithUsersErr error
			paramsClickwrapCreate.UseWithUsers, ClickwrapCreateUseWithUsersErr = lib.FetchKey("use-with-users", paramsClickwrapCreate.UseWithUsers.Enum(), ClickwrapCreateUseWithUsers)
			if ClickwrapCreateUseWithUsers != "" && ClickwrapCreateUseWithUsersErr != nil {
				return ClickwrapCreateUseWithUsersErr
			}
			clickwrap, err = client.Create(ctx, paramsClickwrapCreate)
			lib.HandleResponse(ctx, Profile(cmd), clickwrap, err, formatCreate, fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdCreate.Flags().StringVar(&paramsClickwrapCreate.Name, "name", "", "Name of the Clickwrap agreement (used when selecting from multiple Clickwrap agreements.)")
	cmdCreate.Flags().StringVar(&paramsClickwrapCreate.Body, "body", "", "Body text of Clickwrap (supports Markdown formatting).")
	cmdCreate.Flags().StringVar(&ClickwrapCreateUseWithBundles, "use-with-bundles", "", fmt.Sprintf("Use this Clickwrap for Bundles? %v", reflect.ValueOf(paramsClickwrapCreate.UseWithBundles.Enum()).MapKeys()))
	cmdCreate.Flags().StringVar(&ClickwrapCreateUseWithInboxes, "use-with-inboxes", "", fmt.Sprintf("Use this Clickwrap for Inboxes? %v", reflect.ValueOf(paramsClickwrapCreate.UseWithInboxes.Enum()).MapKeys()))
	cmdCreate.Flags().StringVar(&ClickwrapCreateUseWithUsers, "use-with-users", "", fmt.Sprintf("Use this Clickwrap for User Registrations?  Note: This only applies to User Registrations where the User is invited to your Files.com site using an E-Mail invitation process where they then set their own password. %v", reflect.ValueOf(paramsClickwrapCreate.UseWithUsers.Enum()).MapKeys()))

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	Clickwraps.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	paramsClickwrapUpdate := files_sdk.ClickwrapUpdateParams{}
	ClickwrapUpdateUseWithBundles := ""
	ClickwrapUpdateUseWithInboxes := ""
	ClickwrapUpdateUseWithUsers := ""

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Clickwrap`,
		Long:  `Update Clickwrap`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := clickwrap.Client{Config: *config}

			var clickwrap interface{}
			var err error
			var ClickwrapUpdateUseWithBundlesErr error
			paramsClickwrapUpdate.UseWithBundles, ClickwrapUpdateUseWithBundlesErr = lib.FetchKey("use-with-bundles", paramsClickwrapUpdate.UseWithBundles.Enum(), ClickwrapUpdateUseWithBundles)
			if ClickwrapUpdateUseWithBundles != "" && ClickwrapUpdateUseWithBundlesErr != nil {
				return ClickwrapUpdateUseWithBundlesErr
			}
			var ClickwrapUpdateUseWithInboxesErr error
			paramsClickwrapUpdate.UseWithInboxes, ClickwrapUpdateUseWithInboxesErr = lib.FetchKey("use-with-inboxes", paramsClickwrapUpdate.UseWithInboxes.Enum(), ClickwrapUpdateUseWithInboxes)
			if ClickwrapUpdateUseWithInboxes != "" && ClickwrapUpdateUseWithInboxesErr != nil {
				return ClickwrapUpdateUseWithInboxesErr
			}
			var ClickwrapUpdateUseWithUsersErr error
			paramsClickwrapUpdate.UseWithUsers, ClickwrapUpdateUseWithUsersErr = lib.FetchKey("use-with-users", paramsClickwrapUpdate.UseWithUsers.Enum(), ClickwrapUpdateUseWithUsers)
			if ClickwrapUpdateUseWithUsers != "" && ClickwrapUpdateUseWithUsersErr != nil {
				return ClickwrapUpdateUseWithUsersErr
			}
			clickwrap, err = client.Update(ctx, paramsClickwrapUpdate)
			lib.HandleResponse(ctx, Profile(cmd), clickwrap, err, formatUpdate, fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsClickwrapUpdate.Id, "id", 0, "Clickwrap ID.")
	cmdUpdate.Flags().StringVar(&paramsClickwrapUpdate.Name, "name", "", "Name of the Clickwrap agreement (used when selecting from multiple Clickwrap agreements.)")
	cmdUpdate.Flags().StringVar(&paramsClickwrapUpdate.Body, "body", "", "Body text of Clickwrap (supports Markdown formatting).")
	cmdUpdate.Flags().StringVar(&ClickwrapUpdateUseWithBundles, "use-with-bundles", "", fmt.Sprintf("Use this Clickwrap for Bundles? %v", reflect.ValueOf(paramsClickwrapUpdate.UseWithBundles.Enum()).MapKeys()))
	cmdUpdate.Flags().StringVar(&ClickwrapUpdateUseWithInboxes, "use-with-inboxes", "", fmt.Sprintf("Use this Clickwrap for Inboxes? %v", reflect.ValueOf(paramsClickwrapUpdate.UseWithInboxes.Enum()).MapKeys()))
	cmdUpdate.Flags().StringVar(&ClickwrapUpdateUseWithUsers, "use-with-users", "", fmt.Sprintf("Use this Clickwrap for User Registrations?  Note: This only applies to User Registrations where the User is invited to your Files.com site using an E-Mail invitation process where they then set their own password. %v", reflect.ValueOf(paramsClickwrapUpdate.UseWithUsers.Enum()).MapKeys()))

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	Clickwraps.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsClickwrapDelete := files_sdk.ClickwrapDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Clickwrap`,
		Long:  `Delete Clickwrap`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := clickwrap.Client{Config: *config}

			var err error
			err = client.Delete(ctx, paramsClickwrapDelete)
			if err != nil {
				lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsClickwrapDelete.Id, "id", 0, "Clickwrap ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	Clickwraps.AddCommand(cmdDelete)
	return Clickwraps
}
