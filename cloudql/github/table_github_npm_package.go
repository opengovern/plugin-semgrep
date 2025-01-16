package github

import (
	opengovernance "github.com/opengovern/og-describer-template/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableGitHubNPMPackage() *plugin.Table {
	return &plugin.Table{
		Name: "github_npm_package",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListPackage,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"id"}),
			Hydrate:    opengovernance.GetPackage,
		},
		Columns: commonColumns([]*plugin.Column{
			// Basic details columns
			{Name: "id", Type: proto.ColumnType_INT, Description: "Unique identifier for the package."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the package."},
			{Name: "package_type", Type: proto.ColumnType_STRING, Description: "Type of the package."},
			{Name: "version_count", Type: proto.ColumnType_INT, Description: "Number of versions of the package."},
			{Name: "visibility", Type: proto.ColumnType_STRING, Description: "Visibility of the package (e.g., public or private)."},
			{Name: "url", Type: proto.ColumnType_STRING, Description: "URL to access the package."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the package was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the package was last updated."},
			{Name: "html_url", Type: proto.ColumnType_STRING, Description: "HTML URL for the package."},
			// Nested structure columns
			{Name: "owner", Type: proto.ColumnType_JSON, Description: "Owner details of the package."},
			{Name: "repository", Type: proto.ColumnType_JSON, Description: "Repository details associated with the package."},
		}),
	}
}
