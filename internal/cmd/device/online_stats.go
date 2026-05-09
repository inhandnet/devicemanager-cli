package device

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"

	"github.com/inhandnet/devicemanager-cli/internal/cmdutil"
	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdOnlineStats(f *factory.Factory) *cobra.Command {
	var (
		flags     cmdutil.ListFlags
		startTime string
		endTime   string
		name      string
		model     string
		online    string
	)

	cmd := &cobra.Command{
		Use:   "online-stats",
		Short: "Query device online statistics",
		Long: `Query online statistics for devices. Fetches the device list first, then
queries online stats for those devices and merges the results.
Supports pagination and device filtering.`,
		Example: `  # Query online stats for all devices
  devicemanager device online-stats --start-time 2026-05-01 --end-time 2026-05-09

  # Filter by device name
  devicemanager device online-stats --start-time 2026-05-01 --end-time 2026-05-09 --name router

  # Pagination
  devicemanager device online-stats --start-time 2026-05-01 --end-time 2026-05-09 --limit 50`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			startUnix, err := parseDate(startTime)
			if err != nil {
				return fmt.Errorf("invalid --start-time: %w", err)
			}
			endUnix, err := parseDateEnd(endTime)
			if err != nil {
				return fmt.Errorf("invalid --end-time: %w", err)
			}

			// Step 1: Fetch device list
			q := url.Values{}
			flags.ApplyTo(q)
			cmdutil.SetQueryParam(q, "name", name)
			cmdutil.SetQueryParam(q, "model", model)
			cmdutil.SetQueryParam(q, "online", online)

			devResp, err := client.Get("/api/devices", q)
			if err != nil {
				return fmt.Errorf("fetching device list: %w", err)
			}

			// Extract device IDs and info
			results := gjson.GetBytes(devResp, "result")
			if !results.Exists() || len(results.Array()) == 0 {
				fmt.Fprintln(f.IO.ErrOut, "No results.")
				return nil
			}

			var deviceIDs []string
			deviceMap := make(map[string]gjson.Result)
			for _, dev := range results.Array() {
				id := dev.Get("_id").String()
				deviceIDs = append(deviceIDs, id)
				deviceMap[id] = dev
			}

			// Step 2: Query online stats
			statPath := fmt.Sprintf("/api/online_stat/list?start_time=%s&end_time=%s",
				startUnix, endUnix)
			statResp, err := client.Post(statPath, map[string]any{
				"resourceIds": deviceIDs,
			})
			if err != nil {
				return fmt.Errorf("fetching online stats: %w", err)
			}

			// Parse stats into map
			statMap := make(map[string]gjson.Result)
			statsData := gjson.GetBytes(statResp, "result")
			if !statsData.Exists() {
				statsData = gjson.ParseBytes(statResp)
			}
			for _, stat := range statsData.Array() {
				statMap[stat.Get("deviceId").String()] = stat
			}

			// Step 3: Merge device info + stats
			var merged []map[string]any
			for _, id := range deviceIDs {
				dev := deviceMap[id]
				row := map[string]any{
					"_id":          id,
					"name":         dev.Get("name").String(),
					"serialNumber": dev.Get("serialNumber").String(),
					"model":        dev.Get("model").String(),
					"online":       dev.Get("online").Int(),
				}
				if stat, ok := statMap[id]; ok {
					row["onlineRate"] = stat.Get("onlineRate").Float()
					row["maxOnline"] = stat.Get("maxOnline").Int()
					row["maxOffline"] = stat.Get("maxOffline").Int()
					row["totalOnline"] = stat.Get("totalOnline").Int()
					row["totalOffline"] = stat.Get("totalOffline").Int()
					row["login"] = stat.Get("login").Int()
				}
				merged = append(merged, row)
			}

			// Build response with pagination from device list
			total := gjson.GetBytes(devResp, "total").Int()
			cursor := gjson.GetBytes(devResp, "cursor").Int()
			limit := gjson.GetBytes(devResp, "limit").Int()

			resp := map[string]any{
				"result": merged,
				"total":  total,
				"cursor": cursor,
				"limit":  limit,
			}

			respJSON, _ := json.Marshal(resp)
			output, _ := cmd.Flags().GetString("output")

			return iostreams.FormatOutput(respJSON, f.IO, output,
				iostreams.WithColumns("_id", "name", "serialNumber", "model",
					"onlineRate", "maxOnline", "maxOffline", "login"),
				iostreams.WithFormatters(iostreams.ColumnFormatters{
					"maxOnline":  iostreams.FormatDuration,
					"maxOffline": iostreams.FormatDuration,
				}))
		},
	}

	flags.Register(cmd)
	cmd.Flags().StringVar(&startTime, "start-time", "", "Start date inclusive (YYYY-MM-DD) (required)")
	cmd.Flags().StringVar(&endTime, "end-time", "", "End date (YYYY-MM-DD, set to 23:59:59; if today, use current time) (required)")
	cmd.Flags().StringVar(&name, "name", "", "Filter by device name")
	cmd.Flags().StringVar(&model, "model", "", "Filter by device model")
	cmd.Flags().StringVar(&online, "online", "", "Filter by online status (0=offline, 1=online)")
	_ = cmd.MarkFlagRequired("start-time")
	_ = cmd.MarkFlagRequired("end-time")

	return cmd
}

func parseDate(s string) (string, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return strconv.FormatInt(t.Unix(), 10), nil
		}
	}
	return "", fmt.Errorf("unsupported format %q (use YYYY-MM-DD or YYYY-MM-DDTHH:MM:SS)", s)
}

// parseDateEnd parses end time: if only date (YYYY-MM-DD), set to 23:59:59 of that day;
// if that day is today, use current time instead.
func parseDateEnd(s string) (string, error) {
	// If already has time component, parse as-is
	for _, f := range []string{time.RFC3339, "2006-01-02T15:04:05"} {
		if t, err := time.Parse(f, s); err == nil {
			return strconv.FormatInt(t.Unix(), 10), nil
		}
	}
	// Date only: YYYY-MM-DD
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return "", fmt.Errorf("unsupported format %q (use YYYY-MM-DD or YYYY-MM-DDTHH:MM:SS)", s)
	}
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, t.Location())
	if t.Equal(today) {
		// End date is today, use current time
		return strconv.FormatInt(now.Unix(), 10), nil
	}
	// Set to 23:59:59 of that day
	endOfDay := t.Add(24*time.Hour - time.Second)
	return strconv.FormatInt(endOfDay.Unix(), 10), nil
}
