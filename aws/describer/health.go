package describer

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/health/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/health"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func HealthEvent(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := health.NewFromConfig(cfg)
	paginator := health.NewDescribeEventsPaginator(client, &health.DescribeEventsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, event := range page.Events {
			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    *event.Arn,
				Description: model.HealthEventDescription{
					Event: event,
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

func HealthAffectedEntity(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := health.NewFromConfig(cfg)
	paginator := health.NewDescribeEventsPaginator(client, &health.DescribeEventsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, event := range page.Events {
			entitiesPaginator := health.NewDescribeAffectedEntitiesPaginator(client, &health.DescribeAffectedEntitiesInput{
				Filter: &types.EntityFilter{
					EventArns: []string{*aws.String(*event.Arn)},
				},
			})
			for entitiesPaginator.HasMorePages() {
				entitiesPage, err := entitiesPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}
				for _, entity := range entitiesPage.Entities {
					resource := Resource{
						Region: describeCtx.KaytuRegion,
						ARN:    *event.Arn,
						Description: model.HealthAffectedEntityDescription{
							Entity: entity,
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
		}
	}
	return values, nil
}
