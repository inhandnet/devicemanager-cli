package device

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
	"github.com/inhandnet/devicemanager-cli/internal/ui"
)

func NewCmdKick(f *factory.Factory) *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "kick <device-id>",
		Short: "Force disconnect a device",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !ui.Confirm(f.IO, fmt.Sprintf("Kick device %s?", args[0]), yes) {
				return nil
			}

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

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip confirmation prompt")

	return cmd
}

func NewCmdReboot(f *factory.Factory) *cobra.Command {
	var (
		timeout int
		yes     bool
	)

	cmd := &cobra.Command{
		Use:   "reboot <device-id>",
		Short: "Reboot a device",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !ui.Confirm(f.IO, fmt.Sprintf("Reboot device %s?", args[0]), yes) {
				return nil
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{
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
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip confirmation prompt")

	return cmd
}
