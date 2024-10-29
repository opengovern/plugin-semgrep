package describer

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/codecommit/types"
	"math"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codecommit"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func CodeCommitRepository(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := codecommit.NewFromConfig(cfg)
	paginator := codecommit.NewListRepositoriesPaginator(client, &codecommit.ListRepositoriesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if !isErr(err, "InvalidParameter") {
				return nil, err
			}
			continue
		}
		var repositoryNames []string
		for _, v := range page.Repositories {
			repositoryNames = append(repositoryNames, *v.RepositoryName)
		}
		if len(repositoryNames) == 0 {
			continue
		}
		// BatchGetRepositories can only get 25 repositories at a time
		for i := 0; i < len(repositoryNames); i += 25 {
			repos, err := client.BatchGetRepositories(ctx, &codecommit.BatchGetRepositoriesInput{
				RepositoryNames: repositoryNames[i:int(math.Min(float64(i+25), float64(len(repositoryNames))))],
			})
			if err != nil {
				if !isErr(err, "InvalidParameter") {
					return nil, err
				}
				continue
			}
			for _, v := range repos.Repositories {
				resource, err := codeCommitRepositoryHandle(ctx, cfg, v)
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
	}

	return values, nil
}
func codeCommitRepositoryHandle(ctx context.Context, cfg aws.Config, v types.RepositoryMetadata) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := codecommit.NewFromConfig(cfg)

	tags, err := client.ListTagsForResource(ctx, &codecommit.ListTagsForResourceInput{
		ResourceArn: v.Arn,
	})
	if err != nil {
		if !isErr(err, "InvalidParameter") {
			return Resource{}, err
		}
		tags = &codecommit.ListTagsForResourceOutput{}
	}

	resource := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *v.Arn,
		Name:   *v.RepositoryName,
		Description: model.CodeCommitRepositoryDescription{
			Repository: v,
			Tags:       tags.Tags,
		},
	}
	return resource, nil
}
func GetCodeCommitRepository(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	repositoryName := fields["repositoryName"]
	var values []Resource

	client := codecommit.NewFromConfig(cfg)
	repos, err := client.BatchGetRepositories(ctx, &codecommit.BatchGetRepositoriesInput{
		RepositoryNames: []string{repositoryName},
	})
	if err != nil {
		if !isErr(err, "InvalidParameter") {
			return nil, err
		}
		return nil, nil
	}

	for _, v := range repos.Repositories {
		resource, err := codeCommitRepositoryHandle(ctx, cfg, v)
		if err != nil {
			return nil, err
		}
		values = append(values, resource)
	}
	return values, nil
}
