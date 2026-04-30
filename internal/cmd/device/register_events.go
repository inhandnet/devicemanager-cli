package device

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/cmdutil"
	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdRegisterEvents(f *factory.Factory) *cobra.Command {
	var flags cmdutil.ListFlags

	cmd := &cobra.Command{
		Use:   "register-events <serial-number>",
		Short: "List device registration events by serial number",
		Args:  cobra.ExactArgs(1),
		Example: `  # View registration history for a device
  devicemanager device register-events GL5021937123456`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			flags.ApplyTo(q)

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get(fmt.Sprintf("/api/devices/%s/register-events", args[0]), q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output)
		},
	}

	flags.Register(cmd)

	return cmd
}
