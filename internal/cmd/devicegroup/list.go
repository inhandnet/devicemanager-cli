package devicegroup

import (
	"net/url"

	"github.com/spf13/cobra"

	"github.com/inhandnet/elements-cli/internal/cmdutil"
	"github.com/inhandnet/elements-cli/internal/factory"
	"github.com/inhandnet/elements-cli/internal/iostreams"
)

type ListOptions struct {
	cmdutil.ListFlags
	Parent string
}

func NewCmdList(f *factory.Factory) *cobra.Command {
	opts := &ListOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List device groups",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			opts.ApplyTo(q)
			cmdutil.SetQueryParam(q, "parent", opts.Parent)

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get("/api/devicegroups", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output,
				iostreams.WithColumns("_id", "name", "parent", "total"))
		},
	}

	opts.Register(cmd)
	cmd.Flags().StringVar(&opts.Parent, "parent", "", "Filter by parent group ID")

	return cmd
}
