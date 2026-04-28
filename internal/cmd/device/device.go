package device

import (
	"github.com/spf13/cobra"

	"github.com/inhandnet/elements-cli/internal/factory"
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

	return cmd
}
