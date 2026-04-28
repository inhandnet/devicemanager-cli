package edge

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/inhandnet/elements-cli/internal/factory"
	"github.com/inhandnet/elements-cli/internal/iostreams"
)

func newCmdConfigList(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list <app-id>",
		Short:   "List configs of an edge app",
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(1),
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

			body, err := client.Get(fmt.Sprintf("/api/edge/apps/%s/configs", args[0]), q)
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
