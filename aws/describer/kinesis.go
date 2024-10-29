package describer

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kinesisanalyticsv2"
	"github.com/aws/aws-sdk-go-v2/service/kinesisvideo"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func KinesisStream(ctx context.Context, cfg aws.Config, streamS *StreamSender) ([]Resource, error) {
	client := kinesis.NewFromConfig(cfg)

	var values []Resource
	var lastStreamName *string = nil
	for {
		streams, err := client.ListStreams(ctx, &kinesis.ListStreamsInput{
			ExclusiveStartStreamName: lastStreamName,
		})
		if err != nil {
			if isErr(err, "ResourceNotFoundException") || isErr(err, "InvalidParameter") {
				return nil, nil
			}
			return nil, err
		}
		for _, streamName := range streams.StreamNames {
			streamName := streamName
			stream, err := client.DescribeStream(ctx, &kinesis.DescribeStreamInput{
				StreamName: &streamName,
			})
			if err != nil {
				if !isErr(err, "ResourceNotFoundException") && !isErr(err, "InvalidParameter") {
					return nil, err
				}
				continue
			}

			streamSummery, err := client.DescribeStreamSummary(ctx, &kinesis.DescribeStreamSummaryInput{
				StreamName: &streamName,
			})
			if err != nil {
				if !isErr(err, "ResourceNotFoundException") && !isErr(err, "InvalidParameter") {
					return nil, err
				}
				continue
			}

			tags, err := client.ListTagsForStream(ctx, &kinesis.ListTagsForStreamInput{
				StreamName: &streamName,
			})
			if err != nil {
				if !isErr(err, "ResourceNotFoundException") && !isErr(err, "InvalidParameter") {
					return nil, err
				}
				tags = &kinesis.ListTagsForStreamOutput{}
			}

			resource := kinesisStreamHandle(ctx, stream, streamSummery, tags)
			if streamS != nil {
				if err := (*streamS)(resource); err != nil {
					return nil, err
				}
			} else {
				values = append(values, resource)
			}
		}

		if streams.HasMoreStreams == nil || !*streams.HasMoreStreams {
			break
		}

		lastStreamName = &streams.StreamNames[len(streams.StreamNames)-1]
	}

	return values, nil
}
func kinesisStreamHandle(ctx context.Context, stream *kinesis.DescribeStreamOutput, streamSummery *kinesis.DescribeStreamSummaryOutput, tags *kinesis.ListTagsForStreamOutput) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *stream.StreamDescription.StreamARN,
		Name:   *stream.StreamDescription.StreamName,
		Description: model.KinesisStreamDescription{
			Stream:             *stream.StreamDescription,
			DescriptionSummary: *streamSummery.StreamDescriptionSummary,
			Tags:               tags.Tags,
		},
	}
	return resource
}
func GetKinesisStream(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	streamName := fields["name"]

	var values []Resource
	client := kinesis.NewFromConfig(cfg)
	stream, err := client.DescribeStream(ctx, &kinesis.DescribeStreamInput{
		StreamName: &streamName,
	})
	if err != nil {
		if !isErr(err, "ResourceNotFoundException") && !isErr(err, "InvalidParameter") {
			return nil, err
		}
		return nil, nil
	}

	streamSummery, err := client.DescribeStreamSummary(ctx, &kinesis.DescribeStreamSummaryInput{
		StreamName: &streamName,
	})
	if err != nil {
		if !isErr(err, "ResourceNotFoundException") && !isErr(err, "InvalidParameter") {
			return nil, err
		}
		return nil, nil
	}

	tags, err := client.ListTagsForStream(ctx, &kinesis.ListTagsForStreamInput{
		StreamName: &streamName,
	})
	if err != nil {
		if !isErr(err, "ResourceNotFoundException") && !isErr(err, "InvalidParameter") {
			return nil, err
		}
		tags = &kinesis.ListTagsForStreamOutput{}
	}

	values = append(values, kinesisStreamHandle(ctx, stream, streamSummery, tags))
	return values, nil
}

