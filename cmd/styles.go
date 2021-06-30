package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

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
		Run:  func(cmd *cobra.Command, args []string) {},
	}
	var fieldsFind string
	paramsStyleFind := files_sdk.StyleFindParams{}

	cmdFind := &cobra.Command{
		Use: "find [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := style.Client{Config: *ctx.GetConfig()}

			if len(args) > 0 && args[0] != "" {
				paramsStyleFind.Path = args[0]
			}

			result, err := client.Find(paramsStyleFind)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsFind)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdFind.Flags().StringVarP(&paramsStyleFind.Path, "path", "p", "", "Style path.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	Styles.AddCommand(cmdFind)
	var fieldsUpdate string
	paramsStyleUpdate := files_sdk.StyleUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use: "update [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := style.Client{Config: *ctx.GetConfig()}

			if len(args) > 0 && args[0] != "" {
				paramsStyleUpdate.Path = args[0]
			}

			result, err := client.Update(paramsStyleUpdate)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsUpdate)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdUpdate.Flags().StringVarP(&paramsStyleUpdate.Path, "path", "p", "", "Style path.")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	Styles.AddCommand(cmdUpdate)
	var fieldsDelete string
	paramsStyleDelete := files_sdk.StyleDeleteParams{}

	cmdDelete := &cobra.Command{
		Use: "delete [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := style.Client{Config: *ctx.GetConfig()}

			if len(args) > 0 && args[0] != "" {
				paramsStyleDelete.Path = args[0]
			}

			result, err := client.Delete(paramsStyleDelete)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsDelete)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdDelete.Flags().StringVarP(&paramsStyleDelete.Path, "path", "p", "", "Style path.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	Styles.AddCommand(cmdDelete)
}
