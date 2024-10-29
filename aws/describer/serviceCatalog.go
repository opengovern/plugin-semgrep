package describer

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/servicecatalog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func ServiceCatalogProduct(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := servicecatalog.NewFromConfig(cfg)

	paginator := servicecatalog.NewSearchProductsPaginator(client, &servicecatalog.SearchProductsInput{})
	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, item := range page.ProductViewSummaries {
			productAsA, err := client.DescribeProductAsAdmin(ctx, &servicecatalog.DescribeProductAsAdminInput{
				Id:   item.ProductId,
				Name: item.Name,
			})
			if err != nil {
				return nil, err
			}

			lLP, err := client.ListLaunchPaths(ctx, &servicecatalog.ListLaunchPathsInput{
				ProductId: item.ProductId,
			})
			if err != nil {
				return nil, err
			}
			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ID:     *item.Id,
				Name:   *item.Name,
				Description: model.ServiceCatalogProductDescription{
					ProductViewSummary:    item,
					Budgets:               productAsA.Budgets,
					LunchPaths:            lLP.LaunchPathSummaries,
					ProvisioningArtifacts: productAsA.ProvisioningArtifactSummaries,
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
func ServiceCatalogPortfolio(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := servicecatalog.NewFromConfig(cfg)
	paginator := servicecatalog.NewListPortfoliosPaginator(client, &servicecatalog.ListPortfoliosInput{})
	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, v := range page.PortfolioDetails {
			client.DescribePortfolio(ctx, &servicecatalog.DescribePortfolioInput{
				Id: v.Id,
			})
			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ID:     *v.Id,
				Name:   *v.ProviderName,
				ARN:    *v.ARN,
				Description: model.ServiceCatalogPortfolioDescription{
					Portfolio: v,
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
