// Code is generated by go generate. DO NOT EDIT.
package opengovernance

import (
	"context"
	semgrep "github.com/opengovern/og-describer-semgrep/discovery/provider"
	essdk "github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	steampipesdk "github.com/opengovern/og-util/pkg/steampipe"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"runtime"
)

type Client struct {
	essdk.Client
}

// ==========================  START: Deployment =============================

type Deployment struct {
	ResourceID      string                        `json:"resource_id"`
	PlatformID      string                        `json:"platform_id"`
	Description     semgrep.DeploymentDescription `json:"Description"`
	Metadata        semgrep.Metadata              `json:"metadata"`
	DescribedBy     string                        `json:"described_by"`
	ResourceType    string                        `json:"resource_type"`
	IntegrationType string                        `json:"integration_type"`
	IntegrationID   string                        `json:"integration_id"`
}

type DeploymentHit struct {
	ID      string        `json:"_id"`
	Score   float64       `json:"_score"`
	Index   string        `json:"_index"`
	Type    string        `json:"_type"`
	Version int64         `json:"_version,omitempty"`
	Source  Deployment    `json:"_source"`
	Sort    []interface{} `json:"sort"`
}

type DeploymentHits struct {
	Total essdk.SearchTotal `json:"total"`
	Hits  []DeploymentHit   `json:"hits"`
}

type DeploymentSearchResponse struct {
	PitID string         `json:"pit_id"`
	Hits  DeploymentHits `json:"hits"`
}

type DeploymentPaginator struct {
	paginator *essdk.BaseESPaginator
}

func (k Client) NewDeploymentPaginator(filters []essdk.BoolFilter, limit *int64) (DeploymentPaginator, error) {
	paginator, err := essdk.NewPaginator(k.ES(), "semgrep_deployment", filters, limit)
	if err != nil {
		return DeploymentPaginator{}, err
	}

	p := DeploymentPaginator{
		paginator: paginator,
	}

	return p, nil
}

func (p DeploymentPaginator) HasNext() bool {
	return !p.paginator.Done()
}

func (p DeploymentPaginator) Close(ctx context.Context) error {
	return p.paginator.Deallocate(ctx)
}

func (p DeploymentPaginator) NextPage(ctx context.Context) ([]Deployment, error) {
	var response DeploymentSearchResponse
	err := p.paginator.Search(ctx, &response)
	if err != nil {
		return nil, err
	}

	var values []Deployment
	for _, hit := range response.Hits.Hits {
		values = append(values, hit.Source)
	}

	hits := int64(len(response.Hits.Hits))
	if hits > 0 {
		p.paginator.UpdateState(hits, response.Hits.Hits[hits-1].Sort, response.PitID)
	} else {
		p.paginator.UpdateState(hits, nil, "")
	}

	return values, nil
}

var listDeploymentFilters = map[string]string{
	"findings": "Description.Findings",
	"id":       "Description.ID",
	"name":     "Description.Name",
	"slug":     "Description.Slug",
}

func ListDeployment(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("ListDeployment")
	runtime.GC()

	// create service
	cfg := essdk.GetConfig(d.Connection)
	ke, err := essdk.NewClientCached(cfg, d.ConnectionCache, ctx)
	if err != nil {
		plugin.Logger(ctx).Error("ListDeployment NewClientCached", "error", err)
		return nil, err
	}
	k := Client{Client: ke}

	sc, err := steampipesdk.NewSelfClientCached(ctx, d.ConnectionCache)
	if err != nil {
		plugin.Logger(ctx).Error("ListDeployment NewSelfClientCached", "error", err)
		return nil, err
	}
	integrationId, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyIntegrationID)
	if err != nil {
		plugin.Logger(ctx).Error("ListDeployment GetConfigTableValueOrNil for OpenGovernanceConfigKeyIntegrationID", "error", err)
		return nil, err
	}
	encodedResourceCollectionFilters, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyResourceCollectionFilters)
	if err != nil {
		plugin.Logger(ctx).Error("ListDeployment GetConfigTableValueOrNil for OpenGovernanceConfigKeyResourceCollectionFilters", "error", err)
		return nil, err
	}
	clientType, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyClientType)
	if err != nil {
		plugin.Logger(ctx).Error("ListDeployment GetConfigTableValueOrNil for OpenGovernanceConfigKeyClientType", "error", err)
		return nil, err
	}

	paginator, err := k.NewDeploymentPaginator(essdk.BuildFilter(ctx, d.QueryContext, listDeploymentFilters, integrationId, encodedResourceCollectionFilters, clientType), d.QueryContext.Limit)
	if err != nil {
		plugin.Logger(ctx).Error("ListDeployment NewDeploymentPaginator", "error", err)
		return nil, err
	}

	for paginator.HasNext() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("ListDeployment paginator.NextPage", "error", err)
			return nil, err
		}

		for _, v := range page {
			d.StreamListItem(ctx, v)
		}
	}

	err = paginator.Close(ctx)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

