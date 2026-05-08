package firmware

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdUpload(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload <file-path>",
		Short: "Upload a firmware file",
		Long: `Upload a firmware binary file to the platform. Returns a file ID (fid) that
is needed when creating a firmware record with 'firmware create --fid <fid>'.`,
		Example: `  # Upload firmware and note the returned file ID
  devicemanager firmware upload ./IG502-V2.0.0.rl4207.bin`,
		Args: cobra.ExactArgs(1),
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

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.UploadWithFields("/api/file/form", "file", filePath, file, map[string]string{
				"filename": filepath.Base(filePath),
			})
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "Firmware file uploaded\n")
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	return cmd
}
