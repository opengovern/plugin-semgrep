package maps

import (
	model "github.com/opengovern/og-describer-semgrep/discovery/pkg/models"
	"github.com/opengovern/og-describer-semgrep/platform/constants"
	"github.com/opengovern/og-util/pkg/integration/interfaces"
)

var ResourceTypes = map[string]model.ResourceType{

	"Github/Artifact/DockerFile": {
		IntegrationType: constants.IntegrationName,
		ResourceName:    "Github/Artifact/DockerFile",
		Tags: map[string][]string{
			"category": {"artifact_dockerfile"},
		},
		Labels:        map[string]string{},
		Annotations:   map[string]string{},
		ListDescriber: nil,
		GetDescriber:  nil,
	},
}

var ResourceTypeConfigs = map[string]*interfaces.ResourceTypeConfiguration{

	"Github/Artifact/DockerFile": {
		Name:            "Github/Artifact/DockerFile",
		IntegrationType: constants.IntegrationName,
		Description:     "",
		Params: []interfaces.Param{
			{
				Name:        "organization",
				Description: "Please provide the organization name",
				Required:    false,
				Default:     nil,
			},

			{
				Name:        "repository",
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
