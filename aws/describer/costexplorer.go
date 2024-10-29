package describer

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/opengovern/og-aws-describer/aws/model"
	"github.com/opengovern/og-util/pkg/describe/enums"
	pointerUtil "github.com/opengovern/og-util/pkg/pointer"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func getTrimmedStringPointer(str *string) *string {
	if str == nil {
		return nil
	}
	return aws.String(strings.Trim(strings.TrimSpace(*str), "\x00"))
}

func getParsedFloat64Pointer(str *string) *float64 {
	str = getTrimmedStringPointer(str)
	if str == nil {
		return nil
	}
	amount, err := strconv.ParseFloat(*str, 64)
	if err != nil {
		return nil
	}
	return aws.Float64(amount)
}

func setRowMetrics(row *model.CostExplorerRow, metrics map[string]types.MetricValue) {
	if _, ok := metrics["BlendedCost"]; ok {
		row.BlendedCostAmount = getParsedFloat64Pointer(metrics["BlendedCost"].Amount)
		row.BlendedCostUnit = getTrimmedStringPointer(metrics["BlendedCost"].Unit)
	}
	if _, ok := metrics["UnblendedCost"]; ok {
		row.UnblendedCostAmount = getParsedFloat64Pointer(metrics["UnblendedCost"].Amount)
		row.UnblendedCostUnit = getTrimmedStringPointer(metrics["UnblendedCost"].Unit)
	}
	if _, ok := metrics["NetUnblendedCost"]; ok {
		row.NetUnblendedCostAmount = getParsedFloat64Pointer(metrics["NetUnblendedCost"].Amount)
		row.NetUnblendedCostUnit = getTrimmedStringPointer(metrics["NetUnblendedCost"].Unit)
	}
	if _, ok := metrics["AmortizedCost"]; ok {
		row.AmortizedCostAmount = getParsedFloat64Pointer(metrics["AmortizedCost"].Amount)
		row.AmortizedCostUnit = getTrimmedStringPointer(metrics["AmortizedCost"].Unit)
	}
	if _, ok := metrics["NetAmortizedCost"]; ok {
		row.NetAmortizedCostAmount = getParsedFloat64Pointer(metrics["NetAmortizedCost"].Amount)
		row.NetAmortizedCostUnit = getTrimmedStringPointer(metrics["NetAmortizedCost"].Unit)
	}
	if _, ok := metrics["UsageQuantity"]; ok {
		row.UsageQuantityAmount = getParsedFloat64Pointer(metrics["UsageQuantity"].Amount)
		row.UsageQuantityUnit = getTrimmedStringPointer(metrics["UsageQuantity"].Unit)
	}
	if _, ok := metrics["NormalizedUsageAmount"]; ok {
		row.NormalizedUsageAmount = getParsedFloat64Pointer(metrics["NormalizedUsageAmount"].Amount)
		row.NormalizedUsageUnit = getTrimmedStringPointer(metrics["NormalizedUsageAmount"].Unit)
	}
}

