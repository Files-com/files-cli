package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	metadata_category "github.com/Files-com/files-sdk-go/v3/metadatacategory"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(MetadataCategories())
}

func MetadataCategories() *cobra.Command {
	MetadataCategories := &cobra.Command{
		Use:  "metadata-categories [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command metadata-categories\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsMetadataCategoryList := files_sdk.MetadataCategoryListParams{}
	var MaxPagesList int64
	var listSortByArgs string

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Metadata Categories",
		Long:    `List Metadata Categories`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsMetadataCategoryList
			params.MaxPages = MaxPagesList

			parsedListSortBy, parseListSortByErr := lib.ParseAPIListSortFlag("sort-by", listSortByArgs)
			if parseListSortByErr != nil {
				return parseListSortByErr
			}
			if parsedListSortBy != nil {
				params.SortBy = parsedListSortBy
			}

			client := metadata_category.Client{Config: config}
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
	cmdList.Flags().StringVar(&listSortByArgs, "sort-by", "", "Sort metadata categories by field in ascending or descending order.")
	lib.SetFlagDisplayType(cmdList.Flags(), "sort-by", "field=asc|desc")

	cmdList.Flags().StringVar(&paramsMetadataCategoryList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsMetadataCategoryList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	MetadataCategories.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsMetadataCategoryFind := files_sdk.MetadataCategoryFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Metadata Category`,
		Long:  `Show Metadata Category`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := metadata_category.Client{Config: config}

			var metadataCategory interface{}
			var err error
			metadataCategory, err = client.Find(paramsMetadataCategoryFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), metadataCategory, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsMetadataCategoryFind.Id, "id", 0, "Metadata Category ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	MetadataCategories.AddCommand(cmdFind)
	var fieldsListFor []string
	var formatListFor []string
	usePagerListFor := true
	filterbyListFor := make(map[string]string)
	paramsMetadataCategoryListFor := files_sdk.MetadataCategoryListForParams{}
	var MaxPagesListFor int64

	cmdListFor := &cobra.Command{
		Use:     "list-for [path]",
		Short:   "List Metadata Categories by Path",
		Long:    `List Metadata Categories by Path`,
		Args:    cobra.RangeArgs(0, 1),
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsMetadataCategoryListFor
			params.MaxPages = MaxPagesListFor
			if len(args) > 0 && args[0] != "" {
				params.Path = args[0]
			}

			client := metadata_category.Client{Config: config}
			it, err := client.ListFor(params, files_sdk.WithContext(ctx))
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
			if len(filterbyListFor) > 0 {
				listFilter = func(i interface{}) (interface{}, bool, error) {
					matchOk, err := lib.MatchFilter(filterbyListFor, i)
					return i, matchOk, err
				}
			}
			err = lib.FormatIter(ctx, it, Profile(cmd).Current().SetResourceFormat(cmd, formatListFor), fieldsListFor, usePagerListFor, listFilter, cmd.OutOrStdout())
			return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdListFor.Flags().StringToStringVar(&filterbyListFor, "filter-by", filterbyListFor, "Client-side wildcard filtering, for example field-name=*.jpg or field-name=?ello")
	lib.SetFlagDisplayType(cmdListFor.Flags(), "filter-by", "field=pattern")

	cmdListFor.Flags().StringVar(&paramsMetadataCategoryListFor.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdListFor.Flags().Int64Var(&paramsMetadataCategoryListFor.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdListFor.Flags().StringVar(&paramsMetadataCategoryListFor.Path, "path", "", "Path to operate on.")

	cmdListFor.Flags().Int64VarP(&MaxPagesListFor, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdListFor.Flags().StringSliceVar(&fieldsListFor, "fields", []string{}, "comma separated list of field names to include in response")
	cmdListFor.Flags().StringSliceVar(&formatListFor, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdListFor.Flags().BoolVar(&usePagerListFor, "use-pager", usePagerListFor, "Use $PAGER (.ie less, more, etc)")
	MetadataCategories.AddCommand(cmdListFor)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	paramsMetadataCategoryCreate := files_sdk.MetadataCategoryCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Metadata Category`,
		Long:  `Create Metadata Category`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := metadata_category.Client{Config: config}

			var metadataCategory interface{}
			var err error
			metadataCategory, err = client.Create(paramsMetadataCategoryCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), metadataCategory, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().StringVar(&paramsMetadataCategoryCreate.Name, "name", "", "Name of the metadata category.")
	cmdCreate.Flags().StringSliceVar(&paramsMetadataCategoryCreate.DefaultColumns, "default-columns", []string{}, "Metadata keys that should appear as columns in the UI by default.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	MetadataCategories.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	paramsMetadataCategoryUpdate := files_sdk.MetadataCategoryUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Metadata Category`,
		Long:  `Update Metadata Category`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := metadata_category.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.MetadataCategoryUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsMetadataCategoryUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("name") {
				lib.FlagUpdate(cmd, "name", paramsMetadataCategoryUpdate.Name, mapParams)
			}
			if cmd.Flags().Changed("default-columns") {
				lib.FlagUpdateLen(cmd, "default_columns", paramsMetadataCategoryUpdate.DefaultColumns, mapParams)
			}

			var metadataCategory interface{}
			var err error
			metadataCategory, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), metadataCategory, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsMetadataCategoryUpdate.Id, "id", 0, "Metadata Category ID.")
	cmdUpdate.Flags().StringVar(&paramsMetadataCategoryUpdate.Name, "name", "", "Name of the metadata category.")
	cmdUpdate.Flags().StringSliceVar(&paramsMetadataCategoryUpdate.DefaultColumns, "default-columns", []string{}, "Metadata keys that should appear as columns in the UI by default.")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	MetadataCategories.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsMetadataCategoryDelete := files_sdk.MetadataCategoryDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Metadata Category`,
		Long:  `Delete Metadata Category`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := metadata_category.Client{Config: config}

			var err error
			err = client.Delete(paramsMetadataCategoryDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsMetadataCategoryDelete.Id, "id", 0, "Metadata Category ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	MetadataCategories.AddCommand(cmdDelete)
	return MetadataCategories
}
