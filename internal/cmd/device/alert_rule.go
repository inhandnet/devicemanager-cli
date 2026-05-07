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
		name           string
		alertType      string
		forDeviceType  string
		forDeviceValue []string
		notifyUsers    []string
		notifyTypes    []string
		notifyDelay    int
		webhookURL     string
		webhookSecret  string
		locale         string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an alert rule",
		Example: `  # Create an offline alert for all devices
  devicemanager device alert-rule create --name "offline-alert" --alert-type offline

  # Create an alert for specific devices with email notification
  devicemanager device alert-rule create --name "offline-alert" --alert-type offline \
    --for-device-type DEVICE --for-device-value id1,id2 \
    --notify-users uid1,uid2 --notify-types email

  # Create with webhook notification
  devicemanager device alert-rule create --name "traffic-alert" \
    --alert-type daily_traffic_excess \
    --notify-types webhook --webhook-url https://example.com/hook`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			forDevice := map[string]any{"type": forDeviceType}
			if len(forDeviceValue) > 0 {
				forDevice["value"] = forDeviceValue
			}

			notify := map[string]any{}
			if len(notifyUsers) > 0 {
				notify["users"] = notifyUsers
			}
			if len(notifyTypes) > 0 {
				notify["types"] = notifyTypes
			}
			if notifyDelay > 0 {
				notify["delay"] = notifyDelay
			}

			body := map[string]any{
				"name":      name,
				"alertType": alertType,
				"locale":    locale,
				"forDevice": forDevice,
				"notify":    notify,
			}

			if webhookURL != "" {
				wh := map[string]any{"url": webhookURL}
				if webhookSecret != "" {
					wh["secret"] = webhookSecret
				}
				body["webhook"] = wh
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
	cmd.Flags().StringVar(&alertType, "alert-type", "", "Alert type: online, offline, hourly_traffic_excess, daily_traffic_excess, monthly_traffic_excess, sim_switch, link_backup, link_change, power_switch, power_fault, power_recovery (required)")
	cmd.Flags().StringVar(&forDeviceType, "for-device-type", "ALL", "Target scope: ALL, DEVICE_GROUP, or DEVICE")
	cmd.Flags().StringSliceVar(&forDeviceValue, "for-device-value", nil, "Target device/group IDs (comma-separated)")
	cmd.Flags().StringSliceVar(&notifyUsers, "notify-users", nil, "User IDs to notify (comma-separated)")
	cmd.Flags().StringSliceVar(&notifyTypes, "notify-types", nil, "Notification types: email, sms, webhook (comma-separated)")
	cmd.Flags().IntVar(&notifyDelay, "notify-delay", 0, "Delay in minutes before alerting (for online/offline)")
	cmd.Flags().StringVar(&webhookURL, "webhook-url", "", "Webhook URL")
	cmd.Flags().StringVar(&webhookSecret, "webhook-secret", "", "Webhook secret")
	cmd.Flags().StringVar(&locale, "locale", "en", `Notification language: "en" or "zh"`)
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("alert-type")

	return cmd
}

func newCmdAlertRuleUpdate(f *factory.Factory) *cobra.Command {
	var (
		forDeviceValue []string
		notifyUsers    []string
		notifyTypes    []string
	)

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
			if cmd.Flags().Changed("for-device-type") || cmd.Flags().Changed("for-device-value") {
				forDevice := map[string]any{}
				if cmd.Flags().Changed("for-device-type") {
					v, _ := cmd.Flags().GetString("for-device-type")
					forDevice["type"] = v
				}
				if len(forDeviceValue) > 0 {
					forDevice["value"] = forDeviceValue
				}
				body["forDevice"] = forDevice
			}
			if cmd.Flags().Changed("notify-users") || cmd.Flags().Changed("notify-types") || cmd.Flags().Changed("notify-delay") {
				notify := map[string]any{}
				if len(notifyUsers) > 0 {
					notify["users"] = notifyUsers
				}
				if len(notifyTypes) > 0 {
					notify["types"] = notifyTypes
				}
				if cmd.Flags().Changed("notify-delay") {
					v, _ := cmd.Flags().GetInt("notify-delay")
					notify["delay"] = v
				}
				body["notify"] = notify
			}
			if cmd.Flags().Changed("webhook-url") {
				wh := map[string]any{}
				v, _ := cmd.Flags().GetString("webhook-url")
				wh["url"] = v
				if cmd.Flags().Changed("webhook-secret") {
					s, _ := cmd.Flags().GetString("webhook-secret")
					wh["secret"] = s
				}
				body["webhook"] = wh
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
	cmd.Flags().String("for-device-type", "", "Target scope: ALL, DEVICE_GROUP, or DEVICE")
	cmd.Flags().StringSliceVar(&forDeviceValue, "for-device-value", nil, "Target device/group IDs (comma-separated)")
	cmd.Flags().StringSliceVar(&notifyUsers, "notify-users", nil, "User IDs to notify (comma-separated)")
	cmd.Flags().StringSliceVar(&notifyTypes, "notify-types", nil, "Notification types: email, sms, webhook (comma-separated)")
	cmd.Flags().Int("notify-delay", 0, "Delay in minutes before alerting")
	cmd.Flags().String("webhook-url", "", "Webhook URL")
	cmd.Flags().String("webhook-secret", "", "Webhook secret")

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
