package describer

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptune"
	"github.com/aws/aws-sdk-go-v2/service/neptune/types"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func NeptuneDatabase(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := neptune.NewFromConfig(cfg)
	paginator := neptune.NewDescribeDBInstancesPaginator(client, &neptune.DescribeDBInstancesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.DBInstances {
			if v.DBInstanceArn == nil {
				continue
			}
			if *v.Engine != "neptune" {
				continue
			}
			tags, err := client.ListTagsForResource(ctx, &neptune.ListTagsForResourceInput{
				ResourceName: v.DBInstanceArn,
			})
			if err != nil {
				return nil, err
			}

			var name string
			if v.DBName != nil {
				name = *v.DBName
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    *v.DBInstanceArn,
				Name:   name,
				Description: model.NeptuneDatabaseDescription{
					Database: v,
					Tags:     tags.TagList,
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

func NeptuneDatabaseCluster(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := neptune.NewFromConfig(cfg)
	paginator := neptune.NewDescribeDBClustersPaginator(client, &neptune.DescribeDBClustersInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.DBClusters {
			if v.DBClusterArn == nil {
				continue
			}
			if *v.Engine != "neptune" {
				continue
			}
			tags, err := client.ListTagsForResource(ctx, &neptune.ListTagsForResourceInput{
				ResourceName: v.DBClusterArn,
			})
			if err != nil {
				return nil, err
			}

			var name string
			if v.DBClusterIdentifier != nil {
				name = *v.DBClusterIdentifier
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    *v.DBClusterArn,
				Name:   name,
				Description: model.NeptuneDatabaseClusterDescription{
					Cluster: v,
					Tags:    tags.TagList,
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

func NeptuneDatabaseClusterSnapshot(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := neptune.NewFromConfig(cfg)
	paginator := neptune.NewDescribeDBClustersPaginator(client, &neptune.DescribeDBClustersInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.DBClusters {
			if v.DBClusterArn == nil {
				continue
			}
			input := &neptune.DescribeDBClusterSnapshotsInput{
				DBClusterIdentifier: v.DBClusterIdentifier,
			}
			paginator2 := neptune.NewDescribeDBClusterSnapshotsPaginator(client, input)

			for paginator2.HasMorePages() {
				output, err := paginator2.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				for _, item := range output.DBClusterSnapshots {
					if *item.Engine != "neptune" {
						continue
					}
					resource, err := neptuneDatabaseClusterSnapshotHandler(ctx, client, item)
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
		}
	}

	return values, nil
}

func neptuneDatabaseClusterSnapshotHandler(ctx context.Context, client *neptune.Client, snapshot types.DBClusterSnapshot) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)

	params := &neptune.DescribeDBClusterSnapshotAttributesInput{
		DBClusterSnapshotIdentifier: snapshot.DBClusterSnapshotIdentifier,
	}

	dbClusterSnapshotData, err := client.DescribeDBClusterSnapshotAttributes(ctx, params)
	if err != nil {
		return Resource{}, err
	}

	var attributes = make([]map[string]interface{}, 0)

	if dbClusterSnapshotData.DBClusterSnapshotAttributesResult != nil {

		for _, attribute := range dbClusterSnapshotData.DBClusterSnapshotAttributesResult.DBClusterSnapshotAttributes {
			var result = make(map[string]interface{})

			result["AttributeName"] = attribute.AttributeName
			if len(attribute.AttributeValues) == 0 {
				result["AttributeValues"] = nil
			} else {
				result["AttributeValues"] = attribute.AttributeValues
			}

			attributes = append(attributes, result)

		}
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *snapshot.DBClusterSnapshotArn,
		Name:   *snapshot.DBClusterSnapshotIdentifier,
		ID:     *snapshot.DBClusterSnapshotIdentifier,
		Description: model.NeptuneDatabaseClusterSnapshotDescription{
			Snapshot:   snapshot,
			Attributes: attributes,
		},
	}
	return resource, nil
}
