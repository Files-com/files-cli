package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	"github.com/Files-com/files-sdk-go/v2/file"
	flib "github.com/Files-com/files-sdk-go/v2/lib"
)

var (
	Files = &cobra.Command{}
)

func FilesInit() {
	Files = &cobra.Command{
		Use:  "files [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command files\n\t%v", args[0])
		},
	}
	Files.AddCommand(DownloadCmd())
	var fieldsCreate string
	var formatCreate string
	createMkdirParents := false
	createWithRename := false
	paramsFileCreate := files_sdk.FileCreateParams{}

	cmdCreate := &cobra.Command{
		Use: "create [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := file.Client{Config: *config}

			if createMkdirParents {
				paramsFileCreate.MkdirParents = flib.Bool(true)
			}
			if createWithRename {
				paramsFileCreate.WithRename = flib.Bool(true)
			}

			if len(args) > 0 && args[0] != "" {
				paramsFileCreate.Path = args[0]
			}

			result, err := client.Create(ctx, paramsFileCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatCreate, fieldsCreate)
			if err != nil {
				lib.ClientError(ctx, err)
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
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Files.AddCommand(cmdCreate)
	var fieldsUpdate string
	var formatUpdate string
	paramsFileUpdate := files_sdk.FileUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use: "update [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := file.Client{Config: *config}

			if len(args) > 0 && args[0] != "" {
				paramsFileUpdate.Path = args[0]
			}

			result, err := client.Update(ctx, paramsFileUpdate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatUpdate, fieldsUpdate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdUpdate.Flags().StringVarP(&paramsFileUpdate.Path, "path", "p", "", "Path to operate on.")
	lib.TimeVarP(cmdUpdate.Flags(), &paramsFileUpdate.ProvidedMtime, "provided-mtime", "o")
	cmdUpdate.Flags().StringVarP(&paramsFileUpdate.PriorityColor, "priority-color", "r", "", "Priority/Bookmark color of file.")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	cmdUpdate.Flags().StringVarP(&formatUpdate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Files.AddCommand(cmdUpdate)
	var fieldsDelete string
	var formatDelete string
	deleteRecursive := false
	paramsFileDelete := files_sdk.FileDeleteParams{}

	cmdDelete := &cobra.Command{
		Use: "delete [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := file.Client{Config: *config}

			if deleteRecursive {
				paramsFileDelete.Recursive = flib.Bool(true)
			}

			if len(args) > 0 && args[0] != "" {
				paramsFileDelete.Path = args[0]
			}

			result, err := client.Delete(ctx, paramsFileDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatDelete, fieldsDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdDelete.Flags().StringVarP(&paramsFileDelete.Path, "path", "p", "", "Path to operate on.")
	cmdDelete.Flags().BoolVarP(&deleteRecursive, "recursive", "r", deleteRecursive, "If true, will recursively delete folers.  Otherwise, will error on non-empty folders.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Files.AddCommand(cmdDelete)
	var fieldsFind string
	var formatFind string
	findWithPreviews := false
	findWithPriorityColor := false
	paramsFileFind := files_sdk.FileFindParams{}

	cmdFind := &cobra.Command{
		Use: "find [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := file.Client{Config: *config}

			if findWithPreviews {
				paramsFileFind.WithPreviews = flib.Bool(true)
			}
			if findWithPriorityColor {
				paramsFileFind.WithPriorityColor = flib.Bool(true)
			}

			if len(args) > 0 && args[0] != "" {
				paramsFileFind.Path = args[0]
			}

			result, err := client.Find(ctx, paramsFileFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatFind, fieldsFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdFind.Flags().StringVarP(&paramsFileFind.Path, "path", "p", "", "Path to operate on.")
	cmdFind.Flags().StringVarP(&paramsFileFind.PreviewSize, "preview-size", "r", "", "Request a preview size.  Can be `small` (default), `large`, `xlarge`, or `pdf`.")
	cmdFind.Flags().BoolVarP(&findWithPreviews, "with-previews", "e", findWithPreviews, "Include file preview information?")
	cmdFind.Flags().BoolVarP(&findWithPriorityColor, "with-priority-color", "c", findWithPriorityColor, "Include file priority color information?")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Files.AddCommand(cmdFind)
	var fieldsCopy string
	var formatCopy string
	copyStructure := false
	paramsFileCopy := files_sdk.FileCopyParams{}

	cmdCopy := &cobra.Command{
		Use: "copy [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := file.Client{Config: *config}

			if copyStructure {
				paramsFileCopy.Structure = flib.Bool(true)
			}

			if len(args) > 0 && args[0] != "" {
				paramsFileCopy.Path = args[0]
			}

			result, err := client.Copy(ctx, paramsFileCopy)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatCopy, fieldsCopy)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdCopy.Flags().StringVarP(&paramsFileCopy.Path, "path", "p", "", "Path to operate on.")
	cmdCopy.Flags().StringVarP(&paramsFileCopy.Destination, "destination", "d", "", "Copy destination path.")
	cmdCopy.Flags().BoolVarP(&copyStructure, "structure", "s", copyStructure, "Copy structure only?")

	cmdCopy.Flags().StringVarP(&fieldsCopy, "fields", "", "", "comma separated list of field names")
	cmdCopy.Flags().StringVarP(&formatCopy, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Files.AddCommand(cmdCopy)
	var fieldsMove string
	var formatMove string
	paramsFileMove := files_sdk.FileMoveParams{}

	cmdMove := &cobra.Command{
		Use: "move [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := file.Client{Config: *config}

			if len(args) > 0 && args[0] != "" {
				paramsFileMove.Path = args[0]
			}

			result, err := client.Move(ctx, paramsFileMove)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatMove, fieldsMove)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdMove.Flags().StringVarP(&paramsFileMove.Path, "path", "p", "", "Path to operate on.")
	cmdMove.Flags().StringVarP(&paramsFileMove.Destination, "destination", "d", "", "Move destination path.")

	cmdMove.Flags().StringVarP(&fieldsMove, "fields", "", "", "comma separated list of field names")
	cmdMove.Flags().StringVarP(&formatMove, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Files.AddCommand(cmdMove)
	var fieldsBeginUpload string
	var formatBeginUpload string
	beginUploadMkdirParents := false
	beginUploadWithRename := false
	paramsFileBeginUpload := files_sdk.FileBeginUploadParams{}

	cmdBeginUpload := &cobra.Command{
		Use: "begin-upload [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := file.Client{Config: *config}

			if beginUploadMkdirParents {
				paramsFileBeginUpload.MkdirParents = flib.Bool(true)
			}
			if beginUploadWithRename {
				paramsFileBeginUpload.WithRename = flib.Bool(true)
			}

			if len(args) > 0 && args[0] != "" {
				paramsFileBeginUpload.Path = args[0]
			}

			result, err := client.BeginUpload(ctx, paramsFileBeginUpload)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatBeginUpload, fieldsBeginUpload)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdBeginUpload.Flags().StringVarP(&paramsFileBeginUpload.Path, "path", "t", "", "Path to operate on.")
	cmdBeginUpload.Flags().BoolVarP(&beginUploadMkdirParents, "mkdir-parents", "k", beginUploadMkdirParents, "Create parent directories if they do not exist?")
	cmdBeginUpload.Flags().Int64VarP(&paramsFileBeginUpload.Part, "part", "p", 0, "Part if uploading a part.")
	cmdBeginUpload.Flags().Int64VarP(&paramsFileBeginUpload.Parts, "parts", "a", 0, "How many parts to fetch?")
	cmdBeginUpload.Flags().StringVarP(&paramsFileBeginUpload.Ref, "ref", "r", "", "")
	cmdBeginUpload.Flags().Int64VarP(&paramsFileBeginUpload.Restart, "restart", "e", 0, "File byte offset to restart from.")
	cmdBeginUpload.Flags().BoolVarP(&beginUploadWithRename, "with-rename", "w", beginUploadWithRename, "Allow file rename instead of overwrite?")

	cmdBeginUpload.Flags().StringVarP(&fieldsBeginUpload, "fields", "", "", "comma separated list of field names")
	cmdBeginUpload.Flags().StringVarP(&formatBeginUpload, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Files.AddCommand(cmdBeginUpload)
}
