package edge

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func newCmdAgentUpdate(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <agent-id>",
		Short: "Update an edge agent",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]interface{}{}
			if desc, _ := cmd.Flags().GetString("description"); desc != "" {
				body["description"] = desc
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Put(fmt.Sprintf("/api/edge/agents/%s", args[0]), body)
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "Agent %s updated\n", args[0])
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().String("description", "", "New description")
	return cmd
}
