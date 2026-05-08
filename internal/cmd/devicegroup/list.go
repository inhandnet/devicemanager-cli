package devicegroup

import (
	"net/url"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/cmdutil"
	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdList(f *factory.Factory) *cobra.Command {
	var (
		parent   string
		verbose  int
		maxDepth int
	)

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
			if verbose > 0 {
				q.Set("verbose", strconv.Itoa(verbose))
			}
			if maxDepth > 0 {
				q.Set("max_depth", strconv.Itoa(maxDepth))
			}
			cmdutil.SetQueryParam(q, "parent", parent)

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get("/api/devicegroups", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output,
				iostreams.WithColumns("_id", "name", "parent", "total"))
		},
	}

	cmd.Flags().StringVar(&parent, "parent", "", "Filter by parent group ID")
	cmd.Flags().IntVar(&verbose, "verbose", 100, "Detail level (1-100)")
	cmd.Flags().IntVar(&maxDepth, "max-depth", 3, "Max depth of group tree")

	return cmd
}
