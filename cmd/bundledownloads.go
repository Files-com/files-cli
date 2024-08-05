package cmd

import (
	"fmt"

	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	bundle_download "github.com/Files-com/files-sdk-go/v3/bundledownload"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(BundleDownloads())
}

func BundleDownloads() *cobra.Command {
	BundleDownloads := &cobra.Command{
		Use:  "bundle-downloads [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command bundle-downloads\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsBundleDownloadList := files_sdk.BundleDownloadListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Bundle Downloads",
		Long:    `List Bundle Downloads`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsBundleDownloadList
			params.MaxPages = MaxPagesList

			client := bundle_download.Client{Config: config}
			it, err := client.List(params, files_sdk.WithContext(ctx))
			it.OnPageError = func(err error) (*[]interface{}, error) {
				overriddenValues, newErr := lib.ErrorWithOriginalResponse(err, config.Logger)
				values, ok := overriddenValues.([]interface{})
				if ok {
					return &values, newErr
				} else {
					return &[]interface{}{}, newErr
				}
			}
			if err != nil {
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			if len(filterbyList) > 0 {
				listFilter = func(i interface{}) (interface{}, bool, error) {
					matchOk, err := lib.MatchFilter(filterbyList, i)
					return i, matchOk, err
				}
			}
			err = lib.FormatIter(ctx, it, Profile(cmd).Current().SetResourceFormat(cmd, formatList), fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdList.Flags().StringToStringVar(&filterbyList, "filter-by", filterbyList, `Client side filtering: field-name=*.jpg,field-name=?ello`)

	cmdList.Flags().StringVar(&paramsBundleDownloadList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsBundleDownloadList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().Int64Var(&paramsBundleDownloadList.BundleId, "bundle-id", 0, "Bundle ID")
	cmdList.Flags().Int64Var(&paramsBundleDownloadList.BundleRegistrationId, "bundle-registration-id", 0, "BundleRegistration ID")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	BundleDownloads.AddCommand(cmdList)
	return BundleDownloads
}
