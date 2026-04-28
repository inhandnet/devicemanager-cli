package devicegroup

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/elements-cli/internal/factory"
	"github.com/inhandnet/elements-cli/internal/iostreams"
)

func NewCmdGet(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <group-id>",
		Short: "Get device group details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get(fmt.Sprintf("/api/devicegroups/%s", args[0]), nil)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output)
		},
	}

	return cmd
}
