package describer

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/smithy-go"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func LambdaFunction(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	logger := GetLoggerFromContext(ctx)

	logger.Info("LambdaFunction start working")

	client := lambda.NewFromConfig(cfg)
	paginator := lambda.NewListFunctionsPaginator(client, &lambda.ListFunctionsInput{})

	logger.Info("LambdaFunction start getting pages")
	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		logger.Info("LambdaFunction got page")
		for _, v := range page.Functions {
			if v.FunctionName == nil {
				continue
			}
			resource, err := lambdaFunctionHandle(ctx, client, v)
			if err != nil {
				return nil, err
			}
			emptyResource := Resource{}
			if err == nil && resource == emptyResource {
				continue
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
	logger.Info("LambdaFunction finished")

	return values, nil
}
func lambdaFunctionHandle(ctx context.Context, client *lambda.Client, v types.FunctionConfiguration) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)

	listUrlConfig, err := client.ListFunctionUrlConfigs(ctx, &lambda.ListFunctionUrlConfigsInput{
		FunctionName: v.FunctionName,
	})
	if err != nil {
		if isErr(err, "ListFunctionUrlConfigsNotFound") || isErr(err, "InvalidParameterValue") {
			listUrlConfig = &lambda.ListFunctionUrlConfigsOutput{}
		} else {
			return Resource{}, nil
		}
	}

	policy, err := client.GetPolicy(ctx, &lambda.GetPolicyInput{
		FunctionName: v.FunctionName,
	})
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) && ae.ErrorCode() == "ResourceNotFoundException" {
			policy = &lambda.GetPolicyOutput{}
			err = nil
		}

		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "ResourceNotFoundException" {
				policy = &lambda.GetPolicyOutput{}
				err = nil
			}
		}

		if err != nil {
			return Resource{}, err
		}
	}

	function, err := client.GetFunction(ctx, &lambda.GetFunctionInput{
		FunctionName: v.FunctionName,
	})
	if err != nil {
		return Resource{}, err
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *v.FunctionArn,
		Name:   *v.FunctionName,
		Description: model.LambdaFunctionDescription{
			Function:  function,
			UrlConfig: listUrlConfig.FunctionUrlConfigs,
			Policy:    policy,
		},
	}
	return resource, nil
}
func GetLambdaFunction(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	functionName := fields["name"]
	client := lambda.NewFromConfig(cfg)
	out, err := client.GetFunction(ctx, &lambda.GetFunctionInput{
		FunctionName: &functionName,
		Qualifier:    nil,
	})
	if err != nil {
		return nil, err
	}

	var values []Resource

	resource, err := lambdaFunctionHandle(ctx, client, *out.Configuration)
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

func LambdaFunctionVersion(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := lambda.NewFromConfig(cfg)
	paginator := lambda.NewListFunctionsPaginator(client, &lambda.ListFunctionsInput{
		FunctionVersion: types.FunctionVersionAll,
	})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Functions {
			id := fmt.Sprintf("%s:%s", *v.FunctionArn, *v.Version)

			policy, err := client.GetPolicy(ctx, &lambda.GetPolicyInput{
				FunctionName: v.FunctionName,
				Qualifier:    v.Version,
			})
			if err != nil {
				if isErr(err, "ResourceNotFoundException") {
					policy = &lambda.GetPolicyOutput{}
				} else {
					return nil, err
				}
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ID:     id,
				Description: model.LambdaFunctionVersionDescription{
					FunctionVersion: v,
					Policy:          policy,
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

func LambdaAlias(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	fns, err := LambdaFunction(ctx, cfg, nil)
	if err != nil {
		return nil, err
	}

	client := lambda.NewFromConfig(cfg)

	var values []Resource
	for _, f := range fns {
		fn := f.Description.(model.LambdaFunctionDescription).Function.Configuration
		paginator := lambda.NewListAliasesPaginator(client, &lambda.ListAliasesInput{
			FunctionName:    fn.FunctionName,
			FunctionVersion: fn.Version,
		})

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				if isErr(err, "ResourceNotFoundException") {
					continue
				}
				return nil, err
			}

			for _, v := range page.Aliases {
				resource, err := LambdaAliasHandle(ctx, cfg, v, fn)
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
	}
	return values, nil
}
func LambdaAliasHandle(ctx context.Context, cfg aws.Config, v types.AliasConfiguration, fn *types.FunctionConfiguration) (Resource, error) {
	client := lambda.NewFromConfig(cfg)
	describeCtx := GetDescribeContext(ctx)

	policy, err := client.GetPolicy(ctx, &lambda.GetPolicyInput{
		FunctionName: fn.FunctionName,
		Qualifier:    v.Name,
	})
	if err != nil {
		if isErr(err, "ResourceNotFoundException") {
			policy = &lambda.GetPolicyOutput{}
		} else {
			return Resource{}, err
		}
	}

	urlConfig, err := client.GetFunctionUrlConfig(ctx, &lambda.GetFunctionUrlConfigInput{
		FunctionName: fn.FunctionName,
		Qualifier:    v.Name,
	})
	if err != nil {
		if isErr(err, "ResourceNotFoundException") {
			urlConfig = &lambda.GetFunctionUrlConfigOutput{}
		} else {
			return Resource{}, err
		}
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *v.AliasArn,
		Name:   *v.Name,
		Description: model.LambdaAliasDescription{
			FunctionName: *fn.FunctionName,
			Alias:        v,
			Policy:       policy,
			UrlConfig:    *urlConfig,
		},
	}
	return resource, nil
}
func GetLambdaAlias(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	name := fields["aliasName"]
	client := lambda.NewFromConfig(cfg)

	fns, err := client.ListFunctions(ctx, &lambda.ListFunctionsInput{})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, fn := range fns.Functions {

		out, err := client.GetAlias(ctx, &lambda.GetAliasInput{
			Name:         &name,
			FunctionName: fn.FunctionName,
		})
		if err != nil {
			return nil, err
		}

		alias := types.AliasConfiguration{
			AliasArn:        out.AliasArn,
			Name:            out.Name,
			Description:     out.Description,
			FunctionVersion: out.FunctionVersion,
			RevisionId:      out.RevisionId,
			RoutingConfig:   out.RoutingConfig,
		}
		resource, err := LambdaAliasHandle(ctx, cfg, alias, &fn)
		if err != nil {
			return nil, err
		}

		values = append(values, resource)
	}
	return values, nil
}

func LambdaPermission(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	fns, err := LambdaFunction(ctx, cfg, nil)
	if err != nil {
		return nil, err
	}

	client := lambda.NewFromConfig(cfg)

	var values []Resource
	for _, f := range fns {
		fn := f.Description.(model.LambdaFunctionDescription).Function.Configuration
		v, err := client.GetPolicy(ctx, &lambda.GetPolicyInput{
			FunctionName: fn.FunctionArn,
		})
		if err != nil {
			var ae smithy.APIError
			if errors.As(err, &ae) && ae.ErrorCode() == "ResourceNotFoundException" {
				continue
			}

			return nil, err
		}

		resource := Resource{
			Region:      describeCtx.KaytuRegion,
			ID:          CompositeID(*fn.FunctionArn, *v.Policy),
			Name:        *v.Policy,
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

func LambdaEventInvokeConfig(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	fns, err := LambdaFunction(ctx, cfg, nil)
	if err != nil {
		return nil, err
	}

	client := lambda.NewFromConfig(cfg)

	var values []Resource
	for _, f := range fns {
		fn := f.Description.(model.LambdaFunctionDescription).Function.Configuration
		paginator := lambda.NewListFunctionEventInvokeConfigsPaginator(client, &lambda.ListFunctionEventInvokeConfigsInput{
			FunctionName: fn.FunctionName,
		})

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				return nil, err
			}

			for _, v := range page.FunctionEventInvokeConfigs {
				resource := Resource{
					Region:      describeCtx.KaytuRegion,
					ID:          *fn.FunctionName, // Invoke Config is unique per function
					Name:        *fn.FunctionName,
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
	}

	return values, nil
}

func LambdaCodeSigningConfig(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := lambda.NewFromConfig(cfg)
	paginator := lambda.NewListCodeSigningConfigsPaginator(client, &lambda.ListCodeSigningConfigsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.CodeSigningConfigs {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ARN:         *v.CodeSigningConfigArn,
				Name:        *v.CodeSigningConfigArn,
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

func LambdaEventSourceMapping(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := lambda.NewFromConfig(cfg)
	paginator := lambda.NewListEventSourceMappingsPaginator(client, &lambda.ListEventSourceMappingsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.EventSourceMappings {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ARN:         *v.EventSourceArn,
				Name:        *v.UUID,
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

func LambdaLayerVersion(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	layers, err := listLayers(ctx, cfg)
	if err != nil {
		return nil, err
	}

	client := lambda.NewFromConfig(cfg)

	var values []Resource
	for _, layer := range layers {
		paginator := lambda.NewListLayerVersionsPaginator(client, &lambda.ListLayerVersionsInput{
			LayerName: layer.LayerArn,
		})

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				return nil, err
			}

			for _, v := range page.LayerVersions {
				resource, err := lambdaLayerVersionHandle(ctx, cfg, layer, v)
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
	}

	return values, nil
}
func lambdaLayerVersionHandle(ctx context.Context, cfg aws.Config, layer types.LayersListItem, v types.LayerVersionsListItem) (Resource, error) {
	client := lambda.NewFromConfig(cfg)
	describeCtx := GetDescribeContext(ctx)
	layerVersion, err := client.GetLayerVersion(ctx, &lambda.GetLayerVersionInput{
		LayerName:     layer.LayerArn,
		VersionNumber: &v.Version,
	})
	if err != nil {
		return Resource{}, err
	}

	policy, err := client.GetLayerVersionPolicy(ctx, &lambda.GetLayerVersionPolicyInput{
		LayerName:     layer.LayerArn,
		VersionNumber: &v.Version,
	})
	if err != nil {
		if isErr(err, "ResourceNotFoundException") {
			policy = &lambda.GetLayerVersionPolicyOutput{}
		} else {
			return Resource{}, err
		}
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *v.LayerVersionArn,
		Name:   *v.LayerVersionArn,
		Description: model.LambdaLayerVersionDescription{
			LayerName:    *layer.LayerName,
			LayerVersion: *layerVersion,
			Policy:       *policy,
		},
	}
	return resource, nil

}
func GetLambdaLayerVersion(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	layerArn := fields["arn"]
	client := lambda.NewFromConfig(cfg)

	layers, err := client.ListLayers(ctx, &lambda.ListLayersInput{})
	if err != nil {
		if isErr(err, "ListLayersNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	for _, layer := range layers.Layers {
		if *layer.LayerArn != layerArn {
			continue
		}

		out, err := client.ListLayerVersions(ctx, &lambda.ListLayerVersionsInput{
			LayerName: &layerArn,
		})
		if err != nil {
			if isErr(err, "ListLayerVersionsNotFound") || isErr(err, "InvalidParameterValue") {
				return nil, nil
			}
			return nil, err
		}

		for _, v := range out.LayerVersions {

			resource, err := lambdaLayerVersionHandle(ctx, cfg, layer, v)
			if err != nil {
				return nil, err
			}
			values = append(values, resource)

		}
	}
	return values, nil
}

func LambdaLayer(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	layers, err := listLayers(ctx, cfg)
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, layer := range layers {
		resource := Resource{
			Region: describeCtx.KaytuRegion,
			ARN:    *layer.LayerArn,
			Name:   *layer.LayerName,
			Description: model.LambdaLayerDescription{
				Layer: layer,
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

	return values, nil
}

func LambdaLayerVersionPermission(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	lvs, err := LambdaLayerVersion(ctx, cfg, nil)
	if err != nil {
		return nil, err
	}

	client := lambda.NewFromConfig(cfg)

	var values []Resource
	for _, lv := range lvs {
		arn := lv.Description.(model.LambdaLayerVersionDescription).LayerVersion.LayerVersionArn
		version := lv.Description.(model.LambdaLayerVersionDescription).LayerVersion.Version
		v, err := client.GetLayerVersionPolicy(ctx, &lambda.GetLayerVersionPolicyInput{
			LayerName:     arn,
			VersionNumber: &version,
		})
		if err != nil {
			return nil, err
		}

		resource := Resource{
			Region:      describeCtx.KaytuRegion,
			ID:          CompositeID(*arn, fmt.Sprintf("%d", version)),
			Name:        *arn,
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

func listLayers(ctx context.Context, cfg aws.Config) ([]types.LayersListItem, error) {
	client := lambda.NewFromConfig(cfg)
	paginator := lambda.NewListLayersPaginator(client, &lambda.ListLayersInput{})

	var values []types.LayersListItem
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		values = append(values, page.Layers...)
	}

	return values, nil
}
