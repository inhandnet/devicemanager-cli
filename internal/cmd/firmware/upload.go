package firmware

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdUpload(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload <file-path>",
		Short: "Upload a firmware file",
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

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Upload("/api/file/form", "file", filePath, file)
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "Firmware file uploaded\n")
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	return cmd
}
