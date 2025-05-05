package lib

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
)

func SmsResponse(paramsSessionCreate files_sdk.SessionCreateParams, out io.Writer) (files_sdk.SessionCreateParams, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fprintf(out, "sms: ")
	text, _ := reader.ReadString('\n')
	paramsSessionCreate.Otp = strings.ReplaceAll(text, "\n", "")
	return paramsSessionCreate, nil
}
