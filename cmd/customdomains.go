package cmd

import (
	"fmt"
	"reflect"

	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	custom_domain "github.com/Files-com/files-sdk-go/v3/customdomain"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(CustomDomains())
}

func CustomDomains() *cobra.Command {
	CustomDomains := &cobra.Command{
		Use:  "custom-domains [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command custom-domains\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsCustomDomainList := files_sdk.CustomDomainListParams{}
	var MaxPagesList int64
	var listSortByArgs string

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Custom Domains",
		Long:    `List Custom Domains`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsCustomDomainList
			params.MaxPages = MaxPagesList

			parsedListSortBy, parseListSortByErr := lib.ParseAPIListSortFlag("sort-by", listSortByArgs)
			if parseListSortByErr != nil {
				return parseListSortByErr
			}
			if parsedListSortBy != nil {
				params.SortBy = parsedListSortBy
			}

			client := custom_domain.Client{Config: config}
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

	cmdList.Flags().StringToStringVar(&filterbyList, "filter-by", filterbyList, "Client-side wildcard filtering, for example field-name=*.jpg or field-name=?ello")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter-by", "field=pattern")
	cmdList.Flags().StringVar(&listSortByArgs, "sort-by", "", "Sort custom domains by field in ascending or descending order.")
	lib.SetFlagDisplayType(cmdList.Flags(), "sort-by", "field=asc|desc")

	cmdList.Flags().StringVar(&paramsCustomDomainList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsCustomDomainList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	CustomDomains.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsCustomDomainFind := files_sdk.CustomDomainFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Custom Domain`,
		Long:  `Show Custom Domain`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := custom_domain.Client{Config: config}

			var customDomain interface{}
			var err error
			customDomain, err = client.Find(paramsCustomDomainFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), customDomain, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsCustomDomainFind.Id, "id", 0, "Custom Domain ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	CustomDomains.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	paramsCustomDomainCreate := files_sdk.CustomDomainCreateParams{}
	CustomDomainCreateDestination := ""

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Custom Domain`,
		Long:  `Create Custom Domain`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := custom_domain.Client{Config: config}

			var CustomDomainCreateDestinationErr error
			paramsCustomDomainCreate.Destination, CustomDomainCreateDestinationErr = lib.FetchKey("destination", paramsCustomDomainCreate.Destination.Enum(), CustomDomainCreateDestination)
			if CustomDomainCreateDestination != "" && CustomDomainCreateDestinationErr != nil {
				return CustomDomainCreateDestinationErr
			}

			var customDomain interface{}
			var err error
			customDomain, err = client.Create(paramsCustomDomainCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), customDomain, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().StringVar(&CustomDomainCreateDestination, "destination", "", fmt.Sprintf("Where this custom domain routes. Can be `site_alias`, `public_hosting`, `s3_endpoint`, or `unassigned` (not routing traffic). Set to `unassigned` automatically when a bound `public_hosting` folder behavior is deleted, and can be set manually via the API for any reason. %v", reflect.ValueOf(paramsCustomDomainCreate.Destination.Enum()).MapKeys()))
	cmdCreate.Flags().Int64Var(&paramsCustomDomainCreate.FolderBehaviorId, "folder-behavior-id", 0, "Public Hosting behavior ID when this domain routes to a specific Public Hosting behavior.  Preserved as historical context when `destination` becomes `unassigned`.")
	cmdCreate.Flags().Int64Var(&paramsCustomDomainCreate.SslCertificateId, "ssl-certificate-id", 0, "Current SSL certificate ID.")
	cmdCreate.Flags().StringVar(&paramsCustomDomainCreate.Domain, "domain", "", "Customer-owned domain name.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	CustomDomains.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	paramsCustomDomainUpdate := files_sdk.CustomDomainUpdateParams{}
	CustomDomainUpdateDestination := ""

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Custom Domain`,
		Long:  `Update Custom Domain`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := custom_domain.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.CustomDomainUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			var CustomDomainUpdateDestinationErr error
			paramsCustomDomainUpdate.Destination, CustomDomainUpdateDestinationErr = lib.FetchKey("destination", paramsCustomDomainUpdate.Destination.Enum(), CustomDomainUpdateDestination)
			if CustomDomainUpdateDestination != "" && CustomDomainUpdateDestinationErr != nil {
				return CustomDomainUpdateDestinationErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsCustomDomainUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("destination") {
				lib.FlagUpdate(cmd, "destination", paramsCustomDomainUpdate.Destination, mapParams)
			}
			if cmd.Flags().Changed("folder-behavior-id") {
				lib.FlagUpdate(cmd, "folder_behavior_id", paramsCustomDomainUpdate.FolderBehaviorId, mapParams)
			}
			if cmd.Flags().Changed("ssl-certificate-id") {
				lib.FlagUpdate(cmd, "ssl_certificate_id", paramsCustomDomainUpdate.SslCertificateId, mapParams)
			}
			if cmd.Flags().Changed("domain") {
				lib.FlagUpdate(cmd, "domain", paramsCustomDomainUpdate.Domain, mapParams)
			}

			var customDomain interface{}
			var err error
			customDomain, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), customDomain, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsCustomDomainUpdate.Id, "id", 0, "Custom Domain ID.")
	cmdUpdate.Flags().StringVar(&CustomDomainUpdateDestination, "destination", "", fmt.Sprintf("Where this custom domain routes. Can be `site_alias`, `public_hosting`, `s3_endpoint`, or `unassigned` (not routing traffic). Set to `unassigned` automatically when a bound `public_hosting` folder behavior is deleted, and can be set manually via the API for any reason. %v", reflect.ValueOf(paramsCustomDomainUpdate.Destination.Enum()).MapKeys()))
	cmdUpdate.Flags().Int64Var(&paramsCustomDomainUpdate.FolderBehaviorId, "folder-behavior-id", 0, "Public Hosting behavior ID when this domain routes to a specific Public Hosting behavior.  Preserved as historical context when `destination` becomes `unassigned`.")
	cmdUpdate.Flags().Int64Var(&paramsCustomDomainUpdate.SslCertificateId, "ssl-certificate-id", 0, "Current SSL certificate ID.")
	cmdUpdate.Flags().StringVar(&paramsCustomDomainUpdate.Domain, "domain", "", "Customer-owned domain name.")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	CustomDomains.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsCustomDomainDelete := files_sdk.CustomDomainDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Custom Domain`,
		Long:  `Delete Custom Domain`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := custom_domain.Client{Config: config}

			var err error
			err = client.Delete(paramsCustomDomainDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsCustomDomainDelete.Id, "id", 0, "Custom Domain ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	CustomDomains.AddCommand(cmdDelete)
	return CustomDomains
}
