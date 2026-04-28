package tunnel

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inhandnet/elements-cli/internal/factory"
	"github.com/inhandnet/elements-cli/internal/iostreams"
)

type CreateOptions struct {
	Name         string
	DeviceID     string
	Proto        string
	LocalAddress string
	LocalPort    int
}

func NewCmdCreate(f *factory.Factory) *cobra.Command {
	opts := &CreateOptions{}

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create a tunnel",
		Example: `  elements tunnel create --name ssh-tunnel --device-id <id> --proto tcp --local-address 127.0.0.1 --local-port 22`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]interface{}{
				"name":         opts.Name,
				"deviceId":     opts.DeviceID,
				"proto":        opts.Proto,
				"localAddress": opts.LocalAddress,
				"localPort":    opts.LocalPort,
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Post("/api/touch/tunnels", body)
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "Tunnel created\n")
			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "Tunnel name (required)")
	cmd.Flags().StringVar(&opts.DeviceID, "device-id", "", "Device ID (required)")
	cmd.Flags().StringVar(&opts.Proto, "proto", "tcp", "Protocol (http/https/tcp)")
	cmd.Flags().StringVar(&opts.LocalAddress, "local-address", "127.0.0.1", "Local address")
	cmd.Flags().IntVar(&opts.LocalPort, "local-port", 0, "Local port (required)")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("device-id")
	_ = cmd.MarkFlagRequired("local-port")

	return cmd
}
