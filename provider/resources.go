//go:generate go run ../SDK/runable/resourceType/resource_types_generator.go  --output resource_types.go --index-map ../steampipe/table_index_map.go && gofmt -w -s resource_types.go  && goimports -w resource_types.go

package provider

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"go.uber.org/zap"
	"github.com/opengovern/og-util/pkg/describe/enums"
	"github.com/opengovern/og-util/pkg/source"
)
// any types are used to load your provider configuration.
type ResourceDescriber func(context.Context, any, string, []string, string, enums.DescribeTriggerType, *describer.StreamSender) (*Resources, error)
type SingleResourceDescriber func(context.Context, any, string, []string, string, map[string]string, enums.DescribeTriggerType) (*Resources, error)

type ResourceType struct {
	Connector source.Type

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

func (r ResourceType) GetConnector() source.Type {
	return r.Connector
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

func ListResourceTypes() []string {
	var list []string
	for k := range resourceTypes {
		list = append(list, k)
	}

	sort.Strings(list)
	return list
}

func ListFastDiscoveryResourceTypes() []string {
	var list []string
	for k, v := range resourceTypes {
		if v.FastDiscovery {
			list = append(list, k)
		}
	}

	sort.Strings(list)
	return list
}

func ListSummarizeResourceTypes() []string {
	var list []string
	for k, v := range resourceTypes {
		if v.Summarize {
			list = append(list, k)
		}
	}

	sort.Strings(list)
	return list
}

func GetResourceType(resourceType string) (*ResourceType, error) {
	if r, ok := resourceTypes[resourceType]; ok {
		return &r, nil
	}
	resourceType = strings.ToLower(resourceType)
	for k, v := range resourceTypes {
		k := strings.ToLower(k)
		v := v
		if k == resourceType {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("resource type %s not found", resourceType)
}

func GetResourceTypesMap() map[string]ResourceType {
	return resourceTypes
}

type Resources struct {
	Resources map[string][]describer.Resource
	Errors    map[string]string
	ErrorCode string
}


func describe(ctx context.Context, logger *zap.Logger, cfg any, account string, regions []string, resourceType string, triggerType enums.DescribeTriggerType, stream *describer.StreamSender) (*Resources, error) {
	resourceTypeObject, ok := resourceTypes[resourceType]
	if !ok {
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}
	ctx = describer.WithLogger(ctx, logger)

	return resourceTypeObject.ListDescriber(ctx, cfg, account, regions, resourceType, triggerType, stream)
}

func GetResourceTypeByTerraform(terraformType string) string {
	for t, v := range resourceTypes {
		for _, name := range v.TerraformName {
			if name == terraformType {
				return t
			}
		}
	}
	return ""
}


// ResourceTypes is a map of all the resource types supported by the provider.
// TODO: Add your resource types here.
// When you add a new resource type, you should also add a new entry in the resourceTypes map.
// Write a function to describe the resource type and add it to the resourceTypes map.
