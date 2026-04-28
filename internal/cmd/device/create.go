package device

import (
	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdCreate(f *factory.Factory) *cobra.Command {
	var name, serialNumber string

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Add a device",
		Example: `  devicemanager device create --name "test-router" --serial-number GL5022101241734`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]string{
				"name":         name,
				"serialNumber": serialNumber,
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Post("/api/devices", body)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Device name (required)")
	cmd.Flags().StringVar(&serialNumber, "serial-number", "", "Device serial number (required)")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("serial-number")

	return cmd
}
