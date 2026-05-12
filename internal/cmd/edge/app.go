package edge

import (
	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
)

func newCmdApp(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app",
		Short: "Manage edge apps",
	}

	cmd.AddCommand(newCmdAppList(f))
	cmd.AddCommand(newCmdAppGet(f))
	cmd.AddCommand(newCmdAppCreate(f))
	cmd.AddCommand(newCmdAppUpdate(f))
	cmd.AddCommand(newCmdAppDelete(f))
	cmd.AddCommand(newCmdAppLogs(f))
	cmd.AddCommand(newCmdAppDeploy(f))
	cmd.AddCommand(newCmdAppUndeploy(f))

	return cmd
}
