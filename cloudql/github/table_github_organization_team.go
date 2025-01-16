package github

import (
	opengovernance "github.com/opengovern/og-describer-template/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func gitHubTeamColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "organization", Type: proto.ColumnType_STRING, Description: "The organization the team is associated with.",
			Transform: transform.FromField("Description.Organization")},
		{Name: "slug",
			Transform: transform.FromField("Description.Slug"),
			Type:      proto.ColumnType_STRING, Description: "The team slug name."},
		{Name: "name",
			Transform: transform.FromField("Description.Name"),
			Type:      proto.ColumnType_STRING, Description: "The name of the team."},
		{Name: "id", Type: proto.ColumnType_INT, Description: "The ID of the team.",
			Transform: transform.FromField("Description.ID")},
		{Name: "node_id", Type: proto.ColumnType_STRING, Description: "The node id of the team.",
			Transform: transform.FromField("Description.NodeID")},
		{Name: "description", Type: proto.ColumnType_STRING, Description: "The description of the team.",
			Transform: transform.FromField("Description.Description")},
		{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when team was created.",
			Transform: transform.FromField("Description.CreatedAt").NullIfZero().Transform(convertTimestamp)},
		{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when team was last updated.",
			Transform: transform.FromField("Description.UpdatedAt").NullIfZero().Transform(convertTimestamp)},
		{Name: "combined_slug", Type: proto.ColumnType_STRING, Description: "The slug corresponding to the organization and the team.",
			Transform: transform.FromField("Description.CombinedSlug")},
		{Name: "parent_team", Type: proto.ColumnType_JSON, Description: "The teams parent team.",
			Transform: transform.FromField("Description.ParentTeam")},
		{Name: "privacy", Type: proto.ColumnType_STRING, Description: "The privacy setting of the team (VISIBLE or SECRET).",
			Transform: transform.FromField("Description.Privacy")},
		{Name: "ancestors_total_count", Type: proto.ColumnType_INT, Description: "Count of ancestors this team has.",
			Transform: transform.FromField("Description.AncestorsTotalCount")},
		{Name: "child_teams_total_count", Type: proto.ColumnType_INT, Description: "Count of children teams this team has.",
			Transform: transform.FromField("Description.ChildTeamsTotalCount")},
		{Name: "discussions_total_count", Type: proto.ColumnType_INT, Description: "Count of team discussions.",
			Transform: transform.FromField("Description.DiscussionsTotalCount")},
		{Name: "invitations_total_count", Type: proto.ColumnType_INT, Description: "Count of outstanding team member invitations for the team.",
			Transform: transform.FromField("Description.InvitationsTotalCount")},
		{Name: "members_total_count", Type: proto.ColumnType_INT, Description: "Count of team members.",
			Transform: transform.FromField("Description.MembersTotalCount")},
		{Name: "projects_v2_total_count", Type: proto.ColumnType_INT, Description: "Count of the teams v2 projects.",
			Transform: transform.FromField("Description.ProjectsV2TotalCount")},
		{Name: "repositories_total_count", Type: proto.ColumnType_INT, Description: "Count of repositories the team has.",
			Transform: transform.FromField("Description.RepositoriesTotalCount")},
		{Name: "url", Type: proto.ColumnType_STRING, Description: "URL for the team page in GitHub.",
			Transform: transform.FromField("Description.URL")},
		{Name: "avatar_url", Type: proto.ColumnType_STRING, Description: "URL for teams avatar.",
			Transform: transform.FromField("Description.AvatarURL")},
		{Name: "discussions_url", Type: proto.ColumnType_STRING, Description: "URL for team discussions.",
			Transform: transform.FromField("Description.DiscussionsURL")},
		{Name: "edit_team_url", Type: proto.ColumnType_STRING, Description: "URL for editing this team.",
			Transform: transform.FromField("Description.EditTeamURL")},
		{Name: "members_url", Type: proto.ColumnType_STRING, Description: "URL for team members.",
			Transform: transform.FromField("Description.MembersURL")},
		{Name: "new_team_url", Type: proto.ColumnType_STRING, Description: "The HTTP URL creating a new team.",
			Transform: transform.FromField("Description.NewTeamURL")},
		{Name: "repositories_url", Type: proto.ColumnType_STRING, Description: "URL for team repositories.",
			Transform: transform.FromField("Description.RepositoriesURL")},
		{Name: "teams_url", Type: proto.ColumnType_STRING, Description: "URL for this team's teams.",
			Transform: transform.FromField("Description.TeamsURL")},
		{Name: "can_administer", Type: proto.ColumnType_BOOL, Description: "If true, current user can administer the team.",
			Transform: transform.FromField("Description.CanAdminister")},
		{Name: "can_subscribe", Type: proto.ColumnType_BOOL, Description: "If true, current user can subscribe to the team.",
			Transform: transform.FromField("Description.CanSubscribe")},
		{Name: "subscription", Type: proto.ColumnType_STRING, Description: "Subscription status of the current user to the team.",
			Transform: transform.FromField("Description.Subscription")},
	}
}

func tableGitHubOrganizationTeam() *plugin.Table {
	return &plugin.Table{
		Name:        "github_organization_team",
		Description: "GitHub Teams in a given organization. GitHub Teams are groups of organization members that reflect your company or group's structure with cascading access permissions and mentions.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListTeamMembers,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"organization", "slug"}),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           opengovernance.GetTeamMembers,
		},
		Columns: commonColumns(gitHubTeamColumns()),
	}
}
