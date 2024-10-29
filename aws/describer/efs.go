package describer

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/efs/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/efs"
	"github.com/opengovern/og-aws-describer/aws/model"
)

const (
	efsPolicyNotFound = "PolicyNotFound"
)

func EFSAccessPoint(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := efs.NewFromConfig(cfg)
	paginator := efs.NewDescribeAccessPointsPaginator(client, &efs.DescribeAccessPointsInput{})
	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.AccessPoints {
			resource := eFSAccessPointHandle(ctx, v)
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
func eFSAccessPointHandle(ctx context.Context, v types.AccessPointDescription) Resource {
	describeCtx := GetDescribeContext(ctx)
	name := aws.ToString(v.Name)
	if name == "" {
		name = *v.AccessPointId
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *v.AccessPointArn,
		Name:   name,
		Description: model.EFSAccessPointDescription{
			AccessPoint: v,
		},
	}
	return resource
}
func GetEFSAccessPoint(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	accessPointId := fields["id"]
	client := efs.NewFromConfig(cfg)
	out, err := client.DescribeAccessPoints(ctx, &efs.DescribeAccessPointsInput{
		AccessPointId: &accessPointId,
	})
	if err != nil {
		if isErr(err, "DescribeAccessPointsNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	for _, v := range out.AccessPoints {
		resource := eFSAccessPointHandle(ctx, v)
		values = append(values, resource)
	}
	return values, nil
}

func EFSFileSystem(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := efs.NewFromConfig(cfg)
	paginator := efs.NewDescribeFileSystemsPaginator(client, &efs.DescribeFileSystemsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.FileSystems {
			output, err := client.DescribeFileSystemPolicy(ctx, &efs.DescribeFileSystemPolicyInput{
				FileSystemId: v.FileSystemId,
			})
			if err != nil {
				if !isErr(err, efsPolicyNotFound) {
					return nil, err
				}

				output = &efs.DescribeFileSystemPolicyOutput{}
			}
			// Doc: You can add tags to a file system, including a Name tag. For more information,
			// see CreateFileSystem. If the file system has a Name tag, Amazon EFS returns the
			// values in this field.
			resource := eFSFileSystemHandle(ctx, output, v)
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
func eFSFileSystemHandle(ctx context.Context, describeFSPolicy *efs.DescribeFileSystemPolicyOutput, v types.FileSystemDescription) Resource {
	describeCtx := GetDescribeContext(ctx)
	name := aws.ToString(v.Name)
	if name == "" {
		name = *v.FileSystemId
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *v.FileSystemArn,
		Name:   name,
		Description: model.EFSFileSystemDescription{
			FileSystem: v,
			Policy:     describeFSPolicy.Policy,
		},
	}
	return resource
}
func GetEFSFileSystem(ctx context.Context, cfg aws.Config, field map[string]string) ([]Resource, error) {
	fileSystemId := field["id"]
	client := efs.NewFromConfig(cfg)
	var values []Resource
	describeFSPolicy, err := client.DescribeFileSystemPolicy(ctx, &efs.DescribeFileSystemPolicyInput{
		FileSystemId: &fileSystemId,
	})
	if err != nil {
		if !isErr(err, efsPolicyNotFound) {
			return nil, err
		}

		describeFSPolicy = &efs.DescribeFileSystemPolicyOutput{}
	}

	output, err := client.DescribeFileSystems(ctx, &efs.DescribeFileSystemsInput{
		FileSystemId: &fileSystemId,
	})
	if err != nil {
		if isErr(err, "DescribeFileSystemsNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	for _, v := range output.FileSystems {
		resource := eFSFileSystemHandle(ctx, describeFSPolicy, v)
		values = append(values, resource)
	}
	return values, nil
}

func EFSMountTarget(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := efs.NewFromConfig(cfg)

	var values []Resource

	filesystems, err := EFSFileSystem(ctx, cfg, nil)
	if err != nil {
		return nil, err
	}
	for _, fs := range filesystems {
		filesystem := fs.Description.(model.EFSFileSystemDescription)
		err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
			output, err := client.DescribeMountTargets(ctx, &efs.DescribeMountTargetsInput{
				FileSystemId: filesystem.FileSystem.FileSystemId,
				Marker:       prevToken,
			})
			if err != nil {
				return nil, err
			}

			for _, v := range output.MountTargets {
				securityGroups, err := client.DescribeMountTargetSecurityGroups(ctx, &efs.DescribeMountTargetSecurityGroupsInput{
					MountTargetId: v.MountTargetId,
				})
				if err != nil {
					return nil, err
				}
				resource := eFSMountTargetHandle(ctx, securityGroups, v, filesystem.FileSystem.FileSystemId)
				if stream != nil {
					if err := (*stream)(resource); err != nil {
						return nil, err
					}
				} else {
					values = append(values, resource)
				}

			}
			return output.NextMarker, nil
		})
		if err != nil {
			return nil, err
		}
	}

	return values, nil
}
func eFSMountTargetHandle(ctx context.Context, securityGroups *efs.DescribeMountTargetSecurityGroupsOutput, v types.MountTargetDescription, FileSystemId *string) Resource {
	describeCtx := GetDescribeContext(ctx)
	arn := fmt.Sprintf("arn:%s:elasticfilesystem:%s:%s:file-system/%s/mount-target/%s", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID, *FileSystemId, *v.MountTargetId)

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    arn,
		ID:     *v.MountTargetId,
		Description: model.EFSMountTargetDescription{
			MountTarget:    v,
			SecurityGroups: securityGroups.SecurityGroups,
		},
	}
	return resource
}
func GetEFSMountTarget(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	fileSystemId := fields["id"]
	var values []Resource

	client := efs.NewFromConfig(cfg)
	out, err := client.DescribeMountTargets(ctx, &efs.DescribeMountTargetsInput{
		FileSystemId: &fileSystemId,
	})
	if err != nil {
		if isErr(err, "DescribeMountTargetsNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	for _, v := range out.MountTargets {
		securityGroups, err := client.DescribeMountTargetSecurityGroups(ctx, &efs.DescribeMountTargetSecurityGroupsInput{
			MountTargetId: v.MountTargetId,
		})
		if err != nil {
			if isErr(err, "DescribeMountTargetSecurityGroupsNotFound") || isErr(err, "InvalidParameterValue") {
				return nil, nil
			}
			return nil, err
		}

		resource := eFSMountTargetHandle(ctx, securityGroups, v, &fileSystemId)
		values = append(values, resource)
	}
	return values, nil
}
