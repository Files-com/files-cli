package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/transfers"
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/spf13/cobra"

	"fmt"

	"github.com/Files-com/files-sdk-go/v2/file"
	flib "github.com/Files-com/files-sdk-go/v2/lib"
)

func init() {
	RootCmd.AddCommand(Files())
}

func Files() *cobra.Command {
	Files := &cobra.Command{
		Use:  "files [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command files\n\t%v", args[0])
		},
	}
	Files.AddCommand(Download())
	var fieldsCreate string
	var formatCreate string
	usePagerCreate := true
	createMkdirParents := true
	createWithRename := true
	paramsFileCreate := files_sdk.FileCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create [path]",
		Short: `Upload file`,
		Long:  `Upload file`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := file.Client{Config: *config}

			if cmd.Flags().Changed("mkdir-parents") {
				paramsFileCreate.MkdirParents = flib.Bool(createMkdirParents)
			}
			if cmd.Flags().Changed("with-rename") {
				paramsFileCreate.WithRename = flib.Bool(createWithRename)
			}

			if len(args) > 0 && args[0] != "" {
				paramsFileCreate.Path = args[0]
			}
			var file interface{}
			var err error
			file, err = client.Create(ctx, paramsFileCreate)
			lib.HandleResponse(ctx, Profile(cmd), file, err, formatCreate, fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdCreate.Flags().StringVar(&paramsFileCreate.Path, "path", "", "Path to operate on.")
	cmdCreate.Flags().StringVar(&paramsFileCreate.Action, "action", "", "The action to perform.  Can be `append`, `attachment`, `end`, `upload`, `put`, or may not exist")
	cmdCreate.Flags().Int64Var(&paramsFileCreate.Length, "length", 0, "Length of file.")
	cmdCreate.Flags().BoolVar(&createMkdirParents, "mkdir-parents", createMkdirParents, "Create parent directories if they do not exist?")
	cmdCreate.Flags().Int64Var(&paramsFileCreate.Part, "part", 0, "Part if uploading a part.")
	cmdCreate.Flags().Int64Var(&paramsFileCreate.Parts, "parts", 0, "How many parts to fetch?")
	lib.TimeVar(cmdCreate.Flags(), paramsFileCreate.ProvidedMtime, "provided-mtime")
	cmdCreate.Flags().StringVar(&paramsFileCreate.Ref, "ref", "", "")
	cmdCreate.Flags().Int64Var(&paramsFileCreate.Restart, "restart", 0, "File byte offset to restart from.")
	cmdCreate.Flags().Int64Var(&paramsFileCreate.Size, "size", 0, "Size of file.")
	cmdCreate.Flags().StringVar(&paramsFileCreate.Structure, "structure", "", "If copying folder, copy just the structure?")
	cmdCreate.Flags().BoolVar(&createWithRename, "with-rename", createWithRename, "Allow file rename instead of overwrite?")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	Files.AddCommand(cmdCreate)
	var fieldsUpdate string
	var formatUpdate string
	usePagerUpdate := true
	paramsFileUpdate := files_sdk.FileUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update [path]",
		Short: `Update file/folder metadata`,
		Long:  `Update file/folder metadata`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := file.Client{Config: *config}

			if len(args) > 0 && args[0] != "" {
				paramsFileUpdate.Path = args[0]
			}
			var file interface{}
			var err error
			file, err = client.Update(ctx, paramsFileUpdate)
			lib.HandleResponse(ctx, Profile(cmd), file, err, formatUpdate, fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdUpdate.Flags().StringVar(&paramsFileUpdate.Path, "path", "", "Path to operate on.")
	lib.TimeVar(cmdUpdate.Flags(), paramsFileUpdate.ProvidedMtime, "provided-mtime")
	cmdUpdate.Flags().StringVar(&paramsFileUpdate.PriorityColor, "priority-color", "", "Priority/Bookmark color of file.")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	cmdUpdate.Flags().StringVarP(&formatUpdate, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	Files.AddCommand(cmdUpdate)
	var fieldsDelete string
	var formatDelete string
	usePagerDelete := true
	deleteRecursive := true
	paramsFileDelete := files_sdk.FileDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete [path]",
		Short: `Delete file/folder`,
		Long:  `Delete file/folder`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := file.Client{Config: *config}

			if cmd.Flags().Changed("recursive") {
				paramsFileDelete.Recursive = flib.Bool(deleteRecursive)
			}

			if len(args) > 0 && args[0] != "" {
				paramsFileDelete.Path = args[0]
			}
			var err error
			err = client.Delete(ctx, paramsFileDelete)
			if err != nil {
				lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().StringVar(&paramsFileDelete.Path, "path", "", "Path to operate on.")
	cmdDelete.Flags().BoolVar(&deleteRecursive, "recursive", deleteRecursive, "If true, will recursively delete folers.  Otherwise, will error on non-empty folders.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	Files.AddCommand(cmdDelete)
	var fieldsFind string
	var formatFind string
	usePagerFind := true
	findWithPreviews := true
	findWithPriorityColor := true
	paramsFileFind := files_sdk.FileFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find [path]",
		Short: `Find file/folder by path`,
		Long:  `Find file/folder by path`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := file.Client{Config: *config}

			if cmd.Flags().Changed("with-previews") {
				paramsFileFind.WithPreviews = flib.Bool(findWithPreviews)
			}
			if cmd.Flags().Changed("with-priority-color") {
				paramsFileFind.WithPriorityColor = flib.Bool(findWithPriorityColor)
			}

			if len(args) > 0 && args[0] != "" {
				paramsFileFind.Path = args[0]
			}
			var file interface{}
			var err error
			file, err = client.Find(ctx, paramsFileFind)
			lib.HandleResponse(ctx, Profile(cmd), file, err, formatFind, fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdFind.Flags().StringVar(&paramsFileFind.Path, "path", "", "Path to operate on.")
	cmdFind.Flags().StringVar(&paramsFileFind.PreviewSize, "preview-size", "", "Request a preview size.  Can be `small` (default), `large`, `xlarge`, or `pdf`.")
	cmdFind.Flags().BoolVar(&findWithPreviews, "with-previews", findWithPreviews, "Include file preview information?")
	cmdFind.Flags().BoolVar(&findWithPriorityColor, "with-priority-color", findWithPriorityColor, "Include file priority color information?")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	Files.AddCommand(cmdFind)
	var fieldsCopy string
	var formatCopy string
	usePagerCopy := true
	var blockCopy bool
	var noProgressCopy bool
	var eventLogCopy bool
	copyStructure := true
	paramsFileCopy := files_sdk.FileCopyParams{}

	cmdCopy := &cobra.Command{
		Use:   "copy [path]",
		Short: `Copy file/folder`,
		Long:  `Copy file/folder`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := file.Client{Config: *config}

			if cmd.Flags().Changed("structure") {
				paramsFileCopy.Structure = flib.Bool(copyStructure)
			}

			if len(args) > 0 && args[0] != "" {
				paramsFileCopy.Path = args[0]
			}
			var fileAction interface{}
			var err error
			fileAction, err = client.Copy(ctx, paramsFileCopy)
			fileAction, err = transfers.WaitFileMigration(ctx, *config, fileAction, blockCopy, noProgressCopy, eventLogCopy, formatCopy, cmd.OutOrStdout())
			lib.HandleResponse(ctx, Profile(cmd), fileAction, err, formatCopy, fieldsCopy, usePagerCopy, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdCopy.Flags().StringVar(&paramsFileCopy.Path, "path", "", "Path to operate on.")
	cmdCopy.Flags().StringVar(&paramsFileCopy.Destination, "destination", "", "Copy destination path.")
	cmdCopy.Flags().BoolVar(&copyStructure, "structure", copyStructure, "Copy structure only?")

	cmdCopy.Flags().StringVarP(&fieldsCopy, "fields", "", "", "comma separated list of field names")
	cmdCopy.Flags().StringVarP(&formatCopy, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdCopy.Flags().BoolVar(&usePagerCopy, "use-pager", usePagerCopy, "Use $PAGER (.ie less, more, etc)")

	cmdCopy.Flags().BoolVarP(&blockCopy, "block", "b", blockCopy, "Wait on response for async copy with final status")
	cmdCopy.Flags().BoolVar(&noProgressCopy, "no-progress", noProgressCopy, "Don't display progress bars when using block flag")
	cmdCopy.Flags().BoolVar(&eventLogCopy, "event-log", eventLogCopy, "Output full event log for copy when used with block flag")
	Files.AddCommand(cmdCopy)
	var fieldsMove string
	var formatMove string
	usePagerMove := true
	var blockMove bool
	var noProgressMove bool
	var eventLogMove bool
	paramsFileMove := files_sdk.FileMoveParams{}

	cmdMove := &cobra.Command{
		Use:   "move [path]",
		Short: `Move file/folder`,
		Long:  `Move file/folder`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := file.Client{Config: *config}

			if len(args) > 0 && args[0] != "" {
				paramsFileMove.Path = args[0]
			}
			var fileAction interface{}
			var err error
			fileAction, err = client.Move(ctx, paramsFileMove)
			fileAction, err = transfers.WaitFileMigration(ctx, *config, fileAction, blockMove, noProgressMove, eventLogMove, formatMove, cmd.OutOrStdout())
			lib.HandleResponse(ctx, Profile(cmd), fileAction, err, formatMove, fieldsMove, usePagerMove, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdMove.Flags().StringVar(&paramsFileMove.Path, "path", "", "Path to operate on.")
	cmdMove.Flags().StringVar(&paramsFileMove.Destination, "destination", "", "Move destination path.")

	cmdMove.Flags().StringVarP(&fieldsMove, "fields", "", "", "comma separated list of field names")
	cmdMove.Flags().StringVarP(&formatMove, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdMove.Flags().BoolVar(&usePagerMove, "use-pager", usePagerMove, "Use $PAGER (.ie less, more, etc)")

	cmdMove.Flags().BoolVarP(&blockMove, "block", "b", blockMove, "Wait on response for async move with final status")
	cmdMove.Flags().BoolVar(&noProgressMove, "no-progress", noProgressMove, "Don't display progress bars when using block flag")
	cmdMove.Flags().BoolVar(&eventLogMove, "event-log", eventLogMove, "Output full event log for move when used with block flag")
	Files.AddCommand(cmdMove)
	var fieldsBeginUpload string
	var formatBeginUpload string
	usePagerBeginUpload := true
	beginUploadMkdirParents := true
	beginUploadWithRename := true
	paramsFileBeginUpload := files_sdk.FileBeginUploadParams{}

	cmdBeginUpload := &cobra.Command{
		Use:   "begin-upload [path]",
		Short: `Begin file upload`,
		Long:  `Begin file upload`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := file.Client{Config: *config}

			if cmd.Flags().Changed("mkdir-parents") {
				paramsFileBeginUpload.MkdirParents = flib.Bool(beginUploadMkdirParents)
			}
			if cmd.Flags().Changed("with-rename") {
				paramsFileBeginUpload.WithRename = flib.Bool(beginUploadWithRename)
			}

			if len(args) > 0 && args[0] != "" {
				paramsFileBeginUpload.Path = args[0]
			}
			var fileUploadPartCollection interface{}
			var err error
			fileUploadPartCollection, err = client.BeginUpload(ctx, paramsFileBeginUpload)
			lib.HandleResponse(ctx, Profile(cmd), fileUploadPartCollection, err, formatBeginUpload, fieldsBeginUpload, usePagerBeginUpload, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdBeginUpload.Flags().StringVar(&paramsFileBeginUpload.Path, "path", "", "Path to operate on.")
	cmdBeginUpload.Flags().BoolVar(&beginUploadMkdirParents, "mkdir-parents", beginUploadMkdirParents, "Create parent directories if they do not exist?")
	cmdBeginUpload.Flags().Int64Var(&paramsFileBeginUpload.Part, "part", 0, "Part if uploading a part.")
	cmdBeginUpload.Flags().Int64Var(&paramsFileBeginUpload.Parts, "parts", 0, "How many parts to fetch?")
	cmdBeginUpload.Flags().StringVar(&paramsFileBeginUpload.Ref, "ref", "", "")
	cmdBeginUpload.Flags().Int64Var(&paramsFileBeginUpload.Restart, "restart", 0, "File byte offset to restart from.")
	cmdBeginUpload.Flags().Int64Var(&paramsFileBeginUpload.Size, "size", 0, "Total bytes of file being uploaded (include bytes being retained if appending/restarting).")
	cmdBeginUpload.Flags().BoolVar(&beginUploadWithRename, "with-rename", beginUploadWithRename, "Allow file rename instead of overwrite?")

	cmdBeginUpload.Flags().StringVarP(&fieldsBeginUpload, "fields", "", "", "comma separated list of field names")
	cmdBeginUpload.Flags().StringVarP(&formatBeginUpload, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdBeginUpload.Flags().BoolVar(&usePagerBeginUpload, "use-pager", usePagerBeginUpload, "Use $PAGER (.ie less, more, etc)")

	Files.AddCommand(cmdBeginUpload)
	return Files
}
