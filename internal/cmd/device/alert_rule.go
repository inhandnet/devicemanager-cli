package device

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/cmdutil"
	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
	"github.com/inhandnet/devicemanager-cli/internal/ui"
)

func NewCmdAlertRule(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "alert-rule",
		Short:   "Manage alert rules",
		Aliases: []string{"alert-rules"},
	}

	cmd.AddCommand(newCmdAlertRuleList(f))
	cmd.AddCommand(newCmdAlertRuleGet(f))
	cmd.AddCommand(newCmdAlertRuleCreate(f))
	cmd.AddCommand(newCmdAlertRuleUpdate(f))
	cmd.AddCommand(newCmdAlertRuleDelete(f))
	cmd.AddCommand(newCmdAlertRuleEnable(f))
	cmd.AddCommand(newCmdAlertRuleDisable(f))

	return cmd
}

func newCmdAlertRuleList(f *factory.Factory) *cobra.Command {
	var (
		flags      cmdutil.ListFlags
		deviceName string
	)

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List alert rules",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			flags.ApplyTo(q)
			q.Set("verbose", "100")
			cmdutil.SetQueryParam(q, "device_name", deviceName)

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get("/api/alert-rules", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output,
				iostreams.WithColumns("_id", "name", "metric", "condition", "enabled"))
		},
	}

	flags.Register(cmd)
	cmd.Flags().StringVar(&deviceName, "device-name", "", "Filter by device name")

	return cmd
}

func newCmdAlertRuleGet(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get <rule-id>",
		Short: "Get alert rule details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get(fmt.Sprintf("/api/alert-rules/%s?verbose=100", args[0]), nil)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output)
		},
	}
}

func newCmdAlertRuleCreate(f *factory.Factory) *cobra.Command {
	var (
		name      string
		metric    string
		condition string
		threshold string
		duration  string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an alert rule",
		Example: `  devicemanager device alert-rule create --name "offline-alert" \
    --metric online --condition eq --threshold 0 --duration 300`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{
				"name":      name,
				"metric":    metric,
				"condition": condition,
				"threshold": threshold,
			}
			if duration != "" {
				body["duration"] = duration
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Post("/api/alert-rules", body)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Rule name (required)")
	cmd.Flags().StringVar(&metric, "metric", "", "Metric to monitor (required)")
	cmd.Flags().StringVar(&condition, "condition", "", "Condition operator (required)")
	cmd.Flags().StringVar(&threshold, "threshold", "", "Threshold value (required)")
	cmd.Flags().StringVar(&duration, "duration", "", "Duration in seconds")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("metric")
	_ = cmd.MarkFlagRequired("condition")
	_ = cmd.MarkFlagRequired("threshold")

	return cmd
}

func newCmdAlertRuleUpdate(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <rule-id>",
		Short: "Update an alert rule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{}
			if cmd.Flags().Changed("name") {
				v, _ := cmd.Flags().GetString("name")
				body["name"] = v
			}
			if cmd.Flags().Changed("threshold") {
				v, _ := cmd.Flags().GetString("threshold")
				body["threshold"] = v
			}

			if len(body) == 0 {
				return fmt.Errorf("at least one flag is required")
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Put(fmt.Sprintf("/api/alert-rules/%s", args[0]), body)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().String("name", "", "New rule name")
	cmd.Flags().String("threshold", "", "New threshold value")

	return cmd
}

func newCmdAlertRuleDelete(f *factory.Factory) *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <rule-id>",
		Short: "Delete an alert rule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !ui.Confirm(f.IO, fmt.Sprintf("Delete alert rule %s?", args[0]), yes) {
				return nil
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Delete(fmt.Sprintf("/api/alert-rules/%s", args[0]))
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip confirmation prompt")

	return cmd
}

func newCmdAlertRuleEnable(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "enable <rule-id>",
		Short: "Enable an alert rule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Put(fmt.Sprintf("/api/alert-rules/%s/enable", args[0]), nil)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}
}

func newCmdAlertRuleDisable(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "disable <rule-id>",
		Short: "Disable an alert rule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Put(fmt.Sprintf("/api/alert-rules/%s/disable", args[0]), nil)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}
}
