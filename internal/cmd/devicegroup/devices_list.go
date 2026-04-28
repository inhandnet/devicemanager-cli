package devicegroup

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/inhandnet/elements-cli/internal/cmdutil"
	"github.com/inhandnet/elements-cli/internal/factory"
	"github.com/inhandnet/elements-cli/internal/iostreams"
)

type DevicesListOptions struct {
	cmdutil.ListFlags
	Name         string
	SerialNumber string
	Online       string
	Recursive    bool
}

func NewCmdDevicesList(f *factory.Factory) *cobra.Command {
	opts := &DevicesListOptions{}

	cmd := &cobra.Command{
		Use:   "list <group-id>",
		Short: "List devices in a group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			opts.ListFlags.ApplyTo(q)
			cmdutil.SetQueryParam(q, "name", opts.Name)
			cmdutil.SetQueryParam(q, "serial_number", opts.SerialNumber)
			cmdutil.SetQueryParam(q, "online", opts.Online)
			if opts.Recursive {
				q.Set("recursive", "1")
			}

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get(fmt.Sprintf("/api/devicegroups/%s/devices", args[0]), q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output,
				iostreams.WithColumns("_id", "name", "serialNumber", "model", "online"))
		},
	}

	opts.Register(cmd)
	cmd.Flags().StringVar(&opts.Name, "name", "", "Filter by device name")
	cmd.Flags().StringVar(&opts.SerialNumber, "serial-number", "", "Filter by serial number")
	cmd.Flags().StringVar(&opts.Online, "online", "", "Filter by online status (0=offline, 1=online)")
	cmd.Flags().BoolVar(&opts.Recursive, "recursive", false, "Include devices from sub-groups")

	return cmd
}
