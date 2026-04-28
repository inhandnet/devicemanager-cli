package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// OAuthClient holds the client_id and client_secret from the platform.
type OAuthClient struct {
	ClientID     string
	ClientSecret string
}

// FetchOAuthClient retrieves OAuth client_id and client_secret from the platform's
// public configuration API: GET /api/platform/config
func FetchOAuthClient(ctx context.Context, host string) (*OAuthClient, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, host+"/api/platform/config", http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching platform config: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("platform config HTTP %d: %s", resp.StatusCode, string(body))
	}

	var config struct {
		Result struct {
			Auth struct {
				ClientID     json.Number `json:"clientId"`
				ClientSecret string      `json:"clientSecret"`
			} `json:"auth"`
		} `json:"result"`
	}
	if err := json.Unmarshal(body, &config); err != nil {
		return nil, fmt.Errorf("parsing platform config: %w", err)
	}

	clientID := config.Result.Auth.ClientID.String()
	if clientID == "" {
		return nil, fmt.Errorf("clientId not found in platform config")
	}

	return &OAuthClient{
		ClientID:     clientID,
		ClientSecret: config.Result.Auth.ClientSecret,
	}, nil
}

// OAuthToken holds the token response from DM platform.
type OAuthToken struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"`
	ExpiresAt    time.Time `json:"-"`
}

// RefreshAccessToken uses the refresh_token to obtain a new access_token.
func RefreshAccessToken(host, clientID, clientSecret, refreshToken string) (*OAuthToken, error) {
	data := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
	}

	resp, err := http.Post(host+"/oauth2/access_token", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading refresh response: %w", err)
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		Error        string `json:"error"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("parsing refresh response: %w (%s)", err, string(body))
	}
	if tokenResp.Error != "" {
		return nil, fmt.Errorf("refresh failed: %s", tokenResp.Error)
	}
	if tokenResp.AccessToken == "" {
		return nil, fmt.Errorf("refresh returned empty access_token")
	}

	token := &OAuthToken{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		TokenType:    tokenResp.TokenType,
		ExpiresIn:    tokenResp.ExpiresIn,
	}
	if token.ExpiresIn > 0 {
		token.ExpiresAt = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
	}
	return token, nil
}
