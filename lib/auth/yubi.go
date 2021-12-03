package auth

import (
	"io"
	"syscall"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"golang.org/x/crypto/ssh/terminal"

	"fmt"
	"strings"
)

func YubiResponse(paramsSessionCreate files_sdk.SessionCreateParams, responseError files_sdk.ResponseError, out io.Writer) (files_sdk.SessionCreateParams, error) {
	fmt.Fprintf(out, "yubi: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return paramsSessionCreate, err
	}
	fmt.Fprintf(out, "\n")
	password := string(bytePassword)
	paramsSessionCreate.Otp = strings.Replace(password, "\n", "", -1)
	paramsSessionCreate.Password = ""
	paramsSessionCreate.PartialSessionId = responseError.Data.PartialSessionId

	return paramsSessionCreate, nil
}
