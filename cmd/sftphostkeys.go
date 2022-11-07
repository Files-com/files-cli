package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	sftp_host_key "github.com/Files-com/files-sdk-go/v2/sftphostkey"
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
	var fieldsList string
	var formatList string
	usePagerList := true
	paramsSftpHostKeyList := files_sdk.SftpHostKeyListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List Sftp Host Keys",
		Long:  `List Sftp Host Keys`,
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsSftpHostKeyList
			params.MaxPages = MaxPagesList

			client := sftp_host_key.Client{Config: *config}
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
				lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(ctx, it, formatList, fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			if err != nil {
				lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}

	cmdList.Flags().StringVar(&paramsSftpHostKeyList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsSftpHostKeyList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVar(&fieldsList, "fields", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVar(&formatList, "format", "table light", `'{format} {style} {direction}' - formats: {json, csv, table}
        table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
        json-styles: {raw, pretty}
        `)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	SftpHostKeys.AddCommand(cmdList)
	var fieldsFind string
	var formatFind string
	usePagerFind := true
	paramsSftpHostKeyFind := files_sdk.SftpHostKeyFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Sftp Host Key`,
		Long:  `Show Sftp Host Key`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := sftp_host_key.Client{Config: *config}

			var sftpHostKey interface{}
			var err error
			sftpHostKey, err = client.Find(ctx, paramsSftpHostKeyFind)
			lib.HandleResponse(ctx, Profile(cmd), sftpHostKey, err, formatFind, fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdFind.Flags().Int64Var(&paramsSftpHostKeyFind.Id, "id", 0, "Sftp Host Key ID.")

	cmdFind.Flags().StringVar(&fieldsFind, "fields", "", "comma separated list of field names")
	cmdFind.Flags().StringVar(&formatFind, "format", "table light", `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	SftpHostKeys.AddCommand(cmdFind)
	var fieldsCreate string
	var formatCreate string
	usePagerCreate := true
	paramsSftpHostKeyCreate := files_sdk.SftpHostKeyCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Sftp Host Key`,
		Long:  `Create Sftp Host Key`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := sftp_host_key.Client{Config: *config}

			var sftpHostKey interface{}
			var err error
			sftpHostKey, err = client.Create(ctx, paramsSftpHostKeyCreate)
			lib.HandleResponse(ctx, Profile(cmd), sftpHostKey, err, formatCreate, fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdCreate.Flags().StringVar(&paramsSftpHostKeyCreate.Name, "name", "", "The friendly name of this SFTP Host Key.")
	cmdCreate.Flags().StringVar(&paramsSftpHostKeyCreate.PrivateKey, "private-key", "", "The private key data.")

	cmdCreate.Flags().StringVar(&fieldsCreate, "fields", "", "comma separated list of field names")
	cmdCreate.Flags().StringVar(&formatCreate, "format", "table light", `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	SftpHostKeys.AddCommand(cmdCreate)
	var fieldsUpdate string
	var formatUpdate string
	usePagerUpdate := true
	paramsSftpHostKeyUpdate := files_sdk.SftpHostKeyUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Sftp Host Key`,
		Long:  `Update Sftp Host Key`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := sftp_host_key.Client{Config: *config}

			var sftpHostKey interface{}
			var err error
			sftpHostKey, err = client.Update(ctx, paramsSftpHostKeyUpdate)
			lib.HandleResponse(ctx, Profile(cmd), sftpHostKey, err, formatUpdate, fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsSftpHostKeyUpdate.Id, "id", 0, "Sftp Host Key ID.")
	cmdUpdate.Flags().StringVar(&paramsSftpHostKeyUpdate.Name, "name", "", "The friendly name of this SFTP Host Key.")
	cmdUpdate.Flags().StringVar(&paramsSftpHostKeyUpdate.PrivateKey, "private-key", "", "The private key data.")

	cmdUpdate.Flags().StringVar(&fieldsUpdate, "fields", "", "comma separated list of field names")
	cmdUpdate.Flags().StringVar(&formatUpdate, "format", "table light", `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	SftpHostKeys.AddCommand(cmdUpdate)
	var fieldsDelete string
	var formatDelete string
	usePagerDelete := true
	paramsSftpHostKeyDelete := files_sdk.SftpHostKeyDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Sftp Host Key`,
		Long:  `Delete Sftp Host Key`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := sftp_host_key.Client{Config: *config}

			var err error
			err = client.Delete(ctx, paramsSftpHostKeyDelete)
			if err != nil {
				lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsSftpHostKeyDelete.Id, "id", 0, "Sftp Host Key ID.")

	cmdDelete.Flags().StringVar(&fieldsDelete, "fields", "", "comma separated list of field names")
	cmdDelete.Flags().StringVar(&formatDelete, "format", "table light", `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	SftpHostKeys.AddCommand(cmdDelete)
	return SftpHostKeys
}
