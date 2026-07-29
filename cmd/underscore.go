package cmd

import (
	"fmt"

	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/transfers"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file"
	flib "github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/spf13/cobra"
)

type uploadToDestination func(*file.Client, int64, string, ...file.UploadOption) error
type copyToDestination func(*file.Client, files_sdk.FileCopyParams, int64, string, ...files_sdk.RequestResponseOption) (files_sdk.FileAction, error)
type moveToDestination func(*file.Client, files_sdk.FileMoveParams, int64, string, ...files_sdk.RequestResponseOption) (files_sdk.FileAction, error)

func init() {
	RootCmd.AddCommand(underscoreUploadCommands()...)
}

func underscoreUploadCommands() []*cobra.Command {
	return []*cobra.Command{
		uploadToDestinationCommand("remote-server", "Remote Server", "remote-server-id", func(client *file.Client, id int64, destinationPath string, opts ...file.UploadOption) error {
			return client.UploadToRemoteServer(id, destinationPath, opts...)
		}),
		uploadToDestinationCommand("snapshot", "Snapshot", "snapshot-id", func(client *file.Client, id int64, destinationPath string, opts ...file.UploadOption) error {
			return client.UploadToSnapshot(id, destinationPath, opts...)
		}),
		uploadToDestinationCommand("child-site", "Child Site", "site-id", func(client *file.Client, id int64, destinationPath string, opts ...file.UploadOption) error {
			return client.UploadToChildSite(id, destinationPath, opts...)
		}),
	}
}

func addUnderscoreFileCommands(parent *cobra.Command) {
	parent.AddCommand(
		copyToDestinationCommand("remote-server", "Remote Server", "remote-server-id", func(client *file.Client, params files_sdk.FileCopyParams, id int64, destinationPath string, opts ...files_sdk.RequestResponseOption) (files_sdk.FileAction, error) {
			return client.CopyToRemoteServer(params, id, destinationPath, opts...)
		}),
		moveToDestinationCommand("remote-server", "Remote Server", "remote-server-id", func(client *file.Client, params files_sdk.FileMoveParams, id int64, destinationPath string, opts ...files_sdk.RequestResponseOption) (files_sdk.FileAction, error) {
			return client.MoveToRemoteServer(params, id, destinationPath, opts...)
		}),
	)
	parent.AddCommand(
		copyToDestinationCommand("snapshot", "Snapshot", "snapshot-id", func(client *file.Client, params files_sdk.FileCopyParams, id int64, destinationPath string, opts ...files_sdk.RequestResponseOption) (files_sdk.FileAction, error) {
			return client.CopyToSnapshot(params, id, destinationPath, opts...)
		}),
		moveToDestinationCommand("snapshot", "Snapshot", "snapshot-id", func(client *file.Client, params files_sdk.FileMoveParams, id int64, destinationPath string, opts ...files_sdk.RequestResponseOption) (files_sdk.FileAction, error) {
			return client.MoveToSnapshot(params, id, destinationPath, opts...)
		}),
	)
	parent.AddCommand(
		copyToDestinationCommand("child-site", "Child Site", "site-id", func(client *file.Client, params files_sdk.FileCopyParams, id int64, destinationPath string, opts ...files_sdk.RequestResponseOption) (files_sdk.FileAction, error) {
			return client.CopyToChildSite(params, id, destinationPath, opts...)
		}),
		moveToDestinationCommand("child-site", "Child Site", "site-id", func(client *file.Client, params files_sdk.FileMoveParams, id int64, destinationPath string, opts ...files_sdk.RequestResponseOption) (files_sdk.FileAction, error) {
			return client.MoveToChildSite(params, id, destinationPath, opts...)
		}),
	)
}

func uploadToDestinationCommand(name string, label string, idFlag string, upload uploadToDestination) *cobra.Command {
	var destinationID int64
	command := &cobra.Command{
		Use:   fmt.Sprintf("upload-to-%s <source-path> <destination-path>", name),
		Short: fmt.Sprintf("Upload a file to a %s.", label),
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := file.Client{Config: config}
			err := upload(
				&client,
				destinationID,
				args[1],
				file.UploadWithContext(ctx),
				file.UploadWithFile(args[0]),
			)
			return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}
	command.Flags().Int64Var(&destinationID, idFlag, 0, fmt.Sprintf("%s ID. Required.", label))
	command.MarkFlagRequired(idFlag)
	return command
}

func copyToDestinationCommand(name string, label string, idFlag string, copyFile copyToDestination) *cobra.Command {
	var destinationID int64
	var destinationPath string
	var fields []string
	var format []string
	usePager := true
	var block bool
	var noProgress bool
	var eventLog bool
	copyBehaviors := true
	structure := true
	overwrite := true
	params := files_sdk.FileCopyParams{}

	command := &cobra.Command{
		Use:   fmt.Sprintf("copy-to-%s [path]", name),
		Short: fmt.Sprintf("Copy a file or folder to a %s.", label),
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := file.Client{Config: config}

			if cmd.Flags().Changed("copy-behaviors") {
				params.CopyBehaviors = flib.Bool(copyBehaviors)
			}
			if cmd.Flags().Changed("structure") {
				params.Structure = flib.Bool(structure)
			}
			if cmd.Flags().Changed("overwrite") {
				params.Overwrite = flib.Bool(overwrite)
			}
			if len(args) > 0 {
				params.Path = args[0]
			}

			result, err := copyFile(&client, params, destinationID, destinationPath, files_sdk.WithContext(ctx))
			return handleUnderscoreFileAction(cmd, config, result, err, block, noProgress, eventLog, format, fields, usePager)
		},
	}
	command.Flags().StringVar(&params.Path, "path", "", "Path to operate on.")
	command.Flags().StringVar(&destinationPath, "destination", "", "Destination path relative to the external destination. Required.")
	command.Flags().Int64Var(&destinationID, idFlag, 0, fmt.Sprintf("%s ID. Required.", label))
	command.Flags().BoolVar(&copyBehaviors, "copy-behaviors", copyBehaviors, "If copying a folder, also copy supported behaviors to the destination folder tree.")
	command.Flags().BoolVar(&structure, "structure", structure, "Copy structure only.")
	command.Flags().BoolVar(&overwrite, "overwrite", overwrite, "Overwrite existing files in the destination.")
	addUnderscoreFileActionFlags(command, "copy", &fields, &format, &usePager, &block, &noProgress, &eventLog)
	command.MarkFlagRequired("destination")
	command.MarkFlagRequired(idFlag)
	return command
}

func moveToDestinationCommand(name string, label string, idFlag string, moveFile moveToDestination) *cobra.Command {
	var destinationID int64
	var destinationPath string
	var fields []string
	var format []string
	usePager := true
	var block bool
	var noProgress bool
	var eventLog bool
	overwrite := true
	params := files_sdk.FileMoveParams{}

	command := &cobra.Command{
		Use:   fmt.Sprintf("move-to-%s [path]", name),
		Short: fmt.Sprintf("Move a file or folder to a %s.", label),
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := file.Client{Config: config}

			if cmd.Flags().Changed("overwrite") {
				params.Overwrite = flib.Bool(overwrite)
			}
			if len(args) > 0 {
				params.Path = args[0]
			}

			result, err := moveFile(&client, params, destinationID, destinationPath, files_sdk.WithContext(ctx))
			return handleUnderscoreFileAction(cmd, config, result, err, block, noProgress, eventLog, format, fields, usePager)
		},
	}
	command.Flags().StringVar(&params.Path, "path", "", "Path to operate on.")
	command.Flags().StringVar(&destinationPath, "destination", "", "Destination path relative to the external destination. Required.")
	command.Flags().Int64Var(&destinationID, idFlag, 0, fmt.Sprintf("%s ID. Required.", label))
	command.Flags().BoolVar(&overwrite, "overwrite", overwrite, "Overwrite existing files in the destination.")
	addUnderscoreFileActionFlags(command, "move", &fields, &format, &usePager, &block, &noProgress, &eventLog)
	command.MarkFlagRequired("destination")
	command.MarkFlagRequired(idFlag)
	return command
}

func addUnderscoreFileActionFlags(command *cobra.Command, operation string, fields *[]string, format *[]string, usePager *bool, block *bool, noProgress *bool, eventLog *bool) {
	command.Flags().StringSliceVar(fields, "fields", []string{}, "Comma-separated list of field names.")
	command.Flags().StringSliceVar(format, "format", lib.FormatDefaults, lib.FormatHelpText)
	command.Flags().BoolVar(usePager, "use-pager", true, "Use $PAGER (.ie less, more, etc)")
	command.Flags().BoolVar(block, "block", false, fmt.Sprintf("Wait for the asynchronous %s to finish.", operation))
	command.Flags().BoolVar(noProgress, "no-progress", false, "Do not display progress while waiting.")
	command.Flags().BoolVar(eventLog, "event-log", false, fmt.Sprintf("Output the full event log for the %s when waiting.", operation))
}

func handleUnderscoreFileAction(cmd *cobra.Command, config files_sdk.Config, result files_sdk.FileAction, err error, block bool, noProgress bool, eventLog bool, format []string, fields []string, usePager bool) error {
	if err != nil {
		return err
	}
	ctx := cmd.Context()
	value, err := transfers.WaitFileMigration(ctx, config, result, block, noProgress, eventLog, Profile(cmd).Current().SetResourceFormat(cmd, format), cmd.OutOrStdout())
	if err != nil {
		return err
	}
	return lib.HandleResponse(ctx, Profile(cmd), value, nil, Profile(cmd).Current().SetResourceFormat(cmd, format), fields, usePager, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
}
