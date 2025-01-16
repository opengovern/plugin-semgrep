// artifact_dockerfile.go
package describers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/opengovern/og-describer-template/discovery/pkg/models"
	model "github.com/opengovern/og-describer-template/discovery/provider"
	resilientbridge "github.com/opengovern/resilient-bridge"
	"github.com/opengovern/resilient-bridge/adapters"
	"github.com/opengovern/resilient-bridge/utils" // For ExtractExternalBaseImagesFromBase64
)

// MAX_RESULTS is the maximum number of Dockerfiles to collect or stream.
const MAX_RESULTS = 500

// MAX_DOCKERFILE_LEN is the maximum allowed number of lines in a Dockerfile.
const MAX_DOCKERFILE_LEN = 500

// ListArtifactDockerFiles performs a single code search across the organization
// for "filename:Dockerfile" and processes each result. Each Dockerfile is
// streamed immediately upon processing, and also added to the final slice.
func ListArtifactDockerFiles(
	ctx context.Context,
	githubClient model.GitHubClient,
	organizationName string,
	stream *models.StreamSender,
) ([]models.Resource, error) {

	sdk := resilientbridge.NewResilientBridge()
	sdk.SetDebug(false)
	sdk.RegisterProvider("github", adapters.NewGitHubAdapter(githubClient.Token), &resilientbridge.ProviderConfig{
		UseProviderLimits: true,
		MaxRetries:        3,
		BaseBackoff:       time.Second,
	})

	// If org override is in context
	if orgVal := ctx.Value("organization"); orgVal != nil {
		if orgName, ok := orgVal.(string); ok && orgName != "" {
			organizationName = orgName
		}
	}

	// Build a single code search query that searches the entire org for Dockerfiles
	// Example: org:my-org filename:Dockerfile
	finalQuery := fmt.Sprintf("org:%s filename:Dockerfile", organizationName)

	var allValues []models.Resource
	totalCollected := 0
	perPage := 100
	page := 1

	for totalCollected < MAX_RESULTS {
		// Encode the search query
		q := url.QueryEscape(finalQuery)
		searchEndpoint := fmt.Sprintf("/search/code?q=%s&per_page=%d&page=%d", q, perPage, page)

		searchReq := &resilientbridge.NormalizedRequest{
			Method:   "GET",
			Endpoint: searchEndpoint,
			Headers:  map[string]string{"Accept": "application/vnd.github+json"},
		}
		searchResp, err := sdk.Request("github", searchReq)
		if err != nil {
			return allValues, fmt.Errorf("error performing code search in org %s: %w", organizationName, err)
		}
		if searchResp.StatusCode >= 400 {
			return allValues, fmt.Errorf("HTTP error %d searching Dockerfiles in org %s: %s",
				searchResp.StatusCode, organizationName, string(searchResp.Data))
		}

		var result model.CodeSearchResult
		if err := json.Unmarshal(searchResp.Data, &result); err != nil {
			return allValues, fmt.Errorf("error parsing code search response for org %s: %w", organizationName, err)
		}

		if len(result.Items) == 0 {
			// No more results
			break
		}

		// Process each Dockerfile found
		for _, item := range result.Items {
			resource, err := GetDockerfile(
				ctx,
				githubClient,
				organizationName,         // org name
				item.Repository.FullName, // e.g. "my-org/my-repo"
				item.Path,                // e.g. "path/to/Dockerfile"
				stream,
			)
			if err != nil {
				log.Printf("Skipping %s/%s: %v\n", item.Repository.FullName, item.Path, err)
				continue
			}
			if resource == nil {
				continue
			}

			// 1) Add to our local slice
			allValues = append(allValues, *resource)
			totalCollected++

			// 2) Stream the Dockerfile result immediately
			if stream != nil {
				if err := (*stream)(*resource); err != nil {
					return allValues, fmt.Errorf("error streaming resource: %w", err)
				}
			}

			if totalCollected >= MAX_RESULTS {
				break
			}
		}

		if len(result.Items) < perPage {
			break // no more pages
		}
		page++
	}

	// Return everything, even though we streamed each file already
	return allValues, nil
}

