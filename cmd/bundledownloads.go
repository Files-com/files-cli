package cmd

import "github.com/spf13/cobra"
import (
	"fmt"
	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/bundledownload"
	"os"
)

var (
	_ = files_sdk.Config{}
	_ = bundle_download.Client{}
	_ = lib.OnlyFields
	_ = fmt.Println
	_ = os.Exit
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
			it := bundle_download.List(params)

			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().IntVarP(&paramsBundleDownloadList.Page, "page", "p", 0, "List Bundle Downloads")
	cmdList.Flags().IntVarP(&paramsBundleDownloadList.PerPage, "per-page", "e", 0, "List Bundle Downloads")
	cmdList.Flags().StringVarP(&paramsBundleDownloadList.Action, "action", "a", "", "List Bundle Downloads")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "f", "", "comma separated list of field names to include in response")
	BundleDownloads.AddCommand(cmdList)
}
