package github

import (
	opengovernance "github.com/opengovern/og-describer-template/discovery/pkg/es"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func gitHubRepositoryDeploymentColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "repository_full_name", Type: proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.RepoFullName"),
			Description: "The full name of the repository (login/repo-name)."},
		{Name: "id", Type: proto.ColumnType_INT,
			Transform:   transform.FromField("Description.Id"),
			Description: "The ID of the deployment."},
		{Name: "node_id", Type: proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.NodeId"),
			Description: "The node ID of the deployment."},
		{Name: "commit_sha", Type: proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.CommitSha"),
			Description: "SHA of the commit the deployment is using."},
		{Name: "created_at", Type: proto.ColumnType_TIMESTAMP,
			Transform:   transform.FromField("Description.CreatedAt").NullIfZero().Transform(convertTimestamp),
			Description: "Timestamp when the deployment was created."},
		{Name: "creator", Type: proto.ColumnType_JSON,
			Transform:   transform.FromField("Description.Creator"),
			Description: "The deployment creator."},
		{Name: "description", Type: proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.Description"),
			Description: "The description of the deployment."},
		{Name: "environment", Type: proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.Environment"),
			Description: "The name of the environment to which the deployment was made."},
		{Name: "latest_environment", Type: proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.LatestEnvironment"),
			Description: "The name of the latest environment to which the deployment was made."},
		{Name: "latest_status", Type: proto.ColumnType_JSON,
			Transform:   transform.FromField("Description.LatestStatus"),
			Description: "The latest status of the deployment."},
		{Name: "original_environment", Type: proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.OriginalEnvironment"),
			Description: "The original environment to which this deployment was made."},
		{Name: "payload", Type: proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.Payload"),
			Description: "Extra information that a deployment system might need."},
		{Name: "ref", Type: proto.ColumnType_JSON,
			Transform:   transform.FromField("Description.Ref"),
			Description: "Identifies the Ref of the deployment, if the deployment was created by ref."},
		{Name: "state", Type: proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.State"),
			Description: "The current state of the deployment."},
		{Name: "task", Type: proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.Task"),
			Description: "The deployment task."},
		{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP,
			Transform:   transform.FromField("Description.UpdatedAt").NullIfZero().Transform(convertTimestamp),
			Description: "Timestamp when the deployment was last updated."},
	}
}

func tableGitHubRepositoryDeployment() *plugin.Table {
	return &plugin.Table{
		Name:        "github_repository_deployment",
		Description: "GitHub Deployments are releases of the app/service/etc to an environment.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListRepoDeployment,
		},
		Columns: commonColumns(gitHubRepositoryDeploymentColumns()),
	}
}
