package describer

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/timestreamwrite"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func TimestreamDatabase(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := timestreamwrite.NewFromConfig(cfg)
	paginator := timestreamwrite.NewListDatabasesPaginator(client, &timestreamwrite.ListDatabasesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Databases {
			tags, err := client.ListTagsForResource(ctx, &timestreamwrite.ListTagsForResourceInput{ResourceARN: v.Arn})
			if err != nil {
				return nil, err
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    *v.Arn,
				Name:   *v.DatabaseName,
				Description: model.TimestreamDatabaseDescription{
					Database: v,
					Tags:     tags,
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
