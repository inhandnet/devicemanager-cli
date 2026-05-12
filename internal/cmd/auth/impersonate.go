package auth

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"

	"github.com/inhandnet/devicemanager-cli/internal/api"
	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdImpersonate(f *factory.Factory) *cobra.Command {
	var (
		userID string
		orgID  string
		stop   bool
	)

	cmd := &cobra.Command{
		Use:   "impersonate",
		Short: "Impersonate another user (requires ROOT privilege)",
		Example: `  # Impersonate by org ID (auto-resolves org admin)
  devicemanager auth impersonate --org 5e0956c46aa6d10001e931e6

  # Impersonate a specific user in an org
  devicemanager auth impersonate --org <oid> --user <uid>

  # Stop impersonation and restore admin identity
  devicemanager auth impersonate --stop`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := f.Config()
			if err != nil {
				return err
			}
			name := cfg.ActiveContextName()
			ctx, ok := cfg.Contexts[name]
			if !ok {
				return fmt.Errorf("no active context; run 'devicemanager auth login' first")
			}

			// --stop: restore admin token
			if stop {
				if !ctx.IsImpersonating() {
					return fmt.Errorf("not currently impersonating")
				}
				ctx.Token = ctx.AdminToken
				ctx.RefreshToken = ctx.AdminRefreshToken
				ctx.ExpiresAt = ctx.AdminExpiresAt
				ctx.AdminToken = ""
				ctx.AdminRefreshToken = ""
				ctx.AdminExpiresAt = time.Time{}
				if err := f.SaveConfig(); err != nil {
					return err
				}
				fmt.Fprintf(f.IO.Out, "%s Impersonation stopped, admin identity restored\n", iostreams.Green("✓"))
				return nil
			}

			if orgID == "" {
				return fmt.Errorf("--org is required (or use --stop)")
			}

			if userID != "" && orgID == "" {
				return fmt.Errorf("--user requires --org")
			}

			if ctx.IsImpersonating() {
				return fmt.Errorf("already impersonating; run 'devicemanager auth impersonate --stop' first")
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			// Resolve user for org (when only --org is given)
			if userID == "" {
				resolved, err := resolveUserForOrg(client, orgID)
				if err != nil {
					return fmt.Errorf("resolving user for org: %w", err)
				}
				userID = resolved
			}

			// Ensure token is fresh before impersonation.
			// access_token is passed as a query param which TokenTransport cannot update on 401 retry.
			if !ctx.ExpiresAt.IsZero() && time.Now().After(ctx.ExpiresAt) && ctx.RefreshToken != "" {
				newToken, refreshErr := api.RefreshAccessToken(
					ctx.APIURL(), ctx.ClientID, ctx.ClientSecret, ctx.RefreshToken)
				if refreshErr != nil {
					return fmt.Errorf("token expired and refresh failed: %w\nHint: run 'devicemanager auth login' to re-authenticate", refreshErr)
				}
				ctx.Token = newToken.AccessToken
				if newToken.RefreshToken != "" {
					ctx.RefreshToken = newToken.RefreshToken
				}
				ctx.ExpiresAt = newToken.ExpiresAt
				_ = f.SaveConfig()
			}

			// Call impersonate API
			token := ctx.EffectiveToken()
			q := url.Values{}
			q.Set("oid", orgID)
			q.Set("uid", userID)
			q.Set("access_token", token)
			q.Set("verbose", "100")

			body, err := client.Post("/api/token/impersonate?"+q.Encode(), nil)
			if err != nil {
				return fmt.Errorf("impersonate failed: %w", err)
			}

			var tokenResp struct {
				AccessToken  string `json:"access_token"`
				RefreshToken string `json:"refresh_token"`
				ExpiresIn    int64  `json:"expires_in"`
				Error        string `json:"error"`
				ErrorCode    int    `json:"error_code"`
			}
			if err := json.Unmarshal(body, &tokenResp); err != nil {
				return fmt.Errorf("parsing token response: %w", err)
			}
			if tokenResp.Error != "" {
				return fmt.Errorf("impersonate failed: %s (code: %d)", tokenResp.Error, tokenResp.ErrorCode)
			}
			if tokenResp.AccessToken == "" {
				return fmt.Errorf("impersonate failed: no access_token in response")
			}

			// Backup admin token
			ctx.AdminToken = ctx.Token
			ctx.AdminRefreshToken = ctx.RefreshToken
			ctx.AdminExpiresAt = ctx.ExpiresAt

			// Replace with impersonated token
			ctx.Token = tokenResp.AccessToken
			ctx.RefreshToken = tokenResp.RefreshToken

			if err := f.SaveConfig(); err != nil {
				return err
			}

			fmt.Fprintf(f.IO.Out, "%s Impersonating user %s in org %s\n", iostreams.Green("✓"), userID, orgID)
			fmt.Fprintf(f.IO.Out, "Run 'devicemanager auth impersonate --stop' to restore admin identity\n")
			return nil
		},
	}

	cmd.Flags().StringVar(&userID, "user", "", "User ID to impersonate")
	cmd.Flags().StringVar(&orgID, "org", "", "Organization ID")
	cmd.Flags().BoolVar(&stop, "stop", false, "Stop impersonation and restore admin identity")

	return cmd
}

// resolveUserForOrg finds the admin user for an org.
func resolveUserForOrg(client *api.APIClient, oid string) (string, error) {
	q := url.Values{}
	q.Set("oid", oid)
	q.Set("limit", "0")

	body, err := client.Get("/api2/users", q)
	if err != nil {
		return "", err
	}

	results := gjson.GetBytes(body, "result")
	if !results.Exists() {
		return "", fmt.Errorf("no users found for org %s", oid)
	}

	// Find admin user
	for _, user := range results.Array() {
		if user.Get("roleName").String() == "admin" {
			uid := user.Get("_id").String()
			if uid != "" {
				return uid, nil
			}
		}
	}

	// Fallback: first user
	uid := results.Array()[0].Get("_id").String()
	if uid == "" {
		return "", fmt.Errorf("no user ID found for org %s", oid)
	}
	return uid, nil
}
