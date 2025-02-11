package describers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/opengovern/og-describer-semgrep/discovery/pkg/models"
	"github.com/opengovern/og-describer-semgrep/discovery/provider"
	"net/http"
	"strconv"
	"sync"
)

func ListDeployments(ctx context.Context, handler *provider.SemGrepAPIHandler, stream *models.StreamSender) ([]models.Resource, error) {
	var wg sync.WaitGroup
	semGrepChan := make(chan models.Resource)
	errorChan := make(chan error, 1) // Buffered channel to capture errors

	go func() {
		defer close(semGrepChan)
		defer close(errorChan)
		if err := processDeployments(ctx, handler, semGrepChan, &wg); err != nil {
			errorChan <- err // Send error to the error channel
		}
		wg.Wait()
	}()

	var values []models.Resource
	for {
		select {
		case value, ok := <-semGrepChan:
			if !ok {
				return values, nil
			}
			if stream != nil {
				if err := (*stream)(value); err != nil {
					return nil, err
				}
			} else {
				values = append(values, value)
			}
		case err := <-errorChan:
			return nil, err
		}
	}
}

func processDeployments(ctx context.Context, handler *provider.SemGrepAPIHandler, semGrepChan chan<- models.Resource, wg *sync.WaitGroup) error {
	var deploymentListResponse provider.DeploymentsResponse
	var resp *http.Response
	baseURL := "https://semgrep.dev/api/v1/deployments"

	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	requestFunc := func(req *http.Request) (*http.Response, error) {
		var e error
		resp, e = handler.Client.Do(req)
		if e != nil {
			return nil, fmt.Errorf("request execution failed: %w", e)
		}
		defer resp.Body.Close()

		if e = json.NewDecoder(resp.Body).Decode(&deploymentListResponse); e != nil {
			return nil, fmt.Errorf("failed to decode response: %w", e)
		}
		return resp, nil
	}

	err = handler.DoRequest(ctx, req, requestFunc)
	if err != nil {
		return fmt.Errorf("error during request handling: %w", err)
	}

	for _, deployment := range deploymentListResponse.Deployments {
		wg.Add(1)
		go func(deployment provider.DeploymentJSON) {
			defer wg.Done()
			findings := provider.Finding{
				URL: deployment.Findings.URL,
			}
			value := models.Resource{
				ID:   strconv.Itoa(deployment.ID),
				Name: deployment.Name,
				Description: provider.DeploymentDescription{
					Slug:     deployment.Slug,
					ID:       deployment.ID,
					Name:     deployment.Name,
					Findings: findings,
				},
			}
			semGrepChan <- value
		}(deployment)
	}
	return nil
}
