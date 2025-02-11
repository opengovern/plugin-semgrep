package provider

import (
	"errors"
	"github.com/opengovern/og-describer-semgrep/discovery/pkg/models"
	"github.com/opengovern/og-util/pkg/describe/enums"
	"golang.org/x/net/context"
	"golang.org/x/time/rate"
	"time"
)

// DescribeListBySemGrep A wrapper to pass SemGrep authorization to describers functions
func DescribeListBySemGrep(describe func(context.Context, *SemGrepAPIHandler, *models.StreamSender) ([]models.Resource, error)) models.ResourceDescriber {
	return func(ctx context.Context, cfg models.IntegrationCredentials, triggerType enums.DescribeTriggerType, additionalParameters map[string]string, stream *models.StreamSender) ([]models.Resource, error) {
		ctx = WithTriggerType(ctx, triggerType)

		var err error
		// Check for the token
		if cfg.Token == "" {
			return nil, errors.New("token must be configured")
		}

		semGrepAPIHandler := NewSemGrepAPIHandler(cfg.Token, rate.Every(time.Minute/200), 1, 10, 5, 5*time.Minute)

		// Get values from describers
		var values []models.Resource
		result, err := describe(ctx, semGrepAPIHandler, stream)
		if err != nil {
			return nil, err
		}
		values = append(values, result...)
		return values, nil
	}
}

// DescribeSingleBySemGrep A wrapper to pass SemGrep authorization to describers functions
func DescribeSingleBySemGrep(describe func(context.Context, *SemGrepAPIHandler, string) (*models.Resource, error)) models.SingleResourceDescriber {
	return func(ctx context.Context, cfg models.IntegrationCredentials, triggerType enums.DescribeTriggerType, additionalParameters map[string]string, resourceID string, stream *models.StreamSender) (*models.Resource, error) {
		ctx = WithTriggerType(ctx, triggerType)

		var err error
		// Check for the token
		if cfg.Token == "" {
			return nil, errors.New("token must be configured")
		}

		semGrepAPIHandler := NewSemGrepAPIHandler(cfg.Token, rate.Every(time.Minute/200), 1, 10, 5, 5*time.Minute)

		// Get value from describers
		value, err := describe(ctx, semGrepAPIHandler, resourceID)
		if err != nil {
			return nil, err
		}
		return value, nil
	}
}
