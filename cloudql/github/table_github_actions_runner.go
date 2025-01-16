package github

import (
	opengovernance "github.com/opengovern/og-describer-template/discovery/pkg/es"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubActionsRepositoryRunner() *plugin.Table {
	return &plugin.Table{
		Name:        "github_actions_runner",
		Description: "The runner is the application that runs a job from a GitHub Actions workflow",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListRunner,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"repository_full_name", "id"}),
			Hydrate:    opengovernance.GetRunner,
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{
				Name:        "repository_full_name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.RepoFullName"),
				Description: "Full name of the repository that contains the runners.",
			},
			{
				Name:        "id",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Description.ID"),
				Description: "The unique identifier of the runner.",
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.Name"),
				Description: "The name of the runner.",
			},
			{
				Name:        "os",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.OS"),
				Description: "The operating system of the runner."},
			{
				Name:        "status",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.Status"),
				Description: "The status of the runner.",
			},
			{
				Name:        "busy",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Description.Busy"),
				Description: "Indicates whether the runner is currently in use or not.",
			},
			{
				Name:        "labels",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.Labels"),
				Description: "Labels represents a collection of labels attached to each runner.",
			},
		}),
	}
}
