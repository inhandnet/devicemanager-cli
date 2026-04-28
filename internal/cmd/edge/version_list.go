package edge

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/elements-cli/internal/factory"
	"github.com/inhandnet/elements-cli/internal/iostreams"
)

func newCmdVersionList(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list <app-id>",
		Short:   "List versions of an edge app",
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get(fmt.Sprintf("/api/edge/apps/%s/versions", args[0]), nil)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output,
				iostreams.WithColumns("version", "createTime"))
		},
	}

	return cmd
}
