//go:build !noportable
// +build !noportable

package cmd

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/drakkan/sftpgo/v2/config"

	"github.com/sftpgo/sdk"
	"github.com/spf13/cobra"

	"github.com/drakkan/sftpgo/v2/common"
	"github.com/drakkan/sftpgo/v2/dataprovider"
	"github.com/drakkan/sftpgo/v2/kms"
	"github.com/drakkan/sftpgo/v2/service"
	"github.com/drakkan/sftpgo/v2/vfs"
)

const (
	logFilePathFlag     = "log-file-path"
	logVerboseFlag      = "log-verbose"
	logUTCTimeFlag      = "log-utc-time"
	defaultConfigDir    = "."
	defaultConfigFile   = ""
	defaultLogMaxSize   = 10
	defaultLogMaxBackup = 5
	defaultLogMaxAge    = 28
	defaultLogCompress  = false
)

var (
	ConfigDir                          string
	directoryToServe                   string
	portableSFTPDPort                  int
	portableAdvertiseService           bool
	portableAdvertiseCredentials       bool
	portableUsername                   string
	portablePassword                   string
	portableStartDir                   string
	portableLogFile                    string
	portableLogVerbose                 bool
	portableLogUTCTime                 bool
	portablePublicKeys                 []string
	portablePermissions                []string
	portableSSHCommands                []string
	portableAllowedPatterns            []string
	portableDeniedPatterns             []string
	portableFsProvider                 string
	portableFTPDPort                   int
	portableFTPSCert                   string
	portableFTPSKey                    string
	portableCryptPassphrase            string
	portableSFTPEndpoint               string
	portableSFTPUsername               string
	portableSFTPPassword               string
	portableSFTPPrivateKeyPath         string
	portableSFTPFingerprints           []string
	portableSFTPPrefix                 string
	portableSFTPDisableConcurrentReads bool
	portableSFTPDBufferSize            int64
	AgentCmd                           = &cobra.Command{
		Use:   "agent",
		Short: "Start Files.com Agent",
		Long: `To serve the current working directory with auto generated credentials simply
use:

$ files-cli agent

Please take a look at the usage below to customize the serving parameters`,
		Run: func(cmd *cobra.Command, args []string) {
			portableDir := directoryToServe
			fsProvider := sdk.GetProviderByName(portableFsProvider)
			if !filepath.IsAbs(portableDir) {
				if fsProvider == sdk.LocalFilesystemProvider {
					portableDir, _ = filepath.Abs(portableDir)
				} else {
					portableDir = os.TempDir()
				}
			}
			permissions := make(map[string][]string)
			permissions["/"] = portablePermissions
			portableSFTPPrivateKey := ""
			if portableSFTPPrivateKeyPath != "" {
				contents, err := getFileContents(portableSFTPPrivateKeyPath)
				if err != nil {
					fmt.Printf("Unable to get SFTP private key: %v\n", err)
					os.Exit(1)
				}
				portableSFTPPrivateKey = contents
			}
			if portableFTPDPort >= 0 && portableFTPSCert != "" && portableFTPSKey != "" {
				keyPairs := []common.TLSKeyPair{
					{
						Cert: portableFTPSCert,
						Key:  portableFTPSKey,
						ID:   common.DefaultTLSKeyPaidID,
					},
				}
				_, err := common.NewCertManager(keyPairs, filepath.Clean(defaultConfigDir),
					"FTP portable")
				if err != nil {
					fmt.Printf("Unable to load FTPS key pair, cert file %#v key file %#v error: %v\n",
						portableFTPSCert, portableFTPSKey, err)
					os.Exit(1)
				}
			}
			service := service.Service{
				ConfigDir:     filepath.Clean(defaultConfigDir),
				ConfigFile:    defaultConfigFile,
				LogFilePath:   portableLogFile,
				LogMaxSize:    defaultLogMaxSize,
				LogMaxBackups: defaultLogMaxBackup,
				LogMaxAge:     defaultLogMaxAge,
				LogCompress:   defaultLogCompress,
				LogVerbose:    portableLogVerbose,
				LogUTCTime:    portableLogUTCTime,
				Shutdown:      make(chan bool),
				PortableMode:  1,
				PortableUser: dataprovider.User{
					BaseUser: sdk.BaseUser{
						Username:    portableUsername,
						Password:    portablePassword,
						PublicKeys:  portablePublicKeys,
						Permissions: permissions,
						HomeDir:     portableDir,
						Status:      1,
					},
					Filters: dataprovider.UserFilters{
						BaseUserFilters: sdk.BaseUserFilters{
							FilePatterns:   parsePatternsFilesFilters(),
							StartDirectory: portableStartDir,
						},
					},
					FsConfig: vfs.Filesystem{
						Provider: sdk.GetProviderByName(portableFsProvider),
						CryptConfig: vfs.CryptFsConfig{
							Passphrase: kms.NewPlainSecret(portableCryptPassphrase),
						},
						SFTPConfig: vfs.SFTPFsConfig{
							BaseSFTPFsConfig: sdk.BaseSFTPFsConfig{
								Endpoint:                portableSFTPEndpoint,
								Username:                portableSFTPUsername,
								Fingerprints:            portableSFTPFingerprints,
								Prefix:                  portableSFTPPrefix,
								DisableCouncurrentReads: portableSFTPDisableConcurrentReads,
								BufferSize:              portableSFTPDBufferSize,
							},
							Password:   kms.NewPlainSecret(portableSFTPPassword),
							PrivateKey: kms.NewPlainSecret(portableSFTPPrivateKey),
						},
					},
				},
			}
			if err := service.StartPortableMode(portableSFTPDPort, portableFTPDPort, 0, portableSSHCommands, portableAdvertiseService,
				portableAdvertiseCredentials, "", "", "", ""); err == nil {
				service.Wait()
				if service.Error == nil {
					os.Exit(0)
				}
			}
			os.Exit(1)
		},
	}
)

