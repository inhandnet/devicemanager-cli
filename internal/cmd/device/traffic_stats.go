package device

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"

	"github.com/inhandnet/devicemanager-cli/internal/cmdutil"
	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdTrafficStats(f *factory.Factory) *cobra.Command {
	var (
		flags  cmdutil.ListFlags
		after  string
		before string
		name   string
		model  string
		online string
	)

	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Query device traffic statistics",
		Long: `Query traffic statistics for devices. Fetches the device list first, then
queries traffic data for those devices and merges the results.
Supports pagination and device filtering.`,
		Example: `  # Query traffic stats for all devices
  devicemanager device traffic stats --after 2026-05-01 --before 2026-05-09

  # Filter by device name
  devicemanager device traffic stats --after 2026-05-01 --before 2026-05-09 --name router

  # Only online devices
  devicemanager device traffic stats --after 2026-05-01 --before 2026-05-09 --online 1`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			if after == "" || before == "" {
				return fmt.Errorf("--after and --before are required")
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

			// Step 2: Query traffic stats
			tq := url.Values{}
			tq.Set("after", after)
			tq.Set("before", before)

			statResp, err := client.Post("/api/traffic/list?"+tq.Encode(), map[string]any{
				"resourceIds": deviceIDs,
			})
			if err != nil {
				return fmt.Errorf("fetching traffic stats: %w", err)
			}

			statMap := make(map[string]gjson.Result)
			statsData := gjson.GetBytes(statResp, "result")
			if !statsData.Exists() {
				statsData = gjson.ParseBytes(statResp)
			}
			for _, stat := range statsData.Array() {
				statMap[stat.Get("deviceId").String()] = stat
			}

			// Step 3: Merge
			var merged []map[string]any
			for _, id := range deviceIDs {
				dev := deviceMap[id]
				row := map[string]any{
					"_id":          id,
					"name":         dev.Get("name").String(),
					"serialNumber": dev.Get("serialNumber").String(),
					"model":        dev.Get("model").String(),
				}
				if stat, ok := statMap[id]; ok {
					row["send"] = stat.Get("send").Int()
					row["receive"] = stat.Get("receive").Int()
					row["total"] = stat.Get("total").Int()
				} else {
					row["send"] = 0
					row["receive"] = 0
					row["total"] = 0
				}
				merged = append(merged, row)
			}

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
				iostreams.WithColumns("_id", "name", "serialNumber", "model", "send", "receive", "total"),
				iostreams.WithFormatters(trafficFormatters))
		},
	}

	flags.Register(cmd)
	cmd.Flags().StringVar(&after, "after", "", "Start date inclusive (YYYY-MM-DD) (required)")
	cmd.Flags().StringVar(&before, "before", "", "End date exclusive (YYYY-MM-DD) (required)")
	cmd.Flags().StringVar(&name, "name", "", "Filter by device name")
	cmd.Flags().StringVar(&model, "model", "", "Filter by device model")
	cmd.Flags().StringVar(&online, "online", "", "Filter by online status (0=offline, 1=online)")
	_ = cmd.MarkFlagRequired("after")
	_ = cmd.MarkFlagRequired("before")

	return cmd
}
