package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/inhandnet/elements-cli/internal/api"
	cmd "github.com/inhandnet/elements-cli/internal/cmd"
	authCmd "github.com/inhandnet/elements-cli/internal/cmd/auth"
	configCmd "github.com/inhandnet/elements-cli/internal/cmd/config"
	deviceCmd "github.com/inhandnet/elements-cli/internal/cmd/device"
	versionCmd "github.com/inhandnet/elements-cli/internal/cmd/version"
	"github.com/inhandnet/elements-cli/internal/factory"
)

func main() {
	f := factory.New()
	rootCmd := cmd.NewCmdRoot(f)
	rootCmd.AddCommand(authCmd.NewCmdAuth(f))
	rootCmd.AddCommand(configCmd.NewCmdConfig(f))
	rootCmd.AddCommand(deviceCmd.NewCmdDevice(f))
	rootCmd.AddCommand(versionCmd.NewCmdVersion(f))

	// Top-level shortcut: `elements login` → `elements auth login`
	loginAlias := authCmd.NewCmdLogin(f)
	loginAlias.Use = "login"
	loginAlias.Hidden = true
	rootCmd.AddCommand(loginAlias)

	executedCmd, err := rootCmd.ExecuteC()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		if isFlagError(err) && executedCmd != nil {
			fmt.Fprintln(os.Stderr)
			fmt.Fprint(os.Stderr, executedCmd.UsageString())
		}
		var httpErr *api.HTTPError
		if errors.As(err, &httpErr) && httpErr.StatusCode == 401 {
			fmt.Fprintln(os.Stderr, "Hint: run 'elements auth login' to re-authenticate")
		}
		os.Exit(1)
	}
}

func isFlagError(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "required flag") ||
		strings.Contains(msg, "unknown flag") ||
		strings.Contains(msg, "flag needs an argument")
}
