package firmware

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdDevicesAdd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <firmware-id> <device-id>...",
		Short: "Batch upgrade devices with a firmware",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			firmwareID := args[0]
			deviceIDs := args[1:]

			body := map[string]interface{}{
				"deviceIds": deviceIDs,
			}

			groupIDs, _ := cmd.Flags().GetStringSlice("group")
			if len(groupIDs) > 0 {
				body["deviceGroupIds"] = groupIDs
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Post(fmt.Sprintf("/api/firmwares/%s/devices", firmwareID), body)
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "Added %d device(s) to firmware upgrade %s\n", len(deviceIDs), firmwareID)
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().StringSlice("group", nil, "Device group IDs to upgrade")

	return cmd
}
