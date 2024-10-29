package describer

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dlm"
	"github.com/aws/aws-sdk-go-v2/service/dlm/types"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func DLMLifecyclePolicy(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := dlm.NewFromConfig(cfg)

	lifecyclePolicies, err := client.GetLifecyclePolicies(ctx, &dlm.GetLifecyclePoliciesInput{})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, policySummary := range lifecyclePolicies.Policies {

		resource, err := dLMLifecyclePolicyHandle(ctx, cfg, policySummary)
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
	return values, nil
}
func dLMLifecyclePolicyHandle(ctx context.Context, cfg aws.Config, policySummary types.LifecyclePolicySummary) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := dlm.NewFromConfig(cfg)
	policy, err := client.GetLifecyclePolicy(ctx, &dlm.GetLifecyclePolicyInput{
		PolicyId: policySummary.PolicyId,
	})
	if err != nil {
		if isErr(err, "GetLifecyclePolicyNotFound") || isErr(err, "InvalidParameterValue") {
			return Resource{}, nil
		}
		return Resource{}, err
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ID:     *policy.Policy.PolicyId,
		ARN:    *policy.Policy.PolicyArn,
		Description: model.DLMLifecyclePolicyDescription{
			LifecyclePolicy: *policy.Policy,
		},
	}
	return resource, nil
}
func GetDLMLifecyclePolicy(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	policyId := fields["policyId"]
	var values []Resource
	client := dlm.NewFromConfig(cfg)

	policies, err := client.GetLifecyclePolicies(ctx, &dlm.GetLifecyclePoliciesInput{
		PolicyIds: []string{policyId},
	})
	if err != nil {
		if isErr(err, "GetLifecyclePoliciesNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	for _, policySummary := range policies.Policies {
		resource, err := dLMLifecyclePolicyHandle(ctx, cfg, policySummary)
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
