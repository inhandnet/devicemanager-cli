package edge

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func newCmdAgentUpload(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload <file-path>",
		Short: "Upload an edge agent",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			filePath := args[0]
			file, err := os.Open(filePath)
			if err != nil {
				return fmt.Errorf("opening file: %w", err)
			}
			defer file.Close()

			desc, _ := cmd.Flags().GetString("description")

			output, _ := cmd.Flags().GetString("output")

			uploadURL := "/api/edge/agents/upload"
			if desc != "" {
				uploadURL = fmt.Sprintf("%s?description=%s", uploadURL, desc)
			}
			resp, err := client.Upload(uploadURL, "file", filePath, file)
			if err != nil {
				return err
			}
			fmt.Fprintf(f.IO.Out, "Agent uploaded\n")
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().String("description", "", "Agent description")
	return cmd
}
