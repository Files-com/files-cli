package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	"fmt"

	file_migration "github.com/Files-com/files-sdk-go/filemigration"
)

var (
	FileMigrations = &cobra.Command{}
)

func FileMigrationsInit() {
	FileMigrations = &cobra.Command{
		Use:  "file-migrations [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command file-migrations\n\t%v", args[0])
		},
	}
	var fieldsFind string
	var formatFind string
	paramsFileMigrationFind := files_sdk.FileMigrationFindParams{}

	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := file_migration.Client{Config: *config}

			result, err := client.Find(ctx, paramsFileMigrationFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatFind, fieldsFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdFind.Flags().Int64VarP(&paramsFileMigrationFind.Id, "id", "i", 0, "File Migration ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-light")
	FileMigrations.AddCommand(cmdFind)
}
