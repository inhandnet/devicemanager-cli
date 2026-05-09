package device

import (
	"fmt"
	"net/url"
	"time"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

var trafficFormatters = iostreams.ColumnFormatters{
	"send":    iostreams.FormatBytes,
	"receive": iostreams.FormatBytes,
	"total":   iostreams.FormatBytes,
	"max":     iostreams.FormatBytes,
}

func NewCmdTraffic(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "traffic",
		Short: "Query device traffic statistics",
	}

	cmd.AddCommand(NewCmdTrafficMonthly(f))
	cmd.AddCommand(NewCmdTrafficDaily(f))
	cmd.AddCommand(NewCmdTrafficHourly(f))
	cmd.AddCommand(NewCmdTrafficTop(f))

	return cmd
}

func NewCmdTrafficMonthly(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "monthly <month> <device-id>",
		Short:   "Query monthly traffic for a device (YYYYMM)",
		Args:    cobra.ExactArgs(2),
		Example: `  devicemanager device traffic monthly 202604 <device-id>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			month := args[0]
			deviceID := args[1]

			body := map[string]any{
				"resourceIds": []string{deviceID},
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Post(fmt.Sprintf("/api/traffic_month/list?month=%s", month), body)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(resp, f.IO, output,
				iostreams.WithFormatters(trafficFormatters))
		},
	}

	return cmd
}

func NewCmdTrafficDaily(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "daily <month> <device-id>",
		Short:   "Query daily traffic for a device (YYYYMM)",
		Args:    cobra.ExactArgs(2),
		Example: `  devicemanager device traffic daily 202604 <device-id>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			month := args[0]
			deviceID := args[1]

			q := url.Values{}
			q.Set("month", month)
			q.Set("device_id", deviceID)

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get("/api/traffic_day", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output,
				iostreams.WithFormatters(trafficFormatters))
		},
	}

	return cmd
}

func NewCmdTrafficHourly(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hourly <device-id>",
		Short: "Query hourly traffic for a device",
		Args:  cobra.ExactArgs(1),
		Example: `  # Last 24 hours (default)
  devicemanager device traffic hourly <device-id>

  # Custom date range (max 6 days)
  devicemanager device traffic hourly <device-id> --after 2026-04-25 --before 2026-04-27`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			deviceID := args[0]
			after, _ := cmd.Flags().GetString("after")
			before, _ := cmd.Flags().GetString("before")

			// Default to last 2 days (left-closed, right-open: after ≤ data < before)
			// e.g. today is 9th → after=8th, before=10th → returns 8th and 9th data
			if after == "" {
				after = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
			}
			if before == "" {
				before = time.Now().AddDate(0, 0, 1).Format("2006-01-02")
			}

			q := url.Values{}
			q.Set("after", after)
			q.Set("before", before)

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get(fmt.Sprintf("/api/devices/%s/data-usage/raw", deviceID), q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output,
				iostreams.WithFormatters(trafficFormatters))
		},
	}

	cmd.Flags().String("after", "", "Start date inclusive (YYYY-MM-DD, default: yesterday)")
	cmd.Flags().String("before", "", "End date exclusive (YYYY-MM-DD, default: tomorrow, max 6 days from after)")

	return cmd
}