func KinesisConsumer(ctx context.Context, cfg aws.Config, streamS *StreamSender) ([]Resource, error) {
	client := kinesis.NewFromConfig(cfg)
	var values []Resource
	err := PaginateRetrieveAll(func(startName *string) (*string, error) {
		streams, err := client.ListStreams(ctx, &kinesis.ListStreamsInput{
			ExclusiveStartStreamName: startName,
		})
		if err != nil {
			if isErr(err, "ResourceNotFoundException") || isErr(err, "InvalidParameter") {
				return nil, nil
			}
			return nil, err
		}
		var lastStreamName *string = nil
		for _, streamName := range streams.StreamNames {
			lastStreamName = &streamName
			stream, err := client.DescribeStream(ctx, &kinesis.DescribeStreamInput{
				StreamName: &streamName,
			})
			if err != nil {
				if !isErr(err, "ResourceNotFoundException") && !isErr(err, "InvalidParameter") {
					return nil, err
				}
				continue
			}
			err = PaginateRetrieveAll(func(prevToken *string) (*string, error) {
				consumers, err := client.ListStreamConsumers(ctx, &kinesis.ListStreamConsumersInput{
					StreamARN: stream.StreamDescription.StreamARN,
				})
				if err != nil {
					if isErr(err, "ResourceNotFoundException") || isErr(err, "InvalidParameter") {
						return nil, nil
					}
					return nil, err
				}
				for _, consumer := range consumers.Consumers {

					resource := kinesisConsumerHandle(ctx, stream, consumer)
					if streamS != nil {
						if err := (*streamS)(resource); err != nil {
							return nil, err
						}
					} else {
						values = append(values, resource)
					}
				}
				return consumers.NextToken, nil
			})
		}
		if streams.HasMoreStreams == nil || !*streams.HasMoreStreams {
			return nil, nil
		}
		return lastStreamName, nil
	})

	if err != nil {
		return nil, err
	}

	return values, nil
}
func kinesisConsumerHandle(ctx context.Context, stream *kinesis.DescribeStreamOutput, consumer types.Consumer) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *consumer.ConsumerARN,
		Name:   *consumer.ConsumerName,
		Description: model.KinesisConsumerDescription{
			StreamARN: *stream.StreamDescription.StreamARN,
			Consumer:  consumer,
		},
	}
	return resource
}
func GetKinesisConsumer(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	streamName := fields["name"]
	var values []Resource
	client := kinesis.NewFromConfig(cfg)
	stream, err := client.DescribeStream(ctx, &kinesis.DescribeStreamInput{
		StreamName: &streamName,
	})
	if err != nil {
		if isErr(err, "DescribeStreamNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	consumers, err := client.ListStreamConsumers(ctx, &kinesis.ListStreamConsumersInput{
		StreamARN: stream.StreamDescription.StreamARN,
	})
	if err != nil {
		if isErr(err, "ListStreamConsumersNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	for _, consumer := range consumers.Consumers {
		resource := kinesisConsumerHandle(ctx, stream, consumer)
		values = append(values, resource)
	}
	return values, nil
}

func KinesisVideoStream(ctx context.Context, cfg aws.Config, streamS *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := kinesisvideo.NewFromConfig(cfg)
	paginator := kinesisvideo.NewListStreamsPaginator(client, &kinesisvideo.ListStreamsInput{})
	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, stream := range page.StreamInfoList {
			tags, err := client.ListTagsForStream(ctx, &kinesisvideo.ListTagsForStreamInput{
				StreamARN: stream.StreamARN,
			})
			if err != nil {
				tags = &kinesisvideo.ListTagsForStreamOutput{}
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    *stream.StreamARN,
				Name:   *stream.StreamName,
				Description: model.KinesisVideoStreamDescription{
					Stream: stream,
					Tags:   tags.Tags,
				},
			}
			if streamS != nil {
				if err := (*streamS)(resource); err != nil {
					return nil, err
				}
			} else {
				values = append(values, resource)
			}
		}
	}

	return values, nil
}

func KinesisAnalyticsV2Application(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := kinesisanalyticsv2.NewFromConfig(cfg)
	var values []Resource

	err := PaginateRetrieveAll(func(prevToken *string) (*string, error) {
		applications, err := client.ListApplications(ctx, &kinesisanalyticsv2.ListApplicationsInput{
			NextToken: prevToken,
		})
		if err != nil {
			return nil, err
		}
		for _, application := range applications.ApplicationSummaries {
			application := application
			description, err := client.DescribeApplication(ctx, &kinesisanalyticsv2.DescribeApplicationInput{
				ApplicationName: application.ApplicationName,
			})
			if err != nil {
				if !isErr(err, "ResourceNotFoundException") && !isErr(err, "InvalidParameter") {
					return nil, err
				}
				continue
			}

			tags, err := client.ListTagsForResource(ctx, &kinesisanalyticsv2.ListTagsForResourceInput{
				ResourceARN: description.ApplicationDetail.ApplicationARN,
			})

			resource := kinesisAnalyticsV2ApplicationHandle(ctx, description, tags)
			if stream != nil {
				if err := (*stream)(resource); err != nil {
					return nil, err
				}
			} else {
				values = append(values, resource)
			}

		}

		return applications.NextToken, nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}
func kinesisAnalyticsV2ApplicationHandle(ctx context.Context, description *kinesisanalyticsv2.DescribeApplicationOutput, tags *kinesisanalyticsv2.ListTagsForResourceOutput) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *description.ApplicationDetail.ApplicationARN,
		Name:   *description.ApplicationDetail.ApplicationName,
		Description: model.KinesisAnalyticsV2ApplicationDescription{
			Application: *description.ApplicationDetail,
			Tags:        tags.Tags,
		},
	}
	return resource
}
func GetKinesisAnalyticsV2Application(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	applicationName := fields["name"]

	var values []Resource
	client := kinesisanalyticsv2.NewFromConfig(cfg)
	description, err := client.DescribeApplication(ctx, &kinesisanalyticsv2.DescribeApplicationInput{
		ApplicationName: &applicationName,
	})
	if err != nil {
		if isErr(err, "DescribeApplicationNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	tags, err := client.ListTagsForResource(ctx, &kinesisanalyticsv2.ListTagsForResourceInput{
		ResourceARN: description.ApplicationDetail.ApplicationARN,
	})
	if err != nil {
		if isErr(err, "ListTagsForResourceNotFound") || isErr(err, "invalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	resource := kinesisAnalyticsV2ApplicationHandle(ctx, description, tags)
	values = append(values, resource)
	return values, nil
}
