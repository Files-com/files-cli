package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	public_key "github.com/Files-com/files-sdk-go/v2/publickey"
)

func init() {
	RootCmd.AddCommand(PublicKeys())
}

func PublicKeys() *cobra.Command {
	PublicKeys := &cobra.Command{
		Use:  "public-keys [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command public-keys\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	paramsPublicKeyList := files_sdk.PublicKeyListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List Public Keys",
		Long:  `List Public Keys`,
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsPublicKeyList
			params.MaxPages = MaxPagesList

			client := public_key.Client{Config: *config}
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
				return lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(ctx, it, formatList, fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			if err != nil {
				return lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}

	cmdList.Flags().Int64Var(&paramsPublicKeyList.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVar(&paramsPublicKeyList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsPublicKeyList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
        table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
        json-styles: {raw, pretty}
        `)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	PublicKeys.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsPublicKeyFind := files_sdk.PublicKeyFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Public Key`,
		Long:  `Show Public Key`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := public_key.Client{Config: *config}

			var publicKey interface{}
			var err error
			publicKey, err = client.Find(ctx, paramsPublicKeyFind)
			lib.HandleResponse(ctx, Profile(cmd), publicKey, err, formatFind, fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdFind.Flags().Int64Var(&paramsPublicKeyFind.Id, "id", 0, "Public Key ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	PublicKeys.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	paramsPublicKeyCreate := files_sdk.PublicKeyCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Public Key`,
		Long:  `Create Public Key`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := public_key.Client{Config: *config}

			var publicKey interface{}
			var err error
			publicKey, err = client.Create(ctx, paramsPublicKeyCreate)
			lib.HandleResponse(ctx, Profile(cmd), publicKey, err, formatCreate, fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdCreate.Flags().Int64Var(&paramsPublicKeyCreate.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().StringVar(&paramsPublicKeyCreate.Title, "title", "", "Internal reference for key.")
	cmdCreate.Flags().StringVar(&paramsPublicKeyCreate.PublicKey, "public-key", "", "Actual contents of SSH key.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	PublicKeys.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	paramsPublicKeyUpdate := files_sdk.PublicKeyUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Public Key`,
		Long:  `Update Public Key`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := public_key.Client{Config: *config}

			var publicKey interface{}
			var err error
			publicKey, err = client.Update(ctx, paramsPublicKeyUpdate)
			lib.HandleResponse(ctx, Profile(cmd), publicKey, err, formatUpdate, fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsPublicKeyUpdate.Id, "id", 0, "Public Key ID.")
	cmdUpdate.Flags().StringVar(&paramsPublicKeyUpdate.Title, "title", "", "Internal reference for key.")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	PublicKeys.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsPublicKeyDelete := files_sdk.PublicKeyDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Public Key`,
		Long:  `Delete Public Key`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := public_key.Client{Config: *config}

			var err error
			err = client.Delete(ctx, paramsPublicKeyDelete)
			if err != nil {
				return lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsPublicKeyDelete.Id, "id", 0, "Public Key ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	PublicKeys.AddCommand(cmdDelete)
	return PublicKeys
}
