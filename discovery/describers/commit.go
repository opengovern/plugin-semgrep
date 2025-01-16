// commit.go
package describers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/opengovern/og-describer-template/discovery/pkg/models"
	model "github.com/opengovern/og-describer-template/discovery/provider"
	resilientbridge "github.com/opengovern/resilient-bridge"
	"github.com/opengovern/resilient-bridge/adapters"
)

// ListCommits fetches commits from all active repositories under the specified organization.
// If a stream is provided, each commit is sent to the stream as it’s processed.
// Otherwise, commits are collected and returned as a slice.
func ListCommits(ctx context.Context, githubClient model.GitHubClient, organizationName string, stream *models.StreamSender) ([]models.Resource, error) {
	// Retrieve repositories while excluding archived and disabled ones
	//repos, err := GetRepositoryListWithOptions(ctx, githubClient, organizationName, nil, true, true)
	//if err != nil {
	//	return nil, err
	//}
	repos, err := getRepositories(ctx, githubClient.RestClient, organizationName)
	if err != nil {
		return nil, err
	}

	sdk := newResilientSDK(githubClient.Token)

	var values []models.Resource
	for _, r := range repos {
		// r.Name should correspond to the repository name
		repoValues, err := GetRepositoryCommits(ctx, sdk, stream, organizationName, r.GetName())
		if err != nil {
			return nil, err
		}
		values = append(values, repoValues...)
	}

	return values, nil
}

// GetRepositoryCommits fetches up to 50 commits for a single repository.
// If a stream is provided, commits are streamed; otherwise, returns them as a slice.
func GetRepositoryCommits(ctx context.Context, sdk *resilientbridge.ResilientBridge, stream *models.StreamSender, owner, repo string) ([]models.Resource, error) {
	maxCommits := 10
	commits, err := fetchCommitList(sdk, owner, repo, maxCommits)
	if err != nil {
		if strings.Contains(err.Error(), "client error: 409") {
			return []models.Resource{}, nil
		}
		return nil, fmt.Errorf("error fetching commits list for %s/%s: %w", owner, repo, err)
	}

	if len(commits) == 0 {
		log.Printf("No commits found for %s/%s (possibly empty default branch).", owner, repo)
		return nil, nil
	}

	// Determine concurrency level from env or default to 3
	concurrency := 3
	if cStr := os.Getenv("CONCURRENCY"); cStr != "" {
		if cVal, err := strconv.Atoi(cStr); err == nil && cVal > 0 {
			concurrency = cVal
		}
	}
	log.Printf("Fetching commit details with concurrency=%d", concurrency)

	results := make([]models.Resource, len(commits))

	type job struct {
		index int
		sha   string
	}
	jobCh := make(chan job)
	wg := sync.WaitGroup{}

	worker := func() {
		defer wg.Done()
		for j := range jobCh {
			commitJSON, err := fetchCommitDetails(sdk, owner, repo, j.sha)
			if err != nil {
				log.Printf("Error fetching commit %s details: %v", j.sha, err)
				continue
			}

			var commit model.CommitResp
			if err := json.Unmarshal(commitJSON, &commit); err != nil {
				log.Printf("Error unmarshaling JSON for commit %s: %v", j.sha, err)
				continue
			}

			tree := model.Tree{
				SHA: &commit.CommitDetail.Tree.SHA,
				URL: &commit.CommitDetail.Tree.URL,
			}

			verification := model.Verification{
				Verified:   commit.CommitDetail.Verification.Verified,
				Reason:     &commit.CommitDetail.Verification.Reason,
				Signature:  commit.CommitDetail.Verification.Signature,
				Payload:    commit.CommitDetail.Verification.Payload,
				VerifiedAt: commit.CommitDetail.Verification.VerifiedAt,
			}

			commitDetail := model.CommitDetail{
				Message:      &commit.CommitDetail.Message,
				Tree:         tree,
				CommentCount: commit.CommitDetail.CommentCount,
				Verification: verification,
			}

			author := model.User{
				Login:             &commit.Author.Login,
				ID:                commit.Author.ID,
				NodeID:            &commit.Author.NodeID,
				AvatarURL:         &commit.Author.AvatarURL,
				GravatarID:        &commit.Author.GravatarID,
				URL:               &commit.Author.URL,
				HTMLURL:           &commit.Author.HTMLURL,
				FollowersURL:      &commit.Author.FollowersURL,
				FollowingURL:      &commit.Author.FollowingURL,
				GistsURL:          &commit.Author.GistsURL,
				StarredURL:        &commit.Author.StarredURL,
				SubscriptionsURL:  &commit.Author.SubscriptionsURL,
				OrganizationsURL:  &commit.Author.OrganizationsURL,
				ReposURL:          &commit.Author.ReposURL,
				EventsURL:         &commit.Author.EventsURL,
				ReceivedEventsURL: &commit.Author.ReceivedEventsURL,
				Type:              &commit.Author.Type,
				UserViewType:      &commit.Author.UserViewType,
				SiteAdmin:         commit.Author.SiteAdmin,
			}

			commiter := model.User{
				Login:             &commit.Committer.Login,
				ID:                commit.Committer.ID,
				NodeID:            &commit.Committer.NodeID,
				AvatarURL:         &commit.Committer.AvatarURL,
				GravatarID:        &commit.Committer.GravatarID,
				URL:               &commit.Committer.URL,
				HTMLURL:           &commit.Committer.HTMLURL,
				FollowersURL:      &commit.Committer.FollowersURL,
				FollowingURL:      &commit.Committer.FollowingURL,
				GistsURL:          &commit.Committer.GistsURL,
				StarredURL:        &commit.Committer.StarredURL,
				SubscriptionsURL:  &commit.Committer.SubscriptionsURL,
				OrganizationsURL:  &commit.Committer.OrganizationsURL,
				ReposURL:          &commit.Committer.ReposURL,
				EventsURL:         &commit.Committer.EventsURL,
				ReceivedEventsURL: &commit.Committer.ReceivedEventsURL,
				Type:              &commit.Committer.Type,
				UserViewType:      &commit.Committer.UserViewType,
				SiteAdmin:         commit.Committer.SiteAdmin,
			}

			var parents []model.Parent
			for _, parent := range commit.Parents {
				parents = append(parents, model.Parent{
					SHA:     &parent.SHA,
					URL:     &parent.URL,
					HTMLURL: &parent.HTMLURL,
				})
			}

			stats := model.Stats{
				Total:     commit.Stats.Total,
				Additions: commit.Stats.Additions,
				Deletions: commit.Stats.Deletions,
			}

			var files []model.File
			for _, file := range files {
				files = append(files, model.File{
					SHA:         file.SHA,
					Filename:    file.Filename,
					Status:      file.Status,
					Additions:   file.Additions,
					Deletions:   file.Deletions,
					Changes:     file.Changes,
					BlobURL:     file.BlobURL,
					RawURL:      file.RawURL,
					ContentsURL: file.ContentsURL,
					Patch:       file.Patch,
				})
			}

			value := models.Resource{
				ID:   commit.SHA,
				Name: commit.SHA,
				Description: model.CommitDescription{
					SHA:          &commit.SHA,
					NodeID:       &commit.NodeID,
					CommitDetail: commitDetail,
					URL:          &commit.URL,
					HTMLURL:      &commit.HTMLURL,
					CommentsURL:  &commit.CommentsURL,
					Author:       author,
					Committer:    commiter,
					Parents:      parents,
					Stats:        stats,
					Files:        files,
				},
			}
			results[j.index] = value
		}
	}

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go worker()
	}

	for i, c := range commits {
		jobCh <- job{index: i, sha: c.SHA}
	}
	close(jobCh)

	wg.Wait()

	var finalResults []models.Resource
	for _, res := range results {
		if res.ID != "" {
			if stream != nil {
				if err := (*stream)(res); err != nil {
					return nil, err
				}
			} else {
				finalResults = append(finalResults, res)
			}
		}
	}
	return finalResults, nil
}

