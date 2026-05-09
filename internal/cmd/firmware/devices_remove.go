package firmware

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
	"github.com/inhandnet/devicemanager-cli/internal/ui"
)

func NewCmdDevicesRemove(f *factory.Factory) *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "remove <firmware-id> <device-id>",
		Short: "Cancel a device firmware upgrade task",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			firmwareID := args[0]
			deviceID := args[1]

			if !ui.Confirm(f.IO, fmt.Sprintf("Cancel upgrade for device %s?", deviceID), yes) {
				return nil
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Delete(fmt.Sprintf("/api/job/%s/devices/%s", firmwareID, deviceID))
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "Cancelled upgrade for device %s\n", deviceID)
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip confirmation prompt")

	return cmd
}
