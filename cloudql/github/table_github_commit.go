package github

import (
	opengovernance "github.com/opengovern/og-describer-template/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableGitHubCommit() *plugin.Table {
	return &plugin.Table{
		Name:        "github_commit",
		Description: "GitHub Commits bundle project files for download by users.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListCommit,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"sha"}),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           opengovernance.GetCommit,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "sha",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.SHA"),
				Description: "Unique identifier (SHA) of the commit.",
			},
			{
				Name:        "node_id",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.NodeID"),
				Description: "",
			},
			{
				Name:        "commit_detail",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.CommitDetail"),
				Description: "Details of the commit.",
			},
			{
				Name:        "url",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.URL"),
				Description: "URL of the commit.",
			},
			{
				Name:        "html_url",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.HTMLURL"),
				Description: "URL of the commit on the repository.",
			},
			{
				Name:        "comments_url",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.CommentsURL"),
				Description: "",
			},
			{
				Name:        "author",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.Author"),
				Description: "Details of the author of the commit",
			},
			{
				Name:        "commiter",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.Committer"),
				Description: "Details of the commiter of the commit",
			},
			{
				Name:        "parents",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.Parents"),
				Description: "Parent commits of the commit",
			},
			{
				Name:        "stats",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.Stats"),
				Description: "Stats of the commit",
			},
			{
				Name:        "files",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.Files"),
				Description: "List of files changed in the commit.",
			},
		}),
	}
}
