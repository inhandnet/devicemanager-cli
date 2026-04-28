package edge

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func newCmdConfigUpdate(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <app-id> <config-id>",
		Short: "Update an edge app config",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			appID := args[0]
			configID := args[1]

			body := map[string]interface{}{}
			if desc, _ := cmd.Flags().GetString("description"); desc != "" {
				body["description"] = desc
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Put(fmt.Sprintf("/api/edge/apps/%s/configs/%s", appID, configID), body)
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "Config %s updated\n", configID)
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().String("description", "", "New description")
	return cmd
}
