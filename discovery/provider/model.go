
// Implement types for each resource

package provider



type Metadata struct{}

type ArtifactDockerFileDescription struct {
	Sha  *string
	Name *string
	LastUpdatedAt *string
	HTMLURL *string
	DockerfileContent       string
	DockerfileContentBase64 *string
	Repository              map[string]interface{}
	Images                  []string 
}
