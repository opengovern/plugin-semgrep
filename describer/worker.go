package describer

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-errors/errors"
	"github.com/opengovern/og-describer-template/provider"
	"github.com/opengovern/og-describer-template/provider/describer"
	model "github.com/opengovern/og-describer-template/provider/model"
	"github.com/opengovern/og-describer-template/steampipe"
	"github.com/opengovern/og-util/pkg/describe"
	"github.com/opengovern/og-util/pkg/source"
	"github.com/opengovern/og-util/pkg/vault"
	"go.uber.org/zap"
)

type Error struct {
	ErrCode string

	error
}

func Do(ctx context.Context,
	vlt vault.VaultSourceConfig,
	logger *zap.Logger,
	job describe.DescribeJob,
	grpcEndpoint string,
	describeDeliverToken string,
	ingestionPipelineEndpoint string,
	useOpenSearch bool) (resourceIDs []string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("paniced with error: %v", r)
			logger.Error("paniced with error", zap.Error(err), zap.String("stackTrace", errors.Wrap(r, 2).ErrorStack()))
		}
	}()


	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	config, err := vlt.Decrypt(ctx, job.CipherText)
	if err != nil {
		return nil, fmt.Errorf("decrypt error: %w", err)
	}
	logger.Info("decrypted config", zap.Any("config", config))

	return doDescribe(ctx, logger, job, config, grpcEndpoint, ingestionPipelineEndpoint, describeDeliverToken, useOpenSearch)
}
// TODO
func doDescribe(ctx context.Context, logger *zap.Logger, job describe.DescribeJob, config map[string]any, grpcEndpoint, ingestionPipelineEndpoint string, describeToken string, useOpenSearch bool) ([]string, error) {
	logger.Info("Making New Resource Sender")
	rs, err := NewResourceSender(grpcEndpoint, ingestionPipelineEndpoint, describeToken, job.JobID, useOpenSearch, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to resource sender: %w", err)
	}

	logger.Info("Connect to steampipe plugin")
	plg := steampipe.Plugin()
	logger.Info("Account Config From Map")
	creds, err := provider.AccountConfigFromMap(config)
	if err != nil {
		return nil, fmt.Errorf(" account credentials: %w", err)
	}

	f := func(resource describer.Resource) error {
		// Send the resource to the resource sender
	}
	clientStream := (*describer.StreamSender)(&f)

	logger.Info("Created Client Stream")
	// Get the resource type object
	output, err := "",nil
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	logger.Info("Finished getting resources", zap.Any("output", output))

	rs.Finish()

	kerr:= nil
	// Check if there are any errors

	

	

	return rs.GetResourceIDs(), kerr
}
