package describer

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/mwaa"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func MWAAEnvironment(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := mwaa.NewFromConfig(cfg)
	paginator := mwaa.NewListEnvironmentsPaginator(client, &mwaa.ListEnvironmentsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Environments {
			resource, err := mWAAEnvironmentHandle(ctx, cfg, v)
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
func mWAAEnvironmentHandle(ctx context.Context, cfg aws.Config, v string) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := mwaa.NewFromConfig(cfg)
	environment, err := client.GetEnvironment(ctx, &mwaa.GetEnvironmentInput{
		Name: &v,
	})
	if err != nil {
		if isErr(err, "GetEnvironmentNotFound") || isErr(err, "InvalidParameterVaLue") {
			return Resource{}, nil
		}
		return Resource{}, err
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *environment.Environment.Arn,
		Name:   *environment.Environment.Name,
		Description: model.MWAAEnvironmentDescription{
			Environment: *environment.Environment,
		},
	}
	return resource, nil
}
func GetMWAAEnvironment(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	environmentName := fields["name"]
	var values []Resource
	resource, err := mWAAEnvironmentHandle(ctx, cfg, environmentName)
	if err != nil {
		return nil, err
	}
	emptyResource := Resource{}
	if err == nil && resource == emptyResource {
		return nil, nil
	}
	values = append(values, resource)
	return values, nil
}
