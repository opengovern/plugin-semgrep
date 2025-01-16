
// Implement types for each resource

package provider



type Metadata struct{}

type ArtifactDescription struct {
	ID                 int64
	NodeID             *string
	Name               *string
	SizeInBytes        int64
	ArchiveDownloadURL *string
	Expired            bool
	CreatedAt          *string
	ExpiresAt          *string
	RepoFullName       *string
}
