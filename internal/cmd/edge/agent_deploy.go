package edge

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func newCmdAgentDeploy(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy <agent-id> [device-id]...",
		Short: "Deploy an edge agent to devices or device groups",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			agentID := args[0]
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

			resp, err := client.Post(fmt.Sprintf("/api/edge/agents/%s/deploy", agentID), body)
			if err != nil {
				return err
			}

			switch {
			case len(deviceIDs) > 0 && len(groupIDs) > 0:
				fmt.Fprintf(f.IO.Out, "Deployed agent %s to %d device(s) and %d group(s)\n", agentID, len(deviceIDs), len(groupIDs))
			case len(groupIDs) > 0:
				fmt.Fprintf(f.IO.Out, "Deployed agent %s to %d group(s)\n", agentID, len(groupIDs))
			default:
				fmt.Fprintf(f.IO.Out, "Deployed agent %s to %d device(s)\n", agentID, len(deviceIDs))
			}
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().StringSlice("group", nil, "Device group IDs to deploy to")
	return cmd
}
