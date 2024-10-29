package describer

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/opengovern/og-aws-describer/aws/model"
)

// OrganizationOrganization Retrieves information about the organization that the
// user's account belongs to.
func OrganizationOrganization(ctx context.Context, cfg aws.Config) (*types.Organization, error) {
	client := organizations.NewFromConfig(cfg)

	req, err := client.DescribeOrganization(ctx, &organizations.DescribeOrganizationInput{})
	if err != nil {
		return nil, err
	}

	return req.Organization, nil
}

func OrganizationsOrganization(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := organizations.NewFromConfig(cfg)

	req, err := client.DescribeOrganization(ctx, &organizations.DescribeOrganizationInput{})
	if err != nil {
		return nil, err
	}

	var values []Resource
	resource := organizationsOrganizationHandle(ctx, req)
	if stream != nil {
		if err := (*stream)(resource); err != nil {
			return nil, err
		}
	} else {
		values = append(values, resource)
	}

	return values, nil
}
func organizationsOrganizationHandle(ctx context.Context, req *organizations.DescribeOrganizationOutput) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *req.Organization.Arn,
		Name:   *req.Organization.Id,
		Description: model.OrganizationsOrganizationDescription{
			Organization: *req.Organization,
		},
	}
	return resource
}
func GetOrganizationsOrganization(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	client := organizations.NewFromConfig(cfg)
	var values []Resource
	describes, err := client.DescribeOrganization(ctx, &organizations.DescribeOrganizationInput{})
	if err != nil {
		return nil, err
	}
	resource := organizationsOrganizationHandle(ctx, describes)
	values = append(values, resource)
	return values, nil
}

// OrganizationAccount Retrieves AWS Organizations-related information about
// the specified (ID) account .
func OrganizationAccount(ctx context.Context, cfg aws.Config, id string) (*types.Account, error) {
	svc := organizations.NewFromConfig(cfg)

	req, err := svc.DescribeAccount(ctx, &organizations.DescribeAccountInput{AccountId: aws.String(id)})
	if err != nil {
		return nil, err
	}

	return req.Account, nil
}

// DescribeOrganization Retrieves information about the organization that the
// user's account belongs to.
func OrganizationAccounts(ctx context.Context, cfg aws.Config) ([]types.Account, error) {
	client := organizations.NewFromConfig(cfg)

	paginator := organizations.NewListAccountsPaginator(client, &organizations.ListAccountsInput{})

	var values []types.Account
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		values = append(values, page.Accounts...)
	}

	return values, nil
}

//func OrganizationsAccount(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
//	describeCtx := GetDescribeContext(ctx)
//	client := organizations.NewFromConfig(cfg)
//
//	paginator := organizations.NewListAccountsPaginator(client, &organizations.ListAccountsInput{})
//
//	var values []Resource
//	for paginator.HasMorePages() {
//		page, err := paginator.NextPage(ctx)
//		if err != nil {
//			if isErr(err, "AccessDeniedException") {
//				continue
//			}
//			return nil, err
//		}
//
//		for _, acc := range page.Accounts {
//			tags, err := client.ListTagsForResource(ctx, &organizations.ListTagsForResourceInput{
//				ResourceId: acc.Id,
//			})
//			if err != nil {
//				return nil, err
//			}
//
//			resource := Resource{
//				Region: describeCtx.KaytuRegion,
//				ARN:    *acc.Arn,
//				Name:   *acc.Name,
//				Description: model.OrganizationsAccountDescription{
//					Account: acc,
//					Tags:    tags.Tags,
//				},
//			}
//			if stream != nil {
//				if err := (*stream)(resource); err != nil {
//					return nil, err
//				}
//			} else {
//				values = append(values, resource)
//			}
//		}
//	}
//
//	return values, nil
//}

