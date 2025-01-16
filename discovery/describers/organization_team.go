package describers

import (
	"context"
	"strconv"

	"github.com/opengovern/og-describer-template/discovery/pkg/models"
	model "github.com/opengovern/og-describer-template/discovery/provider"
	"github.com/shurcooL/githubv4"
	steampipemodels "github.com/turbot/steampipe-plugin-github/github/models"
)

func GetOrganizationTeamList(ctx context.Context, githubClient model.GitHubClient, organizationName string, stream *models.StreamSender) ([]models.Resource, error) {
	client := githubClient.GraphQLClient
	var query struct {
		RateLimit steampipemodels.RateLimit
		Viewer    struct {
			Organization struct {
				Login string
				Teams struct {
					PageInfo steampipemodels.PageInfo
					Nodes    []steampipemodels.TeamWithCounts
				} `graphql:"teams(first: $pageSize, after: $cursor)"`
			} `graphql:"organization(login: $orgName)"`
		}
	}
	variables := map[string]interface{}{
		"orgName":  githubv4.String(organizationName),
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
	}
	appendTeamColumnIncludes(&variables, teamCols())
	var values []models.Resource
	var teams []steampipemodels.TeamWithCounts
	err := client.Query(ctx, &query, variables)
	if err != nil {
		return nil, err
	}
	teams = append(teams, query.Viewer.Organization.Teams.Nodes...)
	if query.Viewer.Organization.Teams.PageInfo.HasNextPage {
		ts, err := getAdditionalTeams(ctx, client, query.Viewer.Organization.Login, query.Viewer.Organization.Teams.PageInfo.EndCursor)
		if err != nil {
			return nil, err
		}
		teams = append(teams, ts...)
	}

	for _, team := range teams {
		value := models.Resource{
			ID:   strconv.Itoa(team.Id),
			Name: team.Name,
			Description: model.TeamDescription{
				Organization: team.Organization.Name,
				Slug:         team.Slug,
				Name:         team.Name,
				ID:           team.Id,
				NodeID:       team.NodeId,
				Description:  team.Description,
				CreatedAt:    team.CreatedAt,
				UpdatedAt:    team.UpdatedAt,
				CombinedSlug: team.CombinedSlug,
				ParentTeam: struct {
					Id     int
					NodeId string
					Name   string
					Slug   string
				}{Id: team.ParentTeam.Id, NodeId: team.ParentTeam.NodeId, Name: team.ParentTeam.Name, Slug: team.ParentTeam.Slug},
				Privacy:                team.Privacy,
				AncestorsTotalCount:    team.Ancestors.TotalCount,
				ChildTeamsTotalCount:   team.ChildTeams.TotalCount,
				DiscussionsTotalCount:  team.Discussions.TotalCount,
				InvitationsTotalCount:  team.Invitations.TotalCount,
				MembersTotalCount:      team.Members.TotalCount,
				ProjectsV2TotalCount:   team.ProjectsV2.TotalCount,
				RepositoriesTotalCount: team.Repositories.TotalCount,
				URL:                    team.Url,
				AvatarURL:              team.AvatarUrl,
				DiscussionsURL:         team.DiscussionsUrl,
				EditTeamURL:            team.EditTeamUrl,
				MembersURL:             team.MembersUrl,
				NewTeamURL:             team.NewTeamUrl,
				RepositoriesURL:        team.RepositoriesUrl,
				TeamsURL:               team.TeamsUrl,
				CanAdminister:          team.CanAdminister,
				CanSubscribe:           team.CanSubscribe,
				Subscription:           team.Subscription,
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

func getAdditionalTeams(ctx context.Context, client *githubv4.Client, org string, initialCursor githubv4.String) ([]steampipemodels.TeamWithCounts, error) {
	var query struct {
		RateLimit    steampipemodels.RateLimit
		Organization struct {
			Teams struct {
				PageInfo steampipemodels.PageInfo
				Nodes    []steampipemodels.TeamWithCounts
			} `graphql:"teams(first: $pageSize, after: $cursor)"`
		} `graphql:"organization(login: $login)"`
	}
	variables := map[string]interface{}{
		"pageSize": githubv4.Int(100),
		"cursor":   githubv4.NewString(initialCursor),
		"login":    githubv4.String(org),
	}
	appendTeamColumnIncludes(&variables, teamCols())
	var ts []steampipemodels.TeamWithCounts
	for {
		err := client.Query(ctx, &query, variables)
		if err != nil {
			return nil, err
		}
		ts = append(ts, query.Organization.Teams.Nodes...)
		if !query.Organization.Teams.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Organization.Teams.PageInfo.EndCursor)
	}
	return ts, nil
}
