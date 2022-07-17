package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	ip_address "github.com/Files-com/files-sdk-go/v2/ipaddress"
)

var (
	IpAddresses = &cobra.Command{}
)

func IpAddressesInit() {
	IpAddresses = &cobra.Command{
		Use:  "ip-addresses [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command ip-addresses\n\t%v", args[0])
		},
	}
	var fieldsList string
	var formatList string
	paramsIpAddressList := files_sdk.IpAddressListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List IP Addresses associated with the current site",
		Long:  `List IP Addresses associated with the current site`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsIpAddressList
			params.MaxPages = MaxPagesList

			client := ip_address.Client{Config: *config}
			it, err := client.List(ctx, params)
			it.OnPageError = func(err error) (*[]interface{}, error) {
				overriddenValues, newErr := lib.ErrorWithOriginalResponse(err, formatList, config.Logger())
				values, ok := overriddenValues.([]interface{})
				if ok {
					return &values, newErr
				} else {
					return &[]interface{}{}, newErr
				}
			}
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(it, formatList, fieldsList, listFilter, cmd.OutOrStdout())
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
		},
	}

	cmdList.Flags().StringVar(&paramsIpAddressList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsIpAddressList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright")
	IpAddresses.AddCommand(cmdList)
	var fieldsGetReserved string
	var formatGetReserved string
	paramsIpAddressGetReserved := files_sdk.IpAddressGetReservedParams{}

	cmdGetReserved := &cobra.Command{
		Use:   "get-reserved",
		Short: `List all possible public IP addresses`,
		Long:  `List all possible public IP addresses`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := ip_address.Client{Config: *config}

			var publicIpAddressCollection interface{}
			var err error
			publicIpAddressCollection, err = client.GetReserved(ctx, paramsIpAddressGetReserved)
			lib.HandleResponse(ctx, publicIpAddressCollection, err, formatGetReserved, fieldsGetReserved, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdGetReserved.Flags().StringVar(&paramsIpAddressGetReserved.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdGetReserved.Flags().Int64Var(&paramsIpAddressGetReserved.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdGetReserved.Flags().StringVarP(&fieldsGetReserved, "fields", "", "", "comma separated list of field names")
	cmdGetReserved.Flags().StringVarP(&formatGetReserved, "format", "", "table", "json, csv, table, table-dark, table-bright")
	IpAddresses.AddCommand(cmdGetReserved)
}
