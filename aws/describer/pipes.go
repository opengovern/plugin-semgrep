package describer

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/pipes"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func PipesPipe(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := pipes.NewFromConfig(cfg)

	paginator := pipes.NewListPipesPaginator(client, &pipes.ListPipesInput{})
	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Pipes {
			pipe, err := client.DescribePipe(ctx, &pipes.DescribePipeInput{Name: v.Name})
			if err != nil {
				return nil, err
			}

			var name string
			if pipe.Name != nil {
				name = *pipe.Name
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    *pipe.Arn,
				Name:   name,
				Description: model.PipesPipeDescription{
					PipeOutput: pipe,
					Pipe:       v,
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
