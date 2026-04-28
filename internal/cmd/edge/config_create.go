package edge

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func newCmdConfigCreate(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <app-id>",
		Short: "Create an edge app config",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			version, _ := cmd.Flags().GetString("version")
			content, _ := cmd.Flags().GetString("content")
			desc, _ := cmd.Flags().GetString("description")

			body := map[string]interface{}{
				"version": version,
				"content": content,
			}
			if desc != "" {
				body["description"] = desc
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Post(fmt.Sprintf("/api/edge/apps/%s/configs", args[0]), body)
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "Config created\n")
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().String("version", "", "Config version (required)")
	cmd.Flags().String("content", "", "Config content (required)")
	cmd.Flags().String("description", "", "Description")
	_ = cmd.MarkFlagRequired("version")
	_ = cmd.MarkFlagRequired("content")

	return cmd
}
