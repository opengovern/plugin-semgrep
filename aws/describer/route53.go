package describer

import (
	"context"
	"errors"
	"fmt"
	types3 "github.com/aws/aws-sdk-go-v2/service/route53domains/types"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/route53domains"
	"github.com/aws/smithy-go"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/aws/aws-sdk-go-v2/service/route53resolver"
	resolvertypes "github.com/aws/aws-sdk-go-v2/service/route53resolver/types"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func Route53HealthCheck(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := route53.NewFromConfig(cfg)

	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		output, err := client.ListHealthChecks(ctx, &route53.ListHealthChecksInput{Marker: prevToken})
		if err != nil {
			return nil, err
		}

		for _, v := range output.HealthChecks {
			item, err := client.GetHealthCheckStatus(ctx, &route53.GetHealthCheckStatusInput{
				HealthCheckId: v.Id,
			})
			if err != nil {
				var ae smithy.APIError
				if errors.As(err, &ae) {
					if ae.ErrorCode() == "InvalidInput" {
						item = nil
					} else {
						return nil, err
					}
				} else {
					return nil, err
				}
			}
			if item == nil {
				item = &route53.GetHealthCheckStatusOutput{}
			}

			resp, err := client.ListTagsForResource(ctx, &route53.ListTagsForResourceInput{
				ResourceId:   v.Id,
				ResourceType: "healthcheck",
			})
			if err != nil {
				return nil, err
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ID:     *v.Id,
				Description: model.Route53HealthCheckDescription{
					HealthCheck: v,
					Status:      item,
					Tags:        resp,
				},
			}
			if v.HealthCheckConfig != nil && v.HealthCheckConfig.FullyQualifiedDomainName != nil {
				resource.Name = *v.HealthCheckConfig.FullyQualifiedDomainName
			}
			if stream != nil {
				if err := (*stream)(resource); err != nil {
					return nil, err
				}
			} else {
				values = append(values, resource)
			}

		}

		return output.NextMarker, nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}

func Route53HostedZone(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := route53.NewFromConfig(cfg)

	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		output, err := client.ListHostedZones(ctx, &route53.ListHostedZonesInput{Marker: prevToken})
		if err != nil {
			if !isErr(err, "NoSuchHostedZone") {
				return nil, err
			}
			return nil, nil
		}

		for _, v := range output.HostedZones {
			resource, err := route53HostedZoneHandle(ctx, cfg, v)
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

		return output.NextMarker, nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}
func route53HostedZoneHandle(ctx context.Context, cfg aws.Config, v types.HostedZone) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)

	client := route53.NewFromConfig(cfg)

	id := strings.Split(*v.Id, "/")[2]
	arn := fmt.Sprintf("arn:%s:route53:::hostedzone/%s", describeCtx.Partition, id)

	queryLoggingConfigs, err := client.ListQueryLoggingConfigs(ctx, &route53.ListQueryLoggingConfigsInput{
		HostedZoneId: &id,
	})
	if err != nil {
		if !isErr(err, "NoSuchHostedZone") {
			return Resource{}, err
		}
		queryLoggingConfigs = &route53.ListQueryLoggingConfigsOutput{}
	}

	limit, err := client.GetHostedZoneLimit(ctx, &route53.GetHostedZoneLimitInput{
		HostedZoneId: &id,
		Type:         types.HostedZoneLimitTypeMaxRrsetsByZone,
	})
	if err != nil {
		if !isErr(err, "NoSuchHostedZone") {
			return Resource{}, err
		}
		limit = &route53.GetHostedZoneLimitOutput{}
	}

	dnsSec := &route53.GetDNSSECOutput{}
	if !v.Config.PrivateZone {
		dnsSec, err = client.GetDNSSEC(ctx, &route53.GetDNSSECInput{
			HostedZoneId: &id,
		})
		if err != nil {
			if !isErr(err, "NoSuchHostedZone") && !isErr(err, "AccessDenied") {
				return Resource{}, err
			}
			dnsSec = &route53.GetDNSSECOutput{}
		}
	}

	tags, err := client.ListTagsForResource(ctx, &route53.ListTagsForResourceInput{
		ResourceId:   &id,
		ResourceType: types.TagResourceType("hostedzone"),
	})
	if err != nil {
		if !isErr(err, "NoSuchHostedZone") {
			return Resource{}, err
		}
		tags = &route53.ListTagsForResourceOutput{
			ResourceTagSet: &types.ResourceTagSet{},
		}
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *v.Name,
		Description: model.Route53HostedZoneDescription{
			ID:                  id,
			HostedZone:          v,
			QueryLoggingConfigs: queryLoggingConfigs.QueryLoggingConfigs,
			Limit:               limit.Limit,
			DNSSec:              *dnsSec,
			Tags:                tags.ResourceTagSet.Tags,
		},
	}
	return resource, nil
}
func GetRoute53HostedZone(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	hostedZoneID := fields["hostedZoneId"]
	client := route53.NewFromConfig(cfg)

	var values []Resource
	out, err := client.GetHostedZone(ctx, &route53.GetHostedZoneInput{Id: &hostedZoneID})
	if err != nil {
		return nil, err
	}

	v := out.HostedZone

	resource, err := route53HostedZoneHandle(ctx, cfg, *v)
	if err != nil {
		return nil, err
	}

	values = append(values, resource)
	return values, nil
}

