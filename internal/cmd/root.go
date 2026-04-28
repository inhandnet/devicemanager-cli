package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/build"
	"github.com/inhandnet/devicemanager-cli/internal/debug"
	"github.com/inhandnet/devicemanager-cli/internal/factory"
)

func NewCmdRoot(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "devicemanager",
		Short:         "Device Manager Platform CLI",
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       build.Version,
	}

	cmd.PersistentFlags().StringP("output", "o", "json", "Output format: json, table, yaml")
	cmd.PersistentFlags().String("jq", "", `Filter JSON output using a jq expression (implies -o json)`)
	cmd.PersistentFlags().String("context", "", "Override active context (env: DEVICEMANAGER_CONTEXT)")
	cmd.PersistentFlags().Bool("debug", false, "Enable debug output (env: DEVICEMANAGER_DEBUG)")

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if d, _ := cmd.Flags().GetBool("debug"); d {
			debug.Enabled = true
		} else if os.Getenv("DEVICEMANAGER_DEBUG") != "" {
			debug.Enabled = true
		}

		if jqExpr, _ := cmd.Flags().GetString("jq"); jqExpr != "" {
			f.IO.JQExpr = jqExpr
			_ = cmd.Flags().Set("output", "json")
		}

		if ctx, _ := cmd.Flags().GetString("context"); ctx != "" {
			if err := os.Setenv("DEVICEMANAGER_CONTEXT", ctx); err != nil {
				return err
			}
		}
		return nil
	}

	return cmd
}
