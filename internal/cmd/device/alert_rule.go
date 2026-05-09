package device

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	inapi "github.com/inhandnet/devicemanager-cli/internal/api"
	"github.com/inhandnet/devicemanager-cli/internal/cmdutil"
	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
	"github.com/inhandnet/devicemanager-cli/internal/ui"
)

// parseTimeToSecOfDay parses "HH:MM" to seconds of day.
func parseTimeToSecOfDay(s string) (int, error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return 0, fmt.Errorf("expected HH:MM format, got %q", s)
	}
	h, err := strconv.Atoi(parts[0])
	if err != nil || h < 0 || h > 23 {
		return 0, fmt.Errorf("invalid hour: %s", parts[0])
	}
	m, err := strconv.Atoi(parts[1])
	if err != nil || m < 0 || m > 59 {
		return 0, fmt.Errorf("invalid minute: %s", parts[1])
	}
	return h*3600 + m*60, nil
}

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
		alertTimeStart string
		alertTimeEnd   string
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

			if alertTimeStart != "" || alertTimeEnd != "" {
				body["customAlertTime"] = true
				_, offset := time.Now().Zone()
				body["offsetTotalSec"] = offset
				if alertTimeStart != "" {
					sec, err := parseTimeToSecOfDay(alertTimeStart)
					if err != nil {
						return fmt.Errorf("invalid --alert-time-start: %w", err)
					}
					body["alertTimeRangeStartSecOfDay"] = sec
				}
				if alertTimeEnd != "" {
					sec, err := parseTimeToSecOfDay(alertTimeEnd)
					if err != nil {
						return fmt.Errorf("invalid --alert-time-end: %w", err)
					}
					body["alertTimeRangeEndSecOfDay"] = sec
				}
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
	cmd.Flags().IntVar(&notifyDelay, "notify-delay", 0, "Delay in minutes before alerting, 0 = immediate (for online/offline types)")
	cmd.Flags().StringVar(&webhookURL, "webhook-url", "", "Webhook callback URL (required when notify-types includes webhook)")
	cmd.Flags().StringVar(&webhookSecret, "webhook-secret", "", "Webhook HMAC signing secret")
	cmd.Flags().StringVar(&locale, "locale", "en", `Notification language: "en" or "zh"`)
	cmd.Flags().StringVar(&alertTimeStart, "alert-time-start", "", "Alert time range start (HH:MM, e.g. 08:00)")
	cmd.Flags().StringVar(&alertTimeEnd, "alert-time-end", "", "Alert time range end (HH:MM, e.g. 18:00)")
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

			// GET current rule — API requires all mutable fields on PUT.
			// Use raw JSON merge to preserve exact field types and values.
			current, err := client.Get(fmt.Sprintf("/api/alert-rules/%s?verbose=100", args[0]), nil)
			if err != nil {
				return fmt.Errorf("fetching current rule: %w", err)
			}

			// Unwrap "result" envelope if present
			var envelope struct {
				Result json.RawMessage `json:"result"`
			}
			ruleJSON := current
			if json.Unmarshal(current, &envelope) == nil && len(envelope.Result) > 0 {
				ruleJSON = envelope.Result
			}

			// Parse into raw map preserving JSON number types
			var rawRule map[string]json.RawMessage
			if err := json.Unmarshal(ruleJSON, &rawRule); err != nil {
				return fmt.Errorf("parsing current rule: %w", err)
			}

			// Build body from allowed mutable fields, preserving original JSON values
			allowedKeys := []string{"name", "forDevice", "notify", "customAlertTime",
				"offsetTotalSec", "alertTimeRangeStartSecOfDay", "alertTimeRangeEndSecOfDay", "webhook"}
			patchRule := make(map[string]json.RawMessage)
			for _, k := range allowedKeys {
				if v, ok := rawRule[k]; ok {
					patchRule[k] = v
				}
			}
			// Ensure required fields have defaults
			if _, ok := patchRule["forDevice"]; !ok {
				patchRule["forDevice"] = json.RawMessage(`{"type":"ALL"}`)
			}
			if _, ok := patchRule["notify"]; !ok {
				patchRule["notify"] = json.RawMessage(`{}`)
			}
			if _, ok := patchRule["webhook"]; !ok {
				patchRule["webhook"] = json.RawMessage(`{}`)
			}

			// Apply user overrides — marshal changed values into the patch
			body := make(map[string]any) // for tracking overrides only
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
			if cmd.Flags().Changed("alert-time-start") || cmd.Flags().Changed("alert-time-end") {
				body["customAlertTime"] = true
				_, offset := time.Now().Zone()
				body["offsetTotalSec"] = offset
				if cmd.Flags().Changed("alert-time-start") {
					v, _ := cmd.Flags().GetString("alert-time-start")
					sec, err := parseTimeToSecOfDay(v)
					if err != nil {
						return fmt.Errorf("invalid --alert-time-start: %w", err)
					}
					body["alertTimeRangeStartSecOfDay"] = sec
				}
				if cmd.Flags().Changed("alert-time-end") {
					v, _ := cmd.Flags().GetString("alert-time-end")
					sec, err := parseTimeToSecOfDay(v)
					if err != nil {
						return fmt.Errorf("invalid --alert-time-end: %w", err)
					}
					body["alertTimeRangeEndSecOfDay"] = sec
				}
			}

			hasChanges := cmd.Flags().Changed("name") || cmd.Flags().Changed("for-device-type") ||
				cmd.Flags().Changed("for-device-value") || cmd.Flags().Changed("notify-users") ||
				cmd.Flags().Changed("notify-types") || cmd.Flags().Changed("notify-delay") ||
				cmd.Flags().Changed("webhook-url") || cmd.Flags().Changed("alert-time-start") ||
				cmd.Flags().Changed("alert-time-end")
			if !hasChanges {
				return fmt.Errorf("at least one flag is required")
			}

			output, _ := cmd.Flags().GetString("output")

			// Merge overrides into patchRule
			for k, v := range body {
				b, _ := json.Marshal(v)
				patchRule[k] = json.RawMessage(b)
			}

			patchJSON, _ := json.Marshal(patchRule)
			resp, err := client.Do("PUT", fmt.Sprintf("/api/alert-rules/%s", args[0]),
				&inapi.RequestOptions{RawBody: bytes.NewReader(patchJSON), ContentType: "application/json"})
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
	cmd.Flags().String("alert-time-start", "", "Alert time range start (HH:MM, e.g. 08:00)")
	cmd.Flags().String("alert-time-end", "", "Alert time range end (HH:MM, e.g. 18:00)")

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
