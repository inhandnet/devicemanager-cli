package devicegroup

import (
	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

type CreateOptions struct {
	Name   string
	Parent string
}

func NewCmdCreate(f *factory.Factory) *cobra.Command {
	opts := &CreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a device group",
		Example: `  devicemanager devicegroup create --name "Factory A"
  devicemanager devicegroup create --name "Line 1" --parent <group-id>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]interface{}{
				"name": opts.Name,
			}
			if opts.Parent != "" {
				body["parent"] = opts.Parent
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Post("/api/devicegroups", body)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "Group name (required)")
	cmd.Flags().StringVar(&opts.Parent, "parent", "", "Parent group ID")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}
