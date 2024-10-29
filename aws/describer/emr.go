package describer

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/emr"
	"github.com/aws/aws-sdk-go-v2/service/emr/types"
	_ "github.com/aws/aws-sdk-go/service/configservice"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func EMRCluster(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := emr.NewFromConfig(cfg)
	paginator := emr.NewListClustersPaginator(client, &emr.ListClustersInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range page.Clusters {
			resource, err := eMRClusterHandle(ctx, cfg, *item.Id)
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
func eMRClusterHandle(ctx context.Context, cfg aws.Config, clusterId string) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := emr.NewFromConfig(cfg)

	out, err := client.DescribeCluster(ctx, &emr.DescribeClusterInput{
		ClusterId: &clusterId,
	})
	if err != nil {
		if isErr(err, "DescribeClusterNotFound") || isErr(err, "InvalidParameterValue") {
			return Resource{}, nil
		}
		return Resource{}, err
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *out.Cluster.ClusterArn,
		Name:   *out.Cluster.Name,
		Description: model.EMRClusterDescription{
			Cluster: out.Cluster,
		},
	}
	return resource, nil
}
func GetEMRCluster(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	clusterId := fields["id"]
	var values []Resource

	resource, err := eMRClusterHandle(ctx, cfg, clusterId)
	emptyResource := Resource{}
	if err == nil && resource == emptyResource {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	values = append(values, resource)
	return values, nil
}

func EMRInstance(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := emr.NewFromConfig(cfg)
	clusterPaginator := emr.NewListClustersPaginator(client, &emr.ListClustersInput{})

	var values []Resource
	for clusterPaginator.HasMorePages() {
		clusterPage, err := clusterPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, cluster := range clusterPage.Clusters {
			instancePaginator := emr.NewListInstancesPaginator(client, &emr.ListInstancesInput{
				ClusterId: cluster.Id,
			})

			for instancePaginator.HasMorePages() {
				instancePage, err := instancePaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				for _, instance := range instancePage.Instances {
					describeCtx := GetDescribeContext(ctx)
					arn := fmt.Sprintf("arn:%s:emr:%s:%s:instance/%s", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID, *instance.Id)
					resource := Resource{
						Region: describeCtx.KaytuRegion,
						ID:     *instance.Id,
						ARN:    arn,
						Description: model.EMRInstanceDescription{
							Instance:  instance,
							ClusterID: *cluster.Id,
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

func EMRInstanceFleet(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := emr.NewFromConfig(cfg)
	clusterPaginator := emr.NewListClustersPaginator(client, &emr.ListClustersInput{})

	var values []Resource
	for clusterPaginator.HasMorePages() {
		clusterPage, err := clusterPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, cluster := range clusterPage.Clusters {
			instancePaginator := emr.NewListInstanceFleetsPaginator(client, &emr.ListInstanceFleetsInput{
				ClusterId: cluster.Id,
			})

			for instancePaginator.HasMorePages() {
				instancePage, err := instancePaginator.NextPage(ctx)
				if err != nil {
					if isErr(err, "InvalidRequestException") {
						break
					}
					return nil, err
				}

				for _, instanceFleet := range instancePage.InstanceFleets {
					resource := eMRInstanceFleetHandle(ctx, instanceFleet, *cluster.Id)
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
func eMRInstanceFleetHandle(ctx context.Context, instanceFleet types.InstanceFleet, clusterId string) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := fmt.Sprintf("arn:%s:emr:%s:%s:instance-fleet/%s", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID, *instanceFleet.Id)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ID:     *instanceFleet.Id,
		Name:   *instanceFleet.Name,
		ARN:    arn,
		Description: model.EMRInstanceFleetDescription{
			InstanceFleet: instanceFleet,
			ClusterID:     clusterId,
		},
	}
	return resource
}
func GetEMRInstanceFleet(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	var values []Resource
	clusterId := fields["clusterId"]
	client := emr.NewFromConfig(cfg)

	listInstances, err := client.ListInstanceFleets(ctx, &emr.ListInstanceFleetsInput{
		ClusterId: &clusterId,
	})
	if err != nil {
		if isErr(err, "ListInstanceFleetsNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	for _, instance := range listInstances.InstanceFleets {
		values = append(values, eMRInstanceFleetHandle(ctx, instance, clusterId))
	}
	return values, nil
}

func EMRInstanceGroup(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := emr.NewFromConfig(cfg)
	clusterPaginator := emr.NewListClustersPaginator(client, &emr.ListClustersInput{})

	var values []Resource
	for clusterPaginator.HasMorePages() {
		clusterPage, err := clusterPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, cluster := range clusterPage.Clusters {
			instancePaginator := emr.NewListInstanceGroupsPaginator(client, &emr.ListInstanceGroupsInput{
				ClusterId: cluster.Id,
			})

			for instancePaginator.HasMorePages() {
				instancePage, err := instancePaginator.NextPage(ctx)
				if err != nil {
					if isErr(err, "InvalidRequestException") {
						break
					}
					return nil, err
				}

				for _, instanceGroup := range instancePage.InstanceGroups {
					resource := eMRInstanceGroupHandle(ctx, instanceGroup, *cluster.Id)
					if instanceGroup.Name != nil {
						resource.Name = *instanceGroup.Name
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
func eMRInstanceGroupHandle(ctx context.Context, instanceGroup types.InstanceGroup, clusterId string) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := fmt.Sprintf("arn:%s:emr:%s:%s:instance-group/%s", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID, *instanceGroup.Id)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ID:     *instanceGroup.Id,
		ARN:    arn,
		Description: model.EMRInstanceGroupDescription{
			InstanceGroup: instanceGroup,
			ClusterID:     clusterId,
		},
	}
	return resource
}
func GetEMRInstanceGroup(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	clusterId := fields["clusterId"]
	client := emr.NewFromConfig(cfg)
	var values []Resource

	instances, err := client.ListInstanceGroups(ctx, &emr.ListInstanceGroupsInput{
		ClusterId: &clusterId,
	})
	if err != nil {
		if isErr(err, "ListInstanceGroupsNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	for _, instanceGroup := range instances.InstanceGroups {
		values = append(values, eMRInstanceGroupHandle(ctx, instanceGroup, clusterId))
	}
	return values, nil
}

func EMRBlockPublicAccessConfiguration(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := emr.NewFromConfig(cfg)
	describeCtx := GetDescribeContext(ctx)
	op, err := client.GetBlockPublicAccessConfiguration(ctx, &emr.GetBlockPublicAccessConfigurationInput{})
	if err != nil {
		return nil, err
	}
	var values []Resource
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		Description: model.EMRBlockPublicAccessConfigurationDescription{
			Configuration:         *op.BlockPublicAccessConfiguration,
			ConfigurationMetadata: *op.BlockPublicAccessConfigurationMetadata,
		},
	}

	if stream != nil {
		if err := (*stream)(resource); err != nil {
			return nil, err
		}
	} else {
		values = append(values, resource)
	}
	return values, nil
}
