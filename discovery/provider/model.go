// Implement types for each resource

package provider

type Metadata struct{}

type DeploymentsResponse struct {
	Deployments []DeploymentJSON `json:"deployments"`
}

type FindingJSON struct {
	URL string `json:"url"`
}

type Finding struct {
	URL string
}

type DeploymentJSON struct {
	Slug     string      `json:"slug"`
	ID       int         `json:"id"`
	Name     string      `json:"name"`
	Findings FindingJSON `json:"findings"`
}

type DeploymentDescription struct {
	Slug     string
	ID       int
	Name     string
	Findings Finding
}

type ProjectsListResponse struct {
	Projects []ProjectJSON `json:"projects"`
}

type ProjectJSON struct {
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	URL           string   `json:"url"`
	Tags          []string `json:"tags"`
	CreatedAt     string   `json:"created_at"`
	LatestScanAt  string   `json:"latest_scan_at"`
	PrimaryBranch string   `json:"primary_branch"`
	DefaultBranch string   `json:"default_branch"`
}

type ProjectDescription struct {
	ID            int
	Name          string
	URL           string
	Tags          []string
	CreatedAt     string
	LatestScanAt  string
	PrimaryBranch string
	DefaultBranch string
}

type PoliciesListResponse struct {
	Policies []PolicyJSON `json:"policies"`
}

type PolicyJSON struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	ProductType string `json:"productType"`
	IsDefault   bool   `json:"isDefault"`
}

type PolicyDescription struct {
	ID          string
	Name        string
	Slug        string
	ProductType string
	IsDefault   bool
}

type ScansListResponse struct {
	Scans []ScanJSON `json:"scans"`
}

type FindingsCountJSON struct {
	Total       int `json:"total"`
	Code        int `json:"code"`
	SupplyChain int `json:"supply_chain"`
	Secrets     int `json:"secrets"`
}

type FindingsCount struct {
	Total       int
	Code        int
	SupplyChain int
	Secrets     int
}

type ScanJSON struct {
	ID             string        `json:"id"`
	DeploymentID   string        `json:"deployment_id"`
	RepositoryID   string        `json:"repository_id"`
	Branch         string        `json:"branch"`
	Commit         string        `json:"commit"`
	IsFullScan     bool          `json:"is_full_scan"`
	StartedAt      string        `json:"started_at"`
	CompletedAt    string        `json:"completed_at"`
	ExitCode       int           `json:"exit_code"`
	TotalTime      float64       `json:"total_time"`
	FindingsCounts FindingsCount `json:"findings_counts"`
	Status         string        `json:"status"`
}

type ScanDescription struct {
	ID             string
	DeploymentID   string
	RepositoryID   string
	Branch         string
	Commit         string
	IsFullScan     bool
	StartedAt      string
	CompletedAt    string
	ExitCode       int
	TotalTime      float64
	FindingsCounts FindingsCount
	Status         string
}

type FindingsListResponse struct {
	Findings []FindingObject `json:"findings"`
}

type ExternalTicketJSON struct {
	ExternalSlug string `json:"external_slug"`
	URL          string `json:"url"`
}

type ExternalTicket struct {
	ExternalSlug string
	URL          string
}

type RepositoryJSON struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Repository struct {
	Name string
	URL  string
}

type LocationJSON struct {
	FilePath  string `json:"file_path"`
	Line      int    `json:"line"`
	Column    int    `json:"column"`
	EndLine   int    `json:"end_line"`
	EndColumn int    `json:"end_column"`
}

type Location struct {
	FilePath  string
	Line      int
	Column    int
	EndLine   int
	EndColumn int
}

type SourcingPolicyJSON struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type SourcingPolicy struct {
	ID   int
	Name string
	Slug string
}

type RuleJSON struct {
	Name                 string   `json:"name"`
	Message              string   `json:"message"`
	Confidence           string   `json:"confidence"`
	Category             string   `json:"category"`
	Subcategories        []string `json:"subcategories"`
	VulnerabilityClasses []string `json:"vulnerability_classes"`
	CWENames             []string `json:"cwe_names"`
	OWASPNames           []string `json:"owasp_names"`
}

type Rule struct {
	Name                 string
	Message              string
	Confidence           string
	Category             string
	Subcategories        []string
	VulnerabilityClasses []string
	CWENames             []string
	OWASPNames           []string
}

type AssistantJSON struct {
	Autofix    Autofix    `json:"autofix"`
	Guidance   Guidance   `json:"guidance"`
	Autotriage Autotriage `json:"autotriage"`
	Component  Component  `json:"component"`
}

type Assistant struct {
	Autofix    Autofix
	Guidance   Guidance
	Autotriage Autotriage
	Component  Component
}

type AutofixJSON struct {
	FixCode     string `json:"fix_code"`
	Explanation string `json:"explanation"`
}

type Autofix struct {
	FixCode     string
	Explanation string
}

type GuidanceJSON struct {
	Summary      string `json:"summary"`
	Instructions string `json:"instructions"`
}

type Guidance struct {
	Summary      string
	Instructions string
}

type AutotriageJSON struct {
	Verdict string `json:"verdict"`
	Reason  string `json:"reason"`
}

type Autotriage struct {
	Verdict string
	Reason  string
}

type ComponentJSON struct {
	Tag  string `json:"tag"`
	Risk string `json:"risk"`
}

type Component struct {
	Tag  string
	Risk string
}

type FindingObject struct {
	ID              int                `json:"id"`
	Ref             string             `json:"ref"`
	FirstSeenScanID int                `json:"first_seen_scan_id"`
	SyntacticID     string             `json:"syntactic_id"`
	MatchBasedID    string             `json:"match_based_id"`
	ExternalTicket  ExternalTicketJSON `json:"external_ticket"`
	Repository      RepositoryJSON     `json:"repository"`
	LineOfCodeURL   string             `json:"line_of_code_url"`
	TriageState     string             `json:"triage_state"`
	State           string             `json:"state"`
	Status          string             `json:"status"`
	Severity        string             `json:"severity"`
	Confidence      string             `json:"confidence"`
	Categories      []string           `json:"categories"`
	CreatedAt       string             `json:"created_at"`
	RelevantSince   string             `json:"relevant_since"`
	RuleName        string             `json:"rule_name"`
	RuleMessage     string             `json:"rule_message"`
	Location        LocationJSON       `json:"location"`
	SourcingPolicy  SourcingPolicyJSON `json:"sourcing_policy"`
	TriagedAt       string             `json:"triaged_at"`
	TriageComment   string             `json:"triage_comment"`
	TriageReason    string             `json:"triage_reason"`
	StateUpdatedAt  string             `json:"state_updated_at"`
	Rule            RuleJSON           `json:"rule"`
	Assistant       AssistantJSON      `json:"assistant"`
}

type FindingDescription struct {
	ID              int
	Ref             string
	FirstSeenScanID int
	SyntacticID     string
	MatchBasedID    string
	ExternalTicket  ExternalTicket
	Repository      Repository
	LineOfCodeURL   string
	TriageState     string
	State           string
	Status          string
	Severity        string
	Confidence      string
	Categories      []string
	CreatedAt       string
	RelevantSince   string
	RuleName        string
	RuleMessage     string
	Location        Location
	SourcingPolicy  SourcingPolicy
	TriagedAt       string
	TriageComment   string
	TriageReason    string
	StateUpdatedAt  string
	Rule            Rule
	Assistant       Assistant
}
