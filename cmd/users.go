package cmd
import "github.com/spf13/cobra"
import (
         "github.com/Files-com/files-cli/lib"
         files_sdk "github.com/Files-com/files-sdk-go"
         "github.com/Files-com/files-sdk-go/user"
         "fmt"
         "os"
)

var (
      _ = files_sdk.Config{}
      _ = user.Client{}
      _ = lib.OnlyFields
      _ = fmt.Println
      _ = os.Exit
    )

var (
    Users = &cobra.Command{
      Use: "users [command]",
      Args:  cobra.ExactArgs(1),
      Run: func(cmd *cobra.Command, args []string) {},
    }
)
func UsersInit() {
  var fieldsList string
  paramsUserList := files_sdk.UserListParams{}
  var MaxPagesList int
  cmdList := &cobra.Command{
      Use:   "list",
      Short: "list",
      Long:  `list`,
      Args:  cobra.MinimumNArgs(0),
      Run: func(cmd *cobra.Command, args []string) {
        params := paramsUserList
        params.MaxPages = MaxPagesList
        it := user.List(params)

        lib.JsonMarshalIter(it, fieldsList)
      },
  }
        cmdList.Flags().IntVarP(&paramsUserList.Page, "page", "p", 0, "List Users")
        cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
        cmdList.Flags().StringVarP(&fieldsList, "fields", "f", "", "comma separated list of field names to include in response")
        Users.AddCommand(cmdList)
        var fieldsFind string
        paramsUserFind := files_sdk.UserFindParams{}
        cmdFind := &cobra.Command{
            Use:   "find",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := user.Find(paramsUserFind)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsFind)
            },
        }
        cmdFind.Flags().IntVarP(&paramsUserFind.Id, "id", "i", 0, "Show User")
        cmdFind.Flags().StringVarP(&fieldsFind, "fields", "f", "", "comma separated list of field names")
        Users.AddCommand(cmdFind)
        var fieldsCreate string
        paramsUserCreate := files_sdk.UserCreateParams{}
        cmdCreate := &cobra.Command{
            Use:   "create",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := user.Create(paramsUserCreate)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsCreate)
            },
        }
        cmdCreate.Flags().StringVarP(&paramsUserCreate.ChangePassword, "change-password", "w", "", "Create User")
        cmdCreate.Flags().StringVarP(&paramsUserCreate.ChangePasswordConfirmation, "change-password-confirmation", "c", "", "Create User")
        cmdCreate.Flags().StringVarP(&paramsUserCreate.Email, "email", "", "", "Create User")
        cmdCreate.Flags().StringVarP(&paramsUserCreate.GrantPermission, "grant-permission", "g", "", "Create User")
        cmdCreate.Flags().IntVarP(&paramsUserCreate.GroupId, "group-id", "", 0, "Create User")
        cmdCreate.Flags().StringVarP(&paramsUserCreate.GroupIds, "group-ids", "", "", "Create User")
        cmdCreate.Flags().StringVarP(&paramsUserCreate.Password, "password", "", "", "Create User")
        cmdCreate.Flags().StringVarP(&paramsUserCreate.PasswordConfirmation, "password-confirmation", "", "", "Create User")
        cmdCreate.Flags().StringVarP(&paramsUserCreate.AllowedIps, "allowed-ips", "a", "", "Create User")
        cmdCreate.Flags().StringVarP(&paramsUserCreate.AuthenticateUntil, "authenticate-until", "u", "", "Create User")
        cmdCreate.Flags().StringVarP(&paramsUserCreate.AuthenticationMethod, "authentication-method", "e", "", "Create User")
        cmdCreate.Flags().StringVarP(&paramsUserCreate.HeaderText, "header-text", "x", "", "Create User")
        cmdCreate.Flags().StringVarP(&paramsUserCreate.Language, "language", "", "", "Create User")
        cmdCreate.Flags().IntVarP(&paramsUserCreate.NotificationDailySendTime, "notification-daily-send-time", "y", 0, "Create User")
        cmdCreate.Flags().StringVarP(&paramsUserCreate.Name, "name", "", "", "Create User")
        cmdCreate.Flags().StringVarP(&paramsUserCreate.Notes, "notes", "o", "", "Create User")
        cmdCreate.Flags().IntVarP(&paramsUserCreate.PasswordValidityDays, "password-validity-days", "", 0, "Create User")
        cmdCreate.Flags().StringVarP(&paramsUserCreate.SslRequired, "ssl-required", "", "", "Create User")
        cmdCreate.Flags().IntVarP(&paramsUserCreate.SsoStrategyId, "sso-strategy-id", "", 0, "Create User")
        cmdCreate.Flags().StringVarP(&paramsUserCreate.TimeZone, "time-zone", "z", "", "Create User")
        cmdCreate.Flags().StringVarP(&paramsUserCreate.UserRoot, "user-root", "", "", "Create User")
        cmdCreate.Flags().StringVarP(&paramsUserCreate.Username, "username", "", "", "Create User")
        cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "f", "", "comma separated list of field names")
        Users.AddCommand(cmdCreate)
        var fieldsUnlock string
        paramsUserUnlock := files_sdk.UserUnlockParams{}
        cmdUnlock := &cobra.Command{
            Use:   "unlock",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := user.Unlock(paramsUserUnlock)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsUnlock)
            },
        }
        cmdUnlock.Flags().IntVarP(&paramsUserUnlock.Id, "id", "i", 0, "Unlock user who has been locked out due to failed logins")
        cmdUnlock.Flags().StringVarP(&fieldsUnlock, "fields", "f", "", "comma separated list of field names")
        Users.AddCommand(cmdUnlock)
        var fieldsResendWelcomeEmail string
        paramsUserResendWelcomeEmail := files_sdk.UserResendWelcomeEmailParams{}
        cmdResendWelcomeEmail := &cobra.Command{
            Use:   "resend-welcome-email",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := user.ResendWelcomeEmail(paramsUserResendWelcomeEmail)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsResendWelcomeEmail)
            },
        }
        cmdResendWelcomeEmail.Flags().IntVarP(&paramsUserResendWelcomeEmail.Id, "id", "i", 0, "Resend user welcome email")
        cmdResendWelcomeEmail.Flags().StringVarP(&fieldsResendWelcomeEmail, "fields", "f", "", "comma separated list of field names")
        Users.AddCommand(cmdResendWelcomeEmail)
        var fieldsUser2faReset string
        paramsUserUser2faReset := files_sdk.UserUser2faResetParams{}
        cmdUser2faReset := &cobra.Command{
            Use:   "user-2fa-reset",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := user.User2faReset(paramsUserUser2faReset)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsUser2faReset)
            },
        }
        cmdUser2faReset.Flags().IntVarP(&paramsUserUser2faReset.Id, "id", "i", 0, "Trigger 2FA Reset process for user who has lost access to their existing 2FA methods")
        cmdUser2faReset.Flags().StringVarP(&fieldsUser2faReset, "fields", "f", "", "comma separated list of field names")
        Users.AddCommand(cmdUser2faReset)
        var fieldsUpdate string
        paramsUserUpdate := files_sdk.UserUpdateParams{}
        cmdUpdate := &cobra.Command{
            Use:   "update",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := user.Update(paramsUserUpdate)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsUpdate)
            },
        }
        cmdUpdate.Flags().IntVarP(&paramsUserUpdate.Id, "id", "", 0, "Update User")
        cmdUpdate.Flags().StringVarP(&paramsUserUpdate.ChangePassword, "change-password", "w", "", "Update User")
        cmdUpdate.Flags().StringVarP(&paramsUserUpdate.ChangePasswordConfirmation, "change-password-confirmation", "c", "", "Update User")
        cmdUpdate.Flags().StringVarP(&paramsUserUpdate.Email, "email", "", "", "Update User")
        cmdUpdate.Flags().StringVarP(&paramsUserUpdate.GrantPermission, "grant-permission", "g", "", "Update User")
        cmdUpdate.Flags().IntVarP(&paramsUserUpdate.GroupId, "group-id", "", 0, "Update User")
        cmdUpdate.Flags().StringVarP(&paramsUserUpdate.GroupIds, "group-ids", "", "", "Update User")
        cmdUpdate.Flags().StringVarP(&paramsUserUpdate.Password, "password", "", "", "Update User")
        cmdUpdate.Flags().StringVarP(&paramsUserUpdate.PasswordConfirmation, "password-confirmation", "", "", "Update User")
        cmdUpdate.Flags().StringVarP(&paramsUserUpdate.AllowedIps, "allowed-ips", "a", "", "Update User")
        cmdUpdate.Flags().StringVarP(&paramsUserUpdate.AuthenticateUntil, "authenticate-until", "u", "", "Update User")
        cmdUpdate.Flags().StringVarP(&paramsUserUpdate.AuthenticationMethod, "authentication-method", "e", "", "Update User")
        cmdUpdate.Flags().StringVarP(&paramsUserUpdate.HeaderText, "header-text", "x", "", "Update User")
        cmdUpdate.Flags().StringVarP(&paramsUserUpdate.Language, "language", "", "", "Update User")
        cmdUpdate.Flags().IntVarP(&paramsUserUpdate.NotificationDailySendTime, "notification-daily-send-time", "y", 0, "Update User")
        cmdUpdate.Flags().StringVarP(&paramsUserUpdate.Name, "name", "", "", "Update User")
        cmdUpdate.Flags().StringVarP(&paramsUserUpdate.Notes, "notes", "o", "", "Update User")
        cmdUpdate.Flags().IntVarP(&paramsUserUpdate.PasswordValidityDays, "password-validity-days", "", 0, "Update User")
        cmdUpdate.Flags().StringVarP(&paramsUserUpdate.SslRequired, "ssl-required", "", "", "Update User")
        cmdUpdate.Flags().IntVarP(&paramsUserUpdate.SsoStrategyId, "sso-strategy-id", "", 0, "Update User")
        cmdUpdate.Flags().StringVarP(&paramsUserUpdate.TimeZone, "time-zone", "z", "", "Update User")
        cmdUpdate.Flags().StringVarP(&paramsUserUpdate.UserRoot, "user-root", "", "", "Update User")
        cmdUpdate.Flags().StringVarP(&paramsUserUpdate.Username, "username", "", "", "Update User")
        cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "f", "", "comma separated list of field names")
        Users.AddCommand(cmdUpdate)
        var fieldsDelete string
        paramsUserDelete := files_sdk.UserDeleteParams{}
        cmdDelete := &cobra.Command{
            Use:   "delete",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := user.Delete(paramsUserDelete)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsDelete)
            },
        }
        cmdDelete.Flags().IntVarP(&paramsUserDelete.Id, "id", "i", 0, "Delete User")
        cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "f", "", "comma separated list of field names")
        Users.AddCommand(cmdDelete)
}
