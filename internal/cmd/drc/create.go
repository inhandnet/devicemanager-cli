package drc

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/elements-cli/internal/factory"
	"github.com/inhandnet/elements-cli/internal/iostreams"
)

type CreateOptions struct {
	Name        string
	Model       string
	Content     string
	ContentType string
	Desc        string
}

func NewCmdCreate(f *factory.Factory) *cobra.Command {
	opts := &CreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a DRC template",
		Example: `  elements drc create --name "IR615-default" --model IR615 --content "..."`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]interface{}{
				"name":    opts.Name,
				"model":   opts.Model,
				"content": opts.Content,
			}
			if opts.ContentType != "" {
				body["contentType"] = opts.ContentType
			}
			if opts.Desc != "" {
				body["desc"] = opts.Desc
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Post("/api/drc", body)
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "DRC template created\n")
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "Template name (required)")
	cmd.Flags().StringVar(&opts.Model, "model", "", "Device model (required)")
	cmd.Flags().StringVar(&opts.Content, "content", "", "Template content (required)")
	cmd.Flags().StringVar(&opts.ContentType, "content-type", "", "Content type")
	cmd.Flags().StringVar(&opts.Desc, "desc", "", "Description")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("model")
	_ = cmd.MarkFlagRequired("content")

	return cmd
}
