package device

import (
	"net/url"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdCount(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "count",
		Short: "Query device count trends",
	}

	cmd.AddCommand(newCmdCountOnline(f))
	cmd.AddCommand(newCmdCountTotal(f))

	return cmd
}

func newCmdCountOnline(f *factory.Factory) *cobra.Command {
	var (
		startTime string
		endTime   string
	)

	cmd := &cobra.Command{
		Use:   "online",
		Short: "Query online device count over time",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			if startTime != "" {
				q.Set("start_time", startTime)
			}
			if endTime != "" {
				q.Set("end_time", endTime)
			}

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get("/api/device/count/online", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output)
		},
	}

	cmd.Flags().StringVar(&startTime, "start-time", "", "Start time inclusive, unix timestamp (e.g. 1714492800)")
	cmd.Flags().StringVar(&endTime, "end-time", "", "End time exclusive, unix timestamp (e.g. 1717084800)")

	return cmd
}

func newCmdCountTotal(f *factory.Factory) *cobra.Command {
	var (
		startTime string
		endTime   string
	)

	cmd := &cobra.Command{
		Use:   "total",
		Short: "Query total device count over time",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			if startTime != "" {
				q.Set("start_time", startTime)
			}
			if endTime != "" {
				q.Set("end_time", endTime)
			}

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get("/api/device/count/total", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output)
		},
	}

	cmd.Flags().StringVar(&startTime, "start-time", "", "Start date inclusive (YYYY-MM-DD, e.g. 2026-04-01)")
	cmd.Flags().StringVar(&endTime, "end-time", "", "End date exclusive (YYYY-MM-DD, e.g. 2026-05-01)")

	return cmd
}
