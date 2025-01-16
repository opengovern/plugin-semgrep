package describers

import (
	"context"
	"math"
	"strconv"

	"github.com/opengovern/og-describer-template/discovery/pkg/models"
	model "github.com/opengovern/og-describer-template/discovery/provider"
	"github.com/shurcooL/githubv4"
	steampipemodels "github.com/turbot/steampipe-plugin-github/github/models"
)

func GetAllPullRequests(ctx context.Context, githubClient model.GitHubClient, organizationName string, stream *models.StreamSender) ([]models.Resource, error) {
	org := ctx.Value("organization")
	if org != nil {
		orgName := org.(string)
		if orgName != "" {
			organizationName = orgName
		}
	}

	var values []models.Resource
	repoParam := ctx.Value("repository")
	if repoParam != nil {
		repoName := repoParam.(string)
		if repoName != "" {
			repoValues, err := ListRepositoryPullRequests(ctx, githubClient, stream, organizationName, repoName)
			if err != nil {
				return nil, err
			}
			values = append(values, repoValues...)
		}
	} else {
		repositories, err := getRepositories(ctx, githubClient.RestClient, organizationName)
		if err != nil {
			return nil, nil
		}
		for _, repo := range repositories {
			repoValues, err := ListRepositoryPullRequests(ctx, githubClient, stream, organizationName, repo.GetName())
			if err != nil {
				return nil, err
			}
			values = append(values, repoValues...)
		}
	}

	return values, nil
}

func ListRepositoryPullRequests(ctx context.Context, githubClient model.GitHubClient, stream *models.StreamSender, owner, repo string) ([]models.Resource, error) {
	client := githubClient.GraphQLClient
	states := []githubv4.PullRequestState{githubv4.PullRequestStateOpen, githubv4.PullRequestStateClosed, githubv4.PullRequestStateMerged}
	var query struct {
		RateLimit  steampipemodels.RateLimit
		Repository struct {
			PullRequests struct {
				PageInfo   steampipemodels.PageInfo
				TotalCount int
				Nodes      []steampipemodels.PullRequest
			} `graphql:"pullRequests(first: $pageSize, after: $cursor, states: $states)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}
	variables := map[string]interface{}{
		"owner":    githubv4.String(owner),
		"name":     githubv4.String(repo),
		"pageSize": githubv4.Int(pullRequestsPageSize),
		"cursor":   (*githubv4.String)(nil),
		"states":   states,
	}
	appendPullRequestColumnIncludes(&variables, pullRequestCols())
	repoFullName := formRepositoryFullName(owner, repo)
	var values []models.Resource
	for {
		err := client.Query(ctx, &query, variables)
		if err != nil {
			return nil, err
		}
		for _, issue := range query.Repository.PullRequests.Nodes {
			labelsSrcLength := int(math.Min(float64(len(issue.Labels.Nodes)), 100.0))
			labelsSrc := issue.Labels.Nodes[:labelsSrcLength]
			labels := make(map[string]steampipemodels.Label)
			for _, label := range issue.Labels.Nodes {
				labels[label.Name] = label
			}
			value := models.Resource{
				ID:   strconv.Itoa(issue.Id),
				Name: strconv.Itoa(issue.Number),
				Description: model.PullRequestDescription{
					RepoFullName:             repoFullName,
					Id:                       issue.Id,
					NodeId:                   issue.NodeId,
					Number:                   issue.Number,
					ActiveLockReason:         issue.ActiveLockReason,
					Additions:                issue.Additions,
					Author:                   issue.Author,
					AuthorAssociation:        issue.AuthorAssociation,
					BaseRefName:              issue.BaseRefName,
					Body:                     issue.Body,
					ChangedFiles:             issue.ChangedFiles,
					ChecksUrl:                issue.ChecksUrl,
					Closed:                   issue.Closed,
					ClosedAt:                 issue.ClosedAt,
					CreatedAt:                issue.CreatedAt,
					CreatedViaEmail:          issue.CreatedViaEmail,
					Deletions:                issue.Deletions,
					Editor:                   issue.Editor,
					HeadRefName:              issue.HeadRefName,
					HeadRefOid:               issue.HeadRefOid,
					IncludesCreatedEdit:      issue.IncludesCreatedEdit,
					IsCrossRepository:        issue.IsCrossRepository,
					IsDraft:                  issue.IsDraft,
					IsReadByUser:             issue.IsReadByUser,
					LastEditedAt:             issue.LastEditedAt,
					Locked:                   issue.Locked,
					MaintainerCanModify:      issue.MaintainerCanModify,
					Mergeable:                issue.Mergeable,
					Merged:                   issue.Merged,
					MergedAt:                 issue.MergedAt,
					MergedBy:                 issue.MergedBy,
					Milestone:                issue.Milestone,
					Permalink:                issue.Permalink,
					PublishedAt:              issue.PublishedAt,
					RevertUrl:                issue.RevertUrl,
					ReviewDecision:           issue.ReviewDecision,
					State:                    issue.State,
					Title:                    issue.Title,
					TotalCommentsCount:       issue.TotalCommentsCount,
					UpdatedAt:                issue.UpdatedAt,
					Url:                      issue.Url,
					Assignees:                issue.Assignees.Nodes,
					BaseRef:                  issue.BaseRef,
					HeadRef:                  issue.HeadRef,
					MergeCommit:              issue.MergeCommit,
					SuggestedReviewers:       issue.SuggestedReviewers,
					CanApplySuggestion:       issue.CanApplySuggestion,
					CanClose:                 issue.CanClose,
					CanDeleteHeadRef:         issue.CanDeleteHeadRef,
					CanDisableAutoMerge:      issue.CanDisableAutoMerge,
					CanEditFiles:             issue.CanEditFiles,
					CanEnableAutoMerge:       issue.CanEnableAutoMerge,
					CanMergeAsAdmin:          issue.CanMergeAsAdmin,
					CanReact:                 issue.CanReact,
					CanReopen:                issue.CanReopen,
					CanSubscribe:             issue.CanSubscribe,
					CanUpdate:                issue.CanUpdate,
					CanUpdateBranch:          issue.CanUpdateBranch,
					DidAuthor:                issue.DidAuthor,
					CannotUpdateReasons:      issue.CannotUpdateReasons,
					Subscription:             issue.Subscription,
					LabelsSrc:                labelsSrc,
					Labels:                   labels,
					CommitsTotalCount:        issue.Commits.TotalCount,
					ReviewRequestsTotalCount: issue.ReviewRequests.TotalCount,
					ReviewsTotalCount:        issue.Reviews.TotalCount,
					LabelsTotalCount:         issue.Labels.TotalCount,
					AssigneesTotalCount:      issue.Assignees.TotalCount,
				},
			}
			if stream != nil {
				if err := (*stream)(value); err != nil {
					return nil, err
				}
			} else {
				values = append(values, value)
			}
		}
		if !query.Repository.PullRequests.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Repository.PullRequests.PageInfo.EndCursor)
	}
	return values, nil
}
