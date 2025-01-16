package describers

import (
	"context"
	"strconv"

	"github.com/opengovern/og-describer-template/discovery/pkg/models"
	model "github.com/opengovern/og-describer-template/discovery/provider"
	"github.com/shurcooL/githubv4"
	steampipemodels "github.com/turbot/steampipe-plugin-github/github/models"
)

func GetAllRepositoriesCollaborators(ctx context.Context, githubClient model.GitHubClient, organizationName string, stream *models.StreamSender) ([]models.Resource, error) {
	client := githubClient.RestClient

	var repositoryName string
	if value := ctx.Value(paramKeyRepoName); value != nil {
		repositoryName = value.(string)
	}

	if repositoryName != "" {
		repoValues, err := GetRepositoryCollaborators(ctx, githubClient, stream, organizationName, repositoryName)
		if err != nil {
			return nil, err
		}
		return repoValues, nil
	}
	repositories, err := getRepositories(ctx, client, organizationName)
	if err != nil {
		return nil, nil
	}
	var values []models.Resource
	for _, repo := range repositories {
		repoValues, err := GetRepositoryCollaborators(ctx, githubClient, stream, organizationName, repo.GetName())
		if err != nil {
			return nil, err
		}
		values = append(values, repoValues...)
	}
	return values, nil
}

func GetRepositoryCollaborators(ctx context.Context, githubClient model.GitHubClient, stream *models.StreamSender, owner, repo string) ([]models.Resource, error) {
	client := githubClient.GraphQLClient
	affiliation := githubv4.CollaboratorAffiliationAll
	var query struct {
		RateLimit  steampipemodels.RateLimit
		Repository struct {
			Collaborators struct {
				TotalCount int
				PageInfo   steampipemodels.PageInfo
				Edges      []struct {
					Permission githubv4.RepositoryPermission `graphql:"permission @include(if:$includeRCPermission)" json:"permission"`
					Node       steampipemodels.BasicUser     `graphql:"node @include(if:$includeRCNode)" json:"node"`
				}
			} `graphql:"collaborators(first: $pageSize, after: $cursor, affiliation: $affiliation)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}
	variables := map[string]interface{}{
		"owner":       githubv4.String(owner),
		"repo":        githubv4.String(repo),
		"pageSize":    githubv4.Int(pageSize),
		"cursor":      (*githubv4.String)(nil),
		"affiliation": affiliation,
	}
	appendRepoCollaboratorColumnIncludes(&variables, repositoryCollaboratorsCols())
	repoFullName := formRepositoryFullName(owner, repo)
	var values []models.Resource
	for {
		err := client.Query(ctx, &query, variables)
		if err != nil {
			return nil, err
		}
		for _, collaborator := range query.Repository.Collaborators.Edges {
			value := models.Resource{
				ID:   strconv.Itoa(collaborator.Node.Id),
				Name: collaborator.Node.Name,
				Description: model.RepoCollaboratorsDescription{
					Affiliation:  "ALL",
					RepoFullName: repoFullName,
					Permission:   collaborator.Permission,
					UserLogin:    collaborator.Node.Login,
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
		if !query.Repository.Collaborators.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Repository.Collaborators.PageInfo.EndCursor)
	}
	return values, nil
}
