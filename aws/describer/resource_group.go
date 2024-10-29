package describer

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroups"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func ResourceGroups(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := resourcegroups.NewFromConfig(cfg)
	paginator := resourcegroups.NewListGroupsPaginator(client, &resourcegroups.ListGroupsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.GroupIdentifiers {
			resources, err := client.ListGroupResources(ctx, &resourcegroups.ListGroupResourcesInput{Group: v.GroupArn})
			if err != nil {
				return nil, err
			}

			tags, err := client.GetTags(ctx, &resourcegroups.GetTagsInput{Arn: v.GroupArn})
			if err != nil {
				return nil, err
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    *v.GroupArn,
				Name:   *v.GroupName,
				Description: model.ResourceGroupsGroupDescription{
					GroupIdentifier: v,
					Resources:       resources.Resources,
					Tags:            tags,
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
