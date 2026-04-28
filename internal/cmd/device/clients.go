package device

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdClients(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clients",
		Short: "Query gateway clients connected to devices",
	}

	cmd.AddCommand(NewCmdClientsList(f))
	cmd.AddCommand(NewCmdClientsBatch(f))

	return cmd
}

func NewCmdClientsList(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list <device-id>",
		Short: "List clients connected to a device",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get(fmt.Sprintf("/api/devices/%s/gateway/clients", args[0]), nil)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output)
		},
	}

	return cmd
}

func NewCmdClientsBatch(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch <device-id>...",
		Short: "Batch query clients for multiple devices",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]interface{}{
				"deviceIds": args,
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Post("/api/gateway/clients", body)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	return cmd
}
