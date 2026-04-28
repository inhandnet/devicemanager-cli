package device

import (
	"net/url"

	"github.com/spf13/cobra"

	"github.com/inhandnet/elements-cli/internal/cmdutil"
	"github.com/inhandnet/elements-cli/internal/factory"
	"github.com/inhandnet/elements-cli/internal/iostreams"
)

type ListOptions struct {
	cmdutil.ListFlags
	Name         string
	Model        string
	Online       string
	SerialNumber string
}

func NewCmdList(f *factory.Factory) *cobra.Command {
	opts := &ListOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List devices",
		Aliases: []string{"ls"},
		Example: `  # List all online devices
  elements device list --online 1

  # Filter by model
  elements device list --model IR615

  # List with full details
  elements device list --verbose 100 -o json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			opts.ListFlags.ApplyTo(q)
			cmdutil.SetQueryParam(q, "name", opts.Name)
			cmdutil.SetQueryParam(q, "model", opts.Model)
			cmdutil.SetQueryParam(q, "online", opts.Online)
			cmdutil.SetQueryParam(q, "serial_number", opts.SerialNumber)

			// Auto-set sudo for admin/technical support users
			if cfg, _ := f.Config(); cfg != nil {
				if ctx := cfg.Contexts[cfg.CurrentContext]; ctx != nil {
					if ctx.Authority == "root" || ctx.Authority == "TechnicalSupport" {
						q.Set("sudo", "true")
					}
				}
			}

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get("/api/devices", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output,
				iostreams.WithColumns("_id", "name", "serialNumber", "model", "online"))
		},
	}

	opts.Register(cmd)
	cmd.Flags().StringVar(&opts.Name, "name", "", "Filter by device name")
	cmd.Flags().StringVar(&opts.Model, "model", "", "Filter by model (e.g. IR615, IR900)")
	cmd.Flags().StringVar(&opts.Online, "online", "", "Filter by online status (0=offline, 1=online)")
	cmd.Flags().StringVar(&opts.SerialNumber, "serial-number", "", "Filter by serial number")

	return cmd
}
