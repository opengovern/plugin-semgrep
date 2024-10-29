package describer

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func ECSCapacityProvider(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ecs.NewFromConfig(cfg)

	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		output, err := client.DescribeCapacityProviders(ctx, &ecs.DescribeCapacityProvidersInput{NextToken: prevToken})
		if err != nil {
			return nil, err
		}
		if len(output.Failures) != 0 {
			return nil, failuresToError(output.Failures)
		}

		for _, v := range output.CapacityProviders {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ARN:         *v.CapacityProviderArn,
				Name:        *v.Name,
				Description: v,
			}
			if stream != nil {
				if err := (*stream)(resource); err != nil {
					return nil, err
				}
			} else {
				values = append(values, resource)
			}
		}

		return output.NextToken, nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}

func ECSCluster(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	clusters, err := listEcsClusters(ctx, cfg, nil)
	if err != nil {
		return nil, err
	}

	client := ecs.NewFromConfig(cfg)

	var values []Resource
	// Describe in batch of 100 which is the limit
	for i := 0; i < len(clusters); i = i + 100 {
		j := i + 100
		if j > len(clusters) {
			j = len(clusters)
		}

		output, err := client.DescribeClusters(ctx, &ecs.DescribeClustersInput{
			Clusters: clusters[i:j],
		})
		if err != nil {
			return nil, err
		}
		if len(output.Failures) != 0 {
			return nil, failuresToError(output.Failures)
		}

		for _, v := range output.Clusters {
			resource := eCSClusterHandle(ctx, v)
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
func eCSClusterHandle(ctx context.Context, v types.Cluster) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *v.ClusterArn,
		Name:   *v.ClusterName,
		Description: model.ECSClusterDescription{
			Cluster: v,
		},
	}
	return resource
}
func GetECSCluster(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	client := ecs.NewFromConfig(cfg)

	cluster := fields["name"]

	var values []Resource
	output, err := client.DescribeClusters(ctx, &ecs.DescribeClustersInput{
		Clusters: []string{cluster},
	})
	if err != nil {
		return nil, err
	}
	if len(output.Failures) != 0 {
		return nil, failuresToError(output.Failures)
	}

	for _, v := range output.Clusters {
		values = append(values, eCSClusterHandle(ctx, v))
	}

	return values, nil
}

func ECSService(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	clusters, err := listEcsClusters(ctx, cfg, nil)
	if err != nil {
		return nil, err
	}

	client := ecs.NewFromConfig(cfg)

	var values []Resource
	for _, cluster := range clusters {
		// This prevents Implicit memory aliasing in for loop
		cluster := cluster
		services, err := listECsServices(ctx, cfg, cluster)
		if err != nil {
			return nil, err
		}

		// Describe in batch of 10 which is the limit
		for i := 0; i < len(services); i = i + 10 {
			j := i + 10
			if j > len(services) {
				j = len(services)
			}

			output, err := client.DescribeServices(ctx, &ecs.DescribeServicesInput{
				Cluster:  &cluster,
				Services: services[i:j],
			})
			if err != nil {
				return nil, err
			}
			if len(output.Failures) != 0 {
				return nil, failuresToError(output.Failures)
			}

			for _, v := range output.Services {
				resource, err := eCSServiceHandle(ctx, v, client)
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
	}

	return values, nil
}
func eCSServiceHandle(ctx context.Context, v types.Service, client *ecs.Client) (Resource, error) {
	params := &ecs.ListTagsForResourceInput{
		ResourceArn: v.ServiceArn,
	}

	response, err := client.ListTagsForResource(ctx, params)
	if err != nil {
		return Resource{}, err
	}
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *v.ServiceArn,
		Name:   *v.ServiceName,
		Description: model.ECSServiceDescription{
			Service: v,
			Tags:    response.Tags,
		},
	}
	return resource, err
}
func GetECSService(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	cluster := fields["cluster"]
	service := fields["service"]
	client := ecs.NewFromConfig(cfg)

	var values []Resource
	output, err := client.DescribeServices(ctx, &ecs.DescribeServicesInput{
		Cluster:  &cluster,
		Services: []string{service},
	})
	if err != nil {
		return nil, err
	}
	if len(output.Failures) != 0 {
		return nil, failuresToError(output.Failures)
	}

	for _, v := range output.Services {
		resource, err := eCSServiceHandle(ctx, v, client)
		if err != nil {
			return nil, err
		}
		values = append(values, resource)
	}

	return values, nil
}

func ECSTaskDefinition(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ecs.NewFromConfig(cfg)
	paginator := ecs.NewListTaskDefinitionsPaginator(client, &ecs.ListTaskDefinitionsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, arn := range page.TaskDefinitionArns {

			resource, err := eCSTaskDefinitionHandle(ctx, cfg, arn)
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
func eCSTaskDefinitionHandle(ctx context.Context, cfg aws.Config, arn string) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ecs.NewFromConfig(cfg)

	output, err := client.DescribeTaskDefinition(ctx, &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: &arn,
		Include: []types.TaskDefinitionField{
			types.TaskDefinitionFieldTags,
		},
	})
	if err != nil {
		if isErr(err, "DescribeTaskDefinitionNotFound") || isErr(err, "InvalidParameterValue") {
			return Resource{}, nil
		}
		return Resource{}, err
	}

	// From Steampipe
	splitArn := strings.Split(arn, ":")
	name := splitArn[len(splitArn)-1]

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   name,
		Description: model.ECSTaskDefinitionDescription{
			TaskDefinition: output.TaskDefinition,
			Tags:           output.Tags,
		},
	}
	return resource, nil
}
func GetECSTaskDefinition(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	taskDefinitionARN := fields["arn"]
	var values []Resource

	resource, err := eCSTaskDefinitionHandle(ctx, cfg, taskDefinitionARN)
	if err != nil {
		return nil, err
	}
	emptyResource := Resource{}
	if err == nil && resource == emptyResource {
		return nil, nil
	}

	return values, nil
}

func ECSTaskSet(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	clusters, err := listEcsClusters(ctx, cfg, nil)
	if err != nil {
		return nil, err
	}

	client := ecs.NewFromConfig(cfg)
	var values []Resource

	for _, cluster := range clusters {
		cluster := cluster
		services, err := listECsServices(ctx, cfg, cluster)
		if err != nil {
			return nil, err
		}

		for i := 0; i < len(services); i = i + 10 {
			j := i + 10
			if j > len(services) {
				j = len(services)
			}

			serviceOutput, err := client.DescribeServices(ctx, &ecs.DescribeServicesInput{
				Cluster:  &cluster,
				Services: services[i:j],
			})
			if err != nil {
				return nil, err
			}
			if len(serviceOutput.Failures) != 0 {
				return nil, failuresToError(serviceOutput.Failures)
			}

			for _, service := range serviceOutput.Services {
				service := service
				if err != nil {
					return nil, err
				}
				for _, v := range service.TaskSets {
					resource := eCSTaskSetHandle(ctx, v)
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
func eCSTaskSetHandle(ctx context.Context, v types.TaskSet) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *v.TaskSetArn,
		Name:   *v.Id,
		Description: model.ECSTaskSetDescription{
			TaskSet: v,
		},
	}
	return resource
}
func GetECSTaskSet(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	cluster := fields["cluster"]
	service := fields["service"]
	client := ecs.NewFromConfig(cfg)

	taskSets, err := client.DescribeTaskSets(ctx, &ecs.DescribeTaskSetsInput{
		Cluster: &cluster,
		Service: &service,
	})
	if err != nil {
		if isErr(err, "DescribeTaskSetsNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	for _, v := range taskSets.TaskSets {
		resource := eCSTaskSetHandle(ctx, v)
		values = append(values, resource)
	}
	return values, nil
}

func listECsServices(ctx context.Context, cfg aws.Config, cluster string) ([]string, error) {
	client := ecs.NewFromConfig(cfg)
	paginator := ecs.NewListServicesPaginator(client, &ecs.ListServicesInput{
		Cluster: &cluster,
	})

	var services []string
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		services = append(services, page.ServiceArns...)
	}

	return services, nil
}

func listEcsClusters(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]string, error) {
	client := ecs.NewFromConfig(cfg)
	paginator := ecs.NewListClustersPaginator(client, &ecs.ListClustersInput{})

	var clusters []string
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		clusters = append(clusters, page.ClusterArns...)
	}

	return clusters, nil
}

func ECSContainerInstance(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	clusters, err := listEcsClusters(ctx, cfg, nil)
	if err != nil {
		return nil, err
	}

	client := ecs.NewFromConfig(cfg)

	var values []Resource
	for _, cluster := range clusters {
		paginator := ecs.NewListContainerInstancesPaginator(client, &ecs.ListContainerInstancesInput{
			Cluster: &cluster,
		})
		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				return nil, err
			}

			if page.ContainerInstanceArns == nil || len(page.ContainerInstanceArns) == 0 {
				continue
			}
			describeCluster, err := client.DescribeClusters(ctx, &ecs.DescribeClustersInput{
				Clusters: clusters,
			})
			if err != nil {
				if isErr(err, "DescribeClustersNotFound") || isErr(err, "InvalidParameterValue") {
					return nil, nil
				}
				return nil, err
			}

			output, err := client.DescribeContainerInstances(ctx, &ecs.DescribeContainerInstancesInput{
				Cluster:            &cluster,
				ContainerInstances: page.ContainerInstanceArns,
			})
			if err != nil {
				return nil, err
			}
			if len(output.Failures) != 0 {
				return nil, failuresToError(output.Failures)
			}

			for _, v := range output.ContainerInstances {
				for _, c := range describeCluster.Clusters {
					resource := Resource{
						Region: describeCtx.KaytuRegion,
						ARN:    *v.ContainerInstanceArn,
						Name:   *v.ContainerInstanceArn,
						Description: model.ECSContainerInstanceDescription{
							ContainerInstance: v,
							Cluster:           c,
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

func ECSTask(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	clusters, err := listEcsClusters(ctx, cfg, nil)
	if err != nil {
		return nil, err
	}

	client := ecs.NewFromConfig(cfg)
	var values []Resource

	for _, cluster := range clusters {
		cluster := cluster
		services, err := listECsServices(ctx, cfg, cluster)
		if err != nil {
			return nil, err
		}

		for _, service := range services {
			service := service
			paginator := ecs.NewListTasksPaginator(client, &ecs.ListTasksInput{
				Cluster:     &cluster,
				ServiceName: &service,
			})
			for paginator.HasMorePages() {
				page, err := paginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				if page.TaskArns == nil || len(page.TaskArns) == 0 {
					continue
				}
				output, err := client.DescribeTasks(ctx, &ecs.DescribeTasksInput{
					Cluster: &cluster,
					Tasks:   page.TaskArns,
				})
				if err != nil {
					return nil, err
				}
				if len(output.Failures) != 0 {
					return nil, failuresToError(output.Failures)
				}
				taskProtections, err := client.GetTaskProtection(ctx, &ecs.GetTaskProtectionInput{
					Cluster: &cluster,
					Tasks:   page.TaskArns,
				})
				if err != nil {
					return nil, err
				}
				if len(taskProtections.Failures) != 0 {
					return nil, failuresToError(output.Failures)
				}

				taskProtectionMap := make(map[string]types.ProtectedTask)
				for _, taskProtection := range taskProtections.ProtectedTasks {
					taskProtectionMap[*taskProtection.TaskArn] = taskProtection
				}

				for _, v := range output.Tasks {
					description := model.ECSTaskDescription{
						Task:           v,
						ServiceName:    service,
						TaskProtection: nil,
					}
					if taskProtection, ok := taskProtectionMap[*v.TaskArn]; ok {
						description.TaskProtection = &taskProtection
					}
					resource := Resource{
						Region:      describeCtx.KaytuRegion,
						ARN:         *v.TaskArn,
						Name:        *v.TaskArn,
						Description: description,
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

func failuresToError(failures []types.Failure) error {
	var errs []string
	for _, f := range failures {
		errs = append(errs, fmt.Sprintf("Arn=%s, Detail=%s, Reason=%s",
			aws.ToString(f.Arn),
			aws.ToString(f.Detail),
			aws.ToString(f.Reason)))
	}

	return errors.New(strings.Join(errs, ";"))
}
