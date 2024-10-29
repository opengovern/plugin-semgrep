package describer

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/imagebuilder/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/imagebuilder"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func ImageBuilderImage(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := imagebuilder.NewFromConfig(cfg)
	paginator := imagebuilder.NewListImagesPaginator(client, &imagebuilder.ListImagesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.ImageVersionList {
			buildVersionPaginator := imagebuilder.NewListImageBuildVersionsPaginator(client, &imagebuilder.ListImageBuildVersionsInput{
				ImageVersionArn: v.Arn,
			})
			for buildVersionPaginator.HasMorePages() {
				buildVersionPage, err := buildVersionPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				for _, imageSummary := range buildVersionPage.ImageSummaryList {
					resource, err := imageBuilderImageHandle(ctx, cfg, imageSummary)
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
		}
	}

	return values, nil
}
func imageBuilderImageHandle(ctx context.Context, cfg aws.Config, imageSummary types.ImageSummary) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := imagebuilder.NewFromConfig(cfg)
	image, err := client.GetImage(ctx, &imagebuilder.GetImageInput{
		ImageBuildVersionArn: imageSummary.Arn,
	})
	if err != nil {
		if isErr(err, "GetImageNotFound") || isErr(err, "InvalidParameterValue") {
			return Resource{}, nil
		}
		return Resource{}, err
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *image.Image.Arn,
		Name:   *image.Image.Name,
		Description: model.ImageBuilderImageDescription{
			Image: *image.Image,
		},
	}
	return resource, nil
}
func GetImageBuilderImage(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	arn := fields["arn"]
	client := imagebuilder.NewFromConfig(cfg)
	out, err := client.ListImageBuildVersions(ctx, &imagebuilder.ListImageBuildVersionsInput{
		ImageVersionArn: &arn,
	})
	if err != nil {
		if isErr(err, "ListImageBuildVersionsNotfound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}
	var values []Resource
	for _, imageSummery := range out.ImageSummaryList {

		resource, err := imageBuilderImageHandle(ctx, cfg, imageSummery)
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
