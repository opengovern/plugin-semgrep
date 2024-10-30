package provider

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"go.uber.org/zap"
)

const (
	SecurityAuditPolicyARN = "arn:aws:iam::aws:policy/SecurityAudit"
)

func CheckAttachedPolicy(logger *zap.Logger, cfg aws.Config, roleName string, expectedPolicyARN string) (bool, error) {
	if expectedPolicyARN == "" {
		expectedPolicyARN = SecurityAuditPolicyARN
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfgClone := cfg.Copy()
	if cfgClone.Region == "" {
		cfgClone.Region = "us-east-1"
	}

	iamClient := iam.NewFromConfig(cfgClone)
	policyARNs := make([]string, 0)
	if roleName == "" {
		user, err := iamClient.GetUser(ctx, &iam.GetUserInput{})
		if err != nil {
			logger.Warn("failed to get user", zap.Error(err))
			return false, err
		}

		paginator := iam.NewListAttachedUserPoliciesPaginator(iamClient, &iam.ListAttachedUserPoliciesInput{
			UserName: user.User.UserName,
		})
		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				logger.Warn("failed to get policy page", zap.Error(err))
				return false, err
			}
			for _, policy := range page.AttachedPolicies {
				policyARNs = append(policyARNs, *policy.PolicyArn)
			}
		}

		groups, err := iamClient.ListGroupsForUser(ctx, &iam.ListGroupsForUserInput{
			UserName: user.User.UserName,
		})
		if err != nil {
			logger.Warn("failed to get groups", zap.Error(err))
			return false, err
		}
		for _, group := range groups.Groups {
			paginator := iam.NewListAttachedGroupPoliciesPaginator(iamClient, &iam.ListAttachedGroupPoliciesInput{
				GroupName: group.GroupName,
			})
			for paginator.HasMorePages() {
				page, err := paginator.NextPage(ctx)
				if err != nil {
					logger.Warn("failed to get policy page", zap.Error(err))
					return false, err
				}
				for _, policy := range page.AttachedPolicies {
					policyARNs = append(policyARNs, *policy.PolicyArn)
				}
			}
		}
	} else {
		paginator := iam.NewListAttachedRolePoliciesPaginator(iamClient, &iam.ListAttachedRolePoliciesInput{
			RoleName: &roleName,
		})
		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				logger.Warn("failed to get policy page", zap.Error(err))
				return false, err
			}
			for _, policy := range page.AttachedPolicies {
				policyARNs = append(policyARNs, *policy.PolicyArn)
			}
		}
	}

	for _, policyARN := range policyARNs {
		if policyARN == expectedPolicyARN {
			return true, nil
		}
	}

	return false, nil
}

func CheckGetUserPermission(logger *zap.Logger, cfg aws.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfgClone := cfg.Copy()
	if cfgClone.Region == "" {
		cfgClone.Region = "us-east-1"
	}

	stsClient := sts.NewFromConfig(cfgClone)
	_, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		logger.Warn("failed to get called identity", zap.Error(err))
		return err
	}

	return nil
}

func getAllRegions(ctx context.Context, cfg aws.Config, includeDisabledRegions bool) ([]types.Region, error) {
	client := ec2.NewFromConfig(cfg)
	output, err := client.DescribeRegions(ctx, &ec2.DescribeRegionsInput{
		AllRegions: &includeDisabledRegions,
	})
	if err != nil {
		return nil, err
	}

	return output.Regions, nil
}

func PartitionOf(region string) (string, bool) {
	resolver := endpoints.DefaultResolver()
	partitions := resolver.(endpoints.EnumPartitions).Partitions()

	for _, p := range partitions {
		for r := range p.Regions() {
			if r == region {
				return p.ID(), true
			}
		}
	}

	return "", false
}
