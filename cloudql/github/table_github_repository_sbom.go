package github

import (
	opengovernance "github.com/opengovern/og-describer-template/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubRepositorySbom() *plugin.Table {
	return &plugin.Table{
		Name:        "github_repository_sbom",
		Description: "Get the software bill of materials (SBOM) for a repository.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListRepoSBOM,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "repository_full_name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.RepositoryFullName"),
				Description: "The full name of the repository (login/repo-name).",
			},
			{
				Name:        "spdx_id",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.SPDXID"),
				Description: "The SPDX identifier for the SPDX document.",
			},
			{
				Name:        "spdx_version",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.SPDXVersion"),
				Description: "The version of the SPDX specification that this document conforms to.",
			},
			{
				Name:        "creation_info",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.CreationInfo"),
				Description: "It represents when the SBOM was created and who created it.",
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the SPDX document.",
				Transform:   transform.FromField("Description.Name")},
			{
				Name:        "data_license",
				Type:        proto.ColumnType_STRING,
				Description: "The license under which the SPDX document is licensed.",
				Transform:   transform.FromField("Description.DataLicense")},
			{
				Name:        "document_describes",
				Type:        proto.ColumnType_JSON,
				Description: "The name of the repository that the SPDX document describes.",
				Transform:   transform.FromField("Description.DocumentDescribes")},
			{
				Name:        "document_namespace",
				Type:        proto.ColumnType_STRING,
				Description: "The namespace for the SPDX document.",
				Transform:   transform.FromField("Description.DocumentNamespace")},
			{
				Name:        "packages",
				Type:        proto.ColumnType_JSON,
				Description: "Array of packages in SPDX format.",
				Transform:   transform.FromField("Description.Packages")},
		}),
	}
}
