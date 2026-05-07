package device

import (
	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
)

func NewCmdDevice(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "device",
		Short: "Manage devices",
	}

	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdGet(f))
	cmd.AddCommand(NewCmdCreate(f))
	cmd.AddCommand(NewCmdSignal(f))
	cmd.AddCommand(NewCmdKick(f))
	cmd.AddCommand(NewCmdReboot(f))
	cmd.AddCommand(NewCmdTraffic(f))
	cmd.AddCommand(NewCmdClients(f))
	cmd.AddCommand(NewCmdAlert(f))
	cmd.AddCommand(NewCmdConfig(f))
	cmd.AddCommand(NewCmdUpdate(f))
	cmd.AddCommand(NewCmdDelete(f))
	cmd.AddCommand(NewCmdOnlineStats(f))
	cmd.AddCommand(NewCmdAlertRule(f))
	cmd.AddCommand(NewCmdAlertConfirm(f))
	cmd.AddCommand(NewCmdOnlineEvents(f))
	cmd.AddCommand(NewCmdRegisterEvents(f))
	cmd.AddCommand(NewCmdModels(f))
	cmd.AddCommand(NewCmdStats(f))
	cmd.AddCommand(NewCmdCount(f))

	return cmd
}
