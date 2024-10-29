package describer

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/opensearchserverless"
	"github.com/aws/aws-sdk-go-v2/service/opensearchserverless/types"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func OpenSearchServerlessCollection(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := opensearchserverless.NewFromConfig(cfg)
	paginator := opensearchserverless.NewListCollectionsPaginator(client, &opensearchserverless.ListCollectionsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.CollectionSummaries {
			collections, err := client.BatchGetCollection(ctx, &opensearchserverless.BatchGetCollectionInput{Ids: []string{*v.Id}})
			if err != nil {
				return nil, err
			}

			var collection types.CollectionDetail
			if len(collections.CollectionDetails) > 0 {
				collection = collections.CollectionDetails[0]
			}

			tags, err := client.ListTagsForResource(ctx, &opensearchserverless.ListTagsForResourceInput{ResourceArn: v.Arn})
			if err != nil {
				return nil, err
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    *v.Arn,
				Name:   *v.Name,
				Description: model.OpenSearchServerlessCollectionDescription{
					CollectionSummary: v,
					Collection:        collection,
					Tags:              tags,
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
