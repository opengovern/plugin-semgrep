package cohere

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	opengovernance "github.com/opengovern/og-describer-cohereai/pkg/sdk/es"

)

func tableCohereDatasets(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cohereai_datasets",
		Description: "Cohere Ai list of datasets.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListDataset,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    opengovernance.GetDataset,
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.ID"),
				Description: "ID of Dataset."},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.Name"),
				Description: "Name of the dataset."},
			{
				Name:        "created_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Description.CreatedAt").Transform(transform.UnixToTimestamp),
				Description: "Timestamp of when the dataset was created."},
			{
				Name:        "updated_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Description.UpdatedAt").Transform(transform.UnixToTimestamp),
				Description: "Timestamp of when the dataset was updated"},
			// Other columns
			{
				Name:        "dataset_type",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.DatasetType"),
				Description: "Type of the dataset."},
			{
				Name: 	  "validation_status",
				Type: 	  proto.ColumnType_STRING,
				Transform: transform.FromField("Description.ValidationStatus"),
				Description: "Validation status of the dataset."},
			{
			
				Name: "schema",
				Type: proto.ColumnType_STRING,
				Transform: transform.FromField("Description.Schema"),
				Description: "Schema of the dataset.",
			},
			{
				Name: "required_fields",
				Type: proto.ColumnType_JSON,
				Transform: transform.FromField("Description.RequiredFields"),
				Description: "List of required fields.",
			},
			{
				Name: "preserve_fields",
				Type: proto.ColumnType_JSON,
				Transform: transform.FromField("Description.PreserveFields"),
				Description: "List of preserved fields.",
			},
			{
				Name: "dataset_parts",
				Type: proto.ColumnType_JSON,
				Transform: transform.FromField("Description.DatasetParts"),
				Description: "List of dataset parts.",
			},
			{
				Name: "parse_info",
				Type: proto.ColumnType_JSON,
				Transform: transform.FromField("Description.ParseInfo"),
				Description: "Parse info of the dataset.",
			},

			
		}),
	}
}
