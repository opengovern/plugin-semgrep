package describer

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/securitylake"
	"github.com/opengovern/og-aws-describer/aws/model"
)

// SecurityLakeDataLake TODO: new sdk version available but a field is missing
func SecurityLakeDataLake(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := securitylake.NewFromConfig(cfg)

	var values []Resource
	lakes, err := client.ListDataLakes(ctx, &securitylake.ListDataLakesInput{})
	if err != nil {
		if isErr(err, "AccessDeniedException") {
			return nil, nil
		} else {
			return nil, err
		}
	}
	for _, lake := range lakes.DataLakes {
		if lake.DataLakeArn == nil {
			continue
		}
		resource := Resource{
			Region: describeCtx.KaytuRegion,
			Name:   *lake.DataLakeArn,
			ARN:    *lake.DataLakeArn,
			Description: model.SecurityLakeDataLakeDescription{
				DataLake: lake,
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

	return values, nil
}

func SecurityLakeSubscriber(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := securitylake.NewFromConfig(cfg)

	var values []Resource
	paginator := securitylake.NewListSubscribersPaginator(client, &securitylake.ListSubscribersInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if isErr(err, "AccessDeniedException") {
				continue
			}
			return nil, err
		}

		for _, subscriber := range page.Subscribers {
			resource := Resource{
				Region: describeCtx.KaytuRegion,
				Name:   *subscriber.SubscriberName,
				Description: model.SecurityLakeSubscriberDescription{
					Subscriber: subscriber,
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
