package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	partner_site_request "github.com/Files-com/files-sdk-go/v3/partnersiterequest"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(PartnerSiteRequests())
}

func PartnerSiteRequests() *cobra.Command {
	PartnerSiteRequests := &cobra.Command{
		Use:  "partner-site-requests [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command partner-site-requests\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsPartnerSiteRequestList := files_sdk.PartnerSiteRequestListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Partner Site Requests",
		Long:    `List Partner Site Requests`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsPartnerSiteRequestList
			params.MaxPages = MaxPagesList

			client := partner_site_request.Client{Config: config}
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
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			if len(filterbyList) > 0 {
				listFilter = func(i interface{}) (interface{}, bool, error) {
					matchOk, err := lib.MatchFilter(filterbyList, i)
					return i, matchOk, err
				}
			}
			err = lib.FormatIter(ctx, it, Profile(cmd).Current().SetResourceFormat(cmd, formatList), fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdList.Flags().StringToStringVar(&filterbyList, "filter-by", filterbyList, `Client side filtering: field-name=*.jpg,field-name=?ello`)

	cmdList.Flags().StringVar(&paramsPartnerSiteRequestList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsPartnerSiteRequestList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	PartnerSiteRequests.AddCommand(cmdList)
	var fieldsFindByPairingKey []string
	var formatFindByPairingKey []string
	usePagerFindByPairingKey := true
	paramsPartnerSiteRequestFindByPairingKey := files_sdk.PartnerSiteRequestFindByPairingKeyParams{}

	cmdFindByPairingKey := &cobra.Command{
		Use:   "find-by-pairing-key",
		Short: `Find partner site request by pairing key`,
		Long:  `Find partner site request by pairing key`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := partner_site_request.Client{Config: config}

			var err error
			err = client.FindByPairingKey(paramsPartnerSiteRequestFindByPairingKey, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdFindByPairingKey.Flags().StringVar(&paramsPartnerSiteRequestFindByPairingKey.PairingKey, "pairing-key", "", "Pairing key for the partner site request")

	cmdFindByPairingKey.Flags().StringSliceVar(&fieldsFindByPairingKey, "fields", []string{}, "comma separated list of field names")
	cmdFindByPairingKey.Flags().StringSliceVar(&formatFindByPairingKey, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFindByPairingKey.Flags().BoolVar(&usePagerFindByPairingKey, "use-pager", usePagerFindByPairingKey, "Use $PAGER (.ie less, more, etc)")

	PartnerSiteRequests.AddCommand(cmdFindByPairingKey)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	paramsPartnerSiteRequestCreate := files_sdk.PartnerSiteRequestCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Partner Site Request`,
		Long:  `Create Partner Site Request`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := partner_site_request.Client{Config: config}

			var partnerSiteRequest interface{}
			var err error
			partnerSiteRequest, err = client.Create(paramsPartnerSiteRequestCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), partnerSiteRequest, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().Int64Var(&paramsPartnerSiteRequestCreate.PartnerId, "partner-id", 0, "Partner ID to link with")
	cmdCreate.Flags().StringVar(&paramsPartnerSiteRequestCreate.SiteUrl, "site-url", "", "Site URL to link to")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	PartnerSiteRequests.AddCommand(cmdCreate)
	var fieldsReject []string
	var formatReject []string
	usePagerReject := true
	paramsPartnerSiteRequestReject := files_sdk.PartnerSiteRequestRejectParams{}

	cmdReject := &cobra.Command{
		Use:   "reject",
		Short: `Reject partner site request`,
		Long:  `Reject partner site request`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := partner_site_request.Client{Config: config}

			var err error
			err = client.Reject(paramsPartnerSiteRequestReject, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdReject.Flags().Int64Var(&paramsPartnerSiteRequestReject.Id, "id", 0, "Partner Site Request ID.")

	cmdReject.Flags().StringSliceVar(&fieldsReject, "fields", []string{}, "comma separated list of field names")
	cmdReject.Flags().StringSliceVar(&formatReject, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdReject.Flags().BoolVar(&usePagerReject, "use-pager", usePagerReject, "Use $PAGER (.ie less, more, etc)")

	PartnerSiteRequests.AddCommand(cmdReject)
	var fieldsApprove []string
	var formatApprove []string
	usePagerApprove := true
	paramsPartnerSiteRequestApprove := files_sdk.PartnerSiteRequestApproveParams{}

	cmdApprove := &cobra.Command{
		Use:   "approve",
		Short: `Approve partner site request`,
		Long:  `Approve partner site request`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := partner_site_request.Client{Config: config}

			var err error
			err = client.Approve(paramsPartnerSiteRequestApprove, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdApprove.Flags().Int64Var(&paramsPartnerSiteRequestApprove.Id, "id", 0, "Partner Site Request ID.")

	cmdApprove.Flags().StringSliceVar(&fieldsApprove, "fields", []string{}, "comma separated list of field names")
	cmdApprove.Flags().StringSliceVar(&formatApprove, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdApprove.Flags().BoolVar(&usePagerApprove, "use-pager", usePagerApprove, "Use $PAGER (.ie less, more, etc)")

	PartnerSiteRequests.AddCommand(cmdApprove)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsPartnerSiteRequestDelete := files_sdk.PartnerSiteRequestDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Partner Site Request`,
		Long:  `Delete Partner Site Request`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := partner_site_request.Client{Config: config}

			var err error
			err = client.Delete(paramsPartnerSiteRequestDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsPartnerSiteRequestDelete.Id, "id", 0, "Partner Site Request ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	PartnerSiteRequests.AddCommand(cmdDelete)
	return PartnerSiteRequests
}
