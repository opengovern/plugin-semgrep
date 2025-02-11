package main

import (
	"encoding/json"
	"fmt"
	"github.com/opengovern/og-describer-semgrep/discovery/provider"
	"net/http"
)

// Config represents the JSON input configuration
type Config struct {
	Token string `json:"token"`
}

func IntegrationHealthcheck(cfg Config) (bool, error) {
	var deploymentListResponse provider.DeploymentsResponse
	var resp *http.Response
	client := http.DefaultClient
	baseURL := "https://semgrep.dev/api/v1/deployments"

	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.Token))

	resp, err = client.Do(req)
	if err != nil {
		return false, fmt.Errorf("request execution failed: %w", err)
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&deploymentListResponse); err != nil {
		return false, fmt.Errorf("failed to decode response: %w", err)
	}

	if err != nil {
		return false, fmt.Errorf("error during request handling: %w", err)
	}

	return true, nil
}
