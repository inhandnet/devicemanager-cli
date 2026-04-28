package edge

import (
	"github.com/spf13/cobra"

	"github.com/inhandnet/elements-cli/internal/factory"
)

func newCmdVersion(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Manage edge app versions",
	}

	cmd.AddCommand(newCmdVersionList(f))
	cmd.AddCommand(newCmdVersionUpload(f))
	cmd.AddCommand(newCmdVersionUpdate(f))
	cmd.AddCommand(newCmdVersionDelete(f))
	cmd.AddCommand(newCmdVersionDeploy(f))

	return cmd
}
