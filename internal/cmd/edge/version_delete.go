package edge

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/elements-cli/internal/factory"
	"github.com/inhandnet/elements-cli/internal/iostreams"
)

func newCmdVersionDelete(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <app-id> <version>",
		Short: "Delete an edge app version",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			appID := args[0]
			version := args[1]

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Delete(fmt.Sprintf("/api/edge/apps/%s/versions/%s", appID, version))
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "Version %s deleted\n", version)
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	return cmd
}
