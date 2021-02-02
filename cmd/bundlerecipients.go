package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	"fmt"
	"os"

	bundle_recipient "github.com/Files-com/files-sdk-go/bundlerecipient"
)

var (
	BundleRecipients = &cobra.Command{
		Use:  "bundle-recipients [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func BundleRecipientsInit() {
	var fieldsList string
	paramsBundleRecipientList := files_sdk.BundleRecipientListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsBundleRecipientList
			params.MaxPages = MaxPagesList
			client := bundle_recipient.Client{Config: files_sdk.GlobalConfig}
			it, err := client.List(params)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			err = lib.JsonMarshalIter(it, fieldsList)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	cmdList.Flags().Int64VarP(&paramsBundleRecipientList.UserId, "user-id", "u", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVarP(&paramsBundleRecipientList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().IntVarP(&paramsBundleRecipientList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().Int64VarP(&paramsBundleRecipientList.BundleId, "bundle-id", "b", 0, "List recipients for the bundle with this ID.")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	BundleRecipients.AddCommand(cmdList)
	var fieldsCreate string
	paramsBundleRecipientCreate := files_sdk.BundleRecipientCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			client := bundle_recipient.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Create(paramsBundleRecipientCreate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = lib.JsonMarshal(result, fieldsCreate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	cmdCreate.Flags().Int64VarP(&paramsBundleRecipientCreate.UserId, "user-id", "u", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().Int64VarP(&paramsBundleRecipientCreate.BundleId, "bundle-id", "b", 0, "Bundle to share.")
	cmdCreate.Flags().StringVarP(&paramsBundleRecipientCreate.Recipient, "recipient", "r", "", "Email addresses to share this bundle with.")
	cmdCreate.Flags().StringVarP(&paramsBundleRecipientCreate.Name, "name", "n", "", "Name of recipient.")
	cmdCreate.Flags().StringVarP(&paramsBundleRecipientCreate.Company, "company", "c", "", "Company of recipient.")
	cmdCreate.Flags().StringVarP(&paramsBundleRecipientCreate.Note, "note", "o", "", "Note to include in email.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	BundleRecipients.AddCommand(cmdCreate)
}
