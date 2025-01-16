package constants

import "github.com/opengovern/og-util/pkg/integration"
import _ "embed"

//go:embed ui-spec.json
var UISpec []byte

//go:embed manifest.yaml
var Manifest []byte

//go:embed Setup.md
var SetupMd []byte

const (
	IntegrationName = integration.Type("template") // example: aws_cloud, azure_subscription, github_account
)

const (
	DescriberDeploymentName = "og-describer-template"
	DescriberRunCommand     = "/og-describer-template"
)
