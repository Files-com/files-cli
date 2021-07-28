package lib

import (
	"bufio"
	"context"
	"encoding/json"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/Files-com/files-cli/lib/auth"

	"fmt"
	"os"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/session"
	"golang.org/x/crypto/ssh/terminal"
)

type Config struct {
	SessionId          string    `json:"session_id"`
	SessionExpiry      time.Time `json:"session_expiry"`
	Subdomain          string    `json:"subdomain"`
	Username           string    `json:"username"`
	APIKey             string    `json:"api_key"`
	Endpoint           string    `json:"endpoint,omitempty"`
	configPathOverride string    `json:"-"`
}

type ResetConfig struct {
	Subdomain bool
	Username  bool
	APIKey    bool
	Endpoint  bool
	Session   bool
}

var SessionExpiry = time.Hour * 6

func (c Config) ResetWith(reset ResetConfig) error {
	if reset.Subdomain == true {
		c.Subdomain = ""
	}
	if reset.Username == true {
		c.Username = ""
	}
	if reset.APIKey == true {
		c.APIKey = ""
	}
	if reset.Endpoint == true {
		c.Endpoint = ""
	}
	if reset.Session == true {
		c.SessionId = ""
	}
	return c.Save()
}

func (c Config) Reset() error {
	return initConfig()
}

func (c *Config) Load() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}
	configRoot := filepath.Join(usr.HomeDir, ".config")
	configPath := filepath.Join(configRoot, "files-cli")
	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		initConfig()
		return nil
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &c)

	if err != nil {
		return err
	}

	c.SetGlobal()
	return nil
}

func (c *Config) SetGlobal() {
	files_sdk.GlobalConfig.SessionId = c.SessionId
	files_sdk.GlobalConfig.Subdomain = c.Subdomain
	files_sdk.GlobalConfig.APIKey = c.APIKey
	files_sdk.GlobalConfig.Endpoint = c.Endpoint
}

func (c *Config) Save() error {
	file, _ := json.MarshalIndent(c, "", " ")
	path, err := c.configPath()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, file, 0644)
}

func (c *Config) ValidSession() bool {
	return c.SessionId != "" && !c.SessionExpired()
}

func (c *Config) SessionExpired() bool {
	return c.SessionId != "" && time.Now().Local().After(c.SessionExpiry)
}

func (c *Config) configPath() (string, error) {
	if c.configPathOverride != "" {
		return c.configPathOverride, nil
	}
	root, err := configRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "files-cli"), nil
}

func configRoot() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, ".config"), nil
}

func initConfig() error {
	root, err := configRoot()
	if err != nil {
		return err
	}
	_, err = os.Stat(root)
	if os.IsNotExist(err) {
		os.MkdirAll(root, 0600)
	}

	path, err := (&Config{}).configPath()
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	f.Write([]byte("{}"))
	f.Close()
	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func SessionUnauthorizedError(paramsSessionCreate files_sdk.SessionCreateParams, err error) (files_sdk.SessionCreateParams, error) {
	responseError, ok := err.(files_sdk.ResponseError)
	if !ok {
		return paramsSessionCreate, err
	}
	switch eType := responseError.Type; eType {
	case "not-authorized/authentication-required":
		return paramsSessionCreate, err
	case "bad-request/invalid-username-or-password":
		return paramsSessionCreate, err
	case "not-authenticated/two-factor-authentication-error":
		fmt.Println(responseError.ErrorMessage)

		if contains(responseError.Data.TwoFactorAuthenticationMethod, "u2f") {
			return auth.U2fResponse(paramsSessionCreate, responseError)
		}

		if contains(responseError.Data.TwoFactorAuthenticationMethod, "yubi") {
			return auth.YubiResponse(paramsSessionCreate, responseError)
		}

		if contains(responseError.Data.TwoFactorAuthenticationMethod, "totp") {
			return auth.TotpResponse(paramsSessionCreate)
		}

		if contains(responseError.Data.TwoFactorAuthenticationMethod, "sms") {
			return auth.SmsResponse(paramsSessionCreate)
		}

		return paramsSessionCreate, fmt.Errorf("%v is unsupported as login method", responseError.Data.TwoFactorAuthenticationMethod)
	}

	return paramsSessionCreate, err
}

func stringInput(reader *bufio.Reader, label string) string {
	fmt.Printf("%v: ", label)
	text, _ := reader.ReadString('\n')
	return parseTermInput(text)
}

func parseTermInput(text string) string {
	text = strings.Replace(text, "\r", "", -1)
	text = strings.Replace(text, "\n", "", -1) // Windows command prompt
	return text
}

func CreateSession(paramsSessionCreate files_sdk.SessionCreateParams, config Config) error {
	reader := bufio.NewReader(os.Stdin)

	if config.Subdomain == "" {
		config.Subdomain = stringInput(reader, "Subdomain")
	} else {
		fmt.Printf("Subdomain: %v\n", config.Subdomain)
	}

	if paramsSessionCreate.Username == "" && config.Username == "" {
		paramsSessionCreate.Username = stringInput(reader, "Username")
		config.Username = paramsSessionCreate.Username
	} else {
		if paramsSessionCreate.Username != "" {
			config.Username = paramsSessionCreate.Username
		} else {
			paramsSessionCreate.Username = config.Username
		}
		fmt.Printf("Username: %v\n", paramsSessionCreate.Username)
	}

	if paramsSessionCreate.Password == "" {
		fmt.Print("Password: ")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}

		paramsSessionCreate.Password = parseTermInput(string(bytePassword))
		fmt.Println("")
	}

	config.SetGlobal()
	files_sdk.GlobalConfig.SessionId = ""
	client := session.Client{Config: files_sdk.GlobalConfig}

	result, err := client.Create(context.TODO(), paramsSessionCreate)

	if err != nil {
		otpSessionCreate, err := SessionUnauthorizedError(paramsSessionCreate, err)
		if err == nil {
			result, err = client.Create(context.TODO(), otpSessionCreate)
		}

		if err != nil {
			return err
		}
	}
	config.SessionId = result.Id
	config.SessionExpiry = time.Now().Local().Add(SessionExpiry)

	err = config.Save()
	if err != nil {
		return err
	}
	config.SetGlobal()
	return nil
}
