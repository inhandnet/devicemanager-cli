package devicegroup

import (
	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
)

func NewCmdDevices(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "devices <group-id>",
		Short: "Manage devices in a group",
	}

	cmd.AddCommand(NewCmdDevicesList(f))
	cmd.AddCommand(NewCmdDevicesAdd(f))
	cmd.AddCommand(NewCmdDevicesRemove(f))
	cmd.AddCommand(NewCmdDevicesAvailable(f))

	return cmd
}
