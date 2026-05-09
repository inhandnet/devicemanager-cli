package device

import (
	"net/url"
	"time"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdSignal(f *factory.Factory) *cobra.Command {
	var after, before string

	cmd := &cobra.Command{
		Use:   "signal <device-id>",
		Short: "Query historical signal quality",
		Args:  cobra.ExactArgs(1),
		Example: `  # Query signal from a start time to now
  devicemanager device signal <device-id> --after 2026-05-01T00:00:00Z

  # Query signal for a specific time range
  devicemanager device signal <device-id> --after 2026-05-01T00:00:00Z --before 2026-05-02T00:00:00Z`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			if before == "" {
				before = time.Now().UTC().Format(time.RFC3339)
			}

			q := url.Values{}
			q.Set("after", after)
			q.Set("before", before)

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get("/api/devices/"+args[0]+"/signal-quality", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output)
		},
	}

	cmd.Flags().StringVar(&after, "after", "", "Start time inclusive (ISO 8601, e.g. 2026-05-01T00:00:00Z) (required)")
	cmd.Flags().StringVar(&before, "before", "", "End time exclusive (ISO 8601, defaults to now)")
	_ = cmd.MarkFlagRequired("after")

	return cmd
}
