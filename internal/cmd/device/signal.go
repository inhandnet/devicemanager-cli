package device

import (
	"net/url"

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
		Example: `  devicemanager device signal 5e6f222afbcf3e0001e133f4 \
    --after 2024-01-01T00:00:00Z --before 2024-01-02T00:00:00Z`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
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

	cmd.Flags().StringVar(&after, "after", "", "Start time (ISO 8601, e.g. 2024-01-01T00:00:00Z) (required)")
	cmd.Flags().StringVar(&before, "before", "", "End time (ISO 8601, e.g. 2024-01-02T00:00:00Z) (required)")
	_ = cmd.MarkFlagRequired("after")
	_ = cmd.MarkFlagRequired("before")

	return cmd
}
