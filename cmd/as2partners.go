package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	flib "github.com/Files-com/files-sdk-go/v2/lib"

	"fmt"

	as2_partner "github.com/Files-com/files-sdk-go/v2/as2partner"
)

func init() {
	RootCmd.AddCommand(As2Partners())
}

func As2Partners() *cobra.Command {
	As2Partners := &cobra.Command{
		Use:  "as2-partners [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command as2-partners\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsAs2PartnerList := files_sdk.As2PartnerListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List As2 Partners",
		Long:    `List As2 Partners`,
		Args:    cobra.MinimumNArgs(0),
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsAs2PartnerList
			params.MaxPages = MaxPagesList

			client := as2_partner.Client{Config: *config}
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

	cmdList.Flags().StringVar(&paramsAs2PartnerList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsAs2PartnerList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
        table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
        json-styles: {raw, pretty}
        `)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	As2Partners.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsAs2PartnerFind := files_sdk.As2PartnerFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show As2 Partner`,
		Long:  `Show As2 Partner`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := as2_partner.Client{Config: *config}

			var as2Partner interface{}
			var err error
			as2Partner, err = client.Find(ctx, paramsAs2PartnerFind)
			lib.HandleResponse(ctx, Profile(cmd), as2Partner, err, formatFind, fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdFind.Flags().Int64Var(&paramsAs2PartnerFind.Id, "id", 0, "As2 Partner ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	As2Partners.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	createEnableDedicatedIps := true
	paramsAs2PartnerCreate := files_sdk.As2PartnerCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create As2 Partner`,
		Long:  `Create As2 Partner`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := as2_partner.Client{Config: *config}

			if cmd.Flags().Changed("enable-dedicated-ips") {
				paramsAs2PartnerCreate.EnableDedicatedIps = flib.Bool(createEnableDedicatedIps)
			}

			var as2Partner interface{}
			var err error
			as2Partner, err = client.Create(ctx, paramsAs2PartnerCreate)
			lib.HandleResponse(ctx, Profile(cmd), as2Partner, err, formatCreate, fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdCreate.Flags().StringVar(&paramsAs2PartnerCreate.Name, "name", "", "AS2 Name")
	cmdCreate.Flags().StringVar(&paramsAs2PartnerCreate.Uri, "uri", "", "URL base for AS2 responses")
	cmdCreate.Flags().StringVar(&paramsAs2PartnerCreate.PublicCertificate, "public-certificate", "", "")
	cmdCreate.Flags().Int64Var(&paramsAs2PartnerCreate.As2StationId, "as2-station-id", 0, "Id of As2Station for this partner")
	cmdCreate.Flags().StringVar(&paramsAs2PartnerCreate.ServerCertificate, "server-certificate", "", "Remote server certificate security setting")
	cmdCreate.Flags().BoolVar(&createEnableDedicatedIps, "enable-dedicated-ips", createEnableDedicatedIps, "")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	As2Partners.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	updateEnableDedicatedIps := true
	paramsAs2PartnerUpdate := files_sdk.As2PartnerUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update As2 Partner`,
		Long:  `Update As2 Partner`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := as2_partner.Client{Config: *config}

			mapParams, convertErr := lib.StructToMap(files_sdk.As2PartnerUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsAs2PartnerUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("name") {
				lib.FlagUpdate(cmd, "name", paramsAs2PartnerUpdate.Name, mapParams)
			}
			if cmd.Flags().Changed("uri") {
				lib.FlagUpdate(cmd, "uri", paramsAs2PartnerUpdate.Uri, mapParams)
			}
			if cmd.Flags().Changed("server-certificate") {
				lib.FlagUpdate(cmd, "server_certificate", paramsAs2PartnerUpdate.ServerCertificate, mapParams)
			}
			if cmd.Flags().Changed("public-certificate") {
				lib.FlagUpdate(cmd, "public_certificate", paramsAs2PartnerUpdate.PublicCertificate, mapParams)
			}
			if cmd.Flags().Changed("enable-dedicated-ips") {
				mapParams["enable_dedicated_ips"] = updateEnableDedicatedIps
			}

			var as2Partner interface{}
			var err error
			as2Partner, err = client.UpdateWithMap(ctx, mapParams)
			lib.HandleResponse(ctx, Profile(cmd), as2Partner, err, formatUpdate, fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsAs2PartnerUpdate.Id, "id", 0, "As2 Partner ID.")
	cmdUpdate.Flags().StringVar(&paramsAs2PartnerUpdate.Name, "name", "", "AS2 Name")
	cmdUpdate.Flags().StringVar(&paramsAs2PartnerUpdate.Uri, "uri", "", "URL base for AS2 responses")
	cmdUpdate.Flags().StringVar(&paramsAs2PartnerUpdate.ServerCertificate, "server-certificate", "", "Remote server certificate security setting")
	cmdUpdate.Flags().StringVar(&paramsAs2PartnerUpdate.PublicCertificate, "public-certificate", "", "")
	cmdUpdate.Flags().BoolVar(&updateEnableDedicatedIps, "enable-dedicated-ips", updateEnableDedicatedIps, "")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	As2Partners.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsAs2PartnerDelete := files_sdk.As2PartnerDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete As2 Partner`,
		Long:  `Delete As2 Partner`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := as2_partner.Client{Config: *config}

			var err error
			err = client.Delete(ctx, paramsAs2PartnerDelete)
			if err != nil {
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsAs2PartnerDelete.Id, "id", 0, "As2 Partner ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	As2Partners.AddCommand(cmdDelete)
	return As2Partners
}
