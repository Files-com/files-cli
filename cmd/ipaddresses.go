package cmd

import (
	"fmt"

	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	ip_address "github.com/Files-com/files-sdk-go/v3/ipaddress"
	"github.com/spf13/cobra"
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
	filterbyList := make(map[string]string)
	paramsIpAddressList := files_sdk.IpAddressListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List IP Addresses associated with the current site",
		Long:    `List IP Addresses associated with the current site`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsIpAddressList
			params.MaxPages = MaxPagesList

			client := ip_address.Client{Config: config}
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
			err = lib.FormatIter(ctx, it, formatList, fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdList.Flags().StringToStringVar(&filterbyList, "filter-by", filterbyList, `Client side filtering: field-name=*.jpg,field-name=?ello`)

	cmdList.Flags().StringVar(&paramsIpAddressList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
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
	filterbyGetExavaultReserved := make(map[string]string)
	paramsIpAddressGetExavaultReserved := files_sdk.IpAddressGetExavaultReservedParams{}
	var MaxPagesGetExavaultReserved int64

	cmdGetExavaultReserved := &cobra.Command{
		Use:   "get-exavault-reserved",
		Short: "List all possible public ExaVault IP addresses",
		Long:  `List all possible public ExaVault IP addresses`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsIpAddressGetExavaultReserved
			params.MaxPages = MaxPagesGetExavaultReserved

			client := ip_address.Client{Config: config}
			it, err := client.GetExavaultReserved(params, files_sdk.WithContext(ctx))
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
			if len(filterbyGetExavaultReserved) > 0 {
				listFilter = func(i interface{}) (interface{}, bool, error) {
					matchOk, err := lib.MatchFilter(filterbyGetExavaultReserved, i)
					return i, matchOk, err
				}
			}
			err = lib.FormatIter(ctx, it, formatGetExavaultReserved, fieldsGetExavaultReserved, usePagerGetExavaultReserved, listFilter, cmd.OutOrStdout())
			return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdGetExavaultReserved.Flags().StringToStringVar(&filterbyGetExavaultReserved, "filter-by", filterbyGetExavaultReserved, `Client side filtering: field-name=*.jpg,field-name=?ello`)

	cmdGetExavaultReserved.Flags().StringVar(&paramsIpAddressGetExavaultReserved.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
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
	filterbyGetReserved := make(map[string]string)
	paramsIpAddressGetReserved := files_sdk.IpAddressGetReservedParams{}
	var MaxPagesGetReserved int64

	cmdGetReserved := &cobra.Command{
		Use:   "get-reserved",
		Short: "List all possible public IP addresses",
		Long:  `List all possible public IP addresses`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsIpAddressGetReserved
			params.MaxPages = MaxPagesGetReserved

			client := ip_address.Client{Config: config}
			it, err := client.GetReserved(params, files_sdk.WithContext(ctx))
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
			if len(filterbyGetReserved) > 0 {
				listFilter = func(i interface{}) (interface{}, bool, error) {
					matchOk, err := lib.MatchFilter(filterbyGetReserved, i)
					return i, matchOk, err
				}
			}
			err = lib.FormatIter(ctx, it, formatGetReserved, fieldsGetReserved, usePagerGetReserved, listFilter, cmd.OutOrStdout())
			return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdGetReserved.Flags().StringToStringVar(&filterbyGetReserved, "filter-by", filterbyGetReserved, `Client side filtering: field-name=*.jpg,field-name=?ello`)

	cmdGetReserved.Flags().StringVar(&paramsIpAddressGetReserved.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
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
