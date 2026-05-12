package device

import (
	"net/url"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdModels(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "models",
		Short:   "List supported device models",
		Aliases: []string{"model"},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			output, _ := cmd.Flags().GetString("output")

			q := url.Values{}
			q.Set("limit", "0")

			body, err := client.Get("/api/models", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output,
				iostreams.WithColumns("_id", "name", "model"))
		},
	}
}
