package edge

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
	"github.com/inhandnet/devicemanager-cli/internal/ui"
)

func newCmdConfigDelete(f *factory.Factory) *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <app-id> <config-id>",
		Short: "Delete an edge app config",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			appID := args[0]
			configID := args[1]

			if !ui.Confirm(f.IO, fmt.Sprintf("Delete config %s of app %s?", configID, appID), yes) {
				return nil
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Delete(fmt.Sprintf("/api/edge/apps/%s/configs/%s", appID, configID))
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "Config %s deleted\n", configID)
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip confirmation prompt")

	return cmd
}
