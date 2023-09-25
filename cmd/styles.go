package cmd

import (
	"fmt"

	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/style"
	"github.com/spf13/cobra"
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
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := style.Client{Config: config}

			if len(args) > 0 && args[0] != "" {
				paramsStyleFind.Path = args[0]
			}
			var style interface{}
			var err error
			style, err = client.Find(paramsStyleFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), style, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().StringVar(&paramsStyleFind.Path, "path", "", "Style path.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
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
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := style.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.StyleUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			if cmd.Flags().Changed("path") {
				lib.FlagUpdate(cmd, "path", paramsStyleUpdate.Path, mapParams)
			}
			if cmd.Flags().Changed("file") {
			}

			if len(args) > 0 && args[0] != "" {
				mapParams["path"] = args[0]
			}
			var style interface{}
			var err error
			style, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), style, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().StringVar(&paramsStyleUpdate.Path, "path", "", "Style path.")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
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
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := style.Client{Config: config}

			if len(args) > 0 && args[0] != "" {
				paramsStyleDelete.Path = args[0]
			}
			var err error
			err = client.Delete(paramsStyleDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().StringVar(&paramsStyleDelete.Path, "path", "", "Style path.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	Styles.AddCommand(cmdDelete)
	return Styles
}
