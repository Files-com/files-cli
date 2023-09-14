package cmd

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Files-com/files-cli/lib/changelog"
	"github.com/spf13/cobra"
)

func init() {
	changeLog := Changelog()
	RootCmd.AddCommand(changeLog)
	IgnoreCredentialsCheck = append(IgnoreCredentialsCheck, changeLog.Use)
}

func Changelog() *cobra.Command {
	return &cobra.Command{
		Use:   "changelog",
		Short: "Returns the changelog for the current version",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(cmd.Context(), 30*time.Second)
			defer cancel()

			tag := ""
			if len(args) == 1 {
				tag = args[0]
			}

			tagRange := strings.Split(tag, "..")
			if len(tagRange) == 2 {
				versions, err := changelog.GetAllTags(ctx, *Profile(cmd).Config, changelog.ParseTag(tagRange[0]), changelog.ParseTag(tagRange[1]))
				if err != nil {
					return err
				}

				for i, tag := range versions {
					var err error
					if i == 0 {
						err = changelog.GetLog(ctx, cmd, *Profile(cmd).Config, tag, "# ChangeLog\n\n")
					} else {
						err = changelog.GetLog(ctx, cmd, *Profile(cmd).Config, tag, "")
					}

					if err != nil {
						return err
					}
				}
			} else {
				re := regexp.MustCompile("[0-9/.]+")
				result := re.FindAllString(tag, -1)
				if len(result) > 0 {
					tag = fmt.Sprintf("v%v", result[0])
				}
				return changelog.GetLog(ctx, cmd, *Profile(cmd).Config, tag, "# ChangeLog\n\n")
			}

			return nil
		},
	}
}
