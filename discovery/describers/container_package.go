package describers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/opengovern/og-describer-template/discovery/pkg/models"
	model "github.com/opengovern/og-describer-template/discovery/provider"
	resilientbridge "github.com/opengovern/resilient-bridge"
	"github.com/opengovern/resilient-bridge/adapters"
)

// -----------------------------------------------------------------------------
// 1. GetContainerPackageList
// -----------------------------------------------------------------------------
func GetContainerPackageList(
	ctx context.Context,
	githubClient model.GitHubClient,
	organizationName string,
	stream *models.StreamSender,
) ([]models.Resource, error) {
	sdk := resilientbridge.NewResilientBridge()
	sdk.RegisterProvider("github", adapters.NewGitHubAdapter(githubClient.Token), &resilientbridge.ProviderConfig{
		UseProviderLimits: true,
		MaxRetries:        3,
		BaseBackoff:       0,
	})

	org := organizationName

	organization := ctx.Value("organization")
	if organization != nil {
		org = organization.(string)
		if org == "" {
			org = organizationName
		}
	}

	// [UPDATED] fetchPackages now does pagination
	packages := fetchPackages(sdk, org, "container")

	maxVersions := 1
	var allValues []models.Resource

	// Loop through each package and version
	for _, p := range packages {
		packageName := p.Name

		// [UPDATED] fetchVersions now does pagination
		versions := fetchVersions(sdk, org, "container", packageName)

		if len(versions) > maxVersions {
			versions = versions[:maxVersions]
		}

		for _, v := range versions {

			packageValues, err := getVersionOutput(githubClient.Token, org, packageName, v, stream)

			if err != nil {
				// If you want to fail fast, return err. Or just log and continue.
				log.Printf("Error getting version output for %s/%s: %v", packageName, v.Name, err)
				continue
			}
			//fmt.Println(packageValues)
			allValues = append(allValues, packageValues...)
		}
	}
	fmt.Println(allValues)
	return allValues, nil
}

// -----------------------------------------------------------------------------
// 2. getVersionOutput
//   - fetches details for each tag in a given version concurrently
//   - after concurrency, deduplicate by (version.ID, actualDigest)
//
// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------
func getVersionOutput(
	apiToken, org, packageName string,
	version model.PackageVersion,
	stream *models.StreamSender,
) ([]models.Resource, error) {

	var concurrentResults []models.Resource
	normalizedPackageName := strings.ToLower(packageName)
	tags := version.Metadata.Container.Tags
	if len(tags) == 0 {
		return concurrentResults, nil
	}

	// Prepare concurrency
	var wg sync.WaitGroup
	resultsChan := make(chan models.Resource, len(tags))
	errChan := make(chan error, 1)

	// For each tag, fetch the container details in a goroutine
	for _, tag := range tags {
		wg.Add(1)
		go func(t string) {
			defer wg.Done()
			normalizedTag := strings.ToLower(t)
			imageRef := fmt.Sprintf("ghcr.io/%s/%s:%s", org, normalizedPackageName, normalizedTag)

			ov, err := fetchAndAssembleOutput(apiToken, org, normalizedPackageName, version, imageRef)
			if err != nil {
				select {
				case errChan <- err:
				default:
				}
				return
			}

			// Convert ov -> models.Resource
			value := models.Resource{
				ID:          strconv.Itoa(ov.ID),
				Name:        ov.Name,
				Description: ov,
			}
			resultsChan <- value
		}(tag)
	}

	// Wait for goroutines
	go func() {
		wg.Wait()
		close(resultsChan)
		close(errChan)
	}()

	// If any goroutine sends an error, return the first one
	if err := <-errChan; err != nil {
		return nil, err
	}

	// Collect concurrency results
	for res := range resultsChan {
		// Always append to concurrentResults so it can be returned
		concurrentResults = append(concurrentResults, res)

		// If streaming, also stream each item
		if stream != nil {
			if e := (*stream)(res); e != nil {
				return nil, e
			}
		}
	}

	// Deduplicate by (version.ID + actualDigest)
	deduped := deduplicateVersionOutputsByDigest(concurrentResults, version.ID)
	return deduped, nil
}

