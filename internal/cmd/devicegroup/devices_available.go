package devicegroup

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/cmdutil"
	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

type DevicesAvailableOptions struct {
	cmdutil.ListFlags
	Name         string
	SerialNumber string
}

func NewCmdDevicesAvailable(f *factory.Factory) *cobra.Command {
	opts := &DevicesAvailableOptions{}

	cmd := &cobra.Command{
		Use:   "available <group-id>",
		Short: "List devices that can be added to a group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			opts.ApplyTo(q)
			cmdutil.SetQueryParam(q, "name", opts.Name)
			cmdutil.SetQueryParam(q, "serial_number", opts.SerialNumber)

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get(fmt.Sprintf("/api/devicegroups/%s/devices/exclusion", args[0]), q)
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

	return cmd
}
