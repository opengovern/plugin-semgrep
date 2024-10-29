package describer

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/mgn"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func MGNApplication(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := mgn.NewFromConfig(cfg)
	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		applications, err := client.ListApplications(ctx, &mgn.ListApplicationsInput{
			NextToken: prevToken,
		})
		if err != nil {
			if isErr(err, "UninitializedAccountException") {
				return nil, nil
			}
			return nil, err
		}

		for _, application := range applications.Items {
			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    *application.Arn,
				Name:   *application.Name,
				ID:     *application.ApplicationID,
				Description: model.MgnApplicationDescription{
					Application: application,
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

		return applications.NextToken, nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}
