package device

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
)

func resolveNgrokServer(host string, globalConfig string) string {
	// Priority 1: user-defined global config
	if globalConfig != "" {
		return globalConfig
	}
	// Priority 2: auto-detect by API host
	switch host {
	case "iot.inhandnetworks.com":
		return "ngrok.iot.inhandnetworks.com:4443"
	case "iot.inhand.com.cn":
		return "iot.inhand.com.cn:4443"
	default:
		return "10.5.17.52:4443"
	}
}

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

			ctx, err := f.Config()
			if err != nil {
				return err
			}

			deviceID := args[0]
			proto, _ := cmd.Flags().GetString("proto")
			port, _ := cmd.Flags().GetInt("port")
			server, _ := cmd.Flags().GetString("server")

			// If --server not explicitly set, resolve from config or auto-detect
			if !cmd.Flags().Changed("server") {
				active, _ := ctx.ActiveContext()
				host := ""
				if active != nil {
					host = active.Host
				}
				server = resolveNgrokServer(host, ctx.NgrokServer)
			}

			body := map[string]any{
				"objectId":   deviceID,
				"objectName": deviceID,
				"name":       "ngrok connect",
				"type":       "23",
				"priority":   30,
				"timeout":    20000,
				"data": map[string]any{
					"server": server,
					"proto":  proto,
					"port":   port,
				},
			}

			fmt.Fprintf(f.IO.Out, "Starting remote web management for device %s (server: %s)...\n", deviceID, server)
			resp, err := client.Post("/api2/tasks/run", body)
			if err != nil {
				return err
			}

			apiErr := gjson.GetBytes(resp, "error").String()
			if apiErr != "" {
				return fmt.Errorf("API error: %s (code: %d)", apiErr, gjson.GetBytes(resp, "error_code").Int())
			}

			state := gjson.GetBytes(resp, "result.state").Int()
			if state == 3 {
				url := gjson.GetBytes(resp, "result.data.response").String()
				if url != "" {
					fmt.Fprintf(f.IO.Out, "Remote web URL: %s\n", url)
					return nil
				}
				return fmt.Errorf("task completed but no URL returned")
			}

			errMsg := gjson.GetBytes(resp, "result.error").String()
			if errMsg != "" {
				return fmt.Errorf("task failed: %s", errMsg)
			}
			return fmt.Errorf("unexpected task state: %d", state)
		},
	}

	cmd.Flags().String("proto", "http", "Protocol (http/https)")
	cmd.Flags().Int("port", 80, "Device web management port")
	cmd.Flags().String("server", "", "Ngrok server address (overrides auto-detect)")

	return cmd
}
