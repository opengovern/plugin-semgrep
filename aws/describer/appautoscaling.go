package describer

import (
	"context"
	_ "github.com/aws/aws-sdk-go-v2/service/inspector/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/applicationautoscaling"
	"github.com/aws/aws-sdk-go-v2/service/applicationautoscaling/types"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func ApplicationAutoScalingTarget(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := applicationautoscaling.NewFromConfig(cfg)

	var values []Resource
	for _, serviceNameSpace := range types.ServiceNamespaceEcs.Values() {
		paginator := applicationautoscaling.NewDescribeScalableTargetsPaginator(client, &applicationautoscaling.DescribeScalableTargetsInput{
			ServiceNamespace: serviceNameSpace,
		})

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				return nil, err
			}

			for _, item := range page.ScalableTargets {
				resource := applicationAutoScalingTargetHandle(ctx, item)
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

	return values, nil
}
func applicationAutoScalingTargetHandle(ctx context.Context, item types.ScalableTarget) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := "arn:" + describeCtx.Partition + ":application-autoscaling:" + describeCtx.Region + ":" + describeCtx.AccountID + ":service-namespace:" + string(item.ServiceNamespace) + "/target/" + *item.ResourceId

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *item.ResourceId,
		Description: model.ApplicationAutoScalingTargetDescription{
			ScalableTarget: item,
		},
	}
	return resource
}

func GetApplicationAutoScalingTarget(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	resourceId := fields["resourceId"]
	client := applicationautoscaling.NewFromConfig(cfg)
	var values []Resource

	describers, err := client.DescribeScalableTargets(ctx, &applicationautoscaling.DescribeScalableTargetsInput{
		ResourceIds: []string{resourceId},
	})
	if err != nil {
		return nil, err
	}

	for _, item := range describers.ScalableTargets {
		resource := applicationAutoScalingTargetHandle(ctx, item)
		values = append(values, resource)
	}
	return values, nil
}

func ApplicationAutoScalingPolicy(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := applicationautoscaling.NewFromConfig(cfg)

	var values []Resource
	for _, serviceNameSpace := range types.ServiceNamespaceEcs.Values() {
		paginator := applicationautoscaling.NewDescribeScalingPoliciesPaginator(client, &applicationautoscaling.DescribeScalingPoliciesInput{
			ServiceNamespace: serviceNameSpace,
		})

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				return nil, err
			}

			for _, item := range page.ScalingPolicies {
				resource := applicationAutoScalingPolicyHandle(ctx, item)
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

	return values, nil
}
func applicationAutoScalingPolicyHandle(ctx context.Context, item types.ScalingPolicy) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := "arn:" + describeCtx.Partition + ":application-autoscaling:" + describeCtx.Region + ":" + describeCtx.AccountID + ":service-namespace:" + string(item.ServiceNamespace) + "/target/" + *item.ResourceId

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *item.ResourceId,
		Description: model.ApplicationAutoScalingPolicyDescription{
			ScalablePolicy: item,
		},
	}
	return resource
}

func GetApplicationAutoScalingPolicy(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	resourceId := fields["resourceId"]
	client := applicationautoscaling.NewFromConfig(cfg)
	var values []Resource

	describers, err := client.DescribeScalingPolicies(ctx, &applicationautoscaling.DescribeScalingPoliciesInput{
		ResourceId: &resourceId,
	})
	if err != nil {
		return nil, err
	}

	for _, item := range describers.ScalingPolicies {
		resource := applicationAutoScalingPolicyHandle(ctx, item)
		values = append(values, resource)
	}
	return values, nil
}
