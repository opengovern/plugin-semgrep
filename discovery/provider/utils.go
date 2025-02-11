package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/time/rate"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type SemGrepAPIHandler struct {
	Client       *http.Client
	Token        string
	RateLimiter  *rate.Limiter
	Semaphore    chan struct{}
	MaxRetries   int
	RetryBackoff time.Duration
}

func NewSemGrepAPIHandler(token string, rateLimit rate.Limit, burst int, maxConcurrency int, maxRetries int, retryBackoff time.Duration) *SemGrepAPIHandler {
	return &SemGrepAPIHandler{
		Client:       http.DefaultClient,
		Token:        token,
		RateLimiter:  rate.NewLimiter(rateLimit, burst),
		Semaphore:    make(chan struct{}, maxConcurrency),
		MaxRetries:   maxRetries,
		RetryBackoff: retryBackoff,
	}
}

// DoRequest executes the Semgrep API request with rate limiting, retries, and concurrency control.
func (h *SemGrepAPIHandler) DoRequest(ctx context.Context, req *http.Request, requestFunc func(req *http.Request) (*http.Response, error)) error {
	h.Semaphore <- struct{}{}
	defer func() { <-h.Semaphore }()
	var resp *http.Response
	var err error
	for attempt := 0; attempt <= h.MaxRetries; attempt++ {
		// Wait based on rate limiter
		if err = h.RateLimiter.Wait(ctx); err != nil {
			return err
		}
		// Set request headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", h.Token))
		// Execute the request function
		resp, err = requestFunc(req)
		if err == nil {
			return nil
		}
		if resp == nil {
			return err
		}
		// Set rate limiter new value
		var resetDuration int
		if resp != nil {
			retryAfter := resp.Header.Get("Retry-After")
			if retryAfter != "" {
				resetDuration, _ = strconv.Atoi(retryAfter)
			}
			var remainRequests int
			remainRequestsStr := resp.Header.Get("X-RateLimit-Remaining")
			if remainRequestsStr != "" {
				remainRequests, err = strconv.Atoi(remainRequestsStr)
				if err == nil && resetDuration > 0 {
					h.RateLimiter = rate.NewLimiter(rate.Every(time.Duration(resetDuration)/time.Duration(remainRequests)), 1)
				}
			}
		}
		// Handle rate limit errors
		if resp != nil && resp.StatusCode == http.StatusTooManyRequests {
			if resetDuration > 0 {
				time.Sleep(time.Duration(resetDuration))
				continue
			}
			// Exponential backoff if headers are missing
			backoff := h.RetryBackoff * (1 << attempt)
			time.Sleep(backoff)
			continue
		}
		// Handle temporary network errors
		if isTemporary(err) {
			backoff := h.RetryBackoff * (1 << attempt)
			time.Sleep(backoff)
			continue
		}
		break
	}
	return err
}

// isTemporary checks if an error is temporary.
func isTemporary(err error) bool {
	if err == nil {
		return false
	}
	var netErr interface{ Temporary() bool }
	if errors.As(err, &netErr) {
		return netErr.Temporary()
	}
	return false
}

func ListDeployments(ctx context.Context, handler *SemGrepAPIHandler) ([]DeploymentJSON, error) {
	var deploymentListResponse DeploymentsResponse
	var resp *http.Response
	baseURL := "https://semgrep.dev/api/v1/deployments"

	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
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
		return nil, fmt.Errorf("error during request handling: %w", err)
	}

	return deploymentListResponse.Deployments, nil
}

func ListProjects(ctx context.Context, handler *SemGrepAPIHandler, deploymentSlug string) ([]ProjectJSON, error) {
	var projects []ProjectJSON
	var projectListResponse ProjectsListResponse
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
			return nil, fmt.Errorf("failed to create request: %w", err)
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
			return nil, fmt.Errorf("error during request handling: %w", err)
		}

		if len(projectListResponse.Projects) < 3000 {
			break
		}
		page++
	}

	return projects, nil
}
