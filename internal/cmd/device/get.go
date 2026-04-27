package device

import (
	"net/url"

	"github.com/spf13/cobra"

	"github.com/inhandnet/elements-cli/internal/factory"
	"github.com/inhandnet/elements-cli/internal/iostreams"
)

func NewCmdGet(f *factory.Factory) *cobra.Command {
	var verbose int

	cmd := &cobra.Command{
		Use:   "get <device-id>",
		Short: "Get device details",
		Args:  cobra.ExactArgs(1),
		Example: `  # Get full device details
  elements device get 5d6349d6335c8c000178a194 --verbose 100`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			q.Set("verbose", "100")
			if cmd.Flags().Changed("verbose") {
				q.Set("verbose", cmd.Flags().Lookup("verbose").Value.String())
			}

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get("/api/devices/"+args[0], q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output)
		},
	}

	cmd.Flags().IntVar(&verbose, "verbose", 100, "Detail level (1-100)")

	return cmd
}
