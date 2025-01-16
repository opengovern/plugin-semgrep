package github

import (
	"context"
	"encoding/json"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/memoize"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func commonColumns(c []*plugin.Column) []*plugin.Column {
	return append(c, []*plugin.Column{
		{
			Name:        "platform_integration_id",
			Type:        proto.ColumnType_STRING,
			Description: "The Platform Integration ID in which the resource is located.",
			Transform:   transform.FromField("IntegrationID"),
		},
		{
			Name:        "platform_resource_id",
			Type:        proto.ColumnType_STRING,
			Description: "The unique ID of the resource in opengovernance.",
			Transform:   transform.FromField("PlatformID"),
		},
		{
			Name:        "platform_metadata",
			Type:        proto.ColumnType_JSON,
			Description: "The metadata of the resource",
			Transform:   transform.FromField("Metadata").Transform(marshalJSON),
		},
		{
			Name:        "platform_resource_description",
			Type:        proto.ColumnType_JSON,
			Description: "The full model description of the resource",
			Transform:   transform.FromField("Description").Transform(marshalJSON),
		},
	}...)
}

func marshalJSON(_ context.Context, d *transform.TransformData) (interface{}, error) {
	b, err := json.Marshal(d.Value)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

// if the caching is required other than per connection, build a cache key for the call and use it in Memoize.
var getLoginIdMemoized = plugin.HydrateFunc(getLoginIdUncached).Memoize(memoize.WithCacheKeyFunction(getLoginIdCacheKey))

// declare a wrapper hydrate function to call the memoized function
// - this is required when a memoized function is used for a column definition
func getLoginId(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return getLoginIdMemoized(ctx, d, h)
}

// Build a cache key for the call to getLoginIdCacheKey.
func getLoginIdCacheKey(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	key := "getLoginId"
	return key, nil
}

func getLoginIdUncached(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	client := connectV4(ctx, d)

	var query struct {
		Viewer struct {
			Login githubv4.String
			ID    githubv4.ID
		}
	}
	err := client.Query(ctx, &query, nil)
	if err != nil {
		plugin.Logger(ctx).Error("getLoginIdUncached", "api_error", err)
		return nil, err
	}

	return query.Viewer.ID, nil
}
