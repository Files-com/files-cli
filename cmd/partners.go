package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	flib "github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/Files-com/files-sdk-go/v3/partner"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(Partners())
}

func Partners() *cobra.Command {
	Partners := &cobra.Command{
		Use:  "partners [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command partners\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsPartnerList := files_sdk.PartnerListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Partners",
		Long:    `List Partners`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsPartnerList
			params.MaxPages = MaxPagesList

			client := partner.Client{Config: config}
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

	cmdList.Flags().StringVar(&paramsPartnerList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsPartnerList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	Partners.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsPartnerFind := files_sdk.PartnerFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Partner`,
		Long:  `Show Partner`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := partner.Client{Config: config}

			var partner interface{}
			var err error
			partner, err = client.Find(paramsPartnerFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), partner, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsPartnerFind.Id, "id", 0, "Partner ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	Partners.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	createAllowBypassing2faPolicies := true
	createAllowCredentialChanges := true
	createAllowUserCreation := true
	paramsPartnerCreate := files_sdk.PartnerCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Partner`,
		Long:  `Create Partner`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := partner.Client{Config: config}

			if cmd.Flags().Changed("allow-bypassing-2fa-policies") {
				paramsPartnerCreate.AllowBypassing2faPolicies = flib.Bool(createAllowBypassing2faPolicies)
			}
			if cmd.Flags().Changed("allow-credential-changes") {
				paramsPartnerCreate.AllowCredentialChanges = flib.Bool(createAllowCredentialChanges)
			}
			if cmd.Flags().Changed("allow-user-creation") {
				paramsPartnerCreate.AllowUserCreation = flib.Bool(createAllowUserCreation)
			}

			var partner interface{}
			var err error
			partner, err = client.Create(paramsPartnerCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), partner, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().BoolVar(&createAllowBypassing2faPolicies, "allow-bypassing-2fa-policies", createAllowBypassing2faPolicies, "Allow users created under this Partner to bypass Two-Factor Authentication policies.")
	cmdCreate.Flags().BoolVar(&createAllowCredentialChanges, "allow-credential-changes", createAllowCredentialChanges, "Allow Partner Admins to change or reset credentials for users belonging to this Partner.")
	cmdCreate.Flags().BoolVar(&createAllowUserCreation, "allow-user-creation", createAllowUserCreation, "Allow Partner Admins to create users.")
	cmdCreate.Flags().StringVar(&paramsPartnerCreate.Name, "name", "", "The name of the Partner.")
	cmdCreate.Flags().StringVar(&paramsPartnerCreate.Notes, "notes", "", "Notes about this Partner.")
	cmdCreate.Flags().StringVar(&paramsPartnerCreate.RootFolder, "root-folder", "", "The root folder path for this Partner.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	Partners.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	updateAllowBypassing2faPolicies := true
	updateAllowCredentialChanges := true
	updateAllowUserCreation := true
	paramsPartnerUpdate := files_sdk.PartnerUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Partner`,
		Long:  `Update Partner`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := partner.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.PartnerUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsPartnerUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("allow-bypassing-2fa-policies") {
				mapParams["allow_bypassing_2fa_policies"] = updateAllowBypassing2faPolicies
			}
			if cmd.Flags().Changed("allow-credential-changes") {
				mapParams["allow_credential_changes"] = updateAllowCredentialChanges
			}
			if cmd.Flags().Changed("allow-user-creation") {
				mapParams["allow_user_creation"] = updateAllowUserCreation
			}
			if cmd.Flags().Changed("name") {
				lib.FlagUpdate(cmd, "name", paramsPartnerUpdate.Name, mapParams)
			}
			if cmd.Flags().Changed("notes") {
				lib.FlagUpdate(cmd, "notes", paramsPartnerUpdate.Notes, mapParams)
			}
			if cmd.Flags().Changed("root-folder") {
				lib.FlagUpdate(cmd, "root_folder", paramsPartnerUpdate.RootFolder, mapParams)
			}

			var partner interface{}
			var err error
			partner, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), partner, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsPartnerUpdate.Id, "id", 0, "Partner ID.")
	cmdUpdate.Flags().BoolVar(&updateAllowBypassing2faPolicies, "allow-bypassing-2fa-policies", updateAllowBypassing2faPolicies, "Allow users created under this Partner to bypass Two-Factor Authentication policies.")
	cmdUpdate.Flags().BoolVar(&updateAllowCredentialChanges, "allow-credential-changes", updateAllowCredentialChanges, "Allow Partner Admins to change or reset credentials for users belonging to this Partner.")
	cmdUpdate.Flags().BoolVar(&updateAllowUserCreation, "allow-user-creation", updateAllowUserCreation, "Allow Partner Admins to create users.")
	cmdUpdate.Flags().StringVar(&paramsPartnerUpdate.Name, "name", "", "The name of the Partner.")
	cmdUpdate.Flags().StringVar(&paramsPartnerUpdate.Notes, "notes", "", "Notes about this Partner.")
	cmdUpdate.Flags().StringVar(&paramsPartnerUpdate.RootFolder, "root-folder", "", "The root folder path for this Partner.")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	Partners.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsPartnerDelete := files_sdk.PartnerDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Partner`,
		Long:  `Delete Partner`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := partner.Client{Config: config}

			var err error
			err = client.Delete(paramsPartnerDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsPartnerDelete.Id, "id", 0, "Partner ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	Partners.AddCommand(cmdDelete)
	return Partners
}
