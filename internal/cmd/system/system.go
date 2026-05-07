package system

import (
	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
)

func NewCmdSystem(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "system",
		Short:   "System management (users, permissions, org, logs)",
		Aliases: []string{"sys"},
	}

	cmd.AddCommand(NewCmdUser(f))
	cmd.AddCommand(NewCmdRole(f))
	cmd.AddCommand(NewCmdPermission(f))
	cmd.AddCommand(NewCmdOrg(f))
	cmd.AddCommand(NewCmdLog(f))

	return cmd
}
