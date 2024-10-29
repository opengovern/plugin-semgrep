package describer

import (
	"context"
	"errors"
	"fmt"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"strings"

	"github.com/aws/smithy-go"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/waf"
	waftypes "github.com/aws/aws-sdk-go-v2/service/waf/types"
	regionaltypes "github.com/aws/aws-sdk-go-v2/service/wafregional/types"

	"github.com/aws/aws-sdk-go-v2/service/wafregional"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	"github.com/aws/aws-sdk-go-v2/service/wafv2/types"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func WAFv2IPSet(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := wafv2.NewFromConfig(cfg)

	scopes := []types.Scope{
		types.ScopeRegional,
	}
	if strings.EqualFold(cfg.Region, "us-east-1") {
		scopes = append(scopes, types.ScopeCloudfront)
	}

	var values []Resource
	for _, scope := range scopes {
		err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
			output, err := client.ListIPSets(ctx, &wafv2.ListIPSetsInput{
				Scope:      scope,
				NextMarker: prevToken,
			})
			if err != nil {
				return nil, err
			}

			for _, v := range output.IPSets {
				params := &wafv2.GetIPSetInput{
					Id:    v.Id,
					Name:  v.Name,
					Scope: scope,
				}

				op, err := client.GetIPSet(ctx, params)
				if err != nil {
					return nil, err
				}

				param := &wafv2.ListTagsForResourceInput{
					ResourceARN: v.ARN,
				}
				ipSetTags, err := client.ListTagsForResource(ctx, param)
				if err != nil {
					return nil, err
				}

				resource := Resource{
					Region: describeCtx.KaytuRegion,
					ARN:    *v.ARN,
					Name:   *v.Name,
					Description: model.WAFv2IPSetDescription{
						IPSetSummary: v,
						Scope:        scope,
						IPSet:        op.IPSet,
						Tags:         ipSetTags.TagInfoForResource.TagList,
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
			return output.NextMarker, nil
		})
		if err != nil {
			return nil, err
		}
	}

	return values, nil
}

func WAFv2LoggingConfiguration(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := wafv2.NewFromConfig(cfg)

	scopes := []types.Scope{
		types.ScopeRegional,
	}
	if strings.EqualFold(cfg.Region, "us-east-1") {
		scopes = append(scopes, types.ScopeCloudfront)
	}

	var values []Resource
	for _, scope := range scopes {
		err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
			output, err := client.ListLoggingConfigurations(ctx, &wafv2.ListLoggingConfigurationsInput{
				Scope:      scope,
				NextMarker: prevToken,
			})
			if err != nil {
				return nil, err
			}

			for _, v := range output.LoggingConfigurations {
				resource := Resource{
					Region:      describeCtx.KaytuRegion,
					ARN:         *v.ResourceArn, // TODO: might not be the actual ARN
					Name:        *v.ResourceArn,
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
			return output.NextMarker, nil
		})
		if err != nil {
			return nil, err
		}
	}

	return values, nil
}

func WAFv2RegexPatternSet(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := wafv2.NewFromConfig(cfg)

	scopes := []types.Scope{
		types.ScopeRegional,
	}
	if strings.EqualFold(cfg.Region, "us-east-1") {
		scopes = append(scopes, types.ScopeCloudfront)
	}

	var values []Resource
	for _, scope := range scopes {
		err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
			output, err := client.ListRegexPatternSets(ctx, &wafv2.ListRegexPatternSetsInput{
				Scope:      scope,
				NextMarker: prevToken,
			})
			if err != nil {
				if isErr(err, "WAFNonexistentItemException") {
					return nil, nil
				}
				return nil, err
			}

			for _, v := range output.RegexPatternSets {
				loc := strings.Split(strings.Split(*v.ARN, ":")[5], "/")[0]
				var scope types.Scope
				if loc == "regional" {
					scope = types.ScopeRegional
				} else {
					scope = types.ScopeCloudfront
				}

				op, err := client.GetRegexPatternSet(ctx, &wafv2.GetRegexPatternSetInput{
					Id:    v.Id,
					Name:  v.Name,
					Scope: scope,
				})
				if err != nil {
					if isErr(err, "WAFNonexistentItemException") {
						continue
					}
					return nil, err
				}

				regexPatternSetTags, err := client.ListTagsForResource(ctx, &wafv2.ListTagsForResourceInput{
					ResourceARN: v.ARN,
				})
				if err != nil {
					if isErr(err, "WAFNonexistentItemException") {
						regexPatternSetTags = &wafv2.ListTagsForResourceOutput{}
					} else {
						return nil, err
					}
				}

				resource := Resource{
					Region: describeCtx.KaytuRegion,
					ARN:    *v.ARN,
					Name:   *v.Name,
					Description: model.WAFv2RegexPatternSetDescription{
						RegexPatternSetSummary: v,
						Scope:                  types.Scope(scope),
						RegexPatternSet:        op.RegexPatternSet,
						Tags:                   regexPatternSetTags,
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
			return output.NextMarker, nil
		})
		if err != nil {
			return nil, err
		}
	}

	return values, nil
}

func WAFv2RuleGroup(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := wafv2.NewFromConfig(cfg)

	scopes := []types.Scope{
		types.ScopeRegional,
	}
	if strings.EqualFold(cfg.Region, "us-east-1") {
		scopes = append(scopes, types.ScopeCloudfront)
	}

	var values []Resource
	for _, scope := range scopes {
		err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
			output, err := client.ListRuleGroups(ctx, &wafv2.ListRuleGroupsInput{
				Scope:      scope,
				NextMarker: prevToken,
			})
			if err != nil {
				return nil, err
			}

			for _, v := range output.RuleGroups {
				params := &wafv2.GetRuleGroupInput{
					Id:    v.Id,
					Name:  v.Name,
					Scope: scope,
				}

				op, err := client.GetRuleGroup(ctx, params)
				if err != nil {
					return nil, err
				}

				param := &wafv2.ListTagsForResourceInput{
					ResourceARN: v.ARN,
				}

				ruleGroupTags, err := client.ListTagsForResource(ctx, param)

				resource := Resource{
					Region: describeCtx.KaytuRegion,
					ARN:    *v.ARN,
					Name:   *v.Name,
					Description: model.WAFv2RuleGroupDescription{
						RuleGroupSummary: v,
						RuleGroup:        op.RuleGroup,
						Tags:             ruleGroupTags,
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
			return output.NextMarker, nil
		})
		if err != nil {
			return nil, err
		}
	}

	return values, nil
}

func WAFv2WebACL(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	scopes := []types.Scope{
		types.ScopeRegional,
	}
	if strings.EqualFold(cfg.Region, "us-east-1") {
		scopes = append(scopes, types.ScopeCloudfront)
	}

	var values []Resource
	for _, scope := range scopes {
		acls, err := listWAFv2WebACLs(ctx, cfg, scope)
		if err != nil {
			return nil, err
		}

		for _, v := range acls {
			resource, err := wAFv2WebACLHandle(ctx, cfg, v, scope)
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

func wAFv2WebACLHandle(ctx context.Context, cfg aws.Config, v types.WebACLSummary, scope types.Scope) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := wafv2.NewFromConfig(cfg)

	out, err := client.GetWebACL(ctx, &wafv2.GetWebACLInput{
		Id:    v.Id,
		Name:  v.Name,
		Scope: scope,
	})
	if err != nil {
		return Resource{}, err
	}

	logC, err := client.GetLoggingConfiguration(ctx, &wafv2.GetLoggingConfigurationInput{
		ResourceArn: out.WebACL.ARN,
	})
	if err != nil {
		if isErr(err, "WAFNonexistentItemException") {
			logC = &wafv2.GetLoggingConfigurationOutput{}
			err = nil
		}
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "WAFNonexistentItemException" {
				logC = &wafv2.GetLoggingConfigurationOutput{}
				err = nil
			}
		}

		if err != nil {
			return Resource{}, err
		}
	}

	tags, err := client.ListTagsForResource(ctx, &wafv2.ListTagsForResourceInput{
		ResourceARN: out.WebACL.ARN,
	})
	if err != nil {
		return Resource{}, err
	}

	var associatedResources []string
	locationType := strings.Split(strings.Split(*v.ARN, ":")[5], "/")[0]
	if locationType == "global" {
		cfClient := cloudfront.NewFromConfig(cfg)

		output, err := cfClient.ListDistributionsByWebACLId(ctx, &cloudfront.ListDistributionsByWebACLIdInput{
			WebACLId: v.Id,
		})
		if err != nil {
			return Resource{}, err
		}

		if output.DistributionList != nil {
			if len(output.DistributionList.Items) > 0 {
				for _, item := range output.DistributionList.Items {
					associatedResources = append(associatedResources, *item.ARN)
				}
			}
		}
	} else {
		param := &wafv2.ListResourcesForWebACLInput{
			WebACLArn: v.ARN,
		}
		resourceTypes := []types.ResourceType{types.ResourceTypeApplicationLoadBalancer, types.ResourceTypeApiGateway, types.ResourceTypeAppsync, types.ResourceTypeCognitioUserPool}
		for _, resourceType := range resourceTypes {
			param.ResourceType = resourceType
			op, err := client.ListResourcesForWebACL(ctx, param)
			if err != nil {
				plugin.Logger(ctx).Error("aws_wafv2_web_acl.listAssociatedResourcesByResourceType", "api_error", err)
				var ae smithy.APIError
				if errors.As(err, &ae) {
					if ae.ErrorCode() == "WAFNonexistentItemException" {
						op = &wafv2.ListResourcesForWebACLOutput{}
					} else {
						return Resource{}, err
					}
				} else {
					return Resource{}, err
				}
			}
			associatedResources = append(associatedResources, op.ResourceArns...)
		}
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *v.ARN,
		Name:   *v.Name,
		Description: model.WAFv2WebACLDescription{
			WebACL:               out.WebACL,
			Scope:                scope,
			LoggingConfiguration: logC.LoggingConfiguration,
			TagInfoForResource:   tags.TagInfoForResource,
			LockToken:            v.LockToken,
			AssociatedResources:  associatedResources,
		},
	}
	return resource, nil
}

func GetWAFv2WebACL(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	scopeInput := fields["scope"]
	scopes := []types.Scope{
		types.ScopeRegional,
	}
	if strings.EqualFold(cfg.Region, "us-east-1") {
		scopes = append(scopes, types.ScopeCloudfront)
	}

	var values []Resource
	for _, scope := range scopes {
		if string(scope) != scopeInput {
			continue
		}
		acls, err := listWAFv2WebACLs(ctx, cfg, scope)
		if err != nil {
			return nil, err
		}

		for _, v := range acls {
			resource, err := wAFv2WebACLHandle(ctx, cfg, v, scope)
			if err != nil {
				return nil, err
			}
			values = append(values, resource)
		}
	}
	return values, nil
}

func listWAFv2WebACLs(ctx context.Context, cfg aws.Config, scope types.Scope) ([]types.WebACLSummary, error) {
	client := wafv2.NewFromConfig(cfg)

	var acls []types.WebACLSummary
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		output, err := client.ListWebACLs(ctx, &wafv2.ListWebACLsInput{
			Scope:      scope,
			NextMarker: prevToken,
		})
		if err != nil {
			return nil, err
		}

		acls = append(acls, output.WebACLs...)
		return output.NextMarker, nil
	})
	if err != nil {
		return nil, err
	}

	return acls, nil
}

// Returns ResourceArns that have a WebAcl Associated
func WAFv2WebACLAssociation(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	var values []Resource

	regionalACls, err := listWAFv2WebACLs(ctx, cfg, types.ScopeRegional)
	if err != nil {
		return nil, err
	}

	client := wafv2.NewFromConfig(cfg)
	for _, acl := range regionalACls {
		output, err := client.ListResourcesForWebACL(ctx, &wafv2.ListResourcesForWebACLInput{
			WebACLArn: acl.ARN,
		})
		if err != nil {
			return nil, err
		}

		resource := Resource{
			Region: describeCtx.KaytuRegion,
			ID:     *acl.Id, // Unique per WebACL
			Name:   *acl.Name,
			Description: map[string]interface{}{
				"WebACLArn":    *acl.ARN,
				"ResourceArns": output.ResourceArns,
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

	if strings.EqualFold(cfg.Region, "us-east-1") {
		cloudFrontAcls, err := listWAFv2WebACLs(ctx, cfg, types.ScopeCloudfront)
		if err != nil {
			return nil, err
		}

		cfClient := cloudfront.NewFromConfig(cfg)
		for _, acl := range cloudFrontAcls {
			err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
				output, err := cfClient.ListDistributionsByWebACLId(ctx, &cloudfront.ListDistributionsByWebACLIdInput{
					WebACLId: acl.Id,
					Marker:   prevToken,
				})
				if err != nil {
					return nil, err
				}

				resource := Resource{
					Region: describeCtx.KaytuRegion,
					ID:     *acl.Id, // Unique per WebACL
					Name:   *acl.Name,
					Description: map[string]interface{}{
						"WebACLArn":     *acl.ARN,
						"Distributions": output.DistributionList.Items,
					},
				}
				if stream != nil {
					if err := (*stream)(resource); err != nil {
						return nil, err
					}
				} else {
					values = append(values, resource)
				}

				return output.DistributionList.NextMarker, nil
			})
			if err != nil {
				return nil, err
			}
		}
	}

	return values, nil
}

func WAFRegionalByteMatchSet(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := wafregional.NewFromConfig(cfg)

	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		output, err := client.ListByteMatchSets(ctx, &wafregional.ListByteMatchSetsInput{
			NextMarker: prevToken,
		})
		if err != nil {
			return nil, err
		}

		for _, v := range output.ByteMatchSets {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ID:          *v.ByteMatchSetId,
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
		return output.NextMarker, nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}

func WAFRegionalGeoMatchSet(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := wafregional.NewFromConfig(cfg)

	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		output, err := client.ListGeoMatchSets(ctx, &wafregional.ListGeoMatchSetsInput{
			NextMarker: prevToken,
		})
		if err != nil {
			return nil, err
		}

		for _, v := range output.GeoMatchSets {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ID:          *v.GeoMatchSetId,
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
		return output.NextMarker, nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}

func WAFRegionalIPSet(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := wafregional.NewFromConfig(cfg)

	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		output, err := client.ListIPSets(ctx, &wafregional.ListIPSetsInput{
			NextMarker: prevToken,
		})
		if err != nil {
			return nil, err
		}

		for _, v := range output.IPSets {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ID:          *v.IPSetId,
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
		return output.NextMarker, nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}

func WAFRateBasedRule(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := waf.NewFromConfig(cfg)

	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		output, err := client.ListRateBasedRules(ctx, &waf.ListRateBasedRulesInput{
			NextMarker: prevToken,
		})
		if err != nil {
			return nil, err
		}

		for _, v := range output.Rules {
			arn := "arn:" + describeCtx.Partition + ":waf::" + describeCtx.AccountID + ":ratebasedrule" + "/" + *v.RuleId

			data, err := client.GetRateBasedRule(ctx, &waf.GetRateBasedRuleInput{
				RuleId: v.RuleId,
			})
			if err != nil {
				return nil, err
			}

			op, err := client.ListTagsForResource(ctx, &waf.ListTagsForResourceInput{
				ResourceARN: &arn,
			})
			if err != nil {
				return nil, err
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    arn,
				Name:   *v.Name,
				Description: model.WAFRateBasedRuleDescription{
					ARN:         arn,
					RuleSummary: v,
					Rule:        data.Rule,
					Tags:        op.TagInfoForResource,
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
		return output.NextMarker, nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}

func WAFRegionalRegexPatternSet(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := wafregional.NewFromConfig(cfg)

	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		output, err := client.ListRegexPatternSets(ctx, &wafregional.ListRegexPatternSetsInput{
			NextMarker: prevToken,
		})
		if err != nil {
			return nil, err
		}

		for _, v := range output.RegexPatternSets {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ID:          *v.RegexPatternSetId,
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
		return output.NextMarker, nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}

func WAFRegionalRule(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := wafregional.NewFromConfig(cfg)

	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		output, err := client.ListRules(ctx, &wafregional.ListRulesInput{
			NextMarker: prevToken,
		})
		if err != nil {
			return nil, err
		}

		for _, v := range output.Rules {
			resource, err := wAFRegionalRuleHandle(ctx, cfg, v)
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

func wAFRegionalRuleHandle(ctx context.Context, cfg aws.Config, v regionaltypes.RuleSummary) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := wafregional.NewFromConfig(cfg)
	arn := fmt.Sprintf("arn:%s:waf-regional:%s:%s:rule/%s", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID, *v.RuleId)
	tags, err := client.ListTagsForResource(ctx, &wafregional.ListTagsForResourceInput{
		ResourceARN: &arn,
	})
	if err != nil {
		return Resource{}, err
	}
	rule, err := client.GetRule(ctx, &wafregional.GetRuleInput{
		RuleId: v.RuleId,
	})
	if err != nil {
		return Resource{}, err
	}
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		ID:     *v.RuleId,
		Name:   *v.Name,
		Description: model.WAFRegionalRuleDescription{
			Rule: *rule.Rule,
			Tags: tags.TagInfoForResource.TagList,
		},
	}
	return resource, nil
}

func GetWAFRegionalRule(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	roleId := fields["id"]
	client := wafregional.NewFromConfig(cfg)
	var values []Resource
	role, err := client.ListRules(ctx, &wafregional.ListRulesInput{})
	if err != nil {
		return nil, err
	}
	for _, v := range role.Rules {
		if v.RuleId != &roleId {
			continue
		}
		resource, err := wAFRegionalRuleHandle(ctx, cfg, v)
		if err != nil {
			return nil, err
		}
		values = append(values, resource)
	}
	return values, nil
}

func WAFRegionalRuleGroup(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)

	client := wafregional.NewFromConfig(cfg)

	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		output, err := client.ListRuleGroups(ctx, &wafregional.ListRuleGroupsInput{
			NextMarker: prevToken,
		})
		if err != nil {
			return nil, err
		}

		for _, v := range output.RuleGroups {
			rule, err := client.GetRuleGroup(ctx, &wafregional.GetRuleGroupInput{
				RuleGroupId: v.RuleGroupId,
			})
			if err != nil {
				return nil, err
			}

			arn := fmt.Sprintf("arn:%s:waf::%s:rulegroup/%s", describeCtx.Partition, describeCtx.AccountID, *v.RuleGroupId)

			ac, err := client.ListActivatedRulesInRuleGroup(ctx, &wafregional.ListActivatedRulesInRuleGroupInput{
				RuleGroupId: v.RuleGroupId,
			})
			if err != nil {
				return nil, err
			}

			tags, err := client.ListTagsForResource(ctx, &wafregional.ListTagsForResourceInput{
				ResourceARN: &arn,
			})
			if err != nil {
				return nil, err
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    arn,
				Name:   *rule.RuleGroup.Name,
				Description: model.WAFRegionalRuleGroupDescription{
					ARN:              arn,
					RuleGroupSummary: v,
					RuleGroup:        rule.RuleGroup,
					ActivatedRules:   ac.ActivatedRules,
					Tags:             tags.TagInfoForResource.TagList,
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
		return output.NextMarker, nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}

func WAFRegionalSizeConstraintSet(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := wafregional.NewFromConfig(cfg)

	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		output, err := client.ListSizeConstraintSets(ctx, &wafregional.ListSizeConstraintSetsInput{
			NextMarker: prevToken,
		})
		if err != nil {
			return nil, err
		}

		for _, v := range output.SizeConstraintSets {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ID:          *v.SizeConstraintSetId,
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
		return output.NextMarker, nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}

func WAFRegionalSqlInjectionMatchSet(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := wafregional.NewFromConfig(cfg)

	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		output, err := client.ListSqlInjectionMatchSets(ctx, &wafregional.ListSqlInjectionMatchSetsInput{
			NextMarker: prevToken,
		})
		if err != nil {
			return nil, err
		}

		for _, v := range output.SqlInjectionMatchSets {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ID:          *v.SqlInjectionMatchSetId,
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
		return output.NextMarker, nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}

func WAFRegionalWebACL(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := wafregional.NewFromConfig(cfg)

	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		output, err := client.ListWebACLs(ctx, &wafregional.ListWebACLsInput{
			NextMarker: prevToken,
		})
		if err != nil {
			return nil, err
		}

		for _, v := range output.WebACLs {
			webAcl, err := client.GetWebACL(ctx, &wafregional.GetWebACLInput{
				WebACLId: v.WebACLId,
			})
			if err != nil {
				return nil, err
			}

			op2, err := client.GetLoggingConfiguration(ctx, &wafregional.GetLoggingConfigurationInput{
				ResourceArn: webAcl.WebACL.WebACLArn,
			})
			if err != nil {
				var ae smithy.APIError
				if errors.As(err, &ae) {
					if ae.ErrorCode() == "WAFNonexistentItemException" {
						op2 = &wafregional.GetLoggingConfigurationOutput{}
					} else {
						return nil, err
					}
				} else {
					return nil, err
				}
			}

			webAclTags, err := client.ListTagsForResource(ctx, &wafregional.ListTagsForResourceInput{
				ResourceARN: webAcl.WebACL.WebACLArn,
			})
			if err != nil {
				return nil, err
			}

			resources, err := client.ListResourcesForWebACL(ctx, &wafregional.ListResourcesForWebACLInput{
				WebACLId: webAcl.WebACL.WebACLId,
			})
			if err != nil {
				return nil, err
			}
			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ID:     *v.WebACLId,
				Name:   *v.Name,
				Description: model.WAFRegionalWebAclDescription{
					WebACL:               webAcl.WebACL,
					AssociatedResources:  resources.ResourceArns,
					LoggingConfiguration: op2.LoggingConfiguration,
					Tags:                 webAclTags.TagInfoForResource.TagList,
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
		return output.NextMarker, nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}

func WAFWebACL(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := waf.NewFromConfig(cfg)

	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		output, err := client.ListWebACLs(ctx, &waf.ListWebACLsInput{
			NextMarker: prevToken,
		})
		if err != nil {
			return nil, err
		}

		for _, v := range output.WebACLs {
			op, err := client.GetWebACL(ctx, &waf.GetWebACLInput{
				WebACLId: v.WebACLId,
			})
			if err != nil {
				return nil, err
			}

			op2, err := client.GetLoggingConfiguration(ctx, &waf.GetLoggingConfigurationInput{
				ResourceArn: op.WebACL.WebACLArn,
			})
			if err != nil {
				var ae smithy.APIError
				if errors.As(err, &ae) {
					if ae.ErrorCode() == "WAFNonexistentItemException" {
						op2 = &waf.GetLoggingConfigurationOutput{}
					} else {
						return nil, err
					}
				} else {
					return nil, err
				}
			}

			webAclTags, err := client.ListTagsForResource(ctx, &waf.ListTagsForResourceInput{
				ResourceARN: op.WebACL.WebACLArn,
			})
			if err != nil {
				return nil, err
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ID:     *v.WebACLId,
				Name:   *v.Name,
				Description: model.WAFWebAclDescription{
					WebACLSummary:        v,
					WebACL:               op.WebACL,
					LoggingConfiguration: op2.LoggingConfiguration,
					Tags:                 webAclTags.TagInfoForResource,
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
		return output.NextMarker, nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}

func WAFRegionalXssMatchSet(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := wafregional.NewFromConfig(cfg)

	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		output, err := client.ListXssMatchSets(ctx, &wafregional.ListXssMatchSetsInput{
			NextMarker: prevToken,
		})
		if err != nil {
			return nil, err
		}

		for _, v := range output.XssMatchSets {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ID:          *v.XssMatchSetId,
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
		return output.NextMarker, nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}

func WAFRule(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := waf.NewFromConfig(cfg)

	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		output, err := client.ListRules(ctx, &waf.ListRulesInput{
			NextMarker: prevToken,
		})
		if err != nil {
			if !isErr(err, "WAFNonexistentItemException") {
				return nil, err
			}
			return nil, nil
		}

		for _, v := range output.Rules {
			resource, err := wAFRuleHandle(ctx, cfg, *v.RuleId)
			if err != nil {
				return nil, err
			}
			emptyResource := Resource{}
			if err == nil && resource == emptyResource {
				continue
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

func wAFRuleHandle(ctx context.Context, cfg aws.Config, roleId string) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := waf.NewFromConfig(cfg)

	rule, err := client.GetRule(ctx, &waf.GetRuleInput{
		RuleId: &roleId,
	})
	if err != nil {
		if !isErr(err, "WAFNonexistentItemException") {
			return Resource{}, err
		}
		return Resource{}, nil
	}

	arn := fmt.Sprintf("arn:%s:waf::%s:rule/%s", describeCtx.Partition, describeCtx.AccountID, roleId)

	tags, err := client.ListTagsForResource(ctx, &waf.ListTagsForResourceInput{
		ResourceARN: &arn,
	})
	if err != nil {
		if !isErr(err, "WAFNonexistentItemException") {
			return Resource{}, err
		}
		tags = &waf.ListTagsForResourceOutput{
			TagInfoForResource: &waftypes.TagInfoForResource{},
		}
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *rule.Rule.Name,
		Description: model.WAFRuleDescription{
			Rule: *rule.Rule,
			Tags: tags.TagInfoForResource.TagList,
		},
	}
	return resource, nil
}

func GetWAFRule(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	ruleId := fields["id"]
	var values []Resource

	resource, err := wAFRuleHandle(ctx, cfg, ruleId)
	if err != nil {
		return nil, err
	}
	emptyResource := Resource{}
	if err == nil && resource == emptyResource {
		return nil, nil
	}
	values = append(values, resource)
	return values, nil
}

func WAFRuleGroup(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)

	client := waf.NewFromConfig(cfg)

	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		output, err := client.ListRuleGroups(ctx, &waf.ListRuleGroupsInput{
			NextMarker: prevToken,
		})
		if err != nil {
			return nil, err
		}

		for _, v := range output.RuleGroups {
			if v.RuleGroupId == nil {
				continue
			}
			rule, err := client.GetRuleGroup(ctx, &waf.GetRuleGroupInput{
				RuleGroupId: v.RuleGroupId,
			})
			if err != nil {
				return nil, err
			}

			arn := fmt.Sprintf("arn:%s:waf::%s:rulegroup/%s", describeCtx.Partition, describeCtx.AccountID, *v.RuleGroupId)

			ac, err := client.ListActivatedRulesInRuleGroup(ctx, &waf.ListActivatedRulesInRuleGroupInput{
				RuleGroupId: v.RuleGroupId,
			})
			if err != nil {
				return nil, err
			}

			tags, err := client.ListTagsForResource(ctx, &waf.ListTagsForResourceInput{
				ResourceARN: &arn,
			})
			if err != nil {
				return nil, err
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    arn,
				Description: model.WAFRuleGroupDescription{
					ARN:              arn,
					RuleGroupSummary: v,
					RuleGroup:        rule,
					ActivatedRules:   ac,
					Tags:             tags.TagInfoForResource.TagList,
				},
			}
			if rule.RuleGroup.Name != nil {
				resource.Name = *rule.RuleGroup.Name
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
