package device

import (
	"net/url"

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
			cmdutil.SetQueryParam(q, "start_time", opts.StartTime)
			cmdutil.SetQueryParam(q, "end_time", opts.EndTime)
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
	cmd.Flags().StringVar(&opts.StartTime, "start-time", "", "Filter by start time (ISO8601)")
	cmd.Flags().StringVar(&opts.EndTime, "end-time", "", "Filter by end time (ISO8601)")
	cmd.Flags().StringVar(&opts.State, "state", "", "Filter by state (confirmed/unconfirmed)")

	return cmd
}
