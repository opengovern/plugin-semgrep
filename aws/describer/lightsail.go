package describer

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/lightsail/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func LightsailInstance(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := lightsail.NewFromConfig(cfg)

	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		instances, err := client.GetInstances(ctx, &lightsail.GetInstancesInput{
			PageToken: prevToken,
		})
		if err != nil {
			return nil, err
		}

		for _, instance := range instances.Instances {
			resource := lightsailInstanceHandle(ctx, instance)

			if stream != nil {
				if err := (*stream)(resource); err != nil {
					return nil, err
				}
			} else {
				values = append(values, resource)
			}
		}

		return instances.NextPageToken, nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}
func lightsailInstanceHandle(ctx context.Context, instance types.Instance) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *instance.Arn,
		Name:   *instance.Name,
		Description: model.LightsailInstanceDescription{
			Instance: instance,
		},
	}
	return resource
}
func GetLightsailInstance(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	instanceName := fields["name"]
	var values []Resource

	client := lightsail.NewFromConfig(cfg)
	instance, err := client.GetInstance(ctx, &lightsail.GetInstanceInput{
		InstanceName: &instanceName,
	})
	if err != nil {
		if isErr(err, "GetInstanceNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	values = append(values, lightsailInstanceHandle(ctx, *instance.Instance))
	return values, nil
}
