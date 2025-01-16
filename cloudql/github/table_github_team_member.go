package github

import (
	opengovernance "github.com/opengovern/og-describer-template/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubTeamMember() *plugin.Table {
	return &plugin.Table{
		Name:        "github_team_member",
		Description: "GitHub members for a given team. GitHub Users are user accounts in GitHub.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListTeamMembers,
		},
		Columns: commonColumns(gitHubTeamMemberColumns()),
	}
}

func gitHubTeamMemberColumns() []*plugin.Column {
	cols := []*plugin.Column{
		{Name: "organization", Type: proto.ColumnType_STRING, Description: "The organization the team is associated with.",
			Transform: transform.FromField("Description.Organization")},
		{Name: "slug", Type: proto.ColumnType_STRING, Description: "The team slug name.",
			Transform: transform.FromField("Description.Slug")},
		{Name: "role", Type: proto.ColumnType_STRING, Description: "The team member's role (MEMBER, MAINTAINER).",
			Transform: transform.FromField("Description.Role")},
	}

	cols = append(cols, sharedUserColumns()...)
	return cols
}
