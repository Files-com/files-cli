package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	as2_partner "github.com/Files-com/files-sdk-go/v2/as2partner"
)

var (
	As2Partners = &cobra.Command{}
)

func As2PartnersInit() {
	As2Partners = &cobra.Command{
		Use:  "as2-partners [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command as2-partners\n\t%v", args[0])
		},
	}
	var fieldsList string
	var formatList string
	usePagerList := true
	paramsAs2PartnerList := files_sdk.As2PartnerListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List As2 Partners",
		Long:  `List As2 Partners`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
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
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(ctx, it, formatList, fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
		},
	}

	cmdList.Flags().StringVar(&paramsAs2PartnerList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsAs2PartnerList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	As2Partners.AddCommand(cmdList)
	var fieldsFind string
	var formatFind string
	usePagerFind := true
	paramsAs2PartnerFind := files_sdk.As2PartnerFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show As2 Partner`,
		Long:  `Show As2 Partner`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := as2_partner.Client{Config: *config}

			var as2Partner interface{}
			var err error
			as2Partner, err = client.Find(ctx, paramsAs2PartnerFind)
			lib.HandleResponse(ctx, as2Partner, err, formatFind, fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdFind.Flags().Int64Var(&paramsAs2PartnerFind.Id, "id", 0, "As2 Partner ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	As2Partners.AddCommand(cmdFind)
	var fieldsCreate string
	var formatCreate string
	usePagerCreate := true
	paramsAs2PartnerCreate := files_sdk.As2PartnerCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create As2 Partner`,
		Long:  `Create As2 Partner`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := as2_partner.Client{Config: *config}

			var as2Partner interface{}
			var err error
			as2Partner, err = client.Create(ctx, paramsAs2PartnerCreate)
			lib.HandleResponse(ctx, as2Partner, err, formatCreate, fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdCreate.Flags().StringVar(&paramsAs2PartnerCreate.Name, "name", "", "AS2 Name")
	cmdCreate.Flags().StringVar(&paramsAs2PartnerCreate.Uri, "uri", "", "URL base for AS2 responses")
	cmdCreate.Flags().StringVar(&paramsAs2PartnerCreate.PublicCertificate, "public-certificate", "", "")
	cmdCreate.Flags().Int64Var(&paramsAs2PartnerCreate.As2StationId, "as2-station-id", 0, "Id of As2Station for this partner")
	cmdCreate.Flags().StringVar(&paramsAs2PartnerCreate.ServerCertificate, "server-certificate", "", "Remote server certificate security setting")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	As2Partners.AddCommand(cmdCreate)
	var fieldsUpdate string
	var formatUpdate string
	usePagerUpdate := true
	paramsAs2PartnerUpdate := files_sdk.As2PartnerUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update As2 Partner`,
		Long:  `Update As2 Partner`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := as2_partner.Client{Config: *config}

			var as2Partner interface{}
			var err error
			as2Partner, err = client.Update(ctx, paramsAs2PartnerUpdate)
			lib.HandleResponse(ctx, as2Partner, err, formatUpdate, fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsAs2PartnerUpdate.Id, "id", 0, "As2 Partner ID.")
	cmdUpdate.Flags().StringVar(&paramsAs2PartnerUpdate.Name, "name", "", "AS2 Name")
	cmdUpdate.Flags().StringVar(&paramsAs2PartnerUpdate.Uri, "uri", "", "URL base for AS2 responses")
	cmdUpdate.Flags().StringVar(&paramsAs2PartnerUpdate.ServerCertificate, "server-certificate", "", "Remote server certificate security setting")
	cmdUpdate.Flags().StringVar(&paramsAs2PartnerUpdate.PublicCertificate, "public-certificate", "", "")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	cmdUpdate.Flags().StringVarP(&formatUpdate, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	As2Partners.AddCommand(cmdUpdate)
	var fieldsDelete string
	var formatDelete string
	usePagerDelete := true
	paramsAs2PartnerDelete := files_sdk.As2PartnerDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete As2 Partner`,
		Long:  `Delete As2 Partner`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := as2_partner.Client{Config: *config}

			var err error
			err = client.Delete(ctx, paramsAs2PartnerDelete)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
		},
	}
	cmdDelete.Flags().Int64Var(&paramsAs2PartnerDelete.Id, "id", 0, "As2 Partner ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	As2Partners.AddCommand(cmdDelete)
}
