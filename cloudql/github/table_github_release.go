package github

import (
	opengovernance "github.com/opengovern/og-describer-template/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubRelease() *plugin.Table {
	return &plugin.Table{
		Name:        "github_release",
		Description: "GitHub Releases bundle project files for download by users.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListRelease,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"repository_full_name", "id"}),
			Hydrate:    opengovernance.GetRelease,
		},
		Columns: commonColumns([]*plugin.Column{

			// Top columns
			{Name: "repository_full_name", Type: proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.RepositoryFullName"),
				Description: "Full name of the repository that contains the release."},

			// Other columns
			{Name: "assets", Type: proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.Assets"),
				Description: "List of assets contained in the release."},
			{Name: "assets_url", Type: proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.AssetsURL"),
				Description: "Assets URL for the release."},
			{Name: "author_login", Type: proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.AssetsURL"),
				Description: "The login name of the user that created the release."},
			{Name: "body", Type: proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.Body"),
				Description: "Text describing the contents of the tag."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Description.CreatedAt").NullIfZero().Transform(convertTimestamp),
				Description: "Time when the release was created."},
			{Name: "draft", Type: proto.ColumnType_BOOL,
				Transform: transform.FromField("Description.Draft"),

				Description: "True if this is a draft (unpublished) release."},
			{Name: "html_url", Type: proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.HTMLURL"),
				Description: "HTML URL for the release."},
			{Name: "id", Type: proto.ColumnType_INT,
				Transform:   transform.FromField("Description.ID"),
				Description: "Unique ID of the release."},
			{Name: "name", Type: proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.Name"),
				Description: "The name of the release."},
			{Name: "node_id", Type: proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.NodeID"),
				Description: "Node where GitHub stores this data internally."},
			{Name: "prerelease", Type: proto.ColumnType_BOOL,
				Transform: transform.FromField("Description.Prerelease"),

				Description: "True if this is a prerelease version."},
			{Name: "published_at", Type: proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Description.PublishedAt").NullIfZero().Transform(convertTimestamp),
				Description: "Time when the release was published."},
			{Name: "tag_name", Type: proto.ColumnType_STRING,
				Transform: transform.FromField("Description.TagName"),

				Description: "The name of the tag the release is associated with."},
			{Name: "tarball_url", Type: proto.ColumnType_STRING,
				Transform: transform.FromField("Description.TarballURL"),

				Description: "Tarball URL for the release."},
			{Name: "target_commitish", Type: proto.ColumnType_STRING,
				Transform: transform.FromField("Description.TargetCommitish"),

				Description: "Specifies the commitish value that determines where the Git tag is created from. Can be any branch or commit SHA."},
			{Name: "upload_url", Type: proto.ColumnType_STRING,
				Transform: transform.FromField("Description.UploadURL"),

				Description: "Upload URL for the release."},
			{Name: "url", Type: proto.ColumnType_STRING,
				Transform: transform.FromField("Description.URL"),

				Description: "URL of the release."},
			{Name: "zipball_url", Type: proto.ColumnType_STRING,
				Transform: transform.FromField("Description.ZipballURL"),

				Description: "Zipball URL for the release."},
		}),
	}
}
