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

func ListPolicies(ctx context.Context, handler *provider.SemGrepAPIHandler, stream *models.StreamSender) ([]models.Resource, error) {
	var wg sync.WaitGroup
	semGrepChan := make(chan models.Resource)
	errorChan := make(chan error, 1) // Buffered channel to capture errors
	deployments, err := provider.ListDeployments(ctx, handler)
	if err != nil {
		return nil, err
	}

	go func() {
		defer close(semGrepChan)
		defer close(errorChan)
		for _, deployment := range deployments {
			if err := processPolicies(ctx, handler, strconv.Itoa(deployment.ID), semGrepChan, &wg); err != nil {
				errorChan <- err // Send error to the error channel
			}
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

func processPolicies(ctx context.Context, handler *provider.SemGrepAPIHandler, deploymentID string, semGrepChan chan<- models.Resource, wg *sync.WaitGroup) error {
	var policyListResponse provider.PoliciesListResponse
	var resp *http.Response
	baseURL := "https://semgrep.dev/api/v1/deployments/"

	finalURL := fmt.Sprintf("%s%s/policies", baseURL, deploymentID)

	req, err := http.NewRequest("GET", finalURL, nil)
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

		if e = json.NewDecoder(resp.Body).Decode(&policyListResponse); e != nil {
			return nil, fmt.Errorf("failed to decode response: %w", e)
		}
		return resp, nil
	}

	err = handler.DoRequest(ctx, req, requestFunc)
	if err != nil {
		return fmt.Errorf("error during request handling: %w", err)
	}

	for _, policy := range policyListResponse.Policies {
		wg.Add(1)
		go func(policy provider.PolicyJSON) {
			defer wg.Done()
			value := models.Resource{
				ID:   policy.ID,
				Name: policy.Name,
				Description: provider.PolicyDescription{
					ID:          policy.ID,
					Name:        policy.Name,
					Slug:        policy.Slug,
					ProductType: policy.ProductType,
					IsDefault:   policy.IsDefault,
				},
			}
			semGrepChan <- value
		}(policy)
	}
	return nil
}
