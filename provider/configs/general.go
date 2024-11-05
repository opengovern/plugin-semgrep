package configs

import "github.com/opengovern/og-util/pkg/integration"

const (
	Provider        = "provider"                           // example: aws, azure
	Cloud           = "ProviderCloud"                      // example: AWSCloud, AzureCloud
	UpperProvider   = "Provider"                           // example: AWS, Azure
	IntegrationName = integration.Type("INTEGRATION_NAME") // example: AWS_ACCOUNT, AZURE_SUBSCRIPTION
	OGPluginRepoURL = "repo-url"                           // example: github.com/opengovern/og-describer-aws
)
