package cmd

import "github.com/spf13/cobra"
import (
	"fmt"
	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/ipaddress"
	"os"
)

var (
	_ = files_sdk.Config{}
	_ = ip_address.Client{}
	_ = lib.OnlyFields
	_ = fmt.Println
	_ = os.Exit
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
	cmdList.Flags().IntVarP(&paramsIpAddressList.Page, "page", "p", 0, "List IP Addresses associated with the current site")
	cmdList.Flags().IntVarP(&paramsIpAddressList.PerPage, "per-page", "e", 0, "List IP Addresses associated with the current site")
	cmdList.Flags().StringVarP(&paramsIpAddressList.Action, "action", "a", "", "List IP Addresses associated with the current site")
	cmdList.Flags().StringVarP(&paramsIpAddressList.Cursor, "cursor", "c", "", "List IP Addresses associated with the current site")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "f", "", "comma separated list of field names to include in response")
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
	cmdGetReserved.Flags().IntVarP(&paramsIpAddressGetReserved.Page, "page", "p", 0, "List all possible public IP addresses")
	cmdGetReserved.Flags().IntVarP(&paramsIpAddressGetReserved.PerPage, "per-page", "e", 0, "List all possible public IP addresses")
	cmdGetReserved.Flags().StringVarP(&paramsIpAddressGetReserved.Action, "action", "a", "", "List all possible public IP addresses")
	cmdGetReserved.Flags().StringVarP(&paramsIpAddressGetReserved.Cursor, "cursor", "c", "", "List all possible public IP addresses")
	cmdGetReserved.Flags().StringVarP(&fieldsGetReserved, "fields", "f", "", "comma separated list of field names")
	IpAddresses.AddCommand(cmdGetReserved)
}
