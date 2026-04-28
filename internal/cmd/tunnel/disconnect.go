package tunnel

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdDisconnect(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disconnect <tunnel-id>",
		Short: "Disconnect a tunnel",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Put(fmt.Sprintf("/api/touch/tunnels/%s/disconnect", args[0]), nil)
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "Tunnel %s disconnected\n", args[0])
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	return cmd
}
