package cmd

import (
	"bytes"
	"context"
	"sync"
	"testing"

	"github.com/Files-com/files-sdk-go/v3/file"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/stretchr/testify/assert"
)

func TestDownload(t *testing.T) {
	mutex := &sync.Mutex{}
	t.Run("files-cli download", func(t *testing.T) {
		sourceFs := &file.FS{Context: context.Background()}
		destinationFs := lib.ReadWriteFs(lib.LocalFileSystem{})
		for _, tt := range lib.PathSpec(t, sourceFs.PathSeparator(), destinationFs.PathSeparator()) {
			t.Run(tt.Name, func(t *testing.T) {
				r, config, err := CreateConfig(t.Name())
				if err != nil {
					t.Fatal(err)
				}
				sourceFs := (&file.FS{Context: context.Background()}).Init(config, false)
				lib.BuildPathSpecTest(t, mutex, tt, sourceFs, destinationFs, func(args lib.PathSpecArgs) lib.Cmd {
					if args.PreserveTimes {
						return Cmd(config, Download(), []string{args.Src, args.Dest}, []string{"--format", "text", "--times"})

					}
					return Cmd(config, Download(), []string{args.Src, args.Dest}, []string{"--format", "text"})
				})
				r.Stop()
			})
		}

		t.Run("not found", func(t *testing.T) {
			r, config, err := CreateConfig(t.Name())
			if err != nil {
				t.Fatal(err)
			}
			downloadCmd := Cmd(config, Download(), []string{"you will never find me", t.TempDir()}, []string{"--format", "text"})
			stdout := bytes.NewBufferString("")
			stderr := bytes.NewBufferString("")
			downloadCmd.SetErr(stderr)
			downloadCmd.SetOut(stdout)

			err = downloadCmd.Run()

			assert.Error(t, err)
			assert.Equal(t, err.Error(), "open you will never find me: Not Found - `Not Found.  This may be related to your permissions.` - status (7)")

			assert.Equal(t, stdout.String(), "")
			assert.Equal(t, stderr.String(), "Error: open you will never find me: Not Found - `Not Found.  This may be related to your permissions.` - status (7)\n")
			r.Stop()
		})
	})
}
