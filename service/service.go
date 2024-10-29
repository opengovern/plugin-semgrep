package local

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"runtime"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/opengovern/og-describer-template/describer"
	"github.com/opengovern/og-util/pkg/config"
	"github.com/opengovern/og-util/pkg/describe"
	esSinkClient "github.com/opengovern/og-util/pkg/es/ingest/client"
	"github.com/opengovern/og-util/pkg/jq"
	"github.com/opengovern/og-util/pkg/koanf"
	"github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	StreamName           = "og_aws_describer"
	JobQueueTopic        = "og_aws_describer_job_queue"
	ConsumerGroup        = "aws-describer"
	JobQueueTopicManuals = "og_aws_describer_manuals_job_queue"
	ConsumerGroupManuals = "aws-describer-manuals"
)

var (
	ManualTriggers = os.Getenv("MANUAL_TRIGGERS")
)

type Config struct {
	NATS config.NATS `koanf:"nats"`
}

func WorkerCommand() *cobra.Command {
	cmd := &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			cmd.SilenceUsage = true
			cnf := koanf.Provide("aws_describer", Config{})
			logger, err := zap.NewProduction()
			if err != nil {
				return err
			}

			w, err := NewWorker(
				cnf,
				logger,
				cmd.Context(),
			)
			if err != nil {
				return err
			}

			return w.Run(ctx)
		},
	}

	return cmd
}

type Worker struct {
	config   Config
	logger   *zap.Logger
	esClient opengovernance.Client
	jq       *jq.JobQueue

	esSinkClient esSinkClient.EsSinkServiceClient
}

func NewWorker(
	config Config,
	logger *zap.Logger,
	ctx context.Context,
) (*Worker, error) {
	jq, err := jq.New(config.NATS.URL, logger)
	if err != nil {
		logger.Error("failed to create job queue", zap.Error(err), zap.String("url", config.NATS.URL))
		return nil, err
	}

	topic := JobQueueTopic
	if ManualTriggers == "true" {
		topic = JobQueueTopicManuals
	}
	if err := jq.Stream(ctx, StreamName, "aws describe job runner queue", []string{topic}, 200000); err != nil {
		logger.Error("failed to create stream", zap.Error(err))
		return nil, err
	}

	w := &Worker{
		config: config,
		logger: logger,
		jq:     jq,
	}

	return w, nil
}

func (w *Worker) Run(ctx context.Context) error {
	w.logger.Info("starting to consume")
	topic := JobQueueTopic
	consumer := ConsumerGroup
	if ManualTriggers == "true" {
		topic = JobQueueTopicManuals
		consumer = ConsumerGroupManuals
	}
	consumeCtx, err := w.jq.ConsumeWithConfig(ctx, consumer, StreamName, []string{topic}, jetstream.ConsumerConfig{
		Replicas:          1,
		AckPolicy:         jetstream.AckExplicitPolicy,
		DeliverPolicy:     jetstream.DeliverAllPolicy,
		MaxAckPending:     -1,
		AckWait:           time.Minute * 30,
		InactiveThreshold: time.Hour,
	}, []jetstream.PullConsumeOpt{
		jetstream.PullMaxMessages(1),
	}, func(msg jetstream.Msg) {
		w.logger.Info("received a new job")

		defer msg.Ack()

		ctx, cancel := context.WithTimeoutCause(ctx, time.Minute*25, errors.New("describe worker timed out"))
		defer cancel()

		if err := w.ProcessMessage(ctx, msg); err != nil {
			w.logger.Error("failed to process message", zap.Error(err))
		}
		err := msg.Ack()
		if err != nil {
			w.logger.Error("failed to ack message", zap.Error(err))
		}

		w.logger.Info("processing a job completed")
	})
	if err != nil {
		return err
	}

	w.logger.Info("consuming")

	<-ctx.Done()
	consumeCtx.Drain()
	consumeCtx.Stop()

	return nil
}

func (w *Worker) ProcessMessage(ctx context.Context, msg jetstream.Msg) error {
	startTime := time.Now()
	var input describe.DescribeWorkerInput
	err := json.Unmarshal(msg.Data(), &input)
	if err != nil {
		return err
	}
	runtime.GC()

	w.logger.Info("running job", zap.Uint("id", input.DescribeJob.JobID), zap.String("type", input.DescribeJob.ResourceType), zap.String("account", input.DescribeJob.AccountID))

	err = describer.DescribeHandler(ctx, w.logger, describer.TriggeredByLocal, input)
	endTime := time.Now()

	w.logger.Info("job completed", zap.Uint("id", input.DescribeJob.JobID), zap.String("type", input.DescribeJob.ResourceType), zap.String("account", input.DescribeJob.AccountID), zap.Duration("duration", endTime.Sub(startTime)))
	if err != nil {
		w.logger.Error("failure while running job", zap.Error(err))
		return err
	}

	return nil
}
