package cmd

import (
	"fmt"

	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	sftp_host_key "github.com/Files-com/files-sdk-go/v3/sftphostkey"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(SftpHostKeys())
}

func SftpHostKeys() *cobra.Command {
	SftpHostKeys := &cobra.Command{
		Use:  "sftp-host-keys [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command sftp-host-keys\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsSftpHostKeyList := files_sdk.SftpHostKeyListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Sftp Host Keys",
		Long:    `List Sftp Host Keys`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsSftpHostKeyList
			params.MaxPages = MaxPagesList

			client := sftp_host_key.Client{Config: config}
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

	cmdList.Flags().StringVar(&paramsSftpHostKeyList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsSftpHostKeyList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVar(&paramsSftpHostKeyList.Action, "action", "", "")
	cmdList.Flags().Int64Var(&paramsSftpHostKeyList.Page, "page", 0, "")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	SftpHostKeys.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsSftpHostKeyFind := files_sdk.SftpHostKeyFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Sftp Host Key`,
		Long:  `Show Sftp Host Key`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := sftp_host_key.Client{Config: config}

			var sftpHostKey interface{}
			var err error
			sftpHostKey, err = client.Find(paramsSftpHostKeyFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), sftpHostKey, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsSftpHostKeyFind.Id, "id", 0, "Sftp Host Key ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	SftpHostKeys.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	paramsSftpHostKeyCreate := files_sdk.SftpHostKeyCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Sftp Host Key`,
		Long:  `Create Sftp Host Key`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := sftp_host_key.Client{Config: config}

			var sftpHostKey interface{}
			var err error
			sftpHostKey, err = client.Create(paramsSftpHostKeyCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), sftpHostKey, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().StringVar(&paramsSftpHostKeyCreate.Name, "name", "", "The friendly name of this SFTP Host Key.")
	cmdCreate.Flags().StringVar(&paramsSftpHostKeyCreate.PrivateKey, "private-key", "", "The private key data.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	SftpHostKeys.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	paramsSftpHostKeyUpdate := files_sdk.SftpHostKeyUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Sftp Host Key`,
		Long:  `Update Sftp Host Key`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := sftp_host_key.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.SftpHostKeyUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsSftpHostKeyUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("name") {
				lib.FlagUpdate(cmd, "name", paramsSftpHostKeyUpdate.Name, mapParams)
			}
			if cmd.Flags().Changed("private-key") {
				lib.FlagUpdate(cmd, "private_key", paramsSftpHostKeyUpdate.PrivateKey, mapParams)
			}

			var sftpHostKey interface{}
			var err error
			sftpHostKey, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), sftpHostKey, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsSftpHostKeyUpdate.Id, "id", 0, "Sftp Host Key ID.")
	cmdUpdate.Flags().StringVar(&paramsSftpHostKeyUpdate.Name, "name", "", "The friendly name of this SFTP Host Key.")
	cmdUpdate.Flags().StringVar(&paramsSftpHostKeyUpdate.PrivateKey, "private-key", "", "The private key data.")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	SftpHostKeys.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsSftpHostKeyDelete := files_sdk.SftpHostKeyDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Sftp Host Key`,
		Long:  `Delete Sftp Host Key`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := sftp_host_key.Client{Config: config}

			var err error
			err = client.Delete(paramsSftpHostKeyDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsSftpHostKeyDelete.Id, "id", 0, "Sftp Host Key ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	SftpHostKeys.AddCommand(cmdDelete)
	return SftpHostKeys
}
