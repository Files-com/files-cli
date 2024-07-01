package cmd

import (
	"fmt"

	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	form_field_set "github.com/Files-com/files-sdk-go/v3/formfieldset"
	flib "github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(FormFieldSets())
}

func FormFieldSets() *cobra.Command {
	FormFieldSets := &cobra.Command{
		Use:  "form-field-sets [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command form-field-sets\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsFormFieldSetList := files_sdk.FormFieldSetListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Form Field Sets",
		Long:    `List Form Field Sets`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsFormFieldSetList
			params.MaxPages = MaxPagesList

			client := form_field_set.Client{Config: config}
			it, err := client.List(params, files_sdk.WithContext(ctx))
			it.OnPageError = func(err error) (*[]interface{}, error) {
				overriddenValues, newErr := lib.ErrorWithOriginalResponse(err, config.Logger)
				values, ok := overriddenValues.([]interface{})
				if ok {
					return &values, newErr
				} else {
					return &[]interface{}{}, newErr
				}
			}
			if err != nil {
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			if len(filterbyList) > 0 {
				listFilter = func(i interface{}) (interface{}, bool, error) {
					matchOk, err := lib.MatchFilter(filterbyList, i)
					return i, matchOk, err
				}
			}
			err = lib.FormatIter(ctx, it, Profile(cmd).Current().SetResourceFormat(cmd, formatList), fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdList.Flags().StringToStringVar(&filterbyList, "filter-by", filterbyList, `Client side filtering: field-name=*.jpg,field-name=?ello`)

	cmdList.Flags().Int64Var(&paramsFormFieldSetList.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVar(&paramsFormFieldSetList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsFormFieldSetList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVar(&paramsFormFieldSetList.Action, "action", "", "")
	cmdList.Flags().Int64Var(&paramsFormFieldSetList.Page, "page", 0, "")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	FormFieldSets.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsFormFieldSetFind := files_sdk.FormFieldSetFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Form Field Set`,
		Long:  `Show Form Field Set`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := form_field_set.Client{Config: config}

			var formFieldSet interface{}
			var err error
			formFieldSet, err = client.Find(paramsFormFieldSetFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), formFieldSet, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsFormFieldSetFind.Id, "id", 0, "Form Field Set ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	FormFieldSets.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	createSkipEmail := true
	createSkipName := true
	createSkipCompany := true
	paramsFormFieldSetCreate := files_sdk.FormFieldSetCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Form Field Set`,
		Long:  `Create Form Field Set`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := form_field_set.Client{Config: config}

			if cmd.Flags().Changed("skip-email") {
				paramsFormFieldSetCreate.SkipEmail = flib.Bool(createSkipEmail)
			}
			if cmd.Flags().Changed("skip-name") {
				paramsFormFieldSetCreate.SkipName = flib.Bool(createSkipName)
			}
			if cmd.Flags().Changed("skip-company") {
				paramsFormFieldSetCreate.SkipCompany = flib.Bool(createSkipCompany)
			}

			var formFieldSet interface{}
			var err error
			formFieldSet, err = client.Create(paramsFormFieldSetCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), formFieldSet, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().Int64Var(&paramsFormFieldSetCreate.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().StringVar(&paramsFormFieldSetCreate.Title, "title", "", "Title to be displayed")
	cmdCreate.Flags().BoolVar(&createSkipEmail, "skip-email", createSkipEmail, "Skip validating form email")
	cmdCreate.Flags().BoolVar(&createSkipName, "skip-name", createSkipName, "Skip validating form name")
	cmdCreate.Flags().BoolVar(&createSkipCompany, "skip-company", createSkipCompany, "Skip validating company")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	FormFieldSets.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	updateSkipEmail := true
	updateSkipName := true
	updateSkipCompany := true
	paramsFormFieldSetUpdate := files_sdk.FormFieldSetUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Form Field Set`,
		Long:  `Update Form Field Set`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := form_field_set.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.FormFieldSetUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsFormFieldSetUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("title") {
				lib.FlagUpdate(cmd, "title", paramsFormFieldSetUpdate.Title, mapParams)
			}
			if cmd.Flags().Changed("skip-email") {
				mapParams["skip_email"] = updateSkipEmail
			}
			if cmd.Flags().Changed("skip-name") {
				mapParams["skip_name"] = updateSkipName
			}
			if cmd.Flags().Changed("skip-company") {
				mapParams["skip_company"] = updateSkipCompany
			}
			if cmd.Flags().Changed("form-fields") {
				lib.FlagUpdateLen(cmd, "form_fields", paramsFormFieldSetUpdate.FormFields, mapParams)
			}

			var formFieldSet interface{}
			var err error
			formFieldSet, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), formFieldSet, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsFormFieldSetUpdate.Id, "id", 0, "Form Field Set ID.")
	cmdUpdate.Flags().StringVar(&paramsFormFieldSetUpdate.Title, "title", "", "Title to be displayed")
	cmdUpdate.Flags().BoolVar(&updateSkipEmail, "skip-email", updateSkipEmail, "Skip validating form email")
	cmdUpdate.Flags().BoolVar(&updateSkipName, "skip-name", updateSkipName, "Skip validating form name")
	cmdUpdate.Flags().BoolVar(&updateSkipCompany, "skip-company", updateSkipCompany, "Skip validating company")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	FormFieldSets.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsFormFieldSetDelete := files_sdk.FormFieldSetDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Form Field Set`,
		Long:  `Delete Form Field Set`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := form_field_set.Client{Config: config}

			var err error
			err = client.Delete(paramsFormFieldSetDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsFormFieldSetDelete.Id, "id", 0, "Form Field Set ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	FormFieldSets.AddCommand(cmdDelete)
	return FormFieldSets
}
