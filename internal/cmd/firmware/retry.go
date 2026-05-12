package firmware

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdRetry(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "retry <job-id> <device-id>",
		Short: "Retry a firmware upgrade task for a device",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			jobID := args[0]
			deviceID := args[1]

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Put(fmt.Sprintf("/api/jobs/%s/devices/%s/restart", jobID, deviceID), nil)
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "Retried upgrade for device %s in job %s\n", deviceID, jobID)
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	return cmd
}
