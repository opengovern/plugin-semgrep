package describer

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/cloudsearch/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudsearch"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func CloudSearchDomain(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := cloudsearch.NewFromConfig(cfg)
	var values []Resource

	output, err := client.ListDomainNames(ctx, &cloudsearch.ListDomainNamesInput{})
	if err != nil {
		return nil, err
	}

	var domainList []string
	for domainName := range output.DomainNames {
		domainList = append(domainList, domainName)
	}

	domains, err := client.DescribeDomains(ctx, &cloudsearch.DescribeDomainsInput{
		DomainNames: domainList,
	})
	if err != nil {
		return nil, err
	}

	for _, domain := range domains.DomainStatusList {
		resource := cloudSearchDomainHandle(ctx, domain)
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
func cloudSearchDomainHandle(ctx context.Context, domain types.DomainStatus) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *domain.ARN,
		Name:   *domain.DomainName,
		ID:     *domain.DomainId,
		Description: model.CloudSearchDomainDescription{
			DomainStatus: domain,
		},
	}
	return resource
}
func GetCloudSearchDomain(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	domainList := fields["domainList"]
	client := cloudsearch.NewFromConfig(cfg)

	var values []Resource
	domains, err := client.DescribeDomains(ctx, &cloudsearch.DescribeDomainsInput{
		DomainNames: []string{domainList},
	})
	if err != nil {
		return nil, err
	}

	for _, domain := range domains.DomainStatusList {
		resource := cloudSearchDomainHandle(ctx, domain)
		values = append(values, resource)
	}
	return values, nil
}
