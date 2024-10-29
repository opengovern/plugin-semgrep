package describer

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	sesv2types "github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func SESConfigurationSet(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := sesv2.NewFromConfig(cfg)
	paginator := sesv2.NewListConfigurationSetsPaginator(client, &sesv2.ListConfigurationSetsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.ConfigurationSets {

			resource, err := sESConfigurationSetHandle(ctx, cfg, v)
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
func sESConfigurationSetHandle(ctx context.Context, cfg aws.Config, v string) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)

	sesClient := ses.NewFromConfig(cfg)

	output, err := sesClient.DescribeConfigurationSet(ctx, &ses.DescribeConfigurationSetInput{ConfigurationSetName: aws.String(v)})
	if err != nil {
		if isErr(err, "DescribeConfigurationSetNotFound") || isErr(err, "InvalidParameterValue") {
			return Resource{}, nil
		}
		return Resource{}, err
	}

	arn := fmt.Sprintf("arn:%s:ses:%s:%s:configuration-set/%s", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID, *output.ConfigurationSet.Name)

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *output.ConfigurationSet.Name,
		Description: model.SESConfigurationSetDescription{
			ConfigurationSet: *output.ConfigurationSet,
		},
	}
	return resource, nil
}
func GetSESConfigurationSet(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	configurationSetName := fields["name"]
	var values []Resource

	resource, err := sESConfigurationSetHandle(ctx, cfg, configurationSetName)
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

func SESIdentity(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := ses.NewFromConfig(cfg)
	paginator := ses.NewListIdentitiesPaginator(client, &ses.ListIdentitiesInput{})

	var values []Resource
	for paginator.HasMorePages() {

		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Identities {

			resource, err := sESIdentityHandle(ctx, cfg, v)
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

func SESv2EmailIdentities(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := sesv2.NewFromConfig(cfg)
	paginator := sesv2.NewListEmailIdentitiesPaginator(client, &sesv2.ListEmailIdentitiesInput{})

	var values []Resource
	for paginator.HasMorePages() {

		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.EmailIdentities {

			resource, err := sESv2EmailIdentitiesHandle(ctx, cfg, v)
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

func sESv2EmailIdentitiesHandle(ctx context.Context, cfg aws.Config, v sesv2types.IdentityInfo) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)

	client := sesv2.NewFromConfig(cfg)

	arn := fmt.Sprintf("arn:%s:sesv2:%s:%s:identity/%s", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID, v)

	tags, err := client.ListTagsForResource(ctx, &sesv2.ListTagsForResourceInput{
		ResourceArn: &arn,
	})
	if err != nil {
		if isErr(err, "ListTagsForResourceNotFound") || isErr(err, "InvalidParameterValue") {
			return Resource{}, nil
		}
		return Resource{}, err
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *v.IdentityName,
		Description: model.SESv2EmailIdentityDescription{
			ARN:      arn,
			Identity: v,
			Tags:     tags.Tags,
		},
	}
	return resource, nil
}

func sESIdentityHandle(ctx context.Context, cfg aws.Config, v string) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)

	client := ses.NewFromConfig(cfg)

	arn := fmt.Sprintf("arn:%s:ses:%s:%s:identity/%s", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID, v)

	identity, err := client.GetIdentityVerificationAttributes(ctx, &ses.GetIdentityVerificationAttributesInput{
		Identities: []string{v},
	})
	if err != nil {
		if isErr(err, "GetIdentityVerificationAttributesNotFound") || isErr(err, "InvalidParameterValue") {
			return Resource{}, nil
		}
		return Resource{}, err
	}

	notif, err := client.GetIdentityNotificationAttributes(ctx, &ses.GetIdentityNotificationAttributesInput{
		Identities: []string{v},
	})
	if err != nil {
		if isErr(err, "GetIdentityNotificationAttributesNotFound") || isErr(err, "InvalidParameterValue") {
			return Resource{}, nil
		}
		return Resource{}, err
	}

	DkimAtrb, err := client.GetIdentityDkimAttributes(ctx, &ses.GetIdentityDkimAttributesInput{
		Identities: []string{v},
	})
	if err != nil {
		if isErr(err, "GetIdentityDkimAttributesNotFound") || isErr(err, "InvalidParameterValue") {
			return Resource{}, nil
		}
		return Resource{}, err
	}

	identityMail, err := client.GetIdentityMailFromDomainAttributes(ctx, &ses.GetIdentityMailFromDomainAttributesInput{
		Identities: []string{v},
	})
	if err != nil {
		if isErr(err, "GetIdentityMailFromDomainAttributesNotFound") || isErr(err, "InvalidParameterValue") {
			return Resource{}, nil
		}
		return Resource{}, err
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   v,
		Description: model.SESIdentityDescription{
			ARN:                    arn,
			Identity:               v,
			DkimAttributes:         DkimAtrb.DkimAttributes,
			IdentityMail:           identityMail.MailFromDomainAttributes,
			VerificationAttributes: identity.VerificationAttributes[v],
			NotificationAttributes: notif.NotificationAttributes[v],
		},
	}
	return resource, nil
}
func GetSESIdentity(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	IdentityType := fields["identityType"]
	client := ses.NewFromConfig(cfg)

	out, err := client.ListIdentities(ctx, &ses.ListIdentitiesInput{
		IdentityType: types.IdentityType(IdentityType),
	})
	if err != nil {
		if isErr(err, "DescribeConfigurationSetNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	for _, v := range out.Identities {
		resource, err := sESIdentityHandle(ctx, cfg, v)
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

func SESContactList(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := sesv2.NewFromConfig(cfg)
	paginator := sesv2.NewListContactListsPaginator(client, &sesv2.ListContactListsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.ContactLists {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ID:          *v.ContactListName,
				Name:        *v.ContactListName,
				Description: v,
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

func SESReceiptFilter(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ses.NewFromConfig(cfg)

	output, err := client.ListReceiptFilters(ctx, &ses.ListReceiptFiltersInput{})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range output.Filters {
		resource := Resource{
			Region:      describeCtx.KaytuRegion,
			ID:          *v.Name,
			Name:        *v.Name,
			Description: v,
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

func SESReceiptRuleSet(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ses.NewFromConfig(cfg)

	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		output, err := client.ListReceiptRuleSets(ctx, &ses.ListReceiptRuleSetsInput{NextToken: prevToken})
		if err != nil {
			return nil, err
		}

		for _, v := range output.RuleSets {
			output, err := client.DescribeReceiptRuleSet(ctx, &ses.DescribeReceiptRuleSetInput{RuleSetName: v.Name})
			if err != nil {
				return nil, err
			}

			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ID:          *output.Metadata.Name,
				Name:        *output.Metadata.Name,
				Description: output,
			}
			if stream != nil {
				if err := (*stream)(resource); err != nil {
					return nil, err
				}
			} else {
				values = append(values, resource)
			}

		}

		return output.NextToken, nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}

func SESTemplate(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ses.NewFromConfig(cfg)

	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		output, err := client.ListTemplates(ctx, &ses.ListTemplatesInput{NextToken: prevToken})
		if err != nil {
			return nil, err
		}

		for _, v := range output.TemplatesMetadata {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ID:          *v.Name,
				Name:        *v.Name,
				Description: v,
			}
			if stream != nil {
				if err := (*stream)(resource); err != nil {
					return nil, err
				}
			} else {
				values = append(values, resource)
			}

		}

		return output.NextToken, nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}
