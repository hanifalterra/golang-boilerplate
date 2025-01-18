package cacabot

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"golang-boilerplate/internal/pkg/logger"
)

type Client struct {
	URL      string
	Username string
	Password string
	Enabled  bool
}

// NewCacabotClient creates a new Cacabot client.
func NewCacabotClient(url, username, password string, enabled bool) *Client {
	return &Client{
		URL:      url,
		Username: username,
		Password: password,
		Enabled:  enabled,
	}
}

// SendMessage sends a message to Cacabot.
func (c *Client) SendMessage(ctx context.Context, path string, payload interface{}) error {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Exit if the client is not enabled.
	if !c.Enabled {
		logger.FromContext(ctx).Info(ctx, "connections.cacabot", "SendMessage", "DryRun: %s%s, payload: %s", c.URL, path, payloadJSON)
		return nil
	}

	url := fmt.Sprintf("%s%s", c.URL, path)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set basic auth and content type headers.
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.Username, c.Password)))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	return nil
}
