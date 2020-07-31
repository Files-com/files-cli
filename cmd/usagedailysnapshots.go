package cmd
import "github.com/spf13/cobra"
import (
         "github.com/Files-com/files-cli/lib"
         files_sdk "github.com/Files-com/files-sdk-go"
         "github.com/Files-com/files-sdk-go/usagedailysnapshot"
         "fmt"
         "os"
)

var (
      _ = files_sdk.Config{}
      _ = usage_daily_snapshot.Client{}
      _ = lib.OnlyFields
      _ = fmt.Println
      _ = os.Exit
    )

var (
    UsageDailySnapshots = &cobra.Command{
      Use: "usage-daily-snapshots [command]",
      Args:  cobra.ExactArgs(1),
      Run: func(cmd *cobra.Command, args []string) {},
    }
)
func UsageDailySnapshotsInit() {
  var fieldsList string
  paramsUsageDailySnapshotList := files_sdk.UsageDailySnapshotListParams{}
  var MaxPagesList int
  cmdList := &cobra.Command{
      Use:   "list",
      Short: "list",
      Long:  `list`,
      Args:  cobra.MinimumNArgs(0),
      Run: func(cmd *cobra.Command, args []string) {
        params := paramsUsageDailySnapshotList
        params.MaxPages = MaxPagesList
        it := usage_daily_snapshot.List(params)

        lib.JsonMarshalIter(it, fieldsList)
      },
  }
        cmdList.Flags().IntVarP(&paramsUsageDailySnapshotList.Page, "page", "p", 0, "List Usage Daily Snapshots")
        cmdList.Flags().IntVarP(&paramsUsageDailySnapshotList.PerPage, "per-page", "r", 0, "List Usage Daily Snapshots")
        cmdList.Flags().StringVarP(&paramsUsageDailySnapshotList.Action, "action", "a", "", "List Usage Daily Snapshots")
        cmdList.Flags().StringVarP(&paramsUsageDailySnapshotList.Cursor, "cursor", "c", "", "List Usage Daily Snapshots")
        cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
        cmdList.Flags().StringVarP(&fieldsList, "fields", "f", "", "comma separated list of field names to include in response")
        UsageDailySnapshots.AddCommand(cmdList)
}
