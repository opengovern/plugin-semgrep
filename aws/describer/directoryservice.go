package describer

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/smithy-go"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/directoryservice"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func DirectoryServiceDirectory(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)

	client := directoryservice.NewFromConfig(cfg)
	paginator := directoryservice.NewDescribeDirectoriesPaginator(client, &directoryservice.DescribeDirectoriesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if !isErr(err, "InvalidParameterValueException") && !isErr(err, "ResourceNotFoundFault") && !isErr(err, "EntityDoesNotExistException") {
				return nil, err
			}
			continue
		}

		for _, v := range page.DirectoryDescriptions {
			arn := fmt.Sprintf("arn:%s:ds:%s:%s:directory/%s", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID, *v.DirectoryId)

			tags, err := client.ListTagsForResource(ctx, &directoryservice.ListTagsForResourceInput{
				ResourceId: v.DirectoryId,
			})
			if err != nil {
				if !isErr(err, "InvalidParameterValueException") && !isErr(err, "ResourceNotFoundFault") && !isErr(err, "EntityDoesNotExistException") {
					return nil, err
				}
				tags = &directoryservice.ListTagsForResourceOutput{}
			}
			sharedDirectory, err := client.DescribeSharedDirectories(ctx, &directoryservice.DescribeSharedDirectoriesInput{
				OwnerDirectoryId: v.DirectoryId,
			})
			if err != nil {
				if isErr(err, "DescribeSharedDirectoriesNotFound") || isErr(err, "InvalidParameterValue") {
					return nil, nil
				}
				return nil, err
			}

			snapshot, err := client.GetSnapshotLimits(ctx, &directoryservice.GetSnapshotLimitsInput{
				DirectoryId: v.DirectoryId,
			})
			if err != nil {
				if isErr(err, "GetSnapshotLimitsNotFound") || isErr(err, "InvalidParameterValue") {
					return nil, nil
				}
				return nil, err
			}

			eventTopic, err := client.DescribeEventTopics(ctx, &directoryservice.DescribeEventTopicsInput{
				DirectoryId: v.DirectoryId,
			})
			if err != nil {
				if isErr(err, "DescribeEventTopicsNotFound") || isErr(err, "InvalidParameterValue") {
					return nil, nil
				}
				return nil, err
			}
			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    arn,
				Name:   *v.Name,
				Description: model.DirectoryServiceDirectoryDescription{
					Directory:       v,
					Snapshot:        *snapshot.SnapshotLimits,
					EventTopics:     eventTopic.EventTopics,
					SharedDirectory: sharedDirectory.SharedDirectories,
					Tags:            tags.Tags,
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

func DirectoryServiceCertificate(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)

	client := directoryservice.NewFromConfig(cfg)
	paginator := directoryservice.NewDescribeDirectoriesPaginator(client, &directoryservice.DescribeDirectoriesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if !isErr(err, "InvalidParameterValueException") && !isErr(err, "ResourceNotFoundFault") && !isErr(err, "EntityDoesNotExistException") {
				return nil, err
			}
			continue
		}

		for _, v := range page.DirectoryDescriptions {
			certiPaginator := directoryservice.NewListCertificatesPaginator(client, &directoryservice.ListCertificatesInput{
				DirectoryId: v.DirectoryId,
			})
			if err != nil {
				var ae smithy.APIError
				if errors.As(err, &ae) {
					if ae.ErrorCode() == "UnsupportedOperationException" {
						return nil, nil
					}
				}
				if !isErr(err, "InvalidParameterValueException") && !isErr(err, "ResourceNotFoundFault") && !isErr(err, "EntityDoesNotExistException") {
					return nil, err
				}
			}
			for certiPaginator.HasMorePages() {
				certiPage, err := certiPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}
				for _, certi := range certiPage.CertificatesInfo {
					certificate, err := client.DescribeCertificate(ctx, &directoryservice.DescribeCertificateInput{
						CertificateId: certi.CertificateId,
						DirectoryId:   v.DirectoryId,
					})
					if err != nil {
						return nil, err
					}
					resource := Resource{
						Region: describeCtx.KaytuRegion,
						ID:     *certificate.Certificate.CertificateId,
						Name:   *certificate.Certificate.CommonName,
						Description: model.DirectoryServiceCertificateDescription{
							Certificate: *certificate.Certificate,
							DirectoryId: *v.DirectoryId,
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

func DirectoryServiceLogSubscription(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)

	client := directoryservice.NewFromConfig(cfg)
	paginator := directoryservice.NewDescribeDirectoriesPaginator(client, &directoryservice.DescribeDirectoriesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if !isErr(err, "InvalidParameterValueException") && !isErr(err, "ResourceNotFoundFault") && !isErr(err, "EntityDoesNotExistException") {
				return nil, err
			}
			continue
		}

		for _, v := range page.DirectoryDescriptions {
			logPaginator := directoryservice.NewListLogSubscriptionsPaginator(client, &directoryservice.ListLogSubscriptionsInput{
				DirectoryId: v.DirectoryId,
			})
			if err != nil {
				if !isErr(err, "InvalidParameterValueException") && !isErr(err, "ResourceNotFoundFault") && !isErr(err, "EntityDoesNotExistException") {
					return nil, err
				}
			}
			for logPaginator.HasMorePages() {
				logPage, err := logPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}
				for _, logSub := range logPage.LogSubscriptions {
					resource := Resource{
						Region: describeCtx.KaytuRegion,
						Name:   *logSub.LogGroupName,
						Description: model.DirectoryServiceLogSubscriptionDescription{
							LogSubscription: logSub,
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
