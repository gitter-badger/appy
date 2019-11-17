package cmd

import (
	"os"
	"os/exec"

	appyhttp "github.com/appist/appy/internal/http"
	appysupport "github.com/appist/appy/internal/support"
)

// NewSSLSetupCommand generate and install the locally trusted SSL certs using "mkcert".
func NewSSLSetupCommand(logger *appysupport.Logger, s *appyhttp.Server) *Command {
	return &Command{
		Use:   "ssl:setup",
		Short: `Generate and install the locally trusted SSL certs using "mkcert"`,
		Run: func(cmd *Command, args []string) {
			if appysupport.IsConfigErrored(s.Config(), logger) {
				os.Exit(-1)
			}

			_, err := exec.LookPath("mkcert")
			if err != nil {
				logger.Fatal(err)
			}

			os.MkdirAll(s.Config().HTTPSSLCertPath, os.ModePerm)
			setupArgs := []string{"-install", "-cert-file", s.Config().HTTPSSLCertPath + "/cert.pem", "-key-file", s.Config().HTTPSSLCertPath + "/key.pem"}
			hosts, _ := s.Hosts()
			setupArgs = append(setupArgs, hosts...)
			setupCmd := exec.Command("mkcert", setupArgs...)
			setupCmd.Stdout = os.Stdout
			setupCmd.Stderr = os.Stderr
			setupCmd.Run()
		},
	}
}