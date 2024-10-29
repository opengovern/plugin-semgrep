package describer

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/smithy-go"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func SSMManagedInstance(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ssm.NewFromConfig(cfg)
	paginator := ssm.NewDescribeInstanceInformationPaginator(client, &ssm.DescribeInstanceInformationInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range page.InstanceInformationList {
			arn := "arn:" + describeCtx.Partition + ":ssm:" + describeCtx.Region + ":" + describeCtx.AccountID + ":managed-instance/" + *item.InstanceId
			name := ""
			if item.Name != nil {
				name = *item.Name
			} else {
				name = *item.InstanceId
			}
			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    arn,
				Name:   name,
				Description: model.SSMManagedInstanceDescription{
					InstanceInformation: item,
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

func SSMInventory(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ssm.NewFromConfig(cfg)
	paginator := ssm.NewGetInventoryPaginator(client, &ssm.GetInventoryInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, inventory := range page.Entities {
			if inventory.Data != nil {
				for _, data := range inventory.Data {
					var schemas []types.InventoryItemSchema
					schemaPaginator := ssm.NewGetInventorySchemaPaginator(client, &ssm.GetInventorySchemaInput{})
					for schemaPaginator.HasMorePages() {
						schemaPage, err := schemaPaginator.NextPage(ctx)
						if err != nil {
							return nil, err
						}
						schemas = append(schemas, schemaPage.Schemas...)
					}
					resource := Resource{
						Region: describeCtx.KaytuRegion,
						ID:     *inventory.Id,
						Name:   *inventory.Id,
						Description: model.SSMInventoryDescription{
							Id:            inventory.Id,
							CaptureTime:   data.CaptureTime,
							SchemaVersion: data.SchemaVersion,
							TypeName:      data.TypeName,
							Content:       data.Content,
							Schemas:       schemas,
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

func SSMInventoryEntry(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ssm.NewFromConfig(cfg)
	paginator := ssm.NewGetInventoryPaginator(client, &ssm.GetInventoryInput{})
	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, inventory := range page.Entities {
			if inventory.Data != nil {
				for _, data := range inventory.Data {
					op, err := client.ListInventoryEntries(ctx, &ssm.ListInventoryEntriesInput{
						TypeName:   data.TypeName,
						InstanceId: inventory.Id,
					})
					if err != nil {
						return nil, err
					}
					for _, v := range op.Entries {
						resource := Resource{
							Region: describeCtx.KaytuRegion,
							ID:     *op.InstanceId,
							Name:   *op.InstanceId,
							Description: model.SSMInventoryEntryDescription{
								InstanceId:    op.InstanceId,
								TypeName:      op.TypeName,
								CaptureTime:   op.CaptureTime,
								SchemaVersion: op.SchemaVersion,
								Entries:       v,
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

func SSMManagedInstanceCompliance(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ssm.NewFromConfig(cfg)
	paginator := ssm.NewDescribeInstanceInformationPaginator(client, &ssm.DescribeInstanceInformationInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range page.InstanceInformationList {
			cpaginator := ssm.NewListComplianceItemsPaginator(client, &ssm.ListComplianceItemsInput{
				ResourceIds: []string{*item.InstanceId},
			})

			for cpaginator.HasMorePages() {
				cpage, err := cpaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				for _, item := range cpage.ComplianceItems {
					arn := "arn:" + describeCtx.Partition + ":ssm:" + describeCtx.Region + ":" + describeCtx.AccountID + ":managed-instance/" + *item.ResourceId + "/compliance-item/" + *item.Id + ":" + *item.ComplianceType
					resource := Resource{
						Region: describeCtx.KaytuRegion,
						ARN:    arn,
						Name:   *item.Title,
						Description: model.SSMManagedInstanceComplianceDescription{
							ComplianceItem: item,
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

func SSMAssociation(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ssm.NewFromConfig(cfg)
	paginator := ssm.NewListAssociationsPaginator(client, &ssm.ListAssociationsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Associations {
			out, err := client.DescribeAssociation(ctx, &ssm.DescribeAssociationInput{
				AssociationId: v.AssociationId,
			})
			if err != nil {
				return nil, err
			}

			arn := fmt.Sprintf("arn:%s:ssm:%s:%s:association/%s", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID, *v.AssociationId)
			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ID:     *v.AssociationId,
				Name:   *v.Name,
				ARN:    arn,
				Description: model.SSMAssociationDescription{
					AssociationItem: v,
					Association:     out,
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

func SSMDocument(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ssm.NewFromConfig(cfg)
	paginator := ssm.NewListDocumentsPaginator(client, &ssm.ListDocumentsInput{
		Filters: []types.DocumentKeyValuesFilter{
			{
				Key:    aws.String("Owner"),
				Values: []string{"Self"},
			},
		},
	})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.DocumentIdentifiers {
			permissions, err := client.DescribeDocumentPermission(ctx, &ssm.DescribeDocumentPermissionInput{
				Name:           v.Name,
				PermissionType: "Share",
			})
			if err != nil {
				return nil, err
			}

			data, err := client.DescribeDocument(ctx, &ssm.DescribeDocumentInput{
				Name: v.Name,
			})
			if err != nil {
				return nil, err
			}

			arn := fmt.Sprintf("arn:%s:ssm:%s:%s:document", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID)
			if strings.HasPrefix(*v.Name, "/") {
				arn += *v.Name
			} else {
				arn += "/" + *v.Name
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ID:     *v.Name,
				Name:   *v.Name,
				ARN:    arn,
				Description: model.SSMDocumentDescription{
					DocumentIdentifier: v,
					Document:           data,
					Permissions:        permissions,
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

func SSMDocumentPermission(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ssm.NewFromConfig(cfg)
	paginator := ssm.NewListDocumentsPaginator(client, &ssm.ListDocumentsInput{
		Filters: []types.DocumentKeyValuesFilter{
			{
				Key:    aws.String("Owner"),
				Values: []string{"Self"},
			},
		},
	})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, v := range page.DocumentIdentifiers {
			permissions, err := client.DescribeDocumentPermission(ctx, &ssm.DescribeDocumentPermissionInput{
				Name:           v.Name,
				PermissionType: "Share",
			})
			if err != nil {
				return nil, err
			}

			data, err := client.DescribeDocument(ctx, &ssm.DescribeDocumentInput{
				Name: v.Name,
			})
			if err != nil {
				return nil, err
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ID:     *v.Name,
				Name:   *v.Name,
				Description: model.SSMDocumentPermissionDescription{
					Document:    data,
					Permissions: permissions,
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

func SSMMaintenanceWindow(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ssm.NewFromConfig(cfg)
	paginator := ssm.NewDescribeMaintenanceWindowsPaginator(client, &ssm.DescribeMaintenanceWindowsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.WindowIdentities {
			data, err := client.GetMaintenanceWindow(ctx, &ssm.GetMaintenanceWindowInput{
				WindowId: v.WindowId,
			})
			if err != nil {
				return nil, err
			}

			op, err := client.ListTagsForResource(ctx, &ssm.ListTagsForResourceInput{
				ResourceType: "MaintenanceWindow",
				ResourceId:   v.WindowId,
			})
			if err != nil {
				return nil, err
			}

			op2, err := client.DescribeMaintenanceWindowTargets(ctx, &ssm.DescribeMaintenanceWindowTargetsInput{
				WindowId: v.WindowId,
			})
			if err != nil {
				return nil, err
			}

			op3, err := client.DescribeMaintenanceWindowTasks(ctx, &ssm.DescribeMaintenanceWindowTasksInput{
				WindowId: v.WindowId,
			})
			if err != nil {
				return nil, err
			}
			aka := "arn:" + describeCtx.Partition + ":ssm:" + describeCtx.Region + ":" + describeCtx.AccountID + ":maintenancewindow" + "/" + *v.WindowId

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    aka,
				Name:   *v.Name,
				Description: model.SSMMaintenanceWindowDescription{
					ARN:                       aka,
					MaintenanceWindowIdentity: v,
					MaintenanceWindow:         data,
					Tags:                      op.TagList,
					Targets:                   op2.Targets,
					Tasks:                     op3.Tasks,
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

func SSMMaintenanceWindowTarget(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	windows, err := SSMMaintenanceWindow(ctx, cfg, nil)
	if err != nil {
		return nil, err
	}

	client := ssm.NewFromConfig(cfg)

	var values []Resource
	for _, w := range windows {
		window := w.Description.(types.MaintenanceWindowIdentity)
		paginator := ssm.NewDescribeMaintenanceWindowTargetsPaginator(client, &ssm.DescribeMaintenanceWindowTargetsInput{
			WindowId: window.WindowId,
		})

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				return nil, err
			}

			for _, v := range page.Targets {
				resource := Resource{
					Region:      describeCtx.KaytuRegion,
					ID:          *v.WindowTargetId,
					Name:        *v.Name,
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

func SSMMaintenanceWindowTask(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	windows, err := SSMMaintenanceWindow(ctx, cfg, nil)
	if err != nil {
		return nil, err
	}

	client := ssm.NewFromConfig(cfg)

	var values []Resource
	for _, w := range windows {
		window := w.Description.(types.MaintenanceWindowIdentity)
		paginator := ssm.NewDescribeMaintenanceWindowTasksPaginator(client, &ssm.DescribeMaintenanceWindowTasksInput{
			WindowId: window.WindowId,
		})

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				return nil, err
			}

			for _, v := range page.Tasks {
				resource := Resource{
					Region:      describeCtx.KaytuRegion,
					ARN:         *v.TaskArn,
					Name:        *v.Name,
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

func SSMParameter(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ssm.NewFromConfig(cfg)
	paginator := ssm.NewDescribeParametersPaginator(client, &ssm.DescribeParametersInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Parameters {
			op, err := client.GetParameter(ctx, &ssm.GetParameterInput{
				Name:           v.Name,
				WithDecryption: aws.Bool(true),
			})
			if err != nil {
				var ae smithy.APIError
				if errors.As(err, &ae) {
					// In case the KMS key encrypting the SSM Parameter values is disabled, below error is thrown
					// operation error SSM: GetParameter, https response error StatusCode: 400, RequestID: 0965014b-77ab-4847-98d4-2b9e09a68385, InvalidKeyId: arn:aws:kms:us-east-1:111122223333:key/1a2b3c4d-f6b4-4c5b-97e7-123456ab210c is disabled. (Service: AWSKMS; Status Code: 400; Error Code: DisabledException; Request ID: 7b6ae355-c99a-4cad-b2c3-4b40c0abdda9; Proxy: null)
					if ae.ErrorCode() == "InvalidKeyId" {
						op = &ssm.GetParameterOutput{}
					} else {
						return nil, err
					}
				} else {
					return nil, err
				}
			}

			op2, err := client.ListTagsForResource(ctx, &ssm.ListTagsForResourceInput{
				ResourceType: types.ResourceTypeForTagging("Parameter"),
				ResourceId:   v.Name,
			})
			if err != nil {
				return nil, err
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ID:     *v.Name,
				Name:   *v.Name,
				Description: model.SSMParameterDescription{
					ParameterMetadata: v,
					Parameter:         op.Parameter,
					Tags:              op2.TagList,
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

func SSMPatchBaseline(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ssm.NewFromConfig(cfg)
	paginator := ssm.NewDescribePatchBaselinesPaginator(client, &ssm.DescribePatchBaselinesInput{
		Filters: []types.PatchOrchestratorFilter{
			{
				Key:    aws.String("OWNER"),
				Values: []string{"Self"},
			},
		},
	})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.BaselineIdentities {
			aka := "arn:" + describeCtx.Partition + ":ssm:" + describeCtx.Region + ":" + describeCtx.AccountID + ":patchbaseline"
			if strings.HasPrefix(*v.BaselineId, "/") {
				aka = aka + *v.BaselineId
			} else {
				aka = aka + "/" + *v.BaselineId
			}

			data, err := client.GetPatchBaseline(ctx, &ssm.GetPatchBaselineInput{
				BaselineId: v.BaselineId,
			})
			if err != nil {
				return nil, err
			}

			op, err := client.ListTagsForResource(ctx, &ssm.ListTagsForResourceInput{
				ResourceType: "PatchBaseline",
				ResourceId:   v.BaselineId,
			})
			if err != nil {
				return nil, err
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    aka,
				Name:   *v.BaselineName,
				Description: model.SSMPatchBaselineDescription{
					ARN:                   aka,
					PatchBaselineIdentity: v,
					PatchBaseline:         data,
					Tags:                  op.TagList,
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

func SSMResourceDataSync(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ssm.NewFromConfig(cfg)
	paginator := ssm.NewListResourceDataSyncPaginator(client, &ssm.ListResourceDataSyncInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.ResourceDataSyncItems {
			resource := Resource{
				Region:      describeCtx.KaytuRegion,
				ID:          *v.SyncName,
				Name:        *v.SyncName,
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
func SSMManagedInstancePatchState(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := ssm.NewFromConfig(cfg)
	paginator := ssm.NewDescribeInstanceInformationPaginator(client, &ssm.DescribeInstanceInformationInput{})
	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, v := range page.InstanceInformationList {
			paginatorPS := ssm.NewDescribeInstancePatchStatesPaginator(client, &ssm.DescribeInstancePatchStatesInput{
				InstanceIds: []string{*v.InstanceId},
			})

			for paginatorPS.HasMorePages() {
				pagePS, err := paginatorPS.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				for _, item := range pagePS.InstancePatchStates {
					resource := Resource{
						Region: describeCtx.KaytuRegion,
						ID:     *item.InstanceId,
						Description: model.SSMManagedInstancePatchStateDescription{
							PatchState: item,
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
