package semgrep

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-semgrep/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableSemGrepProject(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "semgrep_project",
		Description: "SemGrep projects information.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListProject,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    opengovernance.GetProject,
		},
		Columns: integrationColumns([]*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Description.ID"),
				Description: "The unique identifier of the project.",
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.Name"),
				Description: "The name of the project.",
			},
			{
				Name:        "url",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.URL"),
				Description: "The repository URL of the project.",
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.Tags"),
				Description: "A list of tags associated with the project.",
			},
			{
				Name:        "created_at",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.CreatedAt"),
				Description: "The timestamp when the project was created.",
			},
			{
				Name:        "latest_scan_at",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.LatestScanAt"),
				Description: "The timestamp of the latest scan.",
			},
			{
				Name:        "primary_branch",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.PrimaryBranch"),
				Description: "The primary branch of the project.",
			},
			{
				Name:        "default_branch",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.DefaultBranch"),
				Description: "The default branch of the project.",
			},
		}),
	}
}
