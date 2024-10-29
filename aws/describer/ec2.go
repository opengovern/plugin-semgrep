package describer

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/smithy-go"
	_ "github.com/aws/smithy-go"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func EC2ElasticIP(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	var values []Resource
	describeCtx := GetDescribeContext(ctx)
	client := ec2.NewFromConfig(cfg)

	addrs, err := client.DescribeAddresses(ctx, &ec2.DescribeAddressesInput{})
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs.Addresses {
		resource := Resource{
			ID:     *addr.AllocationId,
			Region: describeCtx.KaytuRegion,
			Description: model.EC2ElasticIPDescription{
				Address: addr,
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

	return values, nil
}

func EC2LocalGateway(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeLocalGatewaysPaginator(client, &ec2.DescribeLocalGatewaysInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.LocalGateways {
			resource := eC2LocalGatewayHandle(ctx, v)
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
func eC2LocalGatewayHandle(ctx context.Context, v types.LocalGateway) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := "arn:" + describeCtx.Partition + ":ec2:" + describeCtx.Region + ":" + describeCtx.AccountID + ":local-gateway/" + *v.LocalGatewayId
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *v.LocalGatewayId,
		Description: model.EC2LocalGatewayDescription{
			LocalGateway: v,
		},
	}
	return resource
}

func EC2VolumeSnapshot(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	var values []Resource
	client := ec2.NewFromConfig(cfg)
	describeCtx := GetDescribeContext(ctx)

	fmt.Println("+++++++++", describeCtx.Region, cfg.Region)

	paginator := ec2.NewDescribeSnapshotsPaginator(client, &ec2.DescribeSnapshotsInput{
		OwnerIds: []string{"self"},
	})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, snapshot := range page.Snapshots {
			// This prevents Implicit memory aliasing in for loop
			snapshot := snapshot
			attrs, err := client.DescribeSnapshotAttribute(ctx, &ec2.DescribeSnapshotAttributeInput{
				Attribute:  types.SnapshotAttributeNameCreateVolumePermission,
				SnapshotId: snapshot.SnapshotId,
			})
			if err != nil {
				if isErr(err, "InvalidSnapshot.NotFound") {
					continue
				}
				return nil, err
			}

			resource := eC2VolumeSnapshotHandle(ctx, snapshot, attrs)
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
func eC2VolumeSnapshotHandle(ctx context.Context, v types.Snapshot, attrs *ec2.DescribeSnapshotAttributeOutput) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := "arn:" + describeCtx.Partition + ":ec2:" + describeCtx.Region + ":" + describeCtx.AccountID + ":snapshot/" + *v.SnapshotId
	fmt.Println("=======", arn)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *v.SnapshotId,
		Description: model.EC2VolumeSnapshotDescription{
			Snapshot:                &v,
			CreateVolumePermissions: attrs.CreateVolumePermissions,
		},
	}
	return resource
}
func GetEC2VolumeSnapshot(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	VolumeSnapshotId := fields["id"]
	client := ec2.NewFromConfig(cfg)

	out, err := client.DescribeSnapshots(ctx, &ec2.DescribeSnapshotsInput{
		SnapshotIds: []string{VolumeSnapshotId},
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, snapshot := range out.Snapshots {
		attrs, err := client.DescribeSnapshotAttribute(ctx, &ec2.DescribeSnapshotAttributeInput{
			Attribute:  types.SnapshotAttributeNameCreateVolumePermission,
			SnapshotId: snapshot.SnapshotId,
		})
		if err != nil {
			if isErr(err, "InvalidSnapshot.NotFound") {
				continue
			}
			return nil, err
		}
		resource := eC2VolumeSnapshotHandle(ctx, snapshot, attrs)
		values = append(values, resource)
	}

	return values, nil
}

func EC2Volume(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	var values []Resource
	client := ec2.NewFromConfig(cfg)

	paginator := ec2.NewDescribeVolumesPaginator(client, &ec2.DescribeVolumesInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, volume := range page.Volumes {
			var resource Resource
			resource, err = eC2VolumeHandle(ctx, volume, client)
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
func eC2VolumeHandle(ctx context.Context, v types.Volume, client *ec2.Client) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	volume := v
	var description model.EC2VolumeDescription
	description.Volume = &volume

	attrs := []types.VolumeAttributeName{
		types.VolumeAttributeNameAutoEnableIO,
		types.VolumeAttributeNameProductCodes,
	}

	for _, attr := range attrs {
		attrs, err := client.DescribeVolumeAttribute(ctx, &ec2.DescribeVolumeAttributeInput{
			Attribute: attr,
			VolumeId:  volume.VolumeId,
		})
		if err != nil {
			return Resource{}, err
		}

		switch attr {
		case types.VolumeAttributeNameAutoEnableIO:
			description.Attributes.AutoEnableIO = *attrs.AutoEnableIO.Value
		case types.VolumeAttributeNameProductCodes:
			description.Attributes.ProductCodes = attrs.ProductCodes
		}
	}

	arn := "arn:" + describeCtx.Partition + ":ec2:" + describeCtx.Region + ":" + describeCtx.AccountID + ":volume/" + *volume.VolumeId
	resource := Resource{
		Region:      describeCtx.KaytuRegion,
		ARN:         arn,
		Name:        *volume.VolumeId,
		Description: description,
	}
	return resource, nil
}
func GetEC2Volume(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	volumeId := fields["id"]
	client := ec2.NewFromConfig(cfg)

	out, err := client.DescribeVolumes(ctx, &ec2.DescribeVolumesInput{
		VolumeIds: []string{volumeId},
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, volume := range out.Volumes {
		var resource Resource
		resource, err = eC2VolumeHandle(ctx, volume, client)
		if err != nil {
			return nil, err
		}
		values = append(values, resource)
	}
	return values, nil
}

func EC2CapacityReservation(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeCapacityReservationsPaginator(client, &ec2.DescribeCapacityReservationsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if isErr(err, "InvalidCapacityReservationId.NotFound") || isErr(err, "InvalidCapacityReservationId.Unavailable") || isErr(err, "InvalidCapacityReservationId.Malformed") {
				continue
			}
			return nil, err
		}

		for _, v := range page.CapacityReservations {
			resource := eC2CapacityReservationHandle(ctx, v)
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
func eC2CapacityReservationHandle(ctx context.Context, v types.CapacityReservation) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *v.CapacityReservationArn,
		Name:   *v.CapacityReservationId,
		Description: model.EC2CapacityReservationDescription{
			CapacityReservation: v,
		},
	}
	return resource
}
func GetEC2CapacityReservation(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	ReservationId := fields["id"]
	client := ec2.NewFromConfig(cfg)

	out, err := client.DescribeCapacityReservations(ctx, &ec2.DescribeCapacityReservationsInput{
		CapacityReservationIds: []string{ReservationId},
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range out.CapacityReservations {
		resource := eC2CapacityReservationHandle(ctx, v)
		values = append(values, resource)
	}
	return values, nil
}

func EC2CapacityReservationFleet(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeCapacityReservationFleetsPaginator(client, &ec2.DescribeCapacityReservationFleetsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.CapacityReservationFleets {
			resource := eC2CapacityReservationFleetHandle(ctx, v)
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
func eC2CapacityReservationFleetHandle(ctx context.Context, v types.CapacityReservationFleet) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *v.CapacityReservationFleetArn,
		Name:   *v.CapacityReservationFleetId,
		Description: model.EC2CapacityReservationFleetDescription{
			CapacityReservationFleet: v,
		},
	}
	return resource
}
func GetEC2CapacityReservationFleet(ctx context.Context, cfg aws.Config, field map[string]string) ([]Resource, error) {
	CapacityReservationFleetId := field["id"]
	client := ec2.NewFromConfig(cfg)

	out, err := client.DescribeCapacityReservationFleets(ctx, &ec2.DescribeCapacityReservationFleetsInput{
		CapacityReservationFleetIds: []string{CapacityReservationFleetId},
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range out.CapacityReservationFleets {
		resource := eC2CapacityReservationFleetHandle(ctx, v)
		values = append(values, resource)
	}
	return values, nil
}

func EC2CarrierGateway(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeCarrierGatewaysPaginator(client, &ec2.DescribeCarrierGatewaysInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.CarrierGateways {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ID:          *v.CarrierGatewayId,
				Name:        *v.CarrierGatewayId,
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
	}

	return values, nil
}

func EC2ClientVpnAuthorizationRule(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	endpoints, err := EC2ClientVpnEndpoint(ctx, cfg, nil)
	if err != nil {
		return nil, err
	}

	client := ec2.NewFromConfig(cfg)

	var values []Resource
	for _, e := range endpoints {
		endpoint := e.Description.(types.ClientVpnEndpoint)
		paginator := ec2.NewDescribeClientVpnAuthorizationRulesPaginator(client, &ec2.DescribeClientVpnAuthorizationRulesInput{
			ClientVpnEndpointId: endpoint.ClientVpnEndpointId,
		})

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				return nil, err
			}

			for _, v := range page.AuthorizationRules {
				resource := Resource{
					Region:      describeCtx.KaytuRegion,
					ID:          CompositeID(*v.ClientVpnEndpointId, *v.DestinationCidr, *v.GroupId),
					Name:        *v.ClientVpnEndpointId,
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
		}
	}

	return values, nil
}

func EC2ClientVpnEndpoint(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeClientVpnEndpointsPaginator(client, &ec2.DescribeClientVpnEndpointsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if isErr(err, "AccessDeniedException") {
				return nil, nil
			} else {
				return nil, err
			}
		}

		for _, v := range page.ClientVpnEndpoints {
			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ID:     *v.ClientVpnEndpointId,
				Name:   *v.ClientVpnEndpointId,
				Description: model.EC2ClientVpnEndpointDescription{
					ClientVpnEndpoint: v,
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

func EC2ClientVpnRoute(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	endpoints, err := EC2ClientVpnEndpoint(ctx, cfg, nil)
	if err != nil {
		return nil, err
	}

	client := ec2.NewFromConfig(cfg)

	var values []Resource
	for _, e := range endpoints {
		endpoint := e.Description.(types.ClientVpnEndpoint)
		paginator := ec2.NewDescribeClientVpnRoutesPaginator(client, &ec2.DescribeClientVpnRoutesInput{
			ClientVpnEndpointId: endpoint.ClientVpnEndpointId,
		})

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				return nil, err
			}

			for _, v := range page.Routes {
				resource := Resource{
					Region:      describeCtx.KaytuRegion,
					ID:          CompositeID(*v.ClientVpnEndpointId, *v.DestinationCidr, *v.TargetSubnet),
					Name:        *v.ClientVpnEndpointId,
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
		}
	}

	return values, nil
}

func EC2ClientVpnTargetNetworkAssociation(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	endpoints, err := EC2ClientVpnEndpoint(ctx, cfg, nil)
	if err != nil {
		return nil, err
	}

	client := ec2.NewFromConfig(cfg)

	var values []Resource
	for _, e := range endpoints {
		endpoint := e.Description.(types.ClientVpnEndpoint)
		paginator := ec2.NewDescribeClientVpnTargetNetworksPaginator(client, &ec2.DescribeClientVpnTargetNetworksInput{
			ClientVpnEndpointId: endpoint.ClientVpnEndpointId,
		})

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				return nil, err
			}

			for _, v := range page.ClientVpnTargetNetworks {
				resource := Resource{
					Region:      describeCtx.KaytuRegion,
					ID:          *v.AssociationId,
					Name:        *v.AssociationId,
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
		}
	}

	return values, nil
}

func EC2CustomerGateway(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	output, err := client.DescribeCustomerGateways(ctx, &ec2.DescribeCustomerGatewaysInput{})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range output.CustomerGateways {
		resource := eC2CustomerGatewayHandle(ctx, v)
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
func eC2CustomerGatewayHandle(ctx context.Context, v types.CustomerGateway) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ID:     *v.CustomerGatewayId,
		Description: model.EC2CustomerGatewayDescription{
			CustomerGateway: v,
		},
	}
	return resource
}
func GetEC2CustomerGateway(ctx context.Context, cfg aws.Config, field map[string]string) ([]Resource, error) {
	EC2CustomerGatewayId := field["id"]
	client := ec2.NewFromConfig(cfg)

	out, err := client.DescribeCustomerGateways(ctx, &ec2.DescribeCustomerGatewaysInput{
		CustomerGatewayIds: []string{EC2CustomerGatewayId},
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range out.CustomerGateways {
		resource := eC2CustomerGatewayHandle(ctx, v)
		values = append(values, resource)
	}
	return values, nil
}

func EC2VerifiedAccessInstance(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)

	input := &ec2.DescribeVerifiedAccessInstancesInput{}

	var values []Resource
	for {
		resp, err := client.DescribeVerifiedAccessInstances(ctx, input)
		if err != nil {
			return nil, nil
		}

		for _, instance := range resp.VerifiedAccessInstances {
			resource := eC2VerifiedAccessInstanceHandle(ctx, instance)
			if stream != nil {
				if err := (*stream)(resource); err != nil {
					return nil, err
				}
			} else {
				values = append(values, resource)
			}
		}
		if resp.NextToken == nil {
			break
		} else {
			input.NextToken = resp.NextToken
		}
	}

	return values, nil
}
func eC2VerifiedAccessInstanceHandle(ctx context.Context, v types.VerifiedAccessInstance) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ID:     *v.VerifiedAccessInstanceId,
		Name:   *v.VerifiedAccessInstanceId,
		Description: model.EC2VerifiedAccessInstanceDescription{
			VerifiedAccountInstance: v,
		},
	}
	return resource
}
func GetEC2VerifiedAccessInstance(ctx context.Context, cfg aws.Config, field map[string]string) ([]Resource, error) {
	EC2VerifiedAccessInstanceId := field["id"]
	client := ec2.NewFromConfig(cfg)

	out, err := client.DescribeVerifiedAccessInstances(ctx, &ec2.DescribeVerifiedAccessInstancesInput{
		VerifiedAccessInstanceIds: []string{EC2VerifiedAccessInstanceId},
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range out.VerifiedAccessInstances {
		resource := eC2VerifiedAccessInstanceHandle(ctx, v)
		values = append(values, resource)
	}

	return values, nil
}

func EC2VerifiedAccessEndpoint(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)

	input := &ec2.DescribeVerifiedAccessEndpointsInput{}

	var values []Resource
	for {
		resp, err := client.DescribeVerifiedAccessEndpoints(ctx, input)
		if err != nil {
			return nil, nil
		}

		for _, instance := range resp.VerifiedAccessEndpoints {
			resource := eC2VerifiedAccessEndpointHandle(ctx, instance)
			if stream != nil {
				if err := (*stream)(resource); err != nil {
					return nil, err
				}
			} else {
				values = append(values, resource)
			}
		}
		if resp.NextToken == nil {
			break
		} else {
			input.NextToken = resp.NextToken
		}
	}

	return values, nil
}
func eC2VerifiedAccessEndpointHandle(ctx context.Context, v types.VerifiedAccessEndpoint) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ID:     *v.VerifiedAccessEndpointId,
		Name:   *v.VerifiedAccessEndpointId,
		Description: model.EC2VerifiedAccessEndpointDescription{
			VerifiedAccountEndpoint: v,
		},
	}
	return resource
}
func GetEC2VerifiedAccessEndpoint(ctx context.Context, cfg aws.Config, field map[string]string) ([]Resource, error) {
	EC2VerifiedAccessEndpointId := field["id"]
	client := ec2.NewFromConfig(cfg)

	out, err := client.DescribeVerifiedAccessEndpoints(ctx, &ec2.DescribeVerifiedAccessEndpointsInput{
		VerifiedAccessEndpointIds: []string{EC2VerifiedAccessEndpointId},
	})

	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range out.VerifiedAccessEndpoints {
		resource := eC2VerifiedAccessEndpointHandle(ctx, v)
		values = append(values, resource)
	}

	return values, nil
}

func EC2VerifiedAccessGroup(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)

	input := &ec2.DescribeVerifiedAccessGroupsInput{}

	var values []Resource
	for {
		resp, err := client.DescribeVerifiedAccessGroups(ctx, input)
		if err != nil {
			return nil, nil
		}

		for _, instance := range resp.VerifiedAccessGroups {
			resource := eC2VerifiedAccessGroupHandle(ctx, instance)
			if stream != nil {
				if err := (*stream)(resource); err != nil {
					return nil, err
				}
			} else {
				values = append(values, resource)
			}
		}
		if resp.NextToken == nil {
			break
		} else {
			input.NextToken = resp.NextToken
		}
	}

	return values, nil
}
func eC2VerifiedAccessGroupHandle(ctx context.Context, v types.VerifiedAccessGroup) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ID:     *v.VerifiedAccessGroupId,
		Name:   *v.VerifiedAccessGroupId,
		Description: model.EC2VerifiedAccessGroupDescription{
			VerifiedAccountGroup: v,
		},
	}
	return resource
}
func GetEC2VerifiedAccessGroup(ctx context.Context, cfg aws.Config, field map[string]string) ([]Resource, error) {
	EC2VerifiedAccessGroupId := field["group_id"]
	client := ec2.NewFromConfig(cfg)

	out, err := client.DescribeVerifiedAccessGroups(ctx, &ec2.DescribeVerifiedAccessGroupsInput{
		VerifiedAccessGroupIds: []string{EC2VerifiedAccessGroupId},
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range out.VerifiedAccessGroups {
		resource := eC2VerifiedAccessGroupHandle(ctx, v)
		values = append(values, resource)
	}

	return values, nil
}

func EC2VerifiedAccessTrustProvider(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)

	input := &ec2.DescribeVerifiedAccessTrustProvidersInput{}

	var values []Resource
	for {
		resp, err := client.DescribeVerifiedAccessTrustProviders(ctx, input)
		if err != nil {
			return nil, nil
		}

		for _, instance := range resp.VerifiedAccessTrustProviders {
			resource := eC2VerifiedAccessTrustProviderHandle(ctx, instance)
			if stream != nil {
				if err := (*stream)(resource); err != nil {
					return nil, err
				}
			} else {
				values = append(values, resource)
			}
		}
		if resp.NextToken == nil {
			break
		} else {
			input.NextToken = resp.NextToken
		}
	}

	return values, nil
}
func eC2VerifiedAccessTrustProviderHandle(ctx context.Context, v types.VerifiedAccessTrustProvider) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ID:     *v.VerifiedAccessTrustProviderId,
		Name:   *v.VerifiedAccessTrustProviderId,
		Description: model.EC2VerifiedAccessTrustProviderDescription{
			VerifiedAccessTrustProvider: v,
		},
	}
	return resource
}
func GetEC2VerifiedAccessTrustProvider(ctx context.Context, cfg aws.Config, field map[string]string) ([]Resource, error) {
	EC2VerifiedAccessTrustProviderId := field["id"]
	client := ec2.NewFromConfig(cfg)

	out, err := client.DescribeVerifiedAccessTrustProviders(ctx, &ec2.DescribeVerifiedAccessTrustProvidersInput{
		VerifiedAccessTrustProviderIds: []string{EC2VerifiedAccessTrustProviderId},
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range out.VerifiedAccessTrustProviders {
		resource := eC2VerifiedAccessTrustProviderHandle(ctx, v)
		values = append(values, resource)
	}

	return values, nil
}

func EC2DHCPOptions(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeDhcpOptionsPaginator(client, &ec2.DescribeDhcpOptionsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if !isErr(err, "InvalidDhcpOptionID.NotFound") {
				return nil, err
			}
			continue
		}

		for _, v := range page.DhcpOptions {
			resource := eC2DHCPOptionsHandle(ctx, v)
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
func eC2DHCPOptionsHandle(ctx context.Context, v types.DhcpOptions) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := fmt.Sprintf("arn:%s:ec2:%s:%s:dhcp-options/%s", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID, *v.DhcpOptionsId)

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *v.DhcpOptionsId,
		Description: model.EC2DhcpOptionsDescription{
			DhcpOptions: v,
		},
	}
	return resource
}
func GetEC2DHCPOptions(ctx context.Context, cfg aws.Config, field map[string]string) ([]Resource, error) {
	EC2DHCPOptionsId := field["id"]
	client := ec2.NewFromConfig(cfg)

	out, err := client.DescribeDhcpOptions(ctx, &ec2.DescribeDhcpOptionsInput{
		DhcpOptionsIds: []string{EC2DHCPOptionsId},
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range out.DhcpOptions {
		resource := eC2DHCPOptionsHandle(ctx, v)
		values = append(values, resource)
	}

	return values, nil
}

func EC2Fleet(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeFleetsPaginator(client, &ec2.DescribeFleetsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Fleets {
			resource := eC2FleetHandle(ctx, v)
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
func eC2FleetHandle(ctx context.Context, v types.FleetData) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := fmt.Sprintf("arn:%s:ec2:%s:%s:fleet/%s", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID, *v.FleetId)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ID:     arn,
		Name:   *v.FleetId,
		Description: model.EC2FleetDescription{
			Fleet: v,
		},
	}
	return resource
}
func GetEC2Fleet(ctx context.Context, cfg aws.Config, field map[string]string) ([]Resource, error) {
	EC2FleetId := field["id"]
	client := ec2.NewFromConfig(cfg)

	out, err := client.DescribeFleets(ctx, &ec2.DescribeFleetsInput{
		FleetIds: []string{EC2FleetId},
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range out.Fleets {
		resource := eC2FleetHandle(ctx, v)
		values = append(values, resource)
	}
	return values, nil
}

func EC2EgressOnlyInternetGateway(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeEgressOnlyInternetGatewaysPaginator(client, &ec2.DescribeEgressOnlyInternetGatewaysInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if !isErr(err, "InvalidEgressOnlyInternetGatewayId.NotFound") && !isErr(err, "InvalidEgressOnlyInternetGatewayId.Malformed") {
				return nil, err
			}
			continue
		}

		for _, v := range page.EgressOnlyInternetGateways {
			resource := eC2EgressOnlyInternetGatewayHandle(ctx, v)
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
func eC2EgressOnlyInternetGatewayHandle(ctx context.Context, v types.EgressOnlyInternetGateway) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := fmt.Sprintf("arn:%s:ec2:%s:%s:egress-only-internet-gateway/%s", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID, *v.EgressOnlyInternetGatewayId)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ID:     arn,
		Name:   *v.EgressOnlyInternetGatewayId,
		Description: model.EC2EgressOnlyInternetGatewayDescription{
			EgressOnlyInternetGateway: v,
		},
	}
	return resource
}
func GetEC2EgressOnlyInternetGateway(ctx context.Context, cfg aws.Config, field map[string]string) ([]Resource, error) {
	EC2EgressOnlyInternetGatewayId := field["id"]
	client := ec2.NewFromConfig(cfg)

	out, err := client.DescribeEgressOnlyInternetGateways(ctx, &ec2.DescribeEgressOnlyInternetGatewaysInput{
		EgressOnlyInternetGatewayIds: []string{EC2EgressOnlyInternetGatewayId},
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range out.EgressOnlyInternetGateways {
		resource := eC2EgressOnlyInternetGatewayHandle(ctx, v)
		values = append(values, resource)
	}
	return values, nil
}

func EC2EIP(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	output, err := client.DescribeAddresses(ctx, &ec2.DescribeAddressesInput{})
	if err != nil {
		if !isErr(err, "InvalidAllocationID.NotFound") && !isErr(err, "InvalidAllocationID.Malformed") {
			return nil, err
		}
		return nil, nil
	}

	var values []Resource
	for _, v := range output.Addresses {
		resource := eC2EIPHandle(ctx, v)
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
func eC2EIPHandle(ctx context.Context, v types.Address) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := "arn:" + describeCtx.Partition + ":ec2:" + describeCtx.Region + ":" + describeCtx.AccountID + ":eip/" + *v.AllocationId
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *v.AllocationId,
		Description: model.EC2EIPDescription{
			Address: v,
		},
	}
	return resource
}
func GetEC2EIP(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	allocationId := fields["id"]

	client := ec2.NewFromConfig(cfg)
	output, err := client.DescribeAddresses(ctx, &ec2.DescribeAddressesInput{
		AllocationIds: []string{allocationId},
	})
	if err != nil {
		if !isErr(err, "InvalidAllocationID.NotFound") && !isErr(err, "InvalidAllocationID.Malformed") {
			return nil, err
		}
		return nil, nil
	}
	var values []Resource
	for _, v := range output.Addresses {
		resource := eC2EIPHandle(ctx, v)
		values = append(values, resource)
	}

	return values, nil
}
func EC2EIPAddressTransfer(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	describeCtx := GetDescribeContext(ctx)
	paginator := ec2.NewDescribeAddressTransfersPaginator(client, &ec2.DescribeAddressTransfersInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, item := range page.AddressTransfers {
			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ID:     *item.TransferAccountId,
				Name:   *item.TransferAccountId,
				Description: model.EC2EIPAddressTransferDescription{
					AddressTransfer: item,
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

func EC2EnclaveCertificateIamRoleAssociation(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	certs, err := CertificateManagerCertificate(ctx, cfg, nil)
	if err != nil {
		return nil, err
	}

	client := ec2.NewFromConfig(cfg)

	var values []Resource
	for _, c := range certs {
		cert := c.Description.(model.CertificateManagerCertificateDescription)

		output, err := client.GetAssociatedEnclaveCertificateIamRoles(ctx, &ec2.GetAssociatedEnclaveCertificateIamRolesInput{
			CertificateArn: cert.Certificate.CertificateArn,
		})
		if err != nil {
			return nil, err
		}

		for _, v := range output.AssociatedRoles {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ID:          *v.AssociatedRoleArn, // Don't set to ARN since that will be the same for the role itself and this association
				Name:        *v.AssociatedRoleArn,
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
	}

	return values, nil
}

func EC2FlowLog(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeFlowLogsPaginator(client, &ec2.DescribeFlowLogsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.FlowLogs {
			resource := eC2FlowLogHandle(ctx, v)
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
func eC2FlowLogHandle(ctx context.Context, v types.FlowLog) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := "arn:" + describeCtx.Partition + ":ec2:" + describeCtx.Region + ":" + describeCtx.AccountID + ":vpc-flow-log/" + *v.FlowLogId
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *v.FlowLogId,
		Description: model.EC2FlowLogDescription{
			FlowLog: v,
		},
	}
	return resource
}
func GetEC2FlowLog(ctx context.Context, cfg aws.Config, field map[string]string) ([]Resource, error) {
	EC2FlowLogId := field["id"]
	client := ec2.NewFromConfig(cfg)

	output, err := client.DescribeFlowLogs(ctx, &ec2.DescribeFlowLogsInput{
		FlowLogIds: []string{EC2FlowLogId},
	})
	if err != nil {
		if !isErr(err, "InvalidAllocationID.NotFound") && !isErr(err, "InvalidAllocationID.Malformed") {
			return nil, err
		}
		return nil, nil
	}

	var values []Resource
	for _, v := range output.FlowLogs {
		resource := eC2FlowLogHandle(ctx, v)
		values = append(values, resource)
	}
	return values, nil
}

func EC2Host(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeHostsPaginator(client, &ec2.DescribeHostsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Hosts {
			resource := eC2HostHandle(ctx, v)
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
func eC2HostHandle(ctx context.Context, v types.Host) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := fmt.Sprintf("arn:%s:ec2:%s:%s:dedicated-host/%s", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID, *v.HostId)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ID:     arn,
		Name:   *v.HostId,
		Description: model.EC2HostDescription{
			Host: v,
		},
	}
	return resource
}
func GetEC2Host(ctx context.Context, cfg aws.Config, field map[string]string) ([]Resource, error) {
	EC2HostId := field["id"]
	client := ec2.NewFromConfig(cfg)

	output, err := client.DescribeHosts(ctx, &ec2.DescribeHostsInput{
		HostIds: []string{EC2HostId},
	})
	if err != nil {
		if !isErr(err, "InvalidAllocationID.NotFound") && !isErr(err, "InvalidAllocationID.Malformed") {
			return nil, err
		}
		return nil, nil
	}

	var values []Resource
	for _, v := range output.Hosts {
		resource := eC2HostHandle(ctx, v)
		values = append(values, resource)
	}
	return values, nil
}

func EC2Instance(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeInstancesPaginator(client, &ec2.DescribeInstancesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, r := range page.Reservations {
			for _, v := range r.Instances {
				resource, err := eC2InstanceHandle(ctx, v, client)
				if resource == nil {
					continue
				}
				if stream != nil {
					m := *stream
					err = m(*resource)
					if err != nil {
						return nil, err
					}
				} else {
					values = append(values, *resource)
				}
			}
		}
	}

	return values, nil
}
func eC2InstanceHandle(ctx context.Context, v types.Instance, client *ec2.Client) (*Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	var desc model.EC2InstanceDescription

	in := v // Do this to avoid the pointer being replaced by the for loop
	desc.Instance = &in

	statusOutput, err := client.DescribeInstanceStatus(ctx, &ec2.DescribeInstanceStatusInput{
		InstanceIds:         []string{*v.InstanceId},
		IncludeAllInstances: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}
	if len(statusOutput.InstanceStatuses) > 0 {
		desc.InstanceStatus = &statusOutput.InstanceStatuses[0]
	}

	if desc.InstanceStatus.InstanceState.Name == types.InstanceStateNameTerminated {
		return nil, nil
	}

	attrs := []types.InstanceAttributeName{
		types.InstanceAttributeNameUserData,
		types.InstanceAttributeNameInstanceInitiatedShutdownBehavior,
		types.InstanceAttributeNameDisableApiTermination,
		types.InstanceAttributeNameSriovNetSupport,
	}

	for _, attr := range attrs {
		output, err := client.DescribeInstanceAttribute(ctx, &ec2.DescribeInstanceAttributeInput{
			InstanceId: v.InstanceId,
			Attribute:  attr,
		})
		if err != nil {
			return nil, err
		}

		switch attr {
		case types.InstanceAttributeNameInstanceInitiatedShutdownBehavior:
			desc.Attributes.InstanceInitiatedShutdownBehavior = aws.ToString(output.InstanceInitiatedShutdownBehavior.Value)
		case types.InstanceAttributeNameDisableApiTermination:
			desc.Attributes.DisableApiTermination = aws.ToBool(output.DisableApiTermination.Value)
		case types.InstanceAttributeNameSriovNetSupport:
			desc.Instance.SriovNetSupport = output.SriovNetSupport.Value
		}
	}
	arn := "arn:" + describeCtx.Partition + ":ec2:" + describeCtx.Region + ":" + describeCtx.AccountID + ":instance/" + *v.InstanceId

	params := &ec2.GetLaunchTemplateDataInput{
		InstanceId: v.InstanceId,
	}

	op, err := client.GetLaunchTemplateData(ctx, params)
	if err != nil {
		return nil, err
	}
	desc.LaunchTemplateData = *op.LaunchTemplateData
	resource := Resource{
		Region:      describeCtx.KaytuRegion,
		ARN:         arn,
		Name:        *v.InstanceId,
		Description: desc,
	}
	return &resource, nil
}
func GetEC2Instance(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	instanceID := fields["id"]
	client := ec2.NewFromConfig(cfg)

	out, err := client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	})
	if err != nil {
		return nil, err
	}

	var values []Resource

	for _, r := range out.Reservations {
		for _, v := range r.Instances {
			resource, err := eC2InstanceHandle(ctx, v, client)
			if err != nil {
				return nil, err
			}
			if resource == nil {
				continue
			}
			values = append(values, *resource)
		}
	}

	return values, nil
}

func EC2InternetGateway(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeInternetGatewaysPaginator(client, &ec2.DescribeInternetGatewaysInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.InternetGateways {
			resource := eC2InternetGatewayHandle(ctx, v)
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
func eC2InternetGatewayHandle(ctx context.Context, v types.InternetGateway) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := "arn:" + describeCtx.Partition + ":ec2:" + describeCtx.Region + ":" + describeCtx.AccountID + ":internet-gateway/" + *v.InternetGatewayId
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *v.InternetGatewayId,
		Description: model.EC2InternetGatewayDescription{
			InternetGateway: v,
		},
	}
	return resource
}
func GetEC2InternetGateway(ctx context.Context, cfg aws.Config, field map[string]string) ([]Resource, error) {
	EC2InternetGatewayId := field["id"]
	client := ec2.NewFromConfig(cfg)

	output, err := client.DescribeInternetGateways(ctx, &ec2.DescribeInternetGatewaysInput{
		InternetGatewayIds: []string{EC2InternetGatewayId},
	})
	if err != nil {
		if !isErr(err, "InvalidAllocationID.NotFound") && !isErr(err, "InvalidAllocationID.Malformed") {
			return nil, err
		}
		return nil, nil
	}

	var values []Resource
	for _, v := range output.InternetGateways {
		resource := eC2InternetGatewayHandle(ctx, v)
		values = append(values, resource)
	}
	return values, nil
}

func EC2NatGateway(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeNatGatewaysPaginator(client, &ec2.DescribeNatGatewaysInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.NatGateways {
			resource := eC2NatGatewayHandle(ctx, v)
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
func eC2NatGatewayHandle(ctx context.Context, v types.NatGateway) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := "arn:" + describeCtx.Partition + ":ec2:" + describeCtx.Region + ":" + describeCtx.AccountID + ":natgateway/" + *v.NatGatewayId
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *v.NatGatewayId,
		Description: model.EC2NatGatewayDescription{
			NatGateway: v,
		},
	}
	return resource
}
func GetEC2NatGateway(ctx context.Context, cfg aws.Config, field map[string]string) ([]Resource, error) {
	EC2NatGatewayId := field["id"]
	client := ec2.NewFromConfig(cfg)

	output, err := client.DescribeNatGateways(ctx, &ec2.DescribeNatGatewaysInput{
		NatGatewayIds: []string{EC2NatGatewayId},
	})
	if err != nil {
		if !isErr(err, "InvalidAllocationID.NotFound") && !isErr(err, "InvalidAllocationID.Malformed") {
			return nil, err
		}
		return nil, nil
	}

	var values []Resource
	for _, v := range output.NatGateways {
		resource := eC2NatGatewayHandle(ctx, v)
		values = append(values, resource)
	}
	return values, nil
}

func EC2NetworkAcl(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeNetworkAclsPaginator(client, &ec2.DescribeNetworkAclsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.NetworkAcls {
			resource := eC2NetworkAclHandle(ctx, v)
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
func eC2NetworkAclHandle(ctx context.Context, v types.NetworkAcl) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := "arn:" + describeCtx.Partition + ":ec2:" + describeCtx.Region + ":" + describeCtx.AccountID + ":network-acl/" + *v.NetworkAclId
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *v.NetworkAclId,
		Description: model.EC2NetworkAclDescription{
			NetworkAcl: v,
		},
	}
	return resource
}
func GetEC2NetworkAcl(ctx context.Context, cfg aws.Config, field map[string]string) ([]Resource, error) {
	EC2NetworkAclId := field["id"]
	client := ec2.NewFromConfig(cfg)

	output, err := client.DescribeNetworkAcls(ctx, &ec2.DescribeNetworkAclsInput{
		NetworkAclIds: []string{EC2NetworkAclId},
	})
	if err != nil {
		if !isErr(err, "InvalidAllocationID.NotFound") && !isErr(err, "InvalidAllocationID.Malformed") {
			return nil, err
		}
		return nil, nil
	}

	var values []Resource
	for _, v := range output.NetworkAcls {
		resource := eC2NetworkAclHandle(ctx, v)
		values = append(values, resource)
	}
	return values, nil
}

func EC2NetworkInsightsAnalysis(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeNetworkInsightsAnalysesPaginator(client, &ec2.DescribeNetworkInsightsAnalysesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.NetworkInsightsAnalyses {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ARN:         *v.NetworkInsightsAnalysisArn,
				Name:        *v.NetworkInsightsAnalysisArn,
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
	}

	return values, nil
}

func EC2NetworkInsightsPath(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeNetworkInsightsPathsPaginator(client, &ec2.DescribeNetworkInsightsPathsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.NetworkInsightsPaths {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ARN:         *v.NetworkInsightsPathArn,
				Name:        *v.NetworkInsightsPathArn,
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
	}

	return values, nil
}

func EC2NetworkInterface(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	logger := GetLoggerFromContext(ctx)

	logger.Info("EC2NetworkInterface start working")

	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeNetworkInterfacesPaginator(client, &ec2.DescribeNetworkInterfacesInput{})

	logger.Info("EC2NetworkInterface start getting pages")
	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		logger.Info("EC2NetworkInterface got page")
		for _, v := range page.NetworkInterfaces {
			resource := eC2NetworkInterfaceHandle(ctx, v)
			if stream != nil {
				if err := (*stream)(resource); err != nil {
					return nil, err
				}
			} else {
				values = append(values, resource)
			}
		}
	}
	logger.Info("EC2NetworkInterface finished")

	return values, nil
}
func eC2NetworkInterfaceHandle(ctx context.Context, v types.NetworkInterface) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := "arn:" + describeCtx.Partition + ":ec2:" + describeCtx.Region + ":" + describeCtx.AccountID + ":network-interface/" + *v.NetworkInterfaceId
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *v.NetworkInterfaceId,
		Description: model.EC2NetworkInterfaceDescription{
			NetworkInterface: v,
		},
	}
	return resource
}
func GetEC2NetworkInterface(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	networkInterfaceID := fields["id"]
	client := ec2.NewFromConfig(cfg)
	out, err := client.DescribeNetworkInterfaces(ctx, &ec2.DescribeNetworkInterfacesInput{
		NetworkInterfaceIds: []string{networkInterfaceID},
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range out.NetworkInterfaces {
		resource := eC2NetworkInterfaceHandle(ctx, v)
		values = append(values, resource)
	}
	return values, nil
}

func EC2NetworkInterfacePermission(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeNetworkInterfacePermissionsPaginator(client, &ec2.DescribeNetworkInterfacePermissionsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.NetworkInterfacePermissions {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ID:          *v.NetworkInterfacePermissionId,
				Name:        *v.NetworkInterfacePermissionId,
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
	}
	return values, nil
}

func EC2PlacementGroup(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	output, err := client.DescribePlacementGroups(ctx, &ec2.DescribePlacementGroupsInput{})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range output.PlacementGroups {
		resource := eC2PlacementGroupHandle(ctx, v)
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
func eC2PlacementGroupHandle(ctx context.Context, v types.PlacementGroup) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := fmt.Sprintf("arn:%s:ec2:%s:%s:placement-group/%s", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID, *v.GroupName)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ID:     arn,
		Name:   *v.GroupName,
		Description: model.EC2PlacementGroupDescription{
			PlacementGroup: v,
		},
	}
	return resource
}
func GetEC2PlacementGroup(ctx context.Context, cfg aws.Config, field map[string]string) ([]Resource, error) {
	EC2PlacementGroupId := field["group_id"]
	client := ec2.NewFromConfig(cfg)

	output, err := client.DescribePlacementGroups(ctx, &ec2.DescribePlacementGroupsInput{
		GroupIds: []string{EC2PlacementGroupId},
	})
	if err != nil {
		if !isErr(err, "InvalidAllocationID.NotFound") && !isErr(err, "InvalidAllocationID.Malformed") {
			return nil, err
		}
		return nil, nil
	}

	var values []Resource
	for _, v := range output.PlacementGroups {
		resource := eC2PlacementGroupHandle(ctx, v)
		values = append(values, resource)
	}
	return values, nil
}

func EC2PrefixList(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribePrefixListsPaginator(client, &ec2.DescribePrefixListsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.PrefixLists {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ID:          *v.PrefixListId,
				Name:        *v.PrefixListName,
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
	}

	return values, nil
}

func EC2RegionalSettings(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ec2.NewFromConfig(cfg)
	out, err := client.GetEbsEncryptionByDefault(ctx, &ec2.GetEbsEncryptionByDefaultInput{})
	if err != nil {
		return nil, err
	}
	outkey, err := client.GetEbsDefaultKmsKeyId(ctx, &ec2.GetEbsDefaultKmsKeyIdInput{})
	if err != nil {
		return nil, err
	}
	outstate, err := client.GetSnapshotBlockPublicAccessState(ctx, &ec2.GetSnapshotBlockPublicAccessStateInput{})
	if err != nil {
		return nil, err
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		// No ID or ARN. Per Account Configuration
		Name: cfg.Region + " EC2 Settings", // Based on Steampipe
		Description: model.EC2RegionalSettingsDescription{
			EbsEncryptionByDefault:         out.EbsEncryptionByDefault,
			KmsKeyId:                       outkey.KmsKeyId,
			SnapshotBlockPublicAccessState: outstate.State,
		},
	}
	var values []Resource
	if stream != nil {
		if err := (*stream)(resource); err != nil {
			return nil, err
		}
	} else {
		values = append(values, resource)
	}

	return values, nil
}

func EC2RouteTable(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeRouteTablesPaginator(client, &ec2.DescribeRouteTablesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.RouteTables {
			resource := eC2RouteTableHandle(ctx, v)
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
func eC2RouteTableHandle(ctx context.Context, v types.RouteTable) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := "arn:" + describeCtx.Partition + ":ec2:" + describeCtx.Region + ":" + describeCtx.AccountID + ":route-table/" + *v.RouteTableId

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *v.RouteTableId,
		Description: model.EC2RouteTableDescription{
			RouteTable: v,
		},
	}
	return resource
}
func GetEC2RouteTable(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)

	routeTableID := fields["id"]

	out, err := client.DescribeRouteTables(ctx, &ec2.DescribeRouteTablesInput{RouteTableIds: []string{routeTableID}})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range out.RouteTables {
		resource := eC2RouteTableHandle(ctx, v)
		values = append(values, resource)
	}

	return values, nil
}

func EC2LocalGatewayRouteTable(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeLocalGatewayRouteTablesPaginator(client, &ec2.DescribeLocalGatewayRouteTablesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.LocalGatewayRouteTables {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ARN:         *v.LocalGatewayRouteTableArn,
				Name:        *v.LocalGatewayRouteTableId,
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
	}

	return values, nil
}

func EC2LocalGatewayRouteTableVPCAssociation(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeLocalGatewayRouteTableVpcAssociationsPaginator(client, &ec2.DescribeLocalGatewayRouteTableVpcAssociationsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.LocalGatewayRouteTableVpcAssociations {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ID:          *v.LocalGatewayRouteTableVpcAssociationId,
				Name:        *v.LocalGatewayRouteTableVpcAssociationId,
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
	}

	return values, nil
}

func EC2TransitGatewayRouteTable(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeTransitGatewayRouteTablesPaginator(client, &ec2.DescribeTransitGatewayRouteTablesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if !isErr(err, "InvalidRouteTableID.NotFound") && !isErr(err, "InvalidRouteTableId.Unavailable") && !isErr(err, "InvalidRouteTableId.Malformed") {
				return nil, err
			}
			continue
		}

		for _, v := range page.TransitGatewayRouteTables {
			resource := eC2TransitGatewayRouteTableHandle(ctx, v)
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
func eC2TransitGatewayRouteTableHandle(ctx context.Context, v types.TransitGatewayRouteTable) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := fmt.Sprintf("arn:%s:ec2:%s:%s:transit-gateway-route-table/%s", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID, *v.TransitGatewayRouteTableId)

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *v.TransitGatewayRouteTableId,
		Description: model.EC2TransitGatewayRouteTableDescription{
			TransitGatewayRouteTable: v,
		},
	}
	return resource
}
func GetEC2TransitGatewayRouteTable(ctx context.Context, cfg aws.Config, field map[string]string) ([]Resource, error) {
	EC2TransitGatewayRouteTableId := field["id"]
	client := ec2.NewFromConfig(cfg)

	out, err := client.DescribeTransitGatewayRouteTables(ctx, &ec2.DescribeTransitGatewayRouteTablesInput{
		TransitGatewayRouteTableIds: []string{EC2TransitGatewayRouteTableId},
	})
	if err != nil {
		if !isErr(err, "InvalidRouteTableID.NotFound") && !isErr(err, "InvalidRouteTableId.Unavailable") && !isErr(err, "InvalidRouteTableId.Malformed") {
			return nil, err
		}
		return nil, nil
	}

	var values []Resource
	for _, v := range out.TransitGatewayRouteTables {
		resource := eC2TransitGatewayRouteTableHandle(ctx, v)
		values = append(values, resource)
	}
	return values, nil
}
func EC2TransitGatewayRouteTableAssociation(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	rts, err := EC2TransitGatewayRouteTable(ctx, cfg, nil)
	if err != nil {
		return nil, err
	}

	client := ec2.NewFromConfig(cfg)

	var values []Resource
	for _, r := range rts {
		routeTable := r.Description.(types.TransitGatewayRouteTable)
		paginator := ec2.NewGetTransitGatewayRouteTableAssociationsPaginator(client, &ec2.GetTransitGatewayRouteTableAssociationsInput{
			TransitGatewayRouteTableId: routeTable.TransitGatewayRouteTableId,
		})
		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				return nil, err
			}

			for _, v := range page.Associations {
				resource := Resource{
					Region:      describeCtx.KaytuRegion,
					ID:          *v.TransitGatewayAttachmentId,
					Name:        *v.TransitGatewayAttachmentId,
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
		}
	}

	return values, nil
}

func EC2TransitGatewayRouteTablePropagation(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	rts, err := EC2TransitGatewayRouteTable(ctx, cfg, nil)
	if err != nil {
		return nil, err
	}

	client := ec2.NewFromConfig(cfg)

	var values []Resource
	for _, r := range rts {
		routeTable := r.Description.(types.TransitGatewayRouteTable)
		paginator := ec2.NewGetTransitGatewayRouteTablePropagationsPaginator(client, &ec2.GetTransitGatewayRouteTablePropagationsInput{
			TransitGatewayRouteTableId: routeTable.TransitGatewayRouteTableId,
		})

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				return nil, err
			}

			for _, v := range page.TransitGatewayRouteTablePropagations {
				resource := Resource{
					Region:      describeCtx.KaytuRegion,
					ID:          CompositeID(*routeTable.TransitGatewayRouteTableId, *v.TransitGatewayAttachmentId),
					Name:        *routeTable.TransitGatewayRouteTableId,
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
		}
	}

	return values, nil
}

func EC2SecurityGroup(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeSecurityGroupsPaginator(client, &ec2.DescribeSecurityGroupsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.SecurityGroups {
			resource := eC2SecurityGroupHandle(ctx, v)
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
func eC2SecurityGroupHandle(ctx context.Context, v types.SecurityGroup) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := fmt.Sprintf("arn:%s:ec2:%s:%s:security-group/%s", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID, *v.GroupId)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *v.GroupName,
		Description: model.EC2SecurityGroupDescription{
			SecurityGroup: v,
		},
	}
	return resource
}
func GetEC2SecurityGroup(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	groupID := fields["group_id"]
	client := ec2.NewFromConfig(cfg)

	out, err := client.DescribeSecurityGroups(ctx, &ec2.DescribeSecurityGroupsInput{
		GroupIds: []string{groupID},
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range out.SecurityGroups {
		resource := eC2SecurityGroupHandle(ctx, v)
		values = append(values, resource)
	}

	return values, nil
}

func getEC2SecurityGroupRuleDescriptionFromIPPermission(group types.SecurityGroup, permission types.IpPermission, groupType string) []model.EC2SecurityGroupRuleDescription {
	var descArr []model.EC2SecurityGroupRuleDescription

	// create 1 row per ip-range
	if permission.IpRanges != nil {
		for _, r := range permission.IpRanges {
			descArr = append(descArr, model.EC2SecurityGroupRuleDescription{
				Group:           group,
				Permission:      permission,
				IPRange:         &r,
				Ipv6Range:       nil,
				UserIDGroupPair: nil,
				PrefixListId:    nil,
				Type:            groupType,
			})
		}
	}

	// create 1 row per prefix-list Id
	if permission.PrefixListIds != nil {
		for _, r := range permission.PrefixListIds {
			descArr = append(descArr, model.EC2SecurityGroupRuleDescription{
				Group:           group,
				Permission:      permission,
				IPRange:         nil,
				Ipv6Range:       nil,
				UserIDGroupPair: nil,
				PrefixListId:    &r,
				Type:            groupType,
			})
		}
	}

	// create 1 row per ipv6-range
	if permission.Ipv6Ranges != nil {
		for _, r := range permission.Ipv6Ranges {
			descArr = append(descArr, model.EC2SecurityGroupRuleDescription{
				Group:           group,
				Permission:      permission,
				IPRange:         nil,
				Ipv6Range:       &r,
				UserIDGroupPair: nil,
				PrefixListId:    nil,
				Type:            groupType,
			})
		}
	}

	// create 1 row per user id group pair
	if permission.UserIdGroupPairs != nil {
		for _, r := range permission.UserIdGroupPairs {
			descArr = append(descArr, model.EC2SecurityGroupRuleDescription{
				Group:           group,
				Permission:      permission,
				IPRange:         nil,
				Ipv6Range:       nil,
				UserIDGroupPair: &r,
				PrefixListId:    nil,
				Type:            groupType,
			})
		}
	}

	return descArr
}

func EC2SecurityGroupRule(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	groups, err := EC2SecurityGroup(ctx, cfg, nil)
	if err != nil {
		return nil, err
	}

	var values []Resource
	descArr := make([]model.EC2SecurityGroupRuleDescription, 0, 128)
	for _, groupWrapper := range groups {
		group := groupWrapper.Description.(model.EC2SecurityGroupDescription).SecurityGroup
		if group.IpPermissions != nil {
			for _, permission := range group.IpPermissions {
				descArr = append(descArr, getEC2SecurityGroupRuleDescriptionFromIPPermission(group, permission, "ingress")...)
			}
		}
		if group.IpPermissionsEgress != nil {
			for _, permission := range group.IpPermissionsEgress {
				descArr = append(descArr, getEC2SecurityGroupRuleDescriptionFromIPPermission(group, permission, "egress")...)
			}
		}
	}
	for _, desc := range descArr {
		resource := eC2SecurityGroupRuleHandle(ctx, desc)
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
func eC2SecurityGroupRuleHandle(ctx context.Context, desc model.EC2SecurityGroupRuleDescription) Resource {
	describeCtx := GetDescribeContext(ctx)

	hashCode := desc.Type + "_" + *desc.Permission.IpProtocol
	if desc.Permission.FromPort != nil {
		hashCode = hashCode + "_" + fmt.Sprint(desc.Permission.FromPort) + "_" + fmt.Sprint(desc.Permission.ToPort)
	}

	if desc.IPRange != nil && desc.IPRange.CidrIp != nil {
		hashCode = hashCode + "_" + *desc.IPRange.CidrIp
	} else if desc.Ipv6Range != nil && desc.Ipv6Range.CidrIpv6 != nil {
		hashCode = hashCode + "_" + *desc.Ipv6Range.CidrIpv6
	} else if desc.UserIDGroupPair != nil && *desc.UserIDGroupPair.GroupId == *desc.Group.GroupId {
		hashCode = hashCode + "_" + *desc.Group.GroupId
	} else if desc.PrefixListId != nil && desc.PrefixListId.PrefixListId != nil {
		hashCode = hashCode + "_" + *desc.PrefixListId.PrefixListId
	}

	arn := fmt.Sprintf("arn:%s:ec2:%s:%s:security-group/%s:%s", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID, *desc.Group.GroupId, hashCode)
	resource := Resource{
		Region:      describeCtx.KaytuRegion,
		ARN:         arn,
		Name:        fmt.Sprintf("%s_%s", *desc.Group.GroupId, hashCode),
		Description: desc,
	}
	return resource
}
func GetEC2SecurityGroupRule(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	groupName := fields["name"]

	out, err := client.DescribeSecurityGroups(ctx, &ec2.DescribeSecurityGroupsInput{
		GroupNames: []string{groupName},
	})
	if err != nil {
		if isErr(err, "DescribeConfigurationSetNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	descArr := make([]model.EC2SecurityGroupRuleDescription, 0, 128)
	for _, group := range out.SecurityGroups {
		if group.IpPermissions != nil {
			for _, permission := range group.IpPermissions {
				descArr = append(descArr, getEC2SecurityGroupRuleDescriptionFromIPPermission(group, permission, "ingress")...)
			}
		}
		if group.IpPermissionsEgress != nil {
			for _, permission := range group.IpPermissionsEgress {
				descArr = append(descArr, getEC2SecurityGroupRuleDescriptionFromIPPermission(group, permission, "egress")...)
			}
		}
	}

	var values []Resource
	for _, desc := range descArr {
		resource := eC2SecurityGroupRuleHandle(ctx, desc)
		values = append(values, resource)
	}
	return values, nil
}

func EC2SpotFleet(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeSpotFleetRequestsPaginator(client, &ec2.DescribeSpotFleetRequestsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.SpotFleetRequestConfigs {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ID:          *v.SpotFleetRequestId,
				Name:        *v.SpotFleetRequestId,
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
	}

	return values, nil
}

func EC2Subnet(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeSubnetsPaginator(client, &ec2.DescribeSubnetsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Subnets {
			resource := eC2SubnetHandle(ctx, v)
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
func eC2SubnetHandle(ctx context.Context, v types.Subnet) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *v.SubnetArn,
		Name:   *v.SubnetId,
		Description: model.EC2SubnetDescription{
			Subnet: v,
		},
	}
	return resource
}
func GetEC2Subnet(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	subnetId := fields["id"]
	client := ec2.NewFromConfig(cfg)
	out, err := client.DescribeSubnets(ctx, &ec2.DescribeSubnetsInput{
		SubnetIds: []string{subnetId},
	})
	if err != nil {
		return nil, err
	}

	var values []Resource

	for _, v := range out.Subnets {
		resource := eC2SubnetHandle(ctx, v)
		values = append(values, resource)
	}

	return values, nil
}

func EC2TrafficMirrorFilter(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeTrafficMirrorFiltersPaginator(client, &ec2.DescribeTrafficMirrorFiltersInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.TrafficMirrorFilters {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ID:          *v.TrafficMirrorFilterId,
				Name:        *v.TrafficMirrorFilterId,
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
	}

	return values, nil
}

func EC2TrafficMirrorSession(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeTrafficMirrorSessionsPaginator(client, &ec2.DescribeTrafficMirrorSessionsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.TrafficMirrorSessions {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ID:          *v.TrafficMirrorSessionId,
				Name:        *v.TrafficMirrorFilterId,
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
	}

	return values, nil
}

func EC2TrafficMirrorTarget(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeTrafficMirrorTargetsPaginator(client, &ec2.DescribeTrafficMirrorTargetsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.TrafficMirrorTargets {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ID:          *v.TrafficMirrorTargetId,
				Name:        *v.TrafficMirrorTargetId,
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
	}

	return values, nil
}

func EC2TransitGateway(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeTransitGatewaysPaginator(client, &ec2.DescribeTransitGatewaysInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if !isErr(err, "InvalidTransitGatewayID.NotFound") && !isErr(err, "InvalidTransitGatewayID.Unavailable") && !isErr(err, "InvalidTransitGatewayID.Malformed") {
				return nil, err
			}
			continue
		}

		for _, v := range page.TransitGateways {
			resource := eC2TransitGatewayHandle(ctx, v)
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
func eC2TransitGatewayHandle(ctx context.Context, v types.TransitGateway) Resource {
	describeCtx := GetDescribeContext(ctx)
	values := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *v.TransitGatewayArn,
		Name:   *v.TransitGatewayId,
		Description: model.EC2TransitGatewayDescription{
			TransitGateway: v,
		},
	}
	return values
}
func GetEC2TransitGateway(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	ransitGatewayId := fields["id"]
	client := ec2.NewFromConfig(cfg)

	out, err := client.DescribeTransitGateways(ctx, &ec2.DescribeTransitGatewaysInput{
		TransitGatewayIds: []string{ransitGatewayId},
	})
	if err != nil {
		if !isErr(err, "InvalidTransitGatewayID.NotFound") && !isErr(err, "InvalidTransitGatewayID.Unavailable") && !isErr(err, "InvalidTransitGatewayID.Malformed") {
			return nil, err
		}
		return nil, nil
	}

	var values []Resource
	for _, v := range out.TransitGateways {
		resource := eC2TransitGatewayHandle(ctx, v)
		values = append(values, resource)
	}
	return values, nil
}
func EC2TransitGatewayConnect(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeTransitGatewayConnectsPaginator(client, &ec2.DescribeTransitGatewayConnectsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.TransitGatewayConnects {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ID:          *v.TransitGatewayAttachmentId,
				Name:        *v.TransitGatewayAttachmentId,
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
	}

	return values, nil
}

func EC2TransitGatewayMulticastDomain(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeTransitGatewayMulticastDomainsPaginator(client, &ec2.DescribeTransitGatewayMulticastDomainsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.TransitGatewayMulticastDomains {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ARN:         *v.TransitGatewayMulticastDomainArn,
				Name:        *v.TransitGatewayMulticastDomainArn,
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
	}

	return values, nil
}

func EC2TransitGatewayMulticastDomainAssociation(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	domains, err := EC2TransitGatewayMulticastDomain(ctx, cfg, nil)
	if err != nil {
		return nil, err
	}

	client := ec2.NewFromConfig(cfg)
	var values []Resource
	for _, domain := range domains {
		paginator := ec2.NewGetTransitGatewayMulticastDomainAssociationsPaginator(client, &ec2.GetTransitGatewayMulticastDomainAssociationsInput{
			TransitGatewayMulticastDomainId: domain.Description.(types.TransitGatewayMulticastDomain).TransitGatewayMulticastDomainId,
		})

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				return nil, err
			}

			for _, v := range page.MulticastDomainAssociations {
				resource := Resource{
					Region:      describeCtx.KaytuRegion,
					ID:          *v.TransitGatewayAttachmentId,
					Name:        *v.TransitGatewayAttachmentId,
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
		}
	}

	return values, nil
}

func EC2TransitGatewayMulticastGroupMember(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	domains, err := EC2TransitGatewayMulticastDomain(ctx, cfg, nil)
	if err != nil {
		return nil, err
	}

	client := ec2.NewFromConfig(cfg)
	var values []Resource
	for _, domain := range domains {
		tgmdID := domain.Description.(types.TransitGatewayMulticastDomain).TransitGatewayMulticastDomainId
		paginator := ec2.NewSearchTransitGatewayMulticastGroupsPaginator(client, &ec2.SearchTransitGatewayMulticastGroupsInput{
			TransitGatewayMulticastDomainId: tgmdID,
			Filters: []types.Filter{
				{
					Name:   aws.String("is-group-member"),
					Values: []string{"true"},
				},
			},
		})

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				return nil, err
			}

			for _, v := range page.MulticastGroups {
				resource := Resource{
					Region:      describeCtx.KaytuRegion,
					ID:          CompositeID(*tgmdID, *v.GroupIpAddress),
					Name:        *v.GroupIpAddress,
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
		}
	}

	return values, nil
}

func EC2TransitGatewayMulticastGroupSource(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	domains, err := EC2TransitGatewayMulticastDomain(ctx, cfg, nil)
	if err != nil {
		return nil, err
	}

	client := ec2.NewFromConfig(cfg)
	var values []Resource
	for _, domain := range domains {
		tgmdID := domain.Description.(types.TransitGatewayMulticastDomain).TransitGatewayMulticastDomainId
		paginator := ec2.NewSearchTransitGatewayMulticastGroupsPaginator(client, &ec2.SearchTransitGatewayMulticastGroupsInput{
			TransitGatewayMulticastDomainId: tgmdID,
			Filters: []types.Filter{
				{
					Name:   aws.String("is-group-source"),
					Values: []string{"true"},
				},
			},
		})

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				return nil, err
			}

			for _, v := range page.MulticastGroups {
				resource := Resource{
					Region:      describeCtx.KaytuRegion,
					ID:          CompositeID(*tgmdID, *v.GroupIpAddress),
					Name:        *v.GroupIpAddress,
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
		}
	}

	return values, nil
}

func EC2TransitGatewayPeeringAttachment(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeTransitGatewayPeeringAttachmentsPaginator(client, &ec2.DescribeTransitGatewayPeeringAttachmentsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.TransitGatewayPeeringAttachments {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ID:          *v.TransitGatewayAttachmentId,
				Name:        *v.TransitGatewayAttachmentId,
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
	}

	return values, nil
}

func EC2VPC(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeVpcsPaginator(client, &ec2.DescribeVpcsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Vpcs {
			resource := eC2VPCHandle(ctx, v)
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
func eC2VPCHandle(ctx context.Context, v types.Vpc) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := "arn:" + describeCtx.Partition + ":ec2:" + describeCtx.Region + ":" + describeCtx.AccountID + ":vpc/" + *v.VpcId
	values := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *v.VpcId,
		Description: model.EC2VpcDescription{
			Vpc: v,
		},
	}
	return values
}
func GetEC2VPC(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)

	vpcID := fields["id"]

	out, err := client.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{
		VpcIds: []string{vpcID},
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range out.Vpcs {
		resource := eC2VPCHandle(ctx, v)
		values = append(values, resource)
	}
	return values, nil
}

func EC2VPCEndpoint(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeVpcEndpointsPaginator(client, &ec2.DescribeVpcEndpointsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.VpcEndpoints {
			resource := eC2VPCEndpointHandle(ctx, v)

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
func eC2VPCEndpointHandle(ctx context.Context, v types.VpcEndpoint) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := "arn:" + describeCtx.Partition + ":ec2:" + describeCtx.Region + ":" + describeCtx.AccountID + ":vpc-endpoint/" + *v.VpcEndpointId
	values := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		ID:     *v.VpcEndpointId,
		Name:   *v.VpcEndpointId,
		Description: model.EC2VPCEndpointDescription{
			VpcEndpoint: v,
		},
	}
	return values
}
func GetEC2VPCEndpoint(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	VPCEndpointId := fields["id"]
	client := ec2.NewFromConfig(cfg)

	out, err := client.DescribeVpcEndpoints(ctx, &ec2.DescribeVpcEndpointsInput{
		VpcEndpointIds: []string{VPCEndpointId},
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range out.VpcEndpoints {
		resource := eC2VPCEndpointHandle(ctx, v)
		values = append(values, resource)
	}
	return values, nil
}

func EC2VPCEndpointConnectionNotification(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeVpcEndpointConnectionNotificationsPaginator(client, &ec2.DescribeVpcEndpointConnectionNotificationsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.ConnectionNotificationSet {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ARN:         *v.ConnectionNotificationArn,
				Name:        *v.ConnectionNotificationArn,
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
	}

	return values, nil
}

func EC2VPCEndpointService(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	fmt.Println("EC2VPCEndpointService")
	describeCtx := GetDescribeContext(ctx)

	client := ec2.NewFromConfig(cfg)
	var values []Resource

	fmt.Println("EC2VPCEndpointService DescribeVpcEndpointServices")
	output, err := client.DescribeVpcEndpointServices(ctx, &ec2.DescribeVpcEndpointServicesInput{})
	fmt.Println("EC2VPCEndpointService DescribeVpcEndpointServices done")
	if err != nil {
		return nil, err
	}
	if output == nil {
		return nil, nil
	}

	for _, v := range output.ServiceDetails {
		arn := ""
		if v.ServiceName != nil {
			splitServiceName := strings.Split(*v.ServiceName, ".")
			arn = fmt.Sprintf("arn:%s:ec2:%s:%s:vpc-endpoint-service/%s", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID, splitServiceName[len(splitServiceName)-1])
		} else if v.ServiceId != nil {
			arn = fmt.Sprintf("arn:%s:ec2:%s:%s:vpc-endpoint-service/%s", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID, *v.ServiceId)
		} else {
			continue
		}

		fmt.Println("EC2VPCEndpointService NewDescribeVpcEndpointServicePermissionsPaginator")
		paginator := ec2.NewDescribeVpcEndpointServicePermissionsPaginator(client, &ec2.DescribeVpcEndpointServicePermissionsInput{
			ServiceId: v.ServiceId,
		}, func(o *ec2.DescribeVpcEndpointServicePermissionsPaginatorOptions) {
			o.Limit = 100
			o.StopOnDuplicateToken = true
		})
		fmt.Println("EC2VPCEndpointService NewDescribeVpcEndpointServicePermissionsPaginator done")

		var allowedPrincipals []types.AllowedPrincipal
		for paginator.HasMorePages() {
			fmt.Println("EC2VPCEndpointService NewDescribeVpcEndpointServicePermissionsPaginator next page")
			permissions, err := paginator.NextPage(ctx)
			fmt.Println("EC2VPCEndpointService NewDescribeVpcEndpointServicePermissionsPaginator got the page", err)
			if err != nil {
				if err != nil {
					var ae smithy.APIError
					if errors.As(err, &ae) && ae.ErrorCode() == "InvalidVpcEndpointServiceId.NotFound" {
						// VpcEndpoint doesn't have permissions set. Move on!
						break
					}
					return nil, err
				}
			}
			if permissions == nil {
				break
			}
			allowedPrincipals = append(allowedPrincipals, permissions.AllowedPrincipals...)
		}
		fmt.Println("EC2VPCEndpointService NewDescribeVpcEndpointServicePermissionsPaginator got all pages")
		var vpcEndpointConnections []types.VpcEndpointConnection
		if v.ServiceId != nil {
			op, err := client.DescribeVpcEndpointConnections(ctx, &ec2.DescribeVpcEndpointConnectionsInput{
				Filters: []types.Filter{
					{
						Name:   aws.String("service-id"),
						Values: []string{*v.ServiceId},
					},
				},
			})
			if err != nil {
				return nil, err
			}
			if op != nil && op.VpcEndpointConnections != nil && len(op.VpcEndpointConnections) > 0 {
				vpcEndpointConnections = op.VpcEndpointConnections
			}
		}
		fmt.Println("EC2VPCEndpointService DescribeVpcEndpointConnections done")

		resource := Resource{
			Region: describeCtx.KaytuRegion,
			ARN:    arn,
			Description: model.EC2VPCEndpointServiceDescription{
				VpcEndpointService:     v,
				AllowedPrincipals:      allowedPrincipals,
				VpcEndpointConnections: vpcEndpointConnections,
			},
		}
		if v.ServiceName != nil {
			resource.Name = *v.ServiceName
		}

		if stream != nil {
			fmt.Println("EC2VPCEndpointService sending to stream")
			if err := (*stream)(resource); err != nil {
				return nil, err
			}
			fmt.Println("EC2VPCEndpointService sending to stream done")
		} else {
			values = append(values, resource)
		}
	}

	if err != nil {
		return nil, err
	}

	fmt.Println("EC2VPCEndpointService finish")

	return values, nil
}

func EC2VPCEndpointServicePermissions(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	services, err := EC2VPCEndpointService(ctx, cfg, nil)
	if err != nil {
		return nil, err
	}

	client := ec2.NewFromConfig(cfg)

	var values []Resource
	for _, s := range services {
		service := s.Description.(model.EC2VPCEndpointServiceDescription).VpcEndpointService

		paginator := ec2.NewDescribeVpcEndpointServicePermissionsPaginator(client, &ec2.DescribeVpcEndpointServicePermissionsInput{
			ServiceId: service.ServiceId,
		})
		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				var ae smithy.APIError
				if errors.As(err, &ae) && ae.ErrorCode() == "InvalidVpcEndpointServiceId.NotFound" {
					// VpcEndpoint doesn't have permissions set. Move on!
					break
				}
				return nil, err
			}

			for _, v := range page.AllowedPrincipals {
				resource := Resource{
					Region:      describeCtx.KaytuRegion,
					ARN:         *v.Principal,
					Name:        *v.Principal,
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
		}
	}

	return values, nil
}

func EC2VPCPeeringConnection(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeVpcPeeringConnectionsPaginator(client, &ec2.DescribeVpcPeeringConnectionsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.VpcPeeringConnections {
			resource := eC2VPCPeeringConnectionHandle(ctx, v)
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
func eC2VPCPeeringConnectionHandle(ctx context.Context, v types.VpcPeeringConnection) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := fmt.Sprintf("arn:%s:ec2:%s:%s:vpc-peering-connection/%s", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID, *v.VpcPeeringConnectionId)
	values := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *v.VpcPeeringConnectionId,
		Description: model.EC2VpcPeeringConnectionDescription{
			VpcPeeringConnection: v,
		},
	}
	return values
}
func GetEC2VPCPeeringConnection(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	VPCPeeringConnectionId := fields["id"]
	client := ec2.NewFromConfig(cfg)

	out, err := client.DescribeVpcPeeringConnections(ctx, &ec2.DescribeVpcPeeringConnectionsInput{
		VpcPeeringConnectionIds: []string{VPCPeeringConnectionId},
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range out.VpcPeeringConnections {
		resource := eC2VPCPeeringConnectionHandle(ctx, v)
		values = append(values, resource)
	}
	return values, nil
}

func EC2VPNConnection(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	output, err := client.DescribeVpnConnections(ctx, &ec2.DescribeVpnConnectionsInput{})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range output.VpnConnections {
		resource := eC2VPNConnectionHandle(ctx, v)
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
func eC2VPNConnectionHandle(ctx context.Context, v types.VpnConnection) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := "arn:" + describeCtx.Partition + ":ec2:" + describeCtx.Region + ":" + describeCtx.AccountID + ":vpn-connection/" + *v.VpnConnectionId
	values := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *v.VpnConnectionId,
		Description: model.EC2VPNConnectionDescription{
			VpnConnection: v,
		},
	}
	return values
}
func GetEC2VPNConnection(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	VPNConnectionId := fields["id"]
	client := ec2.NewFromConfig(cfg)

	out, err := client.DescribeVpnConnections(ctx, &ec2.DescribeVpnConnectionsInput{
		VpnConnectionIds: []string{VPNConnectionId},
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range out.VpnConnections {
		resource := eC2VPNConnectionHandle(ctx, v)
		values = append(values, resource)
	}
	return values, nil
}

func EC2VPNGateway(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	output, err := client.DescribeVpnGateways(ctx, &ec2.DescribeVpnGatewaysInput{})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range output.VpnGateways {
		resource := eC2VPNGatewayHandle(ctx, v)

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
func eC2VPNGatewayHandle(ctx context.Context, v types.VpnGateway) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := fmt.Sprintf("arn:%s:ec2:%s:%s:vpn-gateway/%s", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID, *v.VpnGatewayId)
	values := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		ID:     *v.VpnGatewayId,
		Name:   *v.VpnGatewayId,
		Description: model.EC2VPNGatewayDescription{
			VPNGateway: v,
		},
	}
	return values
}
func GetEC2VPNGateway(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	VPNGatewayId := fields["id"]
	client := ec2.NewFromConfig(cfg)

	out, err := client.DescribeVpnGateways(ctx, &ec2.DescribeVpnGatewaysInput{
		VpnGatewayIds: []string{VPNGatewayId},
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range out.VpnGateways {
		resource := eC2VPNGatewayHandle(ctx, v)
		values = append(values, resource)
	}
	return values, nil
}

func EC2Region(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ec2.NewFromConfig(cfg)
	output, err := client.DescribeRegions(ctx, &ec2.DescribeRegionsInput{
		AllRegions: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range output.Regions {
		arn := "arn:" + describeCtx.Partition + "::" + *v.RegionName + ":" + describeCtx.AccountID
		resource := Resource{
			Region: describeCtx.KaytuRegion,
			ARN:    arn,
			Name:   *v.RegionName,
			Description: model.EC2RegionDescription{
				Region: v,
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

	return values, nil
}

func EC2AvailabilityZone(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)

	regionsOutput, err := client.DescribeRegions(ctx, &ec2.DescribeRegionsInput{
		AllRegions: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, region := range regionsOutput.Regions {
		if region.OptInStatus != nil && *region.OptInStatus != "not-opted-in" {
			continue
		}
		output, err := client.DescribeAvailabilityZones(ctx, &ec2.DescribeAvailabilityZonesInput{
			AllAvailabilityZones: aws.Bool(true),
			Filters: []types.Filter{
				{
					Name:   aws.String("region-name"),
					Values: []string{*region.RegionName},
				},
			},
		})
		if err != nil {
			return nil, err
		}

		for _, v := range output.AvailabilityZones {
			resource := eC2AvailabilityZoneHandle(ctx, v, region)
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
func eC2AvailabilityZoneHandle(ctx context.Context, v types.AvailabilityZone, region types.Region) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := fmt.Sprintf("arn:%s::%s::availability-zone/%s", describeCtx.Partition, *region.RegionName, *v.ZoneName)
	values := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *v.RegionName,
		Description: model.EC2AvailabilityZoneDescription{
			AvailabilityZone: v,
		},
	}
	return values
}
func GetEC2AvailabilityZone(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	ZoneIds := fields["id"]
	client := ec2.NewFromConfig(cfg)

	regionsOutput, err := client.DescribeRegions(ctx, &ec2.DescribeRegionsInput{
		AllRegions: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, region := range regionsOutput.Regions {
		if region.OptInStatus != nil && *region.OptInStatus != "not-opted-in" {
			continue
		}
		output, err := client.DescribeAvailabilityZones(ctx, &ec2.DescribeAvailabilityZonesInput{
			AllAvailabilityZones: aws.Bool(true),
			Filters: []types.Filter{
				{
					Name:   aws.String("region-name"),
					Values: []string{*region.RegionName},
				},
			},
			ZoneIds: []string{ZoneIds},
		})
		if err != nil {
			return nil, err
		}
		for _, v := range output.AvailabilityZones {
			resource := eC2AvailabilityZoneHandle(ctx, v, region)
			values = append(values, resource)
		}
	}
	return values, nil
}

func EC2KeyPair(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	output, err := client.DescribeKeyPairs(ctx, &ec2.DescribeKeyPairsInput{})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range output.KeyPairs {
		resource := eC2KeyPairHandle(ctx, v)
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
func eC2KeyPairHandle(ctx context.Context, v types.KeyPairInfo) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := "arn:" + describeCtx.Partition + ":ec2:" + describeCtx.Region + ":" + describeCtx.AccountID + ":key-pair/" + *v.KeyName
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *v.KeyName,
		Description: model.EC2KeyPairDescription{
			KeyPair: v,
		},
	}
	return resource
}
func GetEC2KeyPair(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	keyPairID := fields["id"]

	client := ec2.NewFromConfig(cfg)
	output, err := client.DescribeKeyPairs(ctx, &ec2.DescribeKeyPairsInput{
		KeyPairIds: []string{keyPairID},
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range output.KeyPairs {
		resource := eC2KeyPairHandle(ctx, v)
		values = append(values, resource)
	}

	return values, nil
}

func EC2AMI(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	output, err := client.DescribeImages(ctx, &ec2.DescribeImagesInput{
		Owners: []string{"self"},
	})
	if err != nil {
		if isErr(err, "InvalidAMIID.NotFound") || isErr(err, "InvalidAMIID.Unavailable") || isErr(err, "InvalidAMIID.Malformed") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource

	for _, v := range output.Images {
		imageAttribute, err := client.DescribeImageAttribute(ctx, &ec2.DescribeImageAttributeInput{
			Attribute: types.ImageAttributeNameLaunchPermission,
			ImageId:   v.ImageId,
		})
		if err != nil {
			if isErr(err, "InvalidAMIID.NotFound") || isErr(err, "InvalidAMIID.Unavailable") || isErr(err, "InvalidAMIID.Malformed") {
				continue
			}
			return nil, err
		}
		resource := eC2AMIHandle(ctx, v, imageAttribute)
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
func eC2AMIHandle(ctx context.Context, v types.Image, imageAttribute *ec2.DescribeImageAttributeOutput) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := "arn:" + describeCtx.Partition + ":ec2:" + describeCtx.Region + ":" + describeCtx.AccountID + ":image/" + *v.ImageId
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *v.ImageId,
		Description: model.EC2AMIDescription{
			AMI:               v,
			LaunchPermissions: *imageAttribute,
		},
	}

	return resource
}
func GetEC2AMI(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	AMIId := fields["id"]
	client := ec2.NewFromConfig(cfg)

	out, err := client.DescribeImages(ctx, &ec2.DescribeImagesInput{
		ImageIds: []string{AMIId},
	})
	if err != nil {
		if isErr(err, "InvalidAMIID.NotFound") || isErr(err, "InvalidAMIID.Unavailable") || isErr(err, "InvalidAMIID.Malformed") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	for _, v := range out.Images {
		imageAttribute, err := client.DescribeImageAttribute(ctx, &ec2.DescribeImageAttributeInput{
			Attribute: types.ImageAttributeNameLaunchPermission,
			ImageId:   v.ImageId,
		})
		if err != nil {
			if isErr(err, "InvalidAMIID.NotFound") || isErr(err, "InvalidAMIID.Unavailable") || isErr(err, "InvalidAMIID.Malformed") {
				continue
			}
			return nil, err
		}
		resource := eC2AMIHandle(ctx, v, imageAttribute)
		values = append(values, resource)
	}
	return values, nil
}

func EC2ReservedInstances(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	output, err := client.DescribeReservedInstances(ctx, &ec2.DescribeReservedInstancesInput{})
	if err != nil {
		if isErr(err, "InvalidParameterValue") || isErr(err, "InvalidInstanceID.Unavailable") || isErr(err, "InvalidInstanceID.Malformed") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource

	filterName := "reserved-instances-id"
	for _, v := range output.ReservedInstances {
		var modifications []types.ReservedInstancesModification
		modificationPaginator := ec2.NewDescribeReservedInstancesModificationsPaginator(client, &ec2.DescribeReservedInstancesModificationsInput{
			Filters: []types.Filter{
				{
					Name:   &filterName,
					Values: []string{*v.ReservedInstancesId},
				},
			},
		})
		for modificationPaginator.HasMorePages() {
			page, err := modificationPaginator.NextPage(ctx)
			if err != nil {
				if isErr(err, "InvalidParameterValue") || isErr(err, "InvalidInstanceID.Unavailable") || isErr(err, "InvalidInstanceID.Malformed") {
					continue
				}
				return nil, err
			}

			modifications = append(modifications, page.ReservedInstancesModifications...)
		}
		resource := eC2ReservedInstancesHandle(ctx, v, modifications)

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
func eC2ReservedInstancesHandle(ctx context.Context, v types.ReservedInstances, modifications []types.ReservedInstancesModification) Resource {
	describeCtx := GetDescribeContext(ctx)

	arn := "arn:" + describeCtx.Partition + ":ec2:" + describeCtx.Region + ":" + describeCtx.AccountID + ":reserved-instances/" + *v.ReservedInstancesId
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *v.ReservedInstancesId,
		Description: model.EC2ReservedInstancesDescription{
			ReservedInstances:   v,
			ModificationDetails: modifications,
		},
	}
	return resource
}
func GetEC2ReservedInstances(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	ReservedInstancesId := fields["id"]
	client := ec2.NewFromConfig(cfg)

	out, err := client.DescribeReservedInstances(ctx, &ec2.DescribeReservedInstancesInput{
		ReservedInstancesIds: []string{ReservedInstancesId},
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range out.ReservedInstances {
		filterName := "reserved-instances-id"
		var modifications []types.ReservedInstancesModification
		modificationPaginator := ec2.NewDescribeReservedInstancesModificationsPaginator(client, &ec2.DescribeReservedInstancesModificationsInput{
			Filters: []types.Filter{
				{
					Name:   &filterName,
					Values: []string{*v.ReservedInstancesId},
				},
			},
		})
		page, err := modificationPaginator.NextPage(ctx)
		if err != nil {
			if isErr(err, "InvalidParameterValue") || isErr(err, "InvalidInstanceID.Unavailable") || isErr(err, "InvalidInstanceID.Malformed") {
				continue
			}
			return nil, err
		}
		modifications = append(modifications, page.ReservedInstancesModifications...)
		resource := eC2ReservedInstancesHandle(ctx, v, modifications)
		values = append(values, resource)
	}
	return values, nil
}

func EC2IpamPool(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeIpamPoolsPaginator(client, &ec2.DescribeIpamPoolsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.IpamPools {
			resource := eC2IpamPoolHandle(ctx, v)
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
func eC2IpamPoolHandle(ctx context.Context, v types.IpamPool) Resource {
	describeCtx := GetDescribeContext(ctx)
	values := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *v.IpamPoolArn,
		Name:   *v.IpamPoolId,
		Description: model.EC2IpamPoolDescription{
			IpamPool: v,
		},
	}
	return values
}
func GetEC2IpamPool(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	IpamPoolId := fields["id"]
	client := ec2.NewFromConfig(cfg)

	out, err := client.DescribeIpamPools(ctx, &ec2.DescribeIpamPoolsInput{
		IpamPoolIds: []string{IpamPoolId},
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range out.IpamPools {
		resource := eC2IpamPoolHandle(ctx, v)
		values = append(values, resource)
	}
	return values, nil
}

func EC2Ipam(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeIpamsPaginator(client, &ec2.DescribeIpamsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Ipams {
			resource := eC2IpamHandle(ctx, v)
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
func eC2IpamHandle(ctx context.Context, v types.Ipam) Resource {
	describeCtx := GetDescribeContext(ctx)
	values := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *v.IpamArn,
		Name:   *v.IpamId,
		Description: model.EC2IpamDescription{
			Ipam: v,
		},
	}
	return values
}
func GetEC2Ipam(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	IpamId := fields["id"]
	client := ec2.NewFromConfig(cfg)

	out, err := client.DescribeIpams(ctx, &ec2.DescribeIpamsInput{
		IpamIds: []string{IpamId},
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range out.Ipams {
		resource := eC2IpamHandle(ctx, v)
		values = append(values, resource)
	}
	return values, nil
}

func EC2InstanceAvailability(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)

	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeInstanceTypeOfferingsPaginator(client, &ec2.DescribeInstanceTypeOfferingsInput{})
	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.InstanceTypeOfferings {
			arn := fmt.Sprintf("arn:%s:ec2:%s::instance-type/%s", describeCtx.Partition, *v.Location, v.InstanceType)
			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    arn,
				Name:   fmt.Sprintf("%s (%s)", v.InstanceType, *v.Location),
				Description: model.EC2InstanceAvailabilityDescription{
					InstanceAvailability: v,
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

func EC2InstanceType(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)

	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeInstanceTypesPaginator(client, &ec2.DescribeInstanceTypesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.InstanceTypes {
			arn := fmt.Sprintf("arn:%s:ec2:::instance-type/%s", describeCtx.Partition, v.InstanceType)
			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    arn,
				Name:   string(v.InstanceType),
				Description: model.EC2InstanceTypeDescription{
					InstanceType: v,
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

func EC2ManagedPrefixList(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeManagedPrefixListsPaginator(client, &ec2.DescribeManagedPrefixListsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.PrefixLists {
			resource := eC2ManagedPrefixListHandle(ctx, v)
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
func eC2ManagedPrefixListHandle(ctx context.Context, v types.ManagedPrefixList) Resource {
	describeCtx := GetDescribeContext(ctx)
	values := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *v.PrefixListArn,
		Name:   *v.PrefixListName,
		Description: model.EC2ManagedPrefixListDescription{
			ManagedPrefixList: v,
		},
	}
	return values
}
func GetEC2ManagedPrefixList(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	ManagedPrefixListId := fields["id"]
	client := ec2.NewFromConfig(cfg)

	out, err := client.DescribeManagedPrefixLists(ctx, &ec2.DescribeManagedPrefixListsInput{
		PrefixListIds: []string{ManagedPrefixListId},
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range out.PrefixLists {
		resource := eC2ManagedPrefixListHandle(ctx, v)
		values = append(values, resource)
	}
	return values, nil
}

func EC2SpotPrice(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ec2.NewFromConfig(cfg)

	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -1)
	paginator := ec2.NewDescribeSpotPriceHistoryPaginator(client, &ec2.DescribeSpotPriceHistoryInput{
		StartTime: &startTime,
		EndTime:   &endTime,
	})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.SpotPriceHistory {
			if v.SpotPrice == nil {
				continue
			}
			avZone := ""
			if v.AvailabilityZone != nil {
				avZone = *v.AvailabilityZone
			}
			resource := Resource{
				Region: describeCtx.KaytuRegion,
				Name:   fmt.Sprintf("%s-%s (%s)", v.InstanceType, *v.SpotPrice, avZone),
				Description: model.EC2SpotPriceDescription{
					SpotPrice: v,
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

// TODO mahan : check that this function implemented correctly
func EC2TransitGatewayRoute(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)

	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeTransitGatewayRouteTablesPaginator(client, &ec2.DescribeTransitGatewayRouteTablesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, transitGatewayRouteTable := range page.TransitGatewayRouteTables {
			routes, err := client.SearchTransitGatewayRoutes(ctx, &ec2.SearchTransitGatewayRoutesInput{
				Filters: []types.Filter{
					{
						Name:   aws.String("state"),
						Values: []string{"active", "blackhole", "pending"},
					},
				},
				TransitGatewayRouteTableId: transitGatewayRouteTable.TransitGatewayRouteTableId,
			})
			if err != nil {
				return nil, err
			}
			for _, route := range routes.Routes {
				arn := fmt.Sprintf("arn:%s:ec2:%s:%s:transit-gateway-route-table/%s:%s", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID, *transitGatewayRouteTable.TransitGatewayRouteTableId, *route.DestinationCidrBlock)
				resource := Resource{
					Region: describeCtx.KaytuRegion,
					ARN:    arn,
					Name:   *route.DestinationCidrBlock,
					Description: model.EC2TransitGatewayRouteDescription{
						TransitGatewayRoute:        route,
						TransitGatewayRouteTableId: *transitGatewayRouteTable.TransitGatewayRouteTableId,
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
	}
	return values, nil
}
func GetEC2TransitGatewayRoute(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	TransitGatewayRouteId := fields["id"]
	client := ec2.NewFromConfig(cfg)

	routTable, err := client.DescribeTransitGatewayRouteTables(ctx, &ec2.DescribeTransitGatewayRouteTablesInput{
		TransitGatewayRouteTableIds: []string{TransitGatewayRouteId},
	})
	if err != nil {
		return nil, err
	}
	if len(routTable.TransitGatewayRouteTables) == 0 {
		return nil, nil
	}
	var values []Resource
	for _, v := range routTable.TransitGatewayRouteTables {
		arn := fmt.Sprintf("arn:%s:ec2:%s:%s:transit-gateway-route-table/%s:%s", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID, *v.TransitGatewayRouteTableId)
		values = append(values, Resource{
			Region: describeCtx.KaytuRegion,
			ARN:    arn,
			Name:   *v.TransitGatewayRouteTableId,
			Description: model.EC2TransitGatewayRouteTableDescription{
				TransitGatewayRouteTable: v,
			},
		})
	}
	return values, nil
}

func EC2TransitGatewayAttachment(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeTransitGatewayAttachmentsPaginator(client, &ec2.DescribeTransitGatewayAttachmentsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.TransitGatewayAttachments {
			resource := eC2TransitGatewayAttachmentHandle(ctx, v)
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
func eC2TransitGatewayAttachmentHandle(ctx context.Context, v types.TransitGatewayAttachment) Resource {
	var values Resource
	describeCtx := GetDescribeContext(ctx)
	arn := fmt.Sprintf("arn:%s:ec2:%s:%s:transit-gateway-attachment/%s", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID, *v.TransitGatewayAttachmentId)
	values = Resource{
		Region: describeCtx.KaytuRegion,
		ID:     *v.TransitGatewayAttachmentId,
		ARN:    arn,
		Description: model.EC2TransitGatewayAttachmentDescription{
			TransitGatewayAttachment: v,
		},
	}
	return values
}
func GetEC2TransitGatewayAttachment(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	TransitGatewayAttachmentId := fields["id"]
	client := ec2.NewFromConfig(cfg)

	out, err := client.DescribeTransitGatewayAttachments(ctx, &ec2.DescribeTransitGatewayAttachmentsInput{
		TransitGatewayAttachmentIds: []string{TransitGatewayAttachmentId},
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range out.TransitGatewayAttachments {
		resource := eC2TransitGatewayAttachmentHandle(ctx, v)
		values = append(values, resource)
	}
	return values, nil
}

func EC2LaunchTemplate(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeLaunchTemplatesPaginator(client, &ec2.DescribeLaunchTemplatesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.LaunchTemplates {
			resource := eC2LaunchTemplateHandle(ctx, v)
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
func eC2LaunchTemplateHandle(ctx context.Context, v types.LaunchTemplate) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := fmt.Sprintf("arn:%s:ec2:%s:%s:launch-template/%s", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID, *v.LaunchTemplateId)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ID:     *v.LaunchTemplateId,
		ARN:    arn,
		Name:   *v.LaunchTemplateName,
		Description: model.EC2LaunchTemplateDescription{
			LaunchTemplate: v,
		},
	}
	return resource
}

func GetEC2LaunchTemplate(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	LaunchTemplateId := fields["id"]
	client := ec2.NewFromConfig(cfg)

	out, err := client.DescribeLaunchTemplates(ctx, &ec2.DescribeLaunchTemplatesInput{
		LaunchTemplateIds: []string{LaunchTemplateId},
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range out.LaunchTemplates {
		resource := eC2LaunchTemplateHandle(ctx, v)
		values = append(values, resource)
	}
	return values, nil
}

func EbsVolumeMetricReadOps(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeVolumesPaginator(client, &ec2.DescribeVolumesInput{
		MaxResults: aws.Int32(500),
	})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Volumes {
			metrics, err := listCloudWatchMetricStatistics(ctx, cfg, "5_MIN", "AWS/EBS", "VolumeReadOps", "VolumeId", *v.VolumeId)
			if err != nil {
				return nil, err
			}
			for _, metric := range metrics {
				resource := Resource{
					Region: describeCtx.KaytuRegion,
					ID:     fmt.Sprintf("%s:%s:%s:%s", *v.VolumeId, metric.Timestamp.Format(time.RFC3339), *metric.DimensionName, *metric.DimensionValue),
					Description: model.EbsVolumeMetricReadOpsDescription{
						CloudWatchMetricRow: metric,
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
	}

	return values, nil
}

func EbsVolumeMetricReadOpsDaily(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeVolumesPaginator(client, &ec2.DescribeVolumesInput{
		MaxResults: aws.Int32(500),
	})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Volumes {
			metrics, err := listCloudWatchMetricStatistics(ctx, cfg, "DAILY", "AWS/EBS", "VolumeReadOps", "VolumeId", *v.VolumeId)
			if err != nil {
				return nil, err
			}
			for _, metric := range metrics {
				resource := Resource{
					Region: describeCtx.KaytuRegion,
					ID:     fmt.Sprintf("%s:%s:%s:%s-daily", *v.VolumeId, metric.Timestamp.Format(time.RFC3339), *metric.DimensionName, *metric.DimensionValue),
					Description: model.EbsVolumeMetricReadOpsDailyDescription{
						CloudWatchMetricRow: metric,
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
	}

	return values, nil
}

func EbsVolumeMetricReadOpsHourly(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeVolumesPaginator(client, &ec2.DescribeVolumesInput{
		MaxResults: aws.Int32(500),
	})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Volumes {
			metrics, err := listCloudWatchMetricStatistics(ctx, cfg, "HOURLY", "AWS/EBS", "VolumeReadOps", "VolumeId", *v.VolumeId)
			if err != nil {
				return nil, err
			}
			for _, metric := range metrics {
				resource := Resource{
					Region: describeCtx.KaytuRegion,
					ID:     fmt.Sprintf("%s:%s:%s:%s-hourly", *v.VolumeId, metric.Timestamp.Format(time.RFC3339), *metric.DimensionName, *metric.DimensionValue),
					Description: model.EbsVolumeMetricReadOpsHourlyDescription{
						CloudWatchMetricRow: metric,
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
	}

	return values, nil
}

func EbsVolumeMetricWriteOps(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeVolumesPaginator(client, &ec2.DescribeVolumesInput{
		MaxResults: aws.Int32(500),
	})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Volumes {
			metrics, err := listCloudWatchMetricStatistics(ctx, cfg, "5_MIN", "AWS/EBS", "VolumeWriteOps", "VolumeId", *v.VolumeId)
			if err != nil {
				return nil, err
			}
			for _, metric := range metrics {
				resource := Resource{
					Region: describeCtx.KaytuRegion,
					ID:     fmt.Sprintf("%s:%s:%s:%s", *v.VolumeId, metric.Timestamp.Format(time.RFC3339), *metric.DimensionName, *metric.DimensionValue),
					Description: model.EbsVolumeMetricWriteOpsDescription{
						CloudWatchMetricRow: metric,
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
	}

	return values, nil
}

func EbsVolumeMetricWriteOpsDaily(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeVolumesPaginator(client, &ec2.DescribeVolumesInput{
		MaxResults: aws.Int32(500),
	})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Volumes {
			metrics, err := listCloudWatchMetricStatistics(ctx, cfg, "DAILY", "AWS/EBS", "VolumeWriteOps", "VolumeId", *v.VolumeId)
			if err != nil {
				return nil, err
			}
			for _, metric := range metrics {
				resource := Resource{
					Region: describeCtx.KaytuRegion,
					ID:     fmt.Sprintf("%s:%s:%s:%s-daily", *v.VolumeId, metric.Timestamp.Format(time.RFC3339), *metric.DimensionName, *metric.DimensionValue),
					Description: model.EbsVolumeMetricWriteOpsDailyDescription{
						CloudWatchMetricRow: metric,
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
	}

	return values, nil
}

func EbsVolumeMetricWriteOpsHourly(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeVolumesPaginator(client, &ec2.DescribeVolumesInput{
		MaxResults: aws.Int32(500),
	})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Volumes {
			metrics, err := listCloudWatchMetricStatistics(ctx, cfg, "HOURLY", "AWS/EBS", "VolumeWriteOps", "VolumeId", *v.VolumeId)
			if err != nil {
				return nil, err
			}
			for _, metric := range metrics {
				resource := Resource{
					Region: describeCtx.KaytuRegion,
					ID:     fmt.Sprintf("%s:%s:%s:%s-hourly", *v.VolumeId, metric.Timestamp.Format(time.RFC3339), *metric.DimensionName, *metric.DimensionValue),
					Description: model.EbsVolumeMetricWriteOpsHourlyDescription{
						CloudWatchMetricRow: metric,
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
	}

	return values, nil
}
func EC2VPCNatGatewayMetricBytesOutToDestination(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeNatGatewaysPaginator(client, &ec2.DescribeNatGatewaysInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, v := range page.NatGateways {
			resource := eC2VPCNatGatewayMetricBytesOutToDestinationHandle(ctx, v)
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
func eC2VPCNatGatewayMetricBytesOutToDestinationHandle(ctx context.Context, v types.NatGateway) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ID:     *v.NatGatewayId,
		Description: model.EC2NatGatewayMetricBytesOutToDestinationDescription{
			NatGateway: v,
		},
	}
	return resource
}
func GetEC2VPCNatGatewayMetricBytesOutToDestination(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	natGatewayId := fields["id"]
	client := ec2.NewFromConfig(cfg)
	natGateway, err := client.DescribeNatGateways(ctx, &ec2.DescribeNatGatewaysInput{
		NatGatewayIds: []string{natGatewayId},
	})
	if err != nil {
		if isErr(err, "DescribeNatGatewaysNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	for _, v := range natGateway.NatGateways {
		resource := eC2VPCNatGatewayMetricBytesOutToDestinationHandle(ctx, v)
		values = append(values, resource)
	}
	return values, nil
}

func EC2LaunchTemplateVersion(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeLaunchTemplatesPaginator(client, &ec2.DescribeLaunchTemplatesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, t := range page.LaunchTemplates {
			version, err := client.DescribeLaunchTemplateVersions(ctx, &ec2.DescribeLaunchTemplateVersionsInput{
				LaunchTemplateId: t.LaunchTemplateId,
			})
			if err != nil {
				return nil, err
			}
			for _, v := range version.LaunchTemplateVersions {
				resource := eC2LaunchTemplateVersionHandle(ctx, v)
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
	return values, nil
}

func eC2LaunchTemplateVersionHandle(ctx context.Context, v types.LaunchTemplateVersion) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		Description: model.EC2LaunchTemplateVersionDescription{
			LaunchTemplateVersion: v,
		},
	}
	return resource
}

func GetEC2LaunchTemplateVersion(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	LaunchTemplateId := fields["id"]
	client := ec2.NewFromConfig(cfg)

	out, err := client.DescribeLaunchTemplateVersions(ctx, &ec2.DescribeLaunchTemplateVersionsInput{
		LaunchTemplateId: aws.String(LaunchTemplateId),
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range out.LaunchTemplateVersions {
		resource := eC2LaunchTemplateVersionHandle(ctx, v)
		values = append(values, resource)
	}
	return values, nil
}

func EC2ManagedPrefixListEntry(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeManagedPrefixListsPaginator(client, &ec2.DescribeManagedPrefixListsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.PrefixLists {
			enPaginator := ec2.NewGetManagedPrefixListEntriesPaginator(client, &ec2.GetManagedPrefixListEntriesInput{
				PrefixListId: v.PrefixListId,
			})
			for enPaginator.HasMorePages() {
				enPage, err := enPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}
				for _, entry := range enPage.Entries {
					resource := eC2ManagedPrefixListEntryHandle(ctx, *v.PrefixListId, entry)
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
func eC2ManagedPrefixListEntryHandle(ctx context.Context, prefixListId string, v types.PrefixListEntry) Resource {
	describeCtx := GetDescribeContext(ctx)
	values := Resource{
		Region: describeCtx.KaytuRegion,
		Description: model.EC2ManagedPrefixListEntryDescription{
			PrefixListEntry: v,
			PrefixListId:    prefixListId,
		},
	}
	return values
}

func Ec2InstanceMetricCpuUtilizationHourly(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeInstancesPaginator(client, &ec2.DescribeInstancesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, r := range page.Reservations {
			for _, v := range r.Instances {
				statistics, err := getEc2InstanceMetricCpuUtilizationHourly(ctx, cfg, v.InstanceId)
				if err != nil {
					return nil, err
				}
				for _, resource := range statistics {
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

func getEc2InstanceMetricCpuUtilizationHourly(ctx context.Context, cfg aws.Config, instanceId *string) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	statistics, err := listCloudWatchMetricStatistics(ctx, cfg, "HOURLY", "AWS/EC2", "CPUUtilization", "InstanceId", *instanceId)
	if err != nil {
		return nil, err
	}
	var values []Resource

	for _, s := range statistics {
		values = append(values, Resource{
			Account: describeCtx.AccountID,
			Region:  describeCtx.KaytuRegion,
			Description: model.EC2InstanceMetricCpuUtilizationHourlyDescription{
				InstanceId:  instanceId,
				Timestamp:   s.Timestamp,
				Sum:         s.Sum,
				Average:     s.Average,
				Maximum:     s.Maximum,
				Minimum:     s.Minimum,
				SampleCount: s.SampleCount,
			},
		})
	}

	return values, nil
}
