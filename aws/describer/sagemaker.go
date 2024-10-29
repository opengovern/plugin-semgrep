package describer

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func SageMakerEndpointConfiguration(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := sagemaker.NewFromConfig(cfg)
	paginator := sagemaker.NewListEndpointConfigsPaginator(client, &sagemaker.ListEndpointConfigsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range page.EndpointConfigs {

			out, err := client.DescribeEndpointConfig(ctx, &sagemaker.DescribeEndpointConfigInput{
				EndpointConfigName: item.EndpointConfigName,
			})
			if err != nil {
				if isErr(err, "DescribeEndpointConfigNotFound") || isErr(err, "InvalidParameterValue") {
					return nil, nil
				}
				return nil, err
			}

			resource, err := sageMakerEndpointConfigurationHandle(ctx, cfg, item.EndpointConfigArn, out)
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
func sageMakerEndpointConfigurationHandle(ctx context.Context, cfg aws.Config, endpointConfigArn *string, out *sagemaker.DescribeEndpointConfigOutput) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := sagemaker.NewFromConfig(cfg)

	tags, err := client.ListTags(ctx, &sagemaker.ListTagsInput{
		ResourceArn: endpointConfigArn,
	})
	if err != nil {
		if isErr(err, "ListTagsNotFound") || isErr(err, "InvalidParameterValue") {
			return Resource{}, nil
		}
		return Resource{}, err
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *out.EndpointConfigArn,
		Name:   *out.EndpointConfigName,
		Description: model.SageMakerEndpointConfigurationDescription{
			EndpointConfig: out,
			Tags:           tags.Tags,
		},
	}
	return resource, nil
}
func GetMakerEndpointConfigurationHandle(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	endpointConfigName := fields["name"]
	var values []Resource
	client := sagemaker.NewFromConfig(cfg)
	out, err := client.DescribeEndpointConfig(ctx, &sagemaker.DescribeEndpointConfigInput{
		EndpointConfigName: &endpointConfigName,
	})
	if err != nil {
		if isErr(err, "DescribeEndpointConfigNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	resource, err := sageMakerEndpointConfigurationHandle(ctx, cfg, out.EndpointConfigArn, out)
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

func SageMakerApp(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := sagemaker.NewFromConfig(cfg)

	var values []Resource
	paginator := sagemaker.NewListDomainsPaginator(client, &sagemaker.ListDomainsInput{})
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, domain := range output.Domains {
			appPaginator := sagemaker.NewListAppsPaginator(client, &sagemaker.ListAppsInput{
				DomainIdEquals: domain.DomainId,
			})

			for appPaginator.HasMorePages() {
				appOutput, err := appPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				for _, items := range appOutput.Apps {
					params := sagemaker.DescribeAppInput{
						AppName:         items.AppName,
						AppType:         items.AppType,
						DomainId:        domain.DomainId,
						UserProfileName: items.UserProfileName,
					}

					data, err := client.DescribeApp(ctx, &params)
					if err != nil {
						return nil, err
					}

					resource := Resource{
						Region: describeCtx.KaytuRegion,
						ARN:    *data.AppArn,
						Name:   *data.AppName,
						Description: model.SageMakerAppDescription{
							AppDetails:        items,
							DescribeAppOutput: data,
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

func SageMakerDomain(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := sagemaker.NewFromConfig(cfg)

	var values []Resource
	paginator := sagemaker.NewListDomainsPaginator(client, &sagemaker.ListDomainsInput{})
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, domain := range output.Domains {
			data, err := client.DescribeDomain(ctx, &sagemaker.DescribeDomainInput{
				DomainId: domain.DomainId,
			})
			if err != nil {
				return nil, err
			}
			params := &sagemaker.ListTagsInput{
				ResourceArn: domain.DomainArn,
			}

			pagesLeft := true
			var tags []types.Tag
			for pagesLeft {
				keyTags, err := client.ListTags(ctx, params)
				if err != nil {
					return nil, err
				}
				tags = append(tags, keyTags.Tags...)

				if keyTags.NextToken != nil {
					params.NextToken = keyTags.NextToken
				} else {
					pagesLeft = false
				}
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    *domain.DomainArn,
				Name:   *domain.DomainName,
				Description: model.SageMakerDomainDescription{
					Domain:     data,
					DomainItem: domain,
					Tags:       tags,
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

func SageMakerNotebookInstance(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := sagemaker.NewFromConfig(cfg)
	paginator := sagemaker.NewListNotebookInstancesPaginator(client, &sagemaker.ListNotebookInstancesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range page.NotebookInstances {
			out, err := client.DescribeNotebookInstance(ctx, &sagemaker.DescribeNotebookInstanceInput{
				NotebookInstanceName: item.NotebookInstanceName,
			})
			if err != nil {
				if isErr(err, "DescribeNotebookInstanceNotFound") || isErr(err, "InvalidParameterValue") {
					return nil, nil
				}
				return nil, err
			}

			resource, err := sageMakerNotebookInstanceHandle(ctx, cfg, out)
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
func sageMakerNotebookInstanceHandle(ctx context.Context, cfg aws.Config, out *sagemaker.DescribeNotebookInstanceOutput) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := sagemaker.NewFromConfig(cfg)

	tags, err := client.ListTags(ctx, &sagemaker.ListTagsInput{
		ResourceArn: out.NotebookInstanceArn,
	})
	if err != nil {
		if isErr(err, "ListTagsNotFound") || isErr(err, "InvalidParameterValue") {
			return Resource{}, nil
		}
		return Resource{}, err
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *out.NotebookInstanceArn,
		Name:   *out.NotebookInstanceName,
		Description: model.SageMakerNotebookInstanceDescription{
			NotebookInstance: out,
			Tags:             tags.Tags,
		},
	}
	return resource, nil
}
func GetSageMakerNotebookInstance(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	describeNetbookInstanceName := fields["name"]
	var values []Resource

	client := sagemaker.NewFromConfig(cfg)
	out, err := client.DescribeNotebookInstance(ctx, &sagemaker.DescribeNotebookInstanceInput{
		NotebookInstanceName: &describeNetbookInstanceName,
	})
	if err != nil {
		if isErr(err, "DescribeNotebookInstanceNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	resource, err := sageMakerNotebookInstanceHandle(ctx, cfg, out)
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

func SageMakerModel(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := sagemaker.NewFromConfig(cfg)

	var values []Resource
	paginator := sagemaker.NewListModelsPaginator(client, &sagemaker.ListModelsInput{})
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, modelSummary := range output.Models {
			sageModel, err := client.DescribeModel(ctx, &sagemaker.DescribeModelInput{
				ModelName: modelSummary.ModelName,
			})
			if err != nil {
				return nil, err
			}
			resource := sageMakerModelHandle(ctx, cfg, sageModel)
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
func sageMakerModelHandle(ctx context.Context, cfg aws.Config, sageModel *sagemaker.DescribeModelOutput) Resource {
	describeCtx := GetDescribeContext(ctx)
	client := sagemaker.NewFromConfig(cfg)

	tags, err := client.ListTags(ctx, &sagemaker.ListTagsInput{
		ResourceArn: sageModel.ModelArn,
	})
	if err != nil {
		tags = &sagemaker.ListTagsOutput{}
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *sageModel.ModelArn,
		Name:   *sageModel.ModelName,
		Description: model.SageMakerModelDescription{
			Model: sageModel,
			Tags:  tags.Tags,
		},
	}
	return resource
}
func GetSageMakerModel(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	modelName := fields["name"]
	var values []Resource

	client := sagemaker.NewFromConfig(cfg)
	out, err := client.DescribeModel(ctx, &sagemaker.DescribeModelInput{
		ModelName: &modelName,
	})
	if err != nil {
		if isErr(err, "DescribeModelNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	values = append(values, sageMakerModelHandle(ctx, cfg, out))
	return values, nil
}

func SageMakerTrainingJob(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := sagemaker.NewFromConfig(cfg)

	var values []Resource
	paginator := sagemaker.NewListTrainingJobsPaginator(client, &sagemaker.ListTrainingJobsInput{})
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, trainingJobSummary := range output.TrainingJobSummaries {
			trainingJob, err := client.DescribeTrainingJob(ctx, &sagemaker.DescribeTrainingJobInput{
				TrainingJobName: trainingJobSummary.TrainingJobName,
			})
			if err != nil {
				return nil, err
			}

			tags, err := client.ListTags(ctx, &sagemaker.ListTagsInput{
				ResourceArn: trainingJob.TrainingJobArn,
			})
			if err != nil {
				tags = &sagemaker.ListTagsOutput{}
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    *trainingJob.TrainingJobArn,
				Name:   *trainingJob.TrainingJobName,
				Description: model.SageMakerTrainingJobDescription{
					TrainingJob: trainingJob,
					Tags:        tags.Tags,
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
