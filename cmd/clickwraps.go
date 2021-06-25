package cmd

import (
	"reflect"

	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/clickwrap"
)

var (
	Clickwraps = &cobra.Command{}
)

func ClickwrapsInit() {
	Clickwraps = &cobra.Command{
		Use:  "clickwraps [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
	var fieldsList string
	paramsClickwrapList := files_sdk.ClickwrapListParams{}
	var MaxPagesList int64
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			params := paramsClickwrapList
			params.MaxPages = MaxPagesList
			client := clickwrap.Client{Config: *ctx.GetConfig()}
			it, err := client.List(params)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
			err = lib.JsonMarshalIter(it, fieldsList)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdList.Flags().StringVarP(&paramsClickwrapList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().Int64VarP(&paramsClickwrapList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	Clickwraps.AddCommand(cmdList)
	var fieldsFind string
	paramsClickwrapFind := files_sdk.ClickwrapFindParams{}
	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := clickwrap.Client{Config: *ctx.GetConfig()}

			result, err := client.Find(paramsClickwrapFind)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsFind)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdFind.Flags().Int64VarP(&paramsClickwrapFind.Id, "id", "i", 0, "Clickwrap ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	Clickwraps.AddCommand(cmdFind)
	var fieldsCreate string
	paramsClickwrapCreate := files_sdk.ClickwrapCreateParams{}
	ClickwrapCreateUseWithBundles := ""
	ClickwrapCreateUseWithInboxes := ""
	ClickwrapCreateUseWithUsers := ""
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := clickwrap.Client{Config: *ctx.GetConfig()}

			paramsClickwrapCreate.UseWithBundles = paramsClickwrapCreate.UseWithBundles.Enum()[ClickwrapCreateUseWithBundles]
			paramsClickwrapCreate.UseWithInboxes = paramsClickwrapCreate.UseWithInboxes.Enum()[ClickwrapCreateUseWithInboxes]
			paramsClickwrapCreate.UseWithUsers = paramsClickwrapCreate.UseWithUsers.Enum()[ClickwrapCreateUseWithUsers]
			result, err := client.Create(paramsClickwrapCreate)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsCreate)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdCreate.Flags().StringVarP(&paramsClickwrapCreate.Name, "name", "n", "", "Name of the Clickwrap agreement (used when selecting from multiple Clickwrap agreements.)")
	cmdCreate.Flags().StringVarP(&paramsClickwrapCreate.Body, "body", "b", "", "Body text of Clickwrap (supports Markdown formatting).")
	cmdCreate.Flags().StringVarP(&ClickwrapCreateUseWithBundles, "use-with-bundles", "u", "", fmt.Sprintf("Use this Clickwrap for Bundles? %v", reflect.ValueOf(paramsClickwrapCreate.UseWithBundles.Enum()).MapKeys()))
	cmdCreate.Flags().StringVarP(&ClickwrapCreateUseWithInboxes, "use-with-inboxes", "i", "", fmt.Sprintf("Use this Clickwrap for Inboxes? %v", reflect.ValueOf(paramsClickwrapCreate.UseWithInboxes.Enum()).MapKeys()))
	cmdCreate.Flags().StringVarP(&ClickwrapCreateUseWithUsers, "use-with-users", "s", "", fmt.Sprintf("Use this Clickwrap for User Registrations?  Note: This only applies to User Registrations where the User is invited to your Files.com site using an E-Mail invitation process where they then set their own password. %v", reflect.ValueOf(paramsClickwrapCreate.UseWithUsers.Enum()).MapKeys()))

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	Clickwraps.AddCommand(cmdCreate)
	var fieldsUpdate string
	paramsClickwrapUpdate := files_sdk.ClickwrapUpdateParams{}
	ClickwrapUpdateUseWithBundles := ""
	ClickwrapUpdateUseWithInboxes := ""
	ClickwrapUpdateUseWithUsers := ""
	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := clickwrap.Client{Config: *ctx.GetConfig()}

			paramsClickwrapUpdate.UseWithBundles = paramsClickwrapUpdate.UseWithBundles.Enum()[ClickwrapUpdateUseWithBundles]
			paramsClickwrapUpdate.UseWithInboxes = paramsClickwrapUpdate.UseWithInboxes.Enum()[ClickwrapUpdateUseWithInboxes]
			paramsClickwrapUpdate.UseWithUsers = paramsClickwrapUpdate.UseWithUsers.Enum()[ClickwrapUpdateUseWithUsers]
			result, err := client.Update(paramsClickwrapUpdate)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsUpdate)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdUpdate.Flags().Int64VarP(&paramsClickwrapUpdate.Id, "id", "i", 0, "Clickwrap ID.")
	cmdUpdate.Flags().StringVarP(&paramsClickwrapUpdate.Name, "name", "n", "", "Name of the Clickwrap agreement (used when selecting from multiple Clickwrap agreements.)")
	cmdUpdate.Flags().StringVarP(&paramsClickwrapUpdate.Body, "body", "b", "", "Body text of Clickwrap (supports Markdown formatting).")
	cmdUpdate.Flags().StringVarP(&ClickwrapUpdateUseWithBundles, "use-with-bundles", "u", "", fmt.Sprintf("Use this Clickwrap for Bundles? %v", reflect.ValueOf(paramsClickwrapUpdate.UseWithBundles.Enum()).MapKeys()))
	cmdUpdate.Flags().StringVarP(&ClickwrapUpdateUseWithInboxes, "use-with-inboxes", "o", "", fmt.Sprintf("Use this Clickwrap for Inboxes? %v", reflect.ValueOf(paramsClickwrapUpdate.UseWithInboxes.Enum()).MapKeys()))
	cmdUpdate.Flags().StringVarP(&ClickwrapUpdateUseWithUsers, "use-with-users", "s", "", fmt.Sprintf("Use this Clickwrap for User Registrations?  Note: This only applies to User Registrations where the User is invited to your Files.com site using an E-Mail invitation process where they then set their own password. %v", reflect.ValueOf(paramsClickwrapUpdate.UseWithUsers.Enum()).MapKeys()))

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	Clickwraps.AddCommand(cmdUpdate)
	var fieldsDelete string
	paramsClickwrapDelete := files_sdk.ClickwrapDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := clickwrap.Client{Config: *ctx.GetConfig()}

			result, err := client.Delete(paramsClickwrapDelete)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsDelete)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsClickwrapDelete.Id, "id", "i", 0, "Clickwrap ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	Clickwraps.AddCommand(cmdDelete)
}
