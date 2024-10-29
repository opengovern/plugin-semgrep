package describer

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/securityhub"
	"github.com/aws/aws-sdk-go-v2/service/securityhub/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"github.com/opengovern/og-aws-describer/aws/model"
	"strconv"
	"strings"
)

func SecurityHubHub(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := securityhub.NewFromConfig(cfg)
	out, err := client.DescribeHub(ctx, &securityhub.DescribeHubInput{})
	if err != nil {
		if isErr(err, "InvalidAccessException") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource

	resource, err := securityHubHubHandle(ctx, cfg, out)
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

	return values, nil
}
func securityHubHubHandle(ctx context.Context, cfg aws.Config, out *securityhub.DescribeHubOutput) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := securityhub.NewFromConfig(cfg)

	tags, err := client.ListTagsForResource(ctx, &securityhub.ListTagsForResourceInput{ResourceArn: out.HubArn})
	if err != nil {
		if isErr(err, "InvalidAccessException") {
			return Resource{}, nil
		}
		return Resource{}, err
	}

	data, err := client.GetAdministratorAccount(ctx, &securityhub.GetAdministratorAccountInput{})
	if err != nil {
		return Resource{}, err
	}

	desc := model.SecurityHubHubDescription{
		Hub:  out,
		Tags: tags.Tags,
	}
	if data.Administrator != nil {
		desc.AdministratorAccount = *data.Administrator
	}
	resource := Resource{
		Region:      describeCtx.KaytuRegion,
		Description: desc,
	}
	if out.HubArn != nil {
		resource.ARN = *out.HubArn
		resource.Name = nameFromArn(*out.HubArn)
	}
	return resource, nil
}
func GetSecurityHubHub(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	arn := fields["arn"]
	var values []Resource
	client := securityhub.NewFromConfig(cfg)

	out, err := client.DescribeHub(ctx, &securityhub.DescribeHubInput{
		HubArn: &arn,
	})
	if err != nil {
		if isErr(err, "DescribeHubNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	resource, err := securityHubHubHandle(ctx, cfg, out)
	if err != nil {
		return nil, err
	}
	values = append(values, resource)
	return values, nil
}

func SecurityHubActionTarget(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := securityhub.NewFromConfig(cfg)

	var values []Resource
	paginator := securityhub.NewDescribeActionTargetsPaginator(client, &securityhub.DescribeActionTargetsInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if isErr(err, "InvalidAccessException") {
				return nil, nil
			}
			return nil, err
		}

		for _, actionTarget := range page.ActionTargets {
			resource := securityHubActionTargetHandle(ctx, actionTarget)

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
func securityHubActionTargetHandle(ctx context.Context, actionTarget types.ActionTarget) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *actionTarget.ActionTargetArn,
		Name:   *actionTarget.Name,
		Description: model.SecurityHubActionTargetDescription{
			ActionTarget: actionTarget,
		},
	}
	return resource
}
func GetSecurityHubActionTarget(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	actionTargetArn := fields["arn"]
	client := securityhub.NewFromConfig(cfg)

	out, err := client.DescribeActionTargets(ctx, &securityhub.DescribeActionTargetsInput{
		ActionTargetArns: []string{actionTargetArn},
	})
	if err != nil {
		if isErr(err, "DescribeActionTargetsNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	for _, actionTarget := range out.ActionTargets {

		values = append(values, securityHubActionTargetHandle(ctx, actionTarget))

	}
	return values, nil
}

func SecurityHubFinding(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := securityhub.NewFromConfig(cfg)

	var values []Resource
	paginator := securityhub.NewGetFindingsPaginator(client, &securityhub.GetFindingsInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if isErr(err, "InvalidAccessException") {
				return nil, nil
			}
			return nil, err
		}

		for _, finding := range page.Findings {
			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ID:     *finding.Id,
				Name:   *finding.Title,
				Description: model.SecurityHubFindingDescription{
					Finding: finding,
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

func SecurityHubFindingAggregator(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := securityhub.NewFromConfig(cfg)

	var values []Resource
	paginator := securityhub.NewListFindingAggregatorsPaginator(client, &securityhub.ListFindingAggregatorsInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if isErr(err, "InvalidAccessException") {
				return nil, nil
			}
			return nil, err
		}

		for _, findingAggregatorSummary := range page.FindingAggregators {
			resource, err := securityHubFindingAggregatorHandle(ctx, cfg, *findingAggregatorSummary.FindingAggregatorArn)
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

	return values, nil
}
func securityHubFindingAggregatorHandle(ctx context.Context, cfg aws.Config, findingAggregatorArn string) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := securityhub.NewFromConfig(cfg)

	findingAggregator, err := client.GetFindingAggregator(ctx, &securityhub.GetFindingAggregatorInput{
		FindingAggregatorArn: &findingAggregatorArn,
	})
	if err != nil {
		if isErr(err, "InvalidAccessException") {
			return Resource{}, nil
		}
		return Resource{}, err
	}
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *findingAggregator.FindingAggregatorArn,
		Description: model.SecurityHubFindingAggregatorDescription{
			FindingAggregator: *findingAggregator,
		},
	}
	return resource, nil
}
func GetSecurityHubFindingAggregator(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	arn := fields["arn"]
	var values []Resource

	resource, err := securityHubFindingAggregatorHandle(ctx, cfg, arn)
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

func SecurityHubInsight(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := securityhub.NewFromConfig(cfg)

	var values []Resource
	paginator := securityhub.NewGetInsightsPaginator(client, &securityhub.GetInsightsInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if isErr(err, "InvalidAccessException") {
				return nil, nil
			}
			return nil, err
		}

		for _, insight := range page.Insights {
			resource := securityHubInsightHandle(ctx, insight)

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
func securityHubInsightHandle(ctx context.Context, insight types.Insight) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *insight.InsightArn,
		Name:   *insight.Name,
		Description: model.SecurityHubInsightDescription{
			Insight: insight,
		},
	}
	return resource
}
func GetSecurityHubInsight(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	arn := fields["arn"]
	client := securityhub.NewFromConfig(cfg)

	out, err := client.GetInsights(ctx, &securityhub.GetInsightsInput{
		InsightArns: []string{arn},
	})
	if err != nil {
		if isErr(err, "GetInsightsNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	for _, insight := range out.Insights {
		resource := securityHubInsightHandle(ctx, insight)
		values = append(values, resource)
	}
	return values, nil
}

func SecurityHubMember(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	params := &securityhub.ListMembersInput{
		OnlyAssociated: aws.Bool(false),
	}
	myMiddleware := middleware.SerializeMiddlewareFunc(
		"AssociatedMembers",
		func(ctx context.Context, input middleware.SerializeInput, next middleware.SerializeHandler) (
			output middleware.SerializeOutput,
			metadata middleware.Metadata,
			err error) {
			req, ok := input.Request.(*smithyhttp.Request)
			if !ok {
				return output, metadata, fmt.Errorf("unexpected transport: %T", input.Request)
			}

			params, ok = input.Parameters.(*securityhub.ListMembersInput)
			if !ok {
				return output, metadata, fmt.Errorf("unexpected input type: %T", input.Parameters)
			}

			query := req.URL.Query()
			query.Set("OnlyAssociated", strconv.FormatBool(false))
			req.URL.RawQuery = query.Encode()
			return next.HandleSerialize(ctx, input)
		},
	)
	client := securityhub.NewFromConfig(cfg, func(options *securityhub.Options) {
		options.APIOptions = append(options.APIOptions, func(stack *middleware.Stack) error {
			return stack.Serialize.Insert(myMiddleware, "OperationSerializer", middleware.After)
		})
	})

	var values []Resource
	paginator := securityhub.NewListMembersPaginator(client, params, func(o *securityhub.ListMembersPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if isErr(err, "InvalidAccessException") || isErr(err, "InvalidInputException") {
				return nil, nil
			}
			if strings.Contains(err.Error(), "no such resource found") {
				return nil, nil
			}
			return nil, err
		}

		for _, member := range page.Members {
			resource := securityHubMemberHandle(ctx, member)
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
func securityHubMemberHandle(ctx context.Context, member types.Member) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		Name:   *member.AccountId,
		Description: model.SecurityHubMemberDescription{
			Member: member,
		},
	}
	return resource
}
func GetSecurityHubMember(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	accountId := fields["accountId"]
	client := securityhub.NewFromConfig(cfg)

	out, err := client.GetMembers(ctx, &securityhub.GetMembersInput{
		AccountIds: []string{accountId},
	})
	if err != nil {
		if isErr(err, "GetMembersNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	for _, member := range out.Members {
		resource := securityHubMemberHandle(ctx, member)
		values = append(values, resource)
	}
	return values, nil
}

func SecurityHubProduct(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := securityhub.NewFromConfig(cfg)

	var values []Resource
	paginator := securityhub.NewDescribeProductsPaginator(client, &securityhub.DescribeProductsInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if isErr(err, "InvalidAccessException") {
				return nil, nil
			}
			return nil, err
		}

		for _, product := range page.Products {
			resource := Resource{
				Region: describeCtx.KaytuRegion,
				Name:   *product.ProductName,
				ARN:    *product.ProductArn,
				Description: model.SecurityHubProductDescription{
					Product: product,
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

func SecurityHubStandardsControl(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := securityhub.NewFromConfig(cfg)

	var values []Resource

	subPaginator := securityhub.NewGetEnabledStandardsPaginator(client, &securityhub.GetEnabledStandardsInput{})
	for subPaginator.HasMorePages() {
		subPage, err := subPaginator.NextPage(ctx)
		if err != nil {
			if isErr(err, "InvalidAccessException") {
				return nil, nil
			}
			return nil, err
		}

		for _, standard := range subPage.StandardsSubscriptions {
			paginator := securityhub.NewDescribeStandardsControlsPaginator(client, &securityhub.DescribeStandardsControlsInput{
				StandardsSubscriptionArn: standard.StandardsArn,
			})
			for paginator.HasMorePages() {
				page, err := paginator.NextPage(ctx)
				if err != nil {
					if isErr(err, "InvalidAccessException") || isErr(err, "InvalidInputException") {
						return nil, nil
					}
					return nil, err
				}

				for _, standardsControl := range page.Controls {
					resource := securityHubStandardsControlHandle(ctx, standardsControl)
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
	}
	return values, nil
}
func securityHubStandardsControlHandle(ctx context.Context, standardsControl types.StandardsControl) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ID:     *standardsControl.ControlId,
		Name:   *standardsControl.Title,
		ARN:    *standardsControl.StandardsControlArn,
		Description: model.SecurityHubStandardsControlDescription{
			StandardsControl: standardsControl,
		},
	}
	return resource
}
func GetSecurityHubStandardsControl(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	standardsSubscriptionArn := fields["arn"]
	client := securityhub.NewFromConfig(cfg)
	out, err := client.DescribeStandardsControls(ctx, &securityhub.DescribeStandardsControlsInput{
		StandardsSubscriptionArn: &standardsSubscriptionArn,
	})
	if err != nil {
		return nil, err
	}
	var values []Resource
	for _, v := range out.Controls {
		resource := securityHubStandardsControlHandle(ctx, v)
		values = append(values, resource)
	}
	return values, nil
}

func SecurityHubStandardsSubscription(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := securityhub.NewFromConfig(cfg)

	var values []Resource
	standardsPaginator := securityhub.NewDescribeStandardsPaginator(client, &securityhub.DescribeStandardsInput{})
	standards := make(map[string]types.Standard)
	for standardsPaginator.HasMorePages() {
		standardsPage, err := standardsPaginator.NextPage(ctx)
		if err != nil {
			if isErr(err, "InvalidAccessException") {
				return nil, nil
			}
			return nil, err
		}
		for _, standard := range standardsPage.Standards {
			standards[*standard.StandardsArn] = standard
		}
	}

	paginator := securityhub.NewGetEnabledStandardsPaginator(client, &securityhub.GetEnabledStandardsInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if isErr(err, "InvalidAccessException") {
				return nil, nil
			}
			return nil, err
		}

		for _, standardSub := range page.StandardsSubscriptions {
			resource := securityHubStandardsSubscriptionHandle(ctx, standardSub, standards)
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
func securityHubStandardsSubscriptionHandle(ctx context.Context, standardSub types.StandardsSubscription, standards map[string]types.Standard) Resource {
	describeCtx := GetDescribeContext(ctx)

	standard, _ := standards[*standardSub.StandardsArn]
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *standardSub.StandardsSubscriptionArn,
		Description: model.SecurityHubStandardsSubscriptionDescription{
			Standard:              standard,
			StandardsSubscription: standardSub,
		},
	}
	return resource
}
func GetSecurityHubStandardsSubscription(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	standardsSubscriptionArn := fields["standardsSubscriptionArn"]
	client := securityhub.NewFromConfig(cfg)

	var values []Resource
	standardsPaginator := securityhub.NewDescribeStandardsPaginator(client, &securityhub.DescribeStandardsInput{})
	standards := make(map[string]types.Standard)
	for standardsPaginator.HasMorePages() {
		standardsPage, err := standardsPaginator.NextPage(ctx)
		if err != nil {
			if isErr(err, "InvalidAccessException") {
				return nil, nil
			}
			return nil, err
		}
		for _, standard := range standardsPage.Standards {
			standards[*standard.StandardsArn] = standard
		}
	}
	out, err := client.GetEnabledStandards(ctx, &securityhub.GetEnabledStandardsInput{
		StandardsSubscriptionArns: []string{standardsSubscriptionArn},
	})
	if err != nil {
		return nil, err
	}

	for _, standardSub := range out.StandardsSubscriptions {
		resource := securityHubStandardsSubscriptionHandle(ctx, standardSub, standards)
		values = append(values, resource)
	}
	return values, nil
}
