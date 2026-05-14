package edge

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/cmdutil"
	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func newCmdAppDevices(f *factory.Factory) *cobra.Command {
	opts := &cmdutil.ListFlags{}

	cmd := &cobra.Command{
		Use:     "devices <app-id>",
		Short:   "List devices with an edge app deployed",
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			opts.ApplyTo(q)

			status, _ := cmd.Flags().GetString("status")
			q.Set("status", status)

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get(fmt.Sprintf("/api/edge/apps/%s/devices", args[0]), q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output,
				iostreams.WithColumns("deviceId", "status", "version", "currentVersion"))
		},
	}

	opts.Register(cmd)
	cmd.Flags().String("status", "PENDING", "Filter by status (PENDING, INSTALLING, DOWNLOADING, READY, FAILED)")

	return cmd
}
