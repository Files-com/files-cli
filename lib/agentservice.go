package lib

import (
	"context"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/spf13/pflag"
)

const (
	logFilePathFlag = "log-file-path"
	logVerboseFlag  = "log-verbose"
	logUTCTimeFlag  = "log-utc-time"
)

type AgentService struct {
	files_sdk.Config
	ConfigPath string
	files_sdk.RemoteServerConfigurationFile
	portableFsProvider                 string
	permissions                        map[string][]string
	portableSFTPFingerprints           []string
	portablePublicKeys                 []string
	shutdown                           chan bool
	PortableSFTPEndpoint               string
	PortableSFTPPrefix                 string
	PortableSFTPDisableConcurrentReads bool
	PortableSFTPDBufferSize            int64
	PortableAllowedPatterns            []string
	PortableDeniedPatterns             []string
	PortableLogFile                    string
	PortableLogVerbose                 bool
	PortableLogUTCTime                 bool
	PortableSFTPFingerprints           []string
	ipWhitelist                        map[string]bool
	appendedIpWhitelist                []string
	sftpGoConfigPath                   string
	subdomain                          string
	NoSourceIpCheck                    bool
	context.Context
}

func (a *AgentService) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&a.subdomain, "subdomain", "", "")
	flags.StringVar(&a.ConfigPath, "config", "./files-agent-config.json", "")
	flags.StringVarP(&a.PortableLogFile, logFilePathFlag, "l", "files-agent.log", "Leave empty to disable logging")
	flags.BoolVarP(&a.PortableLogVerbose, logVerboseFlag, "v", false, "Enable verbose logs")
	flags.BoolVar(&a.PortableLogUTCTime, logUTCTimeFlag, false, "Use UTC time for logging")
	flags.StringArrayVar(&a.PortableAllowedPatterns, "allowed-patterns", []string{},
		`Allowed file patterns case insensitive.
The format is:
/dir::pattern1,pattern2.
For example: "/somedir::*.jpg,a*b?.png"`)
	flags.StringArrayVar(&a.PortableDeniedPatterns, "denied-patterns", []string{},
		`Denied file patterns case insensitive.
The format is:
/dir::pattern1,pattern2.
For example: "/somedir::*.jpg,a*b?.png"`)
	flags.StringVar(&a.PortableSFTPEndpoint, "sftp-endpoint", "", `SFTP endpoint as host:port for SFTP
provider`)
	flags.StringSliceVar(&a.PortableSFTPFingerprints, "sftp-fingerprints", []string{}, `SFTP fingerprints to verify remote host
key for SFTP provider`)
	flags.StringVar(&a.PortableSFTPPrefix, "sftp-prefix", "", `SFTP prefix allows restrict all
operations to a given path within the
remote SFTP server`)
	flags.Int64Var(&a.PortableSFTPDBufferSize, "sftp-buffer-size", 0, `The size of the buffer (in MB) to use
for transfers. By enabling buffering,
the reads and writes, from/to the
remote SFTP server, are split in
multiple concurrent requests and this
allows data to be transferred at a
faster rate, over high latency networks,
by overlapping round-trip times`)
	flags.StringSliceVar(&a.appendedIpWhitelist, "append-to-ip-whitelist", []string{"127.0.0.1"}, "Add additional IPs to whitelist")
	flags.BoolVar(&a.NoSourceIpCheck, "no-source-ip-check", false, "Disable source IP check")
}
