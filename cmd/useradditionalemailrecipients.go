package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	user_additional_email_recipient "github.com/Files-com/files-sdk-go/v3/useradditionalemailrecipient"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(UserAdditionalEmailRecipients())
}

func UserAdditionalEmailRecipients() *cobra.Command {
	UserAdditionalEmailRecipients := &cobra.Command{
		Use:  "user-additional-email-recipients [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command user-additional-email-recipients\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsUserAdditionalEmailRecipientList := files_sdk.UserAdditionalEmailRecipientListParams{}
	var MaxPagesList int64
	var listSortByArgs string
	var listFilterArgs []string
	var listFilterPrefixArgs []string

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List User Additional Email Recipients",
		Long:    `List User Additional Email Recipients`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsUserAdditionalEmailRecipientList
			params.MaxPages = MaxPagesList

			parsedListSortBy, parseListSortByErr := lib.ParseAPIListSortFlag("sort-by", listSortByArgs)
			if parseListSortByErr != nil {
				return parseListSortByErr
			}
			if parsedListSortBy != nil {
				params.SortBy = parsedListSortBy
			}
			parsedListFilter, parseListFilterErr := lib.ParseAPIListQueryFlag("filter", listFilterArgs)
			if parseListFilterErr != nil {
				return parseListFilterErr
			}
			if parsedListFilter != nil {
				params.Filter = parsedListFilter
			}
			parsedListFilterPrefix, parseListFilterPrefixErr := lib.ParseAPIListQueryFlag("filter-prefix", listFilterPrefixArgs)
			if parseListFilterPrefixErr != nil {
				return parseListFilterPrefixErr
			}
			if parsedListFilterPrefix != nil {
				params.FilterPrefix = parsedListFilterPrefix
			}

			client := user_additional_email_recipient.Client{Config: config}
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
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			if len(filterbyList) > 0 {
				listFilter = func(i interface{}) (interface{}, bool, error) {
					matchOk, err := lib.MatchFilter(filterbyList, i)
					return i, matchOk, err
				}
			}
			err = lib.FormatIter(ctx, it, Profile(cmd).Current().SetResourceFormat(cmd, formatList), fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdList.Flags().StringToStringVar(&filterbyList, "filter-by", filterbyList, "Client-side wildcard filtering, for example field-name=*.jpg or field-name=?ello")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter-by", "field=pattern")
	cmdList.Flags().StringVar(&listSortByArgs, "sort-by", "", "Sort user additional email recipients by field in ascending or descending order.")
	lib.SetFlagDisplayType(cmdList.Flags(), "sort-by", "field=asc|desc")
	cmdList.Flags().StringArrayVar(&listFilterArgs, "filter", []string{}, "Find user additional email recipients where field exactly matches value.")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter", "field=value")
	cmdList.Flags().StringArrayVar(&listFilterPrefixArgs, "filter-prefix", []string{}, "Find user additional email recipients where field starts with value.")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter-prefix", "field=value")

	cmdList.Flags().Int64Var(&paramsUserAdditionalEmailRecipientList.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVar(&paramsUserAdditionalEmailRecipientList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsUserAdditionalEmailRecipientList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	UserAdditionalEmailRecipients.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsUserAdditionalEmailRecipientFind := files_sdk.UserAdditionalEmailRecipientFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show User Additional Email Recipient`,
		Long:  `Show User Additional Email Recipient`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := user_additional_email_recipient.Client{Config: config}

			var userAdditionalEmailRecipient interface{}
			var err error
			userAdditionalEmailRecipient, err = client.Find(paramsUserAdditionalEmailRecipientFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), userAdditionalEmailRecipient, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsUserAdditionalEmailRecipientFind.Id, "id", 0, "User Additional Email Recipient ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	UserAdditionalEmailRecipients.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	paramsUserAdditionalEmailRecipientCreate := files_sdk.UserAdditionalEmailRecipientCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create User Additional Email Recipient`,
		Long:  `Create User Additional Email Recipient`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := user_additional_email_recipient.Client{Config: config}

			var userAdditionalEmailRecipient interface{}
			var err error
			userAdditionalEmailRecipient, err = client.Create(paramsUserAdditionalEmailRecipientCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), userAdditionalEmailRecipient, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().Int64Var(&paramsUserAdditionalEmailRecipientCreate.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().StringVar(&paramsUserAdditionalEmailRecipientCreate.Email, "email", "", "Additional email recipient address")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	UserAdditionalEmailRecipients.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	paramsUserAdditionalEmailRecipientUpdate := files_sdk.UserAdditionalEmailRecipientUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update User Additional Email Recipient`,
		Long:  `Update User Additional Email Recipient`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := user_additional_email_recipient.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.UserAdditionalEmailRecipientUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsUserAdditionalEmailRecipientUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("email") {
				lib.FlagUpdate(cmd, "email", paramsUserAdditionalEmailRecipientUpdate.Email, mapParams)
			}

			var userAdditionalEmailRecipient interface{}
			var err error
			userAdditionalEmailRecipient, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), userAdditionalEmailRecipient, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsUserAdditionalEmailRecipientUpdate.Id, "id", 0, "User Additional Email Recipient ID.")
	cmdUpdate.Flags().StringVar(&paramsUserAdditionalEmailRecipientUpdate.Email, "email", "", "Additional email recipient address")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	UserAdditionalEmailRecipients.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsUserAdditionalEmailRecipientDelete := files_sdk.UserAdditionalEmailRecipientDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete User Additional Email Recipient`,
		Long:  `Delete User Additional Email Recipient`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := user_additional_email_recipient.Client{Config: config}

			var err error
			err = client.Delete(paramsUserAdditionalEmailRecipientDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsUserAdditionalEmailRecipientDelete.Id, "id", 0, "User Additional Email Recipient ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	UserAdditionalEmailRecipients.AddCommand(cmdDelete)
	return UserAdditionalEmailRecipients
}
