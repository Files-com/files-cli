package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	"github.com/Files-com/files-sdk-go/file"
	flib "github.com/Files-com/files-sdk-go/lib"
)

var (
	Files = &cobra.Command{}
)

func FilesInit() {
	Files = &cobra.Command{
		Use:  "files [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
	Files.AddCommand(DownloadCmd())
	var fieldsCreate string
	createMkdirParents := false
	createWithRename := false
	paramsFileCreate := files_sdk.FileCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := file.Client{Config: *ctx.GetConfig()}

			if createMkdirParents {
				paramsFileCreate.MkdirParents = flib.Bool(true)
			}
			if createWithRename {
				paramsFileCreate.WithRename = flib.Bool(true)
			}

			if len(args) > 0 && args[0] != "" {
				paramsFileCreate.Path = args[0]
			}
			result, err := client.Create(paramsFileCreate)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsCreate)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdCreate.Flags().StringVarP(&paramsFileCreate.Path, "path", "", "", "Path to operate on.")
	cmdCreate.Flags().StringVarP(&paramsFileCreate.Action, "action", "a", "", "The action to perform.  Can be `append`, `attachment`, `end`, `upload`, `put`, or may not exist")
	cmdCreate.Flags().Int64VarP(&paramsFileCreate.Length, "length", "l", 0, "Length of file.")
	cmdCreate.Flags().BoolVarP(&createMkdirParents, "mkdir-parents", "k", createMkdirParents, "Create parent directories if they do not exist?")
	cmdCreate.Flags().Int64VarP(&paramsFileCreate.Part, "part", "p", 0, "Part if uploading a part.")
	cmdCreate.Flags().Int64VarP(&paramsFileCreate.Parts, "parts", "r", 0, "How many parts to fetch?")
	lib.TimeVarP(cmdCreate.Flags(), &paramsFileCreate.ProvidedMtime, "provided-mtime", "o")
	cmdCreate.Flags().StringVarP(&paramsFileCreate.Ref, "ref", "f", "", "")
	cmdCreate.Flags().Int64VarP(&paramsFileCreate.Restart, "restart", "s", 0, "File byte offset to restart from.")
	cmdCreate.Flags().Int64VarP(&paramsFileCreate.Size, "size", "i", 0, "Size of file.")
	cmdCreate.Flags().StringVarP(&paramsFileCreate.Structure, "structure", "u", "", "If copying folder, copy just the structure?")
	cmdCreate.Flags().BoolVarP(&createWithRename, "with-rename", "w", createWithRename, "Allow file rename instead of overwrite?")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	Files.AddCommand(cmdCreate)
	var fieldsUpdate string
	paramsFileUpdate := files_sdk.FileUpdateParams{}
	cmdUpdate := &cobra.Command{
		Use: "update [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := file.Client{Config: *ctx.GetConfig()}

			if len(args) > 0 && args[0] != "" {
				paramsFileUpdate.Path = args[0]
			}
			result, err := client.Update(paramsFileUpdate)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsUpdate)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdUpdate.Flags().StringVarP(&paramsFileUpdate.Path, "path", "p", "", "Path to operate on.")
	lib.TimeVarP(cmdUpdate.Flags(), &paramsFileUpdate.ProvidedMtime, "provided-mtime", "o")
	cmdUpdate.Flags().StringVarP(&paramsFileUpdate.PriorityColor, "priority-color", "r", "", "Priority/Bookmark color of file.")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	Files.AddCommand(cmdUpdate)
	var fieldsDelete string
	deleteRecursive := false
	paramsFileDelete := files_sdk.FileDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := file.Client{Config: *ctx.GetConfig()}

			if deleteRecursive {
				paramsFileDelete.Recursive = flib.Bool(true)
			}

			if len(args) > 0 && args[0] != "" {
				paramsFileDelete.Path = args[0]
			}
			result, err := client.Delete(paramsFileDelete)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsDelete)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdDelete.Flags().StringVarP(&paramsFileDelete.Path, "path", "p", "", "Path to operate on.")
	cmdDelete.Flags().BoolVarP(&deleteRecursive, "recursive", "r", deleteRecursive, "If true, will recursively delete folers.  Otherwise, will error on non-empty folders.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	Files.AddCommand(cmdDelete)
}
