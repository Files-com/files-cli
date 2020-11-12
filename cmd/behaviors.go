package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	"fmt"
	"os"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/behavior"
)

var (
	Behaviors = &cobra.Command{
		Use:  "behaviors [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func BehaviorsInit() {
	var fieldsList string
	paramsBehaviorList := files_sdk.BehaviorListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsBehaviorList
			params.MaxPages = MaxPagesList
			it, err := behavior.List(params)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().StringVarP(&paramsBehaviorList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().IntVarP(&paramsBehaviorList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVarP(&paramsBehaviorList.Behavior, "behavior", "b", "", "If set, only shows folder behaviors matching this behavior type.")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	Behaviors.AddCommand(cmdList)
	var fieldsFind string
	paramsBehaviorFind := files_sdk.BehaviorFindParams{}
	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := behavior.Find(paramsBehaviorFind)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsFind)
		},
	}
	cmdFind.Flags().Int64VarP(&paramsBehaviorFind.Id, "id", "i", 0, "Behavior ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	Behaviors.AddCommand(cmdFind)
	var fieldsListFor string
	paramsBehaviorListFor := files_sdk.BehaviorListForParams{}
	var MaxPagesListFor int
	cmdListFor := &cobra.Command{
		Use:   "list-for [path]",
		Short: "list-for",
		Long:  `list-for`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsBehaviorListFor
			params.MaxPages = MaxPagesListFor
			if len(args) > 0 && args[0] != "" {
				params.Path = args[0]
			}
			it, err := behavior.ListFor(params)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			lib.JsonMarshalIter(it, fieldsListFor)
		},
	}
	cmdListFor.Flags().StringVarP(&paramsBehaviorListFor.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdListFor.Flags().IntVarP(&paramsBehaviorListFor.PerPage, "per-page", "r", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdListFor.Flags().StringVarP(&paramsBehaviorListFor.Path, "path", "p", "", "Path to operate on.")
	cmdListFor.Flags().StringVarP(&paramsBehaviorListFor.Recursive, "recursive", "u", "", "Show behaviors above this path?")
	cmdListFor.Flags().StringVarP(&paramsBehaviorListFor.Behavior, "behavior", "b", "", "DEPRECATED: If set only shows folder behaviors matching this behavior type. Use `filter[behavior]` instead.")
	cmdListFor.Flags().IntVarP(&MaxPagesListFor, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdListFor.Flags().StringVarP(&fieldsListFor, "fields", "", "", "comma separated list of field names to include in response")
	Behaviors.AddCommand(cmdListFor)
	var fieldsCreate string
	paramsBehaviorCreate := files_sdk.BehaviorCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := behavior.Create(paramsBehaviorCreate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsCreate)
		},
	}
	cmdCreate.Flags().StringVarP(&paramsBehaviorCreate.Value, "value", "v", "", "The value of the folder behavior.  Can be a integer, array, or hash depending on the type of folder behavior.")
	cmdCreate.Flags().StringVarP(&paramsBehaviorCreate.Path, "path", "p", "", "Folder behaviors path.")
	cmdCreate.Flags().StringVarP(&paramsBehaviorCreate.Behavior, "behavior", "b", "", "Behavior type.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	Behaviors.AddCommand(cmdCreate)
	var fieldsWebhookTest string
	paramsBehaviorWebhookTest := files_sdk.BehaviorWebhookTestParams{}
	cmdWebhookTest := &cobra.Command{
		Use: "webhook-test",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := behavior.WebhookTest(paramsBehaviorWebhookTest)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsWebhookTest)
		},
	}
	cmdWebhookTest.Flags().StringVarP(&paramsBehaviorWebhookTest.Url, "url", "u", "", "URL for testing the webhook.")
	cmdWebhookTest.Flags().StringVarP(&paramsBehaviorWebhookTest.Method, "method", "t", "", "HTTP method(GET or POST).")
	cmdWebhookTest.Flags().StringVarP(&paramsBehaviorWebhookTest.Encoding, "encoding", "e", "", "HTTP encoding method.  Can be JSON, XML, or RAW (form data).")
	cmdWebhookTest.Flags().StringVarP(&paramsBehaviorWebhookTest.Action, "action", "a", "", "action for test body")

	cmdWebhookTest.Flags().StringVarP(&fieldsWebhookTest, "fields", "", "", "comma separated list of field names")
	Behaviors.AddCommand(cmdWebhookTest)
	var fieldsUpdate string
	paramsBehaviorUpdate := files_sdk.BehaviorUpdateParams{}
	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := behavior.Update(paramsBehaviorUpdate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsUpdate)
		},
	}
	cmdUpdate.Flags().Int64VarP(&paramsBehaviorUpdate.Id, "id", "i", 0, "Behavior ID.")
	cmdUpdate.Flags().StringVarP(&paramsBehaviorUpdate.Value, "value", "v", "", "The value of the folder behavior.  Can be a integer, array, or hash depending on the type of folder behavior.")
	cmdUpdate.Flags().StringVarP(&paramsBehaviorUpdate.Behavior, "behavior", "b", "", "Behavior type.")
	cmdUpdate.Flags().StringVarP(&paramsBehaviorUpdate.Path, "path", "p", "", "Folder behaviors path.")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	Behaviors.AddCommand(cmdUpdate)
	var fieldsDelete string
	paramsBehaviorDelete := files_sdk.BehaviorDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := behavior.Delete(paramsBehaviorDelete)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsDelete)
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsBehaviorDelete.Id, "id", "i", 0, "Behavior ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	Behaviors.AddCommand(cmdDelete)
}
