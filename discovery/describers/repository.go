package describers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/opengovern/og-describer-template/discovery/pkg/models"
	model "github.com/opengovern/og-describer-template/discovery/provider"
	resilientbridge "github.com/opengovern/resilient-bridge"
	"github.com/opengovern/resilient-bridge/adapters"
)

// MAX_REPOS_TO_LIST is how many repositories to fetch at most when listing.
const MAX_REPOS_TO_LIST = 250

// GetRepositoryList calls GetRepositoryListWithOptions with default excludes (false, false).
func GetRepositoryList(
	ctx context.Context,
	githubClient model.GitHubClient,
	organizationName string,
	stream *models.StreamSender,
) ([]models.Resource, error) {
	return GetRepositoryListWithOptions(ctx, githubClient, organizationName, stream, false, false)
}

// GetRepositoryListWithOptions enumerates org repositories, but *instead of streaming each raw JSON*,
// it calls GetRepository(...) for each repo, so that *only the final detail* from GetRepository is
// streamed and returned.
func GetRepositoryListWithOptions(
	ctx context.Context,
	githubClient model.GitHubClient,
	organizationName string,
	stream *models.StreamSender,
	excludeArchived bool,
	excludeDisabled bool,
) ([]models.Resource, error) {

	sdk := resilientbridge.NewResilientBridge()
	sdk.RegisterProvider("github", adapters.NewGitHubAdapter(githubClient.Token), &resilientbridge.ProviderConfig{
		UseProviderLimits: true,
		MaxRetries:        3,
		BaseBackoff:       0,
	})

	var allFinalResources []models.Resource
	perPage := 100
	page := 1

	org := ctx.Value("organization")
	if org != nil {
		orgName := org.(string)
		if orgName != "" {
			organizationName = orgName
		}
	}

	repo := ctx.Value("repository")
	if repo != nil {
		repoName := repo.(string)
		if repoName != "" {
			finalResource, err := GetRepository(
				ctx,
				githubClient,
				organizationName,
				repoName,
				"", // pass along or just ""
				stream,
			)
			if err != nil {
				return nil, fmt.Errorf("error fetching details for %s/%s: %v", organizationName, repoName, err)
			}

			return []models.Resource{*finalResource}, nil
		}
	}

	for len(allFinalResources) < MAX_REPOS_TO_LIST {
		endpoint := fmt.Sprintf("/orgs/%s/repos?per_page=%d&page=%d", organizationName, perPage, page)
		req := &resilientbridge.NormalizedRequest{
			Method:   "GET",
			Endpoint: endpoint,
			Headers:  map[string]string{"Accept": "application/vnd.github+json"},
		}

		resp, err := sdk.Request("github", req)
		if err != nil {
			return nil, fmt.Errorf("error fetching repos: %w", err)
		}
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(resp.Data))
		}

		// Decode into a slice of generic maps. We'll only parse name, archived, disabled, etc.
		var repos []map[string]interface{}
		if err := json.Unmarshal(resp.Data, &repos); err != nil {
			return nil, fmt.Errorf("error decoding repos list: %w", err)
		}
		if len(repos) == 0 {
			break
		}

		for _, r := range repos {
			// Filter archived, disabled if requested
			if excludeArchived {
				if archived, ok := r["archived"].(bool); ok && archived {
					continue
				}
			}
			if excludeDisabled {
				if disabled, ok := r["disabled"].(bool); ok && disabled {
					continue
				}
			}

			// Grab the repo name from the raw JSON
			nameStr, _ := r["name"].(string)
			if nameStr == "" {
				continue
			}

			// Now call GetRepository to get the *final* detail
			// resourceID can be empty or you can use the raw "id" from r if you like.
			var resourceID string
			if idVal, ok := r["id"]; ok {
				resourceID = fmt.Sprintf("%v", idVal)
			}
			finalResource, err := GetRepository(
				ctx,
				githubClient,
				organizationName,
				nameStr,
				resourceID, // pass along or just ""
				stream,     // same stream pointer
			)
			if err != nil {
				log.Printf("Error fetching details for %s/%s: %v", organizationName, nameStr, err)
				continue
			}

			// Append the final resource from GetRepository into our big slice
			if finalResource != nil {
				allFinalResources = append(allFinalResources, *finalResource)
				if len(allFinalResources) >= MAX_REPOS_TO_LIST {
					break
				}
			}
		}

		if len(repos) < perPage {
			// No more pages
			break
		}
		page++
	}

	return allFinalResources, nil
}

