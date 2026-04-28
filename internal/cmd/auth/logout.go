package auth

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/debug"
	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdLogout(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "logout [context]",
		Short: "Logout and invalidate tokens",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := f.Config()
			if err != nil {
				return err
			}

			name := cfg.ActiveContextName()
			if len(args) > 0 {
				name = args[0]
			}
			if name == "" {
				return fmt.Errorf("no context specified and no current context set")
			}

			ctx, ok := cfg.Contexts[name]
			if !ok {
				return fmt.Errorf("context %q not found", name)
			}

			// Call server-side logout to invalidate the token
			if ctx.Token != "" {
				client, err := f.APIClient()
				if err == nil {
					q := url.Values{}
					q.Set("access_token", ctx.Token)
					_, logoutErr := client.Get("/api2/logout", q)
					if logoutErr != nil {
						debug.Log("server logout failed: %v", logoutErr)
					} else {
						debug.Log("server logout succeeded")
					}
				}
			}

			// Clear local tokens
			ctx.Token = ""
			ctx.RefreshToken = ""

			if err := f.SaveConfig(); err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "%s Logged out of context %q\n", iostreams.Green("✓"), name)
			return nil
		},
	}
}
