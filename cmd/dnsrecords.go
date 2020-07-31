package cmd

import "github.com/spf13/cobra"
import (
	"fmt"
	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/dnsrecord"
	"os"
)

var (
	_ = files_sdk.Config{}
	_ = dns_record.Client{}
	_ = lib.OnlyFields
	_ = fmt.Println
	_ = os.Exit
)

var (
	DnsRecords = &cobra.Command{
		Use:  "dns-records [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func DnsRecordsInit() {
	var fieldsList string
	paramsDnsRecordList := files_sdk.DnsRecordListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsDnsRecordList
			params.MaxPages = MaxPagesList
			it := dns_record.List(params)

			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().IntVarP(&paramsDnsRecordList.Page, "page", "p", 0, "Show site DNS configuration")
	cmdList.Flags().IntVarP(&paramsDnsRecordList.PerPage, "per-page", "e", 0, "Show site DNS configuration")
	cmdList.Flags().StringVarP(&paramsDnsRecordList.Action, "action", "a", "", "Show site DNS configuration")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "f", "", "comma separated list of field names to include in response")
	DnsRecords.AddCommand(cmdList)
}
