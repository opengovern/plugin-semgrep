package describer

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func TaggingResources(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := resourcegroupstaggingapi.NewFromConfig(cfg)
	paginator := resourcegroupstaggingapi.NewGetResourcesPaginator(client, &resourcegroupstaggingapi.GetResourcesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, v := range page.ResourceTagMappingList {
			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ID:     *v.ResourceARN,
				Name:   *v.ResourceARN,
				Description: model.TaggingResourcesDescription{
					TagMapping: v,
				},
			}
			if stream != nil {
				if err := (*stream)(resource); err != nil {
					return nil, err
				}
			} else {
				values = append(values, resource)
			}
		}
	}
	return values, nil
}
