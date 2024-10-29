package describer

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	dms "github.com/aws/aws-sdk-go-v2/service/databasemigrationservice"
	"github.com/aws/aws-sdk-go-v2/service/databasemigrationservice/types"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func DMSReplicationInstance(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := dms.NewFromConfig(cfg)

	paginator := dms.NewDescribeReplicationInstancesPaginator(client,
		&dms.DescribeReplicationInstancesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range page.ReplicationInstances {
			resource, err := dMSReplicationInstanceHandle(ctx, cfg, item)
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
func dMSReplicationInstanceHandle(ctx context.Context, cfg aws.Config, item types.ReplicationInstance) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := dms.NewFromConfig(cfg)
	tags, err := client.ListTagsForResource(ctx, &dms.ListTagsForResourceInput{
		ResourceArn: item.ReplicationInstanceArn,
	})
	if err != nil {
		if isErr(err, "ListTagsForResourceNoFound") || isErr(err, "InvalidParameterValue") {
			return Resource{}, nil
		}
		return Resource{}, err
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *item.ReplicationInstanceArn,
		Name:   *item.ReplicationInstanceIdentifier,
		Description: model.DMSReplicationInstanceDescription{
			ReplicationInstance: item,
			Tags:                tags.TagList,
		},
	}
	return resource, nil
}
func GetDMSReplicationInstance(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	replicationInstanceArn := fields["arn"]
	client := dms.NewFromConfig(cfg)

	out, err := client.DescribeReplicationInstances(ctx, &dms.DescribeReplicationInstancesInput{})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, item := range out.ReplicationInstances {

		if *item.ReplicationInstanceArn != replicationInstanceArn {
			continue
		}

		resource, err := dMSReplicationInstanceHandle(ctx, cfg, item)
		if err != nil {
			return nil, err
		}
		emptyResource := Resource{}
		if err == nil && resource == emptyResource {
			return nil, nil
		}

		values = append(values, resource)
	}
	return values, nil
}

func DMSEndpoint(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := dms.NewFromConfig(cfg)

	paginator := dms.NewDescribeEndpointsPaginator(client,
		&dms.DescribeEndpointsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range page.Endpoints {
			resource, err := dMSEndpointHandle(ctx, cfg, item)
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

func dMSEndpointHandle(ctx context.Context, cfg aws.Config, item types.Endpoint) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := dms.NewFromConfig(cfg)
	tags, err := client.ListTagsForResource(ctx, &dms.ListTagsForResourceInput{
		ResourceArn: item.EndpointArn,
	})
	if err != nil {
		if isErr(err, "ListTagsForResourceNoFound") || isErr(err, "InvalidParameterValue") {
			return Resource{}, nil
		}
		return Resource{}, err
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *item.EndpointArn,
		Name:   *item.EndpointIdentifier,
		Description: model.DMSEndpointDescription{
			Endpoint: item,
			Tags:     tags.TagList,
		},
	}
	return resource, nil
}

func DMSReplicationTask(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := dms.NewFromConfig(cfg)

	paginator := dms.NewDescribeReplicationTasksPaginator(client,
		&dms.DescribeReplicationTasksInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range page.ReplicationTasks {
			resource, err := dMSReplicationTaskHandle(ctx, cfg, item)
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

func dMSReplicationTaskHandle(ctx context.Context, cfg aws.Config, item types.ReplicationTask) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := dms.NewFromConfig(cfg)
	tags, err := client.ListTagsForResource(ctx, &dms.ListTagsForResourceInput{
		ResourceArn: item.ReplicationTaskArn,
	})
	if err != nil {
		if isErr(err, "ListTagsForResourceNoFound") || isErr(err, "InvalidParameterValue") {
			return Resource{}, nil
		}
		return Resource{}, err
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *item.ReplicationTaskArn,
		ID:     *item.ReplicationTaskIdentifier,
		Name:   *item.ReplicationTaskIdentifier,
		Description: model.DMSReplicationTaskDescription{
			ReplicationTask: item,
			Tags:            tags.TagList,
		},
	}
	return resource, nil
}
