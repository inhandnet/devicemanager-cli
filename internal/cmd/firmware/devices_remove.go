package firmware

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdDevicesRemove(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <firmware-id> <device-id>",
		Short: "Cancel a device firmware upgrade task",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			firmwareID := args[0]
			deviceID := args[1]

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Delete(fmt.Sprintf("/api/jobs/%s/devices/%s", firmwareID, deviceID))
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "Cancelled upgrade for device %s\n", deviceID)
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	return cmd
}
