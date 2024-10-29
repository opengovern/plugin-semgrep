package describer

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/storagegateway/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/storagegateway"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func StorageGatewayStorageGateway(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := storagegateway.NewFromConfig(cfg)
	paginator := storagegateway.NewListGatewaysPaginator(client, &storagegateway.ListGatewaysInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Gateways {
			resource, err := storageGatewayStorageGatewayHandle(ctx, cfg, *v.GatewayARN, *v.GatewayId, v)
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
func storageGatewayStorageGatewayHandle(ctx context.Context, cfg aws.Config, gatewayARN string, gatewayId string, v types.GatewayInfo) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := storagegateway.NewFromConfig(cfg)
	tags, err := client.ListTagsForResource(ctx, &storagegateway.ListTagsForResourceInput{
		ResourceARN: &gatewayARN,
	})
	if err != nil {
		return Resource{}, err
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    gatewayARN,
		Name:   gatewayId,
		Description: model.StorageGatewayStorageGatewayDescription{
			StorageGateway: v,
			Tags:           tags.Tags,
		},
	}
	return resource, nil
}
func GetStorageGatewayStorageGateway(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	gatewayArn := fields["arn"]
	client := storagegateway.NewFromConfig(cfg)
	out, err := client.DescribeGatewayInformation(ctx, &storagegateway.DescribeGatewayInformationInput{
		GatewayARN: &gatewayArn,
	})
	if err != nil {
		return nil, err
	}

	storageGateway := types.GatewayInfo{
		Ec2InstanceId:     out.Ec2InstanceId,
		Ec2InstanceRegion: out.Ec2InstanceRegion,
		GatewayARN:        out.GatewayARN,
		GatewayId:         out.GatewayId,
		GatewayName:       out.GatewayName,
		//GatewayOperationalState: out.GatewayOperationalState, //TODO-Saleh
		GatewayType:       out.GatewayType,
		HostEnvironment:   out.HostEnvironment,
		HostEnvironmentId: out.HostEnvironmentId,
	}

	var values []Resource
	resource, err := storageGatewayStorageGatewayHandle(ctx, cfg, *out.GatewayARN, *out.GatewayId, storageGateway)
	if err != nil {
		return nil, err
	}

	values = append(values, resource)
	return values, nil
}
