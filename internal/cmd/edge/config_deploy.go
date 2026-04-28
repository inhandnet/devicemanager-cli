package edge

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func newCmdConfigDeploy(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy <app-id> <version>",
		Short: "Deploy an edge app config to devices",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			appID := args[0]
			version := args[1]

			body := map[string]interface{}{}
			deviceIDs, _ := cmd.Flags().GetStringSlice("device")
			groupIDs, _ := cmd.Flags().GetStringSlice("group")
			if len(deviceIDs) > 0 {
				body["deviceIds"] = deviceIDs
			}
			if len(groupIDs) > 0 {
				body["deviceGroupIds"] = groupIDs
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Post(fmt.Sprintf("/api/edge/apps/%s/configs/%s/deploy", appID, version), body)
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "Config deployed\n")
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().StringSlice("device", nil, "Device IDs to deploy")
	cmd.Flags().StringSlice("group", nil, "Device group IDs to deploy")

	return cmd
}
