package firmware

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/elements-cli/internal/factory"
	"github.com/inhandnet/elements-cli/internal/iostreams"
)

type CreateOptions struct {
	FID        string
	Name       string
	Version    string
	Model      string
	Desc       string
	JobTimeout int
}

func NewCmdCreate(f *factory.Factory) *cobra.Command {
	opts := &CreateOptions{}

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create a firmware entry",
		Example: `  elements firmware create --fid <file-id> --name "IR615-v2.0" --version 2.0.0 --model IR615`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]interface{}{
				"fid":     opts.FID,
				"name":    opts.Name,
				"version": opts.Version,
				"model":   opts.Model,
			}
			if opts.Desc != "" {
				body["desc"] = opts.Desc
			}
			if opts.JobTimeout > 0 {
				body["jobTimeout"] = opts.JobTimeout
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Post("/api/firmwares", body)
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "Firmware created\n")
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().StringVar(&opts.FID, "fid", "", "Uploaded file ID (required)")
	cmd.Flags().StringVar(&opts.Name, "name", "", "Firmware name (required)")
	cmd.Flags().StringVar(&opts.Version, "version", "", "Firmware version (required)")
	cmd.Flags().StringVar(&opts.Model, "model", "", "Device model (required)")
	cmd.Flags().StringVar(&opts.Desc, "desc", "", "Description")
	cmd.Flags().IntVar(&opts.JobTimeout, "job-timeout", 0, "Job timeout (seconds)")
	_ = cmd.MarkFlagRequired("fid")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("version")
	_ = cmd.MarkFlagRequired("model")

	return cmd
}
