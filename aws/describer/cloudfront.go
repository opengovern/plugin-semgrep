package describer

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func CloudFrontDistribution(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := cloudfront.NewFromConfig(cfg)
	paginator := cloudfront.NewListDistributionsPaginator(client, &cloudfront.ListDistributionsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range page.DistributionList.Items {
			tags, err := client.ListTagsForResource(ctx, &cloudfront.ListTagsForResourceInput{
				Resource: item.ARN,
			})
			if err != nil {
				return nil, err
			}

			distribution, err := client.GetDistribution(ctx, &cloudfront.GetDistributionInput{
				Id: item.Id,
			})
			if err != nil {
				return nil, err
			}

			resource := cloudFrontDistributionHandle(ctx, tags, distribution, item.ARN, item.Id)
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
func cloudFrontDistributionHandle(ctx context.Context, tags *cloudfront.ListTagsForResourceOutput, distribution *cloudfront.GetDistributionOutput, ARN *string, Id *string) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *ARN,
		Name:   *Id,
		Description: model.CloudFrontDistributionDescription{
			Distribution: distribution.Distribution,
			ETag:         distribution.ETag,
			Tags:         tags.Tags.Items,
		},
	}
	return resource
}
func GetCloudFrontDistribution(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	client := cloudfront.NewFromConfig(cfg)
	id := fields["id"]
	var values []Resource

	DistributionData, err := client.GetDistribution(ctx, &cloudfront.GetDistributionInput{
		Id: &id,
	})
	if err != nil {
		if isErr(err, "GetDistributionNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	tags, err := client.ListTagsForResource(ctx, &cloudfront.ListTagsForResourceInput{
		Resource: DistributionData.Distribution.ARN,
	})
	if err != nil {
		if isErr(err, "ListTagsForResourceNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	distribution, err := client.GetDistribution(ctx, &cloudfront.GetDistributionInput{
		Id: DistributionData.Distribution.Id,
	})
	if err != nil {
		if isErr(err, "ListTagsForResourceNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	values = append(values, cloudFrontDistributionHandle(ctx, tags, distribution, DistributionData.Distribution.ARN, DistributionData.Distribution.Id))
	return values, nil
}

func CloudFrontStreamingDistribution(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := cloudfront.NewFromConfig(cfg)
	paginator := cloudfront.NewListStreamingDistributionsPaginator(client, &cloudfront.ListStreamingDistributionsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range page.StreamingDistributionList.Items {
			tags, err := client.ListTagsForResource(ctx, &cloudfront.ListTagsForResourceInput{
				Resource: item.ARN,
			})
			if err != nil {
				return nil, err
			}

			distribution, err := client.GetStreamingDistribution(ctx, &cloudfront.GetStreamingDistributionInput{
				Id: item.Id,
			})
			if err != nil {
				return nil, err
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    *item.ARN,
				Name:   *item.Id,
				Description: model.CloudFrontStreamingDistributionDescription{
					StreamingDistribution: distribution.StreamingDistribution,
					ETag:                  distribution.ETag,
					Tags:                  tags.Tags.Items,
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

func CloudFrontOriginAccessControl(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := cloudfront.NewFromConfig(cfg)

	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		output, err := client.ListOriginAccessControls(ctx, &cloudfront.ListOriginAccessControlsInput{
			Marker:   prevToken,
			MaxItems: aws.Int32(100),
		})
		if err != nil {
			return nil, err
		}
		for _, v := range output.OriginAccessControlList.Items {
			resource := cloudFrontOriginAccessControlHandle(ctx, cfg, v)
			if stream != nil {
				if err := (*stream)(resource); err != nil {
					return nil, err
				}
			} else {
				values = append(values, resource)
			}

		}
		return output.OriginAccessControlList.NextMarker, nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}
func cloudFrontOriginAccessControlHandle(ctx context.Context, cfg aws.Config, v types.OriginAccessControlSummary) Resource {
	describeCtx := GetDescribeContext(ctx)
	client := cloudfront.NewFromConfig(cfg)

	var tags []types.Tag
	arn := fmt.Sprintf("arn:%s:cloudfront::%s:origin-access-control/%s", describeCtx.Partition, describeCtx.AccountID, *v.Id)

	tagsOutput, err := client.ListTagsForResource(ctx, &cloudfront.ListTagsForResourceInput{
		Resource: &arn,
	})
	if err == nil {
		tags = tagsOutput.Tags.Items
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *v.Id,
		Description: model.CloudFrontOriginAccessControlDescription{
			OriginAccessControl: v,
			Tags:                tags,
		},
	}
	return resource
}
func GetCloudFrontOriginAccessControl(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	client := cloudfront.NewFromConfig(cfg)
	id := fields["id"]

	var values []Resource
	out, err := client.GetOriginAccessControl(ctx, &cloudfront.GetOriginAccessControlInput{
		Id: &id,
	})
	if err != nil {
		if isErr(err, "GetOriginAccessControlNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	originAccessControl := types.OriginAccessControlSummary{
		Id:                            out.OriginAccessControl.Id,
		Name:                          out.OriginAccessControl.OriginAccessControlConfig.Name,
		Description:                   out.OriginAccessControl.OriginAccessControlConfig.Description,
		OriginAccessControlOriginType: out.OriginAccessControl.OriginAccessControlConfig.OriginAccessControlOriginType,
		SigningBehavior:               out.OriginAccessControl.OriginAccessControlConfig.SigningBehavior,
		SigningProtocol:               out.OriginAccessControl.OriginAccessControlConfig.SigningProtocol,
	}

	values = append(values, cloudFrontOriginAccessControlHandle(ctx, cfg, originAccessControl))
	return values, nil
}

func CloudFrontCachePolicy(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := cloudfront.NewFromConfig(cfg)

	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		output, err := client.ListCachePolicies(ctx, &cloudfront.ListCachePoliciesInput{
			Marker:   prevToken,
			MaxItems: aws.Int32(1000),
		})
		if err != nil {
			return nil, err
		}

		for _, v := range output.CachePolicyList.Items {

			cachePolicy, err := client.GetCachePolicy(ctx, &cloudfront.GetCachePolicyInput{
				Id: v.CachePolicy.Id,
			})
			if err != nil {
				return nil, err
			}

			resource := cloudFrontCachePolicyHandle(ctx, cachePolicy, *v.CachePolicy.Id)
			if stream != nil {
				if err := (*stream)(resource); err != nil {
					return nil, err
				}
			} else {
				values = append(values, resource)
			}

		}
		return output.CachePolicyList.NextMarker, nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}
func cloudFrontCachePolicyHandle(ctx context.Context, cachePolicy *cloudfront.GetCachePolicyOutput, id string) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := fmt.Sprintf("arn:%s:cloudfront::%s:cache-policy/%s", describeCtx.Partition, describeCtx.AccountID, id)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		ID:     id,
		Description: model.CloudFrontCachePolicyDescription{
			CachePolicy: *cachePolicy,
		},
	}
	return resource
}
func GetCloudFrontCachePolicy(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	client := cloudfront.NewFromConfig(cfg)
	id := fields["id"]
	var values []Resource
	cachePolicy, err := client.GetCachePolicy(ctx, &cloudfront.GetCachePolicyInput{
		Id: &id,
	})
	if err != nil {
		if isErr(err, "GetCachePolicyNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	values = append(values, cloudFrontCachePolicyHandle(ctx, cachePolicy, id))
	return values, nil
}

func CloudFrontFunction(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := cloudfront.NewFromConfig(cfg)

	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		output, err := client.ListFunctions(ctx, &cloudfront.ListFunctionsInput{
			Marker:   prevToken,
			MaxItems: aws.Int32(1000),
		})
		if err != nil {
			return nil, err
		}

		for _, v := range output.FunctionList.Items {
			function, err := client.DescribeFunction(ctx, &cloudfront.DescribeFunctionInput{
				Name:  v.Name,
				Stage: v.FunctionMetadata.Stage,
			})
			if err != nil {
				return nil, err
			}

			resource := cloudFrontFunctionHandle(ctx, function)
			if stream != nil {
				if err := (*stream)(resource); err != nil {
					return nil, err
				}
			} else {
				values = append(values, resource)
			}
		}
		return output.FunctionList.NextMarker, nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}
func cloudFrontFunctionHandle(ctx context.Context, function *cloudfront.DescribeFunctionOutput) Resource {
	describeCtx := GetDescribeContext(ctx)

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *function.FunctionSummary.FunctionMetadata.FunctionARN,
		Name:   *function.FunctionSummary.Name,
		Description: model.CloudFrontFunctionDescription{
			Function: *function,
		},
	}
	return resource
}
func GetCloudFrontFunction(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	client := cloudfront.NewFromConfig(cfg)
	name := fields["name"]

	var values []Resource
	function, err := client.DescribeFunction(ctx, &cloudfront.DescribeFunctionInput{
		Name: &name,
	})
	if err != nil {
		if isErr(err, "DescribeFunctionNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	values = append(values, cloudFrontFunctionHandle(ctx, function))
	return values, nil
}

func CloudFrontOriginAccessIdentity(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := cloudfront.NewFromConfig(cfg)
	var values []Resource
	paginator := cloudfront.NewListCloudFrontOriginAccessIdentitiesPaginator(client, &cloudfront.ListCloudFrontOriginAccessIdentitiesInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range page.CloudFrontOriginAccessIdentityList.Items {
			originAccessIdentity, err := client.GetCloudFrontOriginAccessIdentity(ctx, &cloudfront.GetCloudFrontOriginAccessIdentityInput{
				Id: item.Id,
			})
			if err != nil {
				return nil, err
			}

			resource := cloudFrontOriginAccessIdentityHandle(ctx, originAccessIdentity, *item.Id)
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
func cloudFrontOriginAccessIdentityHandle(ctx context.Context, originAccessIdentity *cloudfront.GetCloudFrontOriginAccessIdentityOutput, id string) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := fmt.Sprintf("arn:%s:cloudfront::%s:origin-access-identity/%s", describeCtx.Partition, describeCtx.AccountID, id)

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   id,
		Description: model.CloudFrontOriginAccessIdentityDescription{
			OriginAccessIdentity: *originAccessIdentity,
		},
	}
	return resource
}
func GetCloudFrontOriginAccessIdentity(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	client := cloudfront.NewFromConfig(cfg)
	id := fields["id"]
	var values []Resource

	originAccessIdentity, err := client.GetCloudFrontOriginAccessIdentity(ctx, &cloudfront.GetCloudFrontOriginAccessIdentityInput{
		Id: &id,
	})
	if err != nil {
		return nil, err
	}

	values = append(values, cloudFrontOriginAccessIdentityHandle(ctx, originAccessIdentity, id))

	return values, nil
}

func CloudFrontOriginRequestPolicy(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := cloudfront.NewFromConfig(cfg)
	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		output, err := client.ListOriginRequestPolicies(ctx, &cloudfront.ListOriginRequestPoliciesInput{
			Marker:   prevToken,
			MaxItems: aws.Int32(1000),
		})
		if err != nil {
			return nil, err
		}

		for _, v := range output.OriginRequestPolicyList.Items {

			policy, err := client.GetOriginRequestPolicy(ctx, &cloudfront.GetOriginRequestPolicyInput{
				Id: v.OriginRequestPolicy.Id,
			})
			if err != nil {
				return nil, err
			}

			resource := cloudFrontOriginRequestPolicyHandle(ctx, policy, *v.OriginRequestPolicy.Id)

			if stream != nil {
				if err := (*stream)(resource); err != nil {
					return nil, err
				}
			} else {
				values = append(values, resource)
			}

		}
		return output.OriginRequestPolicyList.NextMarker, nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}
func cloudFrontOriginRequestPolicyHandle(ctx context.Context, policy *cloudfront.GetOriginRequestPolicyOutput, id string) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := fmt.Sprintf("arn:%s:cloudfront::%s:origin-request-policy/%s", describeCtx.Partition, describeCtx.AccountID, &id)

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		ID:     *policy.OriginRequestPolicy.Id,
		Description: model.CloudFrontOriginRequestPolicyDescription{
			OriginRequestPolicy: *policy,
		},
	}
	return resource
}
func GetCloudFrontOriginRequestPolicy(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	id := fields["id"]
	var values []Resource
	client := cloudfront.NewFromConfig(cfg)

	policy, err := client.GetOriginRequestPolicy(ctx, &cloudfront.GetOriginRequestPolicyInput{
		Id: &id,
	})
	if err != nil {
		return nil, err
	}

	resource := cloudFrontOriginRequestPolicyHandle(ctx, policy, id)
	if err != nil {
		return nil, err
	}

	values = append(values, resource)
	return values, nil
}

func CloudFrontResponseHeadersPolicy(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := cloudfront.NewFromConfig(cfg)
	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		output, err := client.ListResponseHeadersPolicies(ctx, &cloudfront.ListResponseHeadersPoliciesInput{
			Marker:   prevToken,
			MaxItems: aws.Int32(1000),
		})
		if err != nil {
			return nil, err
		}

		for _, v := range output.ResponseHeadersPolicyList.Items {
			policy, err := client.GetResponseHeadersPolicy(ctx, &cloudfront.GetResponseHeadersPolicyInput{
				Id: v.ResponseHeadersPolicy.Id,
			})
			if err != nil {
				return nil, err
			}

			resource := cloudFrontResponseHeadersPolicyHandle(ctx, policy, v.ResponseHeadersPolicy.Id)
			if stream != nil {
				if err := (*stream)(resource); err != nil {
					return nil, err
				}
			} else {
				values = append(values, resource)
			}

		}
		return output.ResponseHeadersPolicyList.NextMarker, nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}
func cloudFrontResponseHeadersPolicyHandle(ctx context.Context, policy *cloudfront.GetResponseHeadersPolicyOutput, id *string) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := fmt.Sprintf("arn:%s:cloudfront::%s:response-headers-policy/%s", describeCtx.Partition, describeCtx.AccountID, *id)

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		ID:     *policy.ResponseHeadersPolicy.Id,
		Description: model.CloudFrontResponseHeadersPolicyDescription{
			ResponseHeadersPolicy: *policy,
		},
	}
	return resource
}
func GetCloudFrontResponseHeadersPolicy(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	id := fields["id"]
	var values []Resource
	client := cloudfront.NewFromConfig(cfg)

	policy, err := client.GetResponseHeadersPolicy(ctx, &cloudfront.GetResponseHeadersPolicyInput{
		Id: &id,
	})
	if err != nil {
		return nil, err
	}

	values = append(values, cloudFrontResponseHeadersPolicyHandle(ctx, policy, &id))
	return values, nil
}
