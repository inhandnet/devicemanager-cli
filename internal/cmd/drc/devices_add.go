package drc

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdDevicesAdd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <template-id> [device-id]...",
		Short: "Assign devices or device groups to a DRC template",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			templateID := args[0]
			deviceIDs := args[1:]
			groupIDs, _ := cmd.Flags().GetStringSlice("group")

			if len(deviceIDs) == 0 && len(groupIDs) == 0 {
				return fmt.Errorf("must specify at least one device ID or --group")
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{}
			if len(deviceIDs) > 0 {
				body["deviceIds"] = deviceIDs
			}
			if len(groupIDs) > 0 {
				body["deviceGroupIds"] = groupIDs
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Post(fmt.Sprintf("/api/drc/%s/devices", templateID), body)
			if err != nil {
				return err
			}

			switch {
			case len(deviceIDs) > 0 && len(groupIDs) > 0:
				fmt.Fprintf(f.IO.Out, "Assigned %d device(s) and %d group(s) to template %s\n", len(deviceIDs), len(groupIDs), templateID)
			case len(groupIDs) > 0:
				fmt.Fprintf(f.IO.Out, "Assigned %d group(s) to template %s\n", len(groupIDs), templateID)
			default:
				fmt.Fprintf(f.IO.Out, "Assigned %d device(s) to template %s\n", len(deviceIDs), templateID)
			}
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().StringSlice("group", nil, "Device group IDs to assign")

	return cmd
}
