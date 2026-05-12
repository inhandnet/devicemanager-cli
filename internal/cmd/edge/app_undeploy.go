package edge

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func newCmdAppUndeploy(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "undeploy <app-id> <device-id>...",
		Short: "Cancel an edge app deployment from devices",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			appID := args[0]
			deviceIDs := args[1:]

			body := map[string]any{
				"deviceIds": deviceIDs,
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Post(fmt.Sprintf("/api/edge/apps/%s/delete-deploy", appID), body)
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "Cancelled deployment of app %s from %d device(s)\n", appID, len(deviceIDs))
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	return cmd
}
