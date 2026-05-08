package system

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

type LogListOptions struct {
	cmdutil.ListFlags
	StartTime string
	EndTime   string
	Level     string
}

func NewCmdLog(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "log",
		Short: "View audit logs",
	}

	cmd.AddCommand(newCmdLogList(f))

	return cmd
}

func newCmdLogList(f *factory.Factory) *cobra.Command {
	opts := &LogListOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List audit logs",
		Aliases: []string{"ls"},
		Example: `  # List audit logs (defaults to last 7 days)
  devicemanager system log list

  # Filter by date range
  devicemanager system log list --start-time 2026-04-24 --end-time 2026-04-30

  # Filter by level
  devicemanager system log list --level warning`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			opts.ApplyTo(q)
			q.Set("language", "1")

			// Default to last 7 days if --start-time not specified
			if opts.StartTime == "" {
				opts.StartTime = time.Now().AddDate(0, 0, -7).Format("2006-01-02")
			}
			if opts.EndTime == "" {
				opts.EndTime = time.Now().Format("2006-01-02")
			}

			ts, err := parseDateToUnix(opts.StartTime)
			if err != nil {
				return fmt.Errorf("invalid --start-time: %w", err)
			}
			q.Set("start_time", ts)

			ts, err = parseDateToUnix(opts.EndTime)
			if err != nil {
				return fmt.Errorf("invalid --end-time: %w", err)
			}
			q.Set("end_time", ts)
			cmdutil.SetQueryParam(q, "level", opts.Level)

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get("/api2/behav_log", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output,
				iostreams.WithColumns("user", "message", "ip", "createdAt"))
		},
	}

	opts.Register(cmd)
	cmd.Flags().StringVar(&opts.StartTime, "start-time", "", "Start date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&opts.EndTime, "end-time", "", "End date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&opts.Level, "level", "", "Filter by level (info, warning, error)")

	return cmd
}

func parseDateToUnix(s string) (string, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(t.Unix(), 10), nil
}
