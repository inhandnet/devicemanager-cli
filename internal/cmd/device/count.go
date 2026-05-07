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
		after  string
		before string
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
			if after != "" {
				q.Set("after", after)
			}
			if before != "" {
				q.Set("before", before)
			}

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get("/api/device/count/online", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output)
		},
	}

	cmd.Flags().StringVar(&after, "after", "", "Start date (ISO format)")
	cmd.Flags().StringVar(&before, "before", "", "End date (ISO format)")

	return cmd
}

func newCmdCountTotal(f *factory.Factory) *cobra.Command {
	var (
		after  string
		before string
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
			if after != "" {
				q.Set("after", after)
			}
			if before != "" {
				q.Set("before", before)
			}

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get("/api/device/count/total", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output)
		},
	}

	cmd.Flags().StringVar(&after, "after", "", "Start date (ISO format)")
	cmd.Flags().StringVar(&before, "before", "", "End date (ISO format)")

	return cmd
}
