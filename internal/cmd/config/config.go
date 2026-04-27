package config

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/inhandnet/elements-cli/internal/factory"
	"github.com/inhandnet/elements-cli/internal/iostreams"
)

func NewCmdConfig(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage CLI configuration",
	}

	cmd.AddCommand(newCmdCurrentContext(f))
	cmd.AddCommand(newCmdListContexts(f))
	cmd.AddCommand(newCmdUseContext(f))
	cmd.AddCommand(newCmdDeleteContext(f))

	return cmd
}

func newCmdCurrentContext(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "current-context",
		Short: "Display the current context",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := f.Config()
			if err != nil {
				return err
			}
			name := cfg.ActiveContextName()
			if name == "" {
				fmt.Fprintln(f.IO.Out, "No current context set")
				return nil
			}
			fmt.Fprintln(f.IO.Out, name)
			return nil
		},
	}
}

func newCmdListContexts(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "list-contexts",
		Aliases: []string{"ls"},
		Short:   "List all contexts",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := f.Config()
			if err != nil {
				return err
			}

			names := make([]string, 0, len(cfg.Contexts))
			for name := range cfg.Contexts {
				names = append(names, name)
			}
			sort.Strings(names)

			tp := iostreams.NewTablePrinter(f.IO.Out, f.IO.IsStdoutTTY())
			tp.AddRow("CURRENT", "NAME", "API", "USER")
			for _, name := range names {
				ctx := cfg.Contexts[name]
				current := ""
				if name == cfg.ActiveContextName() {
					current = "*"
				}
				tp.AddRow(current, name, ctx.APIURL(), ctx.User)
			}
			return tp.Render()
		},
	}
}

func newCmdUseContext(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "use-context <name>",
		Short: "Switch the active context",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := f.Config()
			if err != nil {
				return err
			}
			name := args[0]
			if _, ok := cfg.Contexts[name]; !ok {
				return fmt.Errorf("context %q not found", name)
			}
			cfg.CurrentContext = name
			if err := f.SaveConfig(); err != nil {
				return err
			}
			fmt.Fprintf(f.IO.Out, "%s Switched to context %q\n", iostreams.Green("✓"), name)
			return nil
		},
	}
}

func newCmdDeleteContext(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "delete-context <name>",
		Short: "Delete a context",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := f.Config()
			if err != nil {
				return err
			}
			name := args[0]
			if _, ok := cfg.Contexts[name]; !ok {
				return fmt.Errorf("context %q not found", name)
			}
			cfg.DeleteContext(name)
			if err := f.SaveConfig(); err != nil {
				return err
			}
			fmt.Fprintf(f.IO.Out, "%s Deleted context %q\n", iostreams.Green("✓"), name)
			return nil
		},
	}
}
