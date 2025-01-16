package template

import (
	opengovernance "github.com/opengovern/og-describer-template/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubArtifactDockerFile() *plugin.Table {
	return &plugin.Table{
		Name: "github_artifact_dockerfile",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListArtifactDockerFile,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"sha"}),
			Hydrate:    opengovernance.GetArtifactDockerFile,
		},
		Columns: commonColumns([]*plugin.Column{
			// Basic details columns
			{
				Name:        "sha",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.Sha"),
				Description: "SHA hash of the Dockerfile."},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.Name"),
				Description: "Name of the Dockerfile."},
		
			{
				Name:        "last_updated_at",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.LastUpdatedAt"),
				Description: "Last time the dockerfile updated"},
			
			{
				Name:        "html_url",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.HTMLURL"),
				Description: "HTML URL where the Dockerfile can be accessed."},
			
			{
				Name:        "dockerfile_content",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.DockerfileContent"),
				Description: "Dockerfile content."},
			
			{
				Name:        "repository",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.Repository"),
				Description: "Repository metadata associated with the Dockerfile."},
			{
				Name:        "images",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.Images"),
				Description: ""},
		}),
	}
}
