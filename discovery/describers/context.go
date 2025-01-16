package describers

import (
	"context"
	"github.com/opengovern/og-util/pkg/describe/enums"
	"go.uber.org/zap"
)

var (
	triggerTypeKey string = "trigger_type"
)

func WithTriggerType(ctx context.Context, tt enums.DescribeTriggerType) context.Context {
	return context.WithValue(ctx, triggerTypeKey, tt)
}

func GetTriggerTypeFromContext(ctx context.Context) enums.DescribeTriggerType {
	tt, ok := ctx.Value(triggerTypeKey).(enums.DescribeTriggerType)
	if !ok {
		return ""
	}
	return tt
}

func GetParameterFromContext(ctx context.Context, key string) any {
	return ctx.Value(key)
}

func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, "logger", logger)
}

func GetLoggerFromContext(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value("logger").(*zap.Logger)
	if !ok {
		return zap.NewNop()
	}
	return logger
}
