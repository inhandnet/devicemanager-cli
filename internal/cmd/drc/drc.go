package drc

import (
	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
)

func NewCmdDRC(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "drc",
		Short: "Manage DRC configuration templates",
	}

	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdGet(f))
	cmd.AddCommand(NewCmdCreate(f))
	cmd.AddCommand(NewCmdDelete(f))
	cmd.AddCommand(NewCmdDevices(f))

	return cmd
}
