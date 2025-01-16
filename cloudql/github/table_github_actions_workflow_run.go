package github

import (
	opengovernance "github.com/opengovern/og-describer-template/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableGitHubActionsRepositoryWorkflowRun() *plugin.Table {
	return &plugin.Table{
		Name:        "github_actions_workflow_run",
		Description: "WorkflowRun represents a repository action workflow run",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListWorkflowRun,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"id"}),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           opengovernance.GetWorkflowRun,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Description.ID"),
				Description: "Unique identifier of the workflow run.",
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.Name"),
				Description: "Name of the workflow run.",
			},
			{
				Name:        "head_branch",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.HeadBranch"),
				Description: "Branch for which the workflow was run.",
			},
			{
				Name:        "head_sha",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.HeadSHA"),
				Description: "SHA of the head commit of the workflow run.",
			},
			{
				Name:        "status",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.Status"),
				Description: "Status of the workflow run (e.g., queued, in_progress, completed).",
			},
			{
				Name:        "conclusion",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.Conclusion"),
				Description: "Conclusion of the workflow run (e.g., success, failure, neutral).",
			},
			{
				Name:        "html_url",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.HTMLURL"),
				Description: "HTML URL of the workflow run.",
			},
			{
				Name:        "workflow_id",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Description.WorkflowID"),
				Description: "Unique identifier of the workflow.",
			},
			{
				Name:        "run_number",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Description.RunNumber"),
				Description: "Run number of the workflow.",
			},
			{
				Name:        "event",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.Event"),
				Description: "Event that triggered the workflow run.",
			},
			{
				Name:        "created_at",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.CreatedAt"),
				Description: "Timestamp when the workflow run was created.",
			},
			{
				Name:        "updated_at",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.UpdatedAt"),
				Description: "Timestamp when the workflow run was last updated.",
			},
			{
				Name:        "run_attempt",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Description.RunAttempt"),
				Description: "Attempt number of the workflow run.",
			},
			{
				Name:        "run_started_at",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.RunStartedAt"),
				Description: "Timestamp when the workflow run started.",
			},
			{
				Name:        "actor",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.Actor"),
				Description: "Details of the actor who triggered the workflow run.",
			},
			{
				Name:        "head_commit",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.HeadCommit"),
				Description: "Details of the head commit associated with the workflow run.",
			},
			{
				Name:        "repository",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.Repository"),
				Description: "Details of the repository where the workflow run was executed.",
			},
			{
				Name:        "head_repository",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.HeadRepository"),
				Description: "Details of the head repository associated with the workflow run.",
			},
			{
				Name:        "referenced_workflows",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.ReferencedWorkflows"),
				Description: "Referenced workflows in the workflow run.",
			},
			{
				Name:        "artifact_count",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Description.ArtifactCount"),
				Description: "Number of artifacts generated in the workflow run.",
			},
			{
				Name:        "artifacts",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.Artifacts"),
				Description: "Details of the artifacts generated in the workflow run.",
			},
		}),
	}
}
