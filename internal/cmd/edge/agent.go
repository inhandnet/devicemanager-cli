package edge

import (
	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
)

func newCmdAgent(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Manage edge agents (engines)",
	}

	cmd.AddCommand(newCmdAgentList(f))
	cmd.AddCommand(newCmdAgentGet(f))
	cmd.AddCommand(newCmdAgentUpload(f))
	cmd.AddCommand(newCmdAgentUpdate(f))
	cmd.AddCommand(newCmdAgentDelete(f))
	cmd.AddCommand(newCmdAgentDevices(f))
	cmd.AddCommand(newCmdAgentDeploy(f))
	cmd.AddCommand(newCmdAgentUndeploy(f))

	return cmd
}
