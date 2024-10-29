package describer

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/accessanalyzer"
	"github.com/aws/aws-sdk-go-v2/service/accessanalyzer/types"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func GetAccessAnalyzerAnalyzer(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	analyzerName := fields["name"]
	client := accessanalyzer.NewFromConfig(cfg)
	v, err := client.GetAnalyzer(ctx, &accessanalyzer.GetAnalyzerInput{
		AnalyzerName: &analyzerName,
	})
	if err != nil {
		return nil, err
	}

	findings, err := getAnalyzerFindings(ctx, client, v.Analyzer.Arn)
	if err != nil {
		return nil, err
	}

	return []Resource{
		{
			Region: describeCtx.KaytuRegion,
			ARN:    *v.Analyzer.Arn,
			Name:   *v.Analyzer.Name,
			Description: model.AccessAnalyzerAnalyzerDescription{
				Analyzer: *v.Analyzer,
				Findings: findings,
			},
		}}, nil
}

func AccessAnalyzerAnalyzer(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := accessanalyzer.NewFromConfig(cfg)
	paginator := accessanalyzer.NewListAnalyzersPaginator(client, &accessanalyzer.ListAnalyzersInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Analyzers {
			findings, err := getAnalyzerFindings(ctx, client, v.Arn)
			if err != nil {
				return nil, err
			}
			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    *v.Arn,
				Name:   *v.Name,
				Description: model.AccessAnalyzerAnalyzerDescription{
					Analyzer: v,
					Findings: findings,
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

func getAnalyzerFindings(ctx context.Context, client *accessanalyzer.Client, analyzerArn *string) ([]types.FindingSummary, error) {
	paginator := accessanalyzer.NewListFindingsPaginator(client, &accessanalyzer.ListFindingsInput{
		AnalyzerArn: analyzerArn,
	})

	var findings []types.FindingSummary
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		findings = append(findings, page.Findings...)
	}

	return findings, nil
}

func AccessAnalyzerAnalyzerFinding(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := accessanalyzer.NewFromConfig(cfg)
	paginator := accessanalyzer.NewListAnalyzersPaginator(client, &accessanalyzer.ListAnalyzersInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Analyzers {
			findings, err := getAnalyzerFindings(ctx, client, v.Arn)
			if err != nil {
				return nil, err
			}
			for _, finding := range findings {
				resource := Resource{
					Region: describeCtx.KaytuRegion,
					ID:     *finding.Id,
					Description: model.AccessAnalyzerAnalyzerFindingDescription{
						AnalyzerArn: *v.Arn,
						Finding:     finding,
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
	}

	return values, nil
}
