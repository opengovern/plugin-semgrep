package describer

import (
	"context"
	_ "database/sql/driver"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	typesv2 "github.com/aws/aws-sdk-go-v2/service/apigatewayv2/types"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func ApiGatewayStage(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := apigateway.NewFromConfig(cfg)
	paginator := apigateway.NewGetRestApisPaginator(client, &apigateway.GetRestApisInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, restItem := range page.Items {
			out, err := client.GetStages(ctx, &apigateway.GetStagesInput{
				RestApiId: restItem.Id,
			})
			if err != nil {
				return nil, err
			}

			for _, stageItem := range out.Item {
				resource := apiGatewayStageHandle(ctx, stageItem, *restItem.Id, *restItem.Name)
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
	}
	return values, nil
}
func apiGatewayStageHandle(ctx context.Context, stageItem types.Stage, id string, name string) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := "arn:" + describeCtx.Partition + ":apigateway:" + describeCtx.Region + "::/restapis/" + id + "/stages/" + *stageItem.StageName
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   name,
		Description: model.ApiGatewayStageDescription{
			RestApiId: &id,
			Stage:     stageItem,
		},
	}
	return resource
}
func GetApiGatewayStage(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	restAPIID := fields["respApiId"]
	stageName := fields["stageName"]
	client := apigateway.NewFromConfig(cfg)
	var values []Resource

	stage, err := client.GetStage(ctx, &apigateway.GetStageInput{
		RestApiId: &restAPIID,
		StageName: &stageName,
	})
	if err != nil {
		if isErr(err, "GetStageNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	stageItem := types.Stage{
		Tags:                 stage.Tags,
		StageName:            stage.StageName,
		AccessLogSettings:    stage.AccessLogSettings,
		Description:          stage.Description,
		CreatedDate:          stage.CreatedDate,
		CacheClusterEnabled:  stage.CacheClusterEnabled,
		CacheClusterSize:     stage.CacheClusterSize,
		CacheClusterStatus:   stage.CacheClusterStatus,
		CanarySettings:       stage.CanarySettings,
		DeploymentId:         stage.DeploymentId,
		ClientCertificateId:  stage.ClientCertificateId,
		WebAclArn:            stage.WebAclArn,
		Variables:            stage.Variables,
		TracingEnabled:       stage.TracingEnabled,
		MethodSettings:       stage.MethodSettings,
		LastUpdatedDate:      stage.LastUpdatedDate,
		DocumentationVersion: stage.DocumentationVersion,
	}
	values = append(values, apiGatewayStageHandle(ctx, stageItem, restAPIID, *stageItem.StageName))
	return values, nil
}

func ApiGatewayV2Stage(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := apigatewayv2.NewFromConfig(cfg)

	var apis []typesv2.Api
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		output, err := client.GetApis(ctx, &apigatewayv2.GetApisInput{
			NextToken: prevToken,
		})
		if err != nil {
			return nil, err
		}

		apis = append(apis, output.Items...)
		return output.NextToken, nil
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, api := range apis {
		var stages []typesv2.Stage
		err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
			output, err := client.GetStages(ctx, &apigatewayv2.GetStagesInput{
				ApiId:     api.ApiId,
				NextToken: prevToken,
			})
			if err != nil {
				return nil, err
			}

			stages = append(stages, output.Items...)
			return output.NextToken, nil
		})
		if err != nil {
			return nil, err
		}

		for _, stage := range stages {
			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ID:     CompositeID(*api.ApiId, *stage.StageName),
				Name:   *api.Name,
				Description: model.ApiGatewayV2StageDescription{
					ApiId: api.ApiId,
					Stage: stage,
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

func ApiGatewayRestAPI(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := apigateway.NewFromConfig(cfg)
	paginator := apigateway.NewGetRestApisPaginator(client, &apigateway.GetRestApisInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if isErr(err, "NotFoundException") {
				continue
			}
			return nil, err
		}

		for _, restItem := range page.Items {
			resource := apiGatewayRestAPIHandle(ctx, restItem)
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
func apiGatewayRestAPIHandle(ctx context.Context, restItem types.RestApi) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := fmt.Sprintf("arn:%s:apigateway:%s::/restapis/%s", describeCtx.Partition, describeCtx.Region, *restItem.Id)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *restItem.Name,
		Description: model.ApiGatewayRestAPIDescription{
			RestAPI: restItem,
		},
	}
	return resource
}
func GetApiGatewayRestAPI(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	id := fields["id"]
	client := apigateway.NewFromConfig(cfg)

	out, err := client.GetRestApi(ctx, &apigateway.GetRestApiInput{
		RestApiId: &id,
	})
	if err != nil {
		if isErr(err, "NotFoundException") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	values = append(values, apiGatewayRestAPIHandle(ctx, types.RestApi{
		ApiKeySource:              out.ApiKeySource,
		BinaryMediaTypes:          out.BinaryMediaTypes,
		CreatedDate:               out.CreatedDate,
		Description:               out.Description,
		DisableExecuteApiEndpoint: out.DisableExecuteApiEndpoint,
		EndpointConfiguration:     out.EndpointConfiguration,
		Id:                        out.Id,
		MinimumCompressionSize:    out.MinimumCompressionSize,
		Name:                      out.Name,
		Policy:                    out.Policy,
		Tags:                      out.Tags,
		Version:                   out.Version,
		Warnings:                  out.Warnings,
	}))
	return values, nil
}

func ApiGatewayApiKey(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := apigateway.NewFromConfig(cfg)
	paginator := apigateway.NewGetApiKeysPaginator(client, &apigateway.GetApiKeysInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if isErr(err, "NotFoundException") {
				continue
			}
			return nil, err
		}

		for _, apiKey := range page.Items {
			resource := apiGatewayApiKeyHandle(ctx, apiKey)
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
func apiGatewayApiKeyHandle(ctx context.Context, apiKey types.ApiKey) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := fmt.Sprintf("arn:%s:apigateway:%s::/apikeys/%s", describeCtx.Partition, describeCtx.Region, *apiKey.Id)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ID:     *apiKey.Id,
		ARN:    arn,
		Name:   *apiKey.Name,
		Description: model.ApiGatewayApiKeyDescription{
			ApiKey: apiKey,
		},
	}
	return resource
}
func GetApiGatewayApiKey(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	apiKeyVa := fields["apiKeyVa"]
	client := apigateway.NewFromConfig(cfg)

	out, err := client.GetApiKey(ctx, &apigateway.GetApiKeyInput{
		ApiKey: &apiKeyVa,
	})
	if err != nil {
		if isErr(err, "GetApiKeyNotFound") || isErr(err, "invalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	values = append(values, apiGatewayApiKeyHandle(ctx, types.ApiKey{
		StageKeys:       out.StageKeys,
		Name:            out.Name,
		Id:              out.Id,
		Tags:            out.Tags,
		Description:     out.Description,
		LastUpdatedDate: out.LastUpdatedDate,
		CreatedDate:     out.CreatedDate,
		CustomerId:      out.CustomerId,
		Enabled:         out.Enabled,
		Value:           out.Value,
	}))
	return values, nil
}

func ApiGatewayUsagePlan(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := apigateway.NewFromConfig(cfg)
	paginator := apigateway.NewGetUsagePlansPaginator(client, &apigateway.GetUsagePlansInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if isErr(err, "NotFoundException") {
				continue
			}
			return nil, err
		}

		for _, usagePlan := range page.Items {
			resource := apiGatewayUsagePlanHandle(ctx, usagePlan)
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
func apiGatewayUsagePlanHandle(ctx context.Context, usagePlan types.UsagePlan) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := fmt.Sprintf("arn:%s:apigateway:%s::/usageplans/%s", describeCtx.Partition, describeCtx.Region, *usagePlan.Id)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ID:     *usagePlan.Id,
		ARN:    arn,
		Name:   *usagePlan.Name,
		Description: model.ApiGatewayUsagePlanDescription{
			UsagePlan: usagePlan,
		},
	}
	return resource
}
func GetApiGatewayUsagePlan(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	usagePlanId := fields["usagePlanId"]
	client := apigateway.NewFromConfig(cfg)

	out, err := client.GetUsagePlan(ctx, &apigateway.GetUsagePlanInput{
		UsagePlanId: &usagePlanId,
	})
	if err != nil {
		if isErr(err, "GetUsagePlanNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	values = append(values, apiGatewayUsagePlanHandle(ctx, types.UsagePlan{
		Name:        out.Name,
		Tags:        out.Tags,
		Id:          out.Id,
		Description: out.Description,
		ApiStages:   out.ApiStages,
		ProductCode: out.ProductCode,
		Quota:       out.Quota,
		Throttle:    out.Throttle,
	}))
	return values, nil
}

func ApiGatewayAuthorizer(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := apigateway.NewFromConfig(cfg)
	paginator := apigateway.NewGetRestApisPaginator(client, &apigateway.GetRestApisInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if isErr(err, "NotFoundException") {
				continue
			}
			return nil, err
		}

		for _, api := range page.Items {
			authorizers, err := client.GetAuthorizers(ctx, &apigateway.GetAuthorizersInput{
				RestApiId: api.Id,
			})
			if err != nil {
				return nil, err
			}
			for _, authorizer := range authorizers.Items {
				resource := apiGatewayAuthorizerHandle(ctx, authorizer, *api.Id, *api.Name)
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
	}
	return values, nil
}
func apiGatewayAuthorizerHandle(ctx context.Context, authorizer types.Authorizer, apiId string, apiName string) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := fmt.Sprintf("arn:%s:apigateway:%s::/restapis/%s/authorizer/%s", describeCtx.Partition, describeCtx.Region, apiId, *authorizer.Id)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ID:     *authorizer.Id,
		ARN:    arn,
		Name:   apiName,
		Description: model.ApiGatewayAuthorizerDescription{
			Authorizer: authorizer,
			RestApiId:  apiId,
		},
	}
	return resource
}
func GetApiGatewayAuthorizer(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	restApiId := fields["restApiId"]
	authorizerId := fields["authorizerId"]
	var values []Resource
	client := apigateway.NewFromConfig(cfg)

	out, err := client.GetAuthorizer(ctx, &apigateway.GetAuthorizerInput{
		AuthorizerId: &authorizerId,
		RestApiId:    &restApiId,
	})
	if err != nil {
		if isErr(err, "GetAuthorizerNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	authorizer := types.Authorizer{
		AuthorizerCredentials:        out.AuthorizerCredentials,
		Id:                           out.Id,
		AuthorizerUri:                out.AuthorizerUri,
		AuthType:                     out.AuthType,
		AuthorizerResultTtlInSeconds: out.AuthorizerResultTtlInSeconds,
		IdentitySource:               out.IdentitySource,
		IdentityValidationExpression: out.IdentityValidationExpression,
		ProviderARNs:                 out.ProviderARNs,
		Type:                         out.Type,
	}
	values = append(values, apiGatewayAuthorizerHandle(ctx, authorizer, restApiId, *authorizer.Name))
	return values, nil
}

func ApiGatewayV2API(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := apigatewayv2.NewFromConfig(cfg)

	var values []Resource
	pagesLeft := true
	params := &apigatewayv2.GetApisInput{}
	for pagesLeft {
		output, err := client.GetApis(ctx, params)
		if err != nil {
			return nil, err
		}
		for _, api := range output.Items {
			resource := apiGatewayV2APIHandle(ctx, api)
			if stream != nil {
				if err := (*stream)(resource); err != nil {
					return nil, err
				}
			} else {
				values = append(values, resource)
			}
		}
		if output.NextToken != nil {
			pagesLeft = true
			params.NextToken = output.NextToken
		} else {
			pagesLeft = false
		}
	}
	return values, nil
}

func apiGatewayV2APIHandle(ctx context.Context, api typesv2.Api) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := fmt.Sprintf("arn:%s:apigateway:%s::/apis/%s", describeCtx.Partition, describeCtx.Region, *api.ApiId)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *api.Name,
		Description: model.ApiGatewayV2APIDescription{
			API: api,
		},
	}
	return resource
}
func GetApiGatewayV2API(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	apiID := fields["id"]
	client := apigatewayv2.NewFromConfig(cfg)

	out, err := client.GetApi(ctx, &apigatewayv2.GetApiInput{
		ApiId: &apiID,
	})
	if err != nil {
		if isErr(err, "NotFoundException") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	values = append(values, apiGatewayV2APIHandle(ctx, typesv2.Api{
		Name:                      out.Name,
		ProtocolType:              out.ProtocolType,
		RouteSelectionExpression:  out.RouteSelectionExpression,
		ApiEndpoint:               out.ApiEndpoint,
		ApiGatewayManaged:         out.ApiGatewayManaged,
		ApiId:                     out.ApiId,
		ApiKeySelectionExpression: out.ApiKeySelectionExpression,
		CorsConfiguration:         out.CorsConfiguration,
		CreatedDate:               out.CreatedDate,
		Description:               out.Description,
		DisableExecuteApiEndpoint: out.DisableExecuteApiEndpoint,
		DisableSchemaValidation:   out.DisableSchemaValidation,
		ImportInfo:                out.ImportInfo,
		Tags:                      out.Tags,
		Version:                   out.Version,
		Warnings:                  out.Warnings,
	}))
	return values, nil
}

func ApiGatewayV2DomainName(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := apigatewayv2.NewFromConfig(cfg)
	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		output, err := client.GetDomainNames(ctx, &apigatewayv2.GetDomainNamesInput{
			NextToken: prevToken,
		})
		if err != nil {
			if isErr(err, "NotFoundException") {
				return nil, nil
			}
			return nil, err
		}

		for _, domainName := range output.Items {
			resource := apiGatewayV2DomainNameHandle(ctx, domainName)
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
		if isErr(err, "NotFoundException") {
			return nil, nil
		}
		return nil, err
	}

	return values, nil
}
func apiGatewayV2DomainNameHandle(ctx context.Context, domainName typesv2.DomainName) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := fmt.Sprintf("arn:%s:apigateway:%s::/domainnames/%s", describeCtx.Partition, describeCtx.Region, *domainName.DomainName)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *domainName.DomainName,
		Description: model.ApiGatewayV2DomainNameDescription{
			DomainName: domainName,
		},
	}
	return resource
}
func GetApiGatewayV2DomainName(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	domainName := fields["domain_name"]

	client := apigatewayv2.NewFromConfig(cfg)
	out, err := client.GetDomainName(ctx, &apigatewayv2.GetDomainNameInput{
		DomainName: &domainName,
	})
	if err != nil {
		if isErr(err, "NotFoundException") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	values = append(values, apiGatewayV2DomainNameHandle(ctx, typesv2.DomainName{
		DomainName:                    out.DomainName,
		ApiMappingSelectionExpression: out.ApiMappingSelectionExpression,
		DomainNameConfigurations:      out.DomainNameConfigurations,
		MutualTlsAuthentication:       out.MutualTlsAuthentication,
		Tags:                          out.Tags,
	}))
	return values, nil
}

func ApiGatewayV2Integration(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := apigatewayv2.NewFromConfig(cfg)

	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (*string, error) {
		output, err := client.GetApis(ctx, &apigatewayv2.GetApisInput{
			NextToken: prevToken,
		})
		if err != nil {
			if isErr(err, "NotFoundException") {
				return nil, nil
			}
			return nil, err
		}

		err = PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
			for _, api := range output.Items {
				output, err := client.GetIntegrations(ctx, &apigatewayv2.GetIntegrationsInput{
					ApiId:     aws.String(*api.ApiId),
					NextToken: prevToken,
				})

				if err != nil {
					return nil, err
				}

				for _, integration := range output.Items {
					resource := apiGatewayV2IntegrationHandle(ctx, integration, *api.ApiId)
					if stream != nil {
						if err := (*stream)(resource); err != nil {
							return nil, err
						}
					} else {
						values = append(values, resource)
					}
				}
				if err != nil {
					return nil, err
				}
				return output.NextToken, nil
			}
			return output.NextToken, nil
		})

		if err != nil {
			if isErr(err, "NotFoundException") || isErr(err, "TooManyRequestsException") {
				return nil, nil
			}
			return nil, err
		}
		return output.NextToken, nil
	})
	if err != nil {
		if isErr(err, "NotFoundException") {
			return nil, nil
		}
		return nil, err
	}

	return values, nil
}
func apiGatewayV2IntegrationHandle(ctx context.Context, integration typesv2.Integration, apiId string) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := fmt.Sprintf("arn:%s:apigateway:%s::/apis/%s/integrations/%s", describeCtx.Partition, describeCtx.Region, apiId, *integration.IntegrationId)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		ID:     *integration.IntegrationId,
		Description: model.ApiGatewayV2IntegrationDescription{
			Integration: integration,
			ApiId:       apiId,
		},
	}
	return resource
}
func GetApiGatewayV2Integration(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	apiId := fields["api_id"]
	integrationID := fields["id"]

	client := apigatewayv2.NewFromConfig(cfg)

	api, err := client.GetApi(ctx, &apigatewayv2.GetApiInput{
		ApiId: &apiId,
	})
	if err != nil {
		if isErr(err, "NotFoundException") {
			return nil, nil
		}
		return nil, err
	}

	out, err := client.GetIntegration(ctx, &apigatewayv2.GetIntegrationInput{
		ApiId:         aws.String(*api.ApiId),
		IntegrationId: &integrationID,
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	integration := typesv2.Integration{
		ApiGatewayManaged:                      out.ApiGatewayManaged,
		ConnectionId:                           out.ConnectionId,
		ConnectionType:                         out.ConnectionType,
		ContentHandlingStrategy:                out.ContentHandlingStrategy,
		CredentialsArn:                         out.CredentialsArn,
		Description:                            out.Description,
		IntegrationId:                          out.IntegrationId,
		IntegrationMethod:                      out.IntegrationMethod,
		IntegrationResponseSelectionExpression: out.IntegrationResponseSelectionExpression,
		IntegrationSubtype:                     out.IntegrationSubtype,
		IntegrationType:                        out.IntegrationType,
		IntegrationUri:                         out.IntegrationUri,
		PassthroughBehavior:                    out.PassthroughBehavior,
		PayloadFormatVersion:                   out.PayloadFormatVersion,
		RequestParameters:                      out.RequestParameters,
		RequestTemplates:                       out.RequestTemplates,
		ResponseParameters:                     out.ResponseParameters,
		TemplateSelectionExpression:            out.TemplateSelectionExpression,
		TimeoutInMillis:                        out.TimeoutInMillis,
		TlsConfig:                              out.TlsConfig,
	}
	values = append(values, apiGatewayV2IntegrationHandle(ctx, integration, *api.ApiId))
	return values, nil
}

func ApiGatewayDomainName(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := apigateway.NewFromConfig(cfg)
	var values []Resource
	pager := apigateway.NewGetDomainNamesPaginator(client, &apigateway.GetDomainNamesInput{})
	for pager.HasMorePages() {
		output, err := pager.NextPage(ctx)
		if err != nil {
			if isErr(err, "NotFoundException") {
				return nil, nil
			}
			return nil, err
		}

		for _, domainName := range output.Items {
			resource := apiGatewayDomainNameHandle(ctx, domainName)
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
func apiGatewayDomainNameHandle(ctx context.Context, domainName types.DomainName) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := fmt.Sprintf("arn:%s:apigateway:%s::/domainname/%s", describeCtx.Partition, describeCtx.Region, *domainName.DomainName)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *domainName.DomainName,
		Description: model.ApiGatewayDomainNameDescription{
			DomainName: domainName,
		},
	}
	return resource
}
func GetApiGatewayDomainName(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	domainName := fields["domain_name"]

	client := apigateway.NewFromConfig(cfg)
	out, err := client.GetDomainName(ctx, &apigateway.GetDomainNameInput{
		DomainName: &domainName,
	})
	if err != nil {
		if isErr(err, "NotFoundException") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	values = append(values, apiGatewayDomainNameHandle(ctx, types.DomainName{
		DomainName:                          out.DomainName,
		CertificateName:                     out.CertificateName,
		CertificateArn:                      out.CertificateArn,
		CertificateUploadDate:               out.CertificateUploadDate,
		DistributionDomainName:              out.DistributionDomainName,
		DistributionHostedZoneId:            out.DistributionHostedZoneId,
		DomainNameStatus:                    out.DomainNameStatus,
		DomainNameStatusMessage:             out.DomainNameStatusMessage,
		OwnershipVerificationCertificateArn: out.OwnershipVerificationCertificateArn,
		RegionalCertificateName:             out.RegionalCertificateName,
		RegionalCertificateArn:              out.RegionalCertificateArn,
		RegionalDomainName:                  out.RegionalDomainName,
		RegionalHostedZoneId:                out.RegionalHostedZoneId,
		SecurityPolicy:                      out.SecurityPolicy,
		EndpointConfiguration:               out.EndpointConfiguration,
		MutualTlsAuthentication:             out.MutualTlsAuthentication,
		Tags:                                out.Tags,
	}))
	return values, nil
}

func ApiGatewayV2Route(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := apigatewayv2.NewFromConfig(cfg)

	apis, err := ApiGatewayV2API(ctx, cfg, nil)
	if err != nil {
		return nil, err
	}
	var values []Resource
	for _, a := range apis {
		api := a.Description.(model.ApiGatewayV2APIDescription).API
		output, err := client.GetRoutes(ctx, &apigatewayv2.GetRoutesInput{
			ApiId: api.ApiId,
		})
		if err != nil {
			if isErr(err, "NotFoundException") {
				return nil, nil
			}
			return nil, err
		}

		for _, route := range output.Items {
			resource := apiGatewayV2Route(ctx, route)
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
func apiGatewayV2Route(ctx context.Context, route typesv2.Route) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := fmt.Sprintf("arn:%s:apigateway:%s::/apis/%s/routes/%s", describeCtx.Partition, describeCtx.Region, *route.RouteId)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		Name:   *route.RouteId,
		Description: model.ApiGatewayV2RouteDescription{
			Route: route,
		},
	}
	return resource
}
