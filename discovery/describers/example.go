// artifact_dockerfile.go
package describers

import (
	"context"
	"fmt"
	"github.com/opengovern/og-describer-template/discovery/pkg/models"
	model "github.com/opengovern/og-describer-template/discovery/provider"
)



// ListArtifactDockerFiles performs a single code search across the organization
// for "filename:Dockerfile" and processes each result. Each Dockerfile is
// streamed immediately upon processing, and also added to the final slice.
func ListArtifactDockerFiles(
	ctx context.Context,
	client model.Client,
	extra string,
	stream *models.StreamSender,
) ([]models.Resource, error) {

	var allValues []models.Resource
	// TODO implement the logic to get the 
	var resource *models.Resource
	if stream != nil {
				if err := (*stream)(*resource); err != nil {
					return allValues, fmt.Errorf("error streaming resource: %w", err)
				}
			}

	// Return everything, even though we streamed each file already
	return allValues, nil
}

