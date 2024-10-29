package describer

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/backup/types"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/backup"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func BackupPlan(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := backup.NewFromConfig(cfg)
	paginator := backup.NewListBackupPlansPaginator(client, &backup.ListBackupPlansInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.BackupPlansList {
			plan, err := client.GetBackupPlan(ctx, &backup.GetBackupPlanInput{BackupPlanId: v.BackupPlanId})
			if err != nil {
				return nil, err
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    *v.BackupPlanArn,
				Name:   *v.BackupPlanName,
				Description: model.BackupPlanDescription{
					BackupPlan:  v,
					PlanDetails: *plan.BackupPlan,
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

func BackupRecoveryPoint(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := backup.NewFromConfig(cfg)
	paginator := backup.NewListBackupVaultsPaginator(client, &backup.ListBackupVaultsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range page.BackupVaultList {
			recoveryPointPaginator := backup.NewListRecoveryPointsByBackupVaultPaginator(client,
				&backup.ListRecoveryPointsByBackupVaultInput{
					BackupVaultName: item.BackupVaultName,
				})

			for recoveryPointPaginator.HasMorePages() {
				page, err := recoveryPointPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				for _, recoveryPoint := range page.RecoveryPoints {
					out, err := client.DescribeRecoveryPoint(ctx, &backup.DescribeRecoveryPointInput{
						BackupVaultName:  recoveryPoint.BackupVaultName,
						RecoveryPointArn: recoveryPoint.RecoveryPointArn,
					})
					if err != nil {
						return nil, err
					}

					tags := make(map[string]string)
					var arn string
					if out.RecoveryPointArn == nil {
						arn = ""
					} else {
						arn = *out.RecoveryPointArn
					}

					pattern := `arn:aws:backup:[a-z0-9\-]+:[0-9]{12}:recovery-point:.*`

					re := regexp.MustCompile(pattern)

					if re.MatchString(arn) {
						params := &backup.ListTagsInput{
							ResourceArn: aws.String(arn),
						}

						op, err := client.ListTags(ctx, params)
						if err != nil {
							return nil, err
						}
						if op.Tags != nil {
							tags = op.Tags
						}
					}

					resource := Resource{
						Region: describeCtx.KaytuRegion,
						ARN:    *recoveryPoint.RecoveryPointArn,
						Name:   nameFromArn(*out.RecoveryPointArn),
						Description: model.BackupRecoveryPointDescription{
							RecoveryPoint: out,
							Tags:          tags,
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

func BackupProtectedResource(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := backup.NewFromConfig(cfg)
	paginator := backup.NewListProtectedResourcesPaginator(client, &backup.ListProtectedResourcesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, resource := range page.Results {
			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    *resource.ResourceArn,
				Name:   nameFromArn(*resource.ResourceArn),
				Description: model.BackupProtectedResourceDescription{
					ProtectedResource: resource,
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

func BackupSelection(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)

	plans, err := BackupPlan(ctx, cfg, nil)
	if err != nil {
		return nil, err
	}

	client := backup.NewFromConfig(cfg)

	var values []Resource
	for _, plan := range plans {
		paginator := backup.NewListBackupSelectionsPaginator(client, &backup.ListBackupSelectionsInput{
			BackupPlanId: plan.Description.(model.BackupPlanDescription).BackupPlan.BackupPlanId,
		})

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				return nil, err
			}

			for _, v := range page.BackupSelectionsList {
				out, err := client.GetBackupSelection(ctx, &backup.GetBackupSelectionInput{
					BackupPlanId: v.BackupPlanId,
					SelectionId:  v.SelectionId,
				})
				if err != nil {
					return nil, err
				}

				name := "arn:" + describeCtx.Partition + ":backup:" + describeCtx.Region + ":" + describeCtx.AccountID + ":backup-plan:" + *v.BackupPlanId + "/selection/" + *v.SelectionId
				resource := Resource{
					Region: describeCtx.KaytuRegion,
					ARN:    name,
					Name:   *v.SelectionName,
					Description: model.BackupSelectionDescription{
						BackupSelection: v,
						ListOfTags:      out.BackupSelection.ListOfTags,
						Resources:       out.BackupSelection.Resources,
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

func BackupVault(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := backup.NewFromConfig(cfg)
	paginator := backup.NewListBackupVaultsPaginator(client, &backup.ListBackupVaultsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.BackupVaultList {
			notification, err := client.GetBackupVaultNotifications(ctx, &backup.GetBackupVaultNotificationsInput{
				BackupVaultName: v.BackupVaultName,
			})
			if err != nil {
				//Ignore error (otherwise for missing ones it won't work)
				notification = &backup.GetBackupVaultNotificationsOutput{}
			}
			resource, err := backupVaultHandle(ctx, cfg, v, notification)
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
func backupVaultHandle(ctx context.Context, cfg aws.Config, v types.BackupVaultListMember, notification *backup.GetBackupVaultNotificationsOutput) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := backup.NewFromConfig(cfg)
	accessPolicy, err := client.GetBackupVaultAccessPolicy(ctx, &backup.GetBackupVaultAccessPolicyInput{
		BackupVaultName: v.BackupVaultName,
	})
	if err != nil {
		if isErr(err, "ResourceNotFoundException") || isErr(err, "InvalidParameter") {
			accessPolicy = &backup.GetBackupVaultAccessPolicyOutput{}
		} else {
			return Resource{}, err
		}
	}
	tags := make(map[string]string)
	var arn string
	if v.BackupVaultArn != nil {
		arn = *v.BackupVaultArn
	}
	params := &backup.ListTagsInput{
		ResourceArn: aws.String(arn),
	}

	op, err := client.ListTags(ctx, params)
	if err == nil {
		tags = op.Tags
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *v.BackupVaultArn,
		Name:   *v.BackupVaultName,
		Description: model.BackupVaultDescription{
			BackupVault:       v,
			Policy:            accessPolicy.Policy,
			BackupVaultEvents: notification.BackupVaultEvents,
			SNSTopicArn:       notification.SNSTopicArn,
			Tags:              tags,
		},
	}
	return resource, nil
}
func GetBackupVault(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	backupVaultName := fields["name"]
	client := backup.NewFromConfig(cfg)

	listBackup, err := client.ListBackupVaults(ctx, &backup.ListBackupVaultsInput{})
	if err != nil {
		if isErr(err, "ListBackupVaultsNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	for _, v := range listBackup.BackupVaultList {

		if *v.BackupVaultName != backupVaultName {
			continue
		}
		notification := &backup.GetBackupVaultNotificationsOutput{}

		resource, err := backupVaultHandle(ctx, cfg, v, notification)
		if err != nil {
			return nil, err
		}

		values = append(values, resource)
	}
	return values, nil
}

func BackupFramework(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := backup.NewFromConfig(cfg)
	paginator := backup.NewListFrameworksPaginator(client, &backup.ListFrameworksInput{})

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

		for _, v := range page.Frameworks {
			resource, err := backupFrameworkHandle(ctx, cfg, v)
			if err != nil {
				if isErr(err, "AccessDeniedException") {
					return nil, nil
				} else {
					return nil, err
				}
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
func backupFrameworkHandle(ctx context.Context, cfg aws.Config, v types.Framework) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := backup.NewFromConfig(cfg)

	framework, err := client.DescribeFramework(ctx, &backup.DescribeFrameworkInput{
		FrameworkName: v.FrameworkName,
	})
	if err != nil {
		return Resource{}, err
	}

	tags, err := client.ListTags(ctx, &backup.ListTagsInput{
		ResourceArn: v.FrameworkArn,
	})

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *v.FrameworkArn,
		Name:   *v.FrameworkName,
		Description: model.BackupFrameworkDescription{
			Framework: *framework,
			Tags:      tags.Tags,
		},
	}
	return resource, nil
}
func GetBackupFramework(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	var values []Resource
	frameworkName := fields["name"]
	client := backup.NewFromConfig(cfg)

	describe, err := client.ListFrameworks(ctx, &backup.ListFrameworksInput{})
	if err != nil {
		if isErr(err, "ListFrameworksNotFound") || isErr(err, "InvalidParameterValue") || isErr(err, "AccessDeniedException") {
			return nil, nil
		}
		return nil, err
	}
	for _, v := range describe.Frameworks {
		if *v.FrameworkName != frameworkName {
			continue
		}

		resource, err := backupFrameworkHandle(ctx, cfg, v)
		if err != nil {
			return nil, err
		}

		values = append(values, resource)
	}
	if values == nil {
		return nil, nil
	}
	return values, nil
}

func BackupLegalHold(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := backup.NewFromConfig(cfg)
	paginator := backup.NewListLegalHoldsPaginator(client, &backup.ListLegalHoldsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.LegalHolds {
			legalHold, err := client.GetLegalHold(ctx, &backup.GetLegalHoldInput{
				LegalHoldId: v.LegalHoldId,
			})
			if err != nil {
				return nil, err
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				Name:   *v.Title,
				ARN:    *v.LegalHoldArn,
				ID:     *v.LegalHoldId,
				Description: model.BackupLegalHoldDescription{
					LegalHold: *legalHold,
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

func BackupReportPlan(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := backup.NewFromConfig(cfg)
	paginator := backup.NewListReportPlansPaginator(client, &backup.ListReportPlansInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.ReportPlans {
			resource, err := backupReportPlanHandle(ctx, cfg, v)
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
func backupReportPlanHandle(ctx context.Context, cfg aws.Config, v types.ReportPlan) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := backup.NewFromConfig(cfg)

	reportPlan, err := client.DescribeReportPlan(ctx, &backup.DescribeReportPlanInput{
		ReportPlanName: v.ReportPlanName,
	})
	if err != nil {
		return Resource{}, err
	}

	tags, err := client.ListTags(ctx, &backup.ListTagsInput{
		ResourceArn: v.ReportPlanArn,
	})

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *v.ReportPlanArn,
		Name:   *v.ReportPlanName,
		Description: model.BackupReportPlanDescription{
			ReportPlan: *reportPlan.ReportPlan,
			Tags:       tags.Tags,
		},
	}
	return resource, nil
}
func GetBackupReportPlan(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	var values []Resource
	name := fields["name"]
	client := backup.NewFromConfig(cfg)

	describe, err := client.ListReportPlans(ctx, &backup.ListReportPlansInput{})
	if err != nil {
		if isErr(err, "ListReportPlansNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}
	for _, v := range describe.ReportPlans {
		if *v.ReportPlanName != name {
			continue
		}

		resource, err := backupReportPlanHandle(ctx, cfg, v)
		if err != nil {
			return nil, err
		}

		values = append(values, resource)
	}
	if values == nil {
		return nil, nil
	}
	return values, nil
}

func BackupRegionSetting(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := backup.NewFromConfig(cfg)
	regionSetting, err := client.DescribeRegionSettings(ctx, &backup.DescribeRegionSettingsInput{})
	if err != nil {
		return nil, err
	}

	var values []Resource

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		Description: model.BackupRegionSettingDescription{
			Region:                           describeCtx.KaytuRegion,
			ResourceTypeManagementPreference: regionSetting.ResourceTypeManagementPreference,
			ResourceTypeOptInPreference:      regionSetting.ResourceTypeOptInPreference,
		},
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
