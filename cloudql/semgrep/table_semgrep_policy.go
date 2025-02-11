package semgrep

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-semgrep/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableSemGrepPolicy(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "semgrep_policy",
		Description: "SemGrep policies information.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListPolicy,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    opengovernance.GetPolicy,
		},
		Columns: integrationColumns([]*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.ID"),
				Description: "The unique identifier of the policy.",
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.Name"),
				Description: "The name of the policy.",
			},
			{
				Name:        "slug",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.Slug"),
				Description: "The slug of the policy.",
			},
			{
				Name:        "product_type",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.ProductType"),
				Description: "The product type associated with the policy.",
			},
			{
				Name:        "is_default",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Description.IsDefault"),
				Description: "Indicates whether this policy is the default one.",
			},
		}),
	}
}
