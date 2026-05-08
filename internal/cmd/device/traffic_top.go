package device

import (
	"net/url"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdTrafficTop(f *factory.Factory) *cobra.Command {
	var (
		date string
		limit string
	)

	cmd := &cobra.Command{
		Use:   "top",
		Short: "Show top devices by monthly traffic",
		Example: `  devicemanager device traffic top --date 202604
  devicemanager device traffic top --date 202604 --limit 10`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			if date != "" {
				q.Set("date", date)
			}
			if limit != "" {
				q.Set("limit", limit)
			}

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get("/api/device/traffic/monthly", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output)
		},
	}

	cmd.Flags().StringVar(&date, "date", "", "Month in YYYYMM format (e.g. 202604)")
	cmd.Flags().StringVar(&limit, "limit", "", "Number of top devices to return")

	return cmd
}
