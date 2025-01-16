package main

import (
	"encoding/json"
	"github.com/opengovern/og-describer-template/platform/constants"
	"github.com/opengovern/og-describer-template/global"
	"github.com/opengovern/og-describer-template/global/maps"
	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/og-util/pkg/integration/interfaces"
)

type Integration struct{}

func (i *Integration) GetConfiguration() interfaces.IntegrationConfiguration {
	return interfaces.IntegrationConfiguration{
		NatsScheduledJobsTopic:   global.JobQueueTopic,
		NatsManualJobsTopic:      global.JobQueueTopicManuals,
		NatsStreamName:           global.StreamName,
		NatsConsumerGroup:        global.ConsumerGroup,
		NatsConsumerGroupManuals: global.ConsumerGroupManuals,

		SteampipePluginName: "template",

		UISpec:   constants.UISpec,
		Manifest: constants.Manifest,
		SetupMD:  constants.SetupMd,

		DescriberDeploymentName: constants.DescriberDeploymentName,
		DescriberRunCommand:     constants.DescriberRunCommand,
	}
}

func (i *Integration) HealthCheck(jsonData []byte, providerId string, labels map[string]string, annotations map[string]string) (bool, error) {
	var credentials global.IntegrationCredentials
	err := json.Unmarshal(jsonData, &credentials)
	if err != nil {
		return false, err
	}
	// TODO add credentials
	isHealthy, err := IntegrationHealthcheck(Config{
	})
	
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
	_, err = IntegrationDiscovery(Config{
	})
	
	return integrations, nil
}

func (i *Integration) GetResourceTypesByLabels(labels map[string]string) (map[string]interfaces.ResourceTypeConfiguration, error) {
	resourceTypesMap := make(map[string]interfaces.ResourceTypeConfiguration)
	for _, resourceType := range maps.ResourceTypesList {
		if v, ok := maps.ResourceTypeConfigs[resourceType]; ok {
			resourceTypesMap[resourceType] = *v
		} else {
			resourceTypesMap[resourceType] = interfaces.ResourceTypeConfiguration{}
		}
	}
	return resourceTypesMap, nil
}

func (i *Integration) GetResourceTypeFromTableName(tableName string) string {
	if v, ok := maps.TablesToResourceTypes[tableName]; ok {
		return v
	}

	return ""
}

func (i *Integration) GetIntegrationType() integration.Type {
	return constants.IntegrationName
}

func (i *Integration) ListAllTables() map[string][]interfaces.CloudQLColumn {
	plugin := global.Plugin()
	tables := make(map[string][]interfaces.CloudQLColumn)
	for tableKey, table := range plugin.TableMap {
		columns := make([]interfaces.CloudQLColumn, 0, len(table.Columns))
		for _, column := range table.Columns {
			columns = append(columns, interfaces.CloudQLColumn{Name: column.Name, Type: column.Type.String()})
		}
		tables[tableKey] = columns
	}

	return tables
}
