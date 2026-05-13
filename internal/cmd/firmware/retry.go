package firmware

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdRetry(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "retry <firmware-id> <device-id>",
		Short: "Retry a firmware upgrade task for a device",
		Long:  `Retry a failed firmware upgrade for a device. The firmware-id can be found via 'firmware list'.`,
		Example: `  # Retry upgrade for a device
  devicemanager firmware retry <firmware-id> <device-id>

  # Check which devices failed
  devicemanager firmware devices list <firmware-id> --status failed`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			firmwareID := args[0]
			deviceID := args[1]

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Put(fmt.Sprintf("/api/jobs/%s/devices/%s/restart", firmwareID, deviceID), nil)
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "Retried upgrade for device %s in firmware %s\n", deviceID, firmwareID)
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	return cmd
}
