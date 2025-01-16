package describers

import (
	"context"
	"strconv"
	"strings"

	"github.com/opengovern/og-describer-template/discovery/pkg/models"
	model "github.com/opengovern/og-describer-template/discovery/provider"
	"github.com/shurcooL/githubv4"
	steampipemodels "github.com/turbot/steampipe-plugin-github/github/models"
)

type memberWithRole struct {
	HasTwoFactorEnabled *bool
	Role                *string
	Node                steampipemodels.User
}

func GetAllMembers(ctx context.Context, githubClient model.GitHubClient, organizationName string, stream *models.StreamSender) ([]models.Resource, error) {
	var values []models.Resource
	orgValues, err := GetOrganizationMembers(ctx, githubClient, stream, organizationName)
	if err != nil {
		return nil, err
	}
	values = append(values, orgValues...)

	return values, nil
}

func GetOrganizationMembers(ctx context.Context, githubClient model.GitHubClient, stream *models.StreamSender, org string) ([]models.Resource, error) {
	client := githubClient.GraphQLClient
	var query struct {
		RateLimit    steampipemodels.RateLimit
		Organization struct {
			Login           string
			MembersWithRole struct {
				Edges    []memberWithRole
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage bool
				}
			} `graphql:"membersWithRole(first: $pageSize, after: $cursor)"`
		} `graphql:"organization(login: $login)"`
	}
	variables := map[string]interface{}{
		"login":    githubv4.String(org),
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
	}
	appendUserColumnIncludes(&variables, orgMembersCols())
	var values []models.Resource
	for {
		err := client.Query(ctx, &query, variables)
		if err != nil {
			if strings.Contains(err.Error(), "Could not resolve to an Organization with the login of") {
				return nil, nil
			}
			return nil, err
		}
		for _, member := range query.Organization.MembersWithRole.Edges {
			value := models.Resource{
				ID:   strconv.Itoa(member.Node.Id),
				Name: member.Node.Name,
				Description: model.OrgMembersDescription{
					User:                member.Node,
					Organization:        org,
					HasTwoFactorEnabled: member.HasTwoFactorEnabled,
					Role:                member.Role,
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
		if !query.Organization.MembersWithRole.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Organization.MembersWithRole.PageInfo.EndCursor)
	}
	return values, nil
}
