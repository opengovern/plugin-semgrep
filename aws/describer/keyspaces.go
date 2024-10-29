package describer

import (
	"context"
	_ "database/sql/driver"
	"github.com/aws/aws-sdk-go-v2/service/keyspaces/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/keyspaces"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func KeyspacesKeyspace(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := keyspaces.NewFromConfig(cfg)
	paginator := keyspaces.NewListKeyspacesPaginator(client, &keyspaces.ListKeyspacesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Keyspaces {
			resource, err := keyspacesKeyspaceHandle(ctx, cfg, v)
			emptyResource := Resource{}
			if err == nil && resource == emptyResource {
				return nil, nil
			}
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
func keyspacesKeyspaceHandle(ctx context.Context, cfg aws.Config, v types.KeyspaceSummary) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := keyspaces.NewFromConfig(cfg)

	tags, err := client.ListTagsForResource(ctx, &keyspaces.ListTagsForResourceInput{
		ResourceArn: v.ResourceArn,
	})
	if err != nil {
		if isErr(err, "ListTagsForResourceNotFound") || isErr(err, "InvalidParameterValue") {
			return Resource{}, nil
		}
		return Resource{}, err
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *v.ResourceArn,
		Name:   *v.KeyspaceName,
		Description: model.KeyspacesKeyspaceDescription{
			Keyspace: v,
			Tags:     tags.Tags,
		},
	}
	return resource, nil
}
func GetKeyspacesKeyspace(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	name := fields["name"]
	var values []Resource
	client := keyspaces.NewFromConfig(cfg)
	out, err := client.GetKeyspace(ctx, &keyspaces.GetKeyspaceInput{
		KeyspaceName: &name,
	})
	if err != nil {
		if isErr(err, "GetKeyspaceNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	keyspace := types.KeyspaceSummary{
		KeyspaceName: out.KeyspaceName,
		ResourceArn:  out.ResourceArn,
	}

	resource, err := keyspacesKeyspaceHandle(ctx, cfg, keyspace)
	emptyResource := Resource{}
	if err == nil && resource == emptyResource {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	values = append(values, resource)
	return values, nil
}

func KeyspacesTable(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := keyspaces.NewFromConfig(cfg)
	keyspacePaginator := keyspaces.NewListKeyspacesPaginator(client, &keyspaces.ListKeyspacesInput{})

	var values []Resource
	for keyspacePaginator.HasMorePages() {
		keyspacePage, err := keyspacePaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, keyspace := range keyspacePage.Keyspaces {
			paginator := keyspaces.NewListTablesPaginator(client, &keyspaces.ListTablesInput{
				KeyspaceName: keyspace.KeyspaceName,
			})

			for paginator.HasMorePages() {
				page, err := paginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				for _, v := range page.Tables {
					resource, err := keyspacesTableHandle(ctx, cfg, v)
					emptyResource := Resource{}
					if err == nil && resource == emptyResource {
						return nil, nil
					}
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
func keyspacesTableHandle(ctx context.Context, cfg aws.Config, v types.TableSummary) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := keyspaces.NewFromConfig(cfg)

	tags, err := client.ListTagsForResource(ctx, &keyspaces.ListTagsForResourceInput{
		ResourceArn: v.ResourceArn,
	})
	if err != nil {
		if isErr(err, "ListTagsForResourceNotFound") || isErr(err, "InvalidParameterValue") {
			return Resource{}, nil
		}
		return Resource{}, err
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ID:     *v.ResourceArn,
		Name:   *v.KeyspaceName,
		Description: model.KeyspacesTableDescription{
			Table: v,
			Tags:  tags.Tags,
		},
	}
	return resource, nil
}
func GetKeyspacesTable(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	var values []Resource
	name := fields["name"]
	client := keyspaces.NewFromConfig(cfg)

	list, err := client.ListTables(ctx, &keyspaces.ListTablesInput{
		KeyspaceName: &name,
	})
	if err != nil {
		if isErr(err, "ListTablesNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	for _, v := range list.Tables {
		resource, err := keyspacesTableHandle(ctx, cfg, v)
		emptyResource := Resource{}
		if err == nil && resource == emptyResource {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}

		values = append(values, resource)
	}
	return values, nil
}
