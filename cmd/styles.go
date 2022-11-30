package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/style"
)

func init() {
	RootCmd.AddCommand(Styles())
}

func Styles() *cobra.Command {
	Styles := &cobra.Command{
		Use:  "styles [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command styles\n\t%v", args[0])
		},
	}
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsStyleFind := files_sdk.StyleFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find [path]",
		Short: `Show Style`,
		Long:  `Show Style`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := style.Client{Config: *config}

			if len(args) > 0 && args[0] != "" {
				paramsStyleFind.Path = args[0]
			}
			var style interface{}
			var err error
			style, err = client.Find(ctx, paramsStyleFind)
			lib.HandleResponse(ctx, Profile(cmd), style, err, formatFind, fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdFind.Flags().StringVar(&paramsStyleFind.Path, "path", "", "Style path.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	Styles.AddCommand(cmdFind)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	paramsStyleUpdate := files_sdk.StyleUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update [path]",
		Short: `Update Style`,
		Long:  `Update Style`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := style.Client{Config: *config}

			if len(args) > 0 && args[0] != "" {
				paramsStyleUpdate.Path = args[0]
			}
			var style interface{}
			var err error
			style, err = client.Update(ctx, paramsStyleUpdate)
			lib.HandleResponse(ctx, Profile(cmd), style, err, formatUpdate, fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdUpdate.Flags().StringVar(&paramsStyleUpdate.Path, "path", "", "Style path.")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	Styles.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsStyleDelete := files_sdk.StyleDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete [path]",
		Short: `Delete Style`,
		Long:  `Delete Style`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := style.Client{Config: *config}

			if len(args) > 0 && args[0] != "" {
				paramsStyleDelete.Path = args[0]
			}
			var err error
			err = client.Delete(ctx, paramsStyleDelete)
			if err != nil {
				return lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().StringVar(&paramsStyleDelete.Path, "path", "", "Style path.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	Styles.AddCommand(cmdDelete)
	return Styles
}
