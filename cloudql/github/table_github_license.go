package github

import (
	opengovernance "github.com/opengovern/og-describer-template/discovery/pkg/es"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableGitHubLicense() *plugin.Table {
	return &plugin.Table{
		Name:        "github_license",
		Description: "GitHub Licenses are common software licenses that you can associate with your repository.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListLicense,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("key"),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           opengovernance.GetLicense,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "spdx_id",
				Description: "The Software Package Data Exchange (SPDX) id of the license.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.SpdxId")},
			{
				Name:        "name",
				Description: "The name of the license.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.Name")},
			{
				Name: "url", Description: "The HTML URL of the license.",
				Type:      proto.ColumnType_STRING,
				Transform: transform.FromField("Description.Url")},

			// The body is huge and of limited value, exclude it for now
			// {Name: "body", Type: proto.ColumnType_STRING, Hydrate: tableGitHubLicenseGetData},
			{
				Name:        "conditions",
				Description: "An array of license conditions (include-copyright,disclose-source, etc).",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.Conditions")},
			{
				Name:        "description",
				Description: "The license description.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.Description")},

			{
				Name:        "featured",
				Description: "If true, the license is 'featured' in the GitHub UI.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Description.Featured")},

			{
				Name:        "hidden",
				Description: "Whether the license should be displayed in license pickers.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Description.Hidden")},

			{
				Name:        "implementation",
				Description: "Implementation instructions for the license.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.Implementation")},

			{
				Name:        "key",
				Description: "The unique key of the license.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.Key")},
			{
				Name:        "limitations",
				Description: "An array of limitations for the license (trademark-use, liability,warranty, etc).",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.Limitations")},

			{
				Name:        "permissions",
				Description: "An array of permissions for the license (private-use, commercial-use,modifications, etc).",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.Permissions")},

			{
				Name:        "nickname",
				Description: "The customary short name of the license.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.Nickname")},

			{
				Name:        "pseudo_license",
				Description: "Indicates if the license is a pseudo-license placeholder (e.g. other, no-license).",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Description.PseudoLicense")},
		}),
	}
}