var getDeploymentFilters = map[string]string{
	"findings": "Description.Findings",
	"id":       "Description.ID",
	"name":     "Description.Name",
	"slug":     "Description.Slug",
}

func GetDeployment(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("GetDeployment")
	runtime.GC()
	// create service
	cfg := essdk.GetConfig(d.Connection)
	ke, err := essdk.NewClientCached(cfg, d.ConnectionCache, ctx)
	if err != nil {
		return nil, err
	}
	k := Client{Client: ke}

	sc, err := steampipesdk.NewSelfClientCached(ctx, d.ConnectionCache)
	if err != nil {
		return nil, err
	}
	integrationId, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyIntegrationID)
	if err != nil {
		return nil, err
	}
	encodedResourceCollectionFilters, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyResourceCollectionFilters)
	if err != nil {
		return nil, err
	}
	clientType, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyClientType)
	if err != nil {
		return nil, err
	}

	limit := int64(1)
	paginator, err := k.NewDeploymentPaginator(essdk.BuildFilter(ctx, d.QueryContext, getDeploymentFilters, integrationId, encodedResourceCollectionFilters, clientType), &limit)
	if err != nil {
		return nil, err
	}

	for paginator.HasNext() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page {
			return v, nil
		}
	}

	err = paginator.Close(ctx)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ==========================  END: Deployment =============================

// ==========================  START: Project =============================

type Project struct {
	ResourceID      string                     `json:"resource_id"`
	PlatformID      string                     `json:"platform_id"`
	Description     semgrep.ProjectDescription `json:"Description"`
	Metadata        semgrep.Metadata           `json:"metadata"`
	DescribedBy     string                     `json:"described_by"`
	ResourceType    string                     `json:"resource_type"`
	IntegrationType string                     `json:"integration_type"`
	IntegrationID   string                     `json:"integration_id"`
}

type ProjectHit struct {
	ID      string        `json:"_id"`
	Score   float64       `json:"_score"`
	Index   string        `json:"_index"`
	Type    string        `json:"_type"`
	Version int64         `json:"_version,omitempty"`
	Source  Project       `json:"_source"`
	Sort    []interface{} `json:"sort"`
}

type ProjectHits struct {
	Total essdk.SearchTotal `json:"total"`
	Hits  []ProjectHit      `json:"hits"`
}

type ProjectSearchResponse struct {
	PitID string      `json:"pit_id"`
	Hits  ProjectHits `json:"hits"`
}

type ProjectPaginator struct {
	paginator *essdk.BaseESPaginator
}

func (k Client) NewProjectPaginator(filters []essdk.BoolFilter, limit *int64) (ProjectPaginator, error) {
	paginator, err := essdk.NewPaginator(k.ES(), "semgrep_project", filters, limit)
	if err != nil {
		return ProjectPaginator{}, err
	}

	p := ProjectPaginator{
		paginator: paginator,
	}

	return p, nil
}

func (p ProjectPaginator) HasNext() bool {
	return !p.paginator.Done()
}

func (p ProjectPaginator) Close(ctx context.Context) error {
	return p.paginator.Deallocate(ctx)
}

