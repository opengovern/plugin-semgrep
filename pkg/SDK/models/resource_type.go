package models

import (
	"github.com/opengovern/og-describer-template/provider/configs"
	"github.com/opengovern/og-util/pkg/describe/enums"
	"github.com/opengovern/og-util/pkg/integration"
	"golang.org/x/net/context"
)

// any types are used to load your provider configuration.
type ResourceDescriber func(context.Context, configs.AccountConfig, enums.DescribeTriggerType, map[string]string, *StreamSender) ([]Resource, error)
type SingleResourceDescriber func(context.Context, configs.AccountConfig, enums.DescribeTriggerType, map[string]string) (*Resource, error)

type ResourceType struct {
	IntegrationType integration.Type

	ResourceName  string
	ResourceLabel string
	ServiceName   string

	Tags map[string][]string

	ListDescriber ResourceDescriber
	GetDescriber  SingleResourceDescriber

	TerraformName        []string
	TerraformServiceName string

	FastDiscovery bool
	CostDiscovery bool
	Summarize     bool
}

func (r ResourceType) GetConnector() integration.Type {
	return r.IntegrationType
}

func (r ResourceType) GetResourceName() string {
	return r.ResourceName
}

func (r ResourceType) GetResourceLabel() string {
	return r.ResourceLabel
}

func (r ResourceType) GetServiceName() string {
	return r.ServiceName
}

func (r ResourceType) GetTags() map[string][]string {
	return r.Tags
}

func (r ResourceType) GetTerraformName() []string {
	return r.TerraformName
}

func (r ResourceType) GetTerraformServiceName() string {
	return r.TerraformServiceName
}

func (r ResourceType) IsFastDiscovery() bool {
	return r.FastDiscovery
}

func (r ResourceType) IsCostDiscovery() bool {
	return r.CostDiscovery
}

func (r ResourceType) IsSummarized() bool {
	return r.Summarize
}
