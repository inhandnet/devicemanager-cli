package edge

import (
	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
)

func NewCmdEdge(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edge",
		Short: "Manage edge computing (agents, apps, versions, configs, control)",
	}

	cmd.AddCommand(newCmdAgent(f))
	cmd.AddCommand(newCmdApp(f))
	cmd.AddCommand(newCmdVersion(f))
	cmd.AddCommand(newCmdConfig(f))
	cmd.AddCommand(newCmdControl(f))
	cmd.AddCommand(newCmdDevice(f))

	return cmd
}
