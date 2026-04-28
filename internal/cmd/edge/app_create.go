package edge

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/elements-cli/internal/factory"
	"github.com/inhandnet/elements-cli/internal/iostreams"
)

func newCmdAppCreate(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an edge app",
		Example: `  elements edge app create --name "my-app" --description "My edge app"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			desc, _ := cmd.Flags().GetString("description")

			body := map[string]interface{}{
				"name": name,
			}
			if desc != "" {
				body["description"] = desc
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Post("/api/edge/apps", body)
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "Edge app created\n")
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().String("name", "", "App name (required)")
	cmd.Flags().String("description", "", "App description")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}
