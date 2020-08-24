package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	"fmt"
	"os"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/bundle"
)

var (
	Bundles = &cobra.Command{
		Use:  "bundles [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func BundlesInit() {
	var fieldsList string
	paramsBundleList := files_sdk.BundleListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsBundleList
			params.MaxPages = MaxPagesList
			it := bundle.List(params)

			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().IntVarP(&paramsBundleList.Page, "page", "p", 0, "Current page number.")
	cmdList.Flags().IntVarP(&paramsBundleList.PerPage, "per-page", "r", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVarP(&paramsBundleList.Action, "action", "a", "", "Deprecated: If set to `count` returns a count of matching records rather than the records themselves.")
	cmdList.Flags().StringVarP(&paramsBundleList.Cursor, "cursor", "c", "", "Send cursor to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	Bundles.AddCommand(cmdList)
	var fieldsFind string
	paramsBundleFind := files_sdk.BundleFindParams{}
	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := bundle.Find(paramsBundleFind)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsFind)
		},
	}

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	Bundles.AddCommand(cmdFind)
	var fieldsCreate string
	paramsBundleCreate := files_sdk.BundleCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := bundle.Create(paramsBundleCreate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsCreate)
		},
	}
	cmdCreate.Flags().StringVarP(&paramsBundleCreate.Password, "password", "p", "", "Password for this bundle.")
	lib.TimeVarP(cmdCreate.Flags(), &paramsBundleCreate.ExpiresAt, "expires-at", "e")
	cmdCreate.Flags().IntVarP(&paramsBundleCreate.MaxUses, "max-uses", "a", 0, "Maximum number of times bundle can be accessed")
	cmdCreate.Flags().StringVarP(&paramsBundleCreate.Description, "description", "d", "", "Public description")
	cmdCreate.Flags().StringVarP(&paramsBundleCreate.Note, "note", "n", "", "Bundle internal note")
	cmdCreate.Flags().StringVarP(&paramsBundleCreate.Code, "code", "o", "", "Bundle code.  This code forms the end part of the Public URL.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	Bundles.AddCommand(cmdCreate)
	var fieldsShare string
	paramsBundleShare := files_sdk.BundleShareParams{}
	cmdShare := &cobra.Command{
		Use: "share",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := bundle.Share(paramsBundleShare)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsShare)
		},
	}
	cmdShare.Flags().StringVarP(&paramsBundleShare.Note, "note", "n", "", "Note to include in email.")

	cmdShare.Flags().StringVarP(&fieldsShare, "fields", "", "", "comma separated list of field names")
	Bundles.AddCommand(cmdShare)
	var fieldsUpdate string
	paramsBundleUpdate := files_sdk.BundleUpdateParams{}
	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := bundle.Update(paramsBundleUpdate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsUpdate)
		},
	}
	cmdUpdate.Flags().StringVarP(&paramsBundleUpdate.Password, "password", "p", "", "Password for this bundle.")
	cmdUpdate.Flags().StringVarP(&paramsBundleUpdate.Code, "code", "o", "", "Bundle code.  This code forms the end part of the Public URL.")
	cmdUpdate.Flags().StringVarP(&paramsBundleUpdate.Description, "description", "d", "", "Public description")
	lib.TimeVarP(cmdUpdate.Flags(), &paramsBundleUpdate.ExpiresAt, "expires-at", "e")
	cmdUpdate.Flags().IntVarP(&paramsBundleUpdate.MaxUses, "max-uses", "a", 0, "Maximum number of times bundle can be accessed")
	cmdUpdate.Flags().StringVarP(&paramsBundleUpdate.Note, "note", "t", "", "Bundle internal note")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	Bundles.AddCommand(cmdUpdate)
	var fieldsDelete string
	paramsBundleDelete := files_sdk.BundleDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := bundle.Delete(paramsBundleDelete)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsDelete)
		},
	}

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	Bundles.AddCommand(cmdDelete)
}
