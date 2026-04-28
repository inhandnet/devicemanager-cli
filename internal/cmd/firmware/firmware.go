package firmware

import (
	"github.com/spf13/cobra"

	"github.com/inhandnet/elements-cli/internal/factory"
)

func NewCmdFirmware(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "firmware",
		Short: "Manage firmware and upgrades",
	}

	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdCreate(f))
	cmd.AddCommand(NewCmdUpload(f))
	cmd.AddCommand(NewCmdUpgrade(f))
	cmd.AddCommand(NewCmdDevices(f))

	return cmd
}
