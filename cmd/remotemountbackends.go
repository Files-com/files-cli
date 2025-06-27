package cmd

import (
	"fmt"
	"reflect"

	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	flib "github.com/Files-com/files-sdk-go/v3/lib"
	remote_mount_backend "github.com/Files-com/files-sdk-go/v3/remotemountbackend"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(RemoteMountBackends())
}

func RemoteMountBackends() *cobra.Command {
	RemoteMountBackends := &cobra.Command{
		Use:  "remote-mount-backends [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command remote-mount-backends\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsRemoteMountBackendList := files_sdk.RemoteMountBackendListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Remote Mount Backends",
		Long:    `List Remote Mount Backends`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsRemoteMountBackendList
			params.MaxPages = MaxPagesList

			client := remote_mount_backend.Client{Config: config}
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

	cmdList.Flags().StringToStringVar(&filterbyList, "filter-by", filterbyList, `Client side filtering: field-name=*.jpg,field-name=?ello`)

	cmdList.Flags().StringVar(&paramsRemoteMountBackendList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsRemoteMountBackendList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	RemoteMountBackends.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsRemoteMountBackendFind := files_sdk.RemoteMountBackendFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Remote Mount Backend`,
		Long:  `Show Remote Mount Backend`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := remote_mount_backend.Client{Config: config}

			var remoteMountBackend interface{}
			var err error
			remoteMountBackend, err = client.Find(paramsRemoteMountBackendFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), remoteMountBackend, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsRemoteMountBackendFind.Id, "id", 0, "Remote Mount Backend ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	RemoteMountBackends.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	createEnabled := true
	createHealthCheckEnabled := true
	paramsRemoteMountBackendCreate := files_sdk.RemoteMountBackendCreateParams{}
	RemoteMountBackendCreateHealthCheckType := ""

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Remote Mount Backend`,
		Long:  `Create Remote Mount Backend`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := remote_mount_backend.Client{Config: config}

			var RemoteMountBackendCreateHealthCheckTypeErr error
			paramsRemoteMountBackendCreate.HealthCheckType, RemoteMountBackendCreateHealthCheckTypeErr = lib.FetchKey("health-check-type", paramsRemoteMountBackendCreate.HealthCheckType.Enum(), RemoteMountBackendCreateHealthCheckType)
			if RemoteMountBackendCreateHealthCheckType != "" && RemoteMountBackendCreateHealthCheckTypeErr != nil {
				return RemoteMountBackendCreateHealthCheckTypeErr
			}

			if cmd.Flags().Changed("enabled") {
				paramsRemoteMountBackendCreate.Enabled = flib.Bool(createEnabled)
			}
			if cmd.Flags().Changed("health-check-enabled") {
				paramsRemoteMountBackendCreate.HealthCheckEnabled = flib.Bool(createHealthCheckEnabled)
			}

			var remoteMountBackend interface{}
			var err error
			remoteMountBackend, err = client.Create(paramsRemoteMountBackendCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), remoteMountBackend, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().BoolVar(&createEnabled, "enabled", createEnabled, "True if this backend is enabled.")
	cmdCreate.Flags().Int64Var(&paramsRemoteMountBackendCreate.Fall, "fall", 0, "Number of consecutive failures before considering the backend unhealthy.")
	cmdCreate.Flags().BoolVar(&createHealthCheckEnabled, "health-check-enabled", createHealthCheckEnabled, "True if health checks are enabled for this backend.")
	cmdCreate.Flags().StringVar(&RemoteMountBackendCreateHealthCheckType, "health-check-type", "", fmt.Sprintf("Type of health check to perform. %v", reflect.ValueOf(paramsRemoteMountBackendCreate.HealthCheckType.Enum()).MapKeys()))
	cmdCreate.Flags().Int64Var(&paramsRemoteMountBackendCreate.Interval, "interval", 0, "Interval in seconds between health checks.")
	cmdCreate.Flags().StringVar(&paramsRemoteMountBackendCreate.MinFreeCpu, "min-free-cpu", "", "Minimum free CPU percentage required for this backend to be considered healthy.")
	cmdCreate.Flags().StringVar(&paramsRemoteMountBackendCreate.MinFreeMem, "min-free-mem", "", "Minimum free memory percentage required for this backend to be considered healthy.")
	cmdCreate.Flags().Int64Var(&paramsRemoteMountBackendCreate.Priority, "priority", 0, "Priority of this backend.")
	cmdCreate.Flags().StringVar(&paramsRemoteMountBackendCreate.RemotePath, "remote-path", "", "Path on the remote server to treat as the root of this mount.")
	cmdCreate.Flags().Int64Var(&paramsRemoteMountBackendCreate.Rise, "rise", 0, "Number of consecutive successes before considering the backend healthy.")
	cmdCreate.Flags().StringVar(&paramsRemoteMountBackendCreate.CanaryFilePath, "canary-file-path", "", "Path to the canary file used for health checks.")
	cmdCreate.Flags().Int64Var(&paramsRemoteMountBackendCreate.RemoteServerMountId, "remote-server-mount-id", 0, "The mount ID of the Remote Server Mount that this backend is associated with.")
	cmdCreate.Flags().Int64Var(&paramsRemoteMountBackendCreate.RemoteServerId, "remote-server-id", 0, "The remote server that this backend is associated with.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	RemoteMountBackends.AddCommand(cmdCreate)
	var fieldsResetStatus []string
	var formatResetStatus []string
	usePagerResetStatus := true
	paramsRemoteMountBackendResetStatus := files_sdk.RemoteMountBackendResetStatusParams{}

	cmdResetStatus := &cobra.Command{
		Use:   "reset-status",
		Short: `Reset backend status to healthy`,
		Long:  `Reset backend status to healthy`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := remote_mount_backend.Client{Config: config}

			var err error
			err = client.ResetStatus(paramsRemoteMountBackendResetStatus, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdResetStatus.Flags().Int64Var(&paramsRemoteMountBackendResetStatus.Id, "id", 0, "Remote Mount Backend ID.")

	cmdResetStatus.Flags().StringSliceVar(&fieldsResetStatus, "fields", []string{}, "comma separated list of field names")
	cmdResetStatus.Flags().StringSliceVar(&formatResetStatus, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdResetStatus.Flags().BoolVar(&usePagerResetStatus, "use-pager", usePagerResetStatus, "Use $PAGER (.ie less, more, etc)")

	RemoteMountBackends.AddCommand(cmdResetStatus)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	updateEnabled := true
	updateHealthCheckEnabled := true
	paramsRemoteMountBackendUpdate := files_sdk.RemoteMountBackendUpdateParams{}
	RemoteMountBackendUpdateHealthCheckType := ""

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Remote Mount Backend`,
		Long:  `Update Remote Mount Backend`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := remote_mount_backend.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.RemoteMountBackendUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			var RemoteMountBackendUpdateHealthCheckTypeErr error
			paramsRemoteMountBackendUpdate.HealthCheckType, RemoteMountBackendUpdateHealthCheckTypeErr = lib.FetchKey("health-check-type", paramsRemoteMountBackendUpdate.HealthCheckType.Enum(), RemoteMountBackendUpdateHealthCheckType)
			if RemoteMountBackendUpdateHealthCheckType != "" && RemoteMountBackendUpdateHealthCheckTypeErr != nil {
				return RemoteMountBackendUpdateHealthCheckTypeErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsRemoteMountBackendUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("enabled") {
				mapParams["enabled"] = updateEnabled
			}
			if cmd.Flags().Changed("fall") {
				lib.FlagUpdate(cmd, "fall", paramsRemoteMountBackendUpdate.Fall, mapParams)
			}
			if cmd.Flags().Changed("health-check-enabled") {
				mapParams["health_check_enabled"] = updateHealthCheckEnabled
			}
			if cmd.Flags().Changed("health-check-type") {
				lib.FlagUpdate(cmd, "health_check_type", paramsRemoteMountBackendUpdate.HealthCheckType, mapParams)
			}
			if cmd.Flags().Changed("interval") {
				lib.FlagUpdate(cmd, "interval", paramsRemoteMountBackendUpdate.Interval, mapParams)
			}
			if cmd.Flags().Changed("min-free-cpu") {
				lib.FlagUpdate(cmd, "min_free_cpu", paramsRemoteMountBackendUpdate.MinFreeCpu, mapParams)
			}
			if cmd.Flags().Changed("min-free-mem") {
				lib.FlagUpdate(cmd, "min_free_mem", paramsRemoteMountBackendUpdate.MinFreeMem, mapParams)
			}
			if cmd.Flags().Changed("priority") {
				lib.FlagUpdate(cmd, "priority", paramsRemoteMountBackendUpdate.Priority, mapParams)
			}
			if cmd.Flags().Changed("remote-path") {
				lib.FlagUpdate(cmd, "remote_path", paramsRemoteMountBackendUpdate.RemotePath, mapParams)
			}
			if cmd.Flags().Changed("rise") {
				lib.FlagUpdate(cmd, "rise", paramsRemoteMountBackendUpdate.Rise, mapParams)
			}
			if cmd.Flags().Changed("canary-file-path") {
				lib.FlagUpdate(cmd, "canary_file_path", paramsRemoteMountBackendUpdate.CanaryFilePath, mapParams)
			}
			if cmd.Flags().Changed("remote-server-id") {
				lib.FlagUpdate(cmd, "remote_server_id", paramsRemoteMountBackendUpdate.RemoteServerId, mapParams)
			}

			var remoteMountBackend interface{}
			var err error
			remoteMountBackend, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), remoteMountBackend, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsRemoteMountBackendUpdate.Id, "id", 0, "Remote Mount Backend ID.")
	cmdUpdate.Flags().BoolVar(&updateEnabled, "enabled", updateEnabled, "True if this backend is enabled.")
	cmdUpdate.Flags().Int64Var(&paramsRemoteMountBackendUpdate.Fall, "fall", 0, "Number of consecutive failures before considering the backend unhealthy.")
	cmdUpdate.Flags().BoolVar(&updateHealthCheckEnabled, "health-check-enabled", updateHealthCheckEnabled, "True if health checks are enabled for this backend.")
	cmdUpdate.Flags().StringVar(&RemoteMountBackendUpdateHealthCheckType, "health-check-type", "", fmt.Sprintf("Type of health check to perform. %v", reflect.ValueOf(paramsRemoteMountBackendUpdate.HealthCheckType.Enum()).MapKeys()))
	cmdUpdate.Flags().Int64Var(&paramsRemoteMountBackendUpdate.Interval, "interval", 0, "Interval in seconds between health checks.")
	cmdUpdate.Flags().StringVar(&paramsRemoteMountBackendUpdate.MinFreeCpu, "min-free-cpu", "", "Minimum free CPU percentage required for this backend to be considered healthy.")
	cmdUpdate.Flags().StringVar(&paramsRemoteMountBackendUpdate.MinFreeMem, "min-free-mem", "", "Minimum free memory percentage required for this backend to be considered healthy.")
	cmdUpdate.Flags().Int64Var(&paramsRemoteMountBackendUpdate.Priority, "priority", 0, "Priority of this backend.")
	cmdUpdate.Flags().StringVar(&paramsRemoteMountBackendUpdate.RemotePath, "remote-path", "", "Path on the remote server to treat as the root of this mount.")
	cmdUpdate.Flags().Int64Var(&paramsRemoteMountBackendUpdate.Rise, "rise", 0, "Number of consecutive successes before considering the backend healthy.")
	cmdUpdate.Flags().StringVar(&paramsRemoteMountBackendUpdate.CanaryFilePath, "canary-file-path", "", "Path to the canary file used for health checks.")
	cmdUpdate.Flags().Int64Var(&paramsRemoteMountBackendUpdate.RemoteServerId, "remote-server-id", 0, "The remote server that this backend is associated with.")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	RemoteMountBackends.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsRemoteMountBackendDelete := files_sdk.RemoteMountBackendDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Remote Mount Backend`,
		Long:  `Delete Remote Mount Backend`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := remote_mount_backend.Client{Config: config}

			var err error
			err = client.Delete(paramsRemoteMountBackendDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsRemoteMountBackendDelete.Id, "id", 0, "Remote Mount Backend ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	RemoteMountBackends.AddCommand(cmdDelete)
	return RemoteMountBackends
}
