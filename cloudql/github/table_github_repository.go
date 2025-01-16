package github

import (
	opengovernance "github.com/opengovern/og-describer-template/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func sharedRepositoryColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "GitHubRepoID",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromField("Description.GitHubRepoID"),
			Description: "Unique identifier of the GitHub repository.",
		},
		{
			Name:        "node_id",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.NodeID"),
			Description: "Node ID of the repository.",
		},
		{
			Name:        "name",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.Name"),
			Description: "Name of the repository.",
		},
		{
			Name:        "name_with_owner",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.NameWithOwner"),
			Description: "Full name of the repository including the owner.",
		},
		{
			Name:        "description",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.Description"),
			Description: "Description of the repository.",
		},
		{
			Name:        "created_at",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.CreatedAt"),
			Description: "Timestamp when the repository was created.",
		},
		{
			Name:        "updated_at",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.UpdatedAt"),
			Description: "Timestamp when the repository was last updated.",
		},
		{
			Name:        "pushed_at",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.PushedAt"),
			Description: "Timestamp when the repository was last pushed.",
		},
		{
			Name:        "is_active",
			Type:        proto.ColumnType_BOOL,
			Transform:   transform.FromField("Description.IsActive"),
			Description: "Indicates if the repository is active.",
		},
		{
			Name:        "is_empty",
			Type:        proto.ColumnType_BOOL,
			Transform:   transform.FromField("Description.IsEmpty"),
			Description: "Indicates if the repository is empty.",
		},
		{
			Name:        "is_fork",
			Type:        proto.ColumnType_BOOL,
			Transform:   transform.FromField("Description.IsFork"),
			Description: "Indicates if the repository is a fork.",
		},
		{
			Name:        "is_security_policy_enabled",
			Type:        proto.ColumnType_BOOL,
			Transform:   transform.FromField("Description.IsSecurityPolicyEnabled"),
			Description: "Indicates if the repository has a security policy enabled.",
		},
		{
			Name:        "owner",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("Description.Owner"),
			Description: "Owner details of the repository.",
		},
		{
			Name:        "homepage_url",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.HomepageURL"),
			Description: "Homepage URL of the repository.",
		},
		{
			Name:        "license_info",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("Description.LicenseInfo"),
			Description: "License information of the repository.",
		},
		{
			Name:        "topics",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("Description.Topics"),
			Description: "List of topics associated with the repository.",
		},
		{
			Name:        "visibility",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.Visibility"),
			Description: "Visibility status of the repository.",
		},
		{
			Name:        "default_branch_ref",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("Description.DefaultBranchRef"),
			Description: "Details of the default branch of the repository.",
		},
		{
			Name:        "permissions",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("Description.Permissions"),
			Description: "Permissions associated with the repository.",
		},
		{
			Name:        "organization",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("Description.Organization"),
			Description: "Organization details of the repository.",
		},
		{
			Name:        "parent",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("Description.Parent"),
			Description: "Parent repository details if the repository is forked.",
		},
		{
			Name:        "source",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("Description.Source"),
			Description: "Source repository details if the repository is forked.",
		},
		{
			Name:        "primary_language",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.PrimaryLanguage"),
			Description: "Primary language used in the repository.",
		},
		{
			Name:        "languages",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("Description.Languages"),
			Description: "Languages used in the repository along with their usage statistics.",
		},
		{
			Name:        "repo_settings",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("Description.RepositorySettings"),
			Description: "Settings of the repository.",
		},
		{
			Name:        "security_settings",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("Description.SecuritySettings"),
			Description: "Security settings of the repository.",
		},
		{
			Name:        "repo_urls",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("Description.RepoURLs"),
			Description: "Repository URLs for different purposes (e.g., clone URLs).",
		},
		{
			Name:        "metrics",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("Description.Metrics"),
			Description: "Metrics and statistics of the repository.",
		},
	}
}

func tableGitHubRepository() *plugin.Table {
	return &plugin.Table{
		Name:        "github_repository",
		Description: "GitHub Repositories contain all of your project's files and each file's revision history.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListRepository,
		},
		Columns: commonColumns(sharedRepositoryColumns()),
	}
}
