package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	"fmt"
	"os"

	bundle_download "github.com/Files-com/files-sdk-go/bundledownload"
)

var (
	BundleDownloads = &cobra.Command{
		Use:  "bundle-downloads [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func BundleDownloadsInit() {
	var fieldsList string
	paramsBundleDownloadList := files_sdk.BundleDownloadListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsBundleDownloadList
			params.MaxPages = MaxPagesList
			it, err := bundle_download.List(params)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().IntVarP(&paramsBundleDownloadList.Page, "page", "p", 0, "Current page number.")
	cmdList.Flags().IntVarP(&paramsBundleDownloadList.PerPage, "per-page", "e", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVarP(&paramsBundleDownloadList.Action, "action", "a", "", "Deprecated: If set to `count` returns a count of matching records rather than the records themselves.")
	cmdList.Flags().StringVarP(&paramsBundleDownloadList.Cursor, "cursor", "c", "", "Send cursor to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().Int64VarP(&paramsBundleDownloadList.BundleRegistrationId, "bundle-registration-id", "b", 0, "BundleRegistration ID")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	BundleDownloads.AddCommand(cmdList)
}
