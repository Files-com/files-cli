package auth

import (
	"bufio"
	"os"

	files_sdk "github.com/Files-com/files-sdk-go"

	"fmt"
	"strings"
)

func TotpResponse(paramsSessionCreate files_sdk.SessionCreateParams) (files_sdk.SessionCreateParams, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("totp: ")
	text, _ := reader.ReadString('\n')
	paramsSessionCreate.Otp = strings.Replace(text, "\n", "", -1)
	return paramsSessionCreate, nil
}
