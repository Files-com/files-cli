package cmd

import "github.com/spf13/cobra"
import (
	"fmt"
	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/payment"
	"os"
)

var (
	_ = files_sdk.Config{}
	_ = payment.Client{}
	_ = lib.OnlyFields
	_ = fmt.Println
	_ = os.Exit
)

var (
	Payments = &cobra.Command{
		Use:  "payments [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func PaymentsInit() {
	var fieldsList string
	paramsPaymentList := files_sdk.PaymentListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsPaymentList
			params.MaxPages = MaxPagesList
			it := payment.List(params)

			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().IntVarP(&paramsPaymentList.Page, "page", "p", 0, "List Payments")
	cmdList.Flags().IntVarP(&paramsPaymentList.PerPage, "per-page", "e", 0, "List Payments")
	cmdList.Flags().StringVarP(&paramsPaymentList.Action, "action", "a", "", "List Payments")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "f", "", "comma separated list of field names to include in response")
	Payments.AddCommand(cmdList)
	var fieldsFind string
	paramsPaymentFind := files_sdk.PaymentFindParams{}
	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := payment.Find(paramsPaymentFind)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsFind)
		},
	}
	cmdFind.Flags().IntVarP(&paramsPaymentFind.Id, "id", "i", 0, "Show Payment")
	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "f", "", "comma separated list of field names")
	Payments.AddCommand(cmdFind)
}
