package maps

import (
	"github.com/opengovern/og-describer-template/discovery/pkg/es"
)

var ResourceTypesToTables = map[string]string{
  "Github/Artifact/DockerFile": "template_artifact_dockerfile",
}

var ResourceTypeToDescription = map[string]interface{}{
  "Github/Artifact/DockerFile": opengovernance.ArtifactDockerFile{},
}

var TablesToResourceTypes = map[string]string{
  "template_artifact_dockerfile": "Github/Artifact/DockerFile",
}
