package github

import (
	opengovernance "github.com/opengovern/og-describer-template/discovery/pkg/es"

	"time"

	"github.com/opengovern/og-describer-template/cloudql/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubTag() *plugin.Table {
	return &plugin.Table{
		Name:        "github_tag",
		Description: "Tags for commits in the given repository.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListTag,
		},
		Columns: commonColumns([]*plugin.Column{
			{Name: "repository_full_name", Type: proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.RepositoryFullName"),
				Description: "Full name of the repository that contains the tag."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the tag.",
				Transform: transform.FromField("Description.Name")},
			{Name: "tagger_date", Type: proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Description.TaggerDate").NullIfZero().Transform(convertTimestamp),
				Description: "Date the tag was created."},
			{Name: "tagger_name", Type: proto.ColumnType_STRING, Description: "Name of user whom created the tag.",
				Transform: transform.FromField("Description.TaggerName")},
			{Name: "tagger_login", Type: proto.ColumnType_STRING, Description: "Login of user whom created the tag.",
				Transform: transform.FromField("Description.TaggerLogin")},
			{Name: "message", Type: proto.ColumnType_STRING, Description: "Message associated with the tag.",
				Transform: transform.FromField("Description.Message")},
			{Name: "commit", Type: proto.ColumnType_JSON, Description: "Commit the tag is associated with.",
				Transform: transform.FromField("Description.Commit")},
		}),
	}
}

// tagRow is a struct to flatten returned information.
type tagRow struct {
	Name        string
	TaggerDate  time.Time
	TaggerName  string
	TaggerLogin string
	Message     string
	Commit      models.BaseCommit
}
