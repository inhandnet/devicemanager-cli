package device

import (
	"net/url"
	"time"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdTrafficTop(f *factory.Factory) *cobra.Command {
	var (
		date  string
		limit string
	)

	cmd := &cobra.Command{
		Use:   "top",
		Short: "Show top devices by monthly traffic",
		Example: `  # Current month
  devicemanager device traffic top

  # Specific month
  devicemanager device traffic top --date 2026-04
  devicemanager device traffic top --date 2026-04 --limit 20`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			// Default to current month
			if date == "" {
				date = time.Now().Format("2006-01")
			}
			// Accept YYYY-MM and append -01
			if len(date) == 7 {
				date += "-01"
			}

			q := url.Values{}
			q.Set("date", date)
			if limit != "" {
				q.Set("limit", limit)
			}

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get("/api/device/traffic/monthly", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output,
				iostreams.WithFormatters(trafficFormatters))
		},
	}

	cmd.Flags().StringVar(&date, "date", "", "Month (YYYY-MM, e.g. 2026-04)")
	cmd.Flags().StringVar(&limit, "limit", "10", "Number of top devices to return (default 10)")

	return cmd
}
