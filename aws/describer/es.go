package describer

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/elasticsearchservice/types"
	"math"

	"github.com/aws/aws-sdk-go-v2/aws"
	es "github.com/aws/aws-sdk-go-v2/service/elasticsearchservice"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func ESDomain(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	var values []Resource

	client := es.NewFromConfig(cfg)
	listNamesOut, err := client.ListDomainNames(ctx, &es.ListDomainNamesInput{})
	if err != nil {
		return nil, err
	}

	var domainNamesList []string
	for _, dn := range listNamesOut.DomainNames {
		domainNamesList = append(domainNamesList, *dn.DomainName)
	}

	if len(domainNamesList) == 0 {
		return values, nil
	}

	for i := 0; i < len(domainNamesList); i += 5 {
		out, err := client.DescribeElasticsearchDomains(ctx, &es.DescribeElasticsearchDomainsInput{
			DomainNames: domainNamesList[i:int(math.Min(float64(i+5), float64(len(domainNamesList))))],
		})
		if err != nil {
			return nil, err
		}

		for _, v := range out.DomainStatusList {
			resource, err := ESDomainHandle(ctx, cfg, v)
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
func ESDomainHandle(ctx context.Context, cfg aws.Config, v types.ElasticsearchDomainStatus) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := es.NewFromConfig(cfg)

	out, err := client.ListTags(ctx, &es.ListTagsInput{
		ARN: v.ARN,
	})
	if err != nil {
		if isErr(err, "ListTagsNotFound") || isErr(err, "InvalidParameterValue") {
			return Resource{}, nil
		}
		return Resource{}, err
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *v.ARN,
		Name:   *v.DomainName,
		Description: model.ESDomainDescription{
			Domain: v,
			Tags:   out.TagList,
		},
	}
	return resource, nil
}
func GetESDomain(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	domainName := fields["domainName"]
	var values []Resource

	client := es.NewFromConfig(cfg)
	describeElasticSearch, err := client.DescribeElasticsearchDomain(ctx, &es.DescribeElasticsearchDomainInput{
		DomainName: &domainName,
	})
	if err != nil {
		return nil, err
	}

	resource, err := ESDomainHandle(ctx, cfg, *describeElasticSearch.DomainStatus)
	emptyResource := Resource{}
	if err == nil && resource == emptyResource {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	values = append(values, resource)
	return values, nil
}
