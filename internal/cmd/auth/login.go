package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/api"
	"github.com/inhandnet/devicemanager-cli/internal/config"
	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

const defaultCallbackPort = 18920

type LoginOptions struct {
	ContextName string
	Host        string
	Port        int
	Timeout     time.Duration
}

func NewCmdLogin(f *factory.Factory) *cobra.Command {
	opts := &LoginOptions{}

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login via browser",
		Example: `  # Login via browser (default) — opens DM login page
  devicemanager auth login

  # Login to global region
  devicemanager auth login --host global

  # Login with a custom domain
  devicemanager auth login --host iot.example.com`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBrowserLogin(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.ContextName, "context", "default", "Context name to create/update")
	cmd.Flags().StringVar(&opts.Host, "host", "global", `Platform region: "global", "cn", or a custom domain`)
	cmd.Flags().IntVar(&opts.Port, "port", defaultCallbackPort, "Local callback server port")
	cmd.Flags().DurationVar(&opts.Timeout, "timeout", 3*time.Minute, "Timeout waiting for browser login")

	return cmd
}

var regionHosts = map[string]string{
	"cn":     "iot.inhand.com.cn",
	"global": "iot.inhandnetworks.com",
}

func resolveHost(host string) (string, error) {
	if host == "" {
		return "", fmt.Errorf("host is required")
	}
	if domain, ok := regionHosts[strings.ToLower(host)]; ok {
		return domain, nil
	}
	return host, nil
}

// runBrowserLogin opens the DM login page and waits for the OAuth2 callback.
func runBrowserLogin(f *factory.Factory, opts *LoginOptions) error {
	out := f.IO.Out

	host, err := resolveHost(opts.Host)
	if err != nil {
		return err
	}

	apiURL := "https://" + host

	fmt.Fprintln(out, "Fetching OAuth client configuration...")
	oauthClient, err := api.FetchOAuthClient(context.Background(), apiURL)
	if err != nil {
		return fmt.Errorf("fetching OAuth config from %s: %w", apiURL, err)
	}

	state := fmt.Sprintf("devicemanager-cli-%d", time.Now().UnixNano())
	redirectURI := fmt.Sprintf("http://localhost:%d/callback", opts.Port)

	loginURL := fmt.Sprintf("https://%s/user/login?redirect_uri=%s&state=%s",
		host, url.QueryEscape(redirectURI), state)

	fmt.Fprintln(out, "Opening browser for authentication...")
	fmt.Fprintln(out, iostreams.Gray("If the browser doesn't open, visit:"))
	fmt.Fprintln(out, iostreams.Gray(loginURL))
	fmt.Fprintln(out)

	openBrowser(loginURL)

	fmt.Fprintln(out, "Waiting for login...")

	result, err := api.WaitForCallback(opts.Port, opts.Timeout)
	if err != nil {
		return err
	}

	// Verify state to prevent CSRF
	if result.State != state {
		return fmt.Errorf("state mismatch: expected %s, got %s", state, result.State)
	}

	fmt.Fprintln(out, "Exchanging authorization code...")

	token, err := api.ExchangeCodeForToken("https://"+host, result.Code, oauthClient.ClientID, oauthClient.ClientSecret, redirectURI)
	if err != nil {
		return err
	}

	return saveLogin(f, opts, host, "", oauthClient.ClientID, oauthClient.ClientSecret, token)
}

func saveLogin(f *factory.Factory, opts *LoginOptions, host, username, clientID, clientSecret string, token *api.OAuthToken) error {
	cfg, err := f.Config()
	if err != nil {
		return err
	}

	ctx := &config.Context{
		Host:         host,
		Token:        token.AccessToken,
		RefreshToken: token.RefreshToken,
		User:         username,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
	if !token.ExpiresAt.IsZero() {
		ctx.ExpiresAt = token.ExpiresAt
	}

	// Fetch current user info from API
	if user, authority := fetchCurrentUser(ctx); user != "" {
		ctx.User = user
		ctx.Authority = authority
	}

	cfg.SetContext(opts.ContextName, ctx)
	cfg.CurrentContext = opts.ContextName

	if err := f.SaveConfig(); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	apiURL := "https://" + host
	fmt.Fprintf(f.IO.Out, "%s Logged in to %s (context: %s) as %s\n",
		iostreams.Green("✓"), apiURL, opts.ContextName, iostreams.Bold(ctx.User))
	return nil
}

// fetchCurrentUser calls /api/users/this to get the logged-in user's display name and authority.
func fetchCurrentUser(ctx *config.Context) (string, string) {
	transport := &api.TokenTransport{
		Token: ctx.Token,
		Base:  http.DefaultTransport,
	}
	client := api.NewAPIClient(ctx.APIURL(), transport, 100)
	q := url.Values{}
	body, err := client.Get("/api/users/this", q)
	if err != nil {
		return "", ""
	}
	var resp struct {
		Result struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Phone    string `json:"phone"`
			RoleName string `json:"roleName"`
			IsRoot   bool   `json:"isRoot"`
		} `json:"result"`
	}
	if json.Unmarshal(body, &resp) != nil {
		return "", ""
	}
	name := resp.Result.Name
	if name == "" {
		name = resp.Result.Email
	}
	if name == "" {
		name = resp.Result.Phone
	}
	if resp.Result.Email != "" {
		name = fmt.Sprintf("%s (%s)", name, resp.Result.Email)
	}

	authority := resp.Result.RoleName
	if resp.Result.IsRoot {
		authority = "root"
	}

	return name, authority
}

// openBrowser tries to open a URL in the default browser.
func openBrowser(targetURL string) {
	// Use platform-specific command
	cmd := browserCmd(targetURL)
	if cmd != nil {
		_ = cmd.Start()
	}
}
