package cmd

import "github.com/spf13/cobra"
import (
	"fmt"
	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/file"
	"os"
)

var (
	_ = files_sdk.Config{}
	_ = file.Client{}
	_ = lib.OnlyFields
	_ = fmt.Println
	_ = os.Exit
)

var (
	Files = &cobra.Command{
		Use:  "files [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func FilesInit() {
	var fieldsDownload string
	paramsFileDownload := files_sdk.FileDownloadParams{}
	cmdDownload := &cobra.Command{
		Use: "download",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := file.Download(paramsFileDownload)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsDownload)
		},
	}
	cmdDownload.Flags().StringVarP(&paramsFileDownload.Path, "path", "p", "", "Download file")
	cmdDownload.Flags().StringVarP(&paramsFileDownload.Action, "action", "a", "", "Download file")
	cmdDownload.Flags().StringVarP(&paramsFileDownload.PreviewSize, "preview-size", "r", "", "Download file")
	cmdDownload.Flags().StringVarP(&fieldsDownload, "fields", "f", "", "comma separated list of field names")
	Files.AddCommand(cmdDownload)
	var fieldsCreate string
	paramsFileCreate := files_sdk.FileCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := file.Create(paramsFileCreate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsCreate)
		},
	}
	cmdCreate.Flags().StringVarP(&paramsFileCreate.Path, "path", "", "", "Upload file")
	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "f", "", "comma separated list of field names")
	Files.AddCommand(cmdCreate)
	var fieldsUpdate string
	paramsFileUpdate := files_sdk.FileUpdateParams{}
	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := file.Update(paramsFileUpdate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsUpdate)
		},
	}
	cmdUpdate.Flags().StringVarP(&paramsFileUpdate.Path, "path", "p", "", "Update file/folder metadata")
	cmdUpdate.Flags().StringVarP(&paramsFileUpdate.ProvidedMtime, "provided-mtime", "o", "", "Update file/folder metadata")
	cmdUpdate.Flags().StringVarP(&paramsFileUpdate.PriorityColor, "priority-color", "r", "", "Update file/folder metadata")
	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "f", "", "comma separated list of field names")
	Files.AddCommand(cmdUpdate)
	var fieldsDelete string
	paramsFileDelete := files_sdk.FileDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := file.Delete(paramsFileDelete)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsDelete)
		},
	}
	cmdDelete.Flags().StringVarP(&paramsFileDelete.Path, "path", "p", "", "Delete file/folder")
	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "f", "", "comma separated list of field names")
	Files.AddCommand(cmdDelete)
}
