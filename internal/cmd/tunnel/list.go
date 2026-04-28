package tunnel

import (
	"net/url"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/cmdutil"
	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

type ListOptions struct {
	cmdutil.ListFlags
	Name     string
	DeviceID string
}

func NewCmdList(f *factory.Factory) *cobra.Command {
	opts := &ListOptions{}

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
			opts.ApplyTo(q)
			cmdutil.SetQueryParam(q, "name", opts.Name)
			cmdutil.SetQueryParam(q, "device_id", opts.DeviceID)

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get("/api/touch/tunnels", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output,
				iostreams.WithColumns("_id", "name", "deviceId", "proto", "localAddress", "localPort", "status"))
		},
	}

	opts.Register(cmd)
	cmd.Flags().StringVar(&opts.Name, "name", "", "Filter by tunnel name")
	cmd.Flags().StringVar(&opts.DeviceID, "device-id", "", "Filter by device ID")

	return cmd
}
