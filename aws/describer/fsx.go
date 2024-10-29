package describer

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/fsx"
	"github.com/aws/aws-sdk-go-v2/service/fsx/types"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func FSXFileSystem(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := fsx.NewFromConfig(cfg)
	paginator := fsx.NewDescribeFileSystemsPaginator(client, &fsx.DescribeFileSystemsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range page.FileSystems {
			resource := fSXFileSystemHandle(ctx, item)
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
func fSXFileSystemHandle(ctx context.Context, item types.FileSystem) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *item.ResourceARN,
		Name:   *item.FileSystemId,
		Description: model.FSXFileSystemDescription{
			FileSystem: item,
		},
	}
	return resource
}
func GetFSXFileSystem(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	filesystemId := fields["id"]
	client := fsx.NewFromConfig(cfg)

	filesystem, err := client.DescribeFileSystems(ctx, &fsx.DescribeFileSystemsInput{
		FileSystemIds: []string{filesystemId},
	})
	if err != nil {
		if isErr(err, "DescribeFileSystemsNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	for _, v := range filesystem.FileSystems {
		values = append(values, fSXFileSystemHandle(ctx, v))
	}
	return values, nil
}

func FSXStorageVirtualMachine(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := fsx.NewFromConfig(cfg)
	paginator := fsx.NewDescribeStorageVirtualMachinesPaginator(client, &fsx.DescribeStorageVirtualMachinesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range page.StorageVirtualMachines {
			resource := fSXStorageVirtualMachineHandle(ctx, item)
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
func fSXStorageVirtualMachineHandle(ctx context.Context, item types.StorageVirtualMachine) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *item.ResourceARN,
		Name:   *item.Name,
		Description: model.FSXStorageVirtualMachineDescription{
			StorageVirtualMachine: item,
		},
	}
	return resource
}
func GetFSXStorageVirtualMachine(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	storageVirtualMachineId := fields["id"]

	client := fsx.NewFromConfig(cfg)
	out, err := client.DescribeStorageVirtualMachines(ctx, &fsx.DescribeStorageVirtualMachinesInput{
		StorageVirtualMachineIds: []string{storageVirtualMachineId},
	})
	if err != nil {
		if isErr(err, "DescribeStorageVirtualMachinesNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	for _, v := range out.StorageVirtualMachines {
		values = append(values, fSXStorageVirtualMachineHandle(ctx, v))
	}

	return values, nil
}

func FSXTask(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := fsx.NewFromConfig(cfg)
	paginator := fsx.NewDescribeDataRepositoryTasksPaginator(client, &fsx.DescribeDataRepositoryTasksInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range page.DataRepositoryTasks {
			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    *item.ResourceARN,
				Name:   *item.TaskId,
				Description: model.FSXTaskDescription{
					Task: item,
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

func FSXVolume(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := fsx.NewFromConfig(cfg)
	paginator := fsx.NewDescribeVolumesPaginator(client, &fsx.DescribeVolumesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range page.Volumes {
			resource := fSXVolumeHandle(ctx, item)
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
func fSXVolumeHandle(ctx context.Context, item types.Volume) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *item.ResourceARN,
		Name:   *item.Name,
		Description: model.FSXVolumeDescription{
			Volume: item,
		},
	}
	return resource
}
func GetFSXVolume(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	volumeId := fields["id"]
	client := fsx.NewFromConfig(cfg)

	volumes, err := client.DescribeVolumes(ctx, &fsx.DescribeVolumesInput{
		VolumeIds: []string{volumeId},
	})
	if err != nil {
		if isErr(err, "DescribeVolumesNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	for _, item := range volumes.Volumes {
		values = append(values, fSXVolumeHandle(ctx, item))
	}
	return values, nil
}

func FSXSnapshot(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := fsx.NewFromConfig(cfg)
	paginator := fsx.NewDescribeSnapshotsPaginator(client, &fsx.DescribeSnapshotsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range page.Snapshots {
			resource := fSXSnapshotHandle(ctx, item)
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
func fSXSnapshotHandle(ctx context.Context, item types.Snapshot) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *item.ResourceARN,
		Name:   *item.Name,
		Description: model.FSXSnapshotDescription{
			Snapshot: item,
		},
	}
	return resource
}
func GetFSXSnapshot(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	snapshotId := fields["id"]
	client := fsx.NewFromConfig(cfg)

	snapshot, err := client.DescribeSnapshots(ctx, &fsx.DescribeSnapshotsInput{
		SnapshotIds: []string{snapshotId},
	})
	if err != nil {
		if isErr(err, "DescribeSnapshotsNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	for _, item := range snapshot.Snapshots {
		values = append(values, fSXSnapshotHandle(ctx, item))
	}
	return values, nil
}
