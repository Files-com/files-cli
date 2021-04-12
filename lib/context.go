package lib

import (
	"time"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/lib"
)

type Context struct {
	Values map[string]*files_sdk.Config
	Config *files_sdk.Config
	Test   bool
}

func (c Context) Done() <-chan struct{} {
	return make(chan struct{})
}

func (c Context) Err() error {
	return nil
}

func (c Context) Deadline() (deadline time.Time, ok bool) {
	return time.Now(), false
}

func (c Context) Value(interface{}) interface{} {
	return lib.Interface()
}

func (c Context) GetConfig() *files_sdk.Config {
	return c.Config
}

func (c Context) Testing() bool {
	return c.Test
}
