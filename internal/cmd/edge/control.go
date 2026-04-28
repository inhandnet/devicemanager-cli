package edge

import (
	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
)

func newCmdControl(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "control",
		Short: "Remote control edge apps on devices",
	}

	cmd.AddCommand(newCmdControlStart(f))
	cmd.AddCommand(newCmdControlStop(f))
	cmd.AddCommand(newCmdControlRestart(f))

	return cmd
}
