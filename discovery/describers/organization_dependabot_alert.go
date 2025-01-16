package describers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/go-github/v55/github"
	"github.com/opengovern/og-describer-template/discovery/pkg/models"
	model "github.com/opengovern/og-describer-template/discovery/provider"
)

func GetAllOrganizationsDependabotAlerts(ctx context.Context, githubClient model.GitHubClient, organizationName string, stream *models.StreamSender) ([]models.Resource, error) {

	var values []models.Resource
	orgValues, err := GetOrganizationDependabotAlerts(ctx, githubClient, stream, organizationName)
	if err != nil {
		return nil, err
	}
	values = append(values, orgValues...)
	return values, nil
}

func GetOrganizationDependabotAlerts(ctx context.Context, githubClient model.GitHubClient, stream *models.StreamSender, org string) ([]models.Resource, error) {
	client := githubClient.RestClient
	opt := &github.ListAlertsOptions{
		ListCursorOptions: github.ListCursorOptions{First: maxPagesCount},
	}
	var values []models.Resource
	for {
		alerts, resp, err := client.Dependabot.ListOrgAlerts(ctx, org, opt)
		if err != nil {
			return nil, err
		}
		for _, alert := range alerts {
			var CWEs []string
			for _, cwe := range alert.SecurityAdvisory.CWEs {
				CWEs = append(CWEs, cwe.GetName())
			}
			id := fmt.Sprintf("%s/%s", org, strconv.Itoa(alert.GetNumber()))
			value := models.Resource{
				ID:   id,
				Name: strconv.Itoa(alert.GetNumber()),
				Description: model.OrgAlertDependabotDescription{
					AlertNumber:                 alert.GetNumber(),
					State:                       alert.GetState(),
					DependencyPackageEcosystem:  alert.GetDependency().GetPackage().GetEcosystem(),
					DependencyPackageName:       alert.GetDependency().GetPackage().GetName(),
					DependencyManifestPath:      alert.GetDependency().GetManifestPath(),
					DependencyScope:             alert.GetDependency().GetScope(),
					SecurityAdvisoryGHSAID:      alert.GetSecurityAdvisory().GetGHSAID(),
					SecurityAdvisoryCVEID:       alert.GetSecurityAdvisory().GetCVEID(),
					SecurityAdvisorySummary:     alert.GetSecurityAdvisory().GetSummary(),
					SecurityAdvisoryDescription: alert.GetSecurityAdvisory().GetDescription(),
					SecurityAdvisorySeverity:    alert.GetSecurityAdvisory().GetSeverity(),
					SecurityAdvisoryCVSSScore:   alert.GetSecurityAdvisory().GetCVSS().GetScore(),
					SecurityAdvisoryCVSSVector:  alert.GetSecurityAdvisory().GetCVSS().GetVectorString(),
					SecurityAdvisoryCWEs:        CWEs,
					SecurityAdvisoryPublishedAt: alert.GetSecurityAdvisory().GetPublishedAt(),
					SecurityAdvisoryUpdatedAt:   alert.GetSecurityAdvisory().GetUpdatedAt(),
					SecurityAdvisoryWithdrawnAt: alert.GetSecurityAdvisory().GetWithdrawnAt(),
					URL:                         alert.GetURL(),
					HTMLURL:                     alert.GetHTMLURL(),
					CreatedAt:                   alert.GetCreatedAt(),
					UpdatedAt:                   alert.GetUpdatedAt(),
					DismissedAt:                 alert.GetDismissedAt(),
					DismissedReason:             alert.GetDismissedReason(),
					DismissedComment:            alert.GetDismissedComment(),
					FixedAt:                     alert.GetFixedAt(),
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
		opt.ListCursorOptions.After = resp.After
	}
	return values, nil
}
