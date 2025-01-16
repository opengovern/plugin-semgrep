package maps

import (
	"github.com/opengovern/og-describer-template/discovery/pkg/es"
)

var ResourceTypesToTables = map[string]string{
  "Github/Actions/Artifact": "github_actions_artifact",
  
}

var ResourceTypeToDescription = map[string]interface{}{
  "Github/Actions/Artifact": opengovernance.Artifact{},
  
}

var TablesToResourceTypes = map[string]string{
  "github_actions_artifact": "Github/Actions/Artifact",
 
}
