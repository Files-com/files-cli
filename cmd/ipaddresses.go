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
			client := ip_address.Client{Config: files_sdk.GlobalConfig}
			it, err := client.List(params)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			err = lib.JsonMarshalIter(it, fieldsList)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	cmdList.Flags().StringVarP(&paramsIpAddressList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().IntVarP(&paramsIpAddressList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	IpAddresses.AddCommand(cmdList)
	var fieldsGetReserved string
	paramsIpAddressGetReserved := files_sdk.IpAddressGetReservedParams{}
	cmdGetReserved := &cobra.Command{
		Use: "get-reserved",
		Run: func(cmd *cobra.Command, args []string) {
			client := ip_address.Client{Config: files_sdk.GlobalConfig}
			result, err := client.GetReserved(paramsIpAddressGetReserved)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = lib.JsonMarshal(result, fieldsGetReserved)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	cmdGetReserved.Flags().StringVarP(&paramsIpAddressGetReserved.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdGetReserved.Flags().IntVarP(&paramsIpAddressGetReserved.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdGetReserved.Flags().StringVarP(&fieldsGetReserved, "fields", "", "", "comma separated list of field names")
	IpAddresses.AddCommand(cmdGetReserved)
}
