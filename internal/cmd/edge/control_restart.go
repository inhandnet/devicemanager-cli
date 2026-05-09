package edge

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func newCmdControlRestart(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restart <device-id> <app-id>",
		Short: "Restart an edge app on a device",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			deviceID := args[0]
			appID := args[1]

			output, _ := cmd.Flags().GetString("output")

			body := map[string]any{"apps": []string{appID}}
			resp, err := client.Post(fmt.Sprintf("/api/edge/devices/%s/apps/restart", deviceID), body)
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "App %s restarted on device %s\n", appID, deviceID)
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	return cmd
}
