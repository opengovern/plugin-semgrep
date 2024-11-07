package provider

import (
	model "github.com/opengovern/og-describer-template/pkg/sdk/models"
	"github.com/opengovern/og-describer-template/provider/configs"
	"github.com/opengovern/og-util/pkg/describe/enums"
	"golang.org/x/net/context"
)

// DescribeByIntegration TODO: implement a wrapper to pass integration authorization to describer functions
func DescribeByIntegration(describe func(context.Context, *configs.IntegrationConfig, string, *model.StreamSender) ([]model.Resource, error)) model.ResourceDescriber {
	return func(ctx context.Context, cfg configs.IntegrationConfig, triggerType enums.DescribeTriggerType, additionalParameters map[string]string, stream *model.StreamSender) ([]model.Resource, error) {

		var values []model.Resource

		return values, nil
	}
}
