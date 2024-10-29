package describer

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func SNSSubscription(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := sns.NewFromConfig(cfg)
	paginator := sns.NewListSubscriptionsPaginator(client, &sns.ListSubscriptionsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Subscriptions {
			if v.SubscriptionArn != nil && *v.SubscriptionArn == "PendingConfirmation" {
				continue
			}
			output, err := client.GetSubscriptionAttributes(ctx, &sns.GetSubscriptionAttributesInput{
				SubscriptionArn: v.SubscriptionArn,
			})
			if err != nil {
				if !isErr(err, "NotFound") {
					return nil, err
				}

				output = &sns.GetSubscriptionAttributesOutput{}
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    *v.SubscriptionArn,
				Name:   nameFromArn(*v.SubscriptionArn),
				Description: model.SNSSubscriptionDescription{
					Subscription: v,
					Attributes:   output.Attributes,
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

func SNSTopic(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := sns.NewFromConfig(cfg)
	paginator := sns.NewListTopicsPaginator(client, &sns.ListTopicsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Topics {
			resource, err := sNSTopicHandle(ctx, cfg, v)
			emptyResource := Resource{}
			if err == nil && resource == emptyResource {
				return nil, nil
			}
			if err != nil {
				return nil, err
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
func sNSTopicHandle(ctx context.Context, cfg aws.Config, v types.Topic) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := sns.NewFromConfig(cfg)

	output, err := client.GetTopicAttributes(ctx, &sns.GetTopicAttributesInput{
		TopicArn: v.TopicArn,
	})
	if err != nil {
		if isErr(err, "GetTopicAttributesNotFound") || isErr(err, "invalidParameterValue") {
			return Resource{}, nil
		}
		return Resource{}, err
	}

	tOutput, err := client.ListTagsForResource(ctx, &sns.ListTagsForResourceInput{
		ResourceArn: v.TopicArn,
	})
	if err != nil {
		if isErr(err, "ListTagsForResourceNotFound") || isErr(err, "invalidParameterValue") {
			return Resource{}, nil
		}
		return Resource{}, err
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *v.TopicArn,
		Name:   nameFromArn(*v.TopicArn),
		Description: model.SNSTopicDescription{
			Attributes: output.Attributes,
			Tags:       tOutput.Tags,
		},
	}
	return resource, nil
}
func GetSNSTopic(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	arn := fields["arn"]
	var values []Resource
	client := sns.NewFromConfig(cfg)
	list, err := client.ListTopics(ctx, &sns.ListTopicsInput{})
	if err != nil {
		return nil, err
	}

	for _, v := range list.Topics {
		if v.TopicArn != &arn {
			continue
		}
		resource, err := sNSTopicHandle(ctx, cfg, v)
		emptyResource := Resource{}
		if err == nil && resource == emptyResource {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}

		values = append(values, resource)
	}
	return values, nil
}
