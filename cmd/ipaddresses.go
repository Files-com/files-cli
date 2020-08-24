package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	"fmt"
	"os"

	ip_address "github.com/Files-com/files-sdk-go/ipaddress"
)

var (
	IpAddresses = &cobra.Command{
		Use:  "ip-addresses [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func IpAddressesInit() {
	var fieldsList string
	paramsIpAddressList := files_sdk.IpAddressListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsIpAddressList
			params.MaxPages = MaxPagesList
			it := ip_address.List(params)

			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().IntVarP(&paramsIpAddressList.Page, "page", "p", 0, "Current page number.")
	cmdList.Flags().IntVarP(&paramsIpAddressList.PerPage, "per-page", "e", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVarP(&paramsIpAddressList.Action, "action", "a", "", "Deprecated: If set to `count` returns a count of matching records rather than the records themselves.")
	cmdList.Flags().StringVarP(&paramsIpAddressList.Cursor, "cursor", "c", "", "Send cursor to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	IpAddresses.AddCommand(cmdList)
	var fieldsGetReserved string
	paramsIpAddressGetReserved := files_sdk.IpAddressGetReservedParams{}
	cmdGetReserved := &cobra.Command{
		Use: "get-reserved",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := ip_address.GetReserved(paramsIpAddressGetReserved)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsGetReserved)
		},
	}
	cmdGetReserved.Flags().IntVarP(&paramsIpAddressGetReserved.Page, "page", "p", 0, "Current page number.")
	cmdGetReserved.Flags().IntVarP(&paramsIpAddressGetReserved.PerPage, "per-page", "e", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdGetReserved.Flags().StringVarP(&paramsIpAddressGetReserved.Action, "action", "a", "", "Deprecated: If set to `count` returns a count of matching records rather than the records themselves.")
	cmdGetReserved.Flags().StringVarP(&paramsIpAddressGetReserved.Cursor, "cursor", "c", "", "Send cursor to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")

	cmdGetReserved.Flags().StringVarP(&fieldsGetReserved, "fields", "", "", "comma separated list of field names")
	IpAddresses.AddCommand(cmdGetReserved)
}