func (p ProjectPaginator) NextPage(ctx context.Context) ([]Project, error) {
	var response ProjectSearchResponse
	err := p.paginator.Search(ctx, &response)
	if err != nil {
		return nil, err
	}

	var values []Project
	for _, hit := range response.Hits.Hits {
		values = append(values, hit.Source)
	}

	hits := int64(len(response.Hits.Hits))
	if hits > 0 {
		p.paginator.UpdateState(hits, response.Hits.Hits[hits-1].Sort, response.PitID)
	} else {
		p.paginator.UpdateState(hits, nil, "")
	}

	return values, nil
}

var listProjectFilters = map[string]string{
	"created_at":     "Description.CreatedAt",
	"default_branch": "Description.DefaultBranch",
	"id":             "Description.ID",
	"latest_scan_at": "Description.LatestScanAt",
	"name":           "Description.Name",
	"primary_branch": "Description.PrimaryBranch",
	"tags":           "Description.Tags",
	"url":            "Description.URL",
}

func ListProject(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("ListProject")
	runtime.GC()

	// create service
	cfg := essdk.GetConfig(d.Connection)
	ke, err := essdk.NewClientCached(cfg, d.ConnectionCache, ctx)
	if err != nil {
		plugin.Logger(ctx).Error("ListProject NewClientCached", "error", err)
		return nil, err
	}
	k := Client{Client: ke}

	sc, err := steampipesdk.NewSelfClientCached(ctx, d.ConnectionCache)
	if err != nil {
		plugin.Logger(ctx).Error("ListProject NewSelfClientCached", "error", err)
		return nil, err
	}
	integrationId, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyIntegrationID)
	if err != nil {
		plugin.Logger(ctx).Error("ListProject GetConfigTableValueOrNil for OpenGovernanceConfigKeyIntegrationID", "error", err)
		return nil, err
	}
	encodedResourceCollectionFilters, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyResourceCollectionFilters)
	if err != nil {
		plugin.Logger(ctx).Error("ListProject GetConfigTableValueOrNil for OpenGovernanceConfigKeyResourceCollectionFilters", "error", err)
		return nil, err
	}
	clientType, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyClientType)
	if err != nil {
		plugin.Logger(ctx).Error("ListProject GetConfigTableValueOrNil for OpenGovernanceConfigKeyClientType", "error", err)
		return nil, err
	}

	paginator, err := k.NewProjectPaginator(essdk.BuildFilter(ctx, d.QueryContext, listProjectFilters, integrationId, encodedResourceCollectionFilters, clientType), d.QueryContext.Limit)
	if err != nil {
		plugin.Logger(ctx).Error("ListProject NewProjectPaginator", "error", err)
		return nil, err
	}

	for paginator.HasNext() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("ListProject paginator.NextPage", "error", err)
			return nil, err
		}

		for _, v := range page {
			d.StreamListItem(ctx, v)
		}
	}

	err = paginator.Close(ctx)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

var getProjectFilters = map[string]string{
	"created_at":     "Description.CreatedAt",
	"default_branch": "Description.DefaultBranch",
	"id":             "Description.ID",
	"latest_scan_at": "Description.LatestScanAt",
	"name":           "Description.Name",
	"primary_branch": "Description.PrimaryBranch",
	"tags":           "Description.Tags",
	"url":            "Description.URL",
}

func GetProject(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("GetProject")
	runtime.GC()
	// create service
	cfg := essdk.GetConfig(d.Connection)
	ke, err := essdk.NewClientCached(cfg, d.ConnectionCache, ctx)
	if err != nil {
		return nil, err
	}
	k := Client{Client: ke}

	sc, err := steampipesdk.NewSelfClientCached(ctx, d.ConnectionCache)
	if err != nil {
		return nil, err
	}
	integrationId, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyIntegrationID)
	if err != nil {
		return nil, err
	}
	encodedResourceCollectionFilters, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyResourceCollectionFilters)
	if err != nil {
		return nil, err
	}
	clientType, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyClientType)
	if err != nil {
		return nil, err
	}

	limit := int64(1)
	paginator, err := k.NewProjectPaginator(essdk.BuildFilter(ctx, d.QueryContext, getProjectFilters, integrationId, encodedResourceCollectionFilters, clientType), &limit)
	if err != nil {
		return nil, err
	}

	for paginator.HasNext() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page {
			return v, nil
		}
	}

	err = paginator.Close(ctx)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ==========================  END: Project =============================

