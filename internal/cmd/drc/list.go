package drc

import (
	"net/url"

	"github.com/spf13/cobra"

	"github.com/inhandnet/elements-cli/internal/cmdutil"
	"github.com/inhandnet/elements-cli/internal/factory"
	"github.com/inhandnet/elements-cli/internal/iostreams"
)

type ListOptions struct {
	cmdutil.ListFlags
	Name  string
	Model string
}

func NewCmdList(f *factory.Factory) *cobra.Command {
	opts := &ListOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List DRC templates",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			opts.ListFlags.ApplyTo(q)
			cmdutil.SetQueryParam(q, "name", opts.Name)
			cmdutil.SetQueryParam(q, "model", opts.Model)

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get("/api/drc", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output,
				iostreams.WithColumns("_id", "name", "model", "contentType"))
		},
	}

	opts.Register(cmd)
	cmd.Flags().StringVar(&opts.Name, "name", "", "Filter by template name")
	cmd.Flags().StringVar(&opts.Model, "model", "", "Filter by device model")

	return cmd
}
