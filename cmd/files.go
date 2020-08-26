package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	"fmt"
	"os"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/file"
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
	cmdDownload.Flags().StringVarP(&paramsFileDownload.Path, "path", "p", "", "Path to operate on.")
	cmdDownload.Flags().StringVarP(&paramsFileDownload.Action, "action", "a", "", "Can be blank, `redirect` or `stat`.  If set to `stat`, we will return file information but without a download URL, and without logging a download.  If set to `redirect` we will serve a 302 redirect directly to the file.  This is used for integrations with Zapier, and is not recommended for most integrations.")
	cmdDownload.Flags().StringVarP(&paramsFileDownload.PreviewSize, "preview-size", "r", "", "Request a preview size.  Can be `small` (default), `large`, `xlarge`, or `pdf`.")

	cmdDownload.Flags().StringVarP(&fieldsDownload, "fields", "", "", "comma separated list of field names")
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
	cmdCreate.Flags().StringVarP(&paramsFileCreate.Path, "path", "", "", "Path to operate on.")
	cmdCreate.Flags().StringVarP(&paramsFileCreate.Action, "action", "a", "", "The action to perform.  Can be `append`, `attachment`, `end`, `upload`, `put`, or may not exist")
	cmdCreate.Flags().IntVarP(&paramsFileCreate.Length, "length", "l", 0, "Length of file.")
	cmdCreate.Flags().IntVarP(&paramsFileCreate.Part, "part", "p", 0, "Part if uploading a part.")
	cmdCreate.Flags().IntVarP(&paramsFileCreate.Parts, "parts", "r", 0, "How many parts to fetch?")
	lib.TimeVarP(cmdCreate.Flags(), &paramsFileCreate.ProvidedMtime, "provided-mtime", "o")
	cmdCreate.Flags().StringVarP(&paramsFileCreate.Ref, "ref", "f", "", "")
	cmdCreate.Flags().IntVarP(&paramsFileCreate.Restart, "restart", "s", 0, "File byte offset to restart from.")
	cmdCreate.Flags().IntVarP(&paramsFileCreate.Size, "size", "i", 0, "Size of file.")
	cmdCreate.Flags().StringVarP(&paramsFileCreate.Structure, "structure", "u", "", "If copying folder, copy just the structure?")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
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
	cmdUpdate.Flags().StringVarP(&paramsFileUpdate.Path, "path", "p", "", "Path to operate on.")
	lib.TimeVarP(cmdUpdate.Flags(), &paramsFileUpdate.ProvidedMtime, "provided-mtime", "o")
	cmdUpdate.Flags().StringVarP(&paramsFileUpdate.PriorityColor, "priority-color", "r", "", "Priority/Bookmark color of file.")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
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
	cmdDelete.Flags().StringVarP(&paramsFileDelete.Path, "path", "p", "", "Path to operate on.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	Files.AddCommand(cmdDelete)
}
