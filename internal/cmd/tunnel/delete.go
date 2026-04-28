package tunnel

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/elements-cli/internal/factory"
	"github.com/inhandnet/elements-cli/internal/iostreams"
)

func NewCmdDelete(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <tunnel-id>",
		Short: "Delete a tunnel",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Delete(fmt.Sprintf("/api/touch/tunnels/%s", args[0]))
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "Tunnel %s deleted\n", args[0])
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	return cmd
}
