package describer

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/opengovern/og-util/pkg/es"
	"github.com/opengovern/og-util/pkg/source"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/anypb"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/opengovern/og-util/proto/src/golang"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/metadata"
)

const (
	MinBufferSize   int           = 10
	MaxBufferSize   int           = 100
	ChannelSize     int           = 1000
	BufferEmptyRate time.Duration = 5 * time.Second
)

type ResourceSender struct {
	authToken                 string
	logger                    *zap.Logger
	resourceChannel           chan *golang.AWSResource
	resourceIDs               []string
	doneChannel               chan interface{}
	conn                      *grpc.ClientConn
	grpcEndpoint              string
	ingestionPipelineEndpoint string
	jobID                     uint

	client     golang.EsSinkServiceClient
	httpClient *http.Client

	sendBuffer    []*golang.AWSResource
	useOpenSearch bool
}

func NewResourceSender(grpcEndpoint, ingestionPipelineEndpoint string, describeToken string, jobID uint, useOpenSearch bool, logger *zap.Logger) (*ResourceSender, error) {
	rs := ResourceSender{
		authToken:                 describeToken,
		logger:                    logger,
		resourceChannel:           make(chan *golang.AWSResource, ChannelSize),
		resourceIDs:               nil,
		doneChannel:               make(chan interface{}),
		conn:                      nil,
		grpcEndpoint:              grpcEndpoint,
		ingestionPipelineEndpoint: ingestionPipelineEndpoint,
		jobID:                     jobID,
		useOpenSearch:             useOpenSearch,

		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
	if err := rs.Connect(); err != nil {
		return nil, err
	}

	go rs.ResourceHandler()
	return &rs, nil
}

func (s *ResourceSender) Connect() error {
	var opts []grpc.DialOption
	if s.authToken != "" {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})))
		opts = append(opts, grpc.WithPerRPCCredentials(oauth.TokenSource{
			TokenSource: oauth2.StaticTokenSource(&oauth2.Token{
				AccessToken: s.authToken,
			}),
		}))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.NewClient(
		s.grpcEndpoint,
		opts...,
	)
	if err != nil {
		return err
	}
	s.conn = conn

	client := golang.NewEsSinkServiceClient(conn)
	s.client = client
	return nil
}

func (s *ResourceSender) ResourceHandler() {
	t := time.NewTicker(BufferEmptyRate)
	defer t.Stop()

	for {
		select {
		case resource := <-s.resourceChannel:
			if resource == nil {
				s.flushBuffer(true)
				s.doneChannel <- struct{}{}
				return
			}

			s.resourceIDs = append(s.resourceIDs, resource.UniqueId)
			s.sendBuffer = append(s.sendBuffer, resource)

			if len(s.sendBuffer) > MaxBufferSize {
				s.flushBuffer(true)
			}
		case <-t.C:
			s.flushBuffer(false)
		}
	}
}

func (s *ResourceSender) sendToBackend(resourcesToSend []es.Doc) {
	grpcCtx := metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{
		"resource-job-id": fmt.Sprintf("%d", s.jobID),
	}))

	docs := make([]*anypb.Any, 0, len(resourcesToSend))
	for _, resource := range resourcesToSend {
		docBytes, err := json.Marshal(resource)
		if err != nil {
			s.logger.Error("failed to marshal resource", zap.Error(err))
			continue
		}
		docs = append(docs, &anypb.Any{Value: docBytes})
	}

	_, err := s.client.Ingest(grpcCtx, &golang.IngestRequest{Docs: docs})
	if err != nil {
		s.logger.Error("failed to send resource", zap.Error(err))
		if errors.Is(err, io.EOF) {
			err = s.Connect()
			if err != nil {
				s.logger.Error("failed to reconnect", zap.Error(err))
			}
		}
		return
	}
}

