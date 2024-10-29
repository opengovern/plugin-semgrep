package describer

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/workspaces"
	"github.com/aws/aws-sdk-go-v2/service/workspaces/types"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func WorkSpacesConnectionAlias(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := workspaces.NewFromConfig(cfg)

	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		output, err := client.DescribeConnectionAliases(ctx, &workspaces.DescribeConnectionAliasesInput{
			NextToken: prevToken,
		})
		if err != nil {
			return nil, err
		}

		for _, v := range output.ConnectionAliases {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ID:          *v.AliasId,
				Name:        *v.AliasId,
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

		return output.NextToken, nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}

func WorkspacesWorkspace(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := workspaces.NewFromConfig(cfg)
	paginator := workspaces.NewDescribeWorkspacesPaginator(client, &workspaces.DescribeWorkspacesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if !isErr(err, "ValidationException") {
				return nil, err
			}
			continue
		}

		for _, v := range page.Workspaces {
			resource, err := workspacesWorkspaceHandle(ctx, cfg, v)
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
func workspacesWorkspaceHandle(ctx context.Context, cfg aws.Config, v types.Workspace) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := workspaces.NewFromConfig(cfg)

	arn := fmt.Sprintf("arn:%s:workspaces:%s:%s:workspace/%s", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID, *v.WorkspaceId)
	tags, err := client.DescribeTags(ctx, &workspaces.DescribeTagsInput{
		ResourceId: v.WorkspaceId,
	})
	if err != nil {
		if !isErr(err, "ValidationException") {
			return Resource{}, err
		}
		tags = &workspaces.DescribeTagsOutput{}
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *v.WorkspaceId,
		Description: model.WorkspacesWorkspaceDescription{
			Workspace: v,
			Tags:      tags.TagList,
		},
	}
	return resource, nil
}
func GetWorkspacesWorkspace(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	client := workspaces.NewFromConfig(cfg)
	workspaceId := fields["workspaceId"]
	var values []Resource
	workspace, err := client.DescribeWorkspaces(ctx, &workspaces.DescribeWorkspacesInput{
		WorkspaceIds: []string{workspaceId},
	})
	if err != nil {
		if isErr(err, "DescribeWorkspacesNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	for _, v := range workspace.Workspaces {
		resource, err := workspacesWorkspaceHandle(ctx, cfg, v)
		if err != nil {
			return nil, err
		}

		values = append(values, resource)
	}
	return values, nil
}

func WorkspacesBundle(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := workspaces.NewFromConfig(cfg)
	paginator := workspaces.NewDescribeWorkspaceBundlesPaginator(client, &workspaces.DescribeWorkspaceBundlesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Bundles {
			tags, err := client.DescribeTags(ctx, &workspaces.DescribeTagsInput{
				ResourceId: v.BundleId,
			})
			if err != nil {
				return nil, err
			}

			resource := workspacesBundleHandle(ctx, v, tags)
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
func workspacesBundleHandle(ctx context.Context, v types.WorkspaceBundle, tags *workspaces.DescribeTagsOutput) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := fmt.Sprintf("arn:%s:workspaces:%s:%s:workspacebundle/%s", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID, *v.BundleId)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *v.BundleId,
		Description: model.WorkspacesBundleDescription{
			Bundle: v,
			Tags:   tags.TagList,
		},
	}
	return resource
}
func GetWorkspacesBundle(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	var values []Resource
	BundleId := fields["bundleId"]
	client := workspaces.NewFromConfig(cfg)

	workspace, err := client.DescribeWorkspaceBundles(ctx, &workspaces.DescribeWorkspaceBundlesInput{
		BundleIds: []string{BundleId},
	})
	if err != nil {
		if isErr(err, "DescribeWorkspaceBundlesNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	for _, v := range workspace.Bundles {
		tags, err := client.DescribeTags(ctx, &workspaces.DescribeTagsInput{
			ResourceId: v.BundleId,
		})
		if err != nil {
			if isErr(err, "DescribeTagsNotFound") || isErr(err, "InvalidParameterValue") {
				return nil, nil
			}
			return nil, err
		}

		values = append(values, workspacesBundleHandle(ctx, v, tags))
	}
	return values, nil
}
