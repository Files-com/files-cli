package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	flib "github.com/Files-com/files-sdk-go/lib"

	"fmt"

	file_action "github.com/Files-com/files-sdk-go/fileaction"
)

var (
	FileActions = &cobra.Command{}
)

func FileActionsInit() {
	FileActions = &cobra.Command{
		Use:  "file-actions [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command file-actions\n\t%v", args[0])
		},
	}
	var fieldsCopy string
	copyStructure := false
	paramsFileActionCopy := files_sdk.FileActionCopyParams{}

	cmdCopy := &cobra.Command{
		Use: "copy [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := file_action.Client{Config: *config}

			if copyStructure {
				paramsFileActionCopy.Structure = flib.Bool(true)
			}

			if len(args) > 0 && args[0] != "" {
				paramsFileActionCopy.Path = args[0]
			}

			result, err := client.Copy(ctx, paramsFileActionCopy)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.JsonMarshal(result, fieldsCopy)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdCopy.Flags().StringVarP(&paramsFileActionCopy.Path, "path", "p", "", "Path to operate on.")
	cmdCopy.Flags().StringVarP(&paramsFileActionCopy.Destination, "destination", "d", "", "Copy destination path.")
	cmdCopy.Flags().BoolVarP(&copyStructure, "structure", "s", copyStructure, "Copy structure only?")

	cmdCopy.Flags().StringVarP(&fieldsCopy, "fields", "", "", "comma separated list of field names")
	FileActions.AddCommand(cmdCopy)
	var fieldsMove string
	paramsFileActionMove := files_sdk.FileActionMoveParams{}

	cmdMove := &cobra.Command{
		Use: "move [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := file_action.Client{Config: *config}

			if len(args) > 0 && args[0] != "" {
				paramsFileActionMove.Path = args[0]
			}

			result, err := client.Move(ctx, paramsFileActionMove)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.JsonMarshal(result, fieldsMove)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdMove.Flags().StringVarP(&paramsFileActionMove.Path, "path", "p", "", "Path to operate on.")
	cmdMove.Flags().StringVarP(&paramsFileActionMove.Destination, "destination", "d", "", "Move destination path.")

	cmdMove.Flags().StringVarP(&fieldsMove, "fields", "", "", "comma separated list of field names")
	FileActions.AddCommand(cmdMove)
	var fieldsBeginUpload string
	beginUploadMkdirParents := false
	beginUploadWithRename := false
	paramsFileActionBeginUpload := files_sdk.FileActionBeginUploadParams{}

	cmdBeginUpload := &cobra.Command{
		Use: "begin-upload [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := file_action.Client{Config: *config}

			if beginUploadMkdirParents {
				paramsFileActionBeginUpload.MkdirParents = flib.Bool(true)
			}
			if beginUploadWithRename {
				paramsFileActionBeginUpload.WithRename = flib.Bool(true)
			}

			if len(args) > 0 && args[0] != "" {
				paramsFileActionBeginUpload.Path = args[0]
			}

			result, err := client.BeginUpload(ctx, paramsFileActionBeginUpload)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.JsonMarshal(result, fieldsBeginUpload)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdBeginUpload.Flags().StringVarP(&paramsFileActionBeginUpload.Path, "path", "t", "", "Path to operate on.")
	cmdBeginUpload.Flags().BoolVarP(&beginUploadMkdirParents, "mkdir-parents", "k", beginUploadMkdirParents, "Create parent directories if they do not exist?")
	cmdBeginUpload.Flags().Int64VarP(&paramsFileActionBeginUpload.Part, "part", "p", 0, "Part if uploading a part.")
	cmdBeginUpload.Flags().Int64VarP(&paramsFileActionBeginUpload.Parts, "parts", "a", 0, "How many parts to fetch?")
	cmdBeginUpload.Flags().StringVarP(&paramsFileActionBeginUpload.Ref, "ref", "r", "", "")
	cmdBeginUpload.Flags().Int64VarP(&paramsFileActionBeginUpload.Restart, "restart", "e", 0, "File byte offset to restart from.")
	cmdBeginUpload.Flags().BoolVarP(&beginUploadWithRename, "with-rename", "w", beginUploadWithRename, "Allow file rename instead of overwrite?")

	cmdBeginUpload.Flags().StringVarP(&fieldsBeginUpload, "fields", "", "", "comma separated list of field names")
	FileActions.AddCommand(cmdBeginUpload)
}
