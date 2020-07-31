package cmd

import "github.com/spf13/cobra"
import (
	"fmt"
	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/bandwidthsnapshot"
	"os"
)

var (
	_ = files_sdk.Config{}
	_ = bandwidth_snapshot.Client{}
	_ = lib.OnlyFields
	_ = fmt.Println
	_ = os.Exit
)

var (
	BandwidthSnapshots = &cobra.Command{
		Use:  "bandwidth-snapshots [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func BandwidthSnapshotsInit() {
	var fieldsList string
	paramsBandwidthSnapshotList := files_sdk.BandwidthSnapshotListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsBandwidthSnapshotList
			params.MaxPages = MaxPagesList
			it := bandwidth_snapshot.List(params)

			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().IntVarP(&paramsBandwidthSnapshotList.Page, "page", "p", 0, "List Bandwidth Snapshots")
	cmdList.Flags().IntVarP(&paramsBandwidthSnapshotList.PerPage, "per-page", "e", 0, "List Bandwidth Snapshots")
	cmdList.Flags().StringVarP(&paramsBandwidthSnapshotList.Action, "action", "a", "", "List Bandwidth Snapshots")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "f", "", "comma separated list of field names to include in response")
	BandwidthSnapshots.AddCommand(cmdList)
}
