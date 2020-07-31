package cmd

import "github.com/spf13/cobra"
import (
	"fmt"
	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/as2key"
	"os"
)

var (
	_ = files_sdk.Config{}
	_ = as2_key.Client{}
	_ = lib.OnlyFields
	_ = fmt.Println
	_ = os.Exit
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
	cmdList.Flags().IntVarP(&paramsAs2KeyList.UserId, "user-id", "u", 0, "List As2 Keys")
	cmdList.Flags().IntVarP(&paramsAs2KeyList.Page, "page", "p", 0, "List As2 Keys")
	cmdList.Flags().IntVarP(&paramsAs2KeyList.PerPage, "per-page", "e", 0, "List As2 Keys")
	cmdList.Flags().StringVarP(&paramsAs2KeyList.Action, "action", "a", "", "List As2 Keys")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "f", "", "comma separated list of field names to include in response")
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
	cmdFind.Flags().IntVarP(&paramsAs2KeyFind.Id, "id", "i", 0, "Show As2 Key")
	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "f", "", "comma separated list of field names")
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
	cmdCreate.Flags().IntVarP(&paramsAs2KeyCreate.UserId, "user-id", "u", 0, "Create As2 Key")
	cmdCreate.Flags().StringVarP(&paramsAs2KeyCreate.As2PartnershipName, "as2-partnership-name", "a", "", "Create As2 Key")
	cmdCreate.Flags().StringVarP(&paramsAs2KeyCreate.PublicKey, "public-key", "p", "", "Create As2 Key")
	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "f", "", "comma separated list of field names")
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
	cmdUpdate.Flags().IntVarP(&paramsAs2KeyUpdate.Id, "id", "i", 0, "Update As2 Key")
	cmdUpdate.Flags().StringVarP(&paramsAs2KeyUpdate.As2PartnershipName, "as2-partnership-name", "a", "", "Update As2 Key")
	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "f", "", "comma separated list of field names")
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
	cmdDelete.Flags().IntVarP(&paramsAs2KeyDelete.Id, "id", "i", 0, "Delete As2 Key")
	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "f", "", "comma separated list of field names")
	As2Keys.AddCommand(cmdDelete)
}
