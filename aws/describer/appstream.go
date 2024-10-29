package describer

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/appstream"
	"github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func AppStreamApplication(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := appstream.NewFromConfig(cfg)

	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		output, err := client.DescribeApplications(ctx, &appstream.DescribeApplicationsInput{
			NextToken: prevToken,
		})
		if err != nil {
			return nil, err
		}

		for _, item := range output.Applications {
			tags, err := client.ListTagsForResource(ctx, &appstream.ListTagsForResourceInput{
				ResourceArn: item.Arn,
			})
			if err != nil {
				return nil, err
			}
			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    *item.Arn,
				Name:   *item.Name,
				Description: model.AppStreamApplicationDescription{
					Application: item,
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
		return output.NextToken, nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}

func AppStreamStack(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := appstream.NewFromConfig(cfg)

	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		output, err := client.DescribeStacks(ctx, &appstream.DescribeStacksInput{
			NextToken: prevToken,
		})
		if err != nil {
			return nil, err
		}

		for _, item := range output.Stacks {

			resource, err := appStreamStackHandle(ctx, cfg, item)
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
		return output.NextToken, nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}
func appStreamStackHandle(ctx context.Context, cfg aws.Config, item types.Stack) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := appstream.NewFromConfig(cfg)

	tags, err := client.ListTagsForResource(ctx, &appstream.ListTagsForResourceInput{
		ResourceArn: item.Arn,
	})
	if err != nil {
		if isErr(err, "ListTagsForResourceNotFound") || isErr(err, "InvalidParameterValue") {
			return Resource{}, nil
		}
		return Resource{}, err
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *item.Arn,
		Name:   *item.Name,
		Description: model.AppStreamStackDescription{
			Stack: item,
			Tags:  tags.Tags,
		},
	}
	return resource, nil
}
func GetAppStreamStack(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	name := fields["name"]
	client := appstream.NewFromConfig(cfg)
	out, err := client.DescribeStacks(ctx, &appstream.DescribeStacksInput{
		Names: []string{name},
	})
	if err != nil {
		if isErr(err, "DescribeStacksNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	for _, v := range out.Stacks {

		resource, err := appStreamStackHandle(ctx, cfg, v)
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

func AppStreamFleet(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := appstream.NewFromConfig(cfg)

	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		output, err := client.DescribeFleets(ctx, &appstream.DescribeFleetsInput{
			NextToken: prevToken,
		})
		if err != nil {
			return nil, err
		}

		for _, item := range output.Fleets {
			resource, err := appStreamFleetHandle(ctx, cfg, item)
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
		return output.NextToken, nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}
func appStreamFleetHandle(ctx context.Context, cfg aws.Config, item types.Fleet) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := appstream.NewFromConfig(cfg)
	tags, err := client.ListTagsForResource(ctx, &appstream.ListTagsForResourceInput{
		ResourceArn: item.Arn,
	})
	if err != nil {
		if isErr(err, "ListTagsForResourceNotFound") || isErr(err, "InvalidParameterValue") {
			return Resource{}, nil
		}
		return Resource{}, err
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *item.Arn,
		Name:   *item.Name,
		Description: model.AppStreamFleetDescription{
			Fleet: item,
			Tags:  tags.Tags,
		},
	}

	return resource, nil
}
func GetAppStreamFleet(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	name := fields["name"]
	client := appstream.NewFromConfig(cfg)

	out, err := client.DescribeFleets(ctx, &appstream.DescribeFleetsInput{
		Names: []string{name},
	})
	if err != nil {
		if isErr(err, "DescribeFleetsNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	for _, v := range out.Fleets {

		resource, err := appStreamFleetHandle(ctx, cfg, v)
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

func AppStreamImage(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := appstream.NewFromConfig(cfg)

	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		output, err := client.DescribeImages(ctx, &appstream.DescribeImagesInput{
			NextToken: prevToken,
		})
		if err != nil {
			if isErr(err, "AccessDeniedException") {
				return nil, nil
			} else {
				return nil, err
			}
		}

		for _, item := range output.Images {
			resource, err := appStreamImageHandle(ctx, cfg, item)
			emptyResource := Resource{}
			if err == nil && resource == emptyResource {
				return nil, nil
			}
			if err != nil {
				if isErr(err, "AccessDeniedException") {
					return nil, nil
				} else {
					return nil, err
				}
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
func appStreamImageHandle(ctx context.Context, cfg aws.Config, item types.Image) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := appstream.NewFromConfig(cfg)
	tags, err := client.ListTagsForResource(ctx, &appstream.ListTagsForResourceInput{
		ResourceArn: item.Arn,
	})
	if err != nil {
		if isErr(err, "ListTagsForResourceNotFound") || isErr(err, "InvalidParameterValue") || isErr(err, "AccessDeniedException") {
			return Resource{}, nil
		}
		return Resource{}, err
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *item.Arn,
		Name:   *item.Name,
		Description: model.AppStreamImageDescription{
			Image: item,
			Tags:  tags.Tags,
		},
	}

	return resource, nil
}
func GetAppStreamImage(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	name := fields["name"]
	client := appstream.NewFromConfig(cfg)

	out, err := client.DescribeImages(ctx, &appstream.DescribeImagesInput{
		Names: []string{name},
	})
	if err != nil {
		if isErr(err, "DescribeFleetsNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	for _, v := range out.Images {

		resource, err := appStreamImageHandle(ctx, cfg, v)
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
