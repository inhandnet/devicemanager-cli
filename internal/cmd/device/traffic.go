package device

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/inhandnet/elements-cli/internal/factory"
	"github.com/inhandnet/elements-cli/internal/iostreams"
)

func NewCmdTraffic(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "traffic",
		Short: "Query device traffic statistics",
	}

	cmd.AddCommand(NewCmdTrafficMonthly(f))
	cmd.AddCommand(NewCmdTrafficDaily(f))

	return cmd
}

func NewCmdTrafficMonthly(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "monthly <month>",
		Short: "Query monthly traffic for devices (YYYYMM)",
		Args:  cobra.ExactArgs(1),
		Example: `  elements device traffic monthly 202604 --device <device-id>
  elements device traffic monthly 202604 --device <id1> --device <id2>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			month := args[0]
			deviceIDs, _ := cmd.Flags().GetStringSlice("device")

			body := map[string]interface{}{
				"resourceIds": deviceIDs,
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Post(fmt.Sprintf("/api/traffic_month/list?month=%s", month), body)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().StringSlice("device", nil, "Device IDs (required)")
	_ = cmd.MarkFlagRequired("device")

	return cmd
}

func NewCmdTrafficDaily(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "daily <month> <device-id>",
		Short: "Query daily traffic for a device (YYYYMM)",
		Args:  cobra.ExactArgs(2),
		Example: `  elements device traffic daily 202604 <device-id>`,
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

			return iostreams.FormatOutput(body, f.IO, output)
		},
	}

	return cmd
}
