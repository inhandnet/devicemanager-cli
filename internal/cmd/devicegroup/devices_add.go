package devicegroup

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/elements-cli/internal/factory"
	"github.com/inhandnet/elements-cli/internal/iostreams"
)

func NewCmdDevicesAdd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <group-id> <device-id>...",
		Short: "Add devices to a group",
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

			resp, err := client.Post(fmt.Sprintf("/api/devicegroups/%s/devices", groupID), body)
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "Added %d device(s) to group %s\n", len(deviceIDs), groupID)
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	return cmd
}
