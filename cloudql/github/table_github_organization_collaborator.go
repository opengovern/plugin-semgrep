package github

import (
	opengovernance "github.com/opengovern/og-describer-template/discovery/pkg/es"

	"github.com/opengovern/og-describer-template/cloudql/github/models"

	"github.com/shurcooL/githubv4"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func gitHubOrganizationCollaborators() []*plugin.Column {
	tableCols := []*plugin.Column{
		{
			Name:        "organization",
			Type:        proto.ColumnType_STRING,
			Description: "The organization the member is associated with.",
			Transform:   transform.FromField("Description.Organization")},
		{
			Name:        "affiliation",
			Type:        proto.ColumnType_STRING,
			Description: "Affiliation filter - valid values 'ALL' (default), 'OUTSIDE', 'DIRECT'.",
			Transform:   transform.FromField("Description.Affiliation")},
		{
			Name:        "repository_name",
			Type:        proto.ColumnType_STRING,
			Description: "The name of the repository",
			Transform:   transform.FromField("Description.RepositoryName")},

		{
			Name:        "permission",
			Type:        proto.ColumnType_STRING,
			Description: "The permission the collaborator has on the repository.",
			Transform:   transform.FromField("Description.Permission")},

		{
			Name:        "user_login",
			Type:        proto.ColumnType_JSON,
			Description: "The login details of the collaborator.",
			Transform:   transform.FromField("Description.UserLogin")},
	}

	return tableCols
}

type OrgCollaborators struct {
	RepositoryName githubv4.String
	Permission     githubv4.RepositoryPermission
	Node           models.CollaboratorLogin
}

type CollaboratorEdge struct {
	Permission githubv4.RepositoryPermission `graphql:"permission @include(if:$includeOCPermission)" json:"permission"`
	Node       models.CollaboratorLogin      `graphql:"node @include(if:$includeOCNode)" json:"node"`
}

func tableGitHubOrganizationCollaborator() *plugin.Table {
	return &plugin.Table{
		Name:        "github_organization_collaborator",
		Description: "GitHub members for a given organization. GitHub Users are user accounts in GitHub.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListOrgCollaborators,
		},
		Columns: commonColumns(gitHubOrganizationCollaborators()),
	}
}
