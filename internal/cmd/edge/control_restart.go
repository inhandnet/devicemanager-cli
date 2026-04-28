package edge

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/elements-cli/internal/factory"
	"github.com/inhandnet/elements-cli/internal/iostreams"
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

			resp, err := client.Get(fmt.Sprintf("/api/edge/devices/%s/app/%s/restart", deviceID, appID), nil)
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "App %s restarted on device %s\n", appID, deviceID)
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	return cmd
}
