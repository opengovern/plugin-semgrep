package describer

import (
	"context"
	types2 "github.com/aws/aws-sdk-go-v2/service/docdb/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	_ "google.golang.org/genproto/googleapis/bigtable/admin/cluster/v1"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/docdb"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func DocDBCluster(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := docdb.NewFromConfig(cfg)
	paginator := docdb.NewDescribeDBClustersPaginator(client, &docdb.DescribeDBClustersInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, cluster := range page.DBClusters {
			resource, err := DocDBClusterHandle(ctx, cfg, cluster)
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
func DocDBClusterHandle(ctx context.Context, cfg aws.Config, cluster types2.DBCluster) (Resource, error) {
	client := docdb.NewFromConfig(cfg)
	describeCtx := GetDescribeContext(ctx)

	tags, err := client.ListTagsForResource(ctx, &docdb.ListTagsForResourceInput{
		ResourceName: cluster.DBClusterArn,
	})
	if err != nil {
		if isErr(err, "ListTagsForResourceNotFound") || isErr(err, "InvalidParameterValue") {
			return Resource{}, nil
		}
		return Resource{}, err
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ID:     *cluster.DBClusterIdentifier,
		ARN:    *cluster.DBClusterArn,
		Description: model.DocDBClusterDescription{
			DBCluster: cluster,
			Tags:      tags.TagList,
		},
	}
	return resource, nil
}
func GetDocDBCluster(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	client := docdb.NewFromConfig(cfg)
	dbClusterIdentifier := fields["identifier"]

	out, err := client.DescribeDBClusters(ctx, &docdb.DescribeDBClustersInput{
		DBClusterIdentifier: &dbClusterIdentifier,
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, cluster := range out.DBClusters {

		resource, err := DocDBClusterHandle(ctx, cfg, cluster)
		if err != nil {
			return nil, err
		}
		emptyResource := Resource{}
		if err == nil && resource == emptyResource {
			return nil, nil
		}

		values = append(values, resource)
	}
	return values, nil
}

func DocDBClusterInstance(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := docdb.NewFromConfig(cfg)
	paginator := docdb.NewDescribeDBInstancesPaginator(client, &docdb.DescribeDBInstancesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, cluster := range page.DBInstances {
			resource, err := DocDBClusterInstanceHandle(ctx, cfg, cluster)
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

func DocDBClusterInstanceHandle(ctx context.Context, cfg aws.Config, instance types2.DBInstance) (Resource, error) {
	client := docdb.NewFromConfig(cfg)
	describeCtx := GetDescribeContext(ctx)

	tags, err := client.ListTagsForResource(ctx, &docdb.ListTagsForResourceInput{
		ResourceName: instance.DBInstanceArn,
	})
	if err != nil {
		if isErr(err, "ListTagsForResourceNotFound") || isErr(err, "InvalidParameterValue") {
			return Resource{}, nil
		}
		return Resource{}, err
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ID:     *instance.DBInstanceIdentifier,
		ARN:    *instance.DBInstanceArn,
		Description: model.DocDBClusterInstanceDescription{
			DBInstance: instance,
			Tags:       tags.TagList,
		},
	}
	return resource, nil
}

func GetDocDBClusterInstance(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	client := docdb.NewFromConfig(cfg)
	identifier := fields["identifier"]

	out, err := client.DescribeDBInstances(ctx, &docdb.DescribeDBInstancesInput{
		DBInstanceIdentifier: &identifier,
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, cluster := range out.DBInstances {

		resource, err := DocDBClusterInstanceHandle(ctx, cfg, cluster)
		if err != nil {
			return nil, err
		}
		emptyResource := Resource{}
		if err == nil && resource == emptyResource {
			return nil, nil
		}

		values = append(values, resource)
	}
	return values, nil
}

func DocDBClusterSnapshot(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := docdb.NewFromConfig(cfg)

	paginator := docdb.NewDescribeDBClustersPaginator(client, &docdb.DescribeDBClustersInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, cluster := range page.DBClusters {
			paginator2 := docdb.NewDescribeDBClusterSnapshotsPaginator(client, &docdb.DescribeDBClusterSnapshotsInput{
				DBClusterSnapshotIdentifier: cluster.DBClusterIdentifier,
			})

			for paginator2.HasMorePages() {
				page2, err := paginator2.NextPage(ctx)
				if err != nil && !isErr(err, "DBClusterSnapshotNotFoundFault") {
					return nil, err
				} else if err != nil && !isErr(err, "DBClusterSnapshotNotFoundFault") {
					continue
				}
				for _, snapshot := range page2.DBClusterSnapshots {
					resource, err := DocDBClusterSnapshotHandle(ctx, client, cfg, snapshot)
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

func DocDBClusterSnapshotHandle(ctx context.Context, docdbClient *docdb.Client, cfg aws.Config, snapshot types2.DBClusterSnapshot) (Resource, error) {
	client := docdb.NewFromConfig(cfg)
	describeCtx := GetDescribeContext(ctx)

	tags, err := client.ListTagsForResource(ctx, &docdb.ListTagsForResourceInput{
		ResourceName: snapshot.DBClusterSnapshotArn,
	})
	if err != nil {
		if isErr(err, "ListTagsForResourceNotFound") || isErr(err, "InvalidParameterValue") {
			return Resource{}, nil
		}
		return Resource{}, err
	}

	params := &docdb.DescribeDBClusterSnapshotAttributesInput{
		DBClusterSnapshotIdentifier: snapshot.DBClusterSnapshotIdentifier,
	}

	dbClusterSnapshotData, err := docdbClient.DescribeDBClusterSnapshotAttributes(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_docdb_cluster_snapshot.getAwsDocDBClusterSnapshotAttributes", "api_error", err)
		return Resource{}, err
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ID:     *snapshot.DBClusterSnapshotIdentifier,
		ARN:    *snapshot.DBClusterSnapshotArn,
		Description: model.DocDBClusterSnapshotDescription{
			DBClusterSnapshot: snapshot,
			Tags:              tags.TagList,
			Attributes:        dbClusterSnapshotData.DBClusterSnapshotAttributesResult,
		},
	}
	return resource, nil
}