type commitRef struct {
	SHA string `json:"sha"`
}

// fetchCommitList returns up to maxCommits commit references from the repo’s default branch.
func fetchCommitList(sdk *resilientbridge.ResilientBridge, owner, repo string, maxCommits int) ([]commitRef, error) {
	var allCommits []commitRef
	perPage := 100
	page := 1

	for len(allCommits) < maxCommits {
		remaining := maxCommits - len(allCommits)
		if remaining < perPage {
			perPage = remaining
		}

		endpoint := fmt.Sprintf("/repos/%s/%s/commits?per_page=%d&page=%d", owner, repo, perPage, page)
		req := &resilientbridge.NormalizedRequest{
			Method:   "GET",
			Endpoint: endpoint,
			Headers:  map[string]string{"Accept": "application/vnd.github+json"},
		}

		resp, err := sdk.Request("github", req)
		if err != nil {
			return nil, fmt.Errorf("error fetching commits: %w", err)
		}

		// Handle HTTP errors
		if resp.StatusCode == 409 {
			// 409 typically means no commits on default branch or empty repo
			// Treat this as no commits found.
			return []commitRef{}, nil
		}
		if resp.StatusCode >= 400 {
			return nil, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(resp.Data))
		}

		var commits []commitRef
		if err := json.Unmarshal(resp.Data, &commits); err != nil {
			return nil, fmt.Errorf("error decoding commit list: %w", err)
		}

		if len(commits) == 0 {
			// No more commits
			break
		}

		allCommits = append(allCommits, commits...)
		if len(allCommits) >= maxCommits {
			break
		}
		page++
	}

	if len(allCommits) > maxCommits {
		allCommits = allCommits[:maxCommits]
	}

	return allCommits, nil
}

func fetchCommitDetails(sdk *resilientbridge.ResilientBridge, owner, repo, sha string) ([]byte, error) {
	req := &resilientbridge.NormalizedRequest{
		Method:   "GET",
		Endpoint: fmt.Sprintf("/repos/%s/%s/commits/%s", owner, repo, sha),
		Headers:  map[string]string{"Accept": "application/vnd.github+json"},
	}
	resp, err := sdk.Request("github", req)
	if err != nil {
		return nil, fmt.Errorf("error fetching commit details: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(resp.Data))
	}

	return resp.Data, nil
}

func newResilientSDK(token string) *resilientbridge.ResilientBridge {
	sdk := resilientbridge.NewResilientBridge()
	sdk.RegisterProvider("github", adapters.NewGitHubAdapter(token), &resilientbridge.ProviderConfig{
		UseProviderLimits: true,
		MaxRetries:        3,
		BaseBackoff:       0,
	})
	return sdk
}
