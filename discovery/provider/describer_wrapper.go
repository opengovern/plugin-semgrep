package provider

import (
	model "github.com/opengovern/og-describer-semgrep/discovery/pkg/models"
	"github.com/opengovern/og-util/pkg/describe/enums"
	"golang.org/x/net/context"
)

type Client struct {
}

// DescribeByIntegration TODO: implement a wrapper to pass integration authorization to describer functions
func DescribeByIntegration(describe func(context.Context, Client, string, *model.StreamSender) ([]model.Resource, error)) model.ResourceDescriber {
	return func(ctx context.Context, cfg model.IntegrationCredentials, triggerType enums.DescribeTriggerType, additionalParameters map[string]string, stream *model.StreamSender) ([]model.Resource, error) {
		var values []model.Resource
		return values, nil
	}
}

// DescribeByIntegration TODO: implement a wrapper to pass integration authorization to describer functions
func DescribeSingleByRepo(describe func(context.Context, Client, string, string, string, *model.StreamSender) (*model.Resource, error)) model.SingleResourceDescriber {
	return func(ctx context.Context, cfg model.IntegrationCredentials, triggerType enums.DescribeTriggerType, additionalParameters map[string]string, resourceID string, stream *model.StreamSender) (*model.Resource, error) {
		var result *model.Resource
		return result, nil
	}
}
