package lib

import (
	"bufio"
	"io"
	"os"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"
	"strings"
)

func TotpResponse(paramsSessionCreate files_sdk.SessionCreateParams, out io.Writer) (files_sdk.SessionCreateParams, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fprintf(out, "totp: ")
	text, _ := reader.ReadString('\n')
	paramsSessionCreate.Otp = strings.Replace(text, "\n", "", -1)
	return paramsSessionCreate, nil
}