package edge

import (
	"net/url"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func newCmdAgentList(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List edge agents",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			if version, _ := cmd.Flags().GetString("version"); version != "" {
				q.Set("version", version)
			}

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get("/api/edge/agents", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output,
				iostreams.WithColumns("_id", "version", "description"))
		},
	}

	cmd.Flags().String("version", "", "Filter by version")
	return cmd
}
