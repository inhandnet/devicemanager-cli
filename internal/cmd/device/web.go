package device

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdWeb(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "web <device-id>",
		Short: "Start remote web management for a device",
		Example: `  # Start remote web management
  devicemanager device web <device-id>

  # Use custom port
  devicemanager device web <device-id> --port 443 --proto https`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			deviceID := args[0]
			proto, _ := cmd.Flags().GetString("proto")
			port, _ := cmd.Flags().GetInt("port")
			server, _ := cmd.Flags().GetString("server")

			body := map[string]any{
				"objectId": deviceID,
				"name":     "ngrok connect",
				"type":     "23",
				"priority": 30,
				"timeout":  20000,
				"data": map[string]any{
					"server": server,
					"proto":  proto,
					"port":   port,
				},
			}

			output, _ := cmd.Flags().GetString("output")

			fmt.Fprintf(f.IO.Out, "Starting remote web management for device %s...\n", deviceID)
			resp, err := client.Post("/api2/tasks/run", body)
			if err != nil {
				return err
			}

			state := gjson.GetBytes(resp, "result.state").Int()
			if state == 3 {
				url := gjson.GetBytes(resp, "result.data.response").String()
				if url != "" {
					fmt.Fprintf(f.IO.Out, "Remote web URL: %s\n", url)
					return nil
				}
			}

			errMsg := gjson.GetBytes(resp, "result.error").String()
			if errMsg != "" {
				return fmt.Errorf("task failed: %s", errMsg)
			}

			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().String("proto", "http", "Protocol (http/https)")
	cmd.Flags().Int("port", 80, "Device web management port")
	cmd.Flags().String("server", "ngrok.j3r0lin.com:4443", "Ngrok server address")

	return cmd
}
