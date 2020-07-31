package cmd
import "github.com/spf13/cobra"
import (
         "github.com/Files-com/files-cli/lib"
         files_sdk "github.com/Files-com/files-sdk-go"
         "github.com/Files-com/files-sdk-go/invoice"
         "fmt"
         "os"
)

var (
      _ = files_sdk.Config{}
      _ = invoice.Client{}
      _ = lib.OnlyFields
      _ = fmt.Println
      _ = os.Exit
    )

var (
    Invoices = &cobra.Command{
      Use: "invoices [command]",
      Args:  cobra.ExactArgs(1),
      Run: func(cmd *cobra.Command, args []string) {},
    }
)
func InvoicesInit() {
  var fieldsList string
  paramsInvoiceList := files_sdk.InvoiceListParams{}
  var MaxPagesList int
  cmdList := &cobra.Command{
      Use:   "list",
      Short: "list",
      Long:  `list`,
      Args:  cobra.MinimumNArgs(0),
      Run: func(cmd *cobra.Command, args []string) {
        params := paramsInvoiceList
        params.MaxPages = MaxPagesList
        it := invoice.List(params)

        lib.JsonMarshalIter(it, fieldsList)
      },
  }
        cmdList.Flags().IntVarP(&paramsInvoiceList.Page, "page", "p", 0, "List Invoices")
        cmdList.Flags().IntVarP(&paramsInvoiceList.PerPage, "per-page", "e", 0, "List Invoices")
        cmdList.Flags().StringVarP(&paramsInvoiceList.Action, "action", "a", "", "List Invoices")
        cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
        cmdList.Flags().StringVarP(&fieldsList, "fields", "f", "", "comma separated list of field names to include in response")
        Invoices.AddCommand(cmdList)
        var fieldsFind string
        paramsInvoiceFind := files_sdk.InvoiceFindParams{}
        cmdFind := &cobra.Command{
            Use:   "find",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := invoice.Find(paramsInvoiceFind)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsFind)
            },
        }
        cmdFind.Flags().IntVarP(&paramsInvoiceFind.Id, "id", "i", 0, "Show Invoice")
        cmdFind.Flags().StringVarP(&fieldsFind, "fields", "f", "", "comma separated list of field names")
        Invoices.AddCommand(cmdFind)
}