// GetDockerfile fetches a single Dockerfile from GitHub, decodes the base64 content,
// checks line count, uses `utils.ExtractExternalBaseImagesFromBase64` to parse external images.
// If parse fails, we store an empty Images slice.
func GetDockerfile(
	ctx context.Context,
	githubClient model.GitHubClient,
	organizationName, repoFullName, filePath string,
	stream *models.StreamSender,
) (*models.Resource, error) {

	sdk := resilientbridge.NewResilientBridge()
	sdk.SetDebug(false)
	sdk.RegisterProvider("github", adapters.NewGitHubAdapter(githubClient.Token), &resilientbridge.ProviderConfig{
		UseProviderLimits: true,
		MaxRetries:        3,
		BaseBackoff:       time.Second,
	})

	// 1) Fetch the file content from GitHub
	contentEndpoint := fmt.Sprintf("/repos/%s/contents/%s", repoFullName, url.PathEscape(filePath))
	contentReq := &resilientbridge.NormalizedRequest{
		Method:   "GET",
		Endpoint: contentEndpoint,
		Headers:  map[string]string{"Accept": "application/vnd.github+json"},
	}
	contentResp, err := sdk.Request("github", contentReq)
	if err != nil {
		return nil, fmt.Errorf("error fetching content for %s/%s: %w", repoFullName, filePath, err)
	}
	if contentResp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP error %d fetching content for %s/%s: %s",
			contentResp.StatusCode, repoFullName, filePath, string(contentResp.Data))
	}

	var contentData model.ContentResponse
	if err := json.Unmarshal(contentResp.Data, &contentData); err != nil {
		return nil, fmt.Errorf("error parsing content response for %s/%s: %w", repoFullName, filePath, err)
	}

	// 2) We rely on base64 content
	dockerfileB64 := contentData.Content
	if dockerfileB64 == "" {
		return nil, fmt.Errorf("no base64 content found for %s/%s", repoFullName, filePath)
	}

	// 3) Decode for line count
	decoded, err := base64.StdEncoding.DecodeString(dockerfileB64)
	if err != nil {
		return nil, fmt.Errorf("error decoding base64 for %s/%s: %w", repoFullName, filePath, err)
	}
	lines := strings.Split(string(decoded), "\n")
	if len(lines) > MAX_DOCKERFILE_LEN {
		return nil, fmt.Errorf("skipping %s/%s: Dockerfile has %d lines (> %d)",
			repoFullName, filePath, len(lines), MAX_DOCKERFILE_LEN)
	}

	// 4) Parse via resilient-bridge/utils
	images, parseErr := utils.ExtractExternalBaseImagesFromBase64(dockerfileB64)
	if parseErr != nil {
		log.Printf("Parsing error for Dockerfile at %s/%s: %v\n", repoFullName, filePath, parseErr)
		images = []string{}
	}

	// 5) Last updated date
	var lastUpdatedAt string
	commitsEndpoint := fmt.Sprintf("/repos/%s/commits?path=%s&per_page=1", repoFullName, url.QueryEscape(filePath))
	commitReq := &resilientbridge.NormalizedRequest{
		Method:   "GET",
		Endpoint: commitsEndpoint,
		Headers:  map[string]string{"Accept": "application/vnd.github+json"},
	}
	commitResp, err := sdk.Request("github", commitReq)
	if err == nil && commitResp.StatusCode < 400 {
		var commits []model.CommitResponse
		if json.Unmarshal(commitResp.Data, &commits) == nil && len(commits) > 0 {
			if commits[0].Commit.Author.Date != "" {
				lastUpdatedAt = commits[0].Commit.Author.Date
			} else if commits[0].Commit.Committer.Date != "" {
				lastUpdatedAt = commits[0].Commit.Committer.Date
			}
		}
	}

	repoObj := map[string]interface{}{
		"full_name": repoFullName,
	}

	output := model.ArtifactDockerFileDescription{
		Sha:  &contentData.Sha,
		Name: &contentData.Name,
		//Path:                    &contentData.Path,
		LastUpdatedAt: &lastUpdatedAt,
		//GitURL:                  &contentData.GitURL,
		HTMLURL: &contentData.HTMLURL,
		//URI:                     &contentData.HTMLURL,
		DockerfileContent:       string(decoded),
		DockerfileContentBase64: &dockerfileB64,
		Repository:              repoObj,
		Images:                  images,
	}

	value := models.Resource{
		ID:          contentData.HTMLURL,
		Name:        *output.Name,
		Description: output,
	}
	return &value, nil
}
