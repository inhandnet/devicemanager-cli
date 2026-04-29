package drc

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
	"github.com/inhandnet/devicemanager-cli/internal/ui"
)

func NewCmdDevicesRemove(f *factory.Factory) *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "remove <template-id> <device-id>",
		Short: "Remove a device from a DRC template",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			templateID := args[0]
			deviceID := args[1]

			if !ui.Confirm(f.IO, fmt.Sprintf("Remove device %s from template %s?", deviceID, templateID), yes) {
				return nil
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Delete(fmt.Sprintf("/api/job/%s/devices/%s", templateID, deviceID))
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "Removed device %s from template %s\n", deviceID, templateID)
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip confirmation prompt")

	return cmd
}
