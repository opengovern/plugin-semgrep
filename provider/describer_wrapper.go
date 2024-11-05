package provider

import (
	model "github.com/opengovern/og-describer-template/pkg/SDK/models"
	"github.com/opengovern/og-describer-template/provider/configs"
	"github.com/opengovern/og-util/pkg/describe/enums"
	"golang.org/x/net/context"
)

func DescribeByAccount(describe func(context.Context, *configs.AccountConfig, string, *model.StreamSender) ([]model.Resource, error)) model.ResourceDescriber {
	return func(ctx context.Context, cfg configs.AccountConfig, triggerType enums.DescribeTriggerType, additionalParameters map[string]string, stream *model.StreamSender) ([]model.Resource, error) {

		var values []model.Resource

		return values, nil
	}
}
