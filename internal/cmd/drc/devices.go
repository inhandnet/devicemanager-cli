package drc

import (
	"github.com/spf13/cobra"

	"github.com/inhandnet/elements-cli/internal/factory"
)

func NewCmdDevices(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "devices <template-id>",
		Short: "Manage devices assigned to a DRC template",
	}

	cmd.AddCommand(NewCmdDevicesList(f))
	cmd.AddCommand(NewCmdDevicesAdd(f))
	cmd.AddCommand(NewCmdDevicesRemove(f))
	cmd.AddCommand(NewCmdDevicesRestart(f))

	return cmd
}
