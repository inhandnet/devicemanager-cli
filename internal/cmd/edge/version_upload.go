package edge

import (
	"fmt"
	"os"

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

			resp, err := client.Upload("/api/edge/apps/upload", "file", filePath, file)
			if err != nil {
				return err
			}

			_ = appID
			fmt.Fprintf(f.IO.Out, "Version uploaded\n")
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().String("app", "", "Edge app ID (required)")
	_ = cmd.MarkFlagRequired("app")

	return cmd
}
