package cmd

import (
	"fmt"
	"time"

	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/transfers"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file"
	flib "github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/spf13/cobra"
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
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	createMkdirParents := true
	createWithRename := true
	paramsFileCreate := files_sdk.FileCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create [path]",
		Short: `Upload file`,
		Long:  `Upload file`,
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := file.Client{Config: config}

			if cmd.Flags().Changed("mkdir-parents") {
				paramsFileCreate.MkdirParents = flib.Bool(createMkdirParents)
			}
			if cmd.Flags().Changed("with-rename") {
				paramsFileCreate.WithRename = flib.Bool(createWithRename)
			}

			if paramsFileCreate.ProvidedMtime.IsZero() {
				paramsFileCreate.ProvidedMtime = nil
			}

			if len(args) > 0 && args[0] != "" {
				paramsFileCreate.Path = args[0]
			}
			var file interface{}
			var err error
			file, err = client.Create(paramsFileCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), file, err, formatCreate, fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().StringVar(&paramsFileCreate.Path, "path", "", "Path to operate on.")
	cmdCreate.Flags().StringVar(&paramsFileCreate.Action, "action", "", "The action to perform.  Can be `append`, `attachment`, `end`, `upload`, `put`, or may not exist")
	cmdCreate.Flags().Int64Var(&paramsFileCreate.Length, "length", 0, "Length of file.")
	cmdCreate.Flags().BoolVar(&createMkdirParents, "mkdir-parents", createMkdirParents, "Create parent directories if they do not exist?")
	cmdCreate.Flags().Int64Var(&paramsFileCreate.Part, "part", 0, "Part if uploading a part.")
	cmdCreate.Flags().Int64Var(&paramsFileCreate.Parts, "parts", 0, "How many parts to fetch?")
	paramsFileCreate.ProvidedMtime = &time.Time{}
	lib.TimeVar(cmdCreate.Flags(), paramsFileCreate.ProvidedMtime, "provided-mtime", "User provided modification time.")
	cmdCreate.Flags().StringVar(&paramsFileCreate.Ref, "ref", "", "")
	cmdCreate.Flags().Int64Var(&paramsFileCreate.Restart, "restart", 0, "File byte offset to restart from.")
	cmdCreate.Flags().Int64Var(&paramsFileCreate.Size, "size", 0, "Size of file.")
	cmdCreate.Flags().StringVar(&paramsFileCreate.Structure, "structure", "", "If copying folder, copy just the structure?")
	cmdCreate.Flags().BoolVar(&createWithRename, "with-rename", createWithRename, "Allow file rename instead of overwrite?")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
json-styles: {raw, pretty}`)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	Files.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	paramsFileUpdate := files_sdk.FileUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update [path]",
		Short: `Update file/folder metadata`,
		Long:  `Update file/folder metadata`,
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := file.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.FileUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			if cmd.Flags().Changed("path") {
				lib.FlagUpdate(cmd, "path", paramsFileUpdate.Path, mapParams)
			}
			if cmd.Flags().Changed("provided-mtime") {
				lib.FlagUpdate(cmd, "provided_mtime", paramsFileUpdate.ProvidedMtime, mapParams)
			}
			if cmd.Flags().Changed("priority-color") {
				lib.FlagUpdate(cmd, "priority_color", paramsFileUpdate.PriorityColor, mapParams)
			}

			if paramsFileUpdate.ProvidedMtime.IsZero() {
				paramsFileUpdate.ProvidedMtime = nil
			}

			if len(args) > 0 && args[0] != "" {
				mapParams["path"] = args[0]
			}
			var file interface{}
			var err error
			file, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), file, err, formatUpdate, fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().StringVar(&paramsFileUpdate.Path, "path", "", "Path to operate on.")
	paramsFileUpdate.ProvidedMtime = &time.Time{}
	lib.TimeVar(cmdUpdate.Flags(), paramsFileUpdate.ProvidedMtime, "provided-mtime", "Modified time of file.")
	cmdUpdate.Flags().StringVar(&paramsFileUpdate.PriorityColor, "priority-color", "", "Priority/Bookmark color of file.")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
json-styles: {raw, pretty}`)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	Files.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	deleteRecursive := true
	paramsFileDelete := files_sdk.FileDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete [path]",
		Short: `Delete file/folder`,
		Long:  `Delete file/folder`,
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := file.Client{Config: config}

			if cmd.Flags().Changed("recursive") {
				paramsFileDelete.Recursive = flib.Bool(deleteRecursive)
			}

			if len(args) > 0 && args[0] != "" {
				paramsFileDelete.Path = args[0]
			}
			var err error
			err = client.Delete(paramsFileDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().StringVar(&paramsFileDelete.Path, "path", "", "Path to operate on.")
	cmdDelete.Flags().BoolVar(&deleteRecursive, "recursive", deleteRecursive, "If true, will recursively delete folers.  Otherwise, will error on non-empty folders.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
json-styles: {raw, pretty}`)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	Files.AddCommand(cmdDelete)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	findWithPreviews := true
	findWithPriorityColor := true
	paramsFileFind := files_sdk.FileFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find [path]",
		Short: `Find file/folder by path`,
		Long:  `Find file/folder by path`,
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := file.Client{Config: config}

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
			file, err = client.Find(paramsFileFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), file, err, formatFind, fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().StringVar(&paramsFileFind.Path, "path", "", "Path to operate on.")
	cmdFind.Flags().StringVar(&paramsFileFind.PreviewSize, "preview-size", "", "Request a preview size.  Can be `small` (default), `large`, `xlarge`, or `pdf`.")
	cmdFind.Flags().BoolVar(&findWithPreviews, "with-previews", findWithPreviews, "Include file preview information?")
	cmdFind.Flags().BoolVar(&findWithPriorityColor, "with-priority-color", findWithPriorityColor, "Include file priority color information?")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
json-styles: {raw, pretty}`)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	Files.AddCommand(cmdFind)
	var fieldsCopy []string
	var formatCopy []string
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
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := file.Client{Config: config}

			if cmd.Flags().Changed("structure") {
				paramsFileCopy.Structure = flib.Bool(copyStructure)
			}

			if len(args) > 0 && args[0] != "" {
				paramsFileCopy.Path = args[0]
			}
			var fileAction interface{}
			var err error
			fileAction, err = client.Copy(paramsFileCopy, files_sdk.WithContext(ctx))
			if err != nil {
				return err
			}
			fileAction, err = transfers.WaitFileMigration(ctx, config, fileAction, blockCopy, noProgressCopy, eventLogCopy, formatCopy, cmd.OutOrStdout())
			if err != nil {
				return err
			}
			return lib.HandleResponse(ctx, Profile(cmd), fileAction, err, formatCopy, fieldsCopy, usePagerCopy, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCopy.Flags().StringVar(&paramsFileCopy.Path, "path", "", "Path to operate on.")
	cmdCopy.Flags().StringVar(&paramsFileCopy.Destination, "destination", "", "Copy destination path.")
	cmdCopy.Flags().BoolVar(&copyStructure, "structure", copyStructure, "Copy structure only?")

	cmdCopy.Flags().StringSliceVar(&fieldsCopy, "fields", []string{}, "comma separated list of field names")
	cmdCopy.Flags().StringSliceVar(&formatCopy, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
json-styles: {raw, pretty}`)
	cmdCopy.Flags().BoolVar(&usePagerCopy, "use-pager", usePagerCopy, "Use $PAGER (.ie less, more, etc)")

	cmdCopy.Flags().BoolVarP(&blockCopy, "block", "b", blockCopy, "Wait on response for async copy with final status")
	cmdCopy.Flags().BoolVar(&noProgressCopy, "no-progress", noProgressCopy, "Don't display progress bars when using block flag")
	cmdCopy.Flags().BoolVar(&eventLogCopy, "event-log", eventLogCopy, "Output full event log for copy when used with block flag")
	Files.AddCommand(cmdCopy)
	var fieldsMove []string
	var formatMove []string
	usePagerMove := true
	var blockMove bool
	var noProgressMove bool
	var eventLogMove bool
	paramsFileMove := files_sdk.FileMoveParams{}

	cmdMove := &cobra.Command{
		Use:   "move [path]",
		Short: `Move file/folder`,
		Long:  `Move file/folder`,
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := file.Client{Config: config}

			if len(args) > 0 && args[0] != "" {
				paramsFileMove.Path = args[0]
			}
			var fileAction interface{}
			var err error
			fileAction, err = client.Move(paramsFileMove, files_sdk.WithContext(ctx))
			if err != nil {
				return err
			}
			fileAction, err = transfers.WaitFileMigration(ctx, config, fileAction, blockMove, noProgressMove, eventLogMove, formatMove, cmd.OutOrStdout())
			if err != nil {
				return err
			}
			return lib.HandleResponse(ctx, Profile(cmd), fileAction, err, formatMove, fieldsMove, usePagerMove, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdMove.Flags().StringVar(&paramsFileMove.Path, "path", "", "Path to operate on.")
	cmdMove.Flags().StringVar(&paramsFileMove.Destination, "destination", "", "Move destination path.")

	cmdMove.Flags().StringSliceVar(&fieldsMove, "fields", []string{}, "comma separated list of field names")
	cmdMove.Flags().StringSliceVar(&formatMove, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
json-styles: {raw, pretty}`)
	cmdMove.Flags().BoolVar(&usePagerMove, "use-pager", usePagerMove, "Use $PAGER (.ie less, more, etc)")

	cmdMove.Flags().BoolVarP(&blockMove, "block", "b", blockMove, "Wait on response for async move with final status")
	cmdMove.Flags().BoolVar(&noProgressMove, "no-progress", noProgressMove, "Don't display progress bars when using block flag")
	cmdMove.Flags().BoolVar(&eventLogMove, "event-log", eventLogMove, "Output full event log for move when used with block flag")
	Files.AddCommand(cmdMove)
	var fieldsBeginUpload []string
	var formatBeginUpload []string
	usePagerBeginUpload := true
	beginUploadMkdirParents := true
	beginUploadWithRename := true
	paramsFileBeginUpload := files_sdk.FileBeginUploadParams{}

	cmdBeginUpload := &cobra.Command{
		Use:   "begin-upload [path]",
		Short: `Begin file upload`,
		Long:  `Begin file upload`,
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := file.Client{Config: config}

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
			fileUploadPartCollection, err = client.BeginUpload(paramsFileBeginUpload, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), fileUploadPartCollection, err, formatBeginUpload, fieldsBeginUpload, usePagerBeginUpload, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
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

	cmdBeginUpload.Flags().StringSliceVar(&fieldsBeginUpload, "fields", []string{}, "comma separated list of field names")
	cmdBeginUpload.Flags().StringSliceVar(&formatBeginUpload, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
json-styles: {raw, pretty}`)
	cmdBeginUpload.Flags().BoolVar(&usePagerBeginUpload, "use-pager", usePagerBeginUpload, "Use $PAGER (.ie less, more, etc)")

	Files.AddCommand(cmdBeginUpload)
	return Files
}
