package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	"fmt"
	"os"

	public_key "github.com/Files-com/files-sdk-go/publickey"
)

var (
	PublicKeys = &cobra.Command{
		Use:  "public-keys [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func PublicKeysInit() {
	var fieldsList string
	paramsPublicKeyList := files_sdk.PublicKeyListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsPublicKeyList
			params.MaxPages = MaxPagesList
			it := public_key.List(params)

			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().IntVarP(&paramsPublicKeyList.Page, "page", "p", 0, "Current page number.")
	cmdList.Flags().IntVarP(&paramsPublicKeyList.PerPage, "per-page", "e", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVarP(&paramsPublicKeyList.Action, "action", "a", "", "Deprecated: If set to `count` returns a count of matching records rather than the records themselves.")
	cmdList.Flags().StringVarP(&paramsPublicKeyList.Cursor, "cursor", "c", "", "Send cursor to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	PublicKeys.AddCommand(cmdList)
	var fieldsFind string
	paramsPublicKeyFind := files_sdk.PublicKeyFindParams{}
	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := public_key.Find(paramsPublicKeyFind)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsFind)
		},
	}

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	PublicKeys.AddCommand(cmdFind)
	var fieldsCreate string
	paramsPublicKeyCreate := files_sdk.PublicKeyCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := public_key.Create(paramsPublicKeyCreate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsCreate)
		},
	}
	cmdCreate.Flags().StringVarP(&paramsPublicKeyCreate.Title, "title", "t", "", "Internal reference for key.")
	cmdCreate.Flags().StringVarP(&paramsPublicKeyCreate.PublicKey, "public-key", "p", "", "Actual contents of SSH key.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	PublicKeys.AddCommand(cmdCreate)
	var fieldsUpdate string
	paramsPublicKeyUpdate := files_sdk.PublicKeyUpdateParams{}
	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := public_key.Update(paramsPublicKeyUpdate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsUpdate)
		},
	}
	cmdUpdate.Flags().StringVarP(&paramsPublicKeyUpdate.Title, "title", "t", "", "Internal reference for key.")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	PublicKeys.AddCommand(cmdUpdate)
	var fieldsDelete string
	paramsPublicKeyDelete := files_sdk.PublicKeyDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := public_key.Delete(paramsPublicKeyDelete)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsDelete)
		},
	}

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	PublicKeys.AddCommand(cmdDelete)
}
