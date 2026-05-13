package device

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
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

			fmt.Fprintf(f.IO.Out, "Starting remote web management for device %s...\n", deviceID)
			resp, err := client.Post("/api2/tasks/run", body)
			if err != nil {
				return err
			}

			taskID := gjson.GetBytes(resp, "result._id").String()
			if taskID == "" {
				return fmt.Errorf("failed to get task ID from response")
			}

			const (
				pollInterval = 2 * time.Second
				pollTimeout  = 30 * time.Second
			)
			deadline := time.Now().Add(pollTimeout)

			for {
				state := gjson.GetBytes(resp, "result.state").Int()

				switch {
				case state == 3: // completed
					url := gjson.GetBytes(resp, "result.data.response").String()
					if url != "" {
						fmt.Fprintf(f.IO.Out, "Remote web URL: %s\n", url)
						return nil
					}
					return fmt.Errorf("task completed but no URL returned")

				case state == -1 || state == 2: // failed
					errMsg := gjson.GetBytes(resp, "result.error").String()
					if errMsg != "" {
						return fmt.Errorf("task failed: %s", errMsg)
					}
					return fmt.Errorf("task failed (state=%d)", state)
				}

				if time.Now().After(deadline) {
					return fmt.Errorf("timeout waiting for task to complete (task ID: %s)", taskID)
				}

				time.Sleep(pollInterval)

				resp, err = client.Get(fmt.Sprintf("/api2/tasks/%s", taskID), nil)
				if err != nil {
					return fmt.Errorf("polling task status: %w", err)
				}
			}
		},
	}

	cmd.Flags().String("proto", "http", "Protocol (http/https)")
	cmd.Flags().Int("port", 80, "Device web management port")
	cmd.Flags().String("server", "ngrok.j3r0lin.com:4443", "Ngrok server address")

	return cmd
}
