package tunnel

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdUpdate(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <tunnel-id>",
		Short: "Update a tunnel",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]interface{}{}
			if name, _ := cmd.Flags().GetString("name"); name != "" {
				body["name"] = name
			}
			if proto, _ := cmd.Flags().GetString("proto"); proto != "" {
				body["proto"] = proto
			}
			if addr, _ := cmd.Flags().GetString("local-address"); addr != "" {
				body["localAddress"] = addr
			}
			if port, _ := cmd.Flags().GetInt("local-port"); port > 0 {
				body["localPort"] = port
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Put(fmt.Sprintf("/api/touch/tunnels/%s", args[0]), body)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().String("name", "", "Tunnel name")
	cmd.Flags().String("proto", "", "Protocol (http/https/tcp)")
	cmd.Flags().String("local-address", "", "Local address")
	cmd.Flags().Int("local-port", 0, "Local port")

	return cmd
}
