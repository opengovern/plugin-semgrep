package global

import "github.com/opengovern/og-util/pkg/integration"

const (
	IntegrationTypeLower = "semgrep"                                    // example: aws, azure
	IntegrationName      = integration.Type("semgrep_account")          // example: aws_account, github_account
	OGPluginRepoURL      = "github.com/opengovern/og-describer-semgrep" // example: github.com/opengovern/og-describer-aws
)

type IntegrationCredentials struct {
	// TODO
}
