package cmd

import "github.com/spf13/cobra"
import (
	"fmt"
	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/historyexport"
	"os"
)

var (
	_ = files_sdk.Config{}
	_ = history_export.Client{}
	_ = lib.OnlyFields
	_ = fmt.Println
	_ = os.Exit
)

var (
	HistoryExports = &cobra.Command{
		Use:  "history-exports [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func HistoryExportsInit() {
	var fieldsFind string
	paramsHistoryExportFind := files_sdk.HistoryExportFindParams{}
	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := history_export.Find(paramsHistoryExportFind)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsFind)
		},
	}
	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "f", "", "comma separated list of field names")
	HistoryExports.AddCommand(cmdFind)
	var fieldsCreate string
	paramsHistoryExportCreate := files_sdk.HistoryExportCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := history_export.Create(paramsHistoryExportCreate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsCreate)
		},
	}
	cmdCreate.Flags().StringVarP(&paramsHistoryExportCreate.StartAt, "start-at", "", "", "Create History Export")
	cmdCreate.Flags().StringVarP(&paramsHistoryExportCreate.EndAt, "end-at", "e", "", "Create History Export")
	cmdCreate.Flags().StringVarP(&paramsHistoryExportCreate.QueryAction, "query-action", "a", "", "Create History Export")
	cmdCreate.Flags().StringVarP(&paramsHistoryExportCreate.QueryInterface, "query-interface", "n", "", "Create History Export")
	cmdCreate.Flags().StringVarP(&paramsHistoryExportCreate.QueryPath, "query-path", "", "", "Create History Export")
	cmdCreate.Flags().StringVarP(&paramsHistoryExportCreate.QueryFolder, "query-folder", "o", "", "Create History Export")
	cmdCreate.Flags().StringVarP(&paramsHistoryExportCreate.QuerySrc, "query-src", "s", "", "Create History Export")
	cmdCreate.Flags().StringVarP(&paramsHistoryExportCreate.QueryDestination, "query-destination", "d", "", "Create History Export")
	cmdCreate.Flags().StringVarP(&paramsHistoryExportCreate.QueryIp, "query-ip", "p", "", "Create History Export")
	cmdCreate.Flags().StringVarP(&paramsHistoryExportCreate.QueryUsername, "query-username", "", "", "Create History Export")
	cmdCreate.Flags().StringVarP(&paramsHistoryExportCreate.QueryFailureType, "query-failure-type", "t", "", "Create History Export")
	cmdCreate.Flags().StringVarP(&paramsHistoryExportCreate.QueryTargetName, "query-target-name", "", "", "Create History Export")
	cmdCreate.Flags().StringVarP(&paramsHistoryExportCreate.QueryTargetPermission, "query-target-permission", "r", "", "Create History Export")
	cmdCreate.Flags().StringVarP(&paramsHistoryExportCreate.QueryTargetUsername, "query-target-username", "u", "", "Create History Export")
	cmdCreate.Flags().StringVarP(&paramsHistoryExportCreate.QueryTargetPlatform, "query-target-platform", "l", "", "Create History Export")
	cmdCreate.Flags().StringVarP(&paramsHistoryExportCreate.QueryTargetPermissionSet, "query-target-permission-set", "", "", "Create History Export")
	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "f", "", "comma separated list of field names")
	HistoryExports.AddCommand(cmdCreate)
}
