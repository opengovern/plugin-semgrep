package global

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-hclog"
	semgrep "github.com/opengovern/og-describer-semgrep/cloudql/semgrep"
	"github.com/opengovern/og-describer-semgrep/global/maps"
	"github.com/opengovern/og-util/pkg/steampipe"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/context_key"
	"go.uber.org/zap"
	"strings"
)

func buildContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, context_key.Logger, hclog.New(nil))
	return ctx
}

func ExtractTableName(resourceType string) string {
	resourceType = strings.ToLower(resourceType)
	for k, v := range maps.ResourceTypesToTables {
		if resourceType == strings.ToLower(k) {
			return v
		}
	}
	return ""

}

func Plugin() *plugin.Plugin {
	return semgrep.Plugin(buildContext())
}

func ExtractTagsAndNames(logger *zap.Logger, plg *plugin.Plugin, resourceType string, source interface{}) (map[string]string, string, error) {
	pluginTableName := ExtractTableName(resourceType)
	if pluginTableName == "" {
		return nil, "", fmt.Errorf("cannot find table name for resourceType: %s", resourceType)
	}
	return steampipe.ExtractTagsAndNames(plg, logger, pluginTableName, resourceType, source, maps.ResourceTypeToDescription)
}

func ExtractResourceType(tableName string) string {
	tableName = strings.ToLower(tableName)
	return strings.ToLower(maps.TablesToResourceTypes[tableName])
}

// GetResourceTypeByTableName TODO: use this in integration implementation
func GetResourceTypeByTableName(tableName string) string {
	return ExtractResourceType(tableName)
}
