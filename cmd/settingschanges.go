package cmd

import "github.com/spf13/cobra"
import (
	"fmt"
	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/settingschange"
	"os"
)

var (
	_ = files_sdk.Config{}
	_ = settings_change.Client{}
	_ = lib.OnlyFields
	_ = fmt.Println
	_ = os.Exit
)

var (
	SettingsChanges = &cobra.Command{
		Use:  "settings-changes [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func SettingsChangesInit() {
	var fieldsList string
	paramsSettingsChangeList := files_sdk.SettingsChangeListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsSettingsChangeList
			params.MaxPages = MaxPagesList
			it := settings_change.List(params)

			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().IntVarP(&paramsSettingsChangeList.Page, "page", "p", 0, "List Settings Changes")
	cmdList.Flags().IntVarP(&paramsSettingsChangeList.PerPage, "per-page", "r", 0, "List Settings Changes")
	cmdList.Flags().StringVarP(&paramsSettingsChangeList.Action, "action", "a", "", "List Settings Changes")
	cmdList.Flags().StringVarP(&paramsSettingsChangeList.Cursor, "cursor", "c", "", "List Settings Changes")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "f", "", "comma separated list of field names to include in response")
	SettingsChanges.AddCommand(cmdList)
}
