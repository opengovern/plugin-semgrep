package github

import (
	opengovernance "github.com/opengovern/og-describer-template/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubContainerPackage() *plugin.Table {
	return &plugin.Table{
		Name: "github_container_package",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListContainerPackage,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"id"}),
			Hydrate:    opengovernance.GetContainerPackage,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Description.ID"),
				Description: "Unique identifier for the package."},
			{
				Name:        "digest",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.Digest"),
				Description: "Digest of the package."},
			{
				Name:        "created_at",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.CreatedAt"),
				Description: "Timestamp when the package was created."},
			{
				Name:        "updated_at",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.UpdatedAt"),
				Description: "Timestamp when the package was last updated."},
			{
				Name:        "package_url",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.PackageURL"),
				Description: "HTML URL for the package."},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.Name"),
				Description: "Name of the package."},
			{
				Name:        "media_type",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.MediaType"),
				Description: "Media type of the package."},
			{
				Name:        "total_size",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Description.TotalSize"),
				Description: "Total size of the package."},
			{
				Name:        "metadata",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.Metadata"),
				Description: "Metadata of the package."},
			{
				Name:        "manifest",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.Manifest"),
				Description: "Manifest of the package."},
		}),
	}
}
