package steampipe

import (
	"context"
	"strings"

	"go.uber.org/zap"

	"github.com/hashicorp/go-hclog"

	"fmt"

	"github.com/opengovern/og-util/pkg/steampipe"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/context_key"
)



func  buildContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, context_key.Logger, hclog.New(nil))
	return ctx
}

func  DescriptionToRecord(logger *zap.Logger, resource interface{}, indexName string) (map[string]*proto.Column, error) {
	// To Do
	// Example
	// return steampipe.DescriptionToRecord(logger, aws.Plugin(buildContext()), resource, indexName)
	return nil, nil
	
}

func  Cells(indexName string) ([]string, error) {
	// To Do
	// Example
	// return steampipe.Cells(aws.Plugin(buildContext()), indexName)
	return nil, nil
}

func ExtractTableName(resourceType string) string {
	resourceType = strings.ToLower(resourceType)
	for k, v := range Map {
		if resourceType == strings.ToLower(k) {
			return v
		}
	}
	return ""
	
}

func ExtractResourceType(tableName string) string {
	tableName = strings.ToLower(tableName)
	return strings.ToLower(ReverseMap[tableName])
}

func GetResourceTypeByTableName(tableName string) string {
	return ExtractResourceType(tableName)
}

func Plugin() *plugin.Plugin {
	// To Do
	// Example
	// return aws.Plugin(buildContext())
	return nil
}

func ExtractTagsAndNames(logger *zap.Logger, plg *plugin.Plugin, resourceType string, source interface{}) (map[string]string, string, error) {
	pluginTableName := ExtractTableName(resourceType)
	if pluginTableName == "" {
		return nil, "", fmt.Errorf("cannot find table name for resourceType: %s", resourceType)
	}
	return steampipe.ExtractTagsAndNames(plg, logger, pluginTableName, resourceType, source, DescriptionMap)
}
