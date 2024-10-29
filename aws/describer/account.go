package describer

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/account"
	"github.com/aws/aws-sdk-go-v2/service/account/types"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func AccountAlternateContact(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)

	var values []Resource

	contactTypes := []types.AlternateContactType{types.AlternateContactTypeBilling, types.AlternateContactTypeOperations, types.AlternateContactTypeSecurity}
	input := &account.GetAlternateContactInput{
		AccountId: &describeCtx.AccountID,
	}
	for _, contactType := range contactTypes {
		input.AlternateContactType = contactType
		resource, err := accountAlternateContactHandle(ctx, cfg, *input.AccountId, input.AlternateContactType)
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

	return values, nil
}
func accountAlternateContactHandle(ctx context.Context, cfg aws.Config, accountId string, contactType types.AlternateContactType) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)

	client := account.NewFromConfig(cfg)
	op, err := client.GetAlternateContact(ctx, &account.GetAlternateContactInput{
		AlternateContactType: contactType,
		AccountId:            &accountId,
	})
	if err != nil {
		if isErr(err, "ResourceNotFoundException") {
			op = &account.GetAlternateContactOutput{}
		}
		return Resource{}, err
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		Name:   *op.AlternateContact.Name,
		Description: model.AccountAlternateContactDescription{
			AlternateContact: *op.AlternateContact,
			LinkedAccountID:  describeCtx.AccountID,
		},
	}
	return resource, nil
}
func GetAccountAlternateContact(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	accountId := fields["accountId"]
	contactTypes := []types.AlternateContactType{types.AlternateContactTypeBilling, types.AlternateContactTypeOperations, types.AlternateContactTypeSecurity}
	var values []Resource
	for _, contactType := range contactTypes {
		resource, err := accountAlternateContactHandle(ctx, cfg, accountId, contactType)
		if err != nil {
			return nil, err
		}
		emptyResource := Resource{}
		if err == nil && resource == emptyResource {
			return nil, nil
		}
		values = append(values, resource)
	}
	return values, nil
}

func AccountContact(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)

	client := account.NewFromConfig(cfg)

	var values []Resource

	input := &account.GetContactInformationInput{}
	op, err := client.GetContactInformation(ctx, input)
	if err != nil {
		return nil, err
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		Name:   *op.ContactInformation.FullName,
		Description: model.AccountContactDescription{
			AlternateContact: *op.ContactInformation,
			LinkedAccountID:  describeCtx.AccountID,
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

	return values, nil
}
