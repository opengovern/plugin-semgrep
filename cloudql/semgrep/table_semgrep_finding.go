package semgrep

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-semgrep/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableSemGrepFinding(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "semgrep_finding",
		Description: "SemGrep findings details.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListFinding,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    opengovernance.GetFinding,
		},
		Columns: integrationColumns([]*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Description.ID"),
				Description: "The unique identifier of the finding.",
			},
			{
				Name:        "ref",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.Ref"),
				Description: "Reference to the finding.",
			},
			{
				Name:        "first_seen_scan_id",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Description.FirstSeenScanID"),
				Description: "The ID of the scan where this finding was first seen.",
			},
			{
				Name:        "syntactic_id",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.SyntacticID"),
				Description: "The syntactic ID of the finding.",
			},
			{
				Name:        "match_based_id",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.MatchBasedID"),
				Description: "The match-based ID of the finding.",
			},
			{
				Name:        "external_ticket",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.ExternalTicket"),
				Description: "Details of the external ticket associated with the finding.",
			},
			{
				Name:        "repository",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.Repository"),
				Description: "Repository where the finding was detected.",
			},
			{
				Name:        "line_of_code_url",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.LineOfCodeURL"),
				Description: "URL pointing to the affected line of code.",
			},
			{
				Name:        "triage_state",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.TriageState"),
				Description: "The triage state of the finding.",
			},
			{
				Name:        "state",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.State"),
				Description: "Current state of the finding.",
			},
			{
				Name:        "status",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.Status"),
				Description: "The status of the finding.",
			},
			{
				Name:        "severity",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.Severity"),
				Description: "Severity level of the finding.",
			},
			{
				Name:        "confidence",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.Confidence"),
				Description: "Confidence level of the finding.",
			},
			{
				Name:        "categories",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.Categories"),
				Description: "Categories associated with the finding.",
			},
			{
				Name:        "created_at",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.CreatedAt"),
				Description: "Timestamp when the finding was created.",
			},
			{
				Name:        "relevant_since",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.RelevantSince"),
				Description: "Timestamp since when the finding has been relevant.",
			},
			{
				Name:        "rule_name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.RuleName"),
				Description: "Name of the rule that generated the finding.",
			},
			{
				Name:        "rule_message",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.RuleMessage"),
				Description: "Message from the rule that generated the finding.",
			},
			{
				Name:        "location",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.Location"),
				Description: "Location details of the affected code.",
			},
			{
				Name:        "sourcing_policy",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.SourcingPolicy"),
				Description: "Sourcing policy associated with the finding.",
			},
			{
				Name:        "rule",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.Rule"),
				Description: "The rule that generated the finding.",
			},
			{
				Name:        "assistant",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.Assistant"),
				Description: "AI assistant-generated guidance for the finding.",
			},
		}),
	}
}