// ==========================  START: Policy =============================

type Policy struct {
	ResourceID      string                    `json:"resource_id"`
	PlatformID      string                    `json:"platform_id"`
	Description     semgrep.PolicyDescription `json:"Description"`
	Metadata        semgrep.Metadata          `json:"metadata"`
	DescribedBy     string                    `json:"described_by"`
	ResourceType    string                    `json:"resource_type"`
	IntegrationType string                    `json:"integration_type"`
	IntegrationID   string                    `json:"integration_id"`
}

type PolicyHit struct {
	ID      string        `json:"_id"`
	Score   float64       `json:"_score"`
	Index   string        `json:"_index"`
	Type    string        `json:"_type"`
	Version int64         `json:"_version,omitempty"`
	Source  Policy        `json:"_source"`
	Sort    []interface{} `json:"sort"`
}

type PolicyHits struct {
	Total essdk.SearchTotal `json:"total"`
	Hits  []PolicyHit       `json:"hits"`
}

type PolicySearchResponse struct {
	PitID string     `json:"pit_id"`
	Hits  PolicyHits `json:"hits"`
}

type PolicyPaginator struct {
	paginator *essdk.BaseESPaginator
}

func (k Client) NewPolicyPaginator(filters []essdk.BoolFilter, limit *int64) (PolicyPaginator, error) {
	paginator, err := essdk.NewPaginator(k.ES(), "semgrep_policy", filters, limit)
	if err != nil {
		return PolicyPaginator{}, err
	}

	p := PolicyPaginator{
		paginator: paginator,
	}

	return p, nil
}

func (p PolicyPaginator) HasNext() bool {
	return !p.paginator.Done()
}

func (p PolicyPaginator) Close(ctx context.Context) error {
	return p.paginator.Deallocate(ctx)
}

func (p PolicyPaginator) NextPage(ctx context.Context) ([]Policy, error) {
	var response PolicySearchResponse
	err := p.paginator.Search(ctx, &response)
	if err != nil {
		return nil, err
	}

	var values []Policy
	for _, hit := range response.Hits.Hits {
		values = append(values, hit.Source)
	}

	hits := int64(len(response.Hits.Hits))
	if hits > 0 {
		p.paginator.UpdateState(hits, response.Hits.Hits[hits-1].Sort, response.PitID)
	} else {
		p.paginator.UpdateState(hits, nil, "")
	}

	return values, nil
}

var listPolicyFilters = map[string]string{
	"id":           "Description.ID",
	"is_default":   "Description.IsDefault",
	"name":         "Description.Name",
	"product_type": "Description.ProductType",
	"slug":         "Description.Slug",
}

func ListPolicy(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("ListPolicy")
	runtime.GC()

	// create service
	cfg := essdk.GetConfig(d.Connection)
	ke, err := essdk.NewClientCached(cfg, d.ConnectionCache, ctx)
	if err != nil {
		plugin.Logger(ctx).Error("ListPolicy NewClientCached", "error", err)
		return nil, err
	}
	k := Client{Client: ke}

	sc, err := steampipesdk.NewSelfClientCached(ctx, d.ConnectionCache)
	if err != nil {
		plugin.Logger(ctx).Error("ListPolicy NewSelfClientCached", "error", err)
		return nil, err
	}
	integrationId, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyIntegrationID)
	if err != nil {
		plugin.Logger(ctx).Error("ListPolicy GetConfigTableValueOrNil for OpenGovernanceConfigKeyIntegrationID", "error", err)
		return nil, err
	}
	encodedResourceCollectionFilters, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyResourceCollectionFilters)
	if err != nil {
		plugin.Logger(ctx).Error("ListPolicy GetConfigTableValueOrNil for OpenGovernanceConfigKeyResourceCollectionFilters", "error", err)
		return nil, err
	}
	clientType, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyClientType)
	if err != nil {
		plugin.Logger(ctx).Error("ListPolicy GetConfigTableValueOrNil for OpenGovernanceConfigKeyClientType", "error", err)
		return nil, err
	}

	paginator, err := k.NewPolicyPaginator(essdk.BuildFilter(ctx, d.QueryContext, listPolicyFilters, integrationId, encodedResourceCollectionFilters, clientType), d.QueryContext.Limit)
	if err != nil {
		plugin.Logger(ctx).Error("ListPolicy NewPolicyPaginator", "error", err)
		return nil, err
	}

	for paginator.HasNext() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("ListPolicy paginator.NextPage", "error", err)
			return nil, err
		}

		for _, v := range page {
			d.StreamListItem(ctx, v)
		}
	}

	err = paginator.Close(ctx)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

