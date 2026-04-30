package edge

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func newCmdAppLogs(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "logs <device-id> <app-name>",
		Short: "View edge app runtime logs on a device",
		Args:  cobra.ExactArgs(2),
		Example: `  # View logs for app "data-collector" on a device
  devicemanager edge app logs 639acc6e078a3d0001be2963 data-collector`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			deviceID := args[0]
			appName := args[1]

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get(fmt.Sprintf("/api/edge/devices/%s/apps/%s/logs", deviceID, appName), nil)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output,
				iostreams.WithColumns("ts", "level", "content"))
		},
	}
}
