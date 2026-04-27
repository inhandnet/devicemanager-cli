package device

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/inhandnet/elements-cli/internal/factory"
	"github.com/inhandnet/elements-cli/internal/iostreams"
)

func NewCmdKick(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "kick <device-id>",
		Short: "Force disconnect a device",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			_, err = client.Get("/api/device/"+args[0]+"/kick", url.Values{})
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "%s Device %s kicked\n", iostreams.Green("✓"), args[0])
			return nil
		},
	}
}

func NewCmdReboot(f *factory.Factory) *cobra.Command {
	var timeout int

	cmd := &cobra.Command{
		Use:   "reboot <device-id>",
		Short: "Reboot a device",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]interface{}{
				"method":  "reboot",
				"timeout": timeout,
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Post("/api/device/"+args[0]+"/methods", body)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().IntVar(&timeout, "timeout", 15000, "Timeout in milliseconds")

	return cmd
}
