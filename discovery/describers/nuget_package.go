package describers

import (
	"context"
	"strconv"

	"github.com/google/go-github/v55/github"
	"github.com/opengovern/og-describer-template/discovery/pkg/models"
	model "github.com/opengovern/og-describer-template/discovery/provider"
)

func GetNugetPackageList(ctx context.Context, githubClient model.GitHubClient, organizationName string, stream *models.StreamSender) ([]models.Resource, error) {
	client := githubClient.RestClient
	page := 1
	var values []models.Resource
	for {
		packageType := "nuget"
		var opts = &github.PackageListOptions{
			PackageType: &packageType,
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: packagePageSize,
			},
		}
		respPackages, resp, err := client.Organizations.ListPackages(ctx, organizationName, opts)
		if err != nil {
			return nil, err
		}
		for _, packageData := range respPackages {
			value := models.Resource{
				ID:   strconv.Itoa(int(packageData.GetID())),
				Name: packageData.GetName(),
				Description: model.PackageDescription{
					ID:         strconv.Itoa(int(packageData.GetID())),
					RegistryID: packageData.Registry.GetURL(),
					Name:       packageData.GetName(),
					URL:        packageData.GetURL(),
					CreatedAt:  packageData.GetCreatedAt(),
					UpdatedAt:  packageData.GetUpdatedAt(),
				},
			}
			if stream != nil {
				if err := (*stream)(value); err != nil {
					return nil, err
				}
			} else {
				values = append(values, value)
			}
		}
		if resp.After == "" {
			break
		}
		opts.ListOptions.Page += 1
	}

	return values, nil
}

func GetNugetPackage(ctx context.Context, githubClient model.GitHubClient, organizationName string, repositoryName string, resourceID string, stream *models.StreamSender) (*models.Resource, error) {
	client := githubClient.RestClient
	packageType := "nuget"
	respPackages, _, err := client.Organizations.GetPackage(ctx, organizationName, packageType, resourceID)
	if err != nil {
		return nil, err
	}
	value := models.Resource{
		ID:   strconv.Itoa(int(respPackages.GetID())),
		Name: respPackages.GetName(),
		Description: model.PackageDescription{
			ID:         strconv.Itoa(int(respPackages.GetID())),
			RegistryID: respPackages.Registry.GetURL(),
			Name:       respPackages.GetName(),
			URL:        respPackages.GetURL(),
			CreatedAt:  respPackages.GetCreatedAt(),
			UpdatedAt:  respPackages.GetUpdatedAt(),
		},
	}
	if stream != nil {
		if err := (*stream)(value); err != nil {
			return nil, err
		}
	}

	return &value, nil
}