func init() {
	AgentCmd.Flags().StringVarP(&directoryToServe, "directory", "d", ".", `Path to the directory to serve.
This can be an absolute path or a path
relative to the current directory
`)
	AgentCmd.Flags().StringVar(&portableStartDir, "start-directory", "/", `Alternate start directory.
This is a virtual path not a filesystem
path`)
	AgentCmd.Flags().IntVarP(&portableSFTPDPort, "sftpd-port", "s", 58550, `0 means a random unprivileged port,
< 0 disabled`)
	AgentCmd.Flags().StringVarP(&portableUsername, "username", "u", "admin", `Leave empty to use an auto generated
value`)
	AgentCmd.Flags().StringVarP(&portableLogFile, logFilePathFlag, "l", "", "Leave empty to disable logging")
	AgentCmd.Flags().BoolVarP(&portableLogVerbose, logVerboseFlag, "v", false, "Enable verbose logs")
	AgentCmd.Flags().BoolVar(&portableLogUTCTime, logUTCTimeFlag, false, "Use UTC time for logging")
	AgentCmd.Flags().StringSliceVarP(&portablePublicKeys, "public-key", "k", []string{}, "")
	AgentCmd.Flags().StringSliceVarP(&portablePermissions, "permissions", "g", []string{"list", "download", "upload", "overwrite", "delete", "rename", "create_dirs"},
		`User's permissions. "*" means any
permission`)
	AgentCmd.Flags().StringArrayVar(&portableAllowedPatterns, "allowed-patterns", []string{},
		`Allowed file patterns case insensitive.
The format is:
/dir::pattern1,pattern2.
For example: "/somedir::*.jpg,a*b?.png"`)
	AgentCmd.Flags().StringArrayVar(&portableDeniedPatterns, "denied-patterns", []string{},
		`Denied file patterns case insensitive.
The format is:
/dir::pattern1,pattern2.
For example: "/somedir::*.jpg,a*b?.png"`)
	AgentCmd.Flags().BoolVarP(&portableAdvertiseService, "advertise-service", "S", false,
		`Advertise configured services using
multicast DNS`)
	AgentCmd.Flags().BoolVarP(&portableAdvertiseCredentials, "advertise-credentials", "C", false,
		`If the SFTP/FTP service is
advertised via multicast DNS, this
flag allows to put username/password
inside the advertised TXT record`)
	portableFsProvider = "osfs" // local filesystem
	AgentCmd.Flags().StringVar(&portableSFTPEndpoint, "sftp-endpoint", "", `SFTP endpoint as host:port for SFTP
provider`)
	AgentCmd.Flags().StringVar(&portableSFTPUsername, "sftp-username", "", `SFTP user for SFTP provider`)
	AgentCmd.Flags().StringVar(&portableSFTPPassword, "sftp-password", "", `SFTP password for SFTP provider`)
	AgentCmd.Flags().StringVar(&portableSFTPPrivateKeyPath, "sftp-key-path", "", `SFTP private key path for SFTP provider`)
	AgentCmd.Flags().StringSliceVar(&portableSFTPFingerprints, "sftp-fingerprints", []string{}, `SFTP fingerprints to verify remote host
key for SFTP provider`)
	AgentCmd.Flags().StringVar(&portableSFTPPrefix, "sftp-prefix", "", `SFTP prefix allows restrict all
operations to a given path within the
remote SFTP server`)
	AgentCmd.Flags().BoolVar(&portableSFTPDisableConcurrentReads, "sftp-disable-concurrent-reads", false, `Concurrent reads are safe to use and
disabling them will degrade performance.
Disable for read once servers`)
	AgentCmd.Flags().Int64Var(&portableSFTPDBufferSize, "sftp-buffer-size", 0, `The size of the buffer (in MB) to use
for transfers. By enabling buffering,
the reads and writes, from/to the
remote SFTP server, are split in
multiple concurrent requests and this
allows data to be transferred at a
faster rate, over high latency networks,
by overlapping round-trip times`)
}