func Route53DNSSEC(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	zones, err := Route53HostedZone(ctx, cfg, nil)
	if err != nil {
		return nil, err
	}

	client := route53.NewFromConfig(cfg)

	var values []Resource
	for _, zone := range zones {
		id := zone.Description.(types.HostedZone).Id
		v, err := client.GetDNSSEC(ctx, &route53.GetDNSSECInput{
			HostedZoneId: id,
		})
		if err != nil {
			return nil, err
		}

		resource := Resource{
			Region:      describeCtx.KaytuRegion,
			ID:          *id, // Unique per HostedZone
			Name:        *id,
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

	return values, nil
}

func Route53RecordSet(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	zones, err := Route53HostedZone(ctx, cfg, nil)
	if err != nil {
		return nil, err
	}

	client := route53.NewFromConfig(cfg)

	var values []Resource
	for _, zone := range zones {
		id := zone.Description.(types.HostedZone).Id
		var prevType types.RRType
		err = PaginateRetrieveAll(func(prevName *string) (nextName *string, err error) {
			output, err := client.ListResourceRecordSets(ctx, &route53.ListResourceRecordSetsInput{
				HostedZoneId:    id,
				StartRecordName: prevName,
				StartRecordType: prevType,
			})
			if err != nil {
				return nil, err
			}

			for _, v := range output.ResourceRecordSets {
				resource := Resource{
					Region:      describeCtx.KaytuRegion,
					ID:          CompositeID(*id, *v.Name),
					Name:        *v.Name,
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

			prevType = output.NextRecordType
			return output.NextRecordName, nil
		})
		if err != nil {
			return nil, err
		}
	}

	return values, nil
}

func Route53ResolverFirewallDomainList(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := route53resolver.NewFromConfig(cfg)
	paginator := route53resolver.NewListFirewallDomainListsPaginator(client, &route53resolver.ListFirewallDomainListsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.FirewallDomainLists {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ARN:         *v.Arn,
				Name:        *v.Name,
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

func Route53ResolverFirewallRuleGroup(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := route53resolver.NewFromConfig(cfg)
	paginator := route53resolver.NewListFirewallRuleGroupsPaginator(client, &route53resolver.ListFirewallRuleGroupsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.FirewallRuleGroups {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ARN:         *v.Arn,
				Name:        *v.Name,
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

func Route53ResolverFirewallRuleGroupAssociation(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := route53resolver.NewFromConfig(cfg)
	paginator := route53resolver.NewListFirewallRuleGroupAssociationsPaginator(client, &route53resolver.ListFirewallRuleGroupAssociationsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.FirewallRuleGroupAssociations {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ARN:         *v.Arn,
				Name:        *v.Name,
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

func Route53ResolverResolverDNSSECConfig(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	vpcs, err := EC2VPC(ctx, cfg, nil)
	if err != nil {
		return nil, err
	}

	client := route53resolver.NewFromConfig(cfg)

	var values []Resource
	for _, vpc := range vpcs {
		v, err := client.GetResolverDnssecConfig(ctx, &route53resolver.GetResolverDnssecConfigInput{
			ResourceId: vpc.Description.(model.EC2VpcDescription).Vpc.VpcId,
		})
		if err != nil {
			return nil, err
		}

		resource := Resource{
			Region:      describeCtx.KaytuRegion,
			ID:          *v.ResolverDNSSECConfig.Id,
			Name:        *v.ResolverDNSSECConfig.Id,
			Description: v.ResolverDNSSECConfig,
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

func Route53ResolverResolverQueryLoggingConfig(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := route53resolver.NewFromConfig(cfg)
	paginator := route53resolver.NewListResolverQueryLogConfigsPaginator(client, &route53resolver.ListResolverQueryLogConfigsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.ResolverQueryLogConfigs {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ARN:         *v.Arn,
				Name:        *v.Name,
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

func Route53ResolverResolverQueryLoggingConfigAssociation(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := route53resolver.NewFromConfig(cfg)
	paginator := route53resolver.NewListResolverQueryLogConfigAssociationsPaginator(client, &route53resolver.ListResolverQueryLogConfigAssociationsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.ResolverQueryLogConfigAssociations {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ID:          *v.Id,
				Name:        *v.Id,
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

func Route53ResolverResolverRule(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := route53resolver.NewFromConfig(cfg)
	paginator := route53resolver.NewListResolverRulesPaginator(client, &route53resolver.ListResolverRulesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.ResolverRules {
			defaultID := "rslvr-autodefined-rr-internet-resolver"

			var tags []resolvertypes.Tag
			if *v.Id != defaultID {
				tagsOut, err := client.ListTagsForResource(ctx, &route53resolver.ListTagsForResourceInput{
					ResourceArn: v.Arn,
				})
				if err != nil {
					return nil, err
				}
				tags = tagsOut.Tags
			}

			// Build the params
			params := &route53resolver.ListResolverRuleAssociationsInput{
				Filters: []resolvertypes.Filter{
					{
						Name: aws.String("ResolverRuleId"),
						Values: []string{
							*v.Id,
						},
					},
				},
			}

			ruleass, err := client.ListResolverRuleAssociations(ctx, params)
			if err != nil {
				return nil, err
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    *v.Arn,
				Name:   *v.Name,
				Description: model.Route53ResolverResolverRuleDescription{
					ResolverRole:     v,
					Tags:             tags,
					RuleAssociations: ruleass,
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

func Route53ResolverResolverEndpoint(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := route53resolver.NewFromConfig(cfg)
	paginator := route53resolver.NewListResolverEndpointsPaginator(client, &route53resolver.ListResolverEndpointsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, resolverEndpoint := range page.ResolverEndpoints {
			resource := route53ResolverResolverEndpointHandle(ctx, cfg, resolverEndpoint)

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
func route53ResolverResolverEndpointHandle(ctx context.Context, cfg aws.Config, resolverEndpoint resolvertypes.ResolverEndpoint) Resource {
	describeCtx := GetDescribeContext(ctx)
	client := route53resolver.NewFromConfig(cfg)

	ipAddresesses, err := client.ListResolverEndpointIpAddresses(ctx, &route53resolver.ListResolverEndpointIpAddressesInput{
		ResolverEndpointId: resolverEndpoint.Id,
	})
	if err != nil {
		ipAddresesses = &route53resolver.ListResolverEndpointIpAddressesOutput{}
	}

	tags, err := client.ListTagsForResource(ctx, &route53resolver.ListTagsForResourceInput{
		ResourceArn: resolverEndpoint.Arn,
	})
	if err != nil {
		tags = &route53resolver.ListTagsForResourceOutput{}
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *resolverEndpoint.Arn,
		Name:   *resolverEndpoint.Name,
		ID:     *resolverEndpoint.Id,
		Description: model.Route53ResolverEndpointDescription{
			ResolverEndpoint: resolverEndpoint,
			IpAddresses:      ipAddresesses.IpAddresses,
			Tags:             tags.Tags,
		},
	}
	return resource
}
func GetRoute53ResolverResolverEndpoint(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	resolverEndpointId := fields["id"]
	client := route53resolver.NewFromConfig(cfg)
	var values []Resource

	out, err := client.GetResolverEndpoint(ctx, &route53resolver.GetResolverEndpointInput{
		ResolverEndpointId: &resolverEndpointId,
	})
	if err != nil {
		return nil, err
	}

	resource := route53ResolverResolverEndpointHandle(ctx, cfg, *out.ResolverEndpoint)
	values = append(values, resource)
	return values, nil
}

func Route53ResolverResolverRuleAssociation(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := route53resolver.NewFromConfig(cfg)
	paginator := route53resolver.NewListResolverRuleAssociationsPaginator(client, &route53resolver.ListResolverRuleAssociationsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.ResolverRuleAssociations {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ID:          *v.Id,
				Name:        *v.Name,
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

func Route53Domain(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := route53domains.NewFromConfig(cfg)

	paginator := route53domains.NewListDomainsPaginator(client, &route53domains.ListDomainsInput{})
	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, v := range page.Domains {
			resource, err := Route53DomainHandle(ctx, cfg, v)
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
func Route53DomainHandle(ctx context.Context, cfg aws.Config, v types3.DomainSummary) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := route53domains.NewFromConfig(cfg)
	domain, err := client.GetDomainDetail(ctx, &route53domains.GetDomainDetailInput{
		DomainName: v.DomainName,
	})
	if err != nil {
		return Resource{}, err
	}

	tags, err := client.ListTagsForDomain(ctx, &route53domains.ListTagsForDomainInput{
		DomainName: v.DomainName,
	})
	if err != nil {
		tags = &route53domains.ListTagsForDomainOutput{}
	}

	arn := fmt.Sprintf("arn:%s:route53domains:::domain/%s", describeCtx.Partition, *v.DomainName)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		Name:   *domain.DomainName,
		ARN:    arn,
		Description: model.Route53DomainDescription{
			DomainSummary: v,
			Domain:        *domain,
			Tags:          tags.TagList,
		},
	}
	return resource, nil
}
func GetRoute53Domain(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	domainName := fields["name"]
	client := route53domains.NewFromConfig(cfg)

	list, err := client.ListDomains(ctx, &route53domains.ListDomainsInput{})
	if err != nil {
		if isErr(err, "ListDomainsNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	for _, v := range list.Domains {
		if *v.DomainName != domainName {
			continue
		}
		resource, err := Route53DomainHandle(ctx, cfg, v)
		if err != nil {
			return nil, err
		}

		values = append(values, resource)
	}
	return values, nil
}

func Route53Record(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := route53.NewFromConfig(cfg)
	paginator := route53.NewListHostedZonesPaginator(client, &route53.ListHostedZonesInput{})
	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, v := range page.HostedZones {
			err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
				records, err := client.ListResourceRecordSets(ctx, &route53.ListResourceRecordSetsInput{
					HostedZoneId:    v.Id,
					StartRecordName: prevToken,
				})
				if err != nil {
					return nil, err
				}
				for _, record := range records.ResourceRecordSets {
					resource := route53RecordHandle(ctx, record, *v.Id)

					if stream != nil {
						if err := (*stream)(resource); err != nil {
							return nil, err
						}
					} else {
						values = append(values, resource)
					}
				}
				if records.IsTruncated {
					return records.NextRecordName, nil
				}
				return nil, nil
			})
			if err != nil {
				return nil, err
			}
		}
	}

	return values, nil
}
func route53RecordHandle(ctx context.Context, record types.ResourceRecordSet, hostedZoneId string) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := fmt.Sprintf("arn:%s:route53:::hostedzone/%s/recordset/%s/%s", describeCtx.Partition, hostedZoneId, *record.Name, record.Type)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		Name:   *record.Name,
		ARN:    arn,
		Description: model.Route53RecordDescription{
			ZoneID: hostedZoneId,
			Record: record,
		},
	}
	return resource
}
func GetRoute53Record(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	id := fields["id"]
	client := route53.NewFromConfig(cfg)

	hostedZone, err := client.GetHostedZone(ctx, &route53.GetHostedZoneInput{
		Id: &id,
	})
	if err != nil {
		if isErr(err, "GetHostedZoneNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	list, err := client.ListResourceRecordSets(ctx, &route53.ListResourceRecordSetsInput{
		HostedZoneId: hostedZone.HostedZone.Id,
	})
	if err != nil {
		if isErr(err, "ListResourceRecordSetsNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	for _, record := range list.ResourceRecordSets {

		resource := route53RecordHandle(ctx, record, *hostedZone.HostedZone.Id)
		values = append(values, resource)

	}
	return values, nil
}

func Route53TrafficPolicy(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := route53.NewFromConfig(cfg)
	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		policies, err := client.ListTrafficPolicies(ctx, &route53.ListTrafficPoliciesInput{
			TrafficPolicyIdMarker: prevToken,
		})
		if err != nil {
			return nil, err
		}
		for _, policySummary := range policies.TrafficPolicySummaries {
			resource, err := route53TrafficPolicyHandle(ctx, cfg, policySummary)
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
		if policies.IsTruncated {
			return policies.TrafficPolicyIdMarker, nil
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}
func route53TrafficPolicyHandle(ctx context.Context, cfg aws.Config, policySummary types.TrafficPolicySummary) (Resource, error) {

	describeCtx := GetDescribeContext(ctx)

	client := route53.NewFromConfig(cfg)
	policy, err := client.GetTrafficPolicy(ctx, &route53.GetTrafficPolicyInput{
		Id:      policySummary.Id,
		Version: policySummary.LatestVersion,
	})
	if err != nil {
		return Resource{}, err
	}

	arn := fmt.Sprintf("arn:%s:route53::%s:trafficpolicy/%s/%s", describeCtx.Partition, describeCtx.AccountID, *policy.TrafficPolicy.Id, string(*policy.TrafficPolicy.Version))
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		Name:   *policy.TrafficPolicy.Name,
		ID:     *policy.TrafficPolicy.Id,
		ARN:    arn,
		Description: model.Route53TrafficPolicyDescription{
			TrafficPolicy: *policy.TrafficPolicy,
		},
	}
	return resource, nil
}
func GetRoute53TrafficPolicy(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	id := fields["id"]
	client := route53.NewFromConfig(cfg)

	list, err := client.ListTrafficPolicies(ctx, &route53.ListTrafficPoliciesInput{})
	if err != nil {
		if isErr(err, "ListTrafficPoliciesNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	for _, policySummary := range list.TrafficPolicySummaries {
		if *policySummary.Id != id {
			continue
		}
		resource, err := route53TrafficPolicyHandle(ctx, cfg, policySummary)
		if err != nil {
			return nil, err
		}
		values = append(values, resource)
	}
	return values, nil
}

func Route53TrafficPolicyInstance(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := route53.NewFromConfig(cfg)
	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		policies, err := client.ListTrafficPolicyInstances(ctx, &route53.ListTrafficPolicyInstancesInput{
			TrafficPolicyInstanceNameMarker: prevToken,
		})
		if err != nil {
			return nil, err
		}
		for _, policyInstance := range policies.TrafficPolicyInstances {
			resource := route53TrafficPolicyInstanceHandle(ctx, policyInstance)

			if stream != nil {
				if err := (*stream)(resource); err != nil {
					return nil, err
				}
			} else {
				values = append(values, resource)
			}
		}
		if policies.IsTruncated {
			return policies.TrafficPolicyInstanceNameMarker, nil
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}
func route53TrafficPolicyInstanceHandle(ctx context.Context, policyInstance types.TrafficPolicyInstance) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := fmt.Sprintf("arn:%s:route53::%s:trafficpolicyinstance/%s", describeCtx.Partition, describeCtx.AccountID, *policyInstance.Id)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		Name:   *policyInstance.Name,
		ID:     *policyInstance.Id,
		ARN:    arn,
		Description: model.Route53TrafficPolicyInstanceDescription{
			TrafficPolicyInstance: policyInstance,
		},
	}
	return resource
}
func GetRoute53TrafficPolicyInstance(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	id := fields["id"]
	client := route53.NewFromConfig(cfg)
	var values []Resource

	trafficPolicy, err := client.GetTrafficPolicyInstance(ctx, &route53.GetTrafficPolicyInstanceInput{
		Id: &id,
	})
	if err != nil {
		if isErr(err, "GetTrafficPolicyInstanceNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	resource := route53TrafficPolicyInstanceHandle(ctx, *trafficPolicy.TrafficPolicyInstance)
	values = append(values, resource)
	return values, nil
}

func Route53QueryLog(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := route53.NewFromConfig(cfg)
	paginator := route53.NewListQueryLoggingConfigsPaginator(client, &route53.ListQueryLoggingConfigsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.QueryLoggingConfigs {
			arn := fmt.Sprintf("arn:%s:route53:::query-log/%s/%s", describeCtx.Partition, *v.HostedZoneId, *v.Id)

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ID:     *v.Id,
				ARN:    arn,
				Description: model.Route53QueryLogDescription{
					QueryConfig: v,
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

func Route53ResolverQueryLogConfig(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := route53resolver.NewFromConfig(cfg)
	paginator := route53resolver.NewListResolverQueryLogConfigsPaginator(client, &route53resolver.ListResolverQueryLogConfigsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, queryLogConfig := range page.ResolverQueryLogConfigs {
			resource := route53ResolverQueryLogConfigHandle(ctx, cfg, queryLogConfig)

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
func route53ResolverQueryLogConfigHandle(ctx context.Context, cfg aws.Config, queryLogConfig resolvertypes.ResolverQueryLogConfig) Resource {
	describeCtx := GetDescribeContext(ctx)

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *queryLogConfig.Arn,
		Name:   *queryLogConfig.Name,
		ID:     *queryLogConfig.Id,
		Description: model.Route53ResolverQueryLogConfigDescription{
			QueryConfig: queryLogConfig,
		},
	}
	return resource
}
func GetRoute53ResolverQueryLogConfig(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	id := fields["id"]
	client := route53resolver.NewFromConfig(cfg)
	var values []Resource

	out, err := client.GetResolverQueryLogConfig(ctx, &route53resolver.GetResolverQueryLogConfigInput{
		ResolverQueryLogConfigId: &id,
	})
	if err != nil {
		return nil, err
	}

	resource := route53ResolverQueryLogConfigHandle(ctx, cfg, *out.ResolverQueryLogConfig)
	values = append(values, resource)
	return values, nil
}
