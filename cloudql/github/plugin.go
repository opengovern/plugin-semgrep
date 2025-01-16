package github

import (
	"context"

	essdk "github.com/opengovern/og-util/pkg/opengovernance-es-sdk"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Plugin returns this plugin
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-github",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: essdk.ConfigInstance,
			Schema:      essdk.ConfigSchema(),
		},
		DefaultTransform: transform.FromCamel(),
		TableMap: map[string]*plugin.Table{
			"github_actions_artifact":     tableGitHubActionsArtifact(),
			"github_actions_runner":       tableGitHubActionsRepositoryRunner(),
			"github_actions_secret":       tableGitHubActionsRepositorySecret(),
			"github_actions_workflow_run": tableGitHubActionsRepositoryWorkflowRun(),
			//"github_blob":                 tableGitHubBlob(),
			"github_branch":            tableGitHubBranch(),
			"github_branch_protection": tableGitHubBranchProtection(),
			//"github_code_owner":                      tableGitHubCodeOwner(),
			"github_commit":                        tableGitHubCommit(),
			"github_issue":                         tableGitHubIssue(),
			"github_license":                       tableGitHubLicense(),
			"github_organization":                  tableGitHubOrganization(),
			"github_organization_dependabot_alert": tableGitHubOrganizationDependabotAlert(),
			//"github_organization_external_identity": tableGitHubOrganizationExternalIdentity(),
			"github_organization_member":       tableGitHubOrganizationMember(),
			"github_organization_collaborator": tableGitHubOrganizationCollaborator(),
			"github_organization_team":         tableGitHubOrganizationTeam(),
			"github_pull_request":              tableGitHubPullRequest(),
			"github_pull_request_review":       tableGitHubPullRequestReview(),
			//"github_rate_limit":                     tableGitHubRateLimit(),
			//"github_rate_limit_graphql":             tableGitHubRateLimitGraphQL(),
			"github_release":                 tableGitHubRelease(),
			"github_repository":              tableGitHubRepository(),
			"github_repository_collaborator": tableGitHubRepositoryCollaborator(),
			// "github_repository_content":              tableGitHubRepositoryContent(),
			"github_repository_dependabot_alert":    tableGitHubRepositoryDependabotAlert(),
			"github_repository_deployment":          tableGitHubRepositoryDeployment(),
			"github_repository_environment":         tableGitHubRepositoryEnvironment(),
			"github_repository_ruleset":             tableGitHubRepositoryRuleset(),
			"github_repository_sbom":                tableGitHubRepositorySbom(),
			"github_repository_vulnerability_alert": tableGitHubRepositoryVulnerabilityAlert(),
			"github_tag":                            tableGitHubTag(),
			"github_team_member":                    tableGitHubTeamMember(),
			//"github_tree":                           tableGitHubTree(),
			"github_user":                tableGitHubUser(),
			"github_workflow":            tableGitHubWorkflow(),
			"github_container_package":   tableGitHubContainerPackage(),
			"github_maven_package":       tableGitHubMavenPackage(),
			"github_npm_package":         tableGitHubNPMPackage(),
			"github_nuget_package":       tableGitHubNugetPackage(),
			"github_artifact_dockerfile": tableGitHubArtifactDockerFile(),
		},
	}
	for key, table := range p.TableMap {
		if table == nil {
			continue
		}
		if table.Get != nil && table.Get.Hydrate == nil {
			delete(p.TableMap, key)
			continue
		}
		if table.List != nil && table.List.Hydrate == nil {
			delete(p.TableMap, key)
			continue
		}

		opengovernanceTable := false
		for _, col := range table.Columns {
			if col != nil && col.Name == "platform_integration_id" {
				opengovernanceTable = true
			}
		}

		if opengovernanceTable {
			if table.Get != nil {
				table.Get.KeyColumns = append(table.Get.KeyColumns, plugin.OptionalColumns([]string{"platform_integration_id", "platform_resource_id"})...)
			}

			if table.List != nil {
				table.List.KeyColumns = append(table.List.KeyColumns, plugin.OptionalColumns([]string{"platform_integration_id", "platform_resource_id"})...)
			}
		}
	}
	return p
}
