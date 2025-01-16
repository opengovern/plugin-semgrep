package github

import (
	opengovernance "github.com/opengovern/og-describer-template/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubBranch() *plugin.Table {
	return &plugin.Table{
		Name:        "github_branch",
		Description: "Branches in the given repository.",
		List: &plugin.ListConfig{
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           opengovernance.ListBranch,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "repository_full_name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.RepoFullName"),
				Description: "Full name of the repository that contains the branch.",
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.Name"),
				Description: "Name of the branch."},
			{
				Name: "commit", Type: proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.Commit"),
				Description: "Latest commit on the branch."},
			{
				Name:        "protected",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Description.Protected"),
				Description: "If true, the branch is protected."},
			{
				Name:        "branch_protection_rule",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.BranchProtectionRule"),
				Description: "Branch protection rule if protected."},
		}),
	}
}
