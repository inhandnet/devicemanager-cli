package config

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
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
	cmd.AddCommand(newCmdConfigSet(f))
	cmd.AddCommand(newCmdConfigGet(f))

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

var validConfigKeys = []string{"ngrok-server"}

func isValidConfigKey(key string) bool {
	for _, k := range validConfigKeys {
		if k == key {
			return true
		}
	}
	return false
}

func newCmdConfigSet(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a global configuration value",
		Example: `  # Set custom ngrok server
  devicemanager config set ngrok-server my-ngrok.example.com:4443`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := f.Config()
			if err != nil {
				return err
			}
			key := args[0]
			value := args[1]
			if !isValidConfigKey(key) {
				return fmt.Errorf("unknown config key %q; valid keys: %v", key, validConfigKeys)
			}
			if key == "ngrok-server" {
				cfg.NgrokServer = value
			}
			if err := f.SaveConfig(); err != nil {
				return err
			}
			fmt.Fprintf(f.IO.Out, "%s Set %s = %s\n", iostreams.Green("✓"), key, value)
			return nil
		},
	}
}

func newCmdConfigGet(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get [key]",
		Short: "Get configuration values",
		Example: `  # Get all config values
  devicemanager config get

  # Get specific value
  devicemanager config get ngrok-server`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := f.Config()
			if err != nil {
				return err
			}

			if len(args) == 0 {
				// Show all
				tp := iostreams.NewTablePrinter(f.IO.Out, f.IO.IsStdoutTTY())
				tp.AddRow("KEY", "VALUE")
				for _, key := range validConfigKeys {
					var value string
					if key == "ngrok-server" {
						value = cfg.NgrokServer
					}
					if value == "" {
						value = "(not set)"
					}
					tp.AddRow(key, value)
				}
				return tp.Render()
			}

			key := args[0]
			if !isValidConfigKey(key) {
				return fmt.Errorf("unknown config key %q; valid keys: %v", key, validConfigKeys)
			}
			var value string
			if key == "ngrok-server" {
				value = cfg.NgrokServer
			}
			if value == "" {
				value = "(not set)"
			}
			fmt.Fprintln(f.IO.Out, value)
			return nil
		},
	}
}
