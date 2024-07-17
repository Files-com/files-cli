package cmd

import (
	"fmt"
	"reflect"

	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	as2_partner "github.com/Files-com/files-sdk-go/v3/as2partner"
	flib "github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/spf13/cobra"
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
		Short:   "List AS2 Partners",
		Long:    `List AS2 Partners`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsAs2PartnerList
			params.MaxPages = MaxPagesList

			client := as2_partner.Client{Config: config}
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
			err = lib.FormatIter(ctx, it, Profile(cmd).Current().SetResourceFormat(cmd, formatList), fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdList.Flags().StringToStringVar(&filterbyList, "filter-by", filterbyList, `Client side filtering: field-name=*.jpg,field-name=?ello`)

	cmdList.Flags().StringVar(&paramsAs2PartnerList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsAs2PartnerList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVar(&paramsAs2PartnerList.Action, "action", "", "")
	cmdList.Flags().Int64Var(&paramsAs2PartnerList.Page, "page", 0, "")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	As2Partners.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsAs2PartnerFind := files_sdk.As2PartnerFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show AS2 Partner`,
		Long:  `Show AS2 Partner`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := as2_partner.Client{Config: config}

			var as2Partner interface{}
			var err error
			as2Partner, err = client.Find(paramsAs2PartnerFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), as2Partner, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsAs2PartnerFind.Id, "id", 0, "As2 Partner ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	As2Partners.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	createEnableDedicatedIps := true
	paramsAs2PartnerCreate := files_sdk.As2PartnerCreateParams{}
	As2PartnerCreateMdnValidationLevel := ""
	As2PartnerCreateServerCertificate := ""

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create AS2 Partner`,
		Long:  `Create AS2 Partner`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := as2_partner.Client{Config: config}

			var As2PartnerCreateMdnValidationLevelErr error
			paramsAs2PartnerCreate.MdnValidationLevel, As2PartnerCreateMdnValidationLevelErr = lib.FetchKey("mdn-validation-level", paramsAs2PartnerCreate.MdnValidationLevel.Enum(), As2PartnerCreateMdnValidationLevel)
			if As2PartnerCreateMdnValidationLevel != "" && As2PartnerCreateMdnValidationLevelErr != nil {
				return As2PartnerCreateMdnValidationLevelErr
			}
			var As2PartnerCreateServerCertificateErr error
			paramsAs2PartnerCreate.ServerCertificate, As2PartnerCreateServerCertificateErr = lib.FetchKey("server-certificate", paramsAs2PartnerCreate.ServerCertificate.Enum(), As2PartnerCreateServerCertificate)
			if As2PartnerCreateServerCertificate != "" && As2PartnerCreateServerCertificateErr != nil {
				return As2PartnerCreateServerCertificateErr
			}

			if cmd.Flags().Changed("enable-dedicated-ips") {
				paramsAs2PartnerCreate.EnableDedicatedIps = flib.Bool(createEnableDedicatedIps)
			}

			var as2Partner interface{}
			var err error
			as2Partner, err = client.Create(paramsAs2PartnerCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), as2Partner, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().BoolVar(&createEnableDedicatedIps, "enable-dedicated-ips", createEnableDedicatedIps, "If `true`, we will use your site's dedicated IPs for all outbound connections to this AS2 PArtner.")
	cmdCreate.Flags().StringVar(&paramsAs2PartnerCreate.HttpAuthUsername, "http-auth-username", "", "Username to send to server for HTTP Authentication.")
	cmdCreate.Flags().StringVar(&paramsAs2PartnerCreate.HttpAuthPassword, "http-auth-password", "", "Password to send to server for HTTP Authentication.")
	cmdCreate.Flags().StringVar(&As2PartnerCreateMdnValidationLevel, "mdn-validation-level", "", fmt.Sprintf("How should Files.com evaluate message transfer success based on a partner's MDN response?  This setting does not affect MDN storage; all MDNs received from a partner are always stored. `none`: MDN is stored for informational purposes only, a successful HTTPS transfer is a successful AS2 transfer. `weak`: Inspect the MDN for MIC and Disposition only. `normal`: `weak` plus validate MDN signature matches body, `strict`: `normal` but do not allow signatures from self-signed or incorrectly purposed certificates. %v", reflect.ValueOf(paramsAs2PartnerCreate.MdnValidationLevel.Enum()).MapKeys()))
	cmdCreate.Flags().StringVar(&As2PartnerCreateServerCertificate, "server-certificate", "", fmt.Sprintf("Should we require that the remote HTTP server have a valid SSL Certificate for HTTPS? %v", reflect.ValueOf(paramsAs2PartnerCreate.ServerCertificate.Enum()).MapKeys()))
	cmdCreate.Flags().Int64Var(&paramsAs2PartnerCreate.As2StationId, "as2-station-id", 0, "ID of the AS2 Station associated with this partner.")
	cmdCreate.Flags().StringVar(&paramsAs2PartnerCreate.Name, "name", "", "The partner's formal AS2 name.")
	cmdCreate.Flags().StringVar(&paramsAs2PartnerCreate.Uri, "uri", "", "Public URI where we will send the AS2 messages (via HTTP/HTTPS).")
	cmdCreate.Flags().StringVar(&paramsAs2PartnerCreate.PublicCertificate, "public-certificate", "", "Public certificate for AS2 Partner.  Note: This is the certificate for AS2 message security, not a certificate used for HTTPS authentication.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	As2Partners.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	updateEnableDedicatedIps := true
	paramsAs2PartnerUpdate := files_sdk.As2PartnerUpdateParams{}
	As2PartnerUpdateMdnValidationLevel := ""
	As2PartnerUpdateServerCertificate := ""

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update AS2 Partner`,
		Long:  `Update AS2 Partner`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := as2_partner.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.As2PartnerUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			var As2PartnerUpdateMdnValidationLevelErr error
			paramsAs2PartnerUpdate.MdnValidationLevel, As2PartnerUpdateMdnValidationLevelErr = lib.FetchKey("mdn-validation-level", paramsAs2PartnerUpdate.MdnValidationLevel.Enum(), As2PartnerUpdateMdnValidationLevel)
			if As2PartnerUpdateMdnValidationLevel != "" && As2PartnerUpdateMdnValidationLevelErr != nil {
				return As2PartnerUpdateMdnValidationLevelErr
			}
			var As2PartnerUpdateServerCertificateErr error
			paramsAs2PartnerUpdate.ServerCertificate, As2PartnerUpdateServerCertificateErr = lib.FetchKey("server-certificate", paramsAs2PartnerUpdate.ServerCertificate.Enum(), As2PartnerUpdateServerCertificate)
			if As2PartnerUpdateServerCertificate != "" && As2PartnerUpdateServerCertificateErr != nil {
				return As2PartnerUpdateServerCertificateErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsAs2PartnerUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("enable-dedicated-ips") {
				mapParams["enable_dedicated_ips"] = updateEnableDedicatedIps
			}
			if cmd.Flags().Changed("http-auth-username") {
				lib.FlagUpdate(cmd, "http_auth_username", paramsAs2PartnerUpdate.HttpAuthUsername, mapParams)
			}
			if cmd.Flags().Changed("http-auth-password") {
				lib.FlagUpdate(cmd, "http_auth_password", paramsAs2PartnerUpdate.HttpAuthPassword, mapParams)
			}
			if cmd.Flags().Changed("mdn-validation-level") {
				lib.FlagUpdate(cmd, "mdn_validation_level", paramsAs2PartnerUpdate.MdnValidationLevel, mapParams)
			}
			if cmd.Flags().Changed("server-certificate") {
				lib.FlagUpdate(cmd, "server_certificate", paramsAs2PartnerUpdate.ServerCertificate, mapParams)
			}
			if cmd.Flags().Changed("name") {
				lib.FlagUpdate(cmd, "name", paramsAs2PartnerUpdate.Name, mapParams)
			}
			if cmd.Flags().Changed("uri") {
				lib.FlagUpdate(cmd, "uri", paramsAs2PartnerUpdate.Uri, mapParams)
			}
			if cmd.Flags().Changed("public-certificate") {
				lib.FlagUpdate(cmd, "public_certificate", paramsAs2PartnerUpdate.PublicCertificate, mapParams)
			}

			var as2Partner interface{}
			var err error
			as2Partner, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), as2Partner, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsAs2PartnerUpdate.Id, "id", 0, "As2 Partner ID.")
	cmdUpdate.Flags().BoolVar(&updateEnableDedicatedIps, "enable-dedicated-ips", updateEnableDedicatedIps, "If `true`, we will use your site's dedicated IPs for all outbound connections to this AS2 PArtner.")
	cmdUpdate.Flags().StringVar(&paramsAs2PartnerUpdate.HttpAuthUsername, "http-auth-username", "", "Username to send to server for HTTP Authentication.")
	cmdUpdate.Flags().StringVar(&paramsAs2PartnerUpdate.HttpAuthPassword, "http-auth-password", "", "Password to send to server for HTTP Authentication.")
	cmdUpdate.Flags().StringVar(&As2PartnerUpdateMdnValidationLevel, "mdn-validation-level", "", fmt.Sprintf("How should Files.com evaluate message transfer success based on a partner's MDN response?  This setting does not affect MDN storage; all MDNs received from a partner are always stored. `none`: MDN is stored for informational purposes only, a successful HTTPS transfer is a successful AS2 transfer. `weak`: Inspect the MDN for MIC and Disposition only. `normal`: `weak` plus validate MDN signature matches body, `strict`: `normal` but do not allow signatures from self-signed or incorrectly purposed certificates. %v", reflect.ValueOf(paramsAs2PartnerUpdate.MdnValidationLevel.Enum()).MapKeys()))
	cmdUpdate.Flags().StringVar(&As2PartnerUpdateServerCertificate, "server-certificate", "", fmt.Sprintf("Should we require that the remote HTTP server have a valid SSL Certificate for HTTPS? %v", reflect.ValueOf(paramsAs2PartnerUpdate.ServerCertificate.Enum()).MapKeys()))
	cmdUpdate.Flags().StringVar(&paramsAs2PartnerUpdate.Name, "name", "", "The partner's formal AS2 name.")
	cmdUpdate.Flags().StringVar(&paramsAs2PartnerUpdate.Uri, "uri", "", "Public URI where we will send the AS2 messages (via HTTP/HTTPS).")
	cmdUpdate.Flags().StringVar(&paramsAs2PartnerUpdate.PublicCertificate, "public-certificate", "", "Public certificate for AS2 Partner.  Note: This is the certificate for AS2 message security, not a certificate used for HTTPS authentication.")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	As2Partners.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsAs2PartnerDelete := files_sdk.As2PartnerDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete AS2 Partner`,
		Long:  `Delete AS2 Partner`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := as2_partner.Client{Config: config}

			var err error
			err = client.Delete(paramsAs2PartnerDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsAs2PartnerDelete.Id, "id", 0, "As2 Partner ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	As2Partners.AddCommand(cmdDelete)
	return As2Partners
}
