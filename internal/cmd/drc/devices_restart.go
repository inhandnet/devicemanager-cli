package drc

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/api"
	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdDevicesRestart(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restart <template-id> <device-id>",
		Short: "Restart a device task for a DRC template",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			templateID := args[0]
			deviceID := args[1]

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Do("PUT", fmt.Sprintf("/api/jobs/%s/devices/%s/restart", templateID, deviceID), &api.RequestOptions{
				Query: oidQuery(f),
			})
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "Restarted task for device %s in template %s\n", deviceID, templateID)
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	return cmd
}