// -----------------------------------------------------------------------------
// deduplicateVersionOutputsByDigest
//   - Expects an array of Resources whose .Description is ContainerPackageDescription.
//   - Collapses duplicates by (ID, actualDigest).
//
// -----------------------------------------------------------------------------
func deduplicateVersionOutputsByDigest(resources []models.Resource, versionID int) []models.Resource {
	// Key = "versionID|digest"
	type dedupKey struct {
		versionID string
		digest    string
	}

	dedupMap := make(map[dedupKey]*model.ContainerPackageDescription)
	var finalResults []models.Resource

	for _, r := range resources {
		// Type-assert .Description to JSONAllFieldsMarshaller
		desc := r.Description

		// Marshal then unmarshal so we can read it into ContainerPackageDescription
		rawBytes, err := json.Marshal(desc)
		if err != nil {
			continue
		}

		var cpd model.ContainerPackageDescription
		if err := json.Unmarshal(rawBytes, &cpd); err != nil {
			continue
		}

		dk := dedupKey{
			versionID: fmt.Sprintf("%d", versionID),
			digest:    cpd.ActualDigest(),
		}

		// If we already have this digest, append to AdditionalPackageURIs
		if existing, exists := dedupMap[dk]; exists {
			existing.AdditionalPackageURIs = append(existing.AdditionalPackageURIs, cpd.PackageURL)
		} else {
			dedupMap[dk] = &cpd
		}
	}

	// Convert the deduped map back to a []models.Resource
	for _, cpdPtr := range dedupMap {
		cpd := *cpdPtr
		res := models.Resource{
			ID:          strconv.Itoa(cpd.ID),
			Name:        cpd.Name,
			Description: cpd,
		}
		finalResults = append(finalResults, res)
	}

	return finalResults
}

// -----------------------------------------------------------------------------
// 4. fetchAndAssembleOutput
//   - no signature changes
//   - adds "AdditionalPackageURIs" field to ContainerPackageDescription, plus we
//     store real Docker digest from remote. (We do that by overriding cpd.Digest.)
//
// -----------------------------------------------------------------------------
func fetchAndAssembleOutput(
	apiToken, org, packageName string,
	version model.PackageVersion,
	imageRef string,
) (model.ContainerPackageDescription, error) {

	authOption := remote.WithAuth(&authn.Basic{
		Username: "github",
		Password: apiToken,
	})
	imageRef = strings.ToLower(imageRef)
	ref, err := name.ParseReference(imageRef, name.WeakValidation)
	if err != nil {
		return model.ContainerPackageDescription{},
			fmt.Errorf("error parsing reference %s: %w", imageRef, err)
	}

	desc, err := remote.Get(ref, authOption)
	if err != nil {
		return model.ContainerPackageDescription{},
			fmt.Errorf("error fetching manifest for %s: %w", imageRef, err)
	}

	// [UPDATED] We read the actual Docker registry digest from desc.Descriptor.Digest
	actualDigest := desc.Descriptor.Digest.String()

	var manifestStruct struct {
		SchemaVersion int    `json:"schemaVersion"`
		MediaType     string `json:"mediaType"`
		Config        struct {
			Size      int64  `json:"size"`
			Digest    string `json:"digest"`
			MediaType string `json:"mediaType"`
		} `json:"config"`
		Layers []struct {
			Size      int64  `json:"size"`
			Digest    string `json:"digest"`
			MediaType string `json:"mediaType"`
		} `json:"layers"`
	}
	if err := json.Unmarshal(desc.Manifest, &manifestStruct); err != nil {
		return model.ContainerPackageDescription{},
			fmt.Errorf("error unmarshaling manifest JSON: %w", err)
	}

	totalSize := manifestStruct.Config.Size
	for _, layer := range manifestStruct.Layers {
		totalSize += layer.Size
	}

	var manifestInterface interface{}
	if err := json.Unmarshal(desc.Manifest, &manifestInterface); err != nil {
		return model.ContainerPackageDescription{},
			fmt.Errorf("error unmarshaling manifest for output: %w", err)
	}

	// [UPDATED] The "digest" from the GitHub version might be inaccurate or just a tag name.
	// We'll store the real Docker digest in a separate field (ActualDigest).
	// For uniqueness, we’ll define "Digest" as the GH “version.Name” if you still want that.
	// Or store "Digest" as actualDigest if you prefer. Below we store GH’s version.Name in
	// GHVersionName, and store the actual registry digest in "Digest".
	// We'll also add AdditionalPackageURIs []string (will be deduplicated in a later step).

	ov := model.ContainerPackageDescription{
		ID:                    version.ID,
		Digest:                actualDigest, // store the real Docker digest
		AdditionalPackageURIs: []string{},   // Will be appended after dedup
		CreatedAt:             version.CreatedAt,
		UpdatedAt:             version.UpdatedAt,
		PackageURL:            version.HTMLURL,
		Name:                  imageRef,
		MediaType:             string(desc.Descriptor.MediaType),
		TotalSize:             totalSize,
		Metadata:              version.Metadata,
		Manifest:              manifestInterface,
	}
	return ov, nil
}