func (s *ResourceSender) sendToOpenSearchIngestPipeline(resourcesToSend []es.Doc) {
	if len(resourcesToSend) == 0 {
		return
	}

	jsonResourcesToSend, err := json.Marshal(resourcesToSend)
	if err != nil {
		s.logger.Error("failed to marshal resources", zap.Error(err))
		return
	}

	req, err := http.NewRequest(
		http.MethodPost,
		s.ingestionPipelineEndpoint,
		strings.NewReader(string(jsonResourcesToSend)),
	)
	req.Header.Add("Content-Type", "application/json")

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		s.logger.Error("failed to load configuration", zap.Error(err))
		return
	}
	creds, err := cfg.Credentials.Retrieve(context.Background())
	if err != nil {
		s.logger.Error("failed to retrieve credentials", zap.Error(err))
		return
	}

	signer := v4.NewSigner()
	err = signer.SignHTTP(context.TODO(), creds, req,
		fmt.Sprintf("%x", sha256.Sum256(jsonResourcesToSend)),
		"osis", "us-east-2", time.Now())
	if err != nil {
		s.logger.Error("failed to sign request", zap.Error(err))
		return
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Error("failed to send request", zap.Error(err))
		return
	}
	defer resp.Body.Close()
	// check status
	if resp.StatusCode != http.StatusOK {
		bodyStr := ""
		if resp.Body != nil {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				s.logger.Warn("failed to read response body", zap.Error(err))
			} else {
				bodyStr = string(bodyBytes)
			}
		}
		s.logger.Error("failed to send resources to OpenSearch",
			zap.Int("statusCode", resp.StatusCode),
			zap.String("body", bodyStr),
		)
		return
	}
}

func (s *ResourceSender) flushBuffer(force bool) {
	if len(s.sendBuffer) == 0 {
		return
	}

	if !force && len(s.sendBuffer) < MinBufferSize {
		return
	}

	resourcesToSend := make([]es.Doc, 0, 2*len(s.sendBuffer))
	for _, resource := range s.sendBuffer {
		var description any
		err := json.Unmarshal([]byte(resource.DescriptionJson), &description)
		if err != nil {
			s.logger.Error("failed to parse resource description json", zap.Error(err), zap.Uint32("jobID", resource.Job.JobId), zap.String("resourceID", resource.Id))
			continue
		}

		tags := make([]es.Tag, 0, len(resource.Tags))
		for k, v := range resource.Tags {
			tags = append(tags, es.Tag{
				// tags should be case-insensitive
				Key:   strings.ToLower(k),
				Value: strings.ToLower(v),
			})
		}

		kafkaResource := es.Resource{
			ID:            resource.UniqueId,
			ARN:           resource.Arn,
			Name:          resource.Name,
			SourceType:    source.CloudAWS,
			ResourceType:  strings.ToLower(resource.Job.ResourceType),
			Location:      resource.Region,
			SourceID:      resource.Job.SourceId,
			ResourceJobID: uint(resource.Job.JobId),
			SourceJobID:   uint(resource.Job.ParentJobId),
			ScheduleJobID: uint(resource.Job.ScheduleJobId),
			CreatedAt:     resource.Job.DescribedAt,
			Description:   description,
			Metadata:      resource.Metadata,
			CanonicalTags: tags,
		}
		keys, idx := kafkaResource.KeysAndIndex()
		kafkaResource.EsID = es.HashOf(keys...)
		kafkaResource.EsIndex = idx

		lookupResource := es.LookupResource{
			ResourceID:    resource.UniqueId,
			Name:          resource.Name,
			SourceType:    source.CloudAWS,
			ResourceType:  strings.ToLower(resource.Job.ResourceType),
			Location:      resource.Region,
			SourceID:      resource.Job.SourceId,
			ResourceJobID: uint(resource.Job.JobId),
			SourceJobID:   uint(resource.Job.ParentJobId),
			ScheduleJobID: uint(resource.Job.ScheduleJobId),
			CreatedAt:     resource.Job.DescribedAt,
			Tags:          tags,
		}
		lookupKeys, lookupIdx := lookupResource.KeysAndIndex()
		lookupResource.EsID = es.HashOf(lookupKeys...)
		lookupResource.EsIndex = lookupIdx

		resourcesToSend = append(resourcesToSend, kafkaResource)
		resourcesToSend = append(resourcesToSend, lookupResource)
	}

	s.sendToBackend(resourcesToSend)
	s.sendBuffer = nil
}

func (s *ResourceSender) Finish() {
	s.resourceChannel <- nil
	_ = <-s.doneChannel
	s.conn.Close()
}

func (s *ResourceSender) GetResourceIDs() []string {
	return s.resourceIDs
}

func (s *ResourceSender) Send(resource *golang.AWSResource) {
	s.resourceChannel <- resource
}
