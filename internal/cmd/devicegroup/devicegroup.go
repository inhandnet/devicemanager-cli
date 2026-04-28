package devicegroup

import (
	"github.com/spf13/cobra"

	"github.com/inhandnet/elements-cli/internal/factory"
)

func NewCmdDeviceGroup(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "devicegroup",
		Short:   "Manage device groups",
		Aliases: []string{"dg"},
	}

	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdGet(f))
	cmd.AddCommand(NewCmdCreate(f))
	cmd.AddCommand(NewCmdUpdate(f))
	cmd.AddCommand(NewCmdDelete(f))
	cmd.AddCommand(NewCmdDevices(f))

	return cmd
}
