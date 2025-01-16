package github

import (
	opengovernance "github.com/opengovern/og-describer-template/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubActionsArtifact() *plugin.Table {
	return &plugin.Table{
		Name:        "github_actions_artifact",
		Description: "Artifacts allow you to share data between jobs in a workflow and store data once that workflow has completed.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListArtifact,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"repository_full_name", "id"}),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           opengovernance.GetArtifact,
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{
				Name:        "repository_full_name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.RepoFullName"),
				Description: "Full name of the repository that contains the artifact.",
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the artifact.",
				Transform:   transform.FromField("Description.Name"),
			},
			{
				Name:        "id",
				Type:        proto.ColumnType_INT,
				Description: "Unique ID of the artifact.",
				Transform:   transform.FromField("Description.ID"),
			},
			{
				Name:        "size_in_bytes",
				Type:        proto.ColumnType_INT,
				Description: "Size of the artifact in bytes.",
				Transform:   transform.FromField("Description.SizeInBytes"),
			},

			// Other columns
			{
				Name:        "archive_download_url",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.ArchiveDownloadURL"),
				Description: "Archive download URL for the artifact.",
			},
			{
				Name:        "created_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Description.CreatedAt").NullIfZero().Transform(convertTimestamp),
				Description: "Time when the artifact was created.",
			},
			{
				Name:        "expired",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Description.Expired"),
				Description: "It defines whether the artifact is expires or not.",
			},
			{
				Name:        "expires_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Description.ExpiresAt").NullIfZero().Transform(convertTimestamp),
				Description: "Time when the artifact expires.",
			},
			{
				Name:        "node_id",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.NodeID"),
				Description: "Node where GitHub stores this data internally.",
			},
		}),
	}
}
