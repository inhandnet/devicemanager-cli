package edge

import (
	"github.com/spf13/cobra"

	"github.com/inhandnet/elements-cli/internal/factory"
)

func newCmdConfig(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage edge app configs",
	}

	cmd.AddCommand(newCmdConfigList(f))
	cmd.AddCommand(newCmdConfigGet(f))
	cmd.AddCommand(newCmdConfigCreate(f))
	cmd.AddCommand(newCmdConfigUpdate(f))
	cmd.AddCommand(newCmdConfigDelete(f))
	cmd.AddCommand(newCmdConfigDeploy(f))

	return cmd
}
