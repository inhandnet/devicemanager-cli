package drc

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/elements-cli/internal/factory"
	"github.com/inhandnet/elements-cli/internal/iostreams"
)

func NewCmdDevicesRemove(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <template-id> <device-id>",
		Short: "Remove a device from a DRC template",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			templateID := args[0]
			deviceID := args[1]

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Delete(fmt.Sprintf("/api/job/%s/devices/%s", templateID, deviceID))
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "Removed device %s from template %s\n", deviceID, templateID)
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	return cmd
}
