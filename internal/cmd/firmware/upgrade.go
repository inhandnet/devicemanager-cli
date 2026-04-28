package firmware

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdUpgrade(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "upgrade <device-id>",
		Short:   "Upgrade a single device",
		Args:    cobra.ExactArgs(1),
		Example: `  devicemanager firmware upgrade <device-id> --firmware-id <id> --timeout 600`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			deviceID := args[0]
			firmwareID, _ := cmd.Flags().GetString("firmware-id")
			deviceName, _ := cmd.Flags().GetString("device-name")
			timeout, _ := cmd.Flags().GetInt("timeout")

			if firmwareID == "" {
				return fmt.Errorf("--firmware-id is required")
			}

			body := map[string]interface{}{
				"firmwareId": firmwareID,
			}
			if deviceName != "" {
				body["deviceName"] = deviceName
			}
			if timeout > 0 {
				body["timeout"] = timeout
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Post(fmt.Sprintf("/api/devices/%s/upgrade", deviceID), body)
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "Upgrade started for device %s\n", deviceID)
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().String("firmware-id", "", "Firmware ID (required)")
	cmd.Flags().String("device-name", "", "Device name")
	cmd.Flags().Int("timeout", 0, "Upgrade timeout (seconds)")

	return cmd
}
