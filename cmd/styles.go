package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/style"
)

var (
	Styles = &cobra.Command{}
)

func StylesInit() {
	Styles = &cobra.Command{
		Use:  "styles [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command styles\n\t%v", args[0])
		},
	}
	var fieldsFind string
	var formatFind string
	paramsStyleFind := files_sdk.StyleFindParams{}

	cmdFind := &cobra.Command{
		Use: "find [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := style.Client{Config: *config}

			if len(args) > 0 && args[0] != "" {
				paramsStyleFind.Path = args[0]
			}

			result, err := client.Find(ctx, paramsStyleFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatFind, fieldsFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdFind.Flags().StringVarP(&paramsStyleFind.Path, "path", "p", "", "Style path.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-light")
	Styles.AddCommand(cmdFind)
	var fieldsUpdate string
	var formatUpdate string
	paramsStyleUpdate := files_sdk.StyleUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use: "update [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := style.Client{Config: *config}

			if len(args) > 0 && args[0] != "" {
				paramsStyleUpdate.Path = args[0]
			}

			result, err := client.Update(ctx, paramsStyleUpdate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatUpdate, fieldsUpdate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdUpdate.Flags().StringVarP(&paramsStyleUpdate.Path, "path", "p", "", "Style path.")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	cmdUpdate.Flags().StringVarP(&formatUpdate, "format", "", "table", "json, csv, table, table-dark, table-light")
	Styles.AddCommand(cmdUpdate)
	var fieldsDelete string
	var formatDelete string
	paramsStyleDelete := files_sdk.StyleDeleteParams{}

	cmdDelete := &cobra.Command{
		Use: "delete [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := style.Client{Config: *config}

			if len(args) > 0 && args[0] != "" {
				paramsStyleDelete.Path = args[0]
			}

			result, err := client.Delete(ctx, paramsStyleDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatDelete, fieldsDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdDelete.Flags().StringVarP(&paramsStyleDelete.Path, "path", "p", "", "Style path.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-light")
	Styles.AddCommand(cmdDelete)
}
