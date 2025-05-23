package cmd

import (
	"encoding/json"
	"fmt"
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsers(t *testing.T) {
	user, err := findOrCreateTestUser(files_sdk.UserCreateParams{Username: "test-automations-user", Password: "Foo123456789!"})
	require.NoError(t, err)
	defer deleteUser(user)

	tests := []struct {
		name string
		args []string
		want files_sdk.User
	}{
		{
			name: "update user to authentication-method sso",
			args: []string{"update", "--authentication-method", "sso", "--sso-strategy-id", "12"},
			want: files_sdk.User{SsoStrategyId: 12, AuthenticationMethod: "sso"},
		},
		{
			name: "update user to authentication-method password",
			args: []string{"update", "--authentication-method", "password", "--sso-strategy-id", "0"},
			want: files_sdk.User{SsoStrategyId: 0, AuthenticationMethod: "password"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, config, err := CreateConfig(t.Name())
			if err != nil {
				t.Fatal(err)
			}

			args := append(tt.args, "--format", "json,raw", "--id", fmt.Sprintf("%v", user.Id))
			t.Log(args)
			stdOut, stdErr := callCmd(Users(), config, args)

			j, err := json.Marshal(&tt.want)
			require.NoError(t, err)
			var wantMap map[string]interface{}
			err = json.Unmarshal(j, &wantMap)
			require.NoError(t, err)

			var resultMap map[string]interface{}
			err = json.Unmarshal(stdOut, &resultMap)
			require.NoError(t, err)

			for k, v := range wantMap {
				if v == nil {
					delete(resultMap, k)
					delete(wantMap, k)
				}
			}
			for k, v := range resultMap {
				_, ok := wantMap[k]
				if v == nil || !ok {
					delete(resultMap, k)
				}
			}
			if len(wantMap) != 0 && len(resultMap) != 0 {
				assert.Equal(t, wantMap, resultMap)
			}
			assert.Equal(t, string(stdErr), "")

			err = json.Unmarshal(stdOut, &user)
			require.NoError(t, err)
			r.Stop()
		})
	}
}

func findOrCreateTestUser(p files_sdk.UserCreateParams) (files_sdk.User, error) {
	r, config, err := CreateConfig("findOrCreateTestUser")
	if err != nil {
		return files_sdk.User{}, err
	}

	defer r.Stop()

	client := user.Client{Config: config}
	it, err := client.List(files_sdk.UserListParams{})
	if err != nil {
		return files_sdk.User{}, err
	}

	for it.Next() {
		if it.Err() != nil {
			return files_sdk.User{}, err
		}

		if it.User().Username == p.Username {
			return it.User(), nil
		}
	}

	return client.Create(p)
}

func deleteUser(u files_sdk.User) error {
	r, config, err := CreateConfig("deleteUsers")
	if err != nil {
		return err
	}

	defer r.Stop()

	client := user.Client{Config: config}

	return client.Delete(files_sdk.UserDeleteParams{Id: u.Id})
}
