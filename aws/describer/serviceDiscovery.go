package describer

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/servicediscovery"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func ServiceDiscoveryService(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := servicediscovery.NewFromConfig(cfg)

	paginator := servicediscovery.NewListServicesPaginator(client, &servicediscovery.ListServicesInput{})
	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, item := range page.Services {
			tag, err := client.ListTagsForResource(ctx, &servicediscovery.ListTagsForResourceInput{
				ResourceARN: item.Arn,
			})
			if err != nil {
				return nil, err
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ID:     *item.Id,
				Description: model.ServiceDiscoveryServiceDescription{
					Service: item,
					Tags:    tag.Tags,
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

func ServiceDiscoveryNamespace(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := servicediscovery.NewFromConfig(cfg)

	paginator := servicediscovery.NewListNamespacesPaginator(client, &servicediscovery.ListNamespacesInput{})
	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, v := range page.Namespaces {
			tag, err := client.ListTagsForResource(ctx, &servicediscovery.ListTagsForResourceInput{
				ResourceARN: v.Arn,
			})
			if err != nil {
				return nil, err
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ID:     *v.Id,
				Name:   *v.Name,
				Description: model.ServiceDiscoveryNamespaceDescription{
					Namespace: v,
					Tags:      tag.Tags,
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
func ServiceDiscoveryInstance(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := servicediscovery.NewFromConfig(cfg)

	paginator := servicediscovery.NewListServicesPaginator(client, &servicediscovery.ListServicesInput{})
	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, v := range page.Services {
			resources, err := getServiceDiscoveryInstances(ctx, cfg, v.Id)
			if err != nil {
				return nil, err
			}
			for _, resource := range resources {
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

func getServiceDiscoveryInstances(ctx context.Context, cfg aws.Config, id *string) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := servicediscovery.NewFromConfig(cfg)

	paginator := servicediscovery.NewListInstancesPaginator(client, &servicediscovery.ListInstancesInput{ServiceId: id})
	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, v := range page.Instances {
			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ID:     *v.Id,
				Name:   *v.Id,
				Description: model.ServiceDiscoveryInstanceDescription{
					Instance: v,
				},
			}
			values = append(values, resource)
		}
	}
	return values, nil
}
