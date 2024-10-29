package describer

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/opsworkscm"
	"github.com/aws/aws-sdk-go-v2/service/opsworkscm/types"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func OpsWorksCMServer(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := opsworkscm.NewFromConfig(cfg)
	paginator := opsworkscm.NewDescribeServersPaginator(client, &opsworkscm.DescribeServersInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Servers {
			resource, err := opsWorksCMServerHandle(ctx, cfg, v)
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
func opsWorksCMServerHandle(ctx context.Context, cfg aws.Config, v types.Server) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := opsworkscm.NewFromConfig(cfg)

	tags, err := client.ListTagsForResource(ctx, &opsworkscm.ListTagsForResourceInput{
		ResourceArn: v.ServerArn,
	})
	if err != nil {
		if isErr(err, "ListTagsForResourceNotFound") || isErr(err, "InvalidParameterValue") {
			return Resource{}, nil
		}
		return Resource{}, err
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *v.ServerArn,
		Name:   *v.ServerName,
		Description: model.OpsWorksCMServerDescription{
			Server: v,
			Tags:   tags.Tags,
		},
	}
	return resource, nil
}
func GetOpsWorksCMServer(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	serverName := fields["name"]
	client := opsworkscm.NewFromConfig(cfg)

	server, err := client.DescribeServers(ctx, &opsworkscm.DescribeServersInput{
		ServerName: &serverName,
	})
	if err != nil {
	}

	var values []Resource
	for _, v := range server.Servers {

		resource, err := opsWorksCMServerHandle(ctx, cfg, v)
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
