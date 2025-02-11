package describers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/opengovern/og-describer-semgrep/discovery/pkg/models"
	"github.com/opengovern/og-describer-semgrep/discovery/provider"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

func ListProjects(ctx context.Context, handler *provider.SemGrepAPIHandler, stream *models.StreamSender) ([]models.Resource, error) {
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
			if err := processProjects(ctx, handler, deployment.Slug, semGrepChan, &wg); err != nil {
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

func processProjects(ctx context.Context, handler *provider.SemGrepAPIHandler, deploymentSlug string, semGrepChan chan<- models.Resource, wg *sync.WaitGroup) error {
	var projects []provider.ProjectJSON
	var projectListResponse provider.ProjectsListResponse
	var resp *http.Response
	baseURL := "https://semgrep.dev/api/v1/deployments/"
	page := 0

	for {
		params := url.Values{}
		params.Set("page", strconv.Itoa(page))
		params.Set("page_size", "3000")
		finalURL := fmt.Sprintf("%s%s/projects?%s", baseURL, deploymentSlug, params.Encode())

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

			if e = json.NewDecoder(resp.Body).Decode(&projectListResponse); e != nil {
				return nil, fmt.Errorf("failed to decode response: %w", e)
			}
			projects = append(projects, projectListResponse.Projects...)
			return resp, nil
		}

		err = handler.DoRequest(ctx, req, requestFunc)
		if err != nil {
			return fmt.Errorf("error during request handling: %w", err)
		}

		if len(projectListResponse.Projects) < 3000 {
			break
		}
		page++
	}

	for _, project := range projects {
		wg.Add(1)
		go func(project provider.ProjectJSON) {
			defer wg.Done()
			value := models.Resource{
				ID:   strconv.Itoa(project.ID),
				Name: project.Name,
				Description: provider.ProjectDescription{
					ID:            project.ID,
					Name:          project.Name,
					URL:           project.URL,
					Tags:          project.Tags,
					CreatedAt:     project.CreatedAt,
					LatestScanAt:  project.LatestScanAt,
					PrimaryBranch: project.PrimaryBranch,
					DefaultBranch: project.DefaultBranch,
				},
			}
			semGrepChan <- value
		}(project)
	}
	return nil
}
