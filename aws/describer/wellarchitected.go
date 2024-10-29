package describer

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wellarchitected"
	"github.com/aws/aws-sdk-go-v2/service/wellarchitected/types"
	"github.com/opengovern/og-aws-describer/aws/model"
	"reflect"
)

func WellArchitectedWorkload(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := wellarchitected.NewFromConfig(cfg)
	paginator := wellarchitected.NewListWorkloadsPaginator(client, &wellarchitected.ListWorkloadsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if isErr(err, "AccessDeniedException") {
				return nil, nil
			} else {
				return nil, err
			}
		}
		for _, v := range page.WorkloadSummaries {
			op, err := client.GetWorkload(ctx, &wellarchitected.GetWorkloadInput{
				WorkloadId: v.WorkloadId,
			})
			if err != nil {
				if isErr(err, "AccessDeniedException") {
					return nil, nil
				} else {
					return nil, err
				}
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    *v.WorkloadArn,
				Name:   *v.WorkloadName,
				Description: model.WellArchitectedWorkloadDescription{
					WorkloadSummary: v,
					Workload:        op.Workload,
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

func WellArchitectedAnswer(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := wellarchitected.NewFromConfig(cfg)
	paginator := wellarchitected.NewListWorkloadsPaginator(client, &wellarchitected.ListWorkloadsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, v := range page.WorkloadSummaries {
			for _, lensAlias := range v.Lenses {
				input := &wellarchitected.ListAnswersInput{
					LensAlias:  aws.String(lensAlias),
					WorkloadId: aws.String(*v.WorkloadId),
				}

				paginator := wellarchitected.NewListAnswersPaginator(client, input, func(o *wellarchitected.ListAnswersPaginatorOptions) {
					o.StopOnDuplicateToken = true
				})
				for paginator.HasMorePages() {
					page, err := paginator.NextPage(ctx)
					if err != nil {
						if isErr(err, "AccessDeniedException") {
							return nil, nil
						} else {
							return nil, err
						}
					}
					for _, a := range page.AnswerSummaries {
						params := &wellarchitected.GetAnswerInput{
							QuestionId: aws.String(*a.QuestionId),
							LensAlias:  aws.String(lensAlias),
							WorkloadId: aws.String(*v.WorkloadId),
						}
						op, err := client.GetAnswer(ctx, params)
						if err != nil {
							return nil, err
						}

						resource := Resource{
							Region: describeCtx.KaytuRegion,
							Description: model.WellArchitectedAnswerDescription{
								Answer:          *op.Answer,
								WorkloadId:      *op.WorkloadId,
								WorkloadName:    *v.WorkloadName,
								LensAlias:       lensAlias,
								LensArn:         *op.LensArn,
								MilestoneNumber: op.MilestoneNumber,
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
			}
		}
	}
	return values, nil
}

func WellArchitectedCheckDetail(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := wellarchitected.NewFromConfig(cfg)
	answers, err := WellArchitectedAnswer(ctx, cfg, stream)
	if err != nil {
		if isErr(err, "AccessDeniedException") {
			return nil, nil
		} else {
			return nil, err
		}
	}
	var values []Resource
	for _, answer := range answers {
		des := answer.Description.(model.WellArchitectedAnswerDescription)
		params := &wellarchitected.GetAnswerInput{
			QuestionId: aws.String(*des.Answer.QuestionId),
			LensAlias:  aws.String(des.LensAlias),
			WorkloadId: aws.String(des.WorkloadId),
		}
		op, err := client.GetAnswer(ctx, params)
		if err != nil {
			if isErr(err, "AccessDeniedException") {
				return nil, nil
			} else {
				return nil, err
			}
		}
		answer := op.Answer
		for _, choice := range answer.Choices {
			input := &wellarchitected.ListCheckDetailsInput{
				LensArn:    aws.String(*op.LensArn),
				PillarId:   aws.String(*answer.PillarId),
				QuestionId: aws.String(*answer.QuestionId),
				WorkloadId: aws.String(*op.WorkloadId),
				ChoiceId:   aws.String(*choice.ChoiceId),
			}

			paginator := wellarchitected.NewListCheckDetailsPaginator(client, input, func(o *wellarchitected.ListCheckDetailsPaginatorOptions) {
				o.StopOnDuplicateToken = true
			})
			for paginator.HasMorePages() {
				page, err := paginator.NextPage(ctx)
				if err != nil {
					if isErr(err, "AccessDeniedException") {
						return nil, nil
					} else {
						return nil, err
					}
				}
				for _, c := range page.CheckDetails {

					resource := Resource{
						Region: describeCtx.KaytuRegion,
						Name:   des.WorkloadName,
						Description: model.WellArchitectedCheckDetailDescription{
							CheckDetail: c,
							WorkloadId:  des.WorkloadId,
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
		}
	}

	return values, nil
}

func WellArchitectedCheckSummary(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := wellarchitected.NewFromConfig(cfg)
	answers, err := WellArchitectedAnswer(ctx, cfg, stream)
	if err != nil {
		if isErr(err, "AccessDeniedException") {
			return nil, nil
		} else {
			return nil, err
		}
	}
	var values []Resource
	for _, answer := range answers {
		des := answer.Description.(model.WellArchitectedAnswerDescription)
		params := &wellarchitected.GetAnswerInput{
			QuestionId: aws.String(*des.Answer.QuestionId),
			LensAlias:  aws.String(des.LensAlias),
			WorkloadId: aws.String(des.WorkloadId),
		}
		op, err := client.GetAnswer(ctx, params)
		if err != nil {
			if isErr(err, "AccessDeniedException") {
				return nil, nil
			} else {
				return nil, err
			}
		}
		answer := op.Answer
		for _, choice := range answer.Choices {
			input := &wellarchitected.ListCheckSummariesInput{
				LensArn:    aws.String(*op.LensArn),
				PillarId:   aws.String(*answer.PillarId),
				QuestionId: aws.String(*answer.QuestionId),
				WorkloadId: aws.String(*op.WorkloadId),
				ChoiceId:   aws.String(*choice.ChoiceId),
			}

			paginator := wellarchitected.NewListCheckSummariesPaginator(client, input, func(o *wellarchitected.ListCheckSummariesPaginatorOptions) {
				o.StopOnDuplicateToken = true
			})
			for paginator.HasMorePages() {
				page, err := paginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}
				for _, c := range page.CheckSummaries {

					resource := Resource{
						Region: describeCtx.KaytuRegion,
						Name:   *c.Name,
						Description: model.WellArchitectedCheckSummaryDescription{
							CheckSummary: c,
							WorkloadId:   des.WorkloadId,
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
		}
	}

	return values, nil
}

func WellArchitectedConsolidatedReport(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := wellarchitected.NewFromConfig(cfg)
	var values []Resource
	for _, rFormat := range []types.ReportFormat{types.ReportFormatPdf, types.ReportFormatJson} {
		input := &wellarchitected.GetConsolidatedReportInput{
			IncludeSharedResources: aws.Bool(true),
			Format:                 rFormat,
		}
		sharedValues, err := WellArchitectedConsolidatedReportHelper(ctx, cfg, stream, client, describeCtx, input)
		if err != nil {
			if isErr(err, "AccessDeniedException") {
				return nil, nil
			} else {
				return nil, err
			}
		}
		input2 := &wellarchitected.GetConsolidatedReportInput{
			IncludeSharedResources: aws.Bool(false),
			Format:                 rFormat,
		}
		notSharedValues, err := WellArchitectedConsolidatedReportHelper(ctx, cfg, stream, client, describeCtx, input2)
		if err != nil {
			return nil, err
		}
		for _, value := range sharedValues {
			if !Contains(values, value) {
				values = append(values, value)
			}
		}
		for _, value := range notSharedValues {
			if !Contains(values, value) {
				values = append(values, value)
			}
		}
	}
	return values, nil
}

func Contains(values []Resource, value Resource) bool {
	for _, v := range values {
		val := v.Description.(model.WellArchitectedCheckConsolidatedReportDescription)
		val1 := value.Description.(model.WellArchitectedCheckConsolidatedReportDescription)
		if reflect.DeepEqual(val1, val) {
			return true
		}
	}
	return false
}

func WellArchitectedConsolidatedReportHelper(ctx context.Context, cfg aws.Config, stream *StreamSender, client *wellarchitected.Client, describeCtx DescribeContext, input *wellarchitected.GetConsolidatedReportInput) ([]Resource, error) {
	paginator := wellarchitected.NewGetConsolidatedReportPaginator(client, input)

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if isErr(err, "AccessDeniedException") {
				return nil, nil
			} else {
				return nil, err
			}
		}
		for _, v := range page.Metrics {

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				Description: model.WellArchitectedCheckConsolidatedReportDescription{
					IncludeSharedResources: input.IncludeSharedResources,
					ConsolidateReport:      v,
					Base64:                 *page.Base64String,
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

func WellArchitectedLens(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := wellarchitected.NewFromConfig(cfg)
	paginator := wellarchitected.NewListLensesPaginator(client, &wellarchitected.ListLensesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if isErr(err, "AccessDeniedException") {
				return nil, nil
			} else {
				return nil, err
			}
		}
		for _, v := range page.LensSummaries {
			input := v.LensAlias
			if input == nil {
				input = v.LensArn
			}
			op, err := client.GetLens(ctx, &wellarchitected.GetLensInput{
				LensAlias: input,
			})
			if err != nil {
				if isErr(err, "AccessDeniedException") {
					return nil, nil
				} else {
					return nil, err
				}
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    *v.LensArn,
				Description: model.WellArchitectedLensDescription{
					LensSummary: v,
					Lens:        *op.Lens,
				},
			}
			if v.LensName != nil {
				resource.Name = *v.LensName
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

func WellArchitectedLensReview(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := wellarchitected.NewFromConfig(cfg)
	paginator := wellarchitected.NewListWorkloadsPaginator(client, &wellarchitected.ListWorkloadsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if isErr(err, "AccessDeniedException") {
				return nil, nil
			} else {
				return nil, err
			}
		}
		for _, v := range page.WorkloadSummaries {
			op, err := client.ListLensReviews(ctx, &wellarchitected.ListLensReviewsInput{
				WorkloadId: v.WorkloadId,
			})
			if err != nil {
				if isErr(err, "AccessDeniedException") {
					op = &wellarchitected.ListLensReviewsOutput{}
				} else {
					return nil, err
				}
			}
			for _, r := range op.LensReviewSummaries {
				review, err := client.GetLensReview(ctx, &wellarchitected.GetLensReviewInput{
					LensAlias:  r.LensAlias,
					WorkloadId: v.WorkloadId,
				})
				if err != nil {
					if isErr(err, "AccessDeniedException") {
						op = &wellarchitected.ListLensReviewsOutput{}
					} else {
						return nil, err
					}
				}
				resource := Resource{
					Region: describeCtx.KaytuRegion,
					Description: model.WellArchitectedLensReviewDescription{
						LensReview: *review.LensReview,
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
	}
	return values, nil
}

func WellArchitectedLensReviewImprovement(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := wellarchitected.NewFromConfig(cfg)
	paginator := wellarchitected.NewListWorkloadsPaginator(client, &wellarchitected.ListWorkloadsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if isErr(err, "AccessDeniedException") {
				return nil, nil
			} else {
				return nil, err
			}
		}
		for _, v := range page.WorkloadSummaries {
			for _, lense := range v.Lenses {
				improvementPaginator := wellarchitected.NewListLensReviewImprovementsPaginator(client, &wellarchitected.ListLensReviewImprovementsInput{
					LensAlias:  aws.String(lense),
					WorkloadId: v.WorkloadId,
				})
				for improvementPaginator.HasMorePages() {
					output, err := improvementPaginator.NextPage(ctx)
					if err != nil {
						if isErr(err, "AccessDeniedException") {
							return nil, nil
						} else {
							return nil, err
						}
					}
					for _, improvement := range output.ImprovementSummaries {
						resource := Resource{
							Region: describeCtx.KaytuRegion,
							Description: model.WellArchitectedLensReviewImprovementDescription{
								LensAlias:          *output.LensAlias,
								LensArn:            *output.LensArn,
								MilestoneNumber:    output.MilestoneNumber,
								WorkloadId:         *output.WorkloadId,
								ImprovementSummary: improvement,
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
			}
		}
	}
	return values, nil
}

func WellArchitectedLensReviewReport(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := wellarchitected.NewFromConfig(cfg)
	paginator := wellarchitected.NewListWorkloadsPaginator(client, &wellarchitected.ListWorkloadsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if isErr(err, "AccessDeniedException") {
				return nil, nil
			} else {
				return nil, err
			}
		}
		for _, v := range page.WorkloadSummaries {
			for _, lense := range v.Lenses {
				report, err := client.GetLensReviewReport(ctx, &wellarchitected.GetLensReviewReportInput{
					LensAlias:  aws.String(lense),
					WorkloadId: v.WorkloadId,
				})
				if err != nil {
					if isErr(err, "AccessDeniedException") {
						return nil, nil
					} else {
						return nil, err
					}
				}

				resource := Resource{
					Region: describeCtx.KaytuRegion,
					Description: model.WellArchitectedLensReviewReportDescription{
						Report:          *report.LensReviewReport,
						MilestoneNumber: report.MilestoneNumber,
						WorkloadId:      *v.WorkloadId,
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

	}
	return values, nil
}

func WellArchitectedLensShare(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := wellarchitected.NewFromConfig(cfg)

	lenses, err := WellArchitectedLens(ctx, cfg, nil)
	if err != nil {
		return nil, err
	}
	var values []Resource
	for _, v := range lenses {
		lens := v.Description.(model.WellArchitectedLensDescription).Lens
		input := &wellarchitected.ListLensSharesInput{
			LensAlias: lens.LensArn,
		}
		paginator := wellarchitected.NewListLensSharesPaginator(client, input)
		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				if isErr(err, "AccessDeniedException") {
					return nil, nil
				} else {
					return nil, err
				}
			}
			for _, share := range page.LensShareSummaries {
				resource := Resource{
					Region: describeCtx.KaytuRegion,
					Description: model.WellArchitectedLensShareDescription{
						Lens:  lens,
						Share: share,
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
	}
	return values, nil
}

func WellArchitectedMilestone(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := wellarchitected.NewFromConfig(cfg)
	paginator := wellarchitected.NewListWorkloadsPaginator(client, &wellarchitected.ListWorkloadsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, v := range page.WorkloadSummaries {
			paginator := wellarchitected.NewListMilestonesPaginator(client, &wellarchitected.ListMilestonesInput{
				WorkloadId: aws.String(*v.WorkloadId),
			})
			for paginator.HasMorePages() {
				output, err := paginator.NextPage(ctx)
				if err != nil {
					if isErr(err, "AccessDeniedException") {
						return nil, nil
					} else {
						return nil, err
					}
				}
				for _, m := range output.MilestoneSummaries {
					milestone, err := client.GetMilestone(ctx, &wellarchitected.GetMilestoneInput{
						WorkloadId:      aws.String(*v.WorkloadId),
						MilestoneNumber: m.MilestoneNumber,
					})
					if err != nil {
						return nil, err
					}
					resource := Resource{
						Region: describeCtx.KaytuRegion,
						Description: model.WellArchitectedMilestoneDescription{
							Milestone: *milestone.Milestone,
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
		}
	}
	return values, nil
}

func WellArchitectedNotification(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := wellarchitected.NewFromConfig(cfg)
	paginator := wellarchitected.NewListNotificationsPaginator(client, &wellarchitected.ListNotificationsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, v := range page.NotificationSummaries {

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				Description: model.WellArchitectedNotificationDescription{
					Notification: v,
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

func WellArchitectedShareInvitation(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := wellarchitected.NewFromConfig(cfg)
	paginator := wellarchitected.NewListShareInvitationsPaginator(client, &wellarchitected.ListShareInvitationsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if isErr(err, "AccessDeniedException") {
				return nil, nil
			} else {
				return nil, err
			}
		}
		for _, v := range page.ShareInvitationSummaries {
			resource := Resource{
				Region: describeCtx.KaytuRegion,
				Description: model.WellArchitectedShareInvitationDescription{
					ShareInvitation: v,
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

func WellArchitectedWorkloadShare(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := wellarchitected.NewFromConfig(cfg)
	paginator := wellarchitected.NewListWorkloadsPaginator(client, &wellarchitected.ListWorkloadsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if isErr(err, "AccessDeniedException") {
				return nil, nil
			} else {
				return nil, err
			}
		}
		for _, v := range page.WorkloadSummaries {
			paginator := wellarchitected.NewListWorkloadSharesPaginator(client, &wellarchitected.ListWorkloadSharesInput{
				WorkloadId: aws.String(*v.WorkloadId),
			})
			for paginator.HasMorePages() {
				output, err := paginator.NextPage(ctx)
				if err != nil {
					if isErr(err, "AccessDeniedException") {
						return nil, nil
					} else {
						return nil, err
					}
				}
				for _, m := range output.WorkloadShareSummaries {
					arn := "arn:" + describeCtx.Partition + ":waf::" + describeCtx.AccountID + ":ratebasedrule" + "/" + *m.ShareId

					resource := Resource{
						Region: describeCtx.KaytuRegion,
						Description: model.WellArchitectedWorkloadShareDescription{
							Share:      m,
							WorkloadId: *v.WorkloadId,
							Arn:        arn,
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
		}
	}
	return values, nil
}
