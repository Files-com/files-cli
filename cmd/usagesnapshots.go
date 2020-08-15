package cmd

import "github.com/spf13/cobra"
import (
	"fmt"
	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/usagesnapshot"
	"os"
)

var (
	_ = files_sdk.Config{}
	_ = usage_snapshot.Client{}
	_ = lib.OnlyFields
	_ = fmt.Println
	_ = os.Exit
)

var (
	UsageSnapshots = &cobra.Command{
		Use:  "usage-snapshots [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func UsageSnapshotsInit() {
	var fieldsList string
	paramsUsageSnapshotList := files_sdk.UsageSnapshotListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsUsageSnapshotList
			params.MaxPages = MaxPagesList
			it := usage_snapshot.List(params)

			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().IntVarP(&paramsUsageSnapshotList.Page, "page", "p", 0, "List Usage Snapshots")
	cmdList.Flags().IntVarP(&paramsUsageSnapshotList.PerPage, "per-page", "e", 0, "List Usage Snapshots")
	cmdList.Flags().StringVarP(&paramsUsageSnapshotList.Action, "action", "a", "", "List Usage Snapshots")
	cmdList.Flags().StringVarP(&paramsUsageSnapshotList.Cursor, "cursor", "c", "", "List Usage Snapshots")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "f", "", "comma separated list of field names to include in response")
	UsageSnapshots.AddCommand(cmdList)
}
