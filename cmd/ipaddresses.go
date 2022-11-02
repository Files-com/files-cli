package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	ip_address "github.com/Files-com/files-sdk-go/v2/ipaddress"
)

func init() {
	RootCmd.AddCommand(IpAddresses())
}

func IpAddresses() *cobra.Command {
	IpAddresses := &cobra.Command{
		Use:  "ip-addresses [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command ip-addresses\n\t%v", args[0])
		},
	}
	var fieldsList string
	var formatList string
	usePagerList := true
	paramsIpAddressList := files_sdk.IpAddressListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List IP Addresses associated with the current site",
		Long:  `List IP Addresses associated with the current site`,
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsIpAddressList
			params.MaxPages = MaxPagesList

			client := ip_address.Client{Config: *config}
			it, err := client.List(ctx, params)
			it.OnPageError = func(err error) (*[]interface{}, error) {
				overriddenValues, newErr := lib.ErrorWithOriginalResponse(err, config.Logger())
				values, ok := overriddenValues.([]interface{})
				if ok {
					return &values, newErr
				} else {
					return &[]interface{}{}, newErr
				}
			}
			if err != nil {
				lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(ctx, it, formatList, fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			if err != nil {
				lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}

	cmdList.Flags().StringVar(&paramsIpAddressList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsIpAddressList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	IpAddresses.AddCommand(cmdList)
	var fieldsGetExavaultReserved string
	var formatGetExavaultReserved string
	usePagerGetExavaultReserved := true
	paramsIpAddressGetExavaultReserved := files_sdk.IpAddressGetExavaultReservedParams{}

	cmdGetExavaultReserved := &cobra.Command{
		Use:   "get-exavault-reserved",
		Short: `List all possible public ExaVault IP addresses`,
		Long:  `List all possible public ExaVault IP addresses`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := ip_address.Client{Config: *config}

			var publicIpAddressCollection interface{}
			var err error
			publicIpAddressCollection, err = client.GetExavaultReserved(ctx, paramsIpAddressGetExavaultReserved)
			lib.HandleResponse(ctx, Profile(cmd), publicIpAddressCollection, err, formatGetExavaultReserved, fieldsGetExavaultReserved, usePagerGetExavaultReserved, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdGetExavaultReserved.Flags().StringVar(&paramsIpAddressGetExavaultReserved.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdGetExavaultReserved.Flags().Int64Var(&paramsIpAddressGetExavaultReserved.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdGetExavaultReserved.Flags().StringVarP(&fieldsGetExavaultReserved, "fields", "", "", "comma separated list of field names")
	cmdGetExavaultReserved.Flags().StringVarP(&formatGetExavaultReserved, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdGetExavaultReserved.Flags().BoolVar(&usePagerGetExavaultReserved, "use-pager", usePagerGetExavaultReserved, "Use $PAGER (.ie less, more, etc)")

	IpAddresses.AddCommand(cmdGetExavaultReserved)
	var fieldsGetReserved string
	var formatGetReserved string
	usePagerGetReserved := true
	paramsIpAddressGetReserved := files_sdk.IpAddressGetReservedParams{}

	cmdGetReserved := &cobra.Command{
		Use:   "get-reserved",
		Short: `List all possible public IP addresses`,
		Long:  `List all possible public IP addresses`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := ip_address.Client{Config: *config}

			var publicIpAddressCollection interface{}
			var err error
			publicIpAddressCollection, err = client.GetReserved(ctx, paramsIpAddressGetReserved)
			lib.HandleResponse(ctx, Profile(cmd), publicIpAddressCollection, err, formatGetReserved, fieldsGetReserved, usePagerGetReserved, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdGetReserved.Flags().StringVar(&paramsIpAddressGetReserved.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdGetReserved.Flags().Int64Var(&paramsIpAddressGetReserved.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdGetReserved.Flags().StringVarP(&fieldsGetReserved, "fields", "", "", "comma separated list of field names")
	cmdGetReserved.Flags().StringVarP(&formatGetReserved, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdGetReserved.Flags().BoolVar(&usePagerGetReserved, "use-pager", usePagerGetReserved, "Use $PAGER (.ie less, more, etc)")

	IpAddresses.AddCommand(cmdGetReserved)
	return IpAddresses
}
