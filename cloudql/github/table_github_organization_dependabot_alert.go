package github

import (
	opengovernance "github.com/opengovern/og-describer-template/discovery/pkg/es"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func gitHubDependabotAlertColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "alert_number",
			Type:        proto.ColumnType_INT,
			Description: "The security alert number.",
			Transform:   transform.FromField("Description.AlertNumber"),
		},
		{
			Name:        "state",
			Type:        proto.ColumnType_STRING,
			Description: "The state of the Dependabot alert.",
			Transform:   transform.FromField("Description.State")},
		{
			Name:        "dependency_package_ecosystem",
			Type:        proto.ColumnType_STRING,
			Description: "The package's language or package management ecosystem.",
			Transform:   transform.FromField("Description.DependencyPackageEcosystem")},

		{
			Name:        "dependency_package_name",
			Type:        proto.ColumnType_STRING,
			Description: "The unique package name within its ecosystem.",
			Transform:   transform.FromField("Description.DependencyPackageName"),
		},
		{
			Name:        "dependency_manifest_path",
			Type:        proto.ColumnType_STRING,
			Description: "The unique manifestation path within the ecosystem.",
			Transform:   transform.FromField("Description.DependencyManifestPath"),
		},
		{
			Name:        "dependency_scope",
			Type:        proto.ColumnType_STRING,
			Description: "The execution scope of the vulnerable dependency.",
			Transform:   transform.FromField("Description.DependencyScope"),
		},
		{
			Name:        "security_advisory_ghsa_id",
			Type:        proto.ColumnType_STRING,
			Description: "The unique GitHub Security Advisory ID assigned to the advisory.",
			Transform:   transform.FromField("Description.SecurityAdvisoryGHSAID"),
		},
		{
			Name:        "security_advisory_cve_id",
			Type:        proto.ColumnType_STRING,
			Description: "The unique CVE ID assigned to the advisory.",
			Transform:   transform.FromField("Description.SecurityAdvisoryCVEID"),
		},
		{
			Name:        "security_advisory_summary",
			Type:        proto.ColumnType_STRING,
			Description: "A short, plain text summary of the advisory.",
			Transform:   transform.FromField("Description.SecurityAdvisorySummary"),
		},
		{
			Name:        "security_advisory_description",
			Type:        proto.ColumnType_STRING,
			Description: "A long-form Markdown-supported description of the advisory.",
			Transform:   transform.FromField("Description.SecurityAdvisoryDescription"),
		},
		{
			Name:        "security_advisory_severity",
			Type:        proto.ColumnType_STRING,
			Description: "The severity of the advisory.",
			Transform:   transform.FromField("Description.SecurityAdvisorySeverity"),
		},
		{
			Name:        "security_advisory_cvss_score",
			Type:        proto.ColumnType_DOUBLE,
			Description: "The overall CVSS score of the advisory.",
			Transform:   transform.FromField("Description.SecurityAdvisoryCVSSScore"),
		},
		{
			Name:        "security_advisory_cvss_vector_string",
			Type:        proto.ColumnType_STRING,
			Description: "The full CVSS vector string for the advisory.",
			Transform:   transform.FromField("Description.SecurityAdvisoryCVSSVector"),
		},
		{
			Name:        "security_advisory_cwes",
			Type:        proto.ColumnType_JSON,
			Description: "The associated CWEs",
			Transform:   transform.FromField("Description.SecurityAdvisoryCWEs"),
		},
		{
			Name:        "security_advisory_published_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "The time that the advisory was published.",
			Transform:   transform.FromField("Description.SecurityAdvisoryPublishedAt").NullIfZero().Transform(convertTimestamp),
		},
		{
			Name:        "security_advisory_updated_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "The time that the advisory was last modified.",
			Transform:   transform.FromField("Description.SecurityAdvisoryUpdatedAt").NullIfZero().Transform(convertTimestamp),
		},
		{
			Name:        "security_advisory_withdrawn_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "The time that the advisory was withdrawn.",
			Transform:   transform.FromField("Description.SecurityAdvisoryWithdrawnAt").NullIfZero().Transform(convertTimestamp),
		},
		{
			Name:        "url",
			Type:        proto.ColumnType_STRING,
			Description: "The REST API URL of the alert resource.",
			Transform:   transform.FromField("Description.URL")},
		{
			Name:        "html_url",
			Type:        proto.ColumnType_STRING,
			Description: "The GitHub URL of the alert resource.",
			Transform:   transform.FromField("Description.HTMLURL")},
		{
			Name:        "created_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "The time that the alert was created.",
			Transform:   transform.FromField("Description.CreatedAt").NullIfZero().Transform(convertTimestamp),
		},
		{
			Name:        "updated_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "The time that the alert was last updated.",
			Transform:   transform.FromField("Description.UpdatedAt").NullIfZero().Transform(convertTimestamp),
		},
		{
			Name:        "dismissed_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "The time that the alert was dismissed.",
			Transform:   transform.FromField("Description.DismissedAt").NullIfZero().Transform(convertTimestamp),
		},
		{
			Name:        "dismissed_reason",
			Type:        proto.ColumnType_STRING,
			Description: "The reason that the alert was dismissed.",
			Transform:   transform.FromField("Description.DismissedReason")},
		{
			Name:        "dismissed_comment",
			Type:        proto.ColumnType_STRING,
			Description: "An optional comment associated with the alert's dismissal.",
			Transform:   transform.FromField("Description.DismissedComment")},
		{
			Name:        "fixed_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "The time that the alert was no longer detected and was considered fixed.",
			Transform:   transform.FromField("Description.FixedAt").NullIfZero().Transform(convertTimestamp),
		},
	}
}

func tableGitHubOrganizationDependabotAlert() *plugin.Table {
	return &plugin.Table{
		Name:        "github_organization_dependabot_alert",
		Description: "Dependabot alerts from an organization.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListOrgAlertDependabot,
		},
		Columns: commonColumns(append(
			gitHubDependabotAlertColumns(),
			[]*plugin.Column{
				{
					Name:        "organization",
					Type:        proto.ColumnType_STRING,
					Description: "The login name of the organization.",
					Transform:   transform.FromQual("organization"),
				},
			}...,
		)),
	}
}
