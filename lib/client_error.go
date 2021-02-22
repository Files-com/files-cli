package lib

import (
	"fmt"
	"os"

	files_sdk "github.com/Files-com/files-sdk-go"
)

func ClientError(err error) {
	responseError, ok := err.(files_sdk.ResponseError)
	config := Config{}
	config.Load()
	fmt.Println(responseError)
	if ok && responseError.Type == "not-authorized/authentication-required" && config.ValidSession() {
		fmt.Println("Your session is invalid. Please login with:")
		fmt.Println("\tfiles-cli login")
	}

	if ok && responseError.Type == "not-authorized/authentication-required" && files_sdk.GlobalConfig.GetAPIKey() != "" {
		fmt.Println("Set the correct api-key with: ")
		fmt.Println("\tfiles-cli config set --api-key=")
	}

	os.Exit(1)
}
