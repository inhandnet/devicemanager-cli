package edge

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/elements-cli/internal/factory"
	"github.com/inhandnet/elements-cli/internal/iostreams"
)

func newCmdAppDelete(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <app-id>",
		Short: "Delete an edge app",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Delete(fmt.Sprintf("/api/edge/apps/%s", args[0]))
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "Edge app %s deleted\n", args[0])
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	return cmd
}
