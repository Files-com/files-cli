package lib

import (
	"context"
	"fmt"
	"io"
	"os"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
)

func ClientError(ctx context.Context, err error, out ...io.Writer) {
	if len(out) == 0 {
		out = append(out, os.Stdout)
	}
	if err == nil {
		return
	}
	responseError, ok := err.(files_sdk.ResponseError)
	config := Config{}
	config.Load()

	if ok && responseError.Type == "not-authorized/authentication-required" && config.ValidSession() {
		fmt.Fprintf(out[0], "Your session is invalid. Please login with:\n")
		fmt.Fprintf(out[0], "\tfiles-cli login\n")
	}

	if ok && responseError.Type == "not-authorized/authentication-required" && files_sdk.GlobalConfig.GetAPIKey() != "" {
		fmt.Fprintf(out[0], "Set the correct api-key with:\n")
		fmt.Fprintf(out[0], "\tfiles-cli config set --api-key=\n")
	}

	if ok {
		fmt.Fprintf(out[0], "%v\n", responseError)
	} else {
		fmt.Fprintf(out[0], "%v\n", err)
	}

	if ctx.Value("testing") != nil && !ctx.Value("testing").(bool) {
		os.Exit(1)
	}
}
