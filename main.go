package main

import (
	_ "embed"

	"github.com/Files-com/files-cli/cmd"
	files "github.com/Files-com/files-sdk-go/v2"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.Init(version, commit, date, &files.GlobalConfig)
}
