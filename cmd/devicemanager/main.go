package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/inhandnet/devicemanager-cli/internal/api"
	cmd "github.com/inhandnet/devicemanager-cli/internal/cmd"
	apiCmd "github.com/inhandnet/devicemanager-cli/internal/cmd/api"
	authCmd "github.com/inhandnet/devicemanager-cli/internal/cmd/auth"
	configCmd "github.com/inhandnet/devicemanager-cli/internal/cmd/config"
	deviceCmd "github.com/inhandnet/devicemanager-cli/internal/cmd/device"
	devicegroupCmd "github.com/inhandnet/devicemanager-cli/internal/cmd/devicegroup"
	docsCmd "github.com/inhandnet/devicemanager-cli/internal/cmd/docs"
	drcCmd "github.com/inhandnet/devicemanager-cli/internal/cmd/drc"
	edgeCmd "github.com/inhandnet/devicemanager-cli/internal/cmd/edge"
	firmwareCmd "github.com/inhandnet/devicemanager-cli/internal/cmd/firmware"
	systemCmd "github.com/inhandnet/devicemanager-cli/internal/cmd/system"
	taskCmd "github.com/inhandnet/devicemanager-cli/internal/cmd/task"
	tunnelCmd "github.com/inhandnet/devicemanager-cli/internal/cmd/tunnel"
	updateCmd "github.com/inhandnet/devicemanager-cli/internal/cmd/update"
	versionCmd "github.com/inhandnet/devicemanager-cli/internal/cmd/version"
	"github.com/inhandnet/devicemanager-cli/internal/factory"
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
	rootCmd.AddCommand(taskCmd.NewCmdTask(f))
	rootCmd.AddCommand(systemCmd.NewCmdSystem(f))
	rootCmd.AddCommand(docsCmd.NewCmdDocs(f))
	rootCmd.AddCommand(updateCmd.NewCmdUpdate(f))
	rootCmd.AddCommand(versionCmd.NewCmdVersion(f))

	// Top-level shortcut: `devicemanager login` → `devicemanager auth login`
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
			fmt.Fprintln(os.Stderr, "Hint: run 'devicemanager auth login' to re-authenticate")
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
