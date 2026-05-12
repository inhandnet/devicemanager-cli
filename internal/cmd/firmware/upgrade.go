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
			timeout, _ := cmd.Flags().GetInt("timeout")

			body := map[string]any{
				"firmwareId": firmwareID,
				"timeout":    timeout,
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Post(fmt.Sprintf("/api/device/%s/upgrade", deviceID), body)
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "Upgrade started for device %s\n", deviceID)
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().String("firmware-id", "", "Firmware ID (required)")
	_ = cmd.MarkFlagRequired("firmware-id")
	cmd.Flags().Int("timeout", 600, "Upgrade timeout in seconds (required)")

	return cmd
}
