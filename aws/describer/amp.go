package describer

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/amp/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/amp"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func AMPWorkspace(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := amp.NewFromConfig(cfg)
	paginator := amp.NewListWorkspacesPaginator(client, &amp.ListWorkspacesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Workspaces {
			resource := aMPWorkspaceHandle(ctx, v)

			if stream != nil {
				m := *stream
				err := m(resource)
				if err != nil {
					return nil, err
				}
			} else {
				values = append(values, resource)
			}
		}
	}

	return values, nil
}
func aMPWorkspaceHandle(ctx context.Context, v types.WorkspaceSummary) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *v.Arn,
		Name:   *v.WorkspaceId,
		Description: model.AMPWorkspaceDescription{
			Workspace: v,
		},
	}
	return resource
}
func GetAMPWorkspace(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	workspaceID := fields["id"]
	client := amp.NewFromConfig(cfg)

	out, err := client.ListWorkspaces(ctx, &amp.ListWorkspacesInput{})
	if err != nil {
		if isErr(err, "ListWorkspacesNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	for _, v := range out.Workspaces {
		if *v.WorkspaceId != workspaceID {
			continue
		}
		resource := aMPWorkspaceHandle(ctx, v)
		values = append(values, resource)
	}
	return values, nil
}