var getPolicyFilters = map[string]string{
	"id":           "Description.ID",
	"is_default":   "Description.IsDefault",
	"name":         "Description.Name",
	"product_type": "Description.ProductType",
	"slug":         "Description.Slug",
}

func GetPolicy(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("GetPolicy")
	runtime.GC()
	// create service
	cfg := essdk.GetConfig(d.Connection)
	ke, err := essdk.NewClientCached(cfg, d.ConnectionCache, ctx)
	if err != nil {
		return nil, err
	}
	k := Client{Client: ke}

	sc, err := steampipesdk.NewSelfClientCached(ctx, d.ConnectionCache)
	if err != nil {
		return nil, err
	}
	integrationId, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyIntegrationID)
	if err != nil {
		return nil, err
	}
	encodedResourceCollectionFilters, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyResourceCollectionFilters)
	if err != nil {
		return nil, err
	}
	clientType, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyClientType)
	if err != nil {
		return nil, err
	}

	limit := int64(1)
	paginator, err := k.NewPolicyPaginator(essdk.BuildFilter(ctx, d.QueryContext, getPolicyFilters, integrationId, encodedResourceCollectionFilters, clientType), &limit)
	if err != nil {
		return nil, err
	}

	for paginator.HasNext() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page {
			return v, nil
		}
	}

	err = paginator.Close(ctx)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ==========================  END: Policy =============================

// ==========================  START: Scan =============================

type Scan struct {
	ResourceID      string                  `json:"resource_id"`
	PlatformID      string                  `json:"platform_id"`
	Description     semgrep.ScanDescription `json:"Description"`
	Metadata        semgrep.Metadata        `json:"metadata"`
	DescribedBy     string                  `json:"described_by"`
	ResourceType    string                  `json:"resource_type"`
	IntegrationType string                  `json:"integration_type"`
	IntegrationID   string                  `json:"integration_id"`
}

type ScanHit struct {
	ID      string        `json:"_id"`
	Score   float64       `json:"_score"`
	Index   string        `json:"_index"`
	Type    string        `json:"_type"`
	Version int64         `json:"_version,omitempty"`
	Source  Scan          `json:"_source"`
	Sort    []interface{} `json:"sort"`
}

type ScanHits struct {
	Total essdk.SearchTotal `json:"total"`
	Hits  []ScanHit         `json:"hits"`
}

type ScanSearchResponse struct {
	PitID string   `json:"pit_id"`
	Hits  ScanHits `json:"hits"`
}

type ScanPaginator struct {
	paginator *essdk.BaseESPaginator
}

func (k Client) NewScanPaginator(filters []essdk.BoolFilter, limit *int64) (ScanPaginator, error) {
	paginator, err := essdk.NewPaginator(k.ES(), "semgrep_scan", filters, limit)
	if err != nil {
		return ScanPaginator{}, err
	}

	p := ScanPaginator{
		paginator: paginator,
	}

	return p, nil
}

func (p ScanPaginator) HasNext() bool {
	return !p.paginator.Done()
}

func (p ScanPaginator) Close(ctx context.Context) error {
	return p.paginator.Deallocate(ctx)
}

