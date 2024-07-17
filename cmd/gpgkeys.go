package cmd

import (
	"fmt"

	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	gpg_key "github.com/Files-com/files-sdk-go/v3/gpgkey"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(GpgKeys())
}

func GpgKeys() *cobra.Command {
	GpgKeys := &cobra.Command{
		Use:  "gpg-keys [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command gpg-keys\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsGpgKeyList := files_sdk.GpgKeyListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List GPG Keys",
		Long:    `List GPG Keys`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsGpgKeyList
			params.MaxPages = MaxPagesList

			client := gpg_key.Client{Config: config}
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

	cmdList.Flags().Int64Var(&paramsGpgKeyList.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVar(&paramsGpgKeyList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsGpgKeyList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVar(&paramsGpgKeyList.Action, "action", "", "")
	cmdList.Flags().Int64Var(&paramsGpgKeyList.Page, "page", 0, "")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	GpgKeys.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsGpgKeyFind := files_sdk.GpgKeyFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show GPG Key`,
		Long:  `Show GPG Key`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := gpg_key.Client{Config: config}

			var gpgKey interface{}
			var err error
			gpgKey, err = client.Find(paramsGpgKeyFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), gpgKey, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsGpgKeyFind.Id, "id", 0, "Gpg Key ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	GpgKeys.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	paramsGpgKeyCreate := files_sdk.GpgKeyCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create GPG Key`,
		Long:  `Create GPG Key`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := gpg_key.Client{Config: config}

			var gpgKey interface{}
			var err error
			gpgKey, err = client.Create(paramsGpgKeyCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), gpgKey, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().Int64Var(&paramsGpgKeyCreate.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().StringVar(&paramsGpgKeyCreate.PublicKey, "public-key", "", "Your GPG public key")
	cmdCreate.Flags().StringVar(&paramsGpgKeyCreate.PrivateKey, "private-key", "", "Your GPG private key.")
	cmdCreate.Flags().StringVar(&paramsGpgKeyCreate.PrivateKeyPassword, "private-key-password", "", "Your GPG private key password. Only required for password protected keys.")
	cmdCreate.Flags().StringVar(&paramsGpgKeyCreate.Name, "name", "", "Your GPG key name.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	GpgKeys.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	paramsGpgKeyUpdate := files_sdk.GpgKeyUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update GPG Key`,
		Long:  `Update GPG Key`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := gpg_key.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.GpgKeyUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsGpgKeyUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("public-key") {
				lib.FlagUpdate(cmd, "public_key", paramsGpgKeyUpdate.PublicKey, mapParams)
			}
			if cmd.Flags().Changed("private-key") {
				lib.FlagUpdate(cmd, "private_key", paramsGpgKeyUpdate.PrivateKey, mapParams)
			}
			if cmd.Flags().Changed("private-key-password") {
				lib.FlagUpdate(cmd, "private_key_password", paramsGpgKeyUpdate.PrivateKeyPassword, mapParams)
			}
			if cmd.Flags().Changed("name") {
				lib.FlagUpdate(cmd, "name", paramsGpgKeyUpdate.Name, mapParams)
			}

			var gpgKey interface{}
			var err error
			gpgKey, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), gpgKey, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsGpgKeyUpdate.Id, "id", 0, "Gpg Key ID.")
	cmdUpdate.Flags().StringVar(&paramsGpgKeyUpdate.PublicKey, "public-key", "", "Your GPG public key")
	cmdUpdate.Flags().StringVar(&paramsGpgKeyUpdate.PrivateKey, "private-key", "", "Your GPG private key.")
	cmdUpdate.Flags().StringVar(&paramsGpgKeyUpdate.PrivateKeyPassword, "private-key-password", "", "Your GPG private key password. Only required for password protected keys.")
	cmdUpdate.Flags().StringVar(&paramsGpgKeyUpdate.Name, "name", "", "Your GPG key name.")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	GpgKeys.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsGpgKeyDelete := files_sdk.GpgKeyDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete GPG Key`,
		Long:  `Delete GPG Key`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := gpg_key.Client{Config: config}

			var err error
			err = client.Delete(paramsGpgKeyDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsGpgKeyDelete.Id, "id", 0, "Gpg Key ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	GpgKeys.AddCommand(cmdDelete)
	return GpgKeys
}
