package auth

import (
	"syscall"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"golang.org/x/crypto/ssh/terminal"

	"fmt"
	"strings"
)

func YubiResponse(paramsSessionCreate files_sdk.SessionCreateParams, responseError files_sdk.ResponseError) (files_sdk.SessionCreateParams, error) {
	fmt.Print("yubi: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return paramsSessionCreate, err
	}
	fmt.Println("")
	password := string(bytePassword)
	paramsSessionCreate.Otp = strings.Replace(password, "\n", "", -1)
	paramsSessionCreate.Password = ""
	paramsSessionCreate.PartialSessionId = responseError.Data.PartialSessionId

	return paramsSessionCreate, nil
}
