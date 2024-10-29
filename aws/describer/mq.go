package describer

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/mq/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/mq"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func MQBroker(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := mq.NewFromConfig(cfg)
	paginator := mq.NewListBrokersPaginator(client, &mq.ListBrokersInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.BrokerSummaries {
			resource, err := mQBrokerHandle(ctx, cfg, v)
			if err != nil {
				return nil, err
			}
			emptyResource := Resource{}
			if err == nil && resource == emptyResource {
				continue
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
func mQBrokerHandle(ctx context.Context, cfg aws.Config, v types.BrokerSummary) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := mq.NewFromConfig(cfg)
	tags, err := client.ListTags(ctx, &mq.ListTagsInput{
		ResourceArn: v.BrokerArn,
	})
	if err != nil {
		if isErr(err, "ListTagsNotFound") || isErr(err, "InvalidParameterValue") {
			return Resource{}, nil
		}
		return Resource{}, err
	}

	brokerDescription, err := client.DescribeBroker(ctx, &mq.DescribeBrokerInput{
		BrokerId: v.BrokerId,
	})
	if err != nil {
		return Resource{}, err
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *v.BrokerArn,
		Name:   *v.BrokerName,
		Description: model.MQBrokerDescription{
			BrokerDescription: brokerDescription,
			Tags:              tags.Tags,
		},
	}
	return resource, nil
}
func GetMQBroker(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	brokerId := fields["id"]
	client := mq.NewFromConfig(cfg)

	brokers, err := client.ListBrokers(ctx, &mq.ListBrokersInput{})
	if err != nil {
		if isErr(err, "ListBrokersNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	for _, broker := range brokers.BrokerSummaries {
		if *broker.BrokerId != brokerId {
			continue
		}

		resource, err := mQBrokerHandle(ctx, cfg, broker)
		if err != nil {
			return nil, err
		}
		emptyResource := Resource{}
		if err == nil && resource == emptyResource {
			continue
		}

		values = append(values, resource)
	}
	return values, nil
}
