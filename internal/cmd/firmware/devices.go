package firmware

import (
	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
)

func NewCmdDevices(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "devices <firmware-id>",
		Short: "Manage devices in a firmware upgrade job",
	}

	cmd.AddCommand(NewCmdDevicesList(f))
	cmd.AddCommand(NewCmdDevicesAdd(f))
	cmd.AddCommand(NewCmdDevicesRemove(f))

	return cmd
}
