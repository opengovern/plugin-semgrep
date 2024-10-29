package describer

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/resourceexplorer2"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func ResourceExplorerIndex(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := resourceexplorer2.NewFromConfig(cfg)
	paginator := resourceexplorer2.NewListIndexesPaginator(client, &resourceexplorer2.ListIndexesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Indexes {
			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    *v.Arn,
				Name:   *v.Arn,
				Description: model.ResourceExplorer2IndexDescription{
					Index: v,
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

func ResourceExplorer2SupportedResourceType(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := resourceexplorer2.NewFromConfig(cfg)
	paginator := resourceexplorer2.NewListSupportedResourceTypesPaginator(client, &resourceexplorer2.ListSupportedResourceTypesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.ResourceTypes {
			resource := Resource{
				Region: describeCtx.KaytuRegion,
				Name:   *v.ResourceType,
				Description: model.ResourceExplorer2SupportedResourceTypeDescription{
					SupportedResourceType: v,
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
