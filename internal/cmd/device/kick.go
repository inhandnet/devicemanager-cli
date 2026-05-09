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
		Use:     "kick <device-id>",
		Short:   "Force disconnect a device",
		Long:    "Force disconnect a device from the platform. The device will attempt to reconnect automatically.",
		Example: `  devicemanager device kick 5d6349d6335c8c000178a194`,
		Args:    cobra.ExactArgs(1),
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
		Short: "Reboot a device remotely",
		Long:  "Send a reboot command to a device. The device must be online. It will go offline briefly and reconnect.",
		Example: `  devicemanager device reboot 5d6349d6335c8c000178a194
  devicemanager device reboot 5d6349d6335c8c000178a194 --timeout 30000  # 30 second timeout`,
		Args: cobra.ExactArgs(1),
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

	cmd.Flags().IntVar(&timeout, "timeout", 15000, "Timeout in milliseconds (default 15000ms = 15s)")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip confirmation prompt")

	return cmd
}
