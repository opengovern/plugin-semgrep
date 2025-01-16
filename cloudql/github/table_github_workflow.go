package github

import (
	opengovernance "github.com/opengovern/og-describer-template/discovery/pkg/es"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubWorkflow() *plugin.Table {
	return &plugin.Table{
		Name:        "github_workflow",
		Description: "GitHub Workflows bundle project files for download by users.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListWorkflow,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"repository_full_name", "id"}),
			Hydrate:    opengovernance.GetWorkflow,
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{
				Name: "repository_full_name", Type: proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.RepositoryFullName"),
				Description: "Full name of the repository that contains the workflow."},
			{
				Name: "name", Type: proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.Name"),
				Description: "The name of the workflow."},
			{Name: "id", Type: proto.ColumnType_INT,
				Transform:   transform.FromField("Description.ID"),
				Description: "Unique ID of the workflow."},
			{Name: "path", Type: proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.Path"),
				Description: "Path of the workflow."},

			// Other columns
			{Name: "badge_url", Type: proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.BadgeURL"),
				Description: "Badge URL for the workflow."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Description.CreatedAt").NullIfZero().Transform(convertTimestamp),
				Description: "Time when the workflow was created."},
			{Name: "html_url", Type: proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.HTMLURL"),
				Description: "HTML URL for the workflow."},
			{Name: "node_id", Type: proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.NodeID"),
				Description: "Node where GitHub stores this data internally."},
			{Name: "state", Type: proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.State"),
				Description: "State of the workflow."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Description.UpdatedAt").NullIfZero().Transform(convertTimestamp),
				Description: "Time when the workflow was updated."},
			{Name: "url", Type: proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.URL"),
				Description: "URL of the workflow."},
			{Name: "workflow_file_content", Type: proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.WorkFlowFileContent"),
				Description: "Content of github workflow file in text format."},
			{Name: "workflow_file_content_json",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.WorkFlowFileContentJson"),
				Description: "Content of github workflow file in the JSON format."},
			{Name: "pipeline", Type: proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.Pipeline"),
				Description: "Github workflow in the generic pipeline entity format to be used across CI/CD platforms."},
		}),
	}
}

type FileContent struct {
	Repository string
	FilePath   string
	Content    string
}