func OrganizationsPolicy(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	var values []Resource
	for _, pType := range []types.PolicyType{types.PolicyTypeServiceControlPolicy, types.PolicyTypeTagPolicy,
		types.PolicyTypeBackupPolicy, types.PolicyTypeAiservicesOptOutPolicy} {
		resources, err := getOrganizationsPolicyByType(ctx, cfg, pType)
		if err != nil {
			return nil, err
		}
		for _, resource := range resources {
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

func getOrganizationsPolicyByType(ctx context.Context, cfg aws.Config, policyType types.PolicyType) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := organizations.NewFromConfig(cfg)
	paginator := organizations.NewListPoliciesPaginator(client, &organizations.ListPoliciesInput{Filter: policyType})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, p := range page.Policies {
			policy, err := client.DescribePolicy(ctx, &organizations.DescribePolicyInput{
				PolicyId: p.Id,
			})
			if err != nil {
				return nil, err
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    *p.Arn,
				Name:   *p.Name,
				Description: model.OrganizationsPolicyDescription{
					Policy: *policy.Policy,
				},
			}
			values = append(values, resource)
		}
	}

	return values, nil
}

func OrganizationsRoot(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := organizations.NewFromConfig(cfg)

	paginator := organizations.NewListRootsPaginator(client, &organizations.ListRootsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, root := range output.Roots {
			resource := organizationsRootHandle(ctx, root)
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

func organizationsRootHandle(ctx context.Context, root types.Root) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *root.Arn,
		Name:   *root.Name,
		Description: model.OrganizationsRootDescription{
			Root: root,
		},
	}
	return resource
}

func OrganizationsOrganizationalUnit(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := organizations.NewFromConfig(cfg)

	paginator := organizations.NewListRootsPaginator(client, &organizations.ListRootsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, root := range output.Roots {
			resources, err := listAllNestedOUs(ctx, client, *root.Id, *root.Id)
			if err != nil {
				return nil, err
			}
			for _, r := range resources {
				if stream != nil {
					if err := (*stream)(r); err != nil {
						return nil, err
					}
				} else {
					values = append(values, r)
				}
			}
		}
	}

	return values, nil
}

func listAllNestedOUs(ctx context.Context, svc *organizations.Client, parentId string, currentPath string) ([]Resource, error) {
	params := &organizations.ListOrganizationalUnitsForParentInput{
		ParentId: aws.String(parentId),
	}
	paginator := organizations.NewListOrganizationalUnitsForParentPaginator(svc, params)
	var values []Resource
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, unit := range output.OrganizationalUnits {
			ouPath := strings.Replace(currentPath, "-", "_", -1) + "." + strings.Replace(*unit.Id, "-", "_", -1)
			resource, err := organizationsOrganizationalUnitHandle(ctx, svc, unit, parentId, ouPath)
			if err != nil {
				return nil, err
			}
			values = append(values, *resource)

			// Recursively list units for this child
			resources, err := listAllNestedOUs(ctx, svc, *unit.Id, ouPath)
			if err != nil {
				return nil, err
			}
			for _, r := range resources {
				values = append(values, r)
			}

		}
	}
	return values, nil
}

func organizationsOrganizationalUnitHandle(ctx context.Context, svc *organizations.Client, unit types.OrganizationalUnit, parentId, path string) (*Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	tagsResponse, err := svc.ListTagsForResource(ctx, &organizations.ListTagsForResourceInput{
		ResourceId: unit.Id,
	})
	if err != nil {
		return nil, err
	}
	var tags []types.Tag
	if tagsResponse != nil {
		tags = tagsResponse.Tags
	}
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *unit.Arn,
		Name:   *unit.Name,
		ID:     *unit.Id,
		Description: model.OrganizationsOrganizationalUnitDescription{
			Unit:     unit,
			Path:     path,
			ParentId: parentId,
			Tags:     tags,
		},
	}
	return &resource, nil
}

func OrganizationsPolicyTarget(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := organizations.NewFromConfig(cfg)

	// We should get the policies for different target types

	var values []Resource
	// Accounts
	paginator := organizations.NewListAccountsPaginator(client, &organizations.ListAccountsInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if isErr(err, "AccessDeniedException") {
				continue
			}
			return nil, err
		}

		for _, acc := range page.Accounts {
			policies, err := organizationsPolicyForTarget(ctx, client, *acc.Id)
			if err != nil {
				return nil, err
			}
			for _, policy := range policies {
				if stream != nil {
					if err := (*stream)(policy); err != nil {
						return nil, err
					}
				} else {
					values = append(values, policy)
				}
			}
		}
	}

	// Roots
	paginator2 := organizations.NewListRootsPaginator(client, &organizations.ListRootsInput{})

	for paginator2.HasMorePages() {
		page, err := paginator2.NextPage(ctx)
		if err != nil {
			if isErr(err, "AccessDeniedException") {
				continue
			}
			return nil, err
		}

		for _, root := range page.Roots {
			policies, err := organizationsPolicyForTarget(ctx, client, *root.Id)
			if err != nil {
				return nil, err
			}
			for _, policy := range policies {
				if stream != nil {
					if err := (*stream)(policy); err != nil {
						return nil, err
					}
				} else {
					values = append(values, policy)
				}
			}
		}
	}

	// OUs
	ous, err := OrganizationsOrganizationalUnit(ctx, cfg, nil)
	if err != nil {
		return nil, err
	}
	for _, ou := range ous {
		policies, err := organizationsPolicyForTarget(ctx, client, ou.ID)
		if err != nil {
			return nil, err
		}
		for _, policy := range policies {
			if stream != nil {
				if err := (*stream)(policy); err != nil {
					return nil, err
				}
			} else {
				values = append(values, policy)
			}
		}
	}

	return values, nil
}

func organizationsPolicyForTarget(ctx context.Context, svc *organizations.Client, targetId string) ([]Resource, error) {
	var values []Resource
	for _, pType := range []types.PolicyType{types.PolicyTypeServiceControlPolicy, types.PolicyTypeTagPolicy,
		types.PolicyTypeBackupPolicy, types.PolicyTypeAiservicesOptOutPolicy} {
		resources, err := organizationsPolicyForTargetByPolicyType(ctx, svc, targetId, pType)
		if err != nil {
			return nil, err
		}
		for _, resource := range resources {
			values = append(values, resource)
		}
	}
	return values, nil
}

func organizationsPolicyForTargetByPolicyType(ctx context.Context, svc *organizations.Client, targetId string, policyType types.PolicyType) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)

	var values []Resource

	params := &organizations.ListPoliciesForTargetInput{
		Filter:   policyType,
		TargetId: &targetId,
	}
	paginator := organizations.NewListPoliciesForTargetPaginator(svc, params)

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, policy := range output.Policies {
			op, err := svc.DescribePolicy(ctx, &organizations.DescribePolicyInput{
				PolicyId: policy.Id,
			})
			if err != nil {
				return nil, err
			}
			values = append(values, Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    *policy.Arn,
				Name:   *policy.Name,
				Description: model.OrganizationsPolicyTargetDescription{
					PolicySummary: policy,
					PolicyContent: op.Policy.Content,
					TargetId:      targetId,
				},
			})
		}
	}
	return values, nil
}
