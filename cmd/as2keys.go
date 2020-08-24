package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	"fmt"
	"os"

	as2_key "github.com/Files-com/files-sdk-go/as2key"
)

var (
	As2Keys = &cobra.Command{
		Use:  "as2-keys [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func As2KeysInit() {
	var fieldsList string
	paramsAs2KeyList := files_sdk.As2KeyListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsAs2KeyList
			params.MaxPages = MaxPagesList
			it := as2_key.List(params)

			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().IntVarP(&paramsAs2KeyList.Page, "page", "p", 0, "Current page number.")
	cmdList.Flags().IntVarP(&paramsAs2KeyList.PerPage, "per-page", "e", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVarP(&paramsAs2KeyList.Action, "action", "a", "", "Deprecated: If set to `count` returns a count of matching records rather than the records themselves.")
	cmdList.Flags().StringVarP(&paramsAs2KeyList.Cursor, "cursor", "c", "", "Send cursor to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	As2Keys.AddCommand(cmdList)
	var fieldsFind string
	paramsAs2KeyFind := files_sdk.As2KeyFindParams{}
	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := as2_key.Find(paramsAs2KeyFind)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsFind)
		},
	}

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	As2Keys.AddCommand(cmdFind)
	var fieldsCreate string
	paramsAs2KeyCreate := files_sdk.As2KeyCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := as2_key.Create(paramsAs2KeyCreate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsCreate)
		},
	}
	cmdCreate.Flags().StringVarP(&paramsAs2KeyCreate.As2PartnershipName, "as2-partnership-name", "a", "", "AS2 Partnership Name")
	cmdCreate.Flags().StringVarP(&paramsAs2KeyCreate.PublicKey, "public-key", "p", "", "Actual contents of Public key.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	As2Keys.AddCommand(cmdCreate)
	var fieldsUpdate string
	paramsAs2KeyUpdate := files_sdk.As2KeyUpdateParams{}
	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := as2_key.Update(paramsAs2KeyUpdate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsUpdate)
		},
	}
	cmdUpdate.Flags().StringVarP(&paramsAs2KeyUpdate.As2PartnershipName, "as2-partnership-name", "a", "", "AS2 Partnership Name")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	As2Keys.AddCommand(cmdUpdate)
	var fieldsDelete string
	paramsAs2KeyDelete := files_sdk.As2KeyDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := as2_key.Delete(paramsAs2KeyDelete)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsDelete)
		},
	}

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	As2Keys.AddCommand(cmdDelete)
}
