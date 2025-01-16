package maps
import (
	"github.com/opengovern/og-describer-template/discovery/describers"
	"github.com/opengovern/og-describer-template/discovery/provider"
	"github.com/opengovern/og-describer-template/platform/constants"
	"github.com/opengovern/og-util/pkg/integration/interfaces"
	model "github.com/opengovern/og-describer-template/discovery/pkg/models"
)
var ResourceTypes = map[string]model.ResourceType{

	"Github/Actions/Artifact": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Github/Actions/Artifact",
		Tags:                 map[string][]string{
            "category": {"Action"},
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByGithub(describers.GetAllArtifacts),
		GetDescriber:         provider.DescribeSingleByRepo(describers.GetArtifact),
	},

	"Github/Actions/Runner": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Github/Actions/Runner",
		Tags:                 map[string][]string{
            "category": {"Action"},
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByGithub(describers.GetAllRunners),
		GetDescriber:         provider.DescribeSingleByRepo(describers.GetActionRunner),
	},

	"Github/Actions/Secret": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Github/Actions/Secret",
		Tags:                 map[string][]string{
            "category": {"Action"},
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByGithub(describers.GetAllSecrets),
		GetDescriber:         provider.DescribeSingleByRepo(describers.GetRepoActionSecret),
	},

	"Github/Actions/WorkflowRun": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Github/Actions/WorkflowRun",
		Tags:                 map[string][]string{
            "category": {"Action"},
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByGithub(describers.GetAllWorkflowRuns),
		GetDescriber:         nil,
	},

	"Github/Branch": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Github/Branch",
		Tags:                 map[string][]string{
            "category": {"Branch"},
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByGithub(describers.GetAllBranches),
		GetDescriber:         nil,
	},

	"Github/Branch/Protection": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Github/Branch/Protection",
		Tags:                 map[string][]string{
            "category": {"Branch"},
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByGithub(describers.GetAllBranchProtections),
		GetDescriber:         nil,
	},

	"Github/Commit": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Github/Commit",
		Tags:                 map[string][]string{
            "category": {"Commit"},
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByGithub(describers.ListCommits),
		GetDescriber:         nil,
	},

	"Github/Issue": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Github/Issue",
		Tags:                 map[string][]string{
            "category": {"Issue"},
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByGithub(describers.GetIssueList),
		GetDescriber:         provider.DescribeSingleByRepo(describers.GetIssue),
	},

	"Github/License": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Github/License",
		Tags:                 map[string][]string{
            "category": {"License"},
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByGithub(describers.GetLicenseList),
		GetDescriber:         provider.DescribeSingleByRepo(describers.GetLicense),
	},

	"Github/Organization": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Github/Organization",
		Tags:                 map[string][]string{
            "category": {"Organization"},
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByGithub(describers.GetOrganizationList),
		GetDescriber:         nil,
	},

	"Github/Organization/Collaborator": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Github/Organization/Collaborator",
		Tags:                 map[string][]string{
            "category": {"Organization"},
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByGithub(describers.GetAllOrganizationsCollaborators),
		GetDescriber:         nil,
	},

	"Github/Organization/Dependabot/Alert": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Github/Organization/Dependabot/Alert",
		Tags:                 map[string][]string{
            "category": {"Organization"},
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByGithub(describers.GetAllOrganizationsDependabotAlerts),
		GetDescriber:         nil,
	},

	"Github/Organization/Member": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Github/Organization/Member",
		Tags:                 map[string][]string{
            "category": {"Organization"},
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByGithub(describers.GetAllMembers),
		GetDescriber:         nil,
	},

	"Github/Organization/Team": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Github/Organization/Team",
		Tags:                 map[string][]string{
            "category": {"Organization"},
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByGithub(describers.GetOrganizationTeamList),
		GetDescriber:         nil,
	},

	"Github/PullRequest": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Github/PullRequest",
		Tags:                 map[string][]string{
            "category": {"PullRequest"},
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByGithub(describers.GetAllPullRequests),
		GetDescriber:         nil,
	},

	"Github/Release": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Github/Release",
		Tags:                 map[string][]string{
            "category": {"Release"},
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByGithub(describers.GetReleaseList),
		GetDescriber:         nil,
	},

	"Github/Repository": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Github/Repository",
		Tags:                 map[string][]string{
            "category": {"Repository"},
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByGithub(describers.GetRepositoryList),
		GetDescriber:         provider.DescribeSingleByRepo(describers.GetRepository),
	},

	"Github/Repository/Collaborator": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Github/Repository/Collaborator",
		Tags:                 map[string][]string{
            "category": {"Repository"},
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByGithub(describers.GetAllRepositoriesCollaborators),
		GetDescriber:         nil,
	},

	"Github/Repository/DependabotAlert": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Github/Repository/DependabotAlert",
		Tags:                 map[string][]string{
            "category": {"Repository"},
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByGithub(describers.GetAllRepositoriesDependabotAlerts),
		GetDescriber:         nil,
	},

	"Github/Repository/Deployment": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Github/Repository/Deployment",
		Tags:                 map[string][]string{
            "category": {"Repository"},
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByGithub(describers.GetAllRepositoriesDeployments),
		GetDescriber:         nil,
	},

	"Github/Repository/Environment": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Github/Repository/Environment",
		Tags:                 map[string][]string{
            "category": {"Repository"},
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByGithub(describers.GetAllRepositoriesEnvironments),
		GetDescriber:         nil,
	},

	"Github/Repository/Ruleset": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Github/Repository/Ruleset",
		Tags:                 map[string][]string{
            "category": {"Repository"},
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByGithub(describers.GetAllRepositoriesRuleSets),
		GetDescriber:         nil,
	},

	"Github/Repository/SBOM": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Github/Repository/SBOM",
		Tags:                 map[string][]string{
            "category": {"Repository"},
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByGithub(describers.GetAllRepositoriesSBOMs),
		GetDescriber:         nil,
	},

	"Github/Repository/VulnerabilityAlert": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Github/Repository/VulnerabilityAlert",
		Tags:                 map[string][]string{
            "category": {"Repository"},
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByGithub(describers.GetAllRepositoriesVulnerabilities),
		GetDescriber:         nil,
	},

	"Github/Tag": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Github/Tag",
		Tags:                 map[string][]string{
            "category": {"Tag"},
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByGithub(describers.GetAllTags),
		GetDescriber:         nil,
	},

	"Github/Team/Member": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Github/Team/Member",
		Tags:                 map[string][]string{
            "category": {"Team"},
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByGithub(describers.GetAllTeamsMembers),
		GetDescriber:         nil,
	},

	"Github/User": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Github/User",
		Tags:                 map[string][]string{
            "category": {"user"},
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByGithub(describers.GetUser),
		GetDescriber:         nil,
	},

	"Github/Workflow": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Github/Workflow",
		Tags:                 map[string][]string{
            "category": {"workflow"},
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByGithub(describers.GetAllWorkflows),
		GetDescriber:         provider.DescribeSingleByRepo(describers.GetRepositoryWorkflow),
	},

	"Github/Container/Package": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Github/Container/Package",
		Tags:                 map[string][]string{
            "category": {"package"},
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByGithub(describers.GetContainerPackageList),
		GetDescriber:         nil,
	},

	"Github/Package/Maven": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Github/Package/Maven",
		Tags:                 map[string][]string{
            "category": {"package"},
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByGithub(describers.GetMavenPackageList),
		GetDescriber:         nil,
	},

	"Github/NPM/Package": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Github/NPM/Package",
		Tags:                 map[string][]string{
            "category": {"package"},
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByGithub(describers.GetNPMPackageList),
		GetDescriber:         nil,
	},

	"Github/Nuget/Package": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Github/Nuget/Package",
		Tags:                 map[string][]string{
            "category": {"package"},
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByGithub(describers.GetNugetPackageList),
		GetDescriber:         provider.DescribeSingleByRepo(describers.GetNugetPackage),
	},

	"Github/Artifact/DockerFile": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Github/Artifact/DockerFile",
		Tags:                 map[string][]string{
            "category": {"artifact_dockerfile"},
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByGithub(describers.ListArtifactDockerFiles),
		GetDescriber:         nil,
	},
}


