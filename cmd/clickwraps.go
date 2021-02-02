package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	"fmt"
	"os"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/clickwrap"
)

var (
	Clickwraps = &cobra.Command{
		Use:  "clickwraps [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func ClickwrapsInit() {
	var fieldsList string
	paramsClickwrapList := files_sdk.ClickwrapListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsClickwrapList
			params.MaxPages = MaxPagesList
			client := clickwrap.Client{Config: files_sdk.GlobalConfig}
			it, err := client.List(params)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			err = lib.JsonMarshalIter(it, fieldsList)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	cmdList.Flags().StringVarP(&paramsClickwrapList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().IntVarP(&paramsClickwrapList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	Clickwraps.AddCommand(cmdList)
	var fieldsFind string
	paramsClickwrapFind := files_sdk.ClickwrapFindParams{}
	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			client := clickwrap.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Find(paramsClickwrapFind)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = lib.JsonMarshal(result, fieldsFind)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	cmdFind.Flags().Int64VarP(&paramsClickwrapFind.Id, "id", "i", 0, "Clickwrap ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	Clickwraps.AddCommand(cmdFind)
	var fieldsCreate string
	paramsClickwrapCreate := files_sdk.ClickwrapCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			client := clickwrap.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Create(paramsClickwrapCreate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = lib.JsonMarshal(result, fieldsCreate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	cmdCreate.Flags().StringVarP(&paramsClickwrapCreate.Name, "name", "n", "", "Name of the Clickwrap agreement (used when selecting from multiple Clickwrap agreements.)")
	cmdCreate.Flags().StringVarP(&paramsClickwrapCreate.Body, "body", "b", "", "Body text of Clickwrap (supports Markdown formatting).")
	cmdCreate.Flags().StringVarP(&paramsClickwrapCreate.UseWithBundles, "use-with-bundles", "u", "", "Use this Clickwrap for Bundles?")
	cmdCreate.Flags().StringVarP(&paramsClickwrapCreate.UseWithInboxes, "use-with-inboxes", "i", "", "Use this Clickwrap for Inboxes?")
	cmdCreate.Flags().StringVarP(&paramsClickwrapCreate.UseWithUsers, "use-with-users", "s", "", "Use this Clickwrap for User Registrations?  Note: This only applies to User Registrations where the User is invited to your Files.com site using an E-Mail invitation process where they then set their own password.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	Clickwraps.AddCommand(cmdCreate)
	var fieldsUpdate string
	paramsClickwrapUpdate := files_sdk.ClickwrapUpdateParams{}
	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			client := clickwrap.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Update(paramsClickwrapUpdate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = lib.JsonMarshal(result, fieldsUpdate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	cmdUpdate.Flags().Int64VarP(&paramsClickwrapUpdate.Id, "id", "i", 0, "Clickwrap ID.")
	cmdUpdate.Flags().StringVarP(&paramsClickwrapUpdate.Name, "name", "n", "", "Name of the Clickwrap agreement (used when selecting from multiple Clickwrap agreements.)")
	cmdUpdate.Flags().StringVarP(&paramsClickwrapUpdate.Body, "body", "b", "", "Body text of Clickwrap (supports Markdown formatting).")
	cmdUpdate.Flags().StringVarP(&paramsClickwrapUpdate.UseWithBundles, "use-with-bundles", "u", "", "Use this Clickwrap for Bundles?")
	cmdUpdate.Flags().StringVarP(&paramsClickwrapUpdate.UseWithInboxes, "use-with-inboxes", "o", "", "Use this Clickwrap for Inboxes?")
	cmdUpdate.Flags().StringVarP(&paramsClickwrapUpdate.UseWithUsers, "use-with-users", "s", "", "Use this Clickwrap for User Registrations?  Note: This only applies to User Registrations where the User is invited to your Files.com site using an E-Mail invitation process where they then set their own password.")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	Clickwraps.AddCommand(cmdUpdate)
	var fieldsDelete string
	paramsClickwrapDelete := files_sdk.ClickwrapDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			client := clickwrap.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Delete(paramsClickwrapDelete)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = lib.JsonMarshal(result, fieldsDelete)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsClickwrapDelete.Id, "id", "i", 0, "Clickwrap ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	Clickwraps.AddCommand(cmdDelete)
}
