package cmd

import (
	"bytes"
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Files-com/files-sdk-go/v2/file"
	"github.com/Files-com/files-sdk-go/v2/lib"
)

func TestDownload(t *testing.T) {
	mutex := &sync.Mutex{}
	t.Run("files-cli download", func(t *testing.T) {
		sourceFs := &file.FS{Context: context.Background()}
		destinationFs := lib.ReadWriteFs(lib.LocalFileSystem{})
		for _, tt := range lib.PathSpec(sourceFs.PathSeparator(), destinationFs.PathSeparator()) {
			t.Run(tt.Name, func(t *testing.T) {
				r, config, err := CreateConfig(t.Name())
				if err != nil {
					t.Fatal(err)
				}
				sourceFs := (&file.FS{Context: context.Background()}).Init(*config, false)
				lib.BuildPathSpecTest(t, mutex, tt, sourceFs, destinationFs, func(source, destination string) lib.Cmd {
					return Cmd(config, Download(), []string{source, destination}, []string{"--format", "text"})
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
			require.NoError(t, downloadCmd.Run())
			assert.Equal(t, stdout.String(), ". errored open : Not Found - `Not Found`\n")
			assert.Equal(t, stderr.String(), "")
			r.Stop()
		})
	})
}
