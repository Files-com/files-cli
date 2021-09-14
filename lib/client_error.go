package lib

import (
	"context"
	"fmt"
	"os"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
)

func ClientError(ctx context.Context, err error) {
	if err == nil {
		return
	}
	responseError, ok := err.(files_sdk.ResponseError)
	config := Config{}
	config.Load()

	if ok && responseError.Type == "not-authorized/authentication-required" && config.ValidSession() {
		fmt.Fprintf(os.Stderr, "Your session is invalid. Please login with:\n")
		fmt.Fprintf(os.Stderr, "\tfiles-cli login\n")
	}

	if ok && responseError.Type == "not-authorized/authentication-required" && files_sdk.GlobalConfig.GetAPIKey() != "" {
		fmt.Fprintf(os.Stderr, "Set the correct api-key with:\n")
		fmt.Fprintf(os.Stderr, "\tfiles-cli config set --api-key=\n")
	}

	if ok {
		fmt.Fprintf(os.Stderr, "%v\n", responseError)
	} else {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	if ctx.Value("testing") != nil && !ctx.Value("testing").(bool) {
		os.Exit(1)
	}
}
