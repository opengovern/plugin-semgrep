package maps
import (
	"github.com/opengovern/og-describer-template/discovery/describers"
	"github.com/opengovern/og-describer-template/discovery/provider"
	"github.com/opengovern/og-describer-template/platform/constants"
	"github.com/opengovern/og-util/pkg/integration/interfaces"
	model "github.com/opengovern/og-describer-template/discovery/pkg/models"
)
var ResourceTypes = map[string]model.ResourceType{

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
		ListDescriber:        provider.DescribeByIntegration(describers.ListType),
		GetDescriber:         nil,
	},
}


var ResourceTypeConfigs = map[string]*interfaces.ResourceTypeConfiguration{

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
  "Github/Artifact/DockerFile",
}