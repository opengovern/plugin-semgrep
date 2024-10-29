package describer

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/globalaccelerator"
	"github.com/aws/aws-sdk-go-v2/service/globalaccelerator/types"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func GlobalAcceleratorAccelerator(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := globalaccelerator.NewFromConfig(cfg)
	paginator := globalaccelerator.NewListAcceleratorsPaginator(client, &globalaccelerator.ListAcceleratorsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, accelerator := range page.Accelerators {
			attribute, err := client.DescribeAcceleratorAttributes(ctx, &globalaccelerator.DescribeAcceleratorAttributesInput{
				AcceleratorArn: accelerator.AcceleratorArn,
			})
			if err != nil {
				return nil, err
			}

			tags, err := client.ListTagsForResource(ctx, &globalaccelerator.ListTagsForResourceInput{
				ResourceArn: accelerator.AcceleratorArn,
			})
			if err != nil {
				return nil, err
			}

			resource := globalAcceleratorAcceleratorHandle(ctx, attribute, tags, accelerator)
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
func globalAcceleratorAcceleratorHandle(ctx context.Context, attribute *globalaccelerator.DescribeAcceleratorAttributesOutput, tags *globalaccelerator.ListTagsForResourceOutput, accelerator types.Accelerator) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *accelerator.AcceleratorArn,
		Name:   *accelerator.Name,
		Description: model.GlobalAcceleratorAcceleratorDescription{
			Accelerator:           accelerator,
			AcceleratorAttributes: attribute.AcceleratorAttributes,
			Tags:                  tags.Tags,
		},
	}
	return resource
}
func GetGlobalAcceleratorAccelerator(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	acceleratorArn := fields["arn"]
	client := globalaccelerator.NewFromConfig(cfg)
	var values []Resource

	accelerator, err := client.DescribeAccelerator(ctx, &globalaccelerator.DescribeAcceleratorInput{
		AcceleratorArn: &acceleratorArn,
	})
	if err != nil {
		if isErr(err, "DescribeAcceleratorNotFound") || isErr(err, "invalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	attribute, err := client.DescribeAcceleratorAttributes(ctx, &globalaccelerator.DescribeAcceleratorAttributesInput{
		AcceleratorArn: accelerator.Accelerator.AcceleratorArn,
	})
	if err != nil {
		if isErr(err, "DescribeAcceleratorAttributesNotFound") || isErr(err, "invalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	tags, err := client.ListTagsForResource(ctx, &globalaccelerator.ListTagsForResourceInput{
		ResourceArn: accelerator.Accelerator.AcceleratorArn,
	})
	if err != nil {
		if isErr(err, "ListTagsForResourceNotFound") || isErr(err, "invalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	resource := globalAcceleratorAcceleratorHandle(ctx, attribute, tags, *accelerator.Accelerator)
	values = append(values, resource)
	return values, nil
}

func GlobalAcceleratorListener(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := globalaccelerator.NewFromConfig(cfg)
	paginator := globalaccelerator.NewListAcceleratorsPaginator(client, &globalaccelerator.ListAcceleratorsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, accelerator := range page.Accelerators {
			listenerPaginator := globalaccelerator.NewListListenersPaginator(client, &globalaccelerator.ListListenersInput{
				AcceleratorArn: accelerator.AcceleratorArn,
			})
			for listenerPaginator.HasMorePages() {

				listenerPage, err := listenerPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				for _, listener := range listenerPage.Listeners {

					resource := globalAcceleratorListenerHandle(ctx, listener, *accelerator.AcceleratorArn)
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
func globalAcceleratorListenerHandle(ctx context.Context, listener types.Listener, ARN string) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *listener.ListenerArn,
		Name:   *listener.ListenerArn,
		Description: model.GlobalAcceleratorListenerDescription{
			Listener:       listener,
			AcceleratorArn: ARN,
		},
	}
	return resource
}
func GetGlobalAcceleratorListener(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	var values []Resource
	acceleratorArn := fields["arn"]
	client := globalaccelerator.NewFromConfig(cfg)
	accelerator, err := client.DescribeAccelerator(ctx, &globalaccelerator.DescribeAcceleratorInput{
		AcceleratorArn: &acceleratorArn,
	})
	if err != nil {
		if isErr(err, "DescribeAcceleratorNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	describeListener, err := client.ListListeners(ctx, &globalaccelerator.ListListenersInput{
		AcceleratorArn: accelerator.Accelerator.AcceleratorArn,
	})
	if err != nil {
		if isErr(err, "ListListenersNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	for _, listener := range describeListener.Listeners {
		resource := globalAcceleratorListenerHandle(ctx, listener, acceleratorArn)
		values = append(values, resource)
	}
	return values, nil
}

func GlobalAcceleratorEndpointGroup(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := globalaccelerator.NewFromConfig(cfg)
	paginator := globalaccelerator.NewListAcceleratorsPaginator(client, &globalaccelerator.ListAcceleratorsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, accelerator := range page.Accelerators {
			listenerPaginator := globalaccelerator.NewListListenersPaginator(client, &globalaccelerator.ListListenersInput{
				AcceleratorArn: accelerator.AcceleratorArn,
			})
			for listenerPaginator.HasMorePages() {
				listenerPage, err := listenerPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}
				for _, listener := range listenerPage.Listeners {
					endpointGroupPaginator := globalaccelerator.NewListEndpointGroupsPaginator(client, &globalaccelerator.ListEndpointGroupsInput{
						ListenerArn: listener.ListenerArn,
					})
					for endpointGroupPaginator.HasMorePages() {
						endpointGroupPage, err := endpointGroupPaginator.NextPage(ctx)
						if err != nil {
							return nil, err
						}
						for _, endpointGroup := range endpointGroupPage.EndpointGroups {
							resource := globalAcceleratorEndpointGroupHandle(ctx, endpointGroup, listener, *accelerator.AcceleratorArn)
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
		}
	}

	return values, nil
}
func globalAcceleratorEndpointGroupHandle(ctx context.Context, endpointGroup types.EndpointGroup, listener types.Listener, ARN string) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *endpointGroup.EndpointGroupArn,
		Name:   *endpointGroup.EndpointGroupArn,
		Description: model.GlobalAcceleratorEndpointGroupDescription{
			EndpointGroup:  endpointGroup,
			ListenerArn:    *listener.ListenerArn,
			AcceleratorArn: ARN,
		},
	}
	return resource
}
func GetGlobalAcceleratorEndpointGroup(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	acceleratorArn := fields["arn"]
	client := globalaccelerator.NewFromConfig(cfg)
	var values []Resource

	accelerator, err := client.DescribeAccelerator(ctx, &globalaccelerator.DescribeAcceleratorInput{
		AcceleratorArn: &acceleratorArn,
	})
	if err != nil {
		if isErr(err, "DescribeAcceleratorNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	describeListener, err := client.ListListeners(ctx, &globalaccelerator.ListListenersInput{
		AcceleratorArn: accelerator.Accelerator.AcceleratorArn,
	})
	if err != nil {
		if isErr(err, "ListListenersNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	for _, listener := range describeListener.Listeners {
		listEndpoint, err := client.ListEndpointGroups(ctx, &globalaccelerator.ListEndpointGroupsInput{
			ListenerArn: listener.ListenerArn,
		})
		if err != nil {
			if isErr(err, "ListEndpointGroupsNotFound") || isErr(err, "InvalidParameterValue") {
				return nil, nil
			}
			return nil, err
		}

		for _, endpointGroup := range listEndpoint.EndpointGroups {

			resource := globalAcceleratorEndpointGroupHandle(ctx, endpointGroup, listener, acceleratorArn)
			values = append(values, resource)

		}
	}
	return values, nil
}