var ResourceTypeConfigs = map[string]*interfaces.ResourceTypeConfiguration{

	"Github/Actions/Artifact": {
		Name:         "Github/Actions/Artifact",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Github/Actions/Runner": {
		Name:         "Github/Actions/Runner",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Github/Actions/Secret": {
		Name:         "Github/Actions/Secret",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Github/Actions/WorkflowRun": {
		Name:         "Github/Actions/WorkflowRun",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		Params:           	[]interfaces.Param{
			{
				Name:  "organization",
				Description: "Please provide the organization name",
				Required:    false,
				Default:     nil,
			},
			
			{
				Name:  "repository",
				Description: "Please provide the repo name (i.e. internal-tools)",
				Required:    false,
				Default:     nil,
			},
			
			{
				Name:  "run_number",
				Description: "Please provide the run number",
				Required:    false,
				Default:     nil,
			},
			      },
		
	},

	"Github/Branch": {
		Name:         "Github/Branch",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Github/Branch/Protection": {
		Name:         "Github/Branch/Protection",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Github/Commit": {
		Name:         "Github/Commit",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Github/Issue": {
		Name:         "Github/Issue",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Github/License": {
		Name:         "Github/License",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Github/Organization": {
		Name:         "Github/Organization",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Github/Organization/Collaborator": {
		Name:         "Github/Organization/Collaborator",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Github/Organization/Dependabot/Alert": {
		Name:         "Github/Organization/Dependabot/Alert",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Github/Organization/Member": {
		Name:         "Github/Organization/Member",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Github/Organization/Team": {
		Name:         "Github/Organization/Team",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Github/PullRequest": {
		Name:         "Github/PullRequest",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		Params:           	[]interfaces.Param{
			{
				Name:  "organization",
				Description: "Please provide the organization name",
				Required:    false,
				Default:     nil,
			},
			
			{
				Name:  "repository",
				Description: "Please provide the repo name (i.e. internal-tools)",
				Required:    false,
				Default:     nil,
			},
			      },
		
	},

	"Github/Release": {
		Name:         "Github/Release",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Github/Repository": {
		Name:         "Github/Repository",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		Params:           	[]interfaces.Param{
			{
				Name:  "organization",
				Description: "Please provide the organization name",
				Required:    false,
				Default:     nil,
			},
			
			{
				Name:  "repository",
				Description: "Please provide the repo name (i.e. internal-tools)",
				Required:    false,
				Default:     nil,
			},
			      },
		
	},

	"Github/Repository/Collaborator": {
		Name:         "Github/Repository/Collaborator",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Github/Repository/DependabotAlert": {
		Name:         "Github/Repository/DependabotAlert",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Github/Repository/Deployment": {
		Name:         "Github/Repository/Deployment",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Github/Repository/Environment": {
		Name:         "Github/Repository/Environment",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Github/Repository/Ruleset": {
		Name:         "Github/Repository/Ruleset",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Github/Repository/SBOM": {
		Name:         "Github/Repository/SBOM",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Github/Repository/VulnerabilityAlert": {
		Name:         "Github/Repository/VulnerabilityAlert",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Github/Tag": {
		Name:         "Github/Tag",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Github/Team/Member": {
		Name:         "Github/Team/Member",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Github/User": {
		Name:         "Github/User",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Github/Workflow": {
		Name:         "Github/Workflow",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Github/Container/Package": {
		Name:         "Github/Container/Package",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		Params:           	[]interfaces.Param{
			{
				Name:  "organization",
				Description: "Please provide the organization name",
				Required:    false,
				Default:     nil,
			},
			      },
		
	},

	"Github/Package/Maven": {
		Name:         "Github/Package/Maven",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Github/NPM/Package": {
		Name:         "Github/NPM/Package",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Github/Nuget/Package": {
		Name:         "Github/Nuget/Package",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Github/Artifact/DockerFile": {
		Name:         "Github/Artifact/DockerFile",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		Params:           	[]interfaces.Param{
			{
				Name:  "organization",
				Description: "Please provide the organization name",
				Required:    false,
				Default:     nil,
			},
			
			{
				Name:  "repository",
				Description: "Please provide the repo name (i.e. internal-tools)",
				Required:    false,
				Default:     nil,
			},
			      },
		
	},
}


var ResourceTypesList = []string{
  "Github/Actions/Artifact",
  "Github/Actions/Runner",
  "Github/Actions/Secret",
  "Github/Actions/WorkflowRun",
  "Github/Branch",
  "Github/Branch/Protection",
  "Github/Commit",
  "Github/Issue",
  "Github/License",
  "Github/Organization",
  "Github/Organization/Collaborator",
  "Github/Organization/Dependabot/Alert",
  "Github/Organization/Member",
  "Github/Organization/Team",
  "Github/PullRequest",
  "Github/Release",
  "Github/Repository",
  "Github/Repository/Collaborator",
  "Github/Repository/DependabotAlert",
  "Github/Repository/Deployment",
  "Github/Repository/Environment",
  "Github/Repository/Ruleset",
  "Github/Repository/SBOM",
  "Github/Repository/VulnerabilityAlert",
  "Github/Tag",
  "Github/Team/Member",
  "Github/User",
  "Github/Workflow",
  "Github/Container/Package",
  "Github/Package/Maven",
  "Github/NPM/Package",
  "Github/Nuget/Package",
  "Github/Artifact/DockerFile",
}