// GetRepository fetches a single repo, transforms it, fetches languages, enriches metrics, returns a single Resource.
func GetRepository(
	ctx context.Context,
	githubClient model.GitHubClient,
	organizationName string,
	repositoryName string,
	resourceID string, // optional
	stream *models.StreamSender,
) (*models.Resource, error) {

	sdk := resilientbridge.NewResilientBridge()
	sdk.RegisterProvider("github", adapters.NewGitHubAdapter(githubClient.Token), &resilientbridge.ProviderConfig{
		UseProviderLimits: true,
		MaxRetries:        3,
		BaseBackoff:       0,
	})

	// 1) Fetch RepoDetail
	repoDetail, err := util_fetchRepoDetails(sdk, organizationName, repositoryName)
	if err != nil {
		return nil, fmt.Errorf("error fetching repository details for %s/%s: %w",
			organizationName, repositoryName, err)
	}

	// 2) Transform -> RepositoryDescription
	finalDetail := util_transformToFinalRepoDetail(repoDetail)

	// 3) Fetch /languages => map[string]int
	langs, err := util_fetchLanguages(sdk, organizationName, repositoryName)
	if err == nil && len(langs) > 0 {
		finalDetail.Languages = langs
	}

	// 4) Enrich with metrics
	if err := util_enrichRepoMetrics(sdk, organizationName, repositoryName, finalDetail); err != nil {
		log.Printf("Error enriching repo metrics for %s/%s: %v",
			organizationName, repositoryName, err)
	}

	// 5) Build final Resource
	// If resourceID is empty, use the finalDetail's ID
	if resourceID == "" {
		resourceID = strconv.Itoa(finalDetail.GitHubRepoID)
	}
	value := models.Resource{
		ID:          resourceID,
		Name:        *finalDetail.Name,
		Description: finalDetail,
	}

	// Stream if provided
	if stream != nil {
		if err := (*stream)(value); err != nil {
			return nil, fmt.Errorf("streaming resource failed: %w", err)
		}
	}

	return &value, nil
}

// -----------------------------------------------------------------------------
// Utility / helper functions
// (unchanged from your existing code except for naming and minor comments)
// -----------------------------------------------------------------------------

