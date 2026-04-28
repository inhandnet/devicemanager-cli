package edge

import (
	"net/url"

	"github.com/spf13/cobra"

	"github.com/inhandnet/elements-cli/internal/factory"
	"github.com/inhandnet/elements-cli/internal/iostreams"
)

func newCmdAppList(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List edge apps",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			if name, _ := cmd.Flags().GetString("name"); name != "" {
				q.Set("name", name)
			}

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get("/api/edge/apps", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output,
				iostreams.WithColumns("_id", "name", "description"))
		},
	}

	cmd.Flags().String("name", "", "Filter by app name")
	return cmd
}
