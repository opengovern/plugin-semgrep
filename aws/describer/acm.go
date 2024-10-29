package describer

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/acmpca/types"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acmpca"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func CertificateManagerAccount(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := acm.NewFromConfig(cfg)
	output, err := client.GetAccountConfiguration(ctx, &acm.GetAccountConfigurationInput{})
	if err != nil {
		return nil, err
	}

	return []Resource{
		{
			Region: describeCtx.KaytuRegion,
			// No ID or ARN. Per Account Configuration
			Description: output,
		}}, nil
}

func CertificateManagerCertificate(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := acm.NewFromConfig(cfg)
	paginator := acm.NewListCertificatesPaginator(client, &acm.ListCertificatesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.CertificateSummaryList {
			getOutput, err := client.GetCertificate(ctx, &acm.GetCertificateInput{
				CertificateArn: v.CertificateArn,
			})
			if err != nil {
				if strings.Contains(err.Error(), "not yet issued") {
					continue
				}
				return nil, err
			}

			describeOutput, err := client.DescribeCertificate(ctx, &acm.DescribeCertificateInput{
				CertificateArn: v.CertificateArn,
			})
			if err != nil {
				if isErr(err, "AccessDeniedException") {
					return nil, nil
				}
				return nil, err
			}

			tagsOutput, err := client.ListTagsForCertificate(ctx, &acm.ListTagsForCertificateInput{
				CertificateArn: v.CertificateArn,
			})
			if err != nil {
				return nil, err
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    *v.CertificateArn,
				Name:   nameFromArn(*v.CertificateArn),
				Description: model.CertificateManagerCertificateDescription{
					Certificate: *describeOutput.Certificate,
					Attributes: struct {
						Certificate      *string
						CertificateChain *string
					}{
						Certificate:      getOutput.Certificate,
						CertificateChain: getOutput.CertificateChain,
					},
					Tags: tagsOutput.Tags,
				},
			}
			if stream != nil {
				m := *stream
				err := m(resource)
				if err != nil {
					return nil, err
				}
			} else {
				values = append(values, resource)
			}
		}
	}

	return values, nil
}

func ACMPCACertificateAuthority(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := acmpca.NewFromConfig(cfg)
	paginator := acmpca.NewListCertificateAuthoritiesPaginator(client, &acmpca.ListCertificateAuthoritiesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.CertificateAuthorities {
			resource, err := aCMPCACertificateAuthorityHandle(ctx, cfg, v)
			if err != nil {
				return nil, err
			}
			emptyResource := Resource{}
			if err == nil && resource == emptyResource {
				continue
			}

			if stream != nil {
				m := *stream
				err := m(resource)
				if err != nil {
					return nil, err
				}
			} else {
				values = append(values, resource)
			}
		}
	}

	return values, nil
}
func aCMPCACertificateAuthorityHandle(ctx context.Context, cfg aws.Config, v types.CertificateAuthority) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := acmpca.NewFromConfig(cfg)
	tags, err := client.ListTags(ctx, &acmpca.ListTagsInput{
		CertificateAuthorityArn: v.Arn,
	})
	if err != nil {
		if isErr(err, "ListTagsNotFound") || isErr(err, "InvalidParameterValue") {
			return Resource{}, nil
		}
		return Resource{}, err
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *v.Arn,
		Name:   nameFromArn(*v.Arn),
		Description: model.ACMPCACertificateAuthorityDescription{
			CertificateAuthority: v,
			Tags:                 tags.Tags,
		},
	}
	return resource, nil
}
func GetACMPCACertificateAuthority(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	certificateAuthorityArn := fields["arn"]
	client := acmpca.NewFromConfig(cfg)

	out, err := client.DescribeCertificateAuthority(ctx, &acmpca.DescribeCertificateAuthorityInput{
		CertificateAuthorityArn: &certificateAuthorityArn,
	})
	if err != nil {
		if isErr(err, "DescribeCertificateAuthorityNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	resource, err := aCMPCACertificateAuthorityHandle(ctx, cfg, *out.CertificateAuthority)
	if err != nil {
		return nil, err
	}
	emptyResource := Resource{}
	if err == nil && resource == emptyResource {
		return nil, nil
	}

	values = append(values, resource)
	return values, nil
}
