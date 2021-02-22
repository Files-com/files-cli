package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

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
	Files.AddCommand(DownloadCmd())
	var fieldsCreate string
	paramsFileCreate := files_sdk.FileCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			client := file.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Create(paramsFileCreate)
			if err != nil {
				lib.ClientError(err)
			}

			err = lib.JsonMarshal(result, fieldsCreate)
			if err != nil {
				lib.ClientError(err)
			}
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
			client := file.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Update(paramsFileUpdate)
			if err != nil {
				lib.ClientError(err)
			}

			err = lib.JsonMarshal(result, fieldsUpdate)
			if err != nil {
				lib.ClientError(err)
			}
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
			client := file.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Delete(paramsFileDelete)
			if err != nil {
				lib.ClientError(err)
			}

			err = lib.JsonMarshal(result, fieldsDelete)
			if err != nil {
				lib.ClientError(err)
			}
		},
	}
	cmdDelete.Flags().StringVarP(&paramsFileDelete.Path, "path", "p", "", "Path to operate on.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	Files.AddCommand(cmdDelete)
}