func util_fetchRepoDetails(sdk *resilientbridge.ResilientBridge, owner, repo string) (*model.RepoDetail, error) {
	req := &resilientbridge.NormalizedRequest{
		Method:   "GET",
		Endpoint: fmt.Sprintf("/repos/%s/%s", owner, repo),
		Headers:  map[string]string{"Accept": "application/vnd.github+json"},
	}
	resp, err := sdk.Request("github", req)
	if err != nil {
		return nil, fmt.Errorf("error fetching repo details: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(resp.Data))
	}

	var detail model.RepoDetail
	if err := json.Unmarshal(resp.Data, &detail); err != nil {
		return nil, fmt.Errorf("error decoding repo details: %w", err)
	}
	return &detail, nil
}

func util_transformToFinalRepoDetail(detail *model.RepoDetail) *model.RepositoryDescription {
	var parent *model.RepositoryDescription
	if detail.Parent != nil {
		parent = util_transformToFinalRepoDetail(detail.Parent)
	}
	var source *model.RepositoryDescription
	if detail.Source != nil {
		source = util_transformToFinalRepoDetail(detail.Source)
	}

	var finalOwner *model.Owner
	if detail.Owner != nil {
		finalOwner = &model.Owner{
			Login:   detail.Owner.Login,
			ID:      detail.Owner.ID,
			NodeID:  detail.Owner.NodeID,
			HTMLURL: detail.Owner.HTMLURL,
			Type:    detail.Owner.Type,
		}
	}

	var finalOrg *model.Organization
	if detail.Organization != nil {
		finalOrg = &model.Organization{
			Login:        detail.Organization.Login,
			ID:           detail.Organization.ID,
			NodeID:       detail.Organization.NodeID,
			HTMLURL:      detail.Organization.HTMLURL,
			Type:         detail.Organization.Type,
			UserViewType: detail.Organization.UserViewType,
			SiteAdmin:    detail.Organization.SiteAdmin,
		}
	}

	dbObj := map[string]string{"name": detail.DefaultBranch}
	dbBytes, _ := json.Marshal(dbObj)

	isActive := !(detail.Archived || detail.Disabled)
	isEmpty := (detail.Size == 0)

	var licenseJSON json.RawMessage
	if detail.License != nil {
		if data, err := json.Marshal(detail.License); err == nil {
			licenseJSON = data
		}
	}

	finalDetail := &model.RepositoryDescription{
		GitHubRepoID:            detail.ID,
		NodeID:                  &detail.NodeID,
		Name:                    &detail.Name,
		NameWithOwner:           &detail.FullName,
		Description:             detail.Description,
		CreatedAt:               &detail.CreatedAt,
		UpdatedAt:               &detail.UpdatedAt,
		PushedAt:                &detail.PushedAt,
		IsActive:                isActive,
		IsEmpty:                 isEmpty,
		IsFork:                  detail.Fork,
		IsSecurityPolicyEnabled: false,
		Owner:                   finalOwner,
		HomepageURL:             detail.Homepage,
		LicenseInfo:             licenseJSON,
		Topics:                  detail.Topics,
		Visibility:              detail.Visibility,
		DefaultBranchRef:        dbBytes,
		Permissions:             detail.Permissions,
		Organization:            finalOrg,
		Parent:                  parent,
		Source:                  source,

		// Single primary language from /repos
		PrimaryLanguage: detail.PrimaryLanguage,

		// We'll fill in LanguageBreakdown after calling /languages
		Languages: nil,

		RepositorySettings: model.RepositorySettings{
			HasDiscussionsEnabled:     detail.HasDiscussions,
			HasIssuesEnabled:          detail.HasIssues,
			HasProjectsEnabled:        detail.HasProjects,
			HasWikiEnabled:            detail.HasWiki,
			MergeCommitAllowed:        detail.AllowMergeCommit,
			MergeCommitMessage:        detail.MergeCommitMessage,
			MergeCommitTitle:          detail.MergeCommitTitle,
			SquashMergeAllowed:        detail.AllowSquashMerge,
			SquashMergeCommitMessage:  detail.SquashMergeCommitMessage,
			SquashMergeCommitTitle:    detail.SquashMergeCommitTitle,
			HasDownloads:              detail.HasDownloads,
			HasPages:                  detail.HasPages,
			WebCommitSignoffRequired:  detail.WebCommitSignoffRequired,
			MirrorURL:                 detail.MirrorURL,
			AllowAutoMerge:            detail.AllowAutoMerge,
			DeleteBranchOnMerge:       detail.DeleteBranchOnMerge,
			AllowUpdateBranch:         detail.AllowUpdateBranch,
			UseSquashPRTitleAsDefault: detail.UseSquashPRTitleAsDefault,
			CustomProperties:          detail.CustomProperties,
			ForkingAllowed:            detail.AllowForking,
			IsTemplate:                detail.IsTemplate,
			AllowRebaseMerge:          detail.AllowRebaseMerge,
			Archived:                  detail.Archived,
			Disabled:                  detail.Disabled,
			Locked:                    detail.Locked,
		},
		SecuritySettings: model.SecuritySettings{
			VulnerabilityAlertsEnabled:               false,
			SecretScanningEnabled:                    false,
			SecretScanningPushProtectionEnabled:      false,
			DependabotSecurityUpdatesEnabled:         false,
			SecretScanningNonProviderPatternsEnabled: false,
			SecretScanningValidityChecksEnabled:      false,
		},
		RepoURLs: model.RepoURLs{
			GitURL:   detail.GitURL,
			SSHURL:   detail.SSHURL,
			CloneURL: detail.CloneURL,
			SVNURL:   detail.SVNURL,
			HTMLURL:  detail.HTMLURL,
		},
		Metrics: model.Metrics{
			Stargazers:  detail.StargazersCount,
			Forks:       detail.ForksCount,
			Subscribers: detail.SubscribersCount,
			Size:        detail.Size,
			OpenIssues:  detail.OpenIssuesCount,
		},
	}

	return finalDetail
}

func util_fetchLanguages(
	sdk *resilientbridge.ResilientBridge,
	owner, repo string,
) (map[string]int, error) {

	req := &resilientbridge.NormalizedRequest{
		Method:   "GET",
		Endpoint: fmt.Sprintf("/repos/%s/%s/languages", owner, repo),
		Headers:  map[string]string{"Accept": "application/vnd.github+json"},
	}
	resp, err := sdk.Request("github", req)
	if err != nil {
		return nil, fmt.Errorf("error fetching languages: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(resp.Data))
	}

	var langs map[string]int
	if err := json.Unmarshal(resp.Data, &langs); err != nil {
		return nil, fmt.Errorf("error decoding languages: %w", err)
	}
	return langs, nil
}

func util_enrichRepoMetrics(
	sdk *resilientbridge.ResilientBridge,
	owner, repoName string,
	finalDetail *model.RepositoryDescription,
) error {

	var dbObj map[string]string
	if finalDetail.DefaultBranchRef != nil {
		if err := json.Unmarshal(finalDetail.DefaultBranchRef, &dbObj); err != nil {
			return err
		}
	}
	defaultBranch := dbObj["name"]
	if defaultBranch == "" {
		defaultBranch = "main"
	}

	commitsCount, err := util_countCommits(sdk, owner, repoName, defaultBranch)
	if err != nil {
		return fmt.Errorf("counting commits: %w", err)
	}
	finalDetail.Metrics.Commits = commitsCount

	issuesCount, err := util_countIssues(sdk, owner, repoName)
	if err != nil {
		return fmt.Errorf("counting issues: %w", err)
	}
	finalDetail.Metrics.Issues = issuesCount

	branchesCount, err := util_countBranches(sdk, owner, repoName)
	if err != nil {
		return fmt.Errorf("counting branches: %w", err)
	}
	finalDetail.Metrics.Branches = branchesCount

	prCount, err := util_countPullRequests(sdk, owner, repoName)
	if err != nil {
		return fmt.Errorf("counting PRs: %w", err)
	}
	finalDetail.Metrics.PullRequests = prCount

	releasesCount, err := util_countReleases(sdk, owner, repoName)
	if err != nil {
		return fmt.Errorf("counting releases: %w", err)
	}
	finalDetail.Metrics.Releases = releasesCount

	tagsCount, err := util_countTags(sdk, owner, repoName)
	if err != nil {
		return fmt.Errorf("counting tags: %w", err)
	}
	finalDetail.Metrics.Tags = tagsCount

	return nil
}

func util_countTags(
	sdk *resilientbridge.ResilientBridge,
	owner, repoName string,
) (int, error) {
	endpoint := fmt.Sprintf("/repos/%s/%s/tags?per_page=1", owner, repoName)
	return util_countItemsFromEndpoint(sdk, endpoint)
}

func util_countCommits(
	sdk *resilientbridge.ResilientBridge,
	owner, repoName, defaultBranch string,
) (int, error) {
	endpoint := fmt.Sprintf("/repos/%s/%s/commits?sha=%s&per_page=1", owner, repoName, defaultBranch)
	return util_countItemsFromEndpoint(sdk, endpoint)
}

func util_countIssues(
	sdk *resilientbridge.ResilientBridge,
	owner, repoName string,
) (int, error) {
	endpoint := fmt.Sprintf("/repos/%s/%s/issues?state=all&per_page=1", owner, repoName)
	return util_countItemsFromEndpoint(sdk, endpoint)
}

func util_countBranches(
	sdk *resilientbridge.ResilientBridge,
	owner, repoName string,
) (int, error) {
	endpoint := fmt.Sprintf("/repos/%s/%s/branches?per_page=1", owner, repoName)
	return util_countItemsFromEndpoint(sdk, endpoint)
}

func util_countPullRequests(
	sdk *resilientbridge.ResilientBridge,
	owner, repoName string,
) (int, error) {
	endpoint := fmt.Sprintf("/repos/%s/%s/pulls?state=all&per_page=1", owner, repoName)
	return util_countItemsFromEndpoint(sdk, endpoint)
}

func util_countReleases(
	sdk *resilientbridge.ResilientBridge,
	owner, repoName string,
) (int, error) {
	endpoint := fmt.Sprintf("/repos/%s/%s/releases?per_page=1", owner, repoName)
	return util_countItemsFromEndpoint(sdk, endpoint)
}

func util_countItemsFromEndpoint(
	sdk *resilientbridge.ResilientBridge,
	endpoint string,
) (int, error) {

	req := &resilientbridge.NormalizedRequest{
		Method:   "GET",
		Endpoint: endpoint,
		Headers:  map[string]string{"Accept": "application/vnd.github+json"},
	}
	resp, err := sdk.Request("github", req)
	if err != nil {
		return 0, fmt.Errorf("error fetching data: %w", err)
	}
	if resp.StatusCode == 409 {
		return 0, nil
	}
	if resp.StatusCode >= 400 {
		return 0, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(resp.Data))
	}

	var linkHeader string
	for k, v := range resp.Headers {
		if strings.ToLower(k) == "link" {
			linkHeader = v
			break
		}
	}

	if linkHeader == "" {
		// If there's no Link header, see if the response is an array
		if len(resp.Data) > 2 {
			var items []interface{}
			if err := json.Unmarshal(resp.Data, &items); err != nil {
				return 1, nil
			}
			return len(items), nil
		}
		return 0, nil
	}

	lastPage, err := util_parseLastPage(linkHeader)
	if err != nil {
		return 0, fmt.Errorf("could not parse last page: %w", err)
	}
	return lastPage, nil
}

func util_parseLastPage(linkHeader string) (int, error) {
	re := regexp.MustCompile(`page=(\d+)>; rel="last"`)
	matches := re.FindStringSubmatch(linkHeader)
	if len(matches) < 2 {
		return 1, nil
	}
	var lastPage int
	if _, err := fmt.Sscanf(matches[1], "%d", &lastPage); err != nil {
		return 0, err
	}
	return lastPage, nil
}
