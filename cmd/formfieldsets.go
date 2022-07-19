package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	flib "github.com/Files-com/files-sdk-go/v2/lib"

	"fmt"

	form_field_set "github.com/Files-com/files-sdk-go/v2/formfieldset"
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
	var formatList string
	usePagerList := true
	paramsFormFieldSetList := files_sdk.FormFieldSetListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List Form Field Sets",
		Long:  `List Form Field Sets`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsFormFieldSetList
			params.MaxPages = MaxPagesList

			client := form_field_set.Client{Config: *config}
			it, err := client.List(ctx, params)
			it.OnPageError = func(err error) (*[]interface{}, error) {
				overriddenValues, newErr := lib.ErrorWithOriginalResponse(err, formatList, config.Logger())
				values, ok := overriddenValues.([]interface{})
				if ok {
					return &values, newErr
				} else {
					return &[]interface{}{}, newErr
				}
			}
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(ctx, it, formatList, fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
		},
	}

	cmdList.Flags().Int64Var(&paramsFormFieldSetList.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVar(&paramsFormFieldSetList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsFormFieldSetList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	FormFieldSets.AddCommand(cmdList)
	var fieldsFind string
	var formatFind string
	usePagerFind := true
	paramsFormFieldSetFind := files_sdk.FormFieldSetFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Form Field Set`,
		Long:  `Show Form Field Set`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := form_field_set.Client{Config: *config}

			var formFieldSet interface{}
			var err error
			formFieldSet, err = client.Find(ctx, paramsFormFieldSetFind)
			lib.HandleResponse(ctx, formFieldSet, err, formatFind, fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdFind.Flags().Int64Var(&paramsFormFieldSetFind.Id, "id", 0, "Form Field Set ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	FormFieldSets.AddCommand(cmdFind)
	var fieldsCreate string
	var formatCreate string
	usePagerCreate := true
	createSkipEmail := false
	createSkipName := false
	createSkipCompany := false
	paramsFormFieldSetCreate := files_sdk.FormFieldSetCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Form Field Set`,
		Long:  `Create Form Field Set`,
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

			var formFieldSet interface{}
			var err error
			formFieldSet, err = client.Create(ctx, paramsFormFieldSetCreate)
			lib.HandleResponse(ctx, formFieldSet, err, formatCreate, fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdCreate.Flags().Int64Var(&paramsFormFieldSetCreate.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().StringVar(&paramsFormFieldSetCreate.Title, "title", "", "Title to be displayed")
	cmdCreate.Flags().BoolVar(&createSkipEmail, "skip-email", createSkipEmail, "Skip validating form email")
	cmdCreate.Flags().BoolVar(&createSkipName, "skip-name", createSkipName, "Skip validating form name")
	cmdCreate.Flags().BoolVar(&createSkipCompany, "skip-company", createSkipCompany, "Skip validating company")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	FormFieldSets.AddCommand(cmdCreate)
	var fieldsUpdate string
	var formatUpdate string
	usePagerUpdate := true
	updateSkipEmail := false
	updateSkipName := false
	updateSkipCompany := false
	paramsFormFieldSetUpdate := files_sdk.FormFieldSetUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Form Field Set`,
		Long:  `Update Form Field Set`,
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

			var formFieldSet interface{}
			var err error
			formFieldSet, err = client.Update(ctx, paramsFormFieldSetUpdate)
			lib.HandleResponse(ctx, formFieldSet, err, formatUpdate, fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsFormFieldSetUpdate.Id, "id", 0, "Form Field Set ID.")
	cmdUpdate.Flags().StringVar(&paramsFormFieldSetUpdate.Title, "title", "", "Title to be displayed")
	cmdUpdate.Flags().BoolVar(&updateSkipEmail, "skip-email", updateSkipEmail, "Skip validating form email")
	cmdUpdate.Flags().BoolVar(&updateSkipName, "skip-name", updateSkipName, "Skip validating form name")
	cmdUpdate.Flags().BoolVar(&updateSkipCompany, "skip-company", updateSkipCompany, "Skip validating company")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	cmdUpdate.Flags().StringVarP(&formatUpdate, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	FormFieldSets.AddCommand(cmdUpdate)
	var fieldsDelete string
	var formatDelete string
	usePagerDelete := true
	paramsFormFieldSetDelete := files_sdk.FormFieldSetDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Form Field Set`,
		Long:  `Delete Form Field Set`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := form_field_set.Client{Config: *config}

			var err error
			err = client.Delete(ctx, paramsFormFieldSetDelete)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
		},
	}
	cmdDelete.Flags().Int64Var(&paramsFormFieldSetDelete.Id, "id", 0, "Form Field Set ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	FormFieldSets.AddCommand(cmdDelete)
}
