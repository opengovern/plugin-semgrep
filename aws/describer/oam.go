package describer

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/oam"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func OAMLink(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := oam.NewFromConfig(cfg)
	paginator := oam.NewListLinksPaginator(client, &oam.ListLinksInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Items {
			if v.Arn == nil {
				continue
			}

			out, err := client.GetLink(ctx, &oam.GetLinkInput{
				Identifier: v.Arn,
			})
			if err != nil {
				return nil, err
			}

			var name string
			if out.Id != nil {
				name = *out.Id
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    *v.Arn,
				Name:   name,
				Description: model.OAMLinkDescription{
					Link: out,
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

func OAMSink(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := oam.NewFromConfig(cfg)
	paginator := oam.NewListSinksPaginator(client, &oam.ListSinksInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Items {
			if v.Arn == nil {
				continue
			}

			out, err := client.ListTagsForResource(ctx, &oam.ListTagsForResourceInput{
				ResourceArn: v.Arn,
			})
			if err != nil {
				return nil, err
			}

			var name string
			if v.Name != nil {
				name = *v.Name
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    *v.Arn,
				Name:   name,
				Description: model.OAMSinkDescription{
					Sink: v,
					Tags: out.Tags,
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