func parsePatternsFilesFilters() []sdk.PatternsFilter {
	var patterns []sdk.PatternsFilter
	for _, val := range portableAllowedPatterns {
		p, exts := getPatternsFilterValues(strings.TrimSpace(val))
		if p != "" {
			patterns = append(patterns, sdk.PatternsFilter{
				Path:            path.Clean(p),
				AllowedPatterns: exts,
				DeniedPatterns:  []string{},
			})
		}
	}
	for _, val := range portableDeniedPatterns {
		p, exts := getPatternsFilterValues(strings.TrimSpace(val))
		if p != "" {
			found := false
			for index, e := range patterns {
				if path.Clean(e.Path) == path.Clean(p) {
					patterns[index].DeniedPatterns = append(patterns[index].DeniedPatterns, exts...)
					found = true
					break
				}
			}
			if !found {
				patterns = append(patterns, sdk.PatternsFilter{
					Path:            path.Clean(p),
					AllowedPatterns: []string{},
					DeniedPatterns:  exts,
				})
			}
		}
	}
	return patterns
}

func getPatternsFilterValues(value string) (string, []string) {
	if strings.Contains(value, "::") {
		dirExts := strings.Split(value, "::")
		if len(dirExts) > 1 {
			dir := strings.TrimSpace(dirExts[0])
			exts := []string{}
			for _, e := range strings.Split(dirExts[1], ",") {
				cleanedExt := strings.TrimSpace(e)
				if cleanedExt != "" {
					exts = append(exts, cleanedExt)
				}
			}
			if dir != "" && len(exts) > 0 {
				return dir, exts
			}
		}
	}
	return "", nil
}

func getFileContents(name string) (string, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return "", err
	}
	if fi.Size() > 1048576 {
		return "", fmt.Errorf("%#v is too big %v/1048576 bytes", name, fi.Size())
	}
	contents, err := os.ReadFile(name)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}

func init() {
	AgentCmd.CompletionOptions.DisableDefaultCmd = true
	RootCmd.AddCommand(AgentCmd)
	AgentCmd.Hidden = true
	ConfigDir = os.Getenv("FILES_AGENT_CONFIG_DIR")
	if ConfigDir == "" {
		ConfigDir = "./files-agent"
	}
	config.LoadConfig(ConfigDir, "config.json")
	publicKeys := []string{"id_ed25519.pub"}

	for _, key := range publicKeys {
		keyFile, err := os.Open(filepath.Join(ConfigDir, key))
		if err == nil {
			buf, err := io.ReadAll(keyFile)
			if err == nil {
				portableSFTPFingerprints = append(portableSFTPFingerprints, string(buf))
				portablePublicKeys = append(portablePublicKeys, string(buf))
			}
		}
	}
	portableSFTPPrivateKeyPath = filepath.Join(ConfigDir, "id_ed25519")
}
