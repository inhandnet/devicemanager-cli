package edge

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
	"github.com/inhandnet/devicemanager-cli/internal/ui"
)

func newCmdVersionDelete(f *factory.Factory) *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <app-id> <version>",
		Short: "Delete an edge app version",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			appID := args[0]
			version := args[1]

			if !ui.Confirm(f.IO, fmt.Sprintf("Delete version %s of app %s?", version, appID), yes) {
				return nil
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Delete(fmt.Sprintf("/api/edge/apps/%s/versions/%s", appID, version))
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "Version %s deleted\n", version)
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip confirmation prompt")

	return cmd
}
