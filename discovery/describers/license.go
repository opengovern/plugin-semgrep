package describers

import (
	"context"

	"github.com/opengovern/og-describer-template/discovery/pkg/models"
	model "github.com/opengovern/og-describer-template/discovery/provider"
	"github.com/shurcooL/githubv4"
	steampipemodels "github.com/turbot/steampipe-plugin-github/github/models"
)

func GetLicenseList(ctx context.Context, githubClient model.GitHubClient, organizationName string, stream *models.StreamSender) ([]models.Resource, error) {
	client := githubClient.GraphQLClient

	var query struct {
		RateLimit  steampipemodels.RateLimit
		Repository struct {
			LicenseInfo steampipemodels.License
		} `graphql:"repository(owner: $owner, name: $repoName)"`
	}
	repositories, err := getRepositories(ctx, githubClient.RestClient, organizationName)
	if err != nil {
		return nil, nil
	}

	var values []models.Resource
	for _, r := range repositories {
		variables := map[string]interface{}{
			"owner":    githubv4.String(organizationName),
			"repoName": githubv4.String(r.GetName()),
		}
		appendLicenseColumnIncludes(&variables, licenseCols())
		err := client.Query(ctx, &query, variables)
		if err != nil {
			return nil, err
		}
		value := models.Resource{
			ID:   query.Repository.LicenseInfo.Key,
			Name: query.Repository.LicenseInfo.Name,
			Description: model.LicenseDescription{
				Key:            query.Repository.LicenseInfo.Key,
				Name:           query.Repository.LicenseInfo.Name,
				Nickname:       query.Repository.LicenseInfo.Nickname,
				SpdxId:         query.Repository.LicenseInfo.SpdxId,
				Url:            query.Repository.LicenseInfo.Url,
				Body:           query.Repository.LicenseInfo.Body,
				Conditions:     query.Repository.LicenseInfo.Conditions,
				Description:    query.Repository.LicenseInfo.Description,
				Featured:       query.Repository.LicenseInfo.Featured,
				Hidden:         query.Repository.LicenseInfo.Hidden,
				Implementation: query.Repository.LicenseInfo.Implementation,
				Limitations:    query.Repository.LicenseInfo.Limitations,
				Permissions:    query.Repository.LicenseInfo.Permissions,
				PseudoLicense:  query.Repository.LicenseInfo.PseudoLicense,
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

	return values, nil
}

func GetLicense(ctx context.Context, githubClient model.GitHubClient, organizationName string, repositoryName string, resourceID string, stream *models.StreamSender) (*models.Resource, error) {
	client := githubClient.GraphQLClient

	variables := map[string]interface{}{
		"key": githubv4.String(resourceID),
	}

	var query struct {
		RateLimit steampipemodels.RateLimit
		License   steampipemodels.License `graphql:"license(key: $key)"`
	}
	appendLicenseColumnIncludes(&variables, licenseCols())

	err := client.Query(ctx, &query, variables)
	if err != nil {
		return nil, err
	}
	value := models.Resource{
		ID:   query.License.Key,
		Name: query.License.Name,
		Description: model.LicenseDescription{
			Key:            query.License.Key,
			Name:           query.License.Name,
			Nickname:       query.License.Nickname,
			SpdxId:         query.License.SpdxId,
			Url:            query.License.Url,
			Body:           query.License.Body,
			Conditions:     query.License.Conditions,
			Description:    query.License.Description,
			Featured:       query.License.Featured,
			Hidden:         query.License.Hidden,
			Implementation: query.License.Implementation,
			Limitations:    query.License.Limitations,
			Permissions:    query.License.Permissions,
			PseudoLicense:  query.License.PseudoLicense,
		},
	}
	if stream != nil {
		if err := (*stream)(value); err != nil {
			return nil, err
		}
	}

	return &value, nil
}
