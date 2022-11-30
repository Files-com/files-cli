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
	var fieldsList []string
	var formatList []string
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
				return lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(ctx, it, formatList, fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			if err != nil {
				return lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}

	cmdList.Flags().StringVar(&paramsIpAddressList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsIpAddressList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
        table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
        json-styles: {raw, pretty}
        `)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	IpAddresses.AddCommand(cmdList)
	var fieldsGetExavaultReserved []string
	var formatGetExavaultReserved []string
	usePagerGetExavaultReserved := true
	paramsIpAddressGetExavaultReserved := files_sdk.IpAddressGetExavaultReservedParams{}
	var MaxPagesGetExavaultReserved int64

	cmdGetExavaultReserved := &cobra.Command{
		Use:   "get-exavault-reserved",
		Short: "List all possible public ExaVault IP addresses",
		Long:  `List all possible public ExaVault IP addresses`,
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsIpAddressGetExavaultReserved
			params.MaxPages = MaxPagesGetExavaultReserved

			client := ip_address.Client{Config: *config}
			it, err := client.GetExavaultReserved(ctx, params)
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
				return lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(ctx, it, formatGetExavaultReserved, fieldsGetExavaultReserved, usePagerGetExavaultReserved, listFilter, cmd.OutOrStdout())
			if err != nil {
				return lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}

	cmdGetExavaultReserved.Flags().StringVar(&paramsIpAddressGetExavaultReserved.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdGetExavaultReserved.Flags().Int64Var(&paramsIpAddressGetExavaultReserved.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdGetExavaultReserved.Flags().Int64VarP(&MaxPagesGetExavaultReserved, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdGetExavaultReserved.Flags().StringSliceVar(&fieldsGetExavaultReserved, "fields", []string{}, "comma separated list of field names to include in response")
	cmdGetExavaultReserved.Flags().StringSliceVar(&formatGetExavaultReserved, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
        table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
        json-styles: {raw, pretty}
        `)
	cmdGetExavaultReserved.Flags().BoolVar(&usePagerGetExavaultReserved, "use-pager", usePagerGetExavaultReserved, "Use $PAGER (.ie less, more, etc)")
	IpAddresses.AddCommand(cmdGetExavaultReserved)
	var fieldsGetReserved []string
	var formatGetReserved []string
	usePagerGetReserved := true
	paramsIpAddressGetReserved := files_sdk.IpAddressGetReservedParams{}
	var MaxPagesGetReserved int64

	cmdGetReserved := &cobra.Command{
		Use:   "get-reserved",
		Short: "List all possible public IP addresses",
		Long:  `List all possible public IP addresses`,
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsIpAddressGetReserved
			params.MaxPages = MaxPagesGetReserved

			client := ip_address.Client{Config: *config}
			it, err := client.GetReserved(ctx, params)
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
				return lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(ctx, it, formatGetReserved, fieldsGetReserved, usePagerGetReserved, listFilter, cmd.OutOrStdout())
			if err != nil {
				return lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}

	cmdGetReserved.Flags().StringVar(&paramsIpAddressGetReserved.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdGetReserved.Flags().Int64Var(&paramsIpAddressGetReserved.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdGetReserved.Flags().Int64VarP(&MaxPagesGetReserved, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdGetReserved.Flags().StringSliceVar(&fieldsGetReserved, "fields", []string{}, "comma separated list of field names to include in response")
	cmdGetReserved.Flags().StringSliceVar(&formatGetReserved, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
        table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
        json-styles: {raw, pretty}
        `)
	cmdGetReserved.Flags().BoolVar(&usePagerGetReserved, "use-pager", usePagerGetReserved, "Use $PAGER (.ie less, more, etc)")
	IpAddresses.AddCommand(cmdGetReserved)
	return IpAddresses
}
