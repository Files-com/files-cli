package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	as2_station "github.com/Files-com/files-sdk-go/v2/as2station"
)

var (
	As2Stations = &cobra.Command{}
)

func As2StationsInit() {
	As2Stations = &cobra.Command{
		Use:  "as2-stations [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command as2-stations\n\t%v", args[0])
		},
	}
	var fieldsList string
	var formatList string
	paramsAs2StationList := files_sdk.As2StationListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsAs2StationList
			params.MaxPages = MaxPagesList

			client := as2_station.Client{Config: *config}
			it, err := client.List(ctx, params)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(it, formatList, fieldsList, listFilter, cmd.OutOrStdout())
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
		},
	}

	cmdList.Flags().StringVar(&paramsAs2StationList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsAs2StationList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright")
	As2Stations.AddCommand(cmdList)
	var fieldsFind string
	var formatFind string
	paramsAs2StationFind := files_sdk.As2StationFindParams{}

	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := as2_station.Client{Config: *config}

			var as2Station interface{}
			var err error
			as2Station, err = client.Find(ctx, paramsAs2StationFind)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(as2Station, formatFind, fieldsFind, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}
	cmdFind.Flags().Int64Var(&paramsAs2StationFind.Id, "id", 0, "As2 Station ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-bright")
	As2Stations.AddCommand(cmdFind)
	var fieldsCreate string
	var formatCreate string
	paramsAs2StationCreate := files_sdk.As2StationCreateParams{}

	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := as2_station.Client{Config: *config}

			var as2Station interface{}
			var err error
			as2Station, err = client.Create(ctx, paramsAs2StationCreate)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(as2Station, formatCreate, fieldsCreate, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}
	cmdCreate.Flags().StringVar(&paramsAs2StationCreate.Name, "name", "", "AS2 Name")
	cmdCreate.Flags().StringVar(&paramsAs2StationCreate.PublicCertificate, "public-certificate", "", "")
	cmdCreate.Flags().StringVar(&paramsAs2StationCreate.PrivateKey, "private-key", "", "")
	cmdCreate.Flags().StringVar(&paramsAs2StationCreate.PrivateKeyPassword, "private-key-password", "", "")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	As2Stations.AddCommand(cmdCreate)
	var fieldsUpdate string
	var formatUpdate string
	paramsAs2StationUpdate := files_sdk.As2StationUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := as2_station.Client{Config: *config}

			var as2Station interface{}
			var err error
			as2Station, err = client.Update(ctx, paramsAs2StationUpdate)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(as2Station, formatUpdate, fieldsUpdate, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsAs2StationUpdate.Id, "id", 0, "As2 Station ID.")
	cmdUpdate.Flags().StringVar(&paramsAs2StationUpdate.Name, "name", "", "AS2 Name")
	cmdUpdate.Flags().StringVar(&paramsAs2StationUpdate.PublicCertificate, "public-certificate", "", "")
	cmdUpdate.Flags().StringVar(&paramsAs2StationUpdate.PrivateKey, "private-key", "", "")
	cmdUpdate.Flags().StringVar(&paramsAs2StationUpdate.PrivateKeyPassword, "private-key-password", "", "")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	cmdUpdate.Flags().StringVarP(&formatUpdate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	As2Stations.AddCommand(cmdUpdate)
	var fieldsDelete string
	var formatDelete string
	paramsAs2StationDelete := files_sdk.As2StationDeleteParams{}

	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := as2_station.Client{Config: *config}

			var err error
			err = client.Delete(ctx, paramsAs2StationDelete)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
		},
	}
	cmdDelete.Flags().Int64Var(&paramsAs2StationDelete.Id, "id", 0, "As2 Station ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-bright")
	As2Stations.AddCommand(cmdDelete)
}
