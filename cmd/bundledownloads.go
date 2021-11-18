package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	bundle_download "github.com/Files-com/files-sdk-go/v2/bundledownload"
)

var (
	BundleDownloads = &cobra.Command{}
)

func BundleDownloadsInit() {
	BundleDownloads = &cobra.Command{
		Use:  "bundle-downloads [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command bundle-downloads\n\t%v", args[0])
		},
	}
	var fieldsList string
	var formatList string
	paramsBundleDownloadList := files_sdk.BundleDownloadListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsBundleDownloadList
			params.MaxPages = MaxPagesList

			client := bundle_download.Client{Config: *config}
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

	cmdList.Flags().StringVar(&paramsBundleDownloadList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsBundleDownloadList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().Int64Var(&paramsBundleDownloadList.BundleId, "bundle-id", 0, "Bundle ID")
	cmdList.Flags().Int64Var(&paramsBundleDownloadList.BundleRegistrationId, "bundle-registration-id", 0, "BundleRegistration ID")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright")
	BundleDownloads.AddCommand(cmdList)
}