// -----------------------------------------------------------------------------
// 5. fetchPackages - updated to do pagination
// -----------------------------------------------------------------------------
func fetchPackages(sdk *resilientbridge.ResilientBridge, org, packageType string) []model.Package {
	var allPackages []model.Package
	page := 1
	perPage := 100

	for {
		req := &resilientbridge.NormalizedRequest{
			Method: "GET",
			Endpoint: fmt.Sprintf("/orgs/%s/packages?package_type=%s&per_page=%d&page=%d",
				org, packageType, perPage, page),
			Headers: map[string]string{"Accept": "application/vnd.github+json"},
		}

		resp, err := sdk.Request("github", req)
		if err != nil {
			log.Fatalf("Error listing packages: %v", err)
		}
		if resp.StatusCode >= 400 {
			log.Fatalf("HTTP error %d: %s", resp.StatusCode, string(resp.Data))
		}

		var packages []model.Package
		if err := json.Unmarshal(resp.Data, &packages); err != nil {
			log.Fatalf("Error parsing packages list response: %v", err)
		}
		if len(packages) == 0 {
			// no more data
			break
		}

		allPackages = append(allPackages, packages...)

		if len(packages) < perPage {
			// we got fewer than 100, so no more pages
			break
		}
		page++
	}
	return allPackages
}

// -----------------------------------------------------------------------------
// 6. fetchVersions - updated to do pagination
// -----------------------------------------------------------------------------
func fetchVersions(
	sdk *resilientbridge.ResilientBridge,
	org, packageType, packageName string,
) []model.PackageVersion {

	packageNameEncoded := url.PathEscape(packageName)
	var allVersions []model.PackageVersion
	page := 1
	perPage := 100

	for {
		req := &resilientbridge.NormalizedRequest{
			Method: "GET",
			Endpoint: fmt.Sprintf(
				"/orgs/%s/packages/%s/%s/versions?per_page=%d&page=%d",
				org, packageType, packageNameEncoded, perPage, page,
			),
			Headers: map[string]string{"Accept": "application/vnd.github+json"},
		}

		resp, err := sdk.Request("github", req)
		if err != nil {
			log.Fatalf("Error listing package versions: %v", err)
		}
		if resp.StatusCode >= 400 {
			log.Fatalf("HTTP error %d: %s", resp.StatusCode, string(resp.Data))
		}

		var versions []model.PackageVersion
		if err := json.Unmarshal(resp.Data, &versions); err != nil {
			log.Fatalf("Error parsing package versions response: %v", err)
		}

		if len(versions) == 0 {
			// no more data
			break
		}

		allVersions = append(allVersions, versions...)

		if len(versions) < perPage {
			// we got fewer than 100, so no more pages
			break
		}
		page++
	}

	return allVersions
}
