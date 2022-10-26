package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	file_migration "github.com/Files-com/files-sdk-go/v2/filemigration"
)

func init() {
	RootCmd.AddCommand(FileMigrations())
}

func FileMigrations() *cobra.Command {
	FileMigrations := &cobra.Command{
		Use:  "file-migrations [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command file-migrations\n\t%v", args[0])
		},
	}
	var fieldsFind string
	var formatFind string
	usePagerFind := true
	paramsFileMigrationFind := files_sdk.FileMigrationFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show File Migration`,
		Long:  `Show File Migration`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := file_migration.Client{Config: *config}

			var fileMigration interface{}
			var err error
			fileMigration, err = client.Find(ctx, paramsFileMigrationFind)
			lib.HandleResponse(ctx, Profile(cmd), fileMigration, err, formatFind, fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdFind.Flags().Int64Var(&paramsFileMigrationFind.Id, "id", 0, "File Migration ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	FileMigrations.AddCommand(cmdFind)
	return FileMigrations
}
