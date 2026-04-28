package edge

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/elements-cli/internal/factory"
	"github.com/inhandnet/elements-cli/internal/iostreams"
)

func newCmdVersionUpdate(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <app-id> <version>",
		Short: "Update edge app version release notes",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			appID := args[0]
			version := args[1]

			body := map[string]interface{}{}
			if notes, _ := cmd.Flags().GetString("notes"); notes != "" {
				body["releaseNotes"] = notes
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Put(fmt.Sprintf("/api/edge/apps/%s/versions/%s", appID, version), body)
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "Version %s updated\n", version)
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().String("notes", "", "Release notes")
	return cmd
}
