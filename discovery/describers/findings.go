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

func ListFindings(ctx context.Context, handler *provider.SemGrepAPIHandler, stream *models.StreamSender) ([]models.Resource, error) {
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
			if err := processFindings(ctx, handler, deployment.Slug, semGrepChan, &wg); err != nil {
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

func processFindings(ctx context.Context, handler *provider.SemGrepAPIHandler, deploymentSlug string, semGrepChan chan<- models.Resource, wg *sync.WaitGroup) error {
	var findings []provider.FindingObject
	var findingListResponse provider.FindingsListResponse
	var resp *http.Response
	baseURL := "https://semgrep.dev/api/v1/deployments/"
	page := 0

	for {
		params := url.Values{}
		params.Set("page", strconv.Itoa(page))
		params.Set("page_size", "3000")
		finalURL := fmt.Sprintf("%s%s/findings?%s", baseURL, deploymentSlug, params.Encode())

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

			if e = json.NewDecoder(resp.Body).Decode(&findingListResponse); e != nil {
				return nil, fmt.Errorf("failed to decode response: %w", e)
			}
			findings = append(findings, findingListResponse.Findings...)
			return resp, nil
		}

		err = handler.DoRequest(ctx, req, requestFunc)
		if err != nil {
			return fmt.Errorf("error during request handling: %w", err)
		}

		if len(findingListResponse.Findings) < 3000 {
			break
		}
		page++
	}

	for _, finding := range findings {
		wg.Add(1)
		go func(finding provider.FindingObject) {
			defer wg.Done()
			externalTicket := provider.ExternalTicket{
				ExternalSlug: finding.ExternalTicket.ExternalSlug,
				URL:          finding.ExternalTicket.URL,
			}
			repository := provider.Repository{
				Name: finding.Repository.Name,
				URL:  finding.Repository.URL,
			}
			location := provider.Location{
				FilePath:  finding.Location.FilePath,
				Line:      finding.Location.Line,
				Column:    finding.Location.Column,
				EndLine:   finding.Location.EndLine,
				EndColumn: finding.Location.EndColumn,
			}
			sourcingPolicy := provider.SourcingPolicy{
				ID:   finding.SourcingPolicy.ID,
				Name: finding.SourcingPolicy.Name,
				Slug: finding.SourcingPolicy.Slug,
			}
			rule := provider.Rule{
				Name:                 finding.Rule.Name,
				Message:              finding.Rule.Message,
				Confidence:           finding.Rule.Confidence,
				Category:             finding.Rule.Category,
				Subcategories:        finding.Rule.Subcategories,
				VulnerabilityClasses: finding.Rule.VulnerabilityClasses,
				CWENames:             finding.Rule.CWENames,
				OWASPNames:           finding.Rule.OWASPNames,
			}
			assistant := provider.Assistant{
				Autofix:    finding.Assistant.Autofix,
				Guidance:   finding.Assistant.Guidance,
				Autotriage: finding.Assistant.Autotriage,
				Component:  finding.Assistant.Component,
			}
			value := models.Resource{
				ID:   strconv.Itoa(finding.ID),
				Name: strconv.Itoa(finding.ID),
				Description: provider.FindingDescription{
					ID:              finding.ID,
					Ref:             finding.Ref,
					FirstSeenScanID: finding.FirstSeenScanID,
					SyntacticID:     finding.SyntacticID,
					MatchBasedID:    finding.MatchBasedID,
					ExternalTicket:  externalTicket,
					Repository:      repository,
					LineOfCodeURL:   finding.LineOfCodeURL,
					TriageState:     finding.TriageState,
					State:           finding.State,
					Status:          finding.Status,
					Severity:        finding.Severity,
					Confidence:      finding.Confidence,
					Categories:      finding.Categories,
					CreatedAt:       finding.CreatedAt,
					RelevantSince:   finding.RelevantSince,
					RuleName:        finding.RuleName,
					RuleMessage:     finding.RuleMessage,
					Location:        location,
					SourcingPolicy:  sourcingPolicy,
					TriagedAt:       finding.TriagedAt,
					TriageComment:   finding.TriageComment,
					TriageReason:    finding.TriageReason,
					StateUpdatedAt:  finding.StateUpdatedAt,
					Rule:            rule,
					Assistant:       assistant,
				},
			}
			semGrepChan <- value
		}(finding)
	}
	return nil
}
