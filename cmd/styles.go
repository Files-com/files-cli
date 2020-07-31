package cmd
import "github.com/spf13/cobra"
import (
         "github.com/Files-com/files-cli/lib"
         files_sdk "github.com/Files-com/files-sdk-go"
         "github.com/Files-com/files-sdk-go/style"
         "fmt"
         "os"
)

var (
      _ = files_sdk.Config{}
      _ = style.Client{}
      _ = lib.OnlyFields
      _ = fmt.Println
      _ = os.Exit
    )

var (
    Styles = &cobra.Command{
      Use: "styles [command]",
      Args:  cobra.ExactArgs(1),
      Run: func(cmd *cobra.Command, args []string) {},
    }
)
func StylesInit() {
        var fieldsFind string
        paramsStyleFind := files_sdk.StyleFindParams{}
        cmdFind := &cobra.Command{
            Use:   "find",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := style.Find(paramsStyleFind)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsFind)
            },
        }
        cmdFind.Flags().StringVarP(&paramsStyleFind.Path, "path", "p", "", "Show Style")
        cmdFind.Flags().StringVarP(&fieldsFind, "fields", "f", "", "comma separated list of field names")
        Styles.AddCommand(cmdFind)
        var fieldsUpdate string
        paramsStyleUpdate := files_sdk.StyleUpdateParams{}
        cmdUpdate := &cobra.Command{
            Use:   "update",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := style.Update(paramsStyleUpdate)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsUpdate)
            },
        }
        cmdUpdate.Flags().StringVarP(&paramsStyleUpdate.Path, "path", "p", "", "Update Style")
        cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "f", "", "comma separated list of field names")
        Styles.AddCommand(cmdUpdate)
        var fieldsDelete string
        paramsStyleDelete := files_sdk.StyleDeleteParams{}
        cmdDelete := &cobra.Command{
            Use:   "delete",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := style.Delete(paramsStyleDelete)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsDelete)
            },
        }
        cmdDelete.Flags().StringVarP(&paramsStyleDelete.Path, "path", "p", "", "Delete Style")
        cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "f", "", "comma separated list of field names")
        Styles.AddCommand(cmdDelete)
}
