package describer

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/aws/aws-sdk-go-v2/service/ecrpublic"
	public_types "github.com/aws/aws-sdk-go-v2/service/ecrpublic/types"
	"github.com/aws/smithy-go"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func ECRPublicRepository(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	// Only supported in US-EAST-1
	if !strings.EqualFold(cfg.Region, "us-east-1") {
		return []Resource{}, nil
	}

	client := ecrpublic.NewFromConfig(cfg)
	paginator := ecrpublic.NewDescribeRepositoriesPaginator(client, &ecrpublic.DescribeRepositoriesInput{})
	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if isErr(err, "RepositoryNotFoundException") || isErr(err, "RepositoryPolicyNotFoundException") || isErr(err, "LifecyclePolicyNotFoundException") || isErr(err, "AccessDeniedException") {
				continue
			}
			return nil, err
		}

		for _, v := range page.Repositories {
			var imageDetails []public_types.ImageDetail
			imagePaginator := ecrpublic.NewDescribeImagesPaginator(client, &ecrpublic.DescribeImagesInput{
				RepositoryName: v.RepositoryName,
			})
			for imagePaginator.HasMorePages() {
				imagePage, err := imagePaginator.NextPage(ctx)
				if err != nil {
					if isErr(err, "RepositoryNotFoundException") || isErr(err, "RepositoryPolicyNotFoundException") || isErr(err, "LifecyclePolicyNotFoundException") || isErr(err, "AccessDeniedException") {
						continue
					}
					return nil, err
				}
				imageDetails = append(imageDetails, imagePage.ImageDetails...)
			}

			resource, err := eCRPublicRepositoryHandle(ctx, cfg, v, imageDetails)
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
func eCRPublicRepositoryHandle(ctx context.Context, cfg aws.Config, v public_types.Repository, imageDetails []public_types.ImageDetail) (Resource, error) {
	client := ecrpublic.NewFromConfig(cfg)
	describeCtx := GetDescribeContext(ctx)
	policyOutput, err := client.GetRepositoryPolicy(ctx, &ecrpublic.GetRepositoryPolicyInput{
		RepositoryName: v.RepositoryName,
	})
	if err != nil {
		if !isErr(err, "RepositoryNotFoundException") && !isErr(err, "RepositoryPolicyNotFoundException") && !isErr(err, "LifecyclePolicyNotFoundException") {
			return Resource{}, err
		}
	}

	tagsOutput, err := client.ListTagsForResource(ctx, &ecrpublic.ListTagsForResourceInput{
		ResourceArn: v.RepositoryArn,
	})
	if err != nil {
		if !isErr(err, "RepositoryNotFoundException") && !isErr(err, "RepositoryPolicyNotFoundException") && !isErr(err, "LifecyclePolicyNotFoundException") {
			return Resource{}, err
		} else {
			tagsOutput = &ecrpublic.ListTagsForResourceOutput{}
		}
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *v.RepositoryArn,
		Name:   *v.RepositoryName,
		Description: model.ECRPublicRepositoryDescription{
			PublicRepository: v,
			ImageDetails:     imageDetails,
			Policy:           policyOutput,
			Tags:             tagsOutput.Tags,
		},
	}
	return resource, nil
}
func GetECRPublicRepository(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	repositoryName := fields["name"]
	client := ecrpublic.NewFromConfig(cfg)
	out, err := client.DescribeRepositories(ctx, &ecrpublic.DescribeRepositoriesInput{
		RepositoryNames: []string{repositoryName},
	})
	if err != nil {
		if isErr(err, "DescribeRepositoriesNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	for _, v := range out.Repositories {

		var imageDetails []public_types.ImageDetail
		images, err := client.DescribeImages(ctx, &ecrpublic.DescribeImagesInput{
			RepositoryName: &repositoryName,
		})
		if err != nil {
			if isErr(err, "DescribeImagesNotFound") || isErr(err, "InvalidParameterValue") {
				return nil, nil
			}
			return nil, err
		}
		imageDetails = append(imageDetails, images.ImageDetails...)

		resource, err := eCRPublicRepositoryHandle(ctx, cfg, v, imageDetails)
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

func ECRPublicRegistry(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	// Only supported in US-EAST-1
	if !strings.EqualFold(cfg.Region, "us-east-1") {
		return []Resource{}, nil
	}

	client := ecrpublic.NewFromConfig(cfg)
	paginator := ecrpublic.NewDescribeRegistriesPaginator(client, &ecrpublic.DescribeRegistriesInput{})
	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Registries {
			var tags []public_types.Tag
			tagsOutput, err := client.ListTagsForResource(ctx, &ecrpublic.ListTagsForResourceInput{
				ResourceArn: v.RegistryArn,
			})
			if err == nil {
				tags = tagsOutput.Tags
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    *v.RegistryArn,
				Name:   *v.RegistryId,
				Description: model.ECRPublicRegistryDescription{
					PublicRegistry: v,
					Tags:           tags,
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

func ECRRepository(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ecr.NewFromConfig(cfg)
	paginator := ecr.NewDescribeRepositoriesPaginator(client, &ecr.DescribeRepositoriesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if isErr(err, "RepositoryNotFoundException") || isErr(err, "RepositoryPolicyNotFoundException") || isErr(err, "LifecyclePolicyNotFoundException") || isErr(err, "AccessDeniedException") {
				continue
			}
			return nil, err
		}

		for _, v := range page.Repositories {
			resource, err := eCRRepositoryHandle(ctx, cfg, v)
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
func eCRRepositoryHandle(ctx context.Context, cfg aws.Config, v types.Repository) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ecr.NewFromConfig(cfg)
	lifeCyclePolicyOutput, err := client.GetLifecyclePolicy(ctx, &ecr.GetLifecyclePolicyInput{
		RepositoryName: v.RepositoryName,
	})
	if err != nil {
		if !isErr(err, "RepositoryNotFoundException") && !isErr(err, "RepositoryPolicyNotFoundException") && !isErr(err, "LifecyclePolicyNotFoundException") || isErr(err, "AccessDeniedException") {
			return Resource{}, err
		}
	}

	var imageDetails []types.ImageDetail
	imagePaginator := ecr.NewDescribeImagesPaginator(client, &ecr.DescribeImagesInput{
		RepositoryName: v.RepositoryName,
	})
	for imagePaginator.HasMorePages() {
		imagePage, err := imagePaginator.NextPage(ctx)
		if err != nil {
			if isErr(err, "RepositoryNotFoundException") || isErr(err, "RepositoryPolicyNotFoundException") || isErr(err, "LifecyclePolicyNotFoundException") {
				continue
			}
			return Resource{}, err
		}
		imageDetails = append(imageDetails, imagePage.ImageDetails...)
	}

	policyOutput, err := client.GetRepositoryPolicy(ctx, &ecr.GetRepositoryPolicyInput{
		RepositoryName: v.RepositoryName,
		RegistryId:     v.RegistryId,
	})
	if err != nil {
		if !isErr(err, "RepositoryNotFoundException") && !isErr(err, "RepositoryPolicyNotFoundException") && !isErr(err, "LifecyclePolicyNotFoundException") {
			return Resource{}, err
		}
	}

	repositoryScanningConfiguration, err := client.BatchGetRepositoryScanningConfiguration(ctx, &ecr.BatchGetRepositoryScanningConfigurationInput{
		RepositoryNames: []string{*v.RepositoryName},
	})
	if err != nil {
		if !isErr(err, "RepositoryNotFoundException") && !isErr(err, "RepositoryPolicyNotFoundException") && !isErr(err, "LifecyclePolicyNotFoundException") {
			return Resource{}, err
		}
	}

	tagsOutput, err := client.ListTagsForResource(ctx, &ecr.ListTagsForResourceInput{
		ResourceArn: v.RepositoryArn,
	})
	if err != nil {
		if !isErr(err, "RepositoryNotFoundException") && !isErr(err, "RepositoryPolicyNotFoundException") && !isErr(err, "LifecyclePolicyNotFoundException") {
			return Resource{}, err
		} else {
			tagsOutput = &ecr.ListTagsForResourceOutput{}
		}
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *v.RepositoryArn,
		Name:   *v.RepositoryName,
		Description: model.ECRRepositoryDescription{
			Repository:                      v,
			LifecyclePolicy:                 lifeCyclePolicyOutput,
			ImageDetails:                    imageDetails,
			Policy:                          policyOutput,
			RepositoryScanningConfiguration: repositoryScanningConfiguration,
			Tags:                            tagsOutput.Tags,
		},
	}
	return resource, nil
}
func GetECRRepository(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	client := ecr.NewFromConfig(cfg)
	repositoryName := fields["name"]
	out, err := client.DescribeRepositories(ctx, &ecr.DescribeRepositoriesInput{
		RepositoryNames: []string{repositoryName},
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, repository := range out.Repositories {
		resource, err := eCRRepositoryHandle(ctx, cfg, repository)
		if err != nil {
			return nil, err
		}
		values = append(values, resource)
	}
	return values, nil
}

func ECRRegistryPolicy(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ecr.NewFromConfig(cfg)
	output, err := client.GetRegistryPolicy(ctx, &ecr.GetRegistryPolicyInput{})
	if err != nil {
		var ae smithy.APIError
		e := types.RegistryPolicyNotFoundException{}
		if errors.As(err, &ae) && ae.ErrorCode() == e.ErrorCode() {
			return []Resource{}, nil
		}
		return nil, err
	}

	var values []Resource
	resource := Resource{
		Region:      describeCtx.KaytuRegion,
		ID:          *output.RegistryId,
		Name:        *output.RegistryId,
		Description: output,
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

func ECRRegistry(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ecr.NewFromConfig(cfg)
	output, err := client.DescribeRegistry(ctx, &ecr.DescribeRegistryInput{})
	if err != nil {
		return nil, err
	}

	var values []Resource
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ID:     *output.RegistryId,
		Name:   *output.RegistryId,
		Description: model.ECRRegistryDescription{
			RegistryId:       *output.RegistryId,
			ReplicationRules: output.ReplicationConfiguration.Rules,
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

func ECRImage(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ecr.NewFromConfig(cfg)
	repositoryPaginator := ecr.NewDescribeRepositoriesPaginator(client, &ecr.DescribeRepositoriesInput{})

	var values []Resource
	for repositoryPaginator.HasMorePages() {
		repositoryPage, err := repositoryPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, repository := range repositoryPage.Repositories {
			imagesPaginator := ecr.NewDescribeImagesPaginator(client, &ecr.DescribeImagesInput{
				RepositoryName: repository.RepositoryName,
				RegistryId:     repository.RegistryId,
			})
			if err != nil {
				return nil, err
			}

			for imagesPaginator.HasMorePages() {
				page, err := imagesPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				for _, image := range page.ImageDetails {
					resource, err := eCRImageHandle(ctx, cfg, image, repository)
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
	}

	return values, nil
}
func eCRImageHandle(ctx context.Context, cfg aws.Config, image types.ImageDetail, repository types.Repository) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	var uri string
	if len(image.ImageTags) == 0 {
		uri = describeCtx.AccountID + ".dkr.ecr." + describeCtx.Region + ".amazonaws.com/" + *image.RepositoryName + "@" + *image.ImageDigest
	} else {
		uri = describeCtx.AccountID + ".dkr.ecr." + describeCtx.Region + ".amazonaws.com/" + *image.RepositoryName + ":" + image.ImageTags[0]
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		Name:   fmt.Sprintf("%s:%s", *repository.RepositoryArn, *image.ImageDigest),
		Description: model.ECRImageDescription{
			Image:    image,
			ImageUri: uri,
		},
	}
	return resource, nil
}
func GetECRImage(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	repositoryName := fields["repositoryName"]
	client := ecr.NewFromConfig(cfg)

	out, err := client.DescribeRepositories(ctx, &ecr.DescribeRepositoriesInput{
		RepositoryNames: []string{repositoryName},
	})
	if err != nil {
		if isErr(err, "DescribeRepositoriesNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	for _, repository := range out.Repositories {
		images, err := client.DescribeImages(ctx, &ecr.DescribeImagesInput{
			RepositoryName: repository.RepositoryName,
			RegistryId:     repository.RegistryId,
		})
		if err != nil {
			if isErr(err, "DescribeImagesNotFound") || isErr(err, "InvalidParameterValue") {
				return nil, nil
			}
			return nil, err
		}

		for _, image := range images.ImageDetails {
			resource, err := eCRImageHandle(ctx, cfg, image, repository)
			if err != nil {
				return nil, err
			}
			values = append(values, resource)
		}
	}
	return values, nil
}

func ECRRegistryScanningConfiguration(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ecr.NewFromConfig(cfg)
	describeCtx := GetDescribeContext(ctx)

	scanningConfiguration, err := client.GetRegistryScanningConfiguration(ctx, &ecr.GetRegistryScanningConfigurationInput{})
	if err != nil {
		return nil, err
	}
	var values []Resource

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ID:     *scanningConfiguration.RegistryId,
		Name:   *scanningConfiguration.RegistryId,
		Description: model.ECRRegistryScanningConfigurationDescription{
			RegistryId:            *scanningConfiguration.RegistryId,
			ScanningConfiguration: scanningConfiguration.ScanningConfiguration,
		},
	}
	values = append(values, resource)

	return values, nil
}
