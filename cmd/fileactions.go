package cmd

import "github.com/spf13/cobra"
import (
	"fmt"
	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/fileaction"
	"os"
)

var (
	_ = files_sdk.Config{}
	_ = file_action.Client{}
	_ = lib.OnlyFields
	_ = fmt.Println
	_ = os.Exit
)

var (
	FileActions = &cobra.Command{
		Use:  "file-actions [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func FileActionsInit() {
	var fieldsCopy string
	paramsFileActionCopy := files_sdk.FileActionCopyParams{}
	cmdCopy := &cobra.Command{
		Use: "copy",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := file_action.Copy(paramsFileActionCopy)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsCopy)
		},
	}
	cmdCopy.Flags().StringVarP(&paramsFileActionCopy.Path, "path", "p", "", "Copy file/folder")
	cmdCopy.Flags().StringVarP(&paramsFileActionCopy.Destination, "destination", "d", "", "Copy file/folder")
	cmdCopy.Flags().StringVarP(&fieldsCopy, "fields", "f", "", "comma separated list of field names")
	FileActions.AddCommand(cmdCopy)
	var fieldsMove string
	paramsFileActionMove := files_sdk.FileActionMoveParams{}
	cmdMove := &cobra.Command{
		Use: "move",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := file_action.Move(paramsFileActionMove)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsMove)
		},
	}
	cmdMove.Flags().StringVarP(&paramsFileActionMove.Path, "path", "p", "", "Move file/folder")
	cmdMove.Flags().StringVarP(&paramsFileActionMove.Destination, "destination", "d", "", "Move file/folder")
	cmdMove.Flags().StringVarP(&fieldsMove, "fields", "f", "", "comma separated list of field names")
	FileActions.AddCommand(cmdMove)
	var fieldsBeginUpload string
	paramsFileActionBeginUpload := files_sdk.FileActionBeginUploadParams{}
	cmdBeginUpload := &cobra.Command{
		Use: "begin-upload",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := file_action.BeginUpload(paramsFileActionBeginUpload)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsBeginUpload)
		},
	}
	cmdBeginUpload.Flags().StringVarP(&paramsFileActionBeginUpload.Path, "path", "t", "", "Begin file upload")
	cmdBeginUpload.Flags().IntVarP(&paramsFileActionBeginUpload.Part, "part", "p", 0, "Begin file upload")
	cmdBeginUpload.Flags().IntVarP(&paramsFileActionBeginUpload.Parts, "parts", "a", 0, "Begin file upload")
	cmdBeginUpload.Flags().StringVarP(&paramsFileActionBeginUpload.Ref, "ref", "r", "", "Begin file upload")
	cmdBeginUpload.Flags().IntVarP(&paramsFileActionBeginUpload.Restart, "restart", "e", 0, "Begin file upload")
	cmdBeginUpload.Flags().StringVarP(&fieldsBeginUpload, "fields", "f", "", "comma separated list of field names")
	FileActions.AddCommand(cmdBeginUpload)
}
