package lib

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Files-com/files-sdk-go/v2/lib"
)

func TestHandleResponse(t *testing.T) {
	type args struct {
		ctx    context.Context
		i      interface{}
		errIn  func() error
		format []string
		fields []string
	}
	tests := []struct {
		name       string
		args       args
		wantStdout string
		wantStderr string
		wantLogger string
	}{
		{
			name: "Can't format as json for unexpected server response",
			args: args{
				context.Background(),
				nil,
				func() error {
					var structure map[string]interface{}
					return lib.ErrorWithOriginalResponse{}.ProcessError([]byte(`{}`), &json.UnmarshalTypeError{Value: "number", Type: reflect.TypeOf("")}, structure)
				},
				[]string{"table"},
				[]string{},
			},
			wantStdout: "",
			wantLogger: "Recovering from original error: `json: cannot unmarshal number into Go value of type string`\n",
		},
		{
			name: "unexpected server response prints json",
			args: args{
				context.Background(),
				nil,

				func() error {
					var structure map[string]interface{}
					return lib.ErrorWithOriginalResponse{}.ProcessError([]byte(`{"key":"value"}`), &json.UnmarshalTypeError{Value: "number", Type: reflect.TypeOf("")}, structure)
				},
				[]string{"json"},
				[]string{},
			},
			wantStdout: "{\n    \"key\": \"value\"\n}\n",
			wantLogger: "Recovering from original error: `json: cannot unmarshal number into Go value of type string`\n",
		},
		{
			name: "unexpected server response prints json for collection",
			args: args{
				context.Background(),
				nil,

				func() error {
					var structure []map[string]interface{}
					return lib.ErrorWithOriginalResponse{}.ProcessError([]byte(`[{"key":"value"}]`), &json.UnmarshalTypeError{Value: "number", Type: reflect.TypeOf("")}, structure)
				},
				[]string{"json"},
				[]string{},
			},
			wantStdout: "[{\n    \"key\": \"value\"\n}]\n",
			wantLogger: "Recovering from original error: `json: cannot unmarshal number into Go value of type string`\n",
		},
		{
			name: "other error type",
			args: args{
				context.Background(),
				nil,
				func() error { return fmt.Errorf("some other error") },
				[]string{"json"},
				[]string{},
			},
			wantStderr: "some other error\n",
		},
		{
			name: "prints interface",
			args: args{
				context.Background(),
				map[string]interface{}{"key": "value"},
				func() error { return nil },
				[]string{"json"},
				[]string{},
			},
			wantStdout: "{\n    \"key\": \"value\"\n}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}
			logger := &bytes.Buffer{}

			HandleResponse(tt.args.ctx, &Profiles{}, tt.args.i, tt.args.errIn(), tt.args.format, tt.args.fields, false, stdout, stderr, log.New(logger, "", 0))
			assert.Equalf(t, tt.wantStdout, stdout.String(), "stdout HandleResponse(%v, %v, %v, %v, %v, %v, %v, %v)", tt.args.ctx, tt.args.i, tt.args.errIn(), tt.args.format, tt.args.fields, stdout, stderr, logger)
			assert.Equalf(t, tt.wantStderr, stderr.String(), "stderr HandleResponse(%v, %v, %v, %v, %v, %v, %v, %v)", tt.args.ctx, tt.args.i, tt.args.errIn(), tt.args.format, tt.args.fields, stdout, stderr, logger)
			assert.Equalf(t, tt.wantLogger, logger.String(), "logger HandleResponse(%v, %v, %v, %v, %v, %v, %v, %v)", tt.args.ctx, tt.args.i, tt.args.errIn(), tt.args.format, tt.args.fields, stdout, stderr, logger)
		})
	}
}
