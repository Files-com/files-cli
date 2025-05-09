package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/behavior"
	"github.com/Files-com/files-sdk-go/v3/folder"
	remote_server "github.com/Files-com/files-sdk-go/v3/remoteserver"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(RemoteMounts())
}

type MountValue struct {
	RemotePath     string `json:"remote_path,omitempty"`
	RemoteServerId int64  `json:"remote_server_id,omitempty"`
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
		Args: cobra.NoArgs,
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
				_, err := folderClient.Create(files_sdk.FolderCreateParams{Path: createParams.Path}, files_sdk.WithContext(cmd.Context()))
				if err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "%v\n", err)
				}
			}

			client := behavior.Client{Config: *Profile(cmd).Config}
			resource, err := client.Create(createParams, files_sdk.WithContext(cmd.Context()))
			if err != nil {
				return err
			}
			expandedResource, _, _ := expandBehavior()(resource)
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
				Profile(cmd).Logger,
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
	paramsBehaviorList := files_sdk.BehaviorListParams{Filter: files_sdk.Behavior{Behavior: "remote_server_mount"}}
	var MaxPagesList int64
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List Remote Mounts",
		Long:  `List Remote Mounts`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsBehaviorList
			params.MaxPages = MaxPagesList

			client := behavior.Client{Config: config}
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

			err = lib.FormatIter(ctx, it, Profile(cmd).Current().SetResourceFormat(cmd, formatList), fieldsList, usePagerList, expandBehavior(), cmd.OutOrStdout())
			return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdList.Flags().StringVar(&paramsBehaviorList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsBehaviorList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
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
			config := ctx.Value("config").(files_sdk.Config)
			client := behavior.Client{Config: config}

			var resource interface{}
			var err error
			resource, err = client.Find(paramsBehaviorFind, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			var expandedResource interface{}
			expandedResource, _, err = expandBehavior()(resource)
			return lib.HandleResponse(ctx, Profile(cmd), expandedResource, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, false, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsBehaviorFind.Id, "id", 0, "Behavior ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)

	mounts.AddCommand(cmdFind)

	var fieldsDelete []string
	usePagerDelete := true
	paramsBehaviorDelete := files_sdk.BehaviorDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Remote Mount`,
		Long:  `Delete Remote Mount`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := behavior.Client{Config: config}

			var err error
			err = client.Delete(paramsBehaviorDelete, files_sdk.WithContext(ctx))
			return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}
	cmdDelete.Flags().Int64Var(&paramsBehaviorDelete.Id, "id", 0, "Behavior ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	mounts.AddCommand(cmdDelete)

	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	updateParams := files_sdk.BehaviorUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Behavior`,
		Long:  `Update Behavior`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			client := behavior.Client{Config: *Profile(cmd).Config}
			valueBytes, err := json.Marshal(&mountValue)

			if err != nil {
				return err
			}
			if len(valueBytes) > 2 { // empty json object `{}`
				updateParams.Value = string(valueBytes)
			}

			resource, err := client.Update(updateParams, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			var expandedResource interface{}
			expandedResource, _, err = expandBehavior()(resource)
			return lib.HandleResponse(ctx, Profile(cmd), expandedResource, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), Profile(cmd).Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&updateParams.Id, "id", 0, "Behavior ID.")
	cmdUpdate.Flags().StringVar(&updateParams.Name, "name", updateParams.Name, "")
	cmdUpdate.Flags().StringVar(&mountValue.RemotePath, "remote-path", mountValue.RemotePath, "")
	cmdUpdate.Flags().Int64Var(&mountValue.RemoteServerId, "remote-server-id", mountValue.RemoteServerId, "")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	mounts.AddCommand(cmdUpdate)

	return mounts
}

func findRemoteServer(cmd *cobra.Command, remoteServerName string) (int64, error) {
	remoteServerClient := remote_server.Client{Config: *Profile(cmd).Config}
	it, err := remoteServerClient.List(files_sdk.RemoteServerListParams{}, files_sdk.WithContext(cmd.Context()))
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
	return 0, clierr.Errorf(clierr.ErrorCodeFatal, "no remote server found '%v", remoteServerName)
}

type RemoteMount struct {
	Id             int64  `json:"id"`
	Path           string `json:"path"`
	RemoteServerId int64  `json:"remote_server_id"`
	RemotePath     string `json:"remote_path"`
	Description    string `json:"description"`
}

func expandBehavior() func(i interface{}) (interface{}, bool, error) {
	return func(i interface{}) (interface{}, bool, error) {
		resource := make(map[string]interface{})
		j, err := json.Marshal(&i)
		if err != nil {
			return i, true, nil
		}

		err = json.Unmarshal(j, &resource)
		if err != nil {
			return i, true, nil
		}

		value := resource["value"].(map[string]interface{})
		delete(resource, "value")
		delete(resource, "behavior")
		for k, v := range value {
			resource[k] = v
		}

		j, err = json.Marshal(resource)
		if err != nil {
			return i, true, nil
		}
		var remoteMount RemoteMount
		err = json.Unmarshal(j, &remoteMount)
		if err != nil {
			return i, true, nil
		}

		return remoteMount, true, nil
	}
}
