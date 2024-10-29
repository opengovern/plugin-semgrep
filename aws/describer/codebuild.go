package describer

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/codebuild/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codebuild"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func CodeBuildProject(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := codebuild.NewFromConfig(cfg)
	paginator := codebuild.NewListProjectsPaginator(client, &codebuild.ListProjectsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		if len(page.Projects) == 0 {
			continue
		}

		projects, err := client.BatchGetProjects(ctx, &codebuild.BatchGetProjectsInput{
			Names: page.Projects,
		})
		if err != nil {
			return nil, err
		}

		for _, project := range projects.Projects {

			resource := codeBuildProjectHandle(ctx, project)
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
func codeBuildProjectHandle(ctx context.Context, project types.Project) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *project.Arn,
		Name:   *project.Name,
		Description: model.CodeBuildProjectDescription{
			Project: project,
		},
	}
	return resource
}
func GetCodeBuildProject(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	name := fields["name"]
	client := codebuild.NewFromConfig(cfg)

	out, err := client.BatchGetProjects(ctx, &codebuild.BatchGetProjectsInput{
		Names: []string{name},
	})
	if err != nil {
		if isErr(err, "BatchGetProjectsNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	for _, project := range out.Projects {
		resource := codeBuildProjectHandle(ctx, project)
		values = append(values, resource)
	}
	return values, nil
}

func CodeBuildSourceCredential(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := codebuild.NewFromConfig(cfg)
	out, err := client.ListSourceCredentials(ctx, &codebuild.ListSourceCredentialsInput{})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, item := range out.SourceCredentialsInfos {
		resource := codeBuildSourceCredentialHandle(ctx, item)
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
func codeBuildSourceCredentialHandle(ctx context.Context, item types.SourceCredentialsInfo) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *item.Arn,
		Name:   nameFromArn(*item.Arn),
		Description: model.CodeBuildSourceCredentialDescription{
			SourceCredentialsInfo: item,
		},
	}
	return resource
}
func GetCodeBuildSourceCredential(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	arn := fields["arn"]
	client := codebuild.NewFromConfig(cfg)
	credentials, err := client.ListSourceCredentials(ctx, &codebuild.ListSourceCredentialsInput{})
	if err != nil {
		if isErr(err, "ListSourceCredentialsNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	for _, item := range credentials.SourceCredentialsInfos {

		if *item.Arn != arn {
			continue
		}
		resource := codeBuildSourceCredentialHandle(ctx, item)
		values = append(values, resource)

	}
	return values, nil
}

func CodeBuildBuild(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := codebuild.NewFromConfig(cfg)
	paginator := codebuild.NewListBuildsPaginator(client, &codebuild.ListBuildsInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		if len(page.Ids) == 0 {
			continue
		}

		build, err := client.BatchGetBuilds(ctx, &codebuild.BatchGetBuildsInput{
			Ids: page.Ids,
		})
		if err != nil {
			return nil, err
		}

		for _, project := range build.Builds {

			resource := codeBuildBuildHandle(ctx, project)
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
func codeBuildBuildHandle(ctx context.Context, build types.Build) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *build.Arn,
		ID:     *build.Id,
		Description: model.CodeBuildBuildDescription{
			Build: build,
		},
	}
	return resource
}
func GetCodeBuildBuild(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	id := fields["id"]
	client := codebuild.NewFromConfig(cfg)

	out, err := client.BatchGetBuilds(ctx, &codebuild.BatchGetBuildsInput{
		Ids: []string{id},
	})
	if err != nil {
		if isErr(err, "BatchGetBuildsNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	for _, build := range out.Builds {
		resource := codeBuildBuildHandle(ctx, build)
		values = append(values, resource)
	}
	return values, nil
}
