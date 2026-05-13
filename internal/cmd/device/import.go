package device

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"

	"github.com/inhandnet/devicemanager-cli/internal/cmdutil"
	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdImport(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import <file-path>",
		Short: "Batch import devices from an Excel file",
		Example: `  # Import devices from Excel
  devicemanager device import devices.xlsx

  # Import and assign to a device group
  devicemanager device import devices.xlsx --group <group-id>

  # Import without overwriting existing devices
  devicemanager device import devices.xlsx --no-overwrite`,
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

			// Step 1: upload file
			fmt.Fprintf(f.IO.Out, "Uploading %s...\n", filepath.Base(filePath))
			q := url.Values{}
			q.Set("access_token", f.Token())
			cmdutil.SetQueryParam(q, "oid", f.OrgID())
			uploadURL := fmt.Sprintf("/api/file/form?%s", q.Encode())
			uploadResp, err := client.UploadWithFields(uploadURL, "file", filepath.Base(filePath), file, map[string]string{
				"filename": filepath.Base(filePath),
			})
			if err != nil {
				return fmt.Errorf("uploading file: %w", err)
			}

			fileID := gjson.GetBytes(uploadResp, "result._id").String()
			if fileID == "" {
				return fmt.Errorf("failed to get file ID from upload response")
			}

			// Step 2: batch add
			groupID, _ := cmd.Flags().GetString("group")
			noOverwrite, _ := cmd.Flags().GetBool("no-overwrite")

			body := map[string]any{
				"toBeCovered": !noOverwrite,
			}
			if groupID != "" {
				body["deviceGroupId"] = groupID
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Post(fmt.Sprintf("/api/v2/devices/batch_add?file_id=%s", fileID), body)
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "Devices imported from %s\n", filepath.Base(filePath))
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().String("group", "", "Device group ID to assign imported devices to")
	cmd.Flags().Bool("no-overwrite", false, "Do not overwrite existing devices")

	return cmd
}
