package global

import "github.com/opengovern/og-util/pkg/integration"

const (
	IntegrationTypeLower = "github"                                    // example: aws, azure
	IntegrationName      = integration.Type("github_account")          // example: AWS_ACCOUNT, AZURE_SUBSCRIPTION
	OGPluginRepoURL      = "github.com/opengovern/og-describer-template" // example: github.com/opengovern/og-describer-aws
)

type IntegrationCredentials struct {
	PatToken string `json:"pat_token"`
}