func costMonthly(ctx context.Context, cfg aws.Config, by string, startDate, endDate time.Time) ([]model.CostExplorerRow, error) {
	describeCtx := GetDescribeContext(ctx)
	timeFormat := "2006-01-02"
	endTime := endDate.Format(timeFormat)
	startTime := startDate.Format(timeFormat)

	params := &costexplorer.GetCostAndUsageInput{
		Filter: &types.Expression{
			Dimensions: &types.DimensionValues{
				Key:    types.DimensionLinkedAccount,
				Values: []string{describeCtx.AccountID},
			},
		},
		TimePeriod: &types.DateInterval{
			Start: aws.String(startTime),
			End:   aws.String(endTime),
		},
		Granularity: types.GranularityMonthly,
		Metrics: []string{
			"BlendedCost",
			"UnblendedCost",
			"NetUnblendedCost",
			"AmortizedCost",
			"NetAmortizedCost",
			"UsageQuantity",
			"NormalizedUsageAmount",
		},
		GroupBy: []types.GroupDefinition{
			{
				Type: types.GroupDefinitionTypeDimension,
				Key:  aws.String(by),
			},
			{
				Type: types.GroupDefinitionTypeDimension,
				Key:  aws.String(string(types.DimensionBillingEntity)),
			},
		},
	}

	client := costexplorer.NewFromConfig(cfg)

	var values []model.CostExplorerRow
	for {
		out, err := client.GetCostAndUsage(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, result := range out.ResultsByTime {

			// If there are no groupings, create a row from the totals
			if len(result.Groups) == 0 {
				var row model.CostExplorerRow

				row.Estimated = result.Estimated
				row.PeriodStart = result.TimePeriod.Start
				row.PeriodEnd = result.TimePeriod.End

				setRowMetrics(&row, result.Total)
				values = append(values, row)
			}
			// make a row per group
			for _, group := range result.Groups {
				var row model.CostExplorerRow

				row.Estimated = result.Estimated
				row.PeriodStart = result.TimePeriod.Start
				row.PeriodEnd = result.TimePeriod.End

				if len(group.Keys) > 0 {
					row.Dimension1 = aws.String(group.Keys[0])
					if len(group.Keys) > 1 {
						row.Dimension2 = aws.String(group.Keys[1])
					}
				}
				setRowMetrics(&row, group.Metrics)

				values = append(values, row)
			}
		}

		if out.NextPageToken == nil {
			break
		}

		params.NextPageToken = out.NextPageToken
	}

	return values, nil
}

func CostByServiceLastMonth(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	triggerType := GetTriggerTypeFromContext(ctx)
	startDate := time.Now().AddDate(0, -1, 0)
	if triggerType == enums.DescribeTriggerTypeInitialDiscovery {
		startDate = time.Now().AddDate(0, -3, -7)
	}
	costs, err := costMonthly(ctx, cfg, "SERVICE", startDate, time.Now())
	if err != nil {
		return nil, err
	}
	var values []Resource
	for _, cost := range costs {
		if cost.Dimension1 == nil {
			continue
		}
		resource := Resource{
			Region:      describeCtx.KaytuRegion,
			ID:          "service-" + *cost.Dimension1 + "-cost-monthly",
			Description: model.CostExplorerByServiceMonthlyDescription{CostExplorerRow: cost},
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

func CostByAccountLastMonth(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	triggerType := GetTriggerTypeFromContext(ctx)
	startDate := time.Now().AddDate(0, -1, 0)
	if triggerType == enums.DescribeTriggerTypeInitialDiscovery {
		startDate = time.Now().AddDate(0, -3, -7)
	}

	costs, err := costMonthly(ctx, cfg, "LINKED_ACCOUNT", startDate, time.Now())
	if err != nil {
		return nil, err
	}
	var values []Resource
	for _, cost := range costs {
		if cost.Dimension1 == nil {
			continue
		}
		resource := Resource{
			Region:      describeCtx.KaytuRegion,
			ID:          "account-" + *cost.Dimension1 + "-cost-monthly",
			Description: model.CostExplorerByAccountMonthlyDescription{CostExplorerRow: cost},
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

func costDaily(ctx context.Context, cfg aws.Config, by string, startDate, endDate time.Time) ([]model.CostExplorerRow, error) {
	describeCtx := GetDescribeContext(ctx)
	timeFormat := "2006-01-02"
	endTime := endDate.Format(timeFormat)
	startTime := startDate.Format(timeFormat)

	params := &costexplorer.GetCostAndUsageInput{
		Filter: &types.Expression{
			Dimensions: &types.DimensionValues{
				Key:    types.DimensionLinkedAccount,
				Values: []string{describeCtx.AccountID},
			},
		},
		TimePeriod: &types.DateInterval{
			Start: aws.String(startTime),
			End:   aws.String(endTime),
		},
		Granularity: types.GranularityDaily,
		Metrics: []string{
			"BlendedCost",
			"UnblendedCost",
			"NetUnblendedCost",
			"AmortizedCost",
			"NetAmortizedCost",
			"UsageQuantity",
			"NormalizedUsageAmount",
		},
		GroupBy: []types.GroupDefinition{
			{
				Type: types.GroupDefinitionTypeDimension,
				Key:  aws.String(by),
			},
			{
				Type: types.GroupDefinitionTypeDimension,
				Key:  aws.String(string(types.DimensionBillingEntity)),
			},
		},
	}

	client := costexplorer.NewFromConfig(cfg)
	var values []model.CostExplorerRow
	for {
		out, err := client.GetCostAndUsage(ctx, params)
		if err != nil {
			if isErr(err, "AccessDeniedException") {
				break
			} else {
				return nil, err
			}
		}

		for _, result := range out.ResultsByTime {

			// If there are no groupings, create a row from the totals
			if len(result.Groups) == 0 {
				var row model.CostExplorerRow

				row.Estimated = result.Estimated
				row.PeriodStart = result.TimePeriod.Start
				row.PeriodEnd = result.TimePeriod.End

				setRowMetrics(&row, result.Total)
				values = append(values, row)
			}
			// make a row per group
			for _, group := range result.Groups {
				var row model.CostExplorerRow

				row.Estimated = result.Estimated
				row.PeriodStart = result.TimePeriod.Start
				row.PeriodEnd = result.TimePeriod.End

				if len(group.Keys) > 0 {
					row.Dimension1 = aws.String(group.Keys[0])
					if len(group.Keys) > 1 {
						row.Dimension2 = aws.String(group.Keys[1])
					}
				}
				setRowMetrics(&row, group.Metrics)

				values = append(values, row)
			}
		}

		if out.NextPageToken == nil {
			break
		}

		params.NextPageToken = out.NextPageToken
	}

	return values, nil
}

func getEc2OtherCostKeyFromDimension(dimension string) string {
	switch {
	case strings.Contains(dimension, "EBSOptimized"):
		return "EC2 - EBSOptimized"
	case strings.Contains(dimension, "CPUCredits"):
		return "EC2 - CPUCredits"
	case strings.Contains(dimension, "DataTransfer"):
		return "EC2 - DataTransfer"
	case strings.Contains(dimension, "AWS-In-Bytes"):
		return "EC2 - AWS In"
	case strings.Contains(dimension, "AWS-Out-Bytes"):
		return "EC2 - AWS Out"
	case strings.Contains(dimension, "ElasticIP"):
		return "EC2 - ElasticIP"
	case strings.Contains(dimension, "NatGateway"):
		return "EC2 - NatGateway"
	case strings.Contains(dimension, "EBS:Snapshot"):
		return "EC2 - EBS Snapshot"
	case strings.Contains(dimension, "EBS"):
		return "EC2 - EBS"
	default:
		return "EC2 - Other"
	}
}

func ec2OtherCostDaily(ctx context.Context, cfg aws.Config, startDate, endDate time.Time) ([]model.CostExplorerRow, error) {
	describeCtx := GetDescribeContext(ctx)
	timeFormat := "2006-01-02"
	endTime := endDate.Format(timeFormat)
	startTime := startDate.Format(timeFormat)

	params := &costexplorer.GetCostAndUsageInput{
		Filter: &types.Expression{
			And: []types.Expression{
				{
					Dimensions: &types.DimensionValues{
						Key:    types.DimensionService,
						Values: []string{"EC2 - Other"},
					},
				},
				{
					Dimensions: &types.DimensionValues{
						Key:    types.DimensionLinkedAccount,
						Values: []string{describeCtx.AccountID},
					},
				},
			},
		},
		TimePeriod: &types.DateInterval{
			Start: aws.String(startTime),
			End:   aws.String(endTime),
		},
		Granularity: types.GranularityDaily,
		Metrics: []string{
			"BlendedCost",
			"UnblendedCost",
			"NetUnblendedCost",
			"AmortizedCost",
			"NetAmortizedCost",
			"UsageQuantity",
			"NormalizedUsageAmount",
		},
		GroupBy: []types.GroupDefinition{
			{
				Type: types.GroupDefinitionTypeDimension,
				Key:  aws.String(string(types.DimensionUsageType)),
			},
		},
	}

	client := costexplorer.NewFromConfig(cfg)
	var values []model.CostExplorerRow
	for {
		out, err := client.GetCostAndUsage(ctx, params)
		if err != nil {
			if isErr(err, "AccessDeniedException") {
				break
			} else {
				return nil, err
			}
		}
		for _, result := range out.ResultsByTime {
			valuesMap := make(map[string]model.CostExplorerRow)
			// If there are no groupings, create a row from the totals
			if len(result.Groups) == 0 {
				var row model.CostExplorerRow

				row.Estimated = result.Estimated
				row.PeriodStart = result.TimePeriod.Start
				row.PeriodEnd = result.TimePeriod.End

				setRowMetrics(&row, result.Total)
				values = append(values, row)
			}
			// make a row per group
			for _, group := range result.Groups {
				var row model.CostExplorerRow

				row.Estimated = result.Estimated
				row.PeriodStart = result.TimePeriod.Start
				row.PeriodEnd = result.TimePeriod.End

				if len(group.Keys) > 0 {
					row.Dimension1 = aws.String(getEc2OtherCostKeyFromDimension(group.Keys[0]))
					if len(group.Keys) > 1 {
						row.Dimension2 = aws.String(group.Keys[1])
					}
				} else {
					continue
				}
				setRowMetrics(&row, group.Metrics)

				if v, ok := valuesMap[*row.Dimension1]; ok {
					v.BlendedCostAmount = pointerUtil.PAdd(v.BlendedCostAmount, row.BlendedCostAmount)
					v.UnblendedCostAmount = pointerUtil.PAdd(v.UnblendedCostAmount, row.UnblendedCostAmount)
					v.NetUnblendedCostAmount = pointerUtil.PAdd(v.NetUnblendedCostAmount, row.NetUnblendedCostAmount)
					v.AmortizedCostAmount = pointerUtil.PAdd(v.AmortizedCostAmount, row.AmortizedCostAmount)
					v.NetAmortizedCostAmount = pointerUtil.PAdd(v.NetAmortizedCostAmount, row.NetAmortizedCostAmount)
					v.UsageQuantityAmount = pointerUtil.PAdd(v.UsageQuantityAmount, row.UsageQuantityAmount)
					v.NormalizedUsageAmount = pointerUtil.PAdd(v.NormalizedUsageAmount, row.NormalizedUsageAmount)
					valuesMap[*row.Dimension1] = v
				} else {
					valuesMap[*row.Dimension1] = row
				}
			}
			for _, v := range valuesMap {
				values = append(values, v)
			}
		}

		if out.NextPageToken == nil {
			break
		}

		params.NextPageToken = out.NextPageToken
	}

	return values, nil
}

func CostByServiceLastDay(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	triggerType := GetTriggerTypeFromContext(ctx)
	startDate := time.Now().AddDate(0, 0, -7)
	if time.Now().Day() == 6 {
		y, m, _ := time.Now().Date()
		startDate = time.Date(y, m, 1, 0, 0, 0, 0, time.UTC).AddDate(0, -1, 0)
	}

	if triggerType == enums.DescribeTriggerTypeInitialDiscovery {
		startDate = time.Now().AddDate(0, -1, -7)
	} else if triggerType == enums.DescribeTriggerTypeCostFullDiscovery {
		startDate = time.Date(2022, 12, 1, 0, 0, 0, 0, time.UTC)
	}
	endDate := time.Now()

	costs, err := costDaily(ctx, cfg, "SERVICE", startDate, endDate)
	if err != nil {
		return nil, err
	}
	var values []Resource
	for _, cost := range costs {
		if cost.Dimension1 == nil || *cost.Dimension1 == "EC2 - Other" {
			continue
		}

		tStart, err := time.Parse("2006-01-02", *cost.PeriodStart)
		if err != nil {
			return nil, err
		}
		tEnd, err := time.Parse("2006-01-02", *cost.PeriodEnd)
		if err != nil {
			return nil, err
		}

		diff := tEnd.Sub(tStart) / 2
		costDate := tStart.Add(diff)

		resource := Resource{
			Region:      describeCtx.KaytuRegion,
			ID:          "service-" + *cost.Dimension1 + "-cost-" + *cost.PeriodEnd,
			Description: model.CostExplorerByServiceDailyDescription{CostExplorerRow: cost, CostDateMillis: costDate.UnixMilli()},
		}
		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return nil, err
			}
		} else {
			values = append(values, resource)
		}
	}

	ec2OtherCosts, err := ec2OtherCostDaily(ctx, cfg, startDate, endDate)
	if err != nil {
		return nil, err
	}
	for _, cost := range ec2OtherCosts {
		if cost.Dimension1 == nil {
			continue
		}

		tStart, err := time.Parse("2006-01-02", *cost.PeriodStart)
		if err != nil {
			return nil, err
		}
		tEnd, err := time.Parse("2006-01-02", *cost.PeriodEnd)
		if err != nil {
			return nil, err
		}

		diff := tEnd.Sub(tStart) / 2
		costDate := tStart.Add(diff)

		resource := Resource{
			Region:      describeCtx.KaytuRegion,
			ID:          "service-" + *cost.Dimension1 + "-cost-" + *cost.PeriodEnd,
			Description: model.CostExplorerByServiceDailyDescription{CostExplorerRow: cost, CostDateMillis: costDate.UnixMilli()},
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

func CostByAccountLastDay(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	triggerType := GetTriggerTypeFromContext(ctx)
	startDate := time.Now().AddDate(0, 0, -7)
	if triggerType == enums.DescribeTriggerTypeInitialDiscovery {
		startDate = time.Now().AddDate(0, -3, -7)
	}
	endDate := time.Now()

	costs, err := costDaily(ctx, cfg, "LINKED_ACCOUNT", startDate, endDate)
	if err != nil {
		return nil, err
	}
	var values []Resource
	for _, cost := range costs {
		if cost.Dimension1 == nil {
			continue
		}
		resource := Resource{
			Region:      describeCtx.KaytuRegion,
			ID:          "account-" + *cost.Dimension1 + "-cost-" + *cost.PeriodEnd,
			Description: model.CostExplorerByAccountDailyDescription{CostExplorerRow: cost},
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

func buildCostByRecordTypeInput(granularity string) *costexplorer.GetCostAndUsageInput {
	timeFormat := "2006-01-02"
	if granularity == "HOURLY" {
		timeFormat = "2006-01-02T15:04:05Z"
	}
	endTime := time.Now().Format(timeFormat)
	startTime := time.Now().AddDate(0, -1, 0).Format(timeFormat)

	params := &costexplorer.GetCostAndUsageInput{
		TimePeriod: &types.DateInterval{
			Start: aws.String(startTime),
			End:   aws.String(endTime),
		},
		Granularity: types.Granularity(granularity),
		Metrics: []string{
			"BlendedCost",
			"UnblendedCost",
			"NetUnblendedCost",
			"AmortizedCost",
			"NetAmortizedCost",
			"UsageQuantity",
			"NormalizedUsageAmount",
		},
		GroupBy: []types.GroupDefinition{
			{
				Type: types.GroupDefinitionType("DIMENSION"),
				Key:  aws.String("LINKED_ACCOUNT"),
			},
			{
				Type: types.GroupDefinitionType("DIMENSION"),
				Key:  aws.String("RECORD_TYPE"),
			},
		},
	}

	return params
}

func CostByRecordTypeLastMonth(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := costexplorer.NewFromConfig(cfg)

	params := buildCostByRecordTypeInput("MONTHLY")

	out, err := client.GetCostAndUsage(ctx, params)
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, result := range out.ResultsByTime {
		for _, group := range result.Groups {
			var row model.CostExplorerRow

			row.Estimated = result.Estimated
			row.PeriodStart = result.TimePeriod.Start
			row.PeriodEnd = result.TimePeriod.End

			if len(group.Keys) > 0 {
				row.Dimension1 = aws.String(group.Keys[0])
				if len(group.Keys) > 1 {
					row.Dimension2 = aws.String(group.Keys[1])
				}
			}
			setRowMetrics(&row, group.Metrics)

			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ID:          "account-" + *row.Dimension1 + "-" + *row.Dimension2 + "-cost-monthly",
				Description: model.CostExplorerByRecordTypeMonthlyDescription{CostExplorerRow: row},
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

func CostByRecordTypeLastDay(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := costexplorer.NewFromConfig(cfg)

	params := buildCostByRecordTypeInput("DAILY")

	out, err := client.GetCostAndUsage(ctx, params)
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, result := range out.ResultsByTime {
		for _, group := range result.Groups {
			var row model.CostExplorerRow

			row.Estimated = result.Estimated
			row.PeriodStart = result.TimePeriod.Start
			row.PeriodEnd = result.TimePeriod.End

			if len(group.Keys) > 0 {
				row.Dimension1 = aws.String(group.Keys[0])
				if len(group.Keys) > 1 {
					row.Dimension2 = aws.String(group.Keys[1])
				}
			}
			setRowMetrics(&row, group.Metrics)

			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ID:          "account-" + *row.Dimension1 + "-" + *row.Dimension2 + "-cost-" + *row.PeriodEnd,
				Description: model.CostExplorerByRecordTypeDailyDescription{CostExplorerRow: row},
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

func buildCostByServiceAndUsageInput(granularity string) *costexplorer.GetCostAndUsageInput {
	timeFormat := "2006-01-02"
	if granularity == "HOURLY" {
		timeFormat = "2006-01-02T15:04:05Z"
	}
	endTime := time.Now().Format(timeFormat)
	startTime := time.Now().AddDate(0, -1, 0).Format(timeFormat)

	params := &costexplorer.GetCostAndUsageInput{
		TimePeriod: &types.DateInterval{
			Start: aws.String(startTime),
			End:   aws.String(endTime),
		},
		Granularity: types.Granularity(granularity),
		Metrics: []string{
			"BlendedCost",
			"UnblendedCost",
			"NetUnblendedCost",
			"AmortizedCost",
			"NetAmortizedCost",
			"UsageQuantity",
			"NormalizedUsageAmount",
		},
		GroupBy: []types.GroupDefinition{
			{
				Type: types.GroupDefinitionType("DIMENSION"),
				Key:  aws.String("SERVICE"),
			},
			{
				Type: types.GroupDefinitionType("DIMENSION"),
				Key:  aws.String("USAGE_TYPE"),
			},
		},
	}

	return params
}

func CostByServiceUsageLastMonth(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := costexplorer.NewFromConfig(cfg)

	params := buildCostByServiceAndUsageInput("MONTHLY")

	out, err := client.GetCostAndUsage(ctx, params)
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, result := range out.ResultsByTime {
		for _, group := range result.Groups {
			var row model.CostExplorerRow

			row.Estimated = result.Estimated
			row.PeriodStart = result.TimePeriod.Start
			row.PeriodEnd = result.TimePeriod.End

			if len(group.Keys) > 0 {
				row.Dimension1 = aws.String(group.Keys[0])
				if len(group.Keys) > 1 {
					row.Dimension2 = aws.String(group.Keys[1])
				}
			}
			setRowMetrics(&row, group.Metrics)

			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ID:          "service-" + *row.Dimension1 + "-" + *row.Dimension2 + "-cost-monthly",
				Description: model.CostExplorerByServiceUsageTypeMonthlyDescription{CostExplorerRow: row},
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

func CostByServiceUsageLastDay(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := costexplorer.NewFromConfig(cfg)

	params := buildCostByServiceAndUsageInput("DAILY")

	out, err := client.GetCostAndUsage(ctx, params)
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, result := range out.ResultsByTime {
		for _, group := range result.Groups {
			var row model.CostExplorerRow

			row.Estimated = result.Estimated
			row.PeriodStart = result.TimePeriod.Start
			row.PeriodEnd = result.TimePeriod.End

			if len(group.Keys) > 0 {
				row.Dimension1 = aws.String(group.Keys[0])
				if len(group.Keys) > 1 {
					row.Dimension2 = aws.String(group.Keys[1])
				}
			}
			setRowMetrics(&row, group.Metrics)

			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ID:          "service-" + *row.Dimension1 + "-" + *row.Dimension2 + "-cost-" + *row.PeriodEnd,
				Description: model.CostExplorerByServiceUsageTypeDailyDescription{CostExplorerRow: row},
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

func buildCostForecastInput(granularity string) *costexplorer.GetCostForecastInput {
	metric := "UNBLENDED_COST"

	timeFormat := "2006-01-02"
	startTime := time.Now().UTC()
	endTime := startTime.AddDate(0, 1, 0)

	params := &costexplorer.GetCostForecastInput{
		TimePeriod: &types.DateInterval{
			Start: aws.String(startTime.Format(timeFormat)),
			End:   aws.String(endTime.Format(timeFormat)),
		},
		Granularity: types.Granularity(granularity),
		Metric:      types.Metric(metric),
	}

	return params
}

func CostForecastMonthly(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := costexplorer.NewFromConfig(cfg)

	params := buildCostForecastInput("MONTHLY")
	output, err := client.GetCostForecast(ctx, params)
	if err != nil {
		return nil, err
	}
	var values []Resource
	for _, forecast := range output.ForecastResultsByTime {
		resource := Resource{
			Region: describeCtx.KaytuRegion,
			ID:     "forecast-monthly",
			Description: model.CostExplorerForcastMonthlyDescription{CostExplorerRow: model.CostExplorerRow{
				Estimated:   true,
				PeriodStart: forecast.TimePeriod.Start,
				PeriodEnd:   forecast.TimePeriod.End,
				MeanValue:   forecast.MeanValue}},
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

func CostForecastDaily(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := costexplorer.NewFromConfig(cfg)

	params := buildCostForecastInput("DAILY")
	output, err := client.GetCostForecast(ctx, params)
	if err != nil {
		return nil, err
	}
	var values []Resource
	for _, forecast := range output.ForecastResultsByTime {
		resource := Resource{
			Region: describeCtx.KaytuRegion,
			ID:     "forecast-daily",
			Description: model.CostExplorerForcastDailyDescription{CostExplorerRow: model.CostExplorerRow{
				Estimated:   true,
				PeriodStart: forecast.TimePeriod.Start,
				PeriodEnd:   forecast.TimePeriod.End,
				MeanValue:   forecast.MeanValue,
			}},
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
