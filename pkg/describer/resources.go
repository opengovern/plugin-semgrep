//go:generate go run ../SDK/runable/resourceType/resource_types_generator.go  --output resource_types.go --index-map ../steampipe/table_index_map.go && gofmt -w -s resource_types.go  && goimports -w resource_types.go

package describer

import (
	"context"
	"fmt"
	model "github.com/opengovern/og-describer-template/pkg/SDK/models"
	"github.com/opengovern/og-describer-template/provider"
	"github.com/opengovern/og-describer-template/provider/describer"
	"github.com/opengovern/og-util/pkg/describe/enums"
	"go.uber.org/zap"
	"sort"
	"strings"
)

func ListResourceTypes() []string {
	var list []string
	for k := range provider.ResourceTypes {
		list = append(list, k)
	}

	sort.Strings(list)
	return list
}

func ListSummarizeResourceTypes() []string {
	var list []string
	for k, v := range provider.ResourceTypes {
		if v.Summarize {
			list = append(list, k)
		}
	}

	sort.Strings(list)
	return list
}

func GetResourceType(resourceType string) (*model.ResourceType, error) {
	if r, ok := provider.ResourceTypes[resourceType]; ok {
		return &r, nil
	}
	resourceType = strings.ToLower(resourceType)
	for k, v := range provider.ResourceTypes {
		k := strings.ToLower(k)
		v := v
		if k == resourceType {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("resource type %s not found", resourceType)
}

func GetResourceTypesMap() map[string]model.ResourceType {
	return provider.ResourceTypes
}

type Resources struct {
	Resources map[string][]describer.Resource
	Errors    map[string]string
	ErrorCode string
}

func describe(ctx context.Context, logger *zap.Logger, cfg any, account string, regions []string, resourceType string, triggerType enums.DescribeTriggerType, stream *describer.StreamSender) (*Resources, error) {
	resourceTypeObject, ok := provider.ResourceTypes[resourceType]
	if !ok {
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}
	ctx = describer.WithLogger(ctx, logger)

	return resourceTypeObject.ListDescriber(ctx, cfg, account, regions, resourceType, triggerType, stream)
}

func GetResourceTypeByTerraform(terraformType string) string {
	for t, v := range provider.ResourceTypes {
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
