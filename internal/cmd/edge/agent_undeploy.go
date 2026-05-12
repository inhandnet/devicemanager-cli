package edge

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func newCmdAgentUndeploy(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "undeploy <agent-id> <device-id>...",
		Short: "Remove an edge agent deployment from devices",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			agentID := args[0]
			deviceIDs := args[1:]

			body := map[string]any{
				"deviceIds": deviceIDs,
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Post(fmt.Sprintf("/api/edge/agents/%s/delete-deploy", agentID), body)
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "Removed agent %s from %d device(s)\n", agentID, len(deviceIDs))
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	return cmd
}
