package devicegroup

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdDevicesRemove(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <group-id> <device-id>...",
		Short: "Remove devices from a group",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			groupID := args[0]
			deviceIDs := args[1:]

			body := map[string]interface{}{
				"deviceIds": deviceIDs,
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Post(fmt.Sprintf("/api/devicegroups/%s/devices/remove", groupID), body)
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "Removed %d device(s) from group %s\n", len(deviceIDs), groupID)
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	return cmd
}
