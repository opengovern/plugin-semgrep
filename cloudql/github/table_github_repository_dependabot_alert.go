package github

import (
	opengovernance "github.com/opengovern/og-describer-template/discovery/pkg/es"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubRepositoryDependabotAlert() *plugin.Table {
	return &plugin.Table{
		Name:        "github_repository_dependabot_alert",
		Description: "Dependabot alerts from a repository.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListRepoAlertDependabot,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"repository_full_name", "alert_number"}),
			ShouldIgnoreError: isNotFoundError([]string{"404", "403"}),
			Hydrate:           opengovernance.GetRepoAlertDependabot,
		},
		Columns: commonColumns(append(
			gitHubDependabotAlertColumns(),
			[]*plugin.Column{
				{
					Name:        "repository_full_name",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromQual("Description.RepoFullName"),
					Description: "The full name of the repository (login/repo-name).",
				},
			}...,
		)),
	}
}
