package describer

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func listEksClusters(ctx context.Context, cfg aws.Config) ([]string, error) {
	client := eks.NewFromConfig(cfg)
	paginator := eks.NewListClustersPaginator(client, &eks.ListClustersInput{})

	var values []string
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		values = append(values, page.Clusters...)
	}

	return values, nil
}

type EKSIdentityProviderConfigDescription struct {
	ConfigName             string
	ConfigType             string
	IdentityProviderConfig types.OidcIdentityProviderConfig
}

func EKSCluster(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	clusters, err := listEksClusters(ctx, cfg)
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, cluster := range clusters {
		// This prevents Implicit memory aliasing in for loop
		cluster := cluster
		resource, err := eKSClusterHandle(ctx, cfg, cluster)
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
func eKSClusterHandle(ctx context.Context, cfg aws.Config, cluster string) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := eks.NewFromConfig(cfg)

	output, err := client.DescribeCluster(ctx, &eks.DescribeClusterInput{Name: aws.String(cluster)})
	if err != nil {
		if isErr(err, "DescribeClusterNotFound") || isErr(err, "InvalidParameterValue") {
			return Resource{}, nil
		}
		return Resource{}, err
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *output.Cluster.Arn,
		Name:   *output.Cluster.Name,
		Description: model.EKSClusterDescription{
			Cluster: *output.Cluster,
		},
	}
	return resource, nil
}
func GetEKSCluster(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	clusterName := fields["name"]
	clusters, err := listEksClusters(ctx, cfg)
	if err != nil {
		return nil, err
	}
	var values []Resource
	for _, cluster := range clusters {
		cluster := cluster
		if !strings.EqualFold(*aws.String(cluster), clusterName) {
			continue
		}

		resource, err := eKSClusterHandle(ctx, cfg, cluster)
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

func EKSAddon(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	clusters, err := listEksClusters(ctx, cfg)
	if err != nil {
		return nil, err
	}

	client := eks.NewFromConfig(cfg)

	var values []Resource
	for _, cluster := range clusters {
		var addons []string

		paginator := eks.NewListAddonsPaginator(client, &eks.ListAddonsInput{ClusterName: aws.String(cluster)})
		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				return nil, err
			}

			addons = append(addons, page.Addons...)
		}

		for _, addon := range addons {
			resource, err := eKSAddonHandle(ctx, cfg, addon, cluster)
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
	}

	return values, nil
}
func eKSAddonHandle(ctx context.Context, cfg aws.Config, addon string, clusterName string) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := eks.NewFromConfig(cfg)

	output, err := client.DescribeAddon(ctx, &eks.DescribeAddonInput{
		AddonName:   aws.String(addon),
		ClusterName: aws.String(clusterName),
	})
	if err != nil {
		if isErr(err, "DescribeAddonNotFound") || isErr(err, "InvalidParameterValue") {
			return Resource{}, nil
		}
		return Resource{}, err
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *output.Addon.AddonArn,
		Name:   *output.Addon.AddonName,
		Description: model.EKSAddonDescription{
			Addon: *output.Addon,
		},
	}
	return resource, nil
}
func GetEKSAddon(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	clusterName := fields["clusterName"]
	client := eks.NewFromConfig(cfg)
	addons, err := client.ListAddons(ctx, &eks.ListAddonsInput{
		ClusterName: &clusterName,
	})
	if err != nil {
		if isErr(err, "ListAddonsNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	for _, addon := range addons.Addons {
		resource, err := eKSAddonHandle(ctx, cfg, addon, clusterName)
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

func EKSFargateProfile(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	clusters, err := listEksClusters(ctx, cfg)
	if err != nil {
		return nil, err
	}

	client := eks.NewFromConfig(cfg)

	var values []Resource
	for _, cluster := range clusters {
		var profiles []string

		paginator := eks.NewListFargateProfilesPaginator(client, &eks.ListFargateProfilesInput{ClusterName: aws.String(cluster)})
		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				return nil, err
			}
			profiles = append(profiles, page.FargateProfileNames...)
		}

		for _, profile := range profiles {
			resource, err := eKSFargateProfileHandle(ctx, cfg, profile, cluster)
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
	}

	return values, nil
}
func eKSFargateProfileHandle(ctx context.Context, cfg aws.Config, profile string, clusterName string) (Resource, error) {
	client := eks.NewFromConfig(cfg)
	describeCtx := GetDescribeContext(ctx)
	output, err := client.DescribeFargateProfile(ctx, &eks.DescribeFargateProfileInput{
		FargateProfileName: aws.String(profile),
		ClusterName:        aws.String(clusterName),
	})
	if err != nil {
		if isErr(err, "DescribeFargateProfileNotFound") || isErr(err, "InvalidParameterValue") {
			return Resource{}, nil
		}
		return Resource{}, err
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *output.FargateProfile.FargateProfileArn,
		Name:   *output.FargateProfile.FargateProfileName,
		Description: model.EKSFargateProfileDescription{
			FargateProfile: *output.FargateProfile,
		},
	}
	return resource, nil
}
func GetEKSFargateProfile(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	clusterName := fields["name"]
	client := eks.NewFromConfig(cfg)

	profiles, err := client.ListFargateProfiles(ctx, &eks.ListFargateProfilesInput{
		ClusterName: &clusterName,
	})
	if err != nil {
		if isErr(err, "ListFargateProfilesNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	for _, profile := range profiles.FargateProfileNames {

		resource, err := eKSFargateProfileHandle(ctx, cfg, profile, clusterName)
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

func EKSNodegroup(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	clusters, err := listEksClusters(ctx, cfg)
	if err != nil {
		return nil, err
	}

	client := eks.NewFromConfig(cfg)

	var values []Resource
	for _, cluster := range clusters {
		cluster := cluster
		var groups []string
		paginator := eks.NewListNodegroupsPaginator(client, &eks.ListNodegroupsInput{ClusterName: aws.String(cluster)})

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				return nil, err
			}

			groups = append(groups, page.Nodegroups...)
		}

		for _, profile := range groups {

			resource, err := eKSNodegroupHandle(ctx, cfg, profile, cluster)
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
	}

	return values, nil
}
func eKSNodegroupHandle(ctx context.Context, cfg aws.Config, profile string, clusterName string) (Resource, error) {
	client := eks.NewFromConfig(cfg)
	describeCtx := GetDescribeContext(ctx)

	output, err := client.DescribeNodegroup(ctx, &eks.DescribeNodegroupInput{
		NodegroupName: aws.String(profile),
		ClusterName:   aws.String(clusterName),
	})
	if err != nil {
		if isErr(err, "DescribeNodegroupNotFound") || isErr(err, "InvalidParameterValue") {
			return Resource{}, nil
		}
		return Resource{}, err
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *output.Nodegroup.NodegroupArn,
		Name:   *output.Nodegroup.NodegroupName,
		Description: model.EKSNodegroupDescription{
			Nodegroup: *output.Nodegroup,
		},
	}
	return resource, nil
}
func GetEKSNodegroup(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	clusterName := fields["name"]

	client := eks.NewFromConfig(cfg)
	groups, err := client.ListNodegroups(ctx, &eks.ListNodegroupsInput{
		ClusterName: &clusterName,
	})
	if err != nil {
		if isErr(err, "ListNodegroupsNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	for _, profile := range groups.Nodegroups {
		resource, err := eKSNodegroupHandle(ctx, cfg, profile, clusterName)
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

func EKSIdentityProviderConfig(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	clusters, err := listEksClusters(ctx, cfg)
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, cluster := range clusters {
		cluster := cluster
		client := eks.NewFromConfig(cfg)
		paginator := eks.NewListIdentityProviderConfigsPaginator(client, &eks.ListIdentityProviderConfigsInput{
			ClusterName: &cluster,
		})

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				return nil, err
			}

			for _, config := range page.IdentityProviderConfigs {
				config := config
				output, err := client.DescribeIdentityProviderConfig(ctx, &eks.DescribeIdentityProviderConfigInput{
					ClusterName:            &cluster,
					IdentityProviderConfig: &config,
				})
				if err != nil {
					return nil, err
				}

				resource := Resource{
					Region: describeCtx.KaytuRegion,
					ARN:    *output.IdentityProviderConfig.Oidc.IdentityProviderConfigArn,
					Name:   *config.Name,
					Description: EKSIdentityProviderConfigDescription{
						ConfigName:             *config.Name,
						ConfigType:             *config.Type,
						IdentityProviderConfig: *output.IdentityProviderConfig.Oidc,
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

func EKSAddonVersion(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)

	client := eks.NewFromConfig(cfg)
	paginator := eks.NewDescribeAddonVersionsPaginator(client, &eks.DescribeAddonVersionsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, addon := range page.Addons {
			for _, version := range addon.AddonVersions {
				arn := fmt.Sprintf("arn:%s:eks:%s:%s:addonversion/%s/%s", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID, *addon.AddonName, *version.AddonVersion)

				configuration, err := client.DescribeAddonConfiguration(ctx, &eks.DescribeAddonConfigurationInput{
					AddonName:    addon.AddonName,
					AddonVersion: version.AddonVersion,
				})
				if err != nil {
					return nil, err
				}

				resource := Resource{
					Region: describeCtx.KaytuRegion,
					ARN:    arn,
					Description: model.EKSAddonVersionDescription{
						AddonVersion:       version,
						AddonConfiguration: configuration.ConfigurationSchema,
						AddonName:          addon.AddonName,
						AddonType:          addon.Type,
					},
				}
				if version.AddonVersion != nil {
					resource.Name = *version.AddonVersion
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
