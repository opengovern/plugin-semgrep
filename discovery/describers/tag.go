package describers

import (
	"context"
	"fmt"

	"github.com/opengovern/og-describer-template/discovery/pkg/models"
	model "github.com/opengovern/og-describer-template/discovery/provider"
	"github.com/shurcooL/githubv4"
	steampipemodels "github.com/turbot/steampipe-plugin-github/github/models"
)

func GetAllTags(ctx context.Context, githubClient model.GitHubClient, organizationName string, stream *models.StreamSender) ([]models.Resource, error) {
	client := githubClient.RestClient

	repositories, err := getRepositories(ctx, client, organizationName)
	if err != nil {
		return nil, nil
	}
	var values []models.Resource
	for _, repo := range repositories {
		repoValues, err := GetRepositoryTags(ctx, githubClient, stream, organizationName, repo.GetName())
		if err != nil {
			return nil, err
		}
		values = append(values, repoValues...)
	}
	return values, nil
}

func GetRepositoryTags(ctx context.Context, githubClient model.GitHubClient, stream *models.StreamSender, owner, repo string) ([]models.Resource, error) {
	client := githubClient.GraphQLClient
	var query struct {
		RateLimit  steampipemodels.RateLimit
		Repository struct {
			Refs struct {
				TotalCount int
				PageInfo   steampipemodels.PageInfo
				Nodes      []steampipemodels.TagWithCommits
			} `graphql:"refs(refPrefix: \"refs/tags/\", first: $pageSize, after: $cursor)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}
	variables := map[string]interface{}{
		"owner":    githubv4.String(owner),
		"repo":     githubv4.String(repo),
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
	}
	appendTagColumnIncludes(&variables, tagCols())
	repoFullName := formRepositoryFullName(owner, repo)
	var values []models.Resource
	for {
		err := client.Query(ctx, &query, variables)
		if err != nil {
			return nil, err
		}
		for _, tag := range query.Repository.Refs.Nodes {
			id := fmt.Sprintf("%s/%s/%s", owner, repo, tag.Name)
			value := models.Resource{
				ID:   id,
				Name: tag.Name,
				Description: model.TagDescription{
					RepositoryFullName: repoFullName,
					Name:               tag.Name,
					TaggerDate:         tag.Target.Tag.Tagger.Date,
					TaggerName:         tag.Target.Tag.Tagger.Name,
					TaggerLogin:        tag.Target.Tag.Tagger.Name,
					Message:            tag.Target.Tag.Message,
					Commit:             tag.Target.Commit,
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
		if !query.Repository.Refs.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Repository.Refs.PageInfo.EndCursor)
	}
	return values, nil
}
