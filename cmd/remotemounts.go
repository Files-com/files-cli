package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/behavior"
	"github.com/Files-com/files-sdk-go/v2/folder"
	remote_server "github.com/Files-com/files-sdk-go/v2/remoteserver"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(RemoteMounts())
}

type MountValue struct {
	RemotePath     string `json:"remote_path"`
	RemoteServerId int64  `json:"remote_server_id"`
}

func RemoteMounts() *cobra.Command {
	mounts := &cobra.Command{
		Use:   "remote-mounts",
		Short: "Create a Remote Mount",
		Long:  "Create a remote mount to a remote server",
		Args:  cobra.ExactArgs(1),
	}
	var mountValue MountValue
	createParams := files_sdk.BehaviorCreateParams{Behavior: "remote_server_mount"}
	var createFolder bool
	var remoteServerName string
	create := &cobra.Command{
		Use:  "create",
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if remoteServerName != "" {
				var err error
				mountValue.RemoteServerId, err = findRemoteServer(cmd, remoteServerName)
				if err != nil {
					return err
				}
			}

			valueBytes, err := json.Marshal(&mountValue)
			if err != nil {
				return err
			}
			createParams.Value = string(valueBytes)

			if createFolder {
				folderClient := folder.Client{Config: *Profile(cmd).Config}
				_, err := folderClient.Create(cmd.Context(), files_sdk.FolderCreateParams{Path: createParams.Path})
				if err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "%v\n", err)
				}
			}

			client := behavior.Client{Config: *Profile(cmd).Config}
			resource, err := client.Create(cmd.Context(), createParams)
			if err != nil {
				return err
			}
			expandedResource, _ := expandBehavior()(resource)
			lib.HandleResponse(
				cmd.Context(),
				Profile(cmd),
				expandedResource,
				err,
				[]string{},
				[]string{},
				false,
				cmd.OutOrStdout(),
				cmd.ErrOrStderr(),
				Profile(cmd).Logger(),
			)

			return err
		},
	}

	create.Flags().BoolVar(&createFolder, "create-path", false, "create folder for mount")
	create.Flags().StringVar(&createParams.Name, "name", createParams.Name, "")
	create.MarkFlagRequired("name")
	create.Flags().StringVar(&createParams.Path, "path", createParams.Path, "")
	create.MarkFlagRequired("path")
	create.Flags().StringVar(&remoteServerName, "remote-server-name", remoteServerName, "")
	create.Flags().StringVar(&mountValue.RemotePath, "remote-path", mountValue.RemotePath, "")
	create.Flags().Int64Var(&mountValue.RemoteServerId, "remote-server-id", mountValue.RemoteServerId, "")
	mounts.AddCommand(create)

	var fieldsList []string
	var formatList []string
	usePagerList := true
	paramsBehaviorList := files_sdk.BehaviorListParams{Behavior: "remote_server_mount"}
	var MaxPagesList int64
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List Remote Mounts",
		Long:  `List Remote Mounts`,
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsBehaviorList
			params.MaxPages = MaxPagesList

			client := behavior.Client{Config: *config}
			it, err := client.List(ctx, params)
			it.OnPageError = func(err error) (*[]interface{}, error) {
				overriddenValues, newErr := lib.ErrorWithOriginalResponse(err, config.Logger())
				values, ok := overriddenValues.([]interface{})
				if ok {
					return &values, newErr
				} else {
					return &[]interface{}{}, newErr
				}
			}
			if err != nil {
				lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}

			err = lib.FormatIter(ctx, it, formatList, fieldsList, usePagerList, expandBehavior(), cmd.OutOrStdout())
			if err != nil {
				lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}

	cmdList.Flags().StringVar(&paramsBehaviorList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsBehaviorList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().Int64Var(&paramsBehaviorList.Page, "page", 0, "Current page number.")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{"id", "path", "value"}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
        table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
        json-styles: {raw, pretty}
        `)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	mounts.AddCommand(cmdList)

	var fieldsFind []string
	var formatFind []string
	paramsBehaviorFind := files_sdk.BehaviorFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Remote Mount`,
		Long:  `Show Remote Mount`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := behavior.Client{Config: *config}

			var resource interface{}
			var err error
			resource, err = client.Find(ctx, paramsBehaviorFind)
			if err != nil {
				return err
			}
			expandedResource, _ := expandBehavior()(resource)
			lib.HandleResponse(ctx, Profile(cmd), expandedResource, err, formatFind, fieldsFind, false, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdFind.Flags().Int64Var(&paramsBehaviorFind.Id, "id", 0, "Behavior ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)

	mounts.AddCommand(cmdFind)

	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsBehaviorDelete := files_sdk.BehaviorDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Remote Mount`,
		Long:  `Delete Remote Mount`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := behavior.Client{Config: *config}

			var err error
			err = client.Delete(ctx, paramsBehaviorDelete)
			if err != nil {
				lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsBehaviorDelete.Id, "id", 0, "Behavior ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	mounts.AddCommand(cmdDelete)

	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	updateParams := files_sdk.BehaviorUpdateParams{Behavior: "remote_server_mount"}

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Behavior`,
		Long:  `Update Behavior`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			client := behavior.Client{Config: *Profile(cmd).Config}
			valueBytes, err := json.Marshal(&mountValue)

			if err != nil {
				return err
			}
			updateParams.Value = string(valueBytes)

			resource, err := client.Update(ctx, updateParams)
			if err != nil {
				return nil
			}
			expandedResource, _ := expandBehavior()(resource)
			lib.HandleResponse(ctx, Profile(cmd), expandedResource, err, formatUpdate, fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), Profile(cmd).Logger())
			return nil
		},
	}
	cmdUpdate.Flags().Int64Var(&updateParams.Id, "id", 0, "Behavior ID.")
	cmdUpdate.Flags().StringVar(&updateParams.Path, "path", updateParams.Path, "")
	cmdUpdate.Flags().StringVar(&updateParams.Name, "name", updateParams.Name, "")
	cmdUpdate.Flags().StringVar(&mountValue.RemotePath, "remote-path", mountValue.RemotePath, "")
	cmdUpdate.Flags().Int64Var(&mountValue.RemoteServerId, "remote-server-id", mountValue.RemoteServerId, "")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	mounts.AddCommand(cmdUpdate)

	return mounts
}

func findRemoteServer(cmd *cobra.Command, remoteServerName string) (int64, error) {
	remoteServerClient := remote_server.Client{Config: *Profile(cmd).Config}
	it, err := remoteServerClient.List(cmd.Context(), files_sdk.RemoteServerListParams{})
	if err != nil {
		return 0, err
	}

	for it.Next() {
		if it.Err() != nil {
			return 0, err
		}

		if it.RemoteServer().Name == remoteServerName {
			return it.RemoteServer().Id, nil
		}
	}
	return 0, fmt.Errorf("no remote server found '%v", remoteServerName)
}

func expandBehavior() func(i interface{}) (interface{}, bool) {
	return func(i interface{}) (interface{}, bool) {
		resource := make(map[string]interface{})
		j, err := json.Marshal(&i)
		if err != nil {
			return i, true
		}

		err = json.Unmarshal(j, &resource)
		if err != nil {
			return i, true
		}

		value := resource["value"].(map[string]interface{})
		delete(resource, "value")
		delete(resource, "behavior")
		for k, v := range value {
			resource[k] = v
		}

		return resource, true
	}
}
