package maps

var ResourceTypesToTables = map[string]string{
	"Github/Artifact/DockerFile": "template_artifact_dockerfile",
}

var ResourceTypeToDescription = map[string]interface{}{}

var TablesToResourceTypes = map[string]string{
	"template_artifact_dockerfile": "Github/Artifact/DockerFile",
}
