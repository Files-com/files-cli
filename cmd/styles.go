package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/style"
)

var (
	Styles = &cobra.Command{
		Use:  "styles [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func StylesInit() {
	var fieldsFind string
	paramsStyleFind := files_sdk.StyleFindParams{}
	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			client := style.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Find(paramsStyleFind)
			if err != nil {
				lib.ClientError(err)
			}

			err = lib.JsonMarshal(result, fieldsFind)
			if err != nil {
				lib.ClientError(err)
			}
		},
	}
	cmdFind.Flags().StringVarP(&paramsStyleFind.Path, "path", "p", "", "Style path.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	Styles.AddCommand(cmdFind)
	var fieldsUpdate string
	paramsStyleUpdate := files_sdk.StyleUpdateParams{}
	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			client := style.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Update(paramsStyleUpdate)
			if err != nil {
				lib.ClientError(err)
			}

			err = lib.JsonMarshal(result, fieldsUpdate)
			if err != nil {
				lib.ClientError(err)
			}
		},
	}
	cmdUpdate.Flags().StringVarP(&paramsStyleUpdate.Path, "path", "p", "", "Style path.")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	Styles.AddCommand(cmdUpdate)
	var fieldsDelete string
	paramsStyleDelete := files_sdk.StyleDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			client := style.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Delete(paramsStyleDelete)
			if err != nil {
				lib.ClientError(err)
			}

			err = lib.JsonMarshal(result, fieldsDelete)
			if err != nil {
				lib.ClientError(err)
			}
		},
	}
	cmdDelete.Flags().StringVarP(&paramsStyleDelete.Path, "path", "p", "", "Style path.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	Styles.AddCommand(cmdDelete)
}
