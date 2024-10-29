package describer

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/auditmanager/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/auditmanager"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func AuditManagerAssessment(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := auditmanager.NewFromConfig(cfg)
	paginator := auditmanager.NewListAssessmentsPaginator(client, &auditmanager.ListAssessmentsInput{})

	//

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, assessmentMetadataItem := range page.AssessmentMetadata {

			assessment, err := client.GetAssessment(ctx, &auditmanager.GetAssessmentInput{
				AssessmentId: assessmentMetadataItem.Id,
			})
			if err != nil {
				return nil, err
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    *assessment.Assessment.Arn,
				Name:   *assessment.Assessment.Metadata.Name,
				ID:     *assessment.Assessment.Metadata.Id,
				Description: model.AuditManagerAssessmentDescription{
					Assessment: *assessment.Assessment,
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

func AuditManagerControl(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := auditmanager.NewFromConfig(cfg)

	var values []Resource
	for _, ctype := range []types.ControlType{types.ControlTypeStandard, types.ControlTypeCustom} {
		paginator := auditmanager.NewListControlsPaginator(client, &auditmanager.ListControlsInput{
			ControlType: ctype,
		})

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				return nil, err
			}

			for _, controlMetadata := range page.ControlMetadataList {
				control, err := client.GetControl(ctx, &auditmanager.GetControlInput{
					ControlId: controlMetadata.Id,
				})
				if err != nil {
					return nil, err
				}

				resource := Resource{
					Region: describeCtx.KaytuRegion,
					ARN:    *control.Control.Arn,
					Name:   *control.Control.Name,
					ID:     *control.Control.Id,
					Description: model.AuditManagerControlDescription{
						Control: *control.Control,
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

func GetAuditManagerControl(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	controlID := fields["id"]
	client := auditmanager.NewFromConfig(cfg)
	control, err := client.GetControl(ctx, &auditmanager.GetControlInput{
		ControlId: &controlID,
	})
	if err != nil {
		return nil, err
	}

	var values []Resource
	values = append(values, Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *control.Control.Arn,
		Name:   *control.Control.Name,
		ID:     *control.Control.Id,
		Description: model.AuditManagerControlDescription{
			Control: *control.Control,
		},
	})

	return values, nil
}

func AuditManagerEvidence(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := auditmanager.NewFromConfig(cfg)
	paginator := auditmanager.NewListAssessmentsPaginator(client, &auditmanager.ListAssessmentsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, assessmentMetadataItem := range page.AssessmentMetadata {
			evidenceFolderPaginator := auditmanager.NewGetEvidenceFoldersByAssessmentPaginator(client, &auditmanager.GetEvidenceFoldersByAssessmentInput{
				AssessmentId: assessmentMetadataItem.Id,
			})

			for evidenceFolderPaginator.HasMorePages() {
				evidenceFolderPage, err := evidenceFolderPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				for _, evidenceFolder := range evidenceFolderPage.EvidenceFolders {
					evidencePaginator := auditmanager.NewGetEvidenceByEvidenceFolderPaginator(client, &auditmanager.GetEvidenceByEvidenceFolderInput{
						EvidenceFolderId: evidenceFolder.Id,
					})

					for evidencePaginator.HasMorePages() {
						evidencePage, err := evidencePaginator.NextPage(ctx)
						if err != nil {
							return nil, err
						}

						for _, evidence := range evidencePage.Evidence {
							arn := fmt.Sprintf("arn:%s:auditmanager:%s:%s:evidence/%s", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID, *evidence.Id)
							resource := Resource{
								Region: describeCtx.KaytuRegion,
								ARN:    arn,
								ID:     *evidence.Id,
								Description: model.AuditManagerEvidenceDescription{
									Evidence:     evidence,
									ControlSetID: *evidenceFolder.ControlSetId,
									AssessmentID: *assessmentMetadataItem.Id,
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
	}

	return values, nil
}

func AuditManagerEvidenceFolder(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := auditmanager.NewFromConfig(cfg)
	paginator := auditmanager.NewListAssessmentsPaginator(client, &auditmanager.ListAssessmentsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, assessmentMetadataItem := range page.AssessmentMetadata {
			evidenceFolderPaginator := auditmanager.NewGetEvidenceFoldersByAssessmentPaginator(client, &auditmanager.GetEvidenceFoldersByAssessmentInput{
				AssessmentId: assessmentMetadataItem.Id,
			})

			for evidenceFolderPaginator.HasMorePages() {
				evidenceFolderPage, err := evidenceFolderPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				for _, evidenceFolder := range evidenceFolderPage.EvidenceFolders {
					arn := fmt.Sprintf("arn:%s:auditmanager:%s:%s:evidence-folder/%s", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID, *evidenceFolder.Id)
					resource := Resource{
						Region: describeCtx.KaytuRegion,
						ARN:    arn,
						Name:   *evidenceFolder.Name,
						ID:     *evidenceFolder.Id,
						Description: model.AuditManagerEvidenceFolderDescription{
							EvidenceFolder: evidenceFolder,
							AssessmentID:   *assessmentMetadataItem.Id,
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

func AuditManagerFramework(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := auditmanager.NewFromConfig(cfg)

	var values []Resource
	for _, ftype := range []types.FrameworkType{types.FrameworkTypeStandard, types.FrameworkTypeCustom} {
		paginator := auditmanager.NewListAssessmentFrameworksPaginator(client, &auditmanager.ListAssessmentFrameworksInput{
			FrameworkType: ftype,
		})
		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				return nil, err
			}

			for _, frameworkMetadata := range page.FrameworkMetadataList {
				framework, err := client.GetAssessmentFramework(ctx, &auditmanager.GetAssessmentFrameworkInput{
					FrameworkId: frameworkMetadata.Id,
				})
				if err != nil {
					return nil, err
				}

				resource := Resource{
					Region: describeCtx.KaytuRegion,
					ARN:    *framework.Framework.Arn,
					Name:   *framework.Framework.Name,
					ID:     *framework.Framework.Id,
					Description: model.AuditManagerFrameworkDescription{
						Framework: *framework.Framework,
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
