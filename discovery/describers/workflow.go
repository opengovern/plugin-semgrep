package describers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/opengovern/og-describer-template/discovery/pkg/models"
	model "github.com/opengovern/og-describer-template/discovery/provider"

	goPipeline "github.com/buildkite/go-pipeline"

	"github.com/google/go-github/v55/github"
)

func GetAllWorkflows(ctx context.Context, githubClient model.GitHubClient, organizationName string, stream *models.StreamSender) ([]models.Resource, error) {
	client := githubClient.RestClient

	repositories, err := getRepositories(ctx, client, organizationName)
	if err != nil {
		return nil, err
	}
	var values []models.Resource
	for _, repo := range repositories {
		if repo == nil {
			continue
		}
		repoValues, err := GetRepositoryWorkflows(ctx, githubClient, stream, organizationName, repo.GetName())
		if err != nil {
			return nil, err
		}
		values = append(values, repoValues...)
	}
	return values, nil
}

type FileContent struct {
	Repository string
	FilePath   string
	Content    string
}

func GetRepositoryWorkflows(ctx context.Context, githubClient model.GitHubClient, stream *models.StreamSender, owner, repo string) ([]models.Resource, error) {
	client := githubClient.RestClient
	opts := &github.ListOptions{PerPage: pageSize}
	repoFullName := formRepositoryFullName(owner, repo)
	var values []models.Resource
	for {
		workflows, resp, err := client.Actions.ListWorkflows(ctx, owner, repo, opts)
		if err != nil {
			return nil, err
		}
		for _, workflow := range workflows.Workflows {
			if workflow == nil {
				continue
			}
			var pipeline *goPipeline.Pipeline
			var content string

			fileContent, err := getWorkflowFileContent(ctx, client, workflow, owner, repo)
			if err != nil {
				return nil, err
			}
			if fileContent != nil {
				content, err = fileContent.GetContent()
				if err != nil {
					return nil, err
				}
				fileContentBasic := FileContent{
					Repository: repo,
					FilePath:   fileContent.GetPath(),
					Content:    content,
				}
				pipeline, err = decodeFileContentToPipeline(fileContentBasic)
				if err != nil {
					continue
				}
			}

			value := models.Resource{
				ID:   strconv.Itoa(int(workflow.GetID())),
				Name: workflow.GetName(),
				Description: model.WorkflowDescription{
					ID:                      workflow.ID,
					NodeID:                  workflow.NodeID,
					Name:                    workflow.Name,
					Path:                    workflow.Path,
					State:                   workflow.State,
					CreatedAt:               workflow.CreatedAt,
					UpdatedAt:               workflow.UpdatedAt,
					URL:                     workflow.URL,
					HTMLURL:                 workflow.HTMLURL,
					BadgeURL:                workflow.BadgeURL,
					RepositoryFullName:      &repoFullName,
					WorkFlowFileContent:     &content,
					WorkFlowFileContentJson: fileContent,
					Pipeline:                pipeline,
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

func getWorkflowFileContent(ctx context.Context, client *github.Client, workflow *github.Workflow, owner, repo string) (*github.RepositoryContent, error) {
	if workflow.Path == nil {
		return nil, nil
	}
	workflowUrlParts := strings.Split(*workflow.HTMLURL, "/")
	defaultBranch := "main"
	if len(workflowUrlParts) > 6 {
		defaultBranch = workflowUrlParts[6]
	}
	content, _, _, err := client.Repositories.GetContents(ctx, owner, repo, workflow.GetPath(), &github.RepositoryContentGetOptions{Ref: defaultBranch})
	if err != nil {
		if strings.Contains(err.Error(), "404 Not Found") || strings.Contains(err.Error(), "404 No commit found") {
			return nil, nil
		}
		return nil, err
	}
	return content, nil
}

func decodeFileContentToPipeline(contentDetails FileContent) (*goPipeline.Pipeline, error) {
	pipeline, err := goPipeline.Parse(strings.NewReader(contentDetails.Content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse the workflow file '%s', %v", contentDetails.FilePath, err)
	}
	return pipeline, nil
}

func GetRepositoryWorkflow(ctx context.Context, githubClient model.GitHubClient, organizationName string, repositoryName string, resourceID string, stream *models.StreamSender) (*models.Resource, error) {
	client := githubClient.RestClient
	workflowID, err := strconv.ParseInt(resourceID, 10, 64)
	if err != nil {
		return nil, err
	}
	repoFullName := formRepositoryFullName(organizationName, repositoryName)
	workflow, _, err := client.Actions.GetWorkflowByID(ctx, organizationName, repositoryName, workflowID)
	if err != nil {
		return nil, err
	}
	fileContent, err := getWorkflowFileContent(ctx, client, workflow, organizationName, repositoryName)
	if err != nil {
		return nil, err
	}
	content, err := fileContent.GetContent()
	if err != nil {
		return nil, err
	}
	fileContentBasic := FileContent{
		Repository: repositoryName,
		FilePath:   fileContent.GetPath(),
		Content:    content,
	}
	pipeline, err := decodeFileContentToPipeline(fileContentBasic)
	if err != nil {
		return nil, err
	}
	value := models.Resource{
		ID:   strconv.Itoa(int(workflow.GetID())),
		Name: workflow.GetName(),
		Description: model.WorkflowDescription{
			ID:                      workflow.ID,
			NodeID:                  workflow.NodeID,
			Name:                    workflow.Name,
			Path:                    workflow.Path,
			State:                   workflow.State,
			CreatedAt:               workflow.CreatedAt,
			UpdatedAt:               workflow.UpdatedAt,
			URL:                     workflow.URL,
			HTMLURL:                 workflow.HTMLURL,
			BadgeURL:                workflow.BadgeURL,
			RepositoryFullName:      &repoFullName,
			WorkFlowFileContent:     &content,
			WorkFlowFileContentJson: fileContent,
			Pipeline:                pipeline,
		},
	}
	if stream != nil {
		if err := (*stream)(value); err != nil {
			return nil, err
		}
	}

	return &value, nil
}
