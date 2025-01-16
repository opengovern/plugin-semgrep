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

func GetAllTeamsMembers(ctx context.Context, githubClient model.GitHubClient, organizationName string, stream *models.StreamSender) ([]models.Resource, error) {
	client := githubClient.RestClient
	teams, err := getTeams(ctx, client)
	if err != nil {
		return nil, nil
	}
	var values []models.Resource
	for _, team := range teams {
		teamValues, err := tableGitHubTeamMemberList(ctx, githubClient, stream, team.GetOrganization().GetLogin(), team.GetSlug())
		if err != nil {
			return nil, err
		}
		values = append(values, teamValues...)
	}
	return values, nil
}

func tableGitHubTeamMemberList(ctx context.Context, githubClient model.GitHubClient, stream *models.StreamSender, org, slug string) ([]models.Resource, error) {
	client := githubClient.GraphQLClient
	var query struct {
		RateLimit    steampipemodels.RateLimit
		Organization struct {
			Team struct {
				Members struct {
					TotalCount int
					PageInfo   struct {
						EndCursor   githubv4.String
						HasNextPage bool
					}
					Edges []steampipemodels.TeamMemberWithRole
				} `graphql:"members(first: $pageSize, after: $cursor)"`
			} `graphql:"team(slug: $slug)"`
		} `graphql:"organization(login: $login)"`
	}
	variables := map[string]interface{}{
		"login":    githubv4.String(org),
		"slug":     githubv4.String(slug),
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
	}
	appendUserColumnIncludes(&variables, teamMembersCols())
	var values []models.Resource
	for {
		err := client.Query(ctx, &query, variables)
		if err != nil {
			if strings.Contains(err.Error(), "Could not resolve to an Organization with the login of") {
				return nil, nil
			}
			return nil, err
		}
		for _, member := range query.Organization.Team.Members.Edges {
			value := models.Resource{
				ID:   strconv.Itoa(member.Node.Id),
				Name: member.Node.Name,
				Description: model.TeamMembersDescription{
					User:         member.Node,
					Organization: org,
					Slug:         slug,
					Role:         member.Role,
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
		if !query.Organization.Team.Members.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Organization.Team.Members.PageInfo.EndCursor)
	}
	return values, nil
}
