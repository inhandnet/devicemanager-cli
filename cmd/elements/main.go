package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/inhandnet/elements-cli/internal/api"
	cmd "github.com/inhandnet/elements-cli/internal/cmd"
	apiCmd "github.com/inhandnet/elements-cli/internal/cmd/api"
	authCmd "github.com/inhandnet/elements-cli/internal/cmd/auth"
	configCmd "github.com/inhandnet/elements-cli/internal/cmd/config"
	deviceCmd "github.com/inhandnet/elements-cli/internal/cmd/device"
	devicegroupCmd "github.com/inhandnet/elements-cli/internal/cmd/devicegroup"
	drcCmd "github.com/inhandnet/elements-cli/internal/cmd/drc"
	edgeCmd "github.com/inhandnet/elements-cli/internal/cmd/edge"
	firmwareCmd "github.com/inhandnet/elements-cli/internal/cmd/firmware"
	tunnelCmd "github.com/inhandnet/elements-cli/internal/cmd/tunnel"
	versionCmd "github.com/inhandnet/elements-cli/internal/cmd/version"
	"github.com/inhandnet/elements-cli/internal/factory"
)

func main() {
	f := factory.New()
	rootCmd := cmd.NewCmdRoot(f)
	rootCmd.AddCommand(authCmd.NewCmdAuth(f))
	rootCmd.AddCommand(configCmd.NewCmdConfig(f))
	rootCmd.AddCommand(apiCmd.NewCmdApi(f))
	rootCmd.AddCommand(deviceCmd.NewCmdDevice(f))
	rootCmd.AddCommand(devicegroupCmd.NewCmdDeviceGroup(f))
	rootCmd.AddCommand(tunnelCmd.NewCmdTunnel(f))
	rootCmd.AddCommand(drcCmd.NewCmdDRC(f))
	rootCmd.AddCommand(edgeCmd.NewCmdEdge(f))
	rootCmd.AddCommand(firmwareCmd.NewCmdFirmware(f))
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
