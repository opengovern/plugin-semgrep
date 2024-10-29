package describer

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/aws/aws-sdk-go-v2/service/athena/types"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func AthenaWrokgroup(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := athena.NewFromConfig(cfg)
	pager := athena.NewListWorkGroupsPaginator(client, &athena.ListWorkGroupsInput{})
	var values []Resource
	for pager.HasMorePages() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, item := range page.WorkGroups {
			resource, err := authenaWorkgroupHandle(ctx, cfg, item)
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

func authenaWorkgroupHandle(ctx context.Context, cfg aws.Config, item types.WorkGroupSummary) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := athena.NewFromConfig(cfg)
	output, err := client.GetWorkGroup(ctx, &athena.GetWorkGroupInput{
		WorkGroup: item.Name,
	})
	if err != nil {
		if isErr(err, "GetWorkGroupNotFound") || isErr(err, "InvalidParameterValue") {
			return Resource{}, nil
		}
		return Resource{}, err
	}
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		Name:   *output.WorkGroup.Name,
		Description: model.AthenaWorkGroupDescription{
			WorkGroup: output.WorkGroup,
		},
	}

	return resource, nil
}

func AthenaQueryExecution(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := athena.NewFromConfig(cfg)
	pager := athena.NewListQueryExecutionsPaginator(client, &athena.ListQueryExecutionsInput{})
	var values []Resource
	for pager.HasMorePages() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, item := range page.QueryExecutionIds {
			resource, err := authenaQueryExecutionHandle(ctx, cfg, item)
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

func authenaQueryExecutionHandle(ctx context.Context, cfg aws.Config, id string) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := athena.NewFromConfig(cfg)
	output, err := client.GetQueryExecution(ctx, &athena.GetQueryExecutionInput{
		QueryExecutionId: aws.String(id),
	})
	if err != nil {
		if isErr(err, "GetQueryExecutionNotFound") || isErr(err, "InvalidParameterValue") {
			return Resource{}, nil
		}
		return Resource{}, err
	}
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ID:     *output.QueryExecution.QueryExecutionId,
		Name:   *output.QueryExecution.Query,
		Description: model.AthenaQueryExecutionDescription{
			QueryExecution: output.QueryExecution,
		},
	}

	return resource, nil
}
