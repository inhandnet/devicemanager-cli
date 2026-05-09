package device

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdConfig(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage device configuration",
	}

	cmd.AddCommand(NewCmdConfigGet(f))
	cmd.AddCommand(NewCmdConfigSet(f))
	cmd.AddCommand(NewCmdConfigExport(f))

	return cmd
}

func NewCmdConfigGet(f *factory.Factory) *cobra.Command {
	var skipRefresh bool

	cmd := &cobra.Command{
		Use:   "get <device-id>",
		Short: "Get device running configuration",
		Long: `Fetch the device's running configuration. By default, it first sends a
"GET RUNNING CONFIG" task to the device to retrieve the latest config.
If the device is offline or the task fails, the last known config is returned.
Use --skip-refresh to skip the task and return the cached config directly.`,
		Args: cobra.ExactArgs(1),
		Example: `  # Get latest config (sends task to device first)
  devicemanager device config get <device-id>

  # Get cached config without refreshing
  devicemanager device config get <device-id> --skip-refresh`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			deviceID := args[0]

			if !skipRefresh {
				// Get device name for the task
				objectName := deviceID
				if devBody, err := client.Get(fmt.Sprintf("/api/devices/%s", deviceID), nil); err == nil {
					if name := gjson.GetBytes(devBody, "result.name").String(); name != "" {
						objectName = name
					} else if name := gjson.GetBytes(devBody, "name").String(); name != "" {
						objectName = name
					}
				}

				// Send GET RUNNING CONFIG task and wait for completion
				taskBody := map[string]any{
					"objectId":   deviceID,
					"objectName": objectName,
					"name":       "GET RUNNING CONFIG",
					"type":       "4",
					"priority":   30,
					"timeout":    30000,
				}
				fmt.Fprintf(f.IO.ErrOut, "Fetching latest config from device...\n")
				_, _ = client.Post("/api2/tasks/run", taskBody)
			}

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get(fmt.Sprintf("/api/devices/%s/config", deviceID), nil)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output)
		},
	}

	cmd.Flags().BoolVar(&skipRefresh, "skip-refresh", false, "Skip refreshing config from device, return cached config")

	return cmd
}

func NewCmdConfigSet(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set <device-id>",
		Short: "Set device configuration",
		Args:  cobra.ExactArgs(1),
		Example: `  devicemanager device config set <device-id> --content "..."
  devicemanager device config set <device-id> --content-file config.json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			deviceID := args[0]

			content, _ := cmd.Flags().GetString("content")
			contentFile, _ := cmd.Flags().GetString("content-file")

			if content == "" && contentFile == "" {
				return fmt.Errorf("either --content or --content-file is required")
			}

			if contentFile != "" {
				data, err := os.ReadFile(contentFile)
				if err != nil {
					return fmt.Errorf("reading content file: %w", err)
				}
				content = string(data)
			}

			desc, _ := cmd.Flags().GetString("description")

			body := map[string]interface{}{
				"deviceType":    0,
				"deviceContent": content,
			}
			if desc != "" {
				body["deviceDesc"] = desc
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Post(fmt.Sprintf("/api/devices/%s/config/set", deviceID), body)
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "Configuration sent to device %s\n", deviceID)
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().String("content", "", "Configuration content")
	cmd.Flags().String("content-file", "", "Configuration file path")
	cmd.Flags().String("description", "", "Configuration description")

	return cmd
}

func NewCmdConfigExport(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "export <device-id>",
		Short: "Export device configuration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get(fmt.Sprintf("/api/devices/%s/config/export", args[0]), nil)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output)
		},
	}
}
