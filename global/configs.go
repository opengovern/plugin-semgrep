package global

import "github.com/opengovern/og-util/pkg/integration"

const (
	IntegrationTypeLower = "template"                                    // example: aws, azure
	IntegrationName      = integration.Type("template,github")          // example: aws_account, github_account
	OGPluginRepoURL      = "github.com/opengovern/og-describer-template" // example: github.com/opengovern/og-describer-aws
)

type IntegrationCredentials struct {
	// TODO
}
