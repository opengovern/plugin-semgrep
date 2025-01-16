package provider

import (
	"fmt"

	"github.com/google/go-github/v55/github"
	model "github.com/opengovern/og-describer-template/discovery/pkg/models"
	"github.com/opengovern/og-util/pkg/describe/enums"
	"github.com/shurcooL/githubv4"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)


type GitHubClient struct {
	RestClient    *github.Client
	GraphQLClient *githubv4.Client
	Token         string
}
var (
	triggerTypeKey string = "trigger_type"
)
func WithTriggerType(ctx context.Context, tt enums.DescribeTriggerType) context.Context {
	return context.WithValue(ctx, triggerTypeKey, tt)
}
func DescribeByGithub(describe func(context.Context, GitHubClient, string, *model.StreamSender) ([]model.Resource, error)) model.ResourceDescriber {
	return func(ctx context.Context, cfg model.IntegrationCredentials, triggerType enums.DescribeTriggerType, additionalParameters map[string]string, stream *model.StreamSender) ([]model.Resource, error) {
		ctx = WithTriggerType(ctx, triggerType)

		if cfg.PatToken == "" {
			return nil, fmt.Errorf("'token' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
		}

		// Create an OAuth2 token source
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: cfg.PatToken},
		)

		// Create an OAuth2 client
		tc := oauth2.NewClient(ctx, ts)

		// Create a new GitHub client
		restClient := github.NewClient(tc)
		graphQLClient := githubv4.NewClient(tc)

		client := GitHubClient{
			RestClient:    restClient,
			GraphQLClient: graphQLClient,
			Token:         cfg.PatToken,
		}

		organizationName := additionalParameters["OrganizationName"]
		var values []model.Resource
		result, err := describe(ctx, client, organizationName, stream)
		if err != nil {
			return nil, err
		}
		values = append(values, result...)
		return values, nil
	}
}

func DescribeSingleByRepo(describe func(context.Context, GitHubClient, string, string, string, *model.StreamSender) (*model.Resource, error)) model.SingleResourceDescriber {
	return func(ctx context.Context, cfg model.IntegrationCredentials, triggerType enums.DescribeTriggerType, additionalParameters map[string]string, resourceID string, stream *model.StreamSender) (*model.Resource, error) {
		ctx = WithTriggerType(ctx, triggerType)

		if cfg.PatToken == "" {
			return nil, fmt.Errorf("'token' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
		}

		// Create an OAuth2 token source
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: cfg.PatToken},
		)

		// Create an OAuth2 client
		tc := oauth2.NewClient(ctx, ts)

		// Create a new GitHub client
		restClient := github.NewClient(tc)
		graphQLClient := githubv4.NewClient(tc)

		client := GitHubClient{
			RestClient:    restClient,
			GraphQLClient: graphQLClient,
		}

		organizationName := additionalParameters["OrganizationName"]
		repoName := additionalParameters["RepositoryName"]
		result, err := describe(ctx, client, organizationName, repoName, resourceID, stream)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
}
