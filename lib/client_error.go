package lib

import (
	"context"
	"fmt"
	"os"

	files_sdk "github.com/Files-com/files-sdk-go"
)

func ClientError(ctx context.Context, err error) {
	if err == nil {
		return
	}
	responseError, ok := err.(files_sdk.ResponseError)
	config := Config{}
	config.Load()

	if ok && responseError.Type == "not-authorized/authentication-required" && config.ValidSession() {
		fmt.Println("Your session is invalid. Please login with:")
		fmt.Println("\tfiles-cli login")
	}

	if ok && responseError.Type == "not-authorized/authentication-required" && files_sdk.GlobalConfig.GetAPIKey() != "" {
		fmt.Println("Set the correct api-key with: ")
		fmt.Println("\tfiles-cli config set --api-key=")
	}

	if ok {
		fmt.Println(responseError)
	} else {
		fmt.Println(err)
	}

	if ctx.Value("testing") != nil && !ctx.Value("testing").(bool) {
		os.Exit(1)
	}
}
