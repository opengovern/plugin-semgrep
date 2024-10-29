package describer

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func SecretsManagerSecret(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := secretsmanager.NewFromConfig(cfg)
	paginator := secretsmanager.NewListSecretsPaginator(client, &secretsmanager.ListSecretsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range page.SecretList {
			resource, err := secretsManagerSecretHandle(ctx, cfg, item.ARN, item.Name)
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
func secretsManagerSecretHandle(ctx context.Context, cfg aws.Config, Arn *string, Name *string) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := secretsmanager.NewFromConfig(cfg)
	out, err := client.DescribeSecret(ctx, &secretsmanager.DescribeSecretInput{
		SecretId: Arn,
	})
	if err != nil {
		return Resource{}, err
	}

	policy, err := client.GetResourcePolicy(ctx, &secretsmanager.GetResourcePolicyInput{
		SecretId: Arn,
	})
	if err != nil {
		return Resource{}, err
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *Arn,
		Name:   *Name,
		Description: model.SecretsManagerSecretDescription{
			Secret:         out,
			ResourcePolicy: policy.ResourcePolicy,
		},
	}
	return resource, nil
}
func GetSecretsManagerSecret(ctx context.Context, cfg aws.Config, field map[string]string) ([]Resource, error) {
	secretId := field["id"]
	var values []Resource

	client := secretsmanager.NewFromConfig(cfg)
	secretValue, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: &secretId,
	})
	if err != nil {
		return nil, err
	}

	resource, err := secretsManagerSecretHandle(ctx, cfg, secretValue.ARN, secretValue.Name)
	if err != nil {
		return nil, err
	}

	values = append(values, resource)
	return values, nil
}
