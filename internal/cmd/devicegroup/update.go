package devicegroup

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdUpdate(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <group-id>",
		Short: "Update a device group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			body := map[string]interface{}{}
			if name != "" {
				body["name"] = name
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Put(fmt.Sprintf("/api/devicegroups/%s", args[0]), body)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().String("name", "", "New group name")

	return cmd
}