func (p ScanPaginator) NextPage(ctx context.Context) ([]Scan, error) {
	var response ScanSearchResponse
	err := p.paginator.Search(ctx, &response)
	if err != nil {
		return nil, err
	}

	var values []Scan
	for _, hit := range response.Hits.Hits {
		values = append(values, hit.Source)
	}

	hits := int64(len(response.Hits.Hits))
	if hits > 0 {
		p.paginator.UpdateState(hits, response.Hits.Hits[hits-1].Sort, response.PitID)
	} else {
		p.paginator.UpdateState(hits, nil, "")
	}

	return values, nil
}

var listScanFilters = map[string]string{
	"branch":          "Description.Branch",
	"commit":          "Description.Commit",
	"completed_at":    "Description.CompletedAt",
	"deployment_id":   "Description.DeploymentID",
	"exit_code":       "Description.ExitCode",
	"findings_counts": "Description.FindingsCounts",
	"id":              "Description.ID",
	"is_full_scan":    "Description.IsFullScan",
	"repository_id":   "Description.RepositoryID",
	"started_at":      "Description.StartedAt",
	"status":          "Description.Status",
	"total_time":      "Description.TotalTime",
}

func ListScan(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("ListScan")
	runtime.GC()

	// create service
	cfg := essdk.GetConfig(d.Connection)
	ke, err := essdk.NewClientCached(cfg, d.ConnectionCache, ctx)
	if err != nil {
		plugin.Logger(ctx).Error("ListScan NewClientCached", "error", err)
		return nil, err
	}
	k := Client{Client: ke}

	sc, err := steampipesdk.NewSelfClientCached(ctx, d.ConnectionCache)
	if err != nil {
		plugin.Logger(ctx).Error("ListScan NewSelfClientCached", "error", err)
		return nil, err
	}
	integrationId, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyIntegrationID)
	if err != nil {
		plugin.Logger(ctx).Error("ListScan GetConfigTableValueOrNil for OpenGovernanceConfigKeyIntegrationID", "error", err)
		return nil, err
	}
	encodedResourceCollectionFilters, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyResourceCollectionFilters)
	if err != nil {
		plugin.Logger(ctx).Error("ListScan GetConfigTableValueOrNil for OpenGovernanceConfigKeyResourceCollectionFilters", "error", err)
		return nil, err
	}
	clientType, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyClientType)
	if err != nil {
		plugin.Logger(ctx).Error("ListScan GetConfigTableValueOrNil for OpenGovernanceConfigKeyClientType", "error", err)
		return nil, err
	}

	paginator, err := k.NewScanPaginator(essdk.BuildFilter(ctx, d.QueryContext, listScanFilters, integrationId, encodedResourceCollectionFilters, clientType), d.QueryContext.Limit)
	if err != nil {
		plugin.Logger(ctx).Error("ListScan NewScanPaginator", "error", err)
		return nil, err
	}

	for paginator.HasNext() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("ListScan paginator.NextPage", "error", err)
			return nil, err
		}

		for _, v := range page {
			d.StreamListItem(ctx, v)
		}
	}

	err = paginator.Close(ctx)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

var getScanFilters = map[string]string{
	"branch":          "Description.Branch",
	"commit":          "Description.Commit",
	"completed_at":    "Description.CompletedAt",
	"deployment_id":   "Description.DeploymentID",
	"exit_code":       "Description.ExitCode",
	"findings_counts": "Description.FindingsCounts",
	"id":              "Description.ID",
	"is_full_scan":    "Description.IsFullScan",
	"repository_id":   "Description.RepositoryID",
	"started_at":      "Description.StartedAt",
	"status":          "Description.Status",
	"total_time":      "Description.TotalTime",
}

