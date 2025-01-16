package describers

import (
	"context"
	"strconv"

	"github.com/opengovern/og-describer-template/discovery/pkg/models"
	model "github.com/opengovern/og-describer-template/discovery/provider"
	"github.com/shurcooL/githubv4"
	steampipemodels "github.com/turbot/steampipe-plugin-github/github/models"
)

func GetAllRepositoriesDeployments(ctx context.Context, githubClient model.GitHubClient, organizationName string, stream *models.StreamSender) ([]models.Resource, error) {
	client := githubClient.RestClient

	var repositoryName string
	if value := ctx.Value(paramKeyRepoName); value != nil {
		repositoryName = value.(string)
	}

	if repositoryName != "" {
		repoValues, err := GetRepositoryDeployments(ctx, githubClient, stream, organizationName, repositoryName)
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
		repoValues, err := GetRepositoryDeployments(ctx, githubClient, stream, organizationName, repo.GetName())
		if err != nil {
			return nil, err
		}
		values = append(values, repoValues...)
	}
	return values, nil
}

func GetRepositoryDeployments(ctx context.Context, githubClient model.GitHubClient, stream *models.StreamSender, owner, repo string) ([]models.Resource, error) {
	client := githubClient.GraphQLClient
	var query struct {
		RateLimit  steampipemodels.RateLimit
		Repository struct {
			Deployments struct {
				PageInfo   steampipemodels.PageInfo
				TotalCount int
				Nodes      []steampipemodels.Deployment
			} `graphql:"deployments(first: $pageSize, after: $cursor)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}
	variables := map[string]interface{}{
		"owner":    githubv4.String(owner),
		"name":     githubv4.String(repo),
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
	}
	appendRepoDeploymentColumnIncludes(&variables, repositoryDeploymentsCols())
	repoFullName := formRepositoryFullName(owner, repo)
	var values []models.Resource
	for {
		err := client.Query(ctx, &query, variables)
		if err != nil {
			return nil, err
		}
		for _, deployment := range query.Repository.Deployments.Nodes {
			value := models.Resource{
				ID:   strconv.Itoa(deployment.Id),
				Name: strconv.Itoa(deployment.Id),
				Description: model.RepoDeploymentDescription{
					Deployment:   deployment,
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
		if !query.Repository.Deployments.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Repository.Deployments.PageInfo.EndCursor)
	}
	return values, nil
}
