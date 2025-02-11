package semgrep

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-semgrep/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableSemGrepScan(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "semgrep_scan",
		Description: "SemGrep scan details.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListScan,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    opengovernance.GetScan,
		},
		Columns: integrationColumns([]*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.ID"),
				Description: "The unique identifier of the scan.",
			},
			{
				Name:        "deployment_id",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.DeploymentID"),
				Description: "The deployment ID associated with the scan.",
			},
			{
				Name:        "repository_id",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.RepositoryID"),
				Description: "The repository ID associated with the scan.",
			},
			{
				Name:        "branch",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.Branch"),
				Description: "The branch where the scan was performed.",
			},
			{
				Name:        "commit",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.Commit"),
				Description: "The commit hash of the scanned code.",
			},
			{
				Name:        "is_full_scan",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Description.IsFullScan"),
				Description: "Indicates if the scan was a full repository scan.",
			},
			{
				Name:        "started_at",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.StartedAt"),
				Description: "The timestamp when the scan started.",
			},
			{
				Name:        "completed_at",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.CompletedAt"),
				Description: "The timestamp when the scan completed.",
			},
			{
				Name:        "exit_code",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Description.ExitCode"),
				Description: "The exit code of the scan process.",
			},
			{
				Name:        "total_time",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("Description.TotalTime"),
				Description: "Total time taken for the scan in seconds.",
			},
			{
				Name:        "findings_counts",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.FindingsCounts"),
				Description: "The number of findings categorized by type.",
			},
			{
				Name:        "status",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.Status"),
				Description: "The current status of the scan.",
			},
		}),
	}
}