func GetScan(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("GetScan")
	runtime.GC()
	// create service
	cfg := essdk.GetConfig(d.Connection)
	ke, err := essdk.NewClientCached(cfg, d.ConnectionCache, ctx)
	if err != nil {
		return nil, err
	}
	k := Client{Client: ke}

	sc, err := steampipesdk.NewSelfClientCached(ctx, d.ConnectionCache)
	if err != nil {
		return nil, err
	}
	integrationId, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyIntegrationID)
	if err != nil {
		return nil, err
	}
	encodedResourceCollectionFilters, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyResourceCollectionFilters)
	if err != nil {
		return nil, err
	}
	clientType, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyClientType)
	if err != nil {
		return nil, err
	}

	limit := int64(1)
	paginator, err := k.NewScanPaginator(essdk.BuildFilter(ctx, d.QueryContext, getScanFilters, integrationId, encodedResourceCollectionFilters, clientType), &limit)
	if err != nil {
		return nil, err
	}

	for paginator.HasNext() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page {
			return v, nil
		}
	}

	err = paginator.Close(ctx)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ==========================  END: Scan =============================

// ==========================  START: Finding =============================

type Finding struct {
	ResourceID      string                     `json:"resource_id"`
	PlatformID      string                     `json:"platform_id"`
	Description     semgrep.FindingDescription `json:"Description"`
	Metadata        semgrep.Metadata           `json:"metadata"`
	DescribedBy     string                     `json:"described_by"`
	ResourceType    string                     `json:"resource_type"`
	IntegrationType string                     `json:"integration_type"`
	IntegrationID   string                     `json:"integration_id"`
}

type FindingHit struct {
	ID      string        `json:"_id"`
	Score   float64       `json:"_score"`
	Index   string        `json:"_index"`
	Type    string        `json:"_type"`
	Version int64         `json:"_version,omitempty"`
	Source  Finding       `json:"_source"`
	Sort    []interface{} `json:"sort"`
}

type FindingHits struct {
	Total essdk.SearchTotal `json:"total"`
	Hits  []FindingHit      `json:"hits"`
}

type FindingSearchResponse struct {
	PitID string      `json:"pit_id"`
	Hits  FindingHits `json:"hits"`
}

type FindingPaginator struct {
	paginator *essdk.BaseESPaginator
}

func (k Client) NewFindingPaginator(filters []essdk.BoolFilter, limit *int64) (FindingPaginator, error) {
	paginator, err := essdk.NewPaginator(k.ES(), "semgrep_finding", filters, limit)
	if err != nil {
		return FindingPaginator{}, err
	}

	p := FindingPaginator{
		paginator: paginator,
	}

	return p, nil
}

func (p FindingPaginator) HasNext() bool {
	return !p.paginator.Done()
}

func (p FindingPaginator) Close(ctx context.Context) error {
	return p.paginator.Deallocate(ctx)
}

func (p FindingPaginator) NextPage(ctx context.Context) ([]Finding, error) {
	var response FindingSearchResponse
	err := p.paginator.Search(ctx, &response)
	if err != nil {
		return nil, err
	}

	var values []Finding
	for _, hit := range response.Hits.Hits {
		values = append(values, hit.Source)
	}

	hits := int64(len(response.Hits.Hits))
	if hits > 0 {
		p.paginator.UpdateState(hits, response.Hits.Hits[hits-1].Sort, response.PitID)
	} else {
		p.paginator.UpdateState(hits, nil, "")
	}

	return values, nil
}

var listFindingFilters = map[string]string{
	"assistant":          "Description.Assistant",
	"categories":         "Description.Categories",
	"confidence":         "Description.Confidence",
	"created_at":         "Description.CreatedAt",
	"external_ticket":    "Description.ExternalTicket",
	"first_seen_scan_id": "Description.FirstSeenScanID",
	"id":                 "Description.ID",
	"line_of_code_url":   "Description.LineOfCodeURL",
	"location":           "Description.Location",
	"match_based_id":     "Description.MatchBasedID",
	"ref":                "Description.Ref",
	"relevant_since":     "Description.RelevantSince",
	"repository":         "Description.Repository",
	"rule":               "Description.Rule",
	"rule_message":       "Description.RuleMessage",
	"rule_name":          "Description.RuleName",
	"severity":           "Description.Severity",
	"sourcing_policy":    "Description.SourcingPolicy",
	"state":              "Description.State",
	"status":             "Description.Status",
	"syntactic_id":       "Description.SyntacticID",
	"triage_state":       "Description.TriageState",
}

