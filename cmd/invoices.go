package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/invoice"
)

var (
	Invoices = &cobra.Command{
		Use:  "invoices [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
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
			client := invoice.Client{Config: files_sdk.GlobalConfig}
			it, err := client.List(params)
			if err != nil {
				lib.ClientError(err)
			}
			err = lib.JsonMarshalIter(it, fieldsList)
			if err != nil {
				lib.ClientError(err)
			}
		},
	}
	cmdList.Flags().StringVarP(&paramsInvoiceList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().IntVarP(&paramsInvoiceList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	Invoices.AddCommand(cmdList)
	var fieldsFind string
	paramsInvoiceFind := files_sdk.InvoiceFindParams{}
	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			client := invoice.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Find(paramsInvoiceFind)
			if err != nil {
				lib.ClientError(err)
			}

			err = lib.JsonMarshal(result, fieldsFind)
			if err != nil {
				lib.ClientError(err)
			}
		},
	}
	cmdFind.Flags().Int64VarP(&paramsInvoiceFind.Id, "id", "i", 0, "Invoice ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	Invoices.AddCommand(cmdFind)
}
