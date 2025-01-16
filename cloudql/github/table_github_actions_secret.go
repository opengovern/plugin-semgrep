package github

import (
	opengovernance "github.com/opengovern/og-describer-template/discovery/pkg/es"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubActionsRepositorySecret() *plugin.Table {
	return &plugin.Table{
		Name:        "github_actions_secret",
		Description: "Secrets are encrypted environment variables that you create in a repository",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListSecret,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"repository_full_name", "name"}),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           opengovernance.GetSecret,
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{
				Name:        "repository_full_name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.RepoFullName"),
				Description: "Full name of the repository that contains the secrets.",
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.Name"),
				Description: "The name of the secret.",
			},
			{
				Name:        "visibility",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.Visibility"),
				Description: "The visibility of the secret.",
			},
			{
				Name:        "selected_repositories_url",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.SelectedRepositoriesURL"),
				Description: "The GitHub URL of the repository.",
			},

			// Other columns
			{
				Name:        "created_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Description.CreatedAt").NullIfZero().Transform(convertTimestamp),
				Description: "Time when the secret was created.",
			},
			{
				Name:        "updated_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Description.UpdatedAt").NullIfZero().Transform(convertTimestamp),
				Description: "Time when the secret was updated.",
			},
		}),
	}
}
