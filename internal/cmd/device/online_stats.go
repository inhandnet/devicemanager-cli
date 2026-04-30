package device

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdOnlineStats(f *factory.Factory) *cobra.Command {
	var (
		startTime string
		endTime   string
		deviceIDs []string
	)

	cmd := &cobra.Command{
		Use:   "online-stats",
		Short: "Query device online statistics",
		Example: `  # Query online stats for specific devices in the last 7 days
  devicemanager device online-stats --device-id 639acc6e078a3d0001be2963 \
    --start-time 2026-04-23 --end-time 2026-04-30

  # Query multiple devices
  devicemanager device online-stats \
    --device-id 639acc6e078a3d0001be2963 \
    --device-id 5d6349d6335c8c000178a194 \
    --start-time 2026-04-23 --end-time 2026-04-30`,
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

			body := map[string]any{
				"resourceIds": deviceIDs,
			}

			path := fmt.Sprintf("/api/online_stat/list?start_time=%s&end_time=%s",
				startUnix, endUnix)

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Post(path, body)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(resp, f.IO, output,
				iostreams.WithColumns("deviceId", "maxOnline", "maxOffline",
					"totalOnline", "totalOffline", "login", "onlineRate"))
		},
	}

	cmd.Flags().StringVar(&startTime, "start-time", "", "Start date (YYYY-MM-DD) (required)")
	cmd.Flags().StringVar(&endTime, "end-time", "", "End date (YYYY-MM-DD) (required)")
	cmd.Flags().StringSliceVar(&deviceIDs, "device-id", nil, "Device ID(s) to query (required, repeatable)")
	_ = cmd.MarkFlagRequired("start-time")
	_ = cmd.MarkFlagRequired("end-time")
	_ = cmd.MarkFlagRequired("device-id")

	return cmd
}

func parseDate(s string) (string, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return strconv.FormatInt(t.Unix(), 10), nil
		}
	}
	return "", fmt.Errorf("unsupported format %q (use YYYY-MM-DD or YYYY-MM-DDTHH:MM:SS)", s)
}
