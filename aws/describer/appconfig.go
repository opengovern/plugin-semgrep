package describer

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/appconfig"
	"github.com/aws/aws-sdk-go-v2/service/appconfig/types"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func AppConfigApplication(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := appconfig.NewFromConfig(cfg)
	paginator := appconfig.NewListApplicationsPaginator(client, &appconfig.ListApplicationsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, application := range page.Items {
			resource, err := appConfigApplicationHandle(ctx, cfg, application)
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
func appConfigApplicationHandle(ctx context.Context, cfg aws.Config, application types.Application) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := appconfig.NewFromConfig(cfg)
	arn := fmt.Sprintf("arn:%s:appconfig:%s:%s:application/%s", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID, *application.Id)

	tags, err := client.ListTagsForResource(ctx, &appconfig.ListTagsForResourceInput{
		ResourceArn: aws.String(arn),
	})
	if err != nil {
		return Resource{}, err
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ID:     *application.Id,
		Name:   *application.Name,
		ARN:    arn,
		Description: model.AppConfigApplicationDescription{
			Application: application,
			Tags:        tags.Tags,
		},
	}
	return resource, nil
}
func GetAppConfigApplication(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	var values []Resource
	applicationId := fields["id"]
	client := appconfig.NewFromConfig(cfg)
	applications, err := client.ListApplications(ctx, &appconfig.ListApplicationsInput{})
	if err != nil {
		if isErr(err, "ListApplicationsNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	for _, application := range applications.Items {
		if application.Id != &applicationId {
			continue
		}
		resource, err := appConfigApplicationHandle(ctx, cfg, application)
		if err != nil {
			return nil, err
		}
		values = append(values, resource)
	}
	return values, nil
}
