package device

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdOnlineEvents(f *factory.Factory) *cobra.Command {
	var startTime, endTime string

	cmd := &cobra.Command{
		Use:   "online-events <device-id>",
		Short: "List device online/offline event timeline",
		Args:  cobra.ExactArgs(1),
		Example: `  # Show online/offline events for the last 24 hours
  devicemanager device online-events 639acc6e078a3d0001be2963 \
    --start-time 2026-04-29 --end-time 2026-04-30

  # Investigate overnight disconnections
  devicemanager device online-events 639acc6e078a3d0001be2963 \
    --start-time 2026-04-29T22:00:00 --end-time 2026-04-30T08:00:00`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			startUnix, err := parseDate(startTime)
			if err != nil {
				return fmt.Errorf("invalid --start-time: %w", err)
			}
			endUnix, err := parseDate(endTime)
			if err != nil {
				return fmt.Errorf("invalid --end-time: %w", err)
			}

			q := url.Values{}
			q.Set("start_time", startUnix)
			q.Set("end_time", endUnix)

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get(fmt.Sprintf("/api/devices/%s/online-events", args[0]), q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output)
		},
	}

	cmd.Flags().StringVar(&startTime, "start-time", "", "Start time (YYYY-MM-DD or ISO8601) (required)")
	cmd.Flags().StringVar(&endTime, "end-time", "", "End time (YYYY-MM-DD or ISO8601) (required)")
	_ = cmd.MarkFlagRequired("start-time")
	_ = cmd.MarkFlagRequired("end-time")

	return cmd
}
