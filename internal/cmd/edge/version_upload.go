package edge

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func newCmdVersionUpload(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload <file-path>",
		Short: "Upload an edge app version",
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

			appID, _ := cmd.Flags().GetString("app")
			if appID == "" {
				return fmt.Errorf("--app is required")
			}

			output, _ := cmd.Flags().GetString("output")

			uploadURL := fmt.Sprintf("/api/edge/apps/upload?app=%s", appID)
			resp, err := client.Upload(uploadURL, "file", filepath.Base(filePath), file)
			if err != nil {
				return err
			}
			fmt.Fprintf(f.IO.Out, "Version uploaded\n")
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().String("app", "", "Edge app ID (required)")
	_ = cmd.MarkFlagRequired("app")

	return cmd
}
