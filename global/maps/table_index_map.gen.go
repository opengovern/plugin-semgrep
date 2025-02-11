package maps

import (
	"github.com/opengovern/og-describer-semgrep/discovery/pkg/es"
)

var ResourceTypesToTables = map[string]string{
  "Semgrep/Deployment": "semgrep_deployment",
  "Semgrep/Project": "semgrep_project",
  "Semgrep/Policy": "semgrep_policy",
  "Semgrep/Scan": "semgrep_scan",
  "Semgrep/Finding": "semgrep_finding",
}

var ResourceTypeToDescription = map[string]interface{}{
  "Semgrep/Deployment": opengovernance.Deployment{},
  "Semgrep/Project": opengovernance.Project{},
  "Semgrep/Policy": opengovernance.Policy{},
  "Semgrep/Scan": opengovernance.Scan{},
  "Semgrep/Finding": opengovernance.Finding{},
}

var TablesToResourceTypes = map[string]string{
  "semgrep_deployment": "Semgrep/Deployment",
  "semgrep_project": "Semgrep/Project",
  "semgrep_policy": "Semgrep/Policy",
  "semgrep_scan": "Semgrep/Scan",
  "semgrep_finding": "Semgrep/Finding",
}
