package tunnel

import (
	"net/url"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/cmdutil"
	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdList(f *factory.Factory) *cobra.Command {
	var (
		name     string
		deviceID string
		verbose  int
	)

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List tunnels",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			if verbose > 0 {
				q.Set("verbose", strconv.Itoa(verbose))
			}
			cmdutil.SetQueryParam(q, "name", name)
			cmdutil.SetQueryParam(q, "device_id", deviceID)

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get("/api/touch/tunnels", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output,
				iostreams.WithColumns("_id", "name", "deviceId", "proto", "localAddress", "localPort", "status"))
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Filter by tunnel name")
	cmd.Flags().StringVar(&deviceID, "device-id", "", "Filter by device ID")
	cmd.Flags().IntVar(&verbose, "verbose", 100, "Detail level (1-100)")

	return cmd
}
