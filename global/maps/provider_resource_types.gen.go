package maps

import (
	"github.com/opengovern/og-describer-semgrep/discovery/describers"
	model "github.com/opengovern/og-describer-semgrep/discovery/pkg/models"
	"github.com/opengovern/og-describer-semgrep/discovery/provider"
	"github.com/opengovern/og-describer-semgrep/platform/constants"
	"github.com/opengovern/og-util/pkg/integration/interfaces"
)

var ResourceTypes = map[string]model.ResourceType{

	"Semgrep/Deployment": {
		IntegrationType: constants.IntegrationName,
		ResourceName:    "Semgrep/Deployment",
		Tags:            map[string][]string{},
		Labels:          map[string]string{},
		Annotations:     map[string]string{},
		ListDescriber:   provider.DescribeListBySemGrep(describers.ListDeployments),
		GetDescriber:    nil,
	},

	"Semgrep/Project": {
		IntegrationType: constants.IntegrationName,
		ResourceName:    "Semgrep/Project",
		Tags:            map[string][]string{},
		Labels:          map[string]string{},
		Annotations:     map[string]string{},
		ListDescriber:   provider.DescribeListBySemGrep(describers.ListProjects),
		GetDescriber:    nil,
	},

	"Semgrep/Policy": {
		IntegrationType: constants.IntegrationName,
		ResourceName:    "Semgrep/Policy",
		Tags:            map[string][]string{},
		Labels:          map[string]string{},
		Annotations:     map[string]string{},
		ListDescriber:   provider.DescribeListBySemGrep(describers.ListPolicies),
		GetDescriber:    nil,
	},

	"Semgrep/Scan": {
		IntegrationType: constants.IntegrationName,
		ResourceName:    "Semgrep/Scan",
		Tags:            map[string][]string{},
		Labels:          map[string]string{},
		Annotations:     map[string]string{},
		ListDescriber:   provider.DescribeListBySemGrep(describers.ListScans),
		GetDescriber:    nil,
	},

	"Semgrep/Finding": {
		IntegrationType: constants.IntegrationName,
		ResourceName:    "Semgrep/Finding",
		Tags:            map[string][]string{},
		Labels:          map[string]string{},
		Annotations:     map[string]string{},
		ListDescriber:   provider.DescribeListBySemGrep(describers.ListFindings),
		GetDescriber:    nil,
	},
}

var ResourceTypeConfigs = map[string]*interfaces.ResourceTypeConfiguration{

	"Semgrep/Deployment": {
		Name:            "Semgrep/Deployment",
		IntegrationType: constants.IntegrationName,
		Description:     "",
	},

	"Semgrep/Project": {
		Name:            "Semgrep/Project",
		IntegrationType: constants.IntegrationName,
		Description:     "",
	},

	"Semgrep/Policy": {
		Name:            "Semgrep/Policy",
		IntegrationType: constants.IntegrationName,
		Description:     "",
	},

	"Semgrep/Scan": {
		Name:            "Semgrep/Scan",
		IntegrationType: constants.IntegrationName,
		Description:     "",
	},

	"Semgrep/Finding": {
		Name:            "Semgrep/Finding",
		IntegrationType: constants.IntegrationName,
		Description:     "",
	},
}

var ResourceTypesList = []string{
	"Semgrep/Deployment",
	"Semgrep/Project",
	"Semgrep/Policy",
	"Semgrep/Scan",
	"Semgrep/Finding",
}
