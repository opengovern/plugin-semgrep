package describer

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/grafana"
	"github.com/aws/aws-sdk-go-v2/service/grafana/types"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func GrafanaWorkspace(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := grafana.NewFromConfig(cfg)
	paginator := grafana.NewListWorkspacesPaginator(client, &grafana.ListWorkspacesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Workspaces {
			resource := grafanaWorkspaceHandle(ctx, v)
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
func grafanaWorkspaceHandle(ctx context.Context, v types.WorkspaceSummary) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := fmt.Sprintf("arn:%s:grafana:%s:%s:/workspaces/%s", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID, *v.Id)

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *v.Id,
		Description: model.GrafanaWorkspaceDescription{
			Workspace: v,
		},
	}
	return resource
}
func GetGrafanaWorkspace(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	workspaceId := fields["workspaceId"]
	client := grafana.NewFromConfig(cfg)

	out, err := client.DescribeWorkspace(ctx, &grafana.DescribeWorkspaceInput{
		WorkspaceId: &workspaceId,
	})
	if err != nil {
		if isErr(err, "DescribeWorkspaceNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	workspace := types.WorkspaceSummary{
		Id:                       out.Workspace.Id,
		Name:                     out.Workspace.Name,
		Endpoint:                 out.Workspace.Endpoint,
		Description:              out.Workspace.Description,
		Status:                   out.Workspace.Status,
		Modified:                 out.Workspace.Modified,
		GrafanaVersion:           out.Workspace.GrafanaVersion,
		Created:                  out.Workspace.Created,
		Authentication:           out.Workspace.Authentication,
		NotificationDestinations: out.Workspace.NotificationDestinations,
		Tags:                     out.Workspace.Tags,
	}

	var values []Resource
	values = append(values, grafanaWorkspaceHandle(ctx, workspace))
	return values, nil
}
