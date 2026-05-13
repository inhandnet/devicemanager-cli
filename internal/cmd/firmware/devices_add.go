package firmware

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/api"
	"github.com/inhandnet/devicemanager-cli/internal/cmdutil"
	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdDevicesAdd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <firmware-id> [device-id]...",
		Short: "Batch upgrade devices with a firmware",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			firmwareID := args[0]
			deviceIDs := args[1:]
			groupIDs, _ := cmd.Flags().GetStringSlice("group")

			if len(deviceIDs) == 0 && len(groupIDs) == 0 {
				return fmt.Errorf("must specify at least one device ID or --group")
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{
				"deviceIds":      deviceIDs,
				"deviceGroupIds": groupIDs,
			}

			q := url.Values{}
			cmdutil.SetQueryParam(q, "oid", f.OrgID())

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Do("POST", fmt.Sprintf("/api/firmware/%s/devices", firmwareID), &api.RequestOptions{
				Query: q,
				Body:  body,
			})
			if err != nil {
				return err
			}

			switch {
			case len(deviceIDs) > 0 && len(groupIDs) > 0:
				fmt.Fprintf(f.IO.Out, "Added %d device(s) and %d group(s) to firmware upgrade %s\n", len(deviceIDs), len(groupIDs), firmwareID)
			case len(groupIDs) > 0:
				fmt.Fprintf(f.IO.Out, "Added %d group(s) to firmware upgrade %s\n", len(groupIDs), firmwareID)
			default:
				fmt.Fprintf(f.IO.Out, "Added %d device(s) to firmware upgrade %s\n", len(deviceIDs), firmwareID)
			}
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().StringSlice("group", nil, "Device group IDs to upgrade")

	return cmd
}
