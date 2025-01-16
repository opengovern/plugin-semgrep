package describers

import (
	"context"
	"fmt"
	"strings"

	"github.com/opengovern/og-describer-template/discovery/pkg/models"
	model "github.com/opengovern/og-describer-template/discovery/provider"
	"github.com/shurcooL/githubv4"
	steampipemodels "github.com/turbot/steampipe-plugin-github/github/models"
)

type CollaboratorEdge struct {
	Permission githubv4.RepositoryPermission     `graphql:"permission @include(if:$includeOCPermission)" json:"permission"`
	Node       steampipemodels.CollaboratorLogin `graphql:"node @include(if:$includeOCNode)" json:"node"`
}

func GetAllOrganizationsCollaborators(ctx context.Context, githubClient model.GitHubClient, organizationName string, stream *models.StreamSender) ([]models.Resource, error) {
	var values []models.Resource
	orgValues, err := GetOrganizationCollaborators(ctx, githubClient, stream, organizationName)
	if err != nil {
		return nil, err
	}
	values = append(values, orgValues...)
	return values, nil
}

func GetOrganizationCollaborators(ctx context.Context, githubClient model.GitHubClient, stream *models.StreamSender, org string) ([]models.Resource, error) {
	client := githubClient.GraphQLClient
	affiliation := githubv4.CollaboratorAffiliationAll
	var query struct {
		RateLimit    steampipemodels.RateLimit
		Organization struct {
			URL          githubv4.String
			Login        githubv4.String
			Repositories struct {
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage githubv4.Boolean
				}
				Nodes []struct {
					Name          githubv4.String
					Collaborators struct {
						Edges []CollaboratorEdge
					} `graphql:"collaborators(affiliation: $affiliation)"`
				}
			} `graphql:"repositories(first: $pageSize, after: $cursor)"`
		} `graphql:"organization(login: $login)"`
	}
	variables := map[string]interface{}{
		"login":       githubv4.String(org),
		"pageSize":    githubv4.Int(orgCollaboratorsPageSize),
		"cursor":      (*githubv4.String)(nil),
		"affiliation": affiliation,
	}
	appendOrgCollaboratorColumnIncludes(&variables, orgCollaboratorsCols())
	var values []models.Resource
	for {
		err := client.Query(ctx, &query, variables)
		if err != nil {
			if strings.Contains(err.Error(), "Could not resolve to an Organization with the login of") {
				return nil, nil
			}
			return nil, err
		}
		for _, node := range query.Organization.Repositories.Nodes {
			repoFullName := formRepositoryFullName(org, string(node.Name))
			for _, collaborator := range node.Collaborators.Edges {
				id := fmt.Sprintf("%s/%s", repoFullName, collaborator.Node.Login)
				value := models.Resource{
					ID:   id,
					Name: repoFullName,
					Description: model.OrgCollaboratorsDescription{
						Organization:   org,
						Affiliation:    "ALL",
						RepositoryName: node.Name,
						Permission:     collaborator.Permission,
						UserLogin:      collaborator.Node,
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
		}
		if !query.Organization.Repositories.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Organization.Repositories.PageInfo.EndCursor)
	}
	return values, nil
}
