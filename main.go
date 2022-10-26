package main

import (
	_ "embed"

	"github.com/Files-com/files-cli/cmd"
	files "github.com/Files-com/files-sdk-go/v2"
)

//go:embed _VERSION
var VERSION string

func main() {
	cmd.Init(VERSION, &files.GlobalConfig)
}
