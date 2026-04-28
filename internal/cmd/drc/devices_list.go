package drc

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
	Status       string
}

func NewCmdDevicesList(f *factory.Factory) *cobra.Command {
	opts := &DevicesListOptions{}

	cmd := &cobra.Command{
		Use:   "list <template-id>",
		Short: "List devices assigned to a DRC template",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			opts.ApplyTo(q)
			cmdutil.SetQueryParam(q, "name", opts.Name)
			cmdutil.SetQueryParam(q, "serialNumber", opts.SerialNumber)
			cmdutil.SetQueryParam(q, "status", opts.Status)

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get(fmt.Sprintf("/api/job/%s/devices", args[0]), q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output,
				iostreams.WithColumns("_id", "name", "serialNumber", "model", "status"))
		},
	}

	opts.Register(cmd)
	cmd.Flags().StringVar(&opts.Name, "name", "", "Filter by device name")
	cmd.Flags().StringVar(&opts.SerialNumber, "serial-number", "", "Filter by serial number")
	cmd.Flags().StringVar(&opts.Status, "status", "", "Filter by status (pending/running/failed/completed)")

	return cmd
}
