package describer

import (
	"context"
	"fmt"
	"math"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/directconnect"
	"github.com/aws/aws-sdk-go-v2/service/directconnect/types"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func DirectConnectConnection(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := directconnect.NewFromConfig(cfg)
	connections, err := client.DescribeConnections(ctx, &directconnect.DescribeConnectionsInput{})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range connections.Connections {
		resource := directConnectConnectionHandle(ctx, v)
		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return nil, err
			}
		} else {
			values = append(values, resource)
		}
	}

	return values, nil
}
func directConnectConnectionHandle(ctx context.Context, v types.Connection) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := fmt.Sprintf("arn:%s:directconnect:%s:%s:dxcon/%s", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID, *v.ConnectionId)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *v.ConnectionId,
		Description: model.DirectConnectConnectionDescription{
			Connection: v,
		},
	}
	return resource
}
func GetDirectConnectConnection(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	connectionId := fields["id"]
	client := directconnect.NewFromConfig(cfg)
	out, err := client.DescribeConnections(ctx, &directconnect.DescribeConnectionsInput{
		ConnectionId: &connectionId,
	})
	if err != nil {
		return nil, err
	}
	var values []Resource
	for _, v := range out.Connections {
		resource := directConnectConnectionHandle(ctx, v)
		values = append(values, resource)
	}
	return values, nil
}

func getDirectConnectGatewayArn(describeCtx DescribeContext, directConnectGatewayId string) string {
	return fmt.Sprintf("arn:%s:directconnect::%s:dx-gateway/%s", describeCtx.Partition, describeCtx.AccountID, directConnectGatewayId)
}

func DirectConnectGateway(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := directconnect.NewFromConfig(cfg)

	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		connections, err := client.DescribeDirectConnectGateways(ctx, &directconnect.DescribeDirectConnectGatewaysInput{
			MaxResults: aws.Int32(100),
			NextToken:  prevToken,
		})
		if err != nil {
			return nil, err
		}
		if len(connections.DirectConnectGateways) == 0 {
			return nil, nil
		}
		arns := make([]string, 0, len(connections.DirectConnectGateways))
		for _, v := range connections.DirectConnectGateways {
			arns = append(arns, getDirectConnectGatewayArn(describeCtx, *v.DirectConnectGatewayId))
		}
		// DescribeTags can only handle 20 ARNs at a time
		arnToTagMap := make(map[string][]types.Tag)
		for i := 0; i < len(arns); i += 20 {
			tags, err := client.DescribeTags(ctx, &directconnect.DescribeTagsInput{
				ResourceArns: arns[i:int(math.Min(float64(i+20), float64(len(arns))))],
			})
			if err != nil {
				tags = &directconnect.DescribeTagsOutput{}
			}

			for _, tag := range tags.ResourceTags {
				arnToTagMap[*tag.ResourceArn] = tag.Tags
			}
		}

		for _, v := range connections.DirectConnectGateways {
			resource := directConnectGatewayHandle(ctx, v, arnToTagMap)
			if stream != nil {
				if err := (*stream)(resource); err != nil {
					return nil, err
				}
			} else {
				values = append(values, resource)
			}
		}

		return connections.NextToken, nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}
func directConnectGatewayHandle(ctx context.Context, v types.DirectConnectGateway, arnToTagMap map[string][]types.Tag) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := getDirectConnectGatewayArn(describeCtx, *v.DirectConnectGatewayId)

	tagsList, ok := arnToTagMap[arn]
	if !ok {
		tagsList = []types.Tag{}
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *v.DirectConnectGatewayName,
		Description: model.DirectConnectGatewayDescription{
			Gateway: v,
			Tags:    tagsList,
		},
	}
	return resource
}
func GetDirectConnectGateway(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	DirectConnectGatewayId := fields["id"]

	client := directconnect.NewFromConfig(cfg)
	out, err := client.DescribeDirectConnectGateways(ctx, &directconnect.DescribeDirectConnectGatewaysInput{
		DirectConnectGatewayId: &DirectConnectGatewayId,
	})
	if err != nil {
		if isErr(err, "DirectConnectGatewayNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}
	if len(out.DirectConnectGateways) == 0 {
		return nil, nil
	}
	arns := make([]string, 0, len(out.DirectConnectGateways))
	for _, v := range out.DirectConnectGateways {
		arns = append(arns, getDirectConnectGatewayArn(describeCtx, *v.DirectConnectGatewayId))
	}
	// DescribeTags can only handle 20 ARNs at a time
	arnToTagMap := make(map[string][]types.Tag)
	for i := 0; i < len(arns); i += 20 {
		tags, err := client.DescribeTags(ctx, &directconnect.DescribeTagsInput{
			ResourceArns: arns[i:int(math.Min(float64(i+20), float64(len(arns))))],
		})
		if err != nil {
			tags = &directconnect.DescribeTagsOutput{}
		}

		for _, tag := range tags.ResourceTags {
			arnToTagMap[*tag.ResourceArn] = tag.Tags
		}
	}
	var values []Resource

	for _, v := range out.DirectConnectGateways {
		resource := directConnectGatewayHandle(ctx, v, arnToTagMap)
		values = append(values, resource)
	}
	return values, nil
}
