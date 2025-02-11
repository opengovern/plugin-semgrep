package semgrep

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-semgrep/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableSemGrepDeployment(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "semgrep_deployment",
		Description: "SemGrep deployments information.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListDeployment,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    opengovernance.GetDeployment,
		},
		Columns: integrationColumns([]*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Description.ID"),
				Description: "The unique identifier of the deployment.",
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.Name"),
				Description: "The name of the deployment.",
			},
			{
				Name:        "slug",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.Slug"),
				Description: "The slug of the deployment.",
			},
			{
				Name:        "findings",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.Findings"),
				Description: "",
			},
		}),
	}
}
