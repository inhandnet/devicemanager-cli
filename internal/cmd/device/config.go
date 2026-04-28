package device

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/elements-cli/internal/factory"
	"github.com/inhandnet/elements-cli/internal/iostreams"
)

func NewCmdConfig(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage device configuration",
	}

	cmd.AddCommand(NewCmdConfigGet(f))
	cmd.AddCommand(NewCmdConfigSet(f))

	return cmd
}

func NewCmdConfigGet(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get <device-id>",
		Short:   "Get device running configuration",
		Args:    cobra.ExactArgs(1),
		Example: `  elements device config get <device-id>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get(fmt.Sprintf("/api/devices/%s/config", args[0]), nil)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output)
		},
	}

	return cmd
}

func NewCmdConfigSet(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set <device-id>",
		Short: "Set device configuration",
		Args:  cobra.ExactArgs(1),
		Example: `  elements device config set <device-id> --content "..."
  elements device config set <device-id> --content-file config.json`,
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
				// TODO: read file content
				return fmt.Errorf("--content-file not yet implemented")
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

			resp, err := client.Post(fmt.Sprintf("/api/devices/%s/config/set2", deviceID), body)
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
