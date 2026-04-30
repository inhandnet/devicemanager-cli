package device

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdAlertConfirm(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "alert-ack <alert-id>",
		Short: "Acknowledge (confirm) an alert",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			resp, err := client.Put(fmt.Sprintf("/api/alerts/%s/confirm", args[0]), nil)
			if err != nil {
				return err
			}

			output, _ := cmd.Flags().GetString("output")
			if output == "table" {
				fmt.Fprintf(f.IO.Out, "%s Alert %s confirmed\n", iostreams.Green("✓"), args[0])
				return nil
			}
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}
}
