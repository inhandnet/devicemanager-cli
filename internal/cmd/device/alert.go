package device

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/cmdutil"
	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

type AlertOptions struct {
	cmdutil.ListFlags
	DeviceName string
	RuleName   string
	StartTime  string
	EndTime    string
	State      string
}

func NewCmdAlert(f *factory.Factory) *cobra.Command {
	opts := &AlertOptions{}

	cmd := &cobra.Command{
		Use:     "alert",
		Short:   "List device alerts",
		Aliases: []string{"alerts"},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			opts.ApplyTo(q)
			cmdutil.SetQueryParam(q, "device_name", opts.DeviceName)
			cmdutil.SetQueryParam(q, "rule_name", opts.RuleName)
			if opts.StartTime != "" {
				ts, err := toUnixTimestamp(opts.StartTime)
				if err != nil {
					return fmt.Errorf("invalid --start-time: %w", err)
				}
				q.Set("start_time", ts)
			}
			if opts.EndTime != "" {
				ts, err := toUnixTimestamp(opts.EndTime)
				if err != nil {
					return fmt.Errorf("invalid --end-time: %w", err)
				}
				q.Set("end_time", ts)
			}
			cmdutil.SetQueryParam(q, "state", opts.State)

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get("/api/alerts", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output,
				iostreams.WithColumns("_id", "deviceName", "ruleName", "state", "createTime"))
		},
	}

	opts.Register(cmd)
	cmd.Flags().StringVar(&opts.DeviceName, "device-name", "", "Filter by device name")
	cmd.Flags().StringVar(&opts.RuleName, "rule-name", "", "Filter by rule name")
	cmd.Flags().StringVar(&opts.StartTime, "start-time", "", "Filter by start time (YYYY-MM-DD or unix timestamp)")
	cmd.Flags().StringVar(&opts.EndTime, "end-time", "", "Filter by end time (YYYY-MM-DD or unix timestamp)")
	cmd.Flags().StringVar(&opts.State, "state", "", "Filter by state (confirmed/unconfirmed)")

	return cmd
}

// toUnixTimestamp converts a date string (YYYY-MM-DD) or unix timestamp string to unix timestamp string.
func toUnixTimestamp(s string) (string, error) {
	// Already a unix timestamp?
	if _, err := strconv.ParseInt(s, 10, 64); err == nil {
		return s, nil
	}
	// Try YYYY-MM-DD
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return "", fmt.Errorf("expected YYYY-MM-DD or unix timestamp, got %q", s)
	}
	return strconv.FormatInt(t.Unix(), 10), nil
}
