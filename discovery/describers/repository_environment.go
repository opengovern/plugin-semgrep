package describers

import (
	"context"
	"strconv"

	"github.com/opengovern/og-describer-template/discovery/pkg/models"
	model "github.com/opengovern/og-describer-template/discovery/provider"
	"github.com/shurcooL/githubv4"
	steampipemodels "github.com/turbot/steampipe-plugin-github/github/models"
)

func GetAllRepositoriesEnvironments(ctx context.Context, githubClient model.GitHubClient, organizationName string, stream *models.StreamSender) ([]models.Resource, error) {
	client := githubClient.RestClient

	var repositoryName string
	if value := ctx.Value(paramKeyRepoName); value != nil {
		repositoryName = value.(string)
	}

	if repositoryName != "" {
		repoValues, err := GetRepositoryEnvironments(ctx, githubClient, stream, organizationName, repositoryName)
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
		repoValues, err := GetRepositoryEnvironments(ctx, githubClient, stream, organizationName, repo.GetName())
		if err != nil {
			return nil, err
		}
		values = append(values, repoValues...)
	}
	return values, nil
}

func GetRepositoryEnvironments(ctx context.Context, githubClient model.GitHubClient, stream *models.StreamSender, owner, repo string) ([]models.Resource, error) {
	client := githubClient.GraphQLClient
	var query struct {
		RateLimit  steampipemodels.RateLimit
		Repository struct {
			Environments struct {
				PageInfo   steampipemodels.PageInfo
				TotalCount int
				Nodes      []steampipemodels.Environment
			} `graphql:"environments(first: $pageSize, after: $cursor)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}
	variables := map[string]interface{}{
		"owner":    githubv4.String(owner),
		"name":     githubv4.String(repo),
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
	}
	appendRepoEnvironmentColumnIncludes(&variables, repositoryEnvironmentsCols())
	repoFullName := formRepositoryFullName(owner, repo)
	var values []models.Resource
	for {
		err := client.Query(ctx, &query, variables)
		if err != nil {
			return nil, err
		}
		for _, environment := range query.Repository.Environments.Nodes {
			value := models.Resource{
				ID:   strconv.Itoa(environment.Id),
				Name: environment.Name,
				Description: model.RepoEnvironmentDescription{
					Environment:  environment,
					RepoFullName: repoFullName,
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
		if !query.Repository.Environments.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Repository.Environments.PageInfo.EndCursor)
	}
	return values, nil
}
