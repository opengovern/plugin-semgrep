package describer

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/applicationinsights"
)

func ApplicationInsightsApplication(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := applicationinsights.NewFromConfig(cfg)
	paginator := applicationinsights.NewListApplicationsPaginator(client, &applicationinsights.ListApplicationsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.ApplicationInfoList {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ID:          *v.ResourceGroupName,
				Name:        *v.ResourceGroupName,
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
	}

	return values, nil
}
