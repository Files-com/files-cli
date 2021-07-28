package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	flib "github.com/Files-com/files-sdk-go/lib"

	"fmt"

	form_field_set "github.com/Files-com/files-sdk-go/formfieldset"
)

var (
	FormFieldSets = &cobra.Command{}
)

func FormFieldSetsInit() {
	FormFieldSets = &cobra.Command{
		Use:  "form-field-sets [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command form-field-sets\n\t%v", args[0])
		},
	}
	var fieldsList string
	paramsFormFieldSetList := files_sdk.FormFieldSetListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsFormFieldSetList
			params.MaxPages = MaxPagesList

			client := form_field_set.Client{Config: *config}
			it, err := client.List(ctx, params)
			if err != nil {
				lib.ClientError(ctx, err)
			}
			err = lib.JsonMarshalIter(it, fieldsList)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdList.Flags().Int64VarP(&paramsFormFieldSetList.UserId, "user-id", "u", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVarP(&paramsFormFieldSetList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().Int64VarP(&paramsFormFieldSetList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	FormFieldSets.AddCommand(cmdList)
	var fieldsFind string
	paramsFormFieldSetFind := files_sdk.FormFieldSetFindParams{}

	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := form_field_set.Client{Config: *config}

			result, err := client.Find(ctx, paramsFormFieldSetFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.JsonMarshal(result, fieldsFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdFind.Flags().Int64VarP(&paramsFormFieldSetFind.Id, "id", "i", 0, "Form Field Set ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	FormFieldSets.AddCommand(cmdFind)
	var fieldsCreate string
	createSkipEmail := false
	createSkipName := false
	createSkipCompany := false
	paramsFormFieldSetCreate := files_sdk.FormFieldSetCreateParams{}

	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := form_field_set.Client{Config: *config}

			if createSkipEmail {
				paramsFormFieldSetCreate.SkipEmail = flib.Bool(true)
			}
			if createSkipName {
				paramsFormFieldSetCreate.SkipName = flib.Bool(true)
			}
			if createSkipCompany {
				paramsFormFieldSetCreate.SkipCompany = flib.Bool(true)
			}

			result, err := client.Create(ctx, paramsFormFieldSetCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.JsonMarshal(result, fieldsCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdCreate.Flags().Int64VarP(&paramsFormFieldSetCreate.UserId, "user-id", "u", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().StringVarP(&paramsFormFieldSetCreate.Title, "title", "t", "", "Title to be displayed")
	cmdCreate.Flags().BoolVarP(&createSkipEmail, "skip-email", "e", createSkipEmail, "Skip validating form email")
	cmdCreate.Flags().BoolVarP(&createSkipName, "skip-name", "n", createSkipName, "Skip validating form name")
	cmdCreate.Flags().BoolVarP(&createSkipCompany, "skip-company", "c", createSkipCompany, "Skip validating company")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	FormFieldSets.AddCommand(cmdCreate)
	var fieldsUpdate string
	updateSkipEmail := false
	updateSkipName := false
	updateSkipCompany := false
	paramsFormFieldSetUpdate := files_sdk.FormFieldSetUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := form_field_set.Client{Config: *config}

			if updateSkipEmail {
				paramsFormFieldSetUpdate.SkipEmail = flib.Bool(true)
			}
			if updateSkipName {
				paramsFormFieldSetUpdate.SkipName = flib.Bool(true)
			}
			if updateSkipCompany {
				paramsFormFieldSetUpdate.SkipCompany = flib.Bool(true)
			}

			result, err := client.Update(ctx, paramsFormFieldSetUpdate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.JsonMarshal(result, fieldsUpdate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdUpdate.Flags().Int64VarP(&paramsFormFieldSetUpdate.Id, "id", "i", 0, "Form Field Set ID.")
	cmdUpdate.Flags().StringVarP(&paramsFormFieldSetUpdate.Title, "title", "t", "", "Title to be displayed")
	cmdUpdate.Flags().BoolVarP(&updateSkipEmail, "skip-email", "e", updateSkipEmail, "Skip validating form email")
	cmdUpdate.Flags().BoolVarP(&updateSkipName, "skip-name", "n", updateSkipName, "Skip validating form name")
	cmdUpdate.Flags().BoolVarP(&updateSkipCompany, "skip-company", "c", updateSkipCompany, "Skip validating company")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	FormFieldSets.AddCommand(cmdUpdate)
	var fieldsDelete string
	paramsFormFieldSetDelete := files_sdk.FormFieldSetDeleteParams{}

	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := form_field_set.Client{Config: *config}

			result, err := client.Delete(ctx, paramsFormFieldSetDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.JsonMarshal(result, fieldsDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsFormFieldSetDelete.Id, "id", "i", 0, "Form Field Set ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	FormFieldSets.AddCommand(cmdDelete)
}
