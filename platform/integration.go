package main

import (
	"encoding/json"
	"github.com/opengovern/og-describer-semgrep/global"
	"github.com/opengovern/og-describer-semgrep/global/maps"
	"github.com/opengovern/og-describer-semgrep/platform/constants"
	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/og-util/pkg/integration/interfaces"
)

type Integration struct{}

func (i *Integration) GetConfiguration() (interfaces.IntegrationConfiguration, error) {
	return interfaces.IntegrationConfiguration{
		NatsScheduledJobsTopic:   global.JobQueueTopic,
		NatsManualJobsTopic:      global.JobQueueTopicManuals,
		NatsStreamName:           global.StreamName,
		NatsConsumerGroup:        global.ConsumerGroup,
		NatsConsumerGroupManuals: global.ConsumerGroupManuals,

		SteampipePluginName: "semgrep",

		UISpec:   constants.UISpec,
		Manifest: constants.Manifest,
		SetupMD:  constants.SetupMd,

		DescriberDeploymentName: constants.DescriberDeploymentName,
		DescriberRunCommand:     constants.DescriberRunCommand,
	}, nil
}

func (i *Integration) HealthCheck(jsonData []byte, providerId string, labels map[string]string, annotations map[string]string) (bool, error) {
	var credentials global.IntegrationCredentials
	err := json.Unmarshal(jsonData, &credentials)
	if err != nil {
		return false, err
	}
	// TODO add credentials
	isHealthy, err := IntegrationHealthcheck(Config{})

	return isHealthy, err
}

func (i *Integration) DiscoverIntegrations(jsonData []byte) ([]integration.Integration, error) {
	var credentials global.IntegrationCredentials
	err := json.Unmarshal(jsonData, &credentials)
	if err != nil {
		return nil, err
	}
	var integrations []integration.Integration
	// TODO
	_, err = IntegrationDiscovery(Config{})

	return integrations, nil
}

func (i *Integration) GetResourceTypesByLabels(labels map[string]string) ([]interfaces.ResourceTypeConfiguration, error) {
	var resourceTypesMap []interfaces.ResourceTypeConfiguration
	for _, resourceType := range maps.ResourceTypesList {
		var resource interfaces.ResourceTypeConfiguration
		if v, ok := maps.ResourceTypeConfigs[resourceType]; ok {
			resource.Description = v.Description
			resource.Params = v.Params
			resource.Name = v.Name
			resource.IntegrationType = v.IntegrationType
			resource.Table = maps.ResourceTypesToTables[v.Name]
			resourceTypesMap = append(resourceTypesMap, resource)

		}
	}
	return resourceTypesMap, nil
}
func (i *Integration) GetResourceTypeFromTableName(tableName string) (string, error) {
	if v, ok := maps.TablesToResourceTypes[tableName]; ok {
		return v, nil
	}

	return "", nil
}

func (i *Integration) GetIntegrationType() (integration.Type, error) {
	return constants.IntegrationName, nil
}

func (i *Integration) ListAllTables() (map[string][]interfaces.CloudQLColumn, error) {
	plugin := global.Plugin()
	tables := make(map[string][]interfaces.CloudQLColumn)
	for tableKey, table := range plugin.TableMap {
		columns := make([]interfaces.CloudQLColumn, 0, len(table.Columns))
		for _, column := range table.Columns {
			columns = append(columns, interfaces.CloudQLColumn{Name: column.Name, Type: column.Type.String()})
		}
		tables[tableKey] = columns
	}

	return tables, nil
}

func (i *Integration) Ping() error {
	return nil
}
