package describers

import (
	"context"
	"time"

	"github.com/google/go-github/v55/github"
	"github.com/opengovern/og-describer-template/discovery/pkg/models"
	model "github.com/opengovern/og-describer-template/discovery/provider"
)

func GetAllAuditLogs(ctx context.Context, githubClient model.GitHubClient, organizationName string, stream *models.StreamSender) ([]models.Resource, error) {
	var values []models.Resource
	orgValues, err := GetRepositoryAuditLog(ctx, githubClient, stream, organizationName)
	if err != nil {
		return nil, err
	}
	values = append(values, orgValues...)
	return values, nil
}

func GetRepositoryAuditLog(ctx context.Context, githubClient model.GitHubClient, stream *models.StreamSender, org string) ([]models.Resource, error) {
	client := githubClient.RestClient
	var phrase string
	var include string
	opts := &github.GetAuditLogOptions{
		Phrase:            &phrase,
		Include:           &include,
		ListCursorOptions: github.ListCursorOptions{PerPage: 100},
	}
	var values []models.Resource
	for {
		auditResults, resp, err := client.Organizations.GetAuditLog(ctx, org, opts)
		if err != nil {
			return nil, err
		}
		for _, audit := range auditResults {
			createdAt := audit.CreatedAt.Format(time.RFC3339)
			actorLocation := model.ActorLocation{
				CountryCode: audit.ActorLocation.CountryCode,
			}
			data := model.AuditEntryData{
				OldName:  audit.Data.OldName,
				OldLogin: audit.Data.OldLogin,
			}
			value := models.Resource{
				ID:   audit.GetDocumentID(),
				Name: audit.GetName(),
				Description: model.AuditLogDescription{
					ID:            audit.DocumentID,
					CreatedAt:     &createdAt,
					Organization:  &org,
					Phrase:        &phrase,
					Include:       &include,
					Action:        audit.Action,
					Actor:         audit.Actor,
					ActorLocation: &actorLocation,
					Team:          audit.Team,
					UserLogin:     audit.User,
					Repo:          audit.Repository,
					Data:          &data,
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
		if resp.After == "" {
			break
		}
		opts.After = resp.After
	}
	return values, nil
}
