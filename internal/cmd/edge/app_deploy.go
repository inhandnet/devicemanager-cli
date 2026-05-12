package edge

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func newCmdAppDeploy(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy <app-id> --version <version> [device-id]...",
		Short: "Deploy an edge app version to devices or device groups",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			appID := args[0]
			deviceIDs := args[1:]
			groupIDs, _ := cmd.Flags().GetStringSlice("group")
			version, _ := cmd.Flags().GetString("version")

			if version == "" {
				return fmt.Errorf("--version is required")
			}
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

			resp, err := client.Post(fmt.Sprintf("/api/edge/apps/%s/versions/%s/deploy", appID, version), body)
			if err != nil {
				return err
			}

			switch {
			case len(deviceIDs) > 0 && len(groupIDs) > 0:
				fmt.Fprintf(f.IO.Out, "Deployed app %s (version %s) to %d device(s) and %d group(s)\n", appID, version, len(deviceIDs), len(groupIDs))
			case len(groupIDs) > 0:
				fmt.Fprintf(f.IO.Out, "Deployed app %s (version %s) to %d group(s)\n", appID, version, len(groupIDs))
			default:
				fmt.Fprintf(f.IO.Out, "Deployed app %s (version %s) to %d device(s)\n", appID, version, len(deviceIDs))
			}
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().String("version", "", "App version to deploy (required)")
	cmd.Flags().StringSlice("group", nil, "Device group IDs to deploy to")
	_ = cmd.MarkFlagRequired("version")
	return cmd
}