func ListFinding(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("ListFinding")
	runtime.GC()

	// create service
	cfg := essdk.GetConfig(d.Connection)
	ke, err := essdk.NewClientCached(cfg, d.ConnectionCache, ctx)
	if err != nil {
		plugin.Logger(ctx).Error("ListFinding NewClientCached", "error", err)
		return nil, err
	}
	k := Client{Client: ke}

	sc, err := steampipesdk.NewSelfClientCached(ctx, d.ConnectionCache)
	if err != nil {
		plugin.Logger(ctx).Error("ListFinding NewSelfClientCached", "error", err)
		return nil, err
	}
	integrationId, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyIntegrationID)
	if err != nil {
		plugin.Logger(ctx).Error("ListFinding GetConfigTableValueOrNil for OpenGovernanceConfigKeyIntegrationID", "error", err)
		return nil, err
	}
	encodedResourceCollectionFilters, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyResourceCollectionFilters)
	if err != nil {
		plugin.Logger(ctx).Error("ListFinding GetConfigTableValueOrNil for OpenGovernanceConfigKeyResourceCollectionFilters", "error", err)
		return nil, err
	}
	clientType, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyClientType)
	if err != nil {
		plugin.Logger(ctx).Error("ListFinding GetConfigTableValueOrNil for OpenGovernanceConfigKeyClientType", "error", err)
		return nil, err
	}

	paginator, err := k.NewFindingPaginator(essdk.BuildFilter(ctx, d.QueryContext, listFindingFilters, integrationId, encodedResourceCollectionFilters, clientType), d.QueryContext.Limit)
	if err != nil {
		plugin.Logger(ctx).Error("ListFinding NewFindingPaginator", "error", err)
		return nil, err
	}

	for paginator.HasNext() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("ListFinding paginator.NextPage", "error", err)
			return nil, err
		}

		for _, v := range page {
			d.StreamListItem(ctx, v)
		}
	}

	err = paginator.Close(ctx)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

var getFindingFilters = map[string]string{
	"assistant":          "Description.Assistant",
	"categories":         "Description.Categories",
	"confidence":         "Description.Confidence",
	"created_at":         "Description.CreatedAt",
	"external_ticket":    "Description.ExternalTicket",
	"first_seen_scan_id": "Description.FirstSeenScanID",
	"id":                 "Description.ID",
	"line_of_code_url":   "Description.LineOfCodeURL",
	"location":           "Description.Location",
	"match_based_id":     "Description.MatchBasedID",
	"ref":                "Description.Ref",
	"relevant_since":     "Description.RelevantSince",
	"repository":         "Description.Repository",
	"rule":               "Description.Rule",
	"rule_message":       "Description.RuleMessage",
	"rule_name":          "Description.RuleName",
	"severity":           "Description.Severity",
	"sourcing_policy":    "Description.SourcingPolicy",
	"state":              "Description.State",
	"status":             "Description.Status",
	"syntactic_id":       "Description.SyntacticID",
	"triage_state":       "Description.TriageState",
}

func GetFinding(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("GetFinding")
	runtime.GC()
	// create service
	cfg := essdk.GetConfig(d.Connection)
	ke, err := essdk.NewClientCached(cfg, d.ConnectionCache, ctx)
	if err != nil {
		return nil, err
	}
	k := Client{Client: ke}

	sc, err := steampipesdk.NewSelfClientCached(ctx, d.ConnectionCache)
	if err != nil {
		return nil, err
	}
	integrationId, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyIntegrationID)
	if err != nil {
		return nil, err
	}
	encodedResourceCollectionFilters, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyResourceCollectionFilters)
	if err != nil {
		return nil, err
	}
	clientType, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyClientType)
	if err != nil {
		return nil, err
	}

	limit := int64(1)
	paginator, err := k.NewFindingPaginator(essdk.BuildFilter(ctx, d.QueryContext, getFindingFilters, integrationId, encodedResourceCollectionFilters, clientType), &limit)
	if err != nil {
		return nil, err
	}

	for paginator.HasNext() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page {
			return v, nil
		}
	}

	err = paginator.Close(ctx)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ==========================  END: Finding =============================
