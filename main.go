package main

import (
	"embed"

	"github.com/Files-com/files-cli/cmd"
	files "github.com/Files-com/files-sdk-go/v3"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// guides ships the agent guides with the binary so `files-cli workflows`
// serves the version that matches this build, offline. The embed lives in
// this file because releases build ./main.go on its own.
//
//go:embed CONTEXT.md skills/recipes
var guides embed.FS

func main() {
	cmd.Guides = guides
	cmd.Init(version, commit, date, files.GlobalConfig)
}
