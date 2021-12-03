package cmd

import (
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/stretchr/testify/assert"
)

func TestExternalEventsCreateSuccess(t *testing.T) {
	assert := assert.New(t)
	r, config, err := CreateConfig("TestExternalEventsCreate")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	ExternalEventsInit()
	out, stdErr := callCmd(ExternalEvents, config, []string{"create", "--status", "success", "--body", "this is a success test", "--format", "json"})
	assert.Equal("", string(stdErr))
	event := files_sdk.ExternalEvent{}
	event.UnmarshalJSON(out)
	assert.Equal("this is a success test", event.Body)
	assert.Equal("client_log", event.EventType)
	assert.Equal("success", event.Status)
}

func TestExternalEventsCreateError(t *testing.T) {
	assert := assert.New(t)
	r, config, err := CreateConfig("TestExternalEventsCreateError")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	ExternalEventsInit()
	out, stderr := callCmd(ExternalEvents, config, []string{"create", "--status", "failure", "--body", "this is a error test", "--format", "json"})
	assert.Equal("", string(stderr))
	event := files_sdk.ExternalEvent{}
	err = event.UnmarshalJSON(out)
	assert.NoError(err)
	assert.Equal("this is a error test", event.Body)
	assert.Equal("client_log", event.EventType)
	assert.Equal("failure", event.Status)
}

func TestExternalEventsCreateNoStatus(t *testing.T) {
	assert := assert.New(t)
	r, config, err := CreateConfig("TestExternalEventsCreateNoStatus")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	ExternalEventsInit()
	stdOut, errOut := callCmd(ExternalEvents, config, []string{"create", "--status", "taco", "--body", "this is a error test", "--format", "json"})
	event := files_sdk.ExternalEvent{}
	err = event.UnmarshalJSON(stdOut)
	assert.Error(err)
	assert.Equal("missing required field: ExternalEventCreateParams{}.Status\n", string(errOut))
	assert.Equal("", string(stdOut))
}
