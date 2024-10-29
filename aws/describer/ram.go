package describer

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ram"
	"github.com/aws/aws-sdk-go-v2/service/ram/types"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func RamPrincipalAssociation(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ram.NewFromConfig(cfg)

	var values []Resource
	paginator := ram.NewGetResourceShareAssociationsPaginator(client, &ram.GetResourceShareAssociationsInput{AssociationType: types.ResourceShareAssociationTypePrincipal})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, association := range page.ResourceShareAssociations {
			resource, err := ramPrincipalAssociationHandle(ctx, cfg, association, *association.ResourceShareArn)
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
func GetRamPrincipalAssociation(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	resourceShareArn := fields["ResourceShareArn"]
	client := ram.NewFromConfig(cfg)

	associations, err := client.GetResourceShareAssociations(ctx, &ram.GetResourceShareAssociationsInput{
		ResourceShareArns: []string{resourceShareArn},
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, association := range associations.ResourceShareAssociations {
		resource, err := ramPrincipalAssociationHandle(ctx, cfg, association, resourceShareArn)
		if err != nil {
			return nil, err
		}
		values = append(values, resource)
	}
	return values, nil
}
func ramPrincipalAssociationHandle(ctx context.Context, cfg aws.Config, association types.ResourceShareAssociation, resourceShareArn string) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ram.NewFromConfig(cfg)

	permissionPaginator := ram.NewListResourceSharePermissionsPaginator(client, &ram.ListResourceSharePermissionsInput{
		ResourceShareArn: &resourceShareArn,
	})

	var permissions []types.ResourceSharePermissionSummary
	for permissionPaginator.HasMorePages() {
		permissionPage, err := permissionPaginator.NextPage(ctx)
		if err != nil {
			return Resource{}, err
		}
		permissions = append(permissions, permissionPage.Permissions...)
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		Name:   *association.ResourceShareName,
		ARN:    *association.ResourceShareArn,
		Description: model.RamPrincipalAssociationDescription{
			PrincipalAssociation:    association,
			ResourceSharePermission: permissions,
		},
	}
	return resource, nil
}

func RamResourceAssociation(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ram.NewFromConfig(cfg)

	var values []Resource
	paginator := ram.NewGetResourceShareAssociationsPaginator(client, &ram.GetResourceShareAssociationsInput{AssociationType: types.ResourceShareAssociationTypeResource})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, association := range page.ResourceShareAssociations {
			resource, err := ramResourceAssociationHandle(ctx, cfg, association, *association.ResourceShareArn)
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
func ramResourceAssociationHandle(ctx context.Context, cfg aws.Config, association types.ResourceShareAssociation, resourceShareArn string) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ram.NewFromConfig(cfg)
	permissionPaginator := ram.NewListResourceSharePermissionsPaginator(client, &ram.ListResourceSharePermissionsInput{
		ResourceShareArn: &resourceShareArn,
	})
	var permissions []types.ResourceSharePermissionSummary
	for permissionPaginator.HasMorePages() {
		permissionPage, err := permissionPaginator.NextPage(ctx)
		if err != nil {
			return Resource{}, err
		}
		permissions = append(permissions, permissionPage.Permissions...)
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		Name:   *association.ResourceShareName,
		ARN:    *association.ResourceShareArn,
		Description: model.RamResourceAssociationDescription{
			ResourceAssociation:     association,
			ResourceSharePermission: permissions,
		},
	}
	return resource, nil
}
func GetRamResourceAssociation(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	resourceShareArn := fields["resourceShareArn"]
	client := ram.NewFromConfig(cfg)

	associations, err := client.GetResourceShareAssociations(ctx, &ram.GetResourceShareAssociationsInput{
		ResourceShareArns: []string{resourceShareArn},
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, association := range associations.ResourceShareAssociations {

		resource, err := ramResourceAssociationHandle(ctx, cfg, association, resourceShareArn)
		if err != nil {
			return nil, err
		}

		values = append(values, resource)
	}
	return values, nil
}
