package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	"fmt"
	"os"

	file_action "github.com/Files-com/files-sdk-go/fileaction"
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
			client := file_action.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Copy(paramsFileActionCopy)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = lib.JsonMarshal(result, fieldsCopy)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	cmdCopy.Flags().StringVarP(&paramsFileActionCopy.Path, "path", "p", "", "Path to operate on.")
	cmdCopy.Flags().StringVarP(&paramsFileActionCopy.Destination, "destination", "d", "", "Copy destination path.")

	cmdCopy.Flags().StringVarP(&fieldsCopy, "fields", "", "", "comma separated list of field names")
	FileActions.AddCommand(cmdCopy)
	var fieldsMove string
	paramsFileActionMove := files_sdk.FileActionMoveParams{}
	cmdMove := &cobra.Command{
		Use: "move",
		Run: func(cmd *cobra.Command, args []string) {
			client := file_action.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Move(paramsFileActionMove)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = lib.JsonMarshal(result, fieldsMove)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	cmdMove.Flags().StringVarP(&paramsFileActionMove.Path, "path", "p", "", "Path to operate on.")
	cmdMove.Flags().StringVarP(&paramsFileActionMove.Destination, "destination", "d", "", "Move destination path.")

	cmdMove.Flags().StringVarP(&fieldsMove, "fields", "", "", "comma separated list of field names")
	FileActions.AddCommand(cmdMove)
	var fieldsBeginUpload string
	paramsFileActionBeginUpload := files_sdk.FileActionBeginUploadParams{}
	cmdBeginUpload := &cobra.Command{
		Use: "begin-upload",
		Run: func(cmd *cobra.Command, args []string) {
			client := file_action.Client{Config: files_sdk.GlobalConfig}
			result, err := client.BeginUpload(paramsFileActionBeginUpload)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = lib.JsonMarshal(result, fieldsBeginUpload)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	cmdBeginUpload.Flags().StringVarP(&paramsFileActionBeginUpload.Path, "path", "t", "", "Path to operate on.")
	cmdBeginUpload.Flags().IntVarP(&paramsFileActionBeginUpload.Part, "part", "p", 0, "Part if uploading a part.")
	cmdBeginUpload.Flags().IntVarP(&paramsFileActionBeginUpload.Parts, "parts", "a", 0, "How many parts to fetch?")
	cmdBeginUpload.Flags().StringVarP(&paramsFileActionBeginUpload.Ref, "ref", "r", "", "")
	cmdBeginUpload.Flags().IntVarP(&paramsFileActionBeginUpload.Restart, "restart", "e", 0, "File byte offset to restart from.")

	cmdBeginUpload.Flags().StringVarP(&fieldsBeginUpload, "fields", "", "", "comma separated list of field names")
	FileActions.AddCommand(cmdBeginUpload)
}
