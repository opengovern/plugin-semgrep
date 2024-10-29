package describer

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/serverlessapplicationrepository"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func ServerlessApplicationRepositoryApplication(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := serverlessapplicationrepository.NewFromConfig(cfg)
	paginator := serverlessapplicationrepository.NewListApplicationsPaginator(client, &serverlessapplicationrepository.ListApplicationsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, applicationSummary := range page.Applications {
			resource, err := serverlessApplicationRepositoryApplicationHandle(ctx, cfg, *applicationSummary.ApplicationId)
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
func serverlessApplicationRepositoryApplicationHandle(ctx context.Context, cfg aws.Config, applicationId string) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := serverlessapplicationrepository.NewFromConfig(cfg)

	application, err := client.GetApplication(ctx, &serverlessapplicationrepository.GetApplicationInput{
		ApplicationId: &applicationId,
	})
	if err != nil {
		return Resource{}, err
	}

	policy, err := client.GetApplicationPolicy(ctx, &serverlessapplicationrepository.GetApplicationPolicyInput{
		ApplicationId: &applicationId,
	})
	if err != nil {
		policy = &serverlessapplicationrepository.GetApplicationPolicyOutput{}
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *application.ApplicationId,
		Name:   *application.ApplicationId,
		Description: model.ServerlessApplicationRepositoryApplicationDescription{
			Application: *application,
			Statements:  policy.Statements,
		},
	}
	return resource, nil
}
func GetServerlessApplicationRepositoryApplication(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	applicationId := fields["applicationId"]

	var values []Resource
	resource, err := serverlessApplicationRepositoryApplicationHandle(ctx, cfg, applicationId)
	if err != nil {
		return nil, err
	}

	values = append(values, resource)
	return values, nil
}
