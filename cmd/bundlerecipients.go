package cmd

import "github.com/spf13/cobra"
import (
	"fmt"
	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/bundlerecipient"
	"os"
)

var (
	_ = files_sdk.Config{}
	_ = bundle_recipient.Client{}
	_ = lib.OnlyFields
	_ = fmt.Println
	_ = os.Exit
)

var (
	BundleRecipients = &cobra.Command{
		Use:  "bundle-recipients [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func BundleRecipientsInit() {
	var fieldsList string
	paramsBundleRecipientList := files_sdk.BundleRecipientListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsBundleRecipientList
			params.MaxPages = MaxPagesList
			it := bundle_recipient.List(params)

			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().IntVarP(&paramsBundleRecipientList.Page, "page", "p", 0, "List Bundle Recipients")
	cmdList.Flags().IntVarP(&paramsBundleRecipientList.PerPage, "per-page", "e", 0, "List Bundle Recipients")
	cmdList.Flags().StringVarP(&paramsBundleRecipientList.Action, "action", "a", "", "List Bundle Recipients")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "f", "", "comma separated list of field names to include in response")
	BundleRecipients.AddCommand(cmdList)
}
