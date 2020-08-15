package cmd

import "github.com/spf13/cobra"
import (
	"fmt"
	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/ssostrategy"
	"os"
)

var (
	_ = files_sdk.Config{}
	_ = sso_strategy.Client{}
	_ = lib.OnlyFields
	_ = fmt.Println
	_ = os.Exit
)

var (
	SsoStrategies = &cobra.Command{
		Use:  "sso-strategies [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func SsoStrategiesInit() {
	var fieldsList string
	paramsSsoStrategyList := files_sdk.SsoStrategyListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsSsoStrategyList
			params.MaxPages = MaxPagesList
			it := sso_strategy.List(params)

			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().IntVarP(&paramsSsoStrategyList.Page, "page", "p", 0, "List Sso Strategies")
	cmdList.Flags().IntVarP(&paramsSsoStrategyList.PerPage, "per-page", "e", 0, "List Sso Strategies")
	cmdList.Flags().StringVarP(&paramsSsoStrategyList.Action, "action", "a", "", "List Sso Strategies")
	cmdList.Flags().StringVarP(&paramsSsoStrategyList.Cursor, "cursor", "c", "", "List Sso Strategies")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "f", "", "comma separated list of field names to include in response")
	SsoStrategies.AddCommand(cmdList)
	var fieldsFind string
	paramsSsoStrategyFind := files_sdk.SsoStrategyFindParams{}
	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := sso_strategy.Find(paramsSsoStrategyFind)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsFind)
		},
	}
	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "f", "", "comma separated list of field names")
	SsoStrategies.AddCommand(cmdFind)
}
