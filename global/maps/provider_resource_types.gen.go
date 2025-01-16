package maps
import (
	"github.com/opengovern/og-describer-template/discovery/describers"
	"github.com/opengovern/og-describer-template/discovery/provider"
	"github.com/opengovern/og-describer-template/platform/constants"
	"github.com/opengovern/og-util/pkg/integration/interfaces"
	model "github.com/opengovern/og-describer-template/discovery/pkg/models"
)
var ResourceTypes = map[string]model.ResourceType{
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
		ListDescriber:        provider.DescribeByIntegration(describers.GetContainerPackageList),
		GetDescriber:         nil,
	},
}


var ResourceTypeConfigs = map[string]*interfaces.ResourceTypeConfiguration{

	
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


}


var ResourceTypesList = []string{
  "Github/Actions/Artifact",

}