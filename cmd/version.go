package cmd

import (
	"fmt"
	"strings"

	runtimeDebug "runtime/debug"

	"github.com/spf13/cobra"
)

var (
	VersionCmd *cobra.Command
	sdkVersion bool
)

func init() {
	VersionCmd = &cobra.Command{
		Use:     "version",
		Aliases: []string{"version"},
		Args:    cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("files-cli version %v\n", RootCmd.Version)
			if !sdkVersion {
				return
			}
			bi, ok := runtimeDebug.ReadBuildInfo()
			if ok {
				for _, dep := range bi.Deps {
					if strings.Contains(dep.Path, "github.com/Files-com/files-sdk-go") {
						fmt.Printf("SDK       version %v\n", dep.Version)
						return
					}
				}
			}
		},
	}
	VersionCmd.Flags().BoolVar(&sdkVersion, "verbose", sdkVersion, "list other dependencies")
	RootCmd.AddCommand(VersionCmd)
}
