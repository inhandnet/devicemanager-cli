package tunnel

import (
	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
)

func NewCmdTunnel(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tunnel",
		Short: "Manage remote tunnels",
	}

	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdCreate(f))
	cmd.AddCommand(NewCmdUpdate(f))
	cmd.AddCommand(NewCmdDelete(f))
	cmd.AddCommand(NewCmdConnect(f))
	cmd.AddCommand(NewCmdDisconnect(f))

	return cmd
}
