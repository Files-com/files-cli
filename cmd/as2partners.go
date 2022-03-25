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
	paramsAs2PartnerList := files_sdk.As2PartnerListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsAs2PartnerList
			params.MaxPages = MaxPagesList

			client := as2_partner.Client{Config: *config}
			it, err := client.List(ctx, params)
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

	cmdList.Flags().StringVar(&paramsAs2PartnerList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsAs2PartnerList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright")
	As2Partners.AddCommand(cmdList)
	var fieldsFind string
	var formatFind string
	paramsAs2PartnerFind := files_sdk.As2PartnerFindParams{}

	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := as2_partner.Client{Config: *config}

			var result interface{}
			var err error
			result, err = client.Find(ctx, paramsAs2PartnerFind)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(result, formatFind, fieldsFind, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}
	cmdFind.Flags().Int64Var(&paramsAs2PartnerFind.Id, "id", 0, "As2 Partner ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-bright")
	As2Partners.AddCommand(cmdFind)
	var fieldsCreate string
	var formatCreate string
	paramsAs2PartnerCreate := files_sdk.As2PartnerCreateParams{}

	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := as2_partner.Client{Config: *config}

			var result interface{}
			var err error
			result, err = client.Create(ctx, paramsAs2PartnerCreate)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(result, formatCreate, fieldsCreate, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}
	cmdCreate.Flags().StringVar(&paramsAs2PartnerCreate.Name, "name", "", "AS2 Name")
	cmdCreate.Flags().StringVar(&paramsAs2PartnerCreate.Uri, "uri", "", "URL base for AS2 responses")
	cmdCreate.Flags().StringVar(&paramsAs2PartnerCreate.PublicCertificate, "public-certificate", "", "")
	cmdCreate.Flags().Int64Var(&paramsAs2PartnerCreate.As2StationId, "as2-station-id", 0, "Id of As2Station for this partner")
	cmdCreate.Flags().StringVar(&paramsAs2PartnerCreate.ServerCertificate, "server-certificate", "", "Remote server certificate security setting")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	As2Partners.AddCommand(cmdCreate)
	var fieldsUpdate string
	var formatUpdate string
	paramsAs2PartnerUpdate := files_sdk.As2PartnerUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := as2_partner.Client{Config: *config}

			var result interface{}
			var err error
			result, err = client.Update(ctx, paramsAs2PartnerUpdate)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(result, formatUpdate, fieldsUpdate, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsAs2PartnerUpdate.Id, "id", 0, "As2 Partner ID.")
	cmdUpdate.Flags().StringVar(&paramsAs2PartnerUpdate.Name, "name", "", "AS2 Name")
	cmdUpdate.Flags().StringVar(&paramsAs2PartnerUpdate.Uri, "uri", "", "URL base for AS2 responses")
	cmdUpdate.Flags().StringVar(&paramsAs2PartnerUpdate.ServerCertificate, "server-certificate", "", "Remote server certificate security setting")
	cmdUpdate.Flags().StringVar(&paramsAs2PartnerUpdate.PublicCertificate, "public-certificate", "", "")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	cmdUpdate.Flags().StringVarP(&formatUpdate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	As2Partners.AddCommand(cmdUpdate)
	var fieldsDelete string
	var formatDelete string
	paramsAs2PartnerDelete := files_sdk.As2PartnerDeleteParams{}

	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := as2_partner.Client{Config: *config}

			var result interface{}
			var err error
			result, err = client.Delete(ctx, paramsAs2PartnerDelete)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(result, formatDelete, fieldsDelete, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}
	cmdDelete.Flags().Int64Var(&paramsAs2PartnerDelete.Id, "id", 0, "As2 Partner ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-bright")
	As2Partners.AddCommand(cmdDelete)
}
