package device

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdUpdate(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <device-id>",
		Short: "Update a device",
		Args:  cobra.ExactArgs(1),
		Example: `  # Rename a device
  devicemanager device update 5d6349d6335c8c000178a194 --name "new-name"

  # Update description
  devicemanager device update 5d6349d6335c8c000178a194 --description "office router"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]interface{}{}
			if cmd.Flags().Changed("name") {
				v, _ := cmd.Flags().GetString("name")
				body["name"] = v
			}
			if cmd.Flags().Changed("description") {
				v, _ := cmd.Flags().GetString("description")
				body["description"] = v
			}

			if len(body) == 0 {
				return fmt.Errorf("at least one flag (--name, --description) is required")
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Put(fmt.Sprintf("/api/devices/%s?verbose=100", args[0]), body)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().String("name", "", "New device name")
	cmd.Flags().String("description", "", "New device description")

	return cmd
}
