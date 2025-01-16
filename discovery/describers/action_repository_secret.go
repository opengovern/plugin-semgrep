package describers

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v55/github"
	"github.com/opengovern/og-describer-template/discovery/pkg/models"
	model "github.com/opengovern/og-describer-template/discovery/provider"
)

func GetAllSecrets(ctx context.Context, githubClient model.GitHubClient, organizationName string, stream *models.StreamSender) ([]models.Resource, error) {
	client := githubClient.RestClient
	owner := organizationName
	repositories, err := getRepositories(ctx, client, owner)
	if err != nil {
		return nil, nil
	}
	var values []models.Resource
	for _, repo := range repositories {
		repoValues, err := GetRepositorySecrets(ctx, githubClient, stream, owner, repo.GetName())
		if err != nil {
			return nil, err
		}
		values = append(values, repoValues...)
	}
	return values, nil
}

func GetRepositorySecrets(ctx context.Context, githubClient model.GitHubClient, stream *models.StreamSender, owner, repo string) ([]models.Resource, error) {
	client := githubClient.RestClient
	opts := &github.ListOptions{PerPage: maxPagesCount}
	repoFullName := formRepositoryFullName(owner, repo)
	var values []models.Resource
	for {
		secrets, resp, err := client.Actions.ListRepoSecrets(ctx, owner, repo, opts)
		if err != nil {
			return nil, err
		}
		for _, secret := range secrets.Secrets {
			id := fmt.Sprintf("%s/%s/%s", owner, repo, secret.Name)
			createdAt := secret.CreatedAt.Format(time.RFC3339)
			updatedAt := secret.UpdatedAt.Format(time.RFC3339)
			value := models.Resource{
				ID:   id,
				Name: secret.Name,
				Description: model.SecretDescription{
					Name:                    &secret.Name,
					CreatedAt:               &createdAt,
					UpdatedAt:               &updatedAt,
					Visibility:              &secret.Visibility,
					SelectedRepositoriesURL: &secret.SelectedRepositoriesURL,
					RepoFullName:            &repoFullName,
				},
			}
			if stream != nil {
				if err := (*stream)(value); err != nil {
					return nil, err
				}
			} else {
				values = append(values, value)
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return values, nil
}

func GetRepoActionSecret(ctx context.Context, githubClient model.GitHubClient, organizationName string, repositoryName string, resourceID string, stream *models.StreamSender) (*models.Resource, error) {
	client := githubClient.RestClient
	repoFullName := formRepositoryFullName(organizationName, repositoryName)
	secret, _, err := client.Actions.GetRepoSecret(ctx, organizationName, repositoryName, resourceID)
	if err != nil {
		return nil, err
	}
	id := fmt.Sprintf("%s/%s/%s", organizationName, repositoryName, secret.Name)
	createdAt := secret.CreatedAt.Format(time.RFC3339)
	updatedAt := secret.UpdatedAt.Format(time.RFC3339)
	value := models.Resource{
		ID:   id,
		Name: secret.Name,
		Description: model.SecretDescription{
			Name:                    &secret.Name,
			CreatedAt:               &createdAt,
			UpdatedAt:               &updatedAt,
			Visibility:              &secret.Visibility,
			SelectedRepositoriesURL: &secret.SelectedRepositoriesURL,
			RepoFullName:            &repoFullName,
		},
	}
	if stream != nil {
		if err := (*stream)(value); err != nil {
			return nil, err
		}
	}

	return &value, nil
}
