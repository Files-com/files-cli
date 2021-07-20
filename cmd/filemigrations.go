package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	file_migration "github.com/Files-com/files-sdk-go/filemigration"
)

var (
	FileMigrations = &cobra.Command{}
)

func FileMigrationsInit() {
	FileMigrations = &cobra.Command{
		Use:  "file-migrations [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
	var fieldsFind string
	paramsFileMigrationFind := files_sdk.FileMigrationFindParams{}

	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := file_migration.Client{Config: *ctx.GetConfig()}

			result, err := client.Find(paramsFileMigrationFind)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsFind)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdFind.Flags().Int64VarP(&paramsFileMigrationFind.Id, "id", "i", 0, "File Migration ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	FileMigrations.AddCommand(cmdFind)
}
