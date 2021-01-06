package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	"fmt"
	"os"

	form_field_set "github.com/Files-com/files-sdk-go/formfieldset"
)

var (
	FormFieldSets = &cobra.Command{
		Use:  "form-field-sets [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func FormFieldSetsInit() {
	var fieldsList string
	paramsFormFieldSetList := files_sdk.FormFieldSetListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsFormFieldSetList
			params.MaxPages = MaxPagesList
			it, err := form_field_set.List(params)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().Int64VarP(&paramsFormFieldSetList.UserId, "user-id", "u", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVarP(&paramsFormFieldSetList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().IntVarP(&paramsFormFieldSetList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	FormFieldSets.AddCommand(cmdList)
	var fieldsFind string
	paramsFormFieldSetFind := files_sdk.FormFieldSetFindParams{}
	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := form_field_set.Find(paramsFormFieldSetFind)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsFind)
		},
	}
	cmdFind.Flags().Int64VarP(&paramsFormFieldSetFind.Id, "id", "i", 0, "Form Field Set ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	FormFieldSets.AddCommand(cmdFind)
	var fieldsCreate string
	paramsFormFieldSetCreate := files_sdk.FormFieldSetCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := form_field_set.Create(paramsFormFieldSetCreate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsCreate)
		},
	}
	cmdCreate.Flags().Int64VarP(&paramsFormFieldSetCreate.UserId, "user-id", "u", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().StringVarP(&paramsFormFieldSetCreate.Title, "title", "t", "", "Title to be displayed")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	FormFieldSets.AddCommand(cmdCreate)
	var fieldsUpdate string
	paramsFormFieldSetUpdate := files_sdk.FormFieldSetUpdateParams{}
	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := form_field_set.Update(paramsFormFieldSetUpdate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsUpdate)
		},
	}
	cmdUpdate.Flags().Int64VarP(&paramsFormFieldSetUpdate.Id, "id", "i", 0, "Form Field Set ID.")
	cmdUpdate.Flags().StringVarP(&paramsFormFieldSetUpdate.Title, "title", "t", "", "Title to be displayed")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	FormFieldSets.AddCommand(cmdUpdate)
	var fieldsDelete string
	paramsFormFieldSetDelete := files_sdk.FormFieldSetDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := form_field_set.Delete(paramsFormFieldSetDelete)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsDelete)
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsFormFieldSetDelete.Id, "id", "i", 0, "Form Field Set ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	FormFieldSets.AddCommand(cmdDelete)
}
