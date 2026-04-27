package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/inhandnet/elements-cli/internal/build"
	"github.com/inhandnet/elements-cli/internal/debug"
	"github.com/inhandnet/elements-cli/internal/factory"
)

func NewCmdRoot(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "elements",
		Short:         "Device Manager Platform CLI",
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       build.Version,
	}

	cmd.PersistentFlags().StringP("output", "o", "", "Output format: json, table, yaml (default: table for TTY, json otherwise)")
	cmd.PersistentFlags().String("jq", "", `Filter JSON output using a jq expression (implies -o json)`)
	cmd.PersistentFlags().String("context", "", "Override active context (env: ELEMENTS_CONTEXT)")
	cmd.PersistentFlags().Bool("debug", false, "Enable debug output (env: ELEMENTS_DEBUG)")

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if d, _ := cmd.Flags().GetBool("debug"); d {
			debug.Enabled = true
		} else if os.Getenv("ELEMENTS_DEBUG") != "" {
			debug.Enabled = true
		}

		outputExplicit := cmd.Flags().Changed("output")
		if !outputExplicit {
			if f.IO.IsStdoutTTY() {
				_ = cmd.Flags().Set("output", "table")
			} else {
				_ = cmd.Flags().Set("output", "json")
			}
		}

		if jqExpr, _ := cmd.Flags().GetString("jq"); jqExpr != "" {
			f.IO.JQExpr = jqExpr
			if !outputExplicit {
				_ = cmd.Flags().Set("output", "json")
			}
		}

		if ctx, _ := cmd.Flags().GetString("context"); ctx != "" {
			if err := os.Setenv("ELEMENTS_CONTEXT", ctx); err != nil {
				return err
			}
		}
		return nil
	}

	return cmd
}
