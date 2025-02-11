package describers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/opengovern/og-describer-semgrep/discovery/pkg/models"
	"github.com/opengovern/og-describer-semgrep/discovery/provider"
	"net/http"
	"strconv"
	"sync"
)

type RequestBody struct {
	RepositoryID int `json:"repository_id"`
}

func ListScans(ctx context.Context, handler *provider.SemGrepAPIHandler, stream *models.StreamSender) ([]models.Resource, error) {
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
			projects, err := provider.ListProjects(ctx, handler, deployment.Slug)
			if err != nil {
				errorChan <- err // Send error to the error channel
			}
			for _, project := range projects {
				if err := processScans(ctx, handler, strconv.Itoa(deployment.ID), project.ID, semGrepChan, &wg); err != nil {
					errorChan <- err // Send error to the error channel
				}
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

func processScans(ctx context.Context, handler *provider.SemGrepAPIHandler, deploymentID string, repositoryID int, semGrepChan chan<- models.Resource, wg *sync.WaitGroup) error {
	var scanListResponse provider.ScansListResponse
	var resp *http.Response
	baseURL := "https://semgrep.dev/api/v1/deployments/"

	finalURL := fmt.Sprintf("%s%s/scans/search", baseURL, deploymentID)

	body := RequestBody{
		RepositoryID: repositoryID,
	}

	requestData, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshaling json: %w", err)
	}

	req, err := http.NewRequest("POST", finalURL, bytes.NewBuffer(requestData))
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

		if e = json.NewDecoder(resp.Body).Decode(&scanListResponse); e != nil {
			return nil, fmt.Errorf("failed to decode response: %w", e)
		}
		return resp, nil
	}

	err = handler.DoRequest(ctx, req, requestFunc)
	if err != nil {
		return fmt.Errorf("error during request handling: %w", err)
	}

	for _, scan := range scanListResponse.Scans {
		wg.Add(1)
		go func(scan provider.ScanJSON) {
			defer wg.Done()
			value := models.Resource{
				ID:   scan.ID,
				Name: scan.ID,
				Description: provider.ScanDescription{
					ID:             scan.ID,
					DeploymentID:   scan.DeploymentID,
					RepositoryID:   scan.RepositoryID,
					Branch:         scan.Branch,
					Commit:         scan.Commit,
					IsFullScan:     scan.IsFullScan,
					StartedAt:      scan.StartedAt,
					CompletedAt:    scan.CompletedAt,
					ExitCode:       scan.ExitCode,
					TotalTime:      scan.TotalTime,
					FindingsCounts: scan.FindingsCounts,
					Status:         scan.Status,
				},
			}
			semGrepChan <- value
		}(scan)
	}
	return nil
}
