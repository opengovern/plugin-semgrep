package configs

import "github.com/opengovern/og-util/pkg/integration"

const (
	IntegrationTypeLower = "integrationType"                    // example: aws, azure
	IntegrationName      = integration.Type("INTEGRATION_NAME") // example: AWS_ACCOUNT, AZURE_SUBSCRIPTION
	OGPluginRepoURL      = "repo-url"                           // example: github.com/opengovern/og-describer-aws
)
