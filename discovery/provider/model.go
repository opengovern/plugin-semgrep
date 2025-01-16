//go:generate go run ../../SDK/runnable/models/main.go --file $GOFILE --output ../../SDK/generated/resources_clients.go --type $PROVIDER

// Implement types for each resource

package provider

import (
	"encoding/json"
	"time"
	goPipeline "github.com/buildkite/go-pipeline"
	"github.com/google/go-github/v55/github"
	"github.com/shurcooL/githubv4"
	steampipemodels "github.com/turbot/steampipe-plugin-github/github/models"
)

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

type RunnerLabels struct {
	ID   *int64
	Name *string
	Type *string
}

type RunnerDescription struct {
	ID           *int64
	Name         *string
	OS           *string
	Status       *string
	Busy         *bool
	Labels       []*RunnerLabels
	RepoFullName *string
}

type SecretDescription struct {
	Name                    *string
	CreatedAt               *string
	UpdatedAt               *string
	Visibility              *string
	SelectedRepositoriesURL *string
	RepoFullName            *string
}

type SimpleActor struct {
	Login  *string
	ID     int
	NodeID *string
	Type   *string
}

type SimpleRepo struct {
	ID     int
	NodeID *string
}

type CommitRefWorkflow struct {
	ID *string
}

type WorkflowRunDescription struct {
	ID                  int
	Name                *string
	HeadBranch          *string
	HeadSHA             *string
	Status              *string
	Conclusion          *string
	HTMLURL             *string
	WorkflowID          int
	RunNumber           int
	Event               *string
	CreatedAt           *string
	UpdatedAt           *string
	RunAttempt          int
	RunStartedAt        *string
	Actor               *SimpleActor
	HeadCommit          *CommitRefWorkflow
	Repository          *SimpleRepo
	HeadRepository      *SimpleRepo
	ReferencedWorkflows []interface{}
	ArtifactCount       int
	Artifacts           []WorkflowArtifact
}

type WorkflowRunsResponse struct {
	TotalCount   int                      `json:"total_count"`
	WorkflowRuns []WorkflowRunDescription `json:"workflow_runs"`
}

type WorkflowArtifactJSON struct {
	ID                 int    `json:"id"`
	NodeID             string `json:"node_id"`
	Name               string `json:"name"`
	SizeInBytes        int    `json:"size_in_bytes"`
	URL                string `json:"url"`
	ArchiveDownloadURL string `json:"archive_download_url"`
	Expired            bool   `json:"expired"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
	ExpiresAt          string `json:"expires_at"`
}

type WorkflowArtifact struct {
	ID                 int
	NodeID             *string
	Name               *string
	SizeInBytes        int
	URL                *string
	ArchiveDownloadURL *string
	Expired            bool
	CreatedAt          *string
	UpdatedAt          *string
	ExpiresAt          *string
}

type ArtifactsResponse struct {
	TotalCount int                    `json:"total_count"`
	Artifacts  []WorkflowArtifactJSON `json:"artifacts"`
}

type ActorLocation struct {
	CountryCode *string
}

type AuditEntryData struct {
	OldName  *string
	OldLogin *string
}

type AuditLogDescription struct {
	ID            *string
	CreatedAt     *string
	Organization  *string
	Phrase        *string
	Include       *string
	Action        *string
	Actor         *string
	ActorLocation *ActorLocation
	Team          *string
	UserLogin     *string
	Repo          *string
	Data          *AuditEntryData
}

//type BlobDescription struct {
//	Content      *string
//	Encoding     *string
//	SHA          *string
//	Size         *int
//	URL          *string
//	NodeID       *string
//	RepoFullName string
//}

type BasicUser struct {
	Id        int
	NodeId    *string
	Name      *string
	Login     *string
	Email     *string
	CreatedAt *string
	UpdatedAt *string
	Url       *string
}

type GitActor struct {
	AvatarUrl *string
	Date      *string
	Email     *string
	Name      *string
	User      BasicUser
}

type Signature struct {
	Email             *string
	IsValid           bool
	State             *string
	WasSignedByGitHub bool
	Signer            struct {
		Email *string
		Login *string
	}
}

type CommitStatus struct {
	State *string
}

type BaseCommit struct {
	Sha                 *string
	ShortSha            *string
	AuthoredDate        *string
	Author              GitActor
	CommittedDate       *string
	Committer           GitActor
	Message             *string
	Url                 *string
	Additions           int
	AuthoredByCommitter bool
	ChangedFiles        int
	CommittedViaWeb     bool
	CommitUrl           *string
	Deletions           int
	Signature           Signature
	TarballUrl          *string
	TreeUrl             *string
	CanSubscribe        bool
	Subscription        *string
	ZipballUrl          *string
	MessageHeadline     *string
	Status              CommitStatus
	NodeId              *string
}

type Actor struct {
	AvatarUrl *string
	Login     *string
	Url       *string
}

type BranchProtectionRule struct {
	AllowsDeletions                bool
	AllowsForcePushes              bool
	BlocksCreations                bool
	CreatorLogin                   *string
	Id                             int
	NodeId                         *string
	DismissesStaleReviews          bool
	IsAdminEnforced                bool
	LockAllowsFetchAndMerge        bool
	LockBranch                     bool
	Pattern                        *string
	RequireLastPushApproval        bool
	RequiredApprovingReviewCount   int
	RequiredDeploymentEnvironments []string
	RequiredStatusChecks           []string
	RequiresApprovingReviews       bool
	RequiresConversationResolution bool
	RequiresCodeOwnerReviews       bool
	RequiresCommitSignatures       bool
	RequiresDeployments            bool
	RequiresLinearHistory          bool
	RequiresStatusChecks           bool
	RequiresStrictStatusChecks     bool
	RestrictsPushes                bool
	RestrictsReviewDismissals      bool
	MatchingBranches               int
}

type BranchDescription struct {
	RepoFullName         *string
	Name                 *string
	Commit               BaseCommit
	BranchProtectionRule BranchProtectionRule
	Protected            bool
}

type BranchApp struct {
	Name *string
	Slug *string
}

type BranchTeam struct {
	Name *string
	Slug *string
}

type BranchUser struct {
	Name  *string
	Login *string
}

type BranchProtectionDescription struct {
	AllowsDeletions                 bool
	AllowsForcePushes               bool
	BlocksCreations                 bool
	Id                              int
	NodeId                          *string
	DismissesStaleReviews           bool
	IsAdminEnforced                 bool
	LockAllowsFetchAndMerge         bool
	LockBranch                      bool
	Pattern                         *string
	RequireLastPushApproval         bool
	RequiredApprovingReviewCount    int
	RequiredDeploymentEnvironments  []string
	RequiredStatusChecks            []string
	RequiresApprovingReviews        bool
	RequiresConversationResolution  bool
	RequiresCodeOwnerReviews        bool
	RequiresCommitSignatures        bool
	RequiresDeployments             bool
	RequiresLinearHistory           bool
	RequiresStatusChecks            bool
	RequiresStrictStatusChecks      bool
	RestrictsPushes                 bool
	RestrictsReviewDismissals       bool
	RepoFullName                    *string
	CreatorLogin                    *string
	MatchingBranches                int
	PushAllowanceApps               []BranchApp
	PushAllowanceTeams              []BranchTeam
	PushAllowanceUsers              []BranchUser
	BypassForcePushAllowanceApps    []BranchApp
	BypassForcePushAllowanceTeams   []BranchTeam
	BypassForcePushAllowanceUsers   []BranchUser
	BypassPullRequestAllowanceApps  []BranchApp
	BypassPullRequestAllowanceTeams []BranchTeam
	BypassPullRequestAllowanceUsers []BranchUser
}

type TreeJSON struct {
	SHA string `json:"sha"`
	URL string `json:"url"`
}

type Tree struct {
	SHA *string
	URL *string
}

type FileJSON struct {
	SHA         string  `json:"sha"`
	Filename    string  `json:"filename"`
	Status      string  `json:"status"`
	Additions   int     `json:"additions"`
	Deletions   int     `json:"deletions"`
	Changes     int     `json:"changes"`
	BlobURL     string  `json:"blob_url"`
	RawURL      string  `json:"raw_url"`
	ContentsURL string  `json:"contents_url"`
	Patch       *string `json:"patch"`
}

type File struct {
	SHA         *string
	Filename    *string
	Status      *string
	Additions   int
	Deletions   int
	Changes     int
	BlobURL     *string
	RawURL      *string
	ContentsURL *string
	Patch       *string
}

type VerificationJSON struct {
	Verified   bool    `json:"verified"`
	Reason     string  `json:"reason"`
	Signature  *string `json:"signature"`
	Payload    *string `json:"payload"`
	VerifiedAt *string `json:"verified_at"`
}

type Verification struct {
	Verified   bool
	Reason     *string
	Signature  *string
	Payload    *string
	VerifiedAt *string
}

type UserJSON struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	UserViewType      string `json:"user_view_type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type User struct {
	Login             *string `json:"login"`
	ID                int     `json:"id"`
	NodeID            *string `json:"node_id"`
	AvatarURL         *string `json:"avatar_url"`
	GravatarID        *string `json:"gravatar_id"`
	URL               *string `json:"url"`
	HTMLURL           *string `json:"html_url"`
	FollowersURL      *string `json:"followers_url"`
	FollowingURL      *string `json:"following_url"`
	GistsURL          *string `json:"gists_url"`
	StarredURL        *string `json:"starred_url"`
	SubscriptionsURL  *string `json:"subscriptions_url"`
	OrganizationsURL  *string `json:"organizations_url"`
	ReposURL          *string `json:"repos_url"`
	EventsURL         *string `json:"events_url"`
	ReceivedEventsURL *string `json:"received_events_url"`
	Type              *string `json:"type"`
	UserViewType      *string `json:"user_view_type"`
	SiteAdmin         bool    `json:"site_admin"`
}

type CommitDetailJSON struct {
	//Author       UserMinimalInfo `json:"author"`
	//Committer    UserMinimalInfo `json:"committer"`
	//URL          string          `json:"url"`
	Message      string           `json:"message"`
	Tree         TreeJSON         `json:"tree"`
	CommentCount int              `json:"comment_count"`
	Verification VerificationJSON `json:"verification"`
}

type CommitDetail struct {
	Message      *string
	Tree         Tree
	CommentCount int
	Verification Verification
}

type ParentJSON struct {
	SHA     string `json:"sha"`
	URL     string `json:"url"`
	HTMLURL string `json:"html_url"`
}

type Parent struct {
	SHA     *string
	URL     *string
	HTMLURL *string
}

type StatsJSON struct {
	Total     int `json:"total"`
	Additions int `json:"additions"`
	Deletions int `json:"deletions"`
}

type Stats struct {
	Total     int
	Additions int
	Deletions int
}

type CommitResp struct {
	SHA          string           `json:"sha"`
	NodeID       string           `json:"node_id"`
	CommitDetail CommitDetailJSON `json:"commit"`
	URL          string           `json:"url"`
	HTMLURL      string           `json:"html_url"`
	CommentsURL  string           `json:"comments_url"`
	Author       UserJSON         `json:"author"`
	Committer    UserJSON         `json:"committer"`
	Parents      []ParentJSON     `json:"parents"`
	Stats        StatsJSON        `json:"stats"`
	Files        []FileJSON       `json:"files"`
}

type CommitDescription struct {
	SHA          *string
	NodeID       *string
	CommitDetail CommitDetail
	URL          *string
	HTMLURL      *string
	CommentsURL  *string
	Author       User
	Committer    User
	Parents      []Parent
	Stats        Stats
	Files        []File
}

type Milestone struct {
	Closed             bool
	ClosedAt           *string
	CreatedAt          *string
	Creator            Actor
	Description        *string
	DueOn              *string
	Number             int
	ProgressPercentage float32
	State              *githubv4.MilestoneState
	Title              *string
	UpdatedAt          *string
	UserCanClose       bool
	UserCanReopen      bool
}

type Label struct {
	NodeId      *string
	Name        *string
	Description *string
	IsDefault   bool
	Color       *string
}

type RepositoryInteractionAbility struct {
	ExpiresAt *string
	Limit     *string
	Origin    *string
}

type SponsorsGoal struct {
	Description     *string
	PercentComplete int
	TargetValue     int
	Title           *string
	Kind            *githubv4.SponsorsGoalKind
}

type StripeConnectAccount struct {
	AccountId              *string
	BillingCountryOrRegion *string
	CountryOrRegion        *string
	IsActive               bool
	StripeDashboardUrl     *string
}

type SponsorsListing struct {
	Id                         *string
	ActiveGoal                 SponsorsGoal
	ActiveStripeConnectAccount StripeConnectAccount
	BillingCountryOrRegion     *string
	ContactEmailAddress        *string
	CreatedAt                  *string
	DashboardUrl               *string
	FullDescription            *string
	IsPublic                   bool
	Name                       *string
	NextPayoutDate             *string
	ResidenceCountryOrRegion   *string
	ShortDescription           *string
	Slug                       *string
	Url                        *string
}

type BaseUser struct {
	BasicUser
	AnyPinnableItems                      bool
	AvatarUrl                             *string
	Bio                                   *string
	Company                               *string
	EstimatedNextSponsorsPayoutInCents    int
	HasSponsorsListing                    bool
	InteractionAbility                    RepositoryInteractionAbility
	IsBountyHunter                        bool
	IsCampusExpert                        bool
	IsDeveloperProgramMember              bool
	IsEmployee                            bool
	IsFollowingYou                        bool
	IsGitHubStar                          bool
	IsHireable                            bool
	IsSiteAdmin                           bool
	IsSponsoringYou                       bool
	IsYou                                 bool
	Location                              *string
	MonthlyEstimatedSponsorsIncomeInCents int
	PinnedItemsRemaining                  int
	ProjectsUrl                           *string
	Pronouns                              *string
	SponsorsListing                       SponsorsListing
	Status                                UserStatus
	TwitterUsername                       *string
	CanChangedPinnedItems                 bool
	CanCreateProjects                     bool
	CanFollow                             bool
	CanSponsor                            bool
	IsFollowing                           bool
	IsSponsoring                          bool
	WebsiteUrl                            *string
}

type UserStatus struct {
	CreatedAt                    *string
	UpdatedAt                    *string
	ExpiresAt                    *string
	Emoji                        *string
	Message                      *string
	IndicatesLimitedAvailability bool
}

type IssueDescription struct {
	RepositoryFullName      *string
	Id                      int
	NodeId                  *string
	Number                  int
	ActiveLockReason        *githubv4.LockReason
	Author                  Actor
	AuthorLogin             *string
	AuthorAssociation       *githubv4.CommentAuthorAssociation
	Body                    *string
	BodyUrl                 *string
	Closed                  bool
	ClosedAt                *string
	CreatedAt               *string
	CreatedViaEmail         bool
	Editor                  Actor
	FullDatabaseId          *string
	IncludesCreatedEdit     bool
	IsPinned                bool
	IsReadByUser            bool
	LastEditedAt            *string
	Locked                  bool
	Milestone               Milestone
	PublishedAt             *string
	State                   *githubv4.IssueState
	StateReason             *githubv4.IssueStateReason
	Title                   *string
	UpdatedAt               *string
	Url                     *string
	UserCanClose            bool
	UserCanReact            bool
	UserCanReopen           bool
	UserCanSubscribe        bool
	UserCanUpdate           bool
	UserCannotUpdateReasons []githubv4.CommentCannotUpdateReason
	UserDidAuthor           bool
	UserSubscription        *githubv4.SubscriptionState
	CommentsTotalCount      int
	LabelsTotalCount        int
	LabelsSrc               []Label
	Labels                  map[string]Label
	AssigneesTotalCount     int
	Assignees               []BaseUser
}

//type IssueCommentDescription struct {
//	steampipemodels.IssueComment
//	RepoFullName string
//	Number       int
//	AuthorLogin  string
//	EditorLogin  string
//}

type LicenseDescription struct {
	Key            string
	Name           string
	Nickname       string
	SpdxId         string
	Url            string
	Body           string
	Conditions     []steampipemodels.LicenseRule
	Description    string
	Featured       bool
	Hidden         bool
	Implementation string
	Limitations    []steampipemodels.LicenseRule
	Permissions    []steampipemodels.LicenseRule
	PseudoLicense  bool
}

type OrganizationDescription struct {
	Id                                     int
	NodeId                                 string
	Name                                   string
	Login                                  string
	CreatedAt                              string
	UpdatedAt                              string
	Description                            string
	Email                                  string
	Url                                    string
	Announcement                           string
	AnnouncementExpiresAt                  string
	AnnouncementUserDismissible            bool
	AnyPinnableItems                       bool
	AvatarUrl                              string
	EstimatedNextSponsorsPayoutInCents     int
	HasSponsorsListing                     bool
	InteractionAbility                     steampipemodels.RepositoryInteractionAbility
	IsSponsoringYou                        bool
	IsVerified                             bool
	Location                               string
	MonthlyEstimatedSponsorsIncomeInCents  int
	NewTeamUrl                             string
	PinnedItemsRemaining                   int
	ProjectsUrl                            string
	SamlIdentityProvider                   steampipemodels.OrganizationIdentityProvider
	SponsorsListing                        steampipemodels.SponsorsListing
	TeamsUrl                               string
	TotalSponsorshipAmountAsSponsorInCents int
	TwitterUsername                        string
	CanAdminister                          bool
	CanChangedPinnedItems                  bool
	CanCreateProjects                      bool
	CanCreateRepositories                  bool
	CanCreateTeams                         bool
	CanSponsor                             bool
	IsAMember                              bool
	IsFollowing                            bool
	IsSponsoring                           bool
	WebsiteUrl                             string
	Hooks                                  []*github.Hook
	BillingEmail                           string
	TwoFactorRequirementEnabled            bool
	DefaultRepoPermission                  string
	MembersAllowedRepositoryCreationType   string
	MembersCanCreateInternalRepos          bool
	MembersCanCreatePages                  bool
	MembersCanCreatePrivateRepos           bool
	MembersCanCreatePublicRepos            bool
	MembersCanCreateRepos                  bool
	MembersCanForkPrivateRepos             bool
	PlanFilledSeats                        int
	PlanName                               string
	PlanPrivateRepos                       int
	PlanSeats                              int
	PlanSpace                              int
	Followers                              int
	Following                              int
	Collaborators                          int
	HasOrganizationProjects                bool
	HasRepositoryProjects                  bool
	WebCommitSignoffRequired               bool
	MembersWithRoleTotalCount              int
	PackagesTotalCount                     int
	PinnableItemsTotalCount                int
	PinnedItemsTotalCount                  int
	ProjectsTotalCount                     int
	ProjectsV2TotalCount                   int
	SponsoringTotalCount                   int
	SponsorsTotalCount                     int
	TeamsTotalCount                        int
	PrivateRepositoriesTotalCount          int
	PublicRepositoriesTotalCount           int
	RepositoriesTotalCount                 int
	RepositoriesTotalDiskUsage             int
}

type OrgCollaboratorsDescription struct {
	Organization   string
	Affiliation    string
	RepositoryName githubv4.String
	Permission     githubv4.RepositoryPermission
	UserLogin      steampipemodels.CollaboratorLogin
}

type OrgAlertDependabotDescription struct {
	AlertNumber                 int
	State                       string
	DependencyPackageEcosystem  string
	DependencyPackageName       string
	DependencyManifestPath      string
	DependencyScope             string
	SecurityAdvisoryGHSAID      string
	SecurityAdvisoryCVEID       string
	SecurityAdvisorySummary     string
	SecurityAdvisoryDescription string
	SecurityAdvisorySeverity    string
	SecurityAdvisoryCVSSScore   *float64
	SecurityAdvisoryCVSSVector  string
	SecurityAdvisoryCWEs        []string
	SecurityAdvisoryPublishedAt github.Timestamp
	SecurityAdvisoryUpdatedAt   github.Timestamp
	SecurityAdvisoryWithdrawnAt github.Timestamp
	URL                         string
	HTMLURL                     string
	CreatedAt                   github.Timestamp
	UpdatedAt                   github.Timestamp
	DismissedAt                 github.Timestamp
	DismissedReason             string
	DismissedComment            string
	FixedAt                     github.Timestamp
}

type OrgExternalIdentityDescription struct {
	steampipemodels.OrganizationExternalIdentity
	Organization string
	UserLogin    string
	UserDetail   steampipemodels.BasicUser
}

type OrgMembersDescription struct {
	steampipemodels.User
	Organization        string
	HasTwoFactorEnabled *bool
	Role                *string
}

type PullRequestDescription struct {
	RepoFullName             string
	Id                       int
	NodeId                   string
	Number                   int
	ActiveLockReason         githubv4.LockReason
	Additions                int
	Author                   steampipemodels.Actor
	AuthorAssociation        githubv4.CommentAuthorAssociation
	BaseRefName              string
	Body                     string
	ChangedFiles             int
	ChecksUrl                string
	Closed                   bool
	ClosedAt                 steampipemodels.NullableTime
	CreatedAt                steampipemodels.NullableTime
	CreatedViaEmail          bool
	Deletions                int
	Editor                   steampipemodels.Actor
	HeadRefName              string
	HeadRefOid               string
	IncludesCreatedEdit      bool
	IsCrossRepository        bool
	IsDraft                  bool
	IsReadByUser             bool
	LastEditedAt             steampipemodels.NullableTime
	Locked                   bool
	MaintainerCanModify      bool
	Mergeable                githubv4.MergeableState
	Merged                   bool
	MergedAt                 steampipemodels.NullableTime
	MergedBy                 steampipemodels.Actor
	Milestone                steampipemodels.Milestone
	Permalink                string
	PublishedAt              steampipemodels.NullableTime
	RevertUrl                string
	ReviewDecision           githubv4.PullRequestReviewDecision
	State                    githubv4.PullRequestState
	Title                    string
	TotalCommentsCount       int
	UpdatedAt                steampipemodels.NullableTime
	Url                      string
	Assignees                []steampipemodels.BaseUser
	BaseRef                  *steampipemodels.BasicRef
	HeadRef                  *steampipemodels.BasicRef
	MergeCommit              *steampipemodels.BasicCommit
	SuggestedReviewers       []steampipemodels.SuggestedReviewer
	CanApplySuggestion       bool
	CanClose                 bool
	CanDeleteHeadRef         bool
	CanDisableAutoMerge      bool
	CanEditFiles             bool
	CanEnableAutoMerge       bool
	CanMergeAsAdmin          bool
	CanReact                 bool
	CanReopen                bool
	CanSubscribe             bool
	CanUpdate                bool
	CanUpdateBranch          bool
	DidAuthor                bool
	CannotUpdateReasons      []githubv4.CommentCannotUpdateReason
	Subscription             githubv4.SubscriptionState
	LabelsSrc                []steampipemodels.Label
	Labels                   map[string]steampipemodels.Label
	CommitsTotalCount        int
	ReviewRequestsTotalCount int
	ReviewsTotalCount        int
	LabelsTotalCount         int
	AssigneesTotalCount      int
}

type License struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	SPDXID string `json:"spdx_id"`
	URL    string `json:"url"`
	NodeID string `json:"node_id"`
}

type Permissions struct {
	Admin    bool `json:"admin"`
	Maintain bool `json:"maintain"`
	Push     bool `json:"push"`
	Triage   bool `json:"triage"`
	Pull     bool `json:"pull"`
}

type StatusObj struct {
	Status string `json:"status"`
}

type RepoDetail struct {
	ID              int     `json:"id"`
	NodeID          string  `json:"node_id"`
	Name            string  `json:"name"`
	FullName        string  `json:"full_name"`
	Private         bool    `json:"private"`
	Owner           *Owner  `json:"owner"`
	HTMLURL         string  `json:"html_url"`
	Description     *string `json:"description"`
	Fork            bool    `json:"fork"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
	PushedAt        string  `json:"pushed_at"`
	GitURL          string  `json:"git_url"`
	SSHURL          string  `json:"ssh_url"`
	CloneURL        string  `json:"clone_url"`
	SVNURL          string  `json:"svn_url"`
	Homepage        *string `json:"homepage"`
	Size            int     `json:"size"`
	StargazersCount int     `json:"stargazers_count"`
	WatchersCount   int     `json:"watchers_count"`

	// 1) The single “primary” language returned by the main /repos/:owner/:repo call
	PrimaryLanguage *string `json:"language"`
	// If you want to store the breakdown from /languages in the same struct, you can do:
	LanguageBreakdown map[string]int `json:"-"`

	HasIssues                 bool                   `json:"has_issues"`
	HasProjects               bool                   `json:"has_projects"`
	HasDownloads              bool                   `json:"has_downloads"`
	HasWiki                   bool                   `json:"has_wiki"`
	HasPages                  bool                   `json:"has_pages"`
	HasDiscussions            bool                   `json:"has_discussions"`
	ForksCount                int                    `json:"forks_count"`
	MirrorURL                 *string                `json:"mirror_url"`
	Archived                  bool                   `json:"archived"`
	Disabled                  bool                   `json:"disabled"`
	OpenIssuesCount           int                    `json:"open_issues_count"`
	License                   *License               `json:"license"`
	AllowForking              bool                   `json:"allow_forking"`
	IsTemplate                bool                   `json:"is_template"`
	WebCommitSignoffRequired  bool                   `json:"web_commit_signoff_required"`
	Topics                    []string               `json:"topics"`
	Visibility                string                 `json:"visibility"`
	DefaultBranch             string                 `json:"default_branch"`
	Permissions               *Permissions           `json:"permissions"`
	AllowSquashMerge          bool                   `json:"allow_squash_merge"`
	AllowMergeCommit          bool                   `json:"allow_merge_commit"`
	AllowRebaseMerge          bool                   `json:"allow_rebase_merge"`
	AllowAutoMerge            bool                   `json:"allow_auto_merge"`
	DeleteBranchOnMerge       bool                   `json:"delete_branch_on_merge"`
	AllowUpdateBranch         bool                   `json:"allow_update_branch"`
	UseSquashPRTitleAsDefault bool                   `json:"use_squash_pr_title_as_default"`
	SquashMergeCommitMessage  string                 `json:"squash_merge_commit_message"`
	SquashMergeCommitTitle    string                 `json:"squash_merge_commit_title"`
	MergeCommitMessage        string                 `json:"merge_commit_message"`
	MergeCommitTitle          string                 `json:"merge_commit_title"`
	CustomProperties          map[string]interface{} `json:"custom_properties"`
	Organization              *Organization          `json:"organization"`
	Parent                    *RepoDetail            `json:"parent"`
	Source                    *RepoDetail            `json:"source"`
	NetworkCount              int                    `json:"network_count"`
	SubscribersCount          int                    `json:"subscribers_count"`
	BlankIssuesEnabled        bool                   `json:"blank_issues_enabled"`
	Locked                    bool                   `json:"locked"`

	SecurityAndAnalysis *struct {
		SecretScanning                    *StatusObj `json:"secret_scanning"`
		SecretScanningPushProtection      *StatusObj `json:"secret_scanning_push_protection"`
		DependabotSecurityUpdates         *StatusObj `json:"dependabot_security_updates"`
		SecretScanningNonProviderPatterns *StatusObj `json:"secret_scanning_non_provider_patterns"`
		SecretScanningValidityChecks      *StatusObj `json:"secret_scanning_validity_checks"`
	} `json:"security_and_analysis"`
}

type RepositorySettings struct {
	HasDiscussionsEnabled     bool                   `json:"has_discussions_enabled"`
	HasIssuesEnabled          bool                   `json:"has_issues_enabled"`
	HasProjectsEnabled        bool                   `json:"has_projects_enabled"`
	HasWikiEnabled            bool                   `json:"has_wiki_enabled"`
	MergeCommitAllowed        bool                   `json:"merge_commit_allowed"`
	MergeCommitMessage        string                 `json:"merge_commit_message"`
	MergeCommitTitle          string                 `json:"merge_commit_title"`
	SquashMergeAllowed        bool                   `json:"squash_merge_allowed"`
	SquashMergeCommitMessage  string                 `json:"squash_merge_commit_message"`
	SquashMergeCommitTitle    string                 `json:"squash_merge_commit_title"`
	HasDownloads              bool                   `json:"has_downloads"`
	HasPages                  bool                   `json:"has_pages"`
	WebCommitSignoffRequired  bool                   `json:"web_commit_signoff_required"`
	MirrorURL                 *string                `json:"mirror_url"`
	AllowAutoMerge            bool                   `json:"allow_auto_merge"`
	DeleteBranchOnMerge       bool                   `json:"delete_branch_on_merge"`
	AllowUpdateBranch         bool                   `json:"allow_update_branch"`
	UseSquashPRTitleAsDefault bool                   `json:"use_squash_pr_title_as_default"`
	CustomProperties          map[string]interface{} `json:"custom_properties"`
	ForkingAllowed            bool                   `json:"forking_allowed"`
	IsTemplate                bool                   `json:"is_template"`
	AllowRebaseMerge          bool                   `json:"allow_rebase_merge"`
	// Renamed fields:
	Archived bool `json:"archived"`
	Disabled bool `json:"disabled"`
	Locked   bool `json:"locked"`
}

type SecuritySettings struct {
	VulnerabilityAlertsEnabled               bool `json:"vulnerability_alerts_enabled"`
	SecretScanningEnabled                    bool `json:"secret_scanning_enabled"`
	SecretScanningPushProtectionEnabled      bool `json:"secret_scanning_push_protection_enabled"`
	DependabotSecurityUpdatesEnabled         bool `json:"dependabot_security_updates_enabled"`
	SecretScanningNonProviderPatternsEnabled bool `json:"secret_scanning_non_provider_patterns_enabled"`
	SecretScanningValidityChecksEnabled      bool `json:"secret_scanning_validity_checks_enabled"`

	// New field
	PrivateVulnerabilityReportingEnabled bool `json:"private_vulnerability_reporting_enabled"`
}

type RepoURLs struct {
	GitURL   string `json:"git_url"`
	SSHURL   string `json:"ssh_url"`
	CloneURL string `json:"clone_url"`
	SVNURL   string `json:"svn_url"`
	HTMLURL  string `json:"html_url"`
}
type Owner struct {
	Login   string `json:"login"`
	ID      int    `json:"id"`
	NodeID  string `json:"node_id"`
	HTMLURL string `json:"html_url"`
	Type    string `json:"type"`
}

type Organization struct {
	Login        string `json:"login"`
	ID           int    `json:"id"`
	NodeID       string `json:"node_id"`
	HTMLURL      string `json:"html_url"`
	Type         string `json:"type"`
	UserViewType string `json:"user_view_type"`
	SiteAdmin    bool   `json:"site_admin"`
}

type Metrics struct {
	Stargazers   int `json:"stargazers"`
	Forks        int `json:"forks"`
	Subscribers  int `json:"subscribers"`
	Size         int `json:"size"`
	Tags         int `json:"tags"`
	Commits      int `json:"commits"`
	Issues       int `json:"issues"`
	OpenIssues   int `json:"open_issues"`
	Branches     int `json:"branches"`
	PullRequests int `json:"pull_requests"`
	Releases     int `json:"releases"`
}

//type RepositoryResponse struct {
//	ID                       float64                      `json:"id"`
//	NodeID                   string                       `json:"node_id"`
//	Name                     string                       `json:"name"`
//	NameWithOwner            string                       `json:"name_with_owner"`
//	Description              *string                      `json:"description"`
//	Private                  bool                         `json:"private"`
//	HTMLURL                  string                       `json:"html_url"`
//	Fork                     bool                         `json:"fork"`
//	URL                      string                       `json:"url"`
//	CloneURL                 string                       `json:"clone_url"`
//	GitURL                   string                       `json:"git_url"`
//	SSHURL                   string                       `json:"ssh_url"`
//	Homepage                 *string                      `json:"homepage"`
//	Language                 string                       `json:"language"`
//	ForksCount               int                          `json:"forks_count"`
//	StargazersCount          int                          `json:"stargazers_count"`
//	WatchersCount            int                          `json:"watchers_count"`
//	Size                     int                          `json:"size"`
//	DefaultBranch            string                       `json:"default_branch"`
//	OpenIssuesCount          int                          `json:"open_issues_count"`
//	Topics                   []string                     `json:"topics"`
//	HasIssues                bool                         `json:"has_issues"`
//	HasProjects              bool                         `json:"has_projects"`
//	HasWiki                  bool                         `json:"has_wiki"`
//	HasPages                 bool                         `json:"has_pages"`
//	HasDownloads             bool                         `json:"has_downloads"`
//	HasDiscussions           bool                         `json:"has_discussions"`
//	Archived                 bool                         `json:"archived"`
//	Disabled                 bool                         `json:"disabled"`
//	Visibility               string                       `json:"visibility"`
//	PushedAt                 string                       `json:"pushed_at"`
//	CreatedAt                string                       `json:"created_at"`
//	UpdatedAt                string                       `json:"updated_at"`
//	AllowForking             bool                         `json:"allow_forking"`
//	WebCommitSignoffRequired bool                         `json:"web_commit_signoff_required"`
//	ContentsURL              string                       `json:"contents_url"`
//	PullsURL                 string                       `json:"pulls_url"`
//	CommitsURL               string                       `json:"commits_url"`
//	CompareURL               string                       `json:"compare_url"`
//	DownloadsURL             string                       `json:"downloads_url"`
//	ArchiveURL               string                       `json:"archive_url"`
//	AssigneesURL             string                       `json:"assignees_url"`
//	BlobsURL                 string                       `json:"blobs_url"`
//	BranchesURL              string                       `json:"branches_url"`
//	CollaboratorsURL         string                       `json:"collaborators_url"`
//	CommentsURL              string                       `json:"comments_url"`
//	ContributorsURL          string                       `json:"contributors_url"`
//	DeploymentsURL           string                       `json:"deployments_url"`
//	EventsURL                string                       `json:"events_url"`
//	ForksURL                 string                       `json:"forks_url"`
//	GitCommitsURL            string                       `json:"git_commits_url"`
//	GitRefsURL               string                       `json:"git_refs_url"`
//	GitTagsURL               string                       `json:"git_tags_url"`
//	HooksURL                 string                       `json:"hooks_url"`
//	IssueCommentURL          string                       `json:"issue_comment_url"`
//	IssueEventsURL           string                       `json:"issue_events_url"`
//	IssuesURL                string                       `json:"issues_url"`
//	KeysURL                  string                       `json:"keys_url"`
//	LabelsURL                string                       `json:"labels_url"`
//	LanguagesURL             string                       `json:"languages_url"`
//	MergesURL                string                       `json:"merges_url"`
//	MilestonesURL            string                       `json:"milestones_url"`
//	NotificationsURL         string                       `json:"notifications_url"`
//	ReleasesURL              string                       `json:"releases_url"`
//	StargazersURL            string                       `json:"stargazers_url"`
//	StatusesURL              string                       `json:"statuses_url"`
//	SubscribersURL           string                       `json:"subscribers_url"`
//	SubscriptionURL          string                       `json:"subscription_url"`
//	TagsURL                  string                       `json:"tags_url"`
//	TeamsURL                 string                       `json:"teams_url"`
//	TreesURL                 string                       `json:"trees_url"`
//	Watchers                 int                          `json:"watchers"`
//	IsTemplate               bool                         `json:"is_template"`
//	SecurityAndAnalysis      map[string]map[string]string `json:"security_and_analysis"`
//	Permissions              map[string]bool              `json:"permissions"`
//	Owner                    RepoOwner                    `json:"owner"`
//}
//
//type RepoOwner struct {
//	Login             string `json:"login"`
//	ID                int64  `json:"id"`
//	NodeID            string `json:"node_id"`
//	AvatarURL         string `json:"avatar_url"`
//	GravatarID        string `json:"gravatar_id"`
//	URL               string `json:"url"`
//	HTMLURL           string `json:"html_url"`
//	FollowersURL      string `json:"followers_url"`
//	FollowingURL      string `json:"following_url"`
//	GistsURL          string `json:"gists_url"`
//	StarredURL        string `json:"starred_url"`
//	SubscriptionsURL  string `json:"subscriptions_url"`
//	OrganizationsURL  string `json:"organizations_url"`
//	ReposURL          string `json:"repos_url"`
//	EventsURL         string `json:"events_url"`
//	ReceivedEventsURL string `json:"received_events_url"`
//	Type              string `json:"type"`
//	SiteAdmin         bool   `json:"site_admin"`
//}
//
//type LicenseInfo struct {
//	Key    string `json:"key"`
//	Name   string `json:"name"`
//	SPDXID string `json:"spdx_id"`
//	URL    string `json:"url"`
//	NodeID string `json:"node_id"`
//}
//
//type DefaultBranchRef struct {
//	Name string `json:"name"`
//}

type RepositoryDescription struct {
	GitHubRepoID            int
	NodeID                  *string
	Name                    *string
	NameWithOwner           *string
	Description             *string
	CreatedAt               *string
	UpdatedAt               *string
	PushedAt                *string
	IsActive                bool
	IsEmpty                 bool
	IsFork                  bool
	IsSecurityPolicyEnabled bool
	Owner                   *Owner
	HomepageURL             *string
	LicenseInfo             json.RawMessage
	Topics                  []string
	Visibility              string
	DefaultBranchRef        json.RawMessage
	Permissions             *Permissions
	Organization            *Organization
	Parent                  *RepositoryDescription
	Source                  *RepositoryDescription
	PrimaryLanguage         *string
	Languages               map[string]int
	RepositorySettings      RepositorySettings
	SecuritySettings        SecuritySettings
	RepoURLs                RepoURLs
	Metrics                 Metrics
}

type MinimalRepoInfo struct {
	Name     string `json:"name"`
	Archived bool   `json:"archived"`
	Disabled bool   `json:"disabled"`
	Owner    struct {
		Login string `json:"login"`
	} `json:"owner"`
}

type ReleaseDescription struct {
	github.RepositoryRelease
	RepositoryFullName string
}

type RepoCollaboratorsDescription struct {
	Affiliation  string
	RepoFullName string
	Permission   githubv4.RepositoryPermission
	UserLogin    string
}

type RepoAlertDependabotDescription struct {
	RepoFullName                string
	AlertNumber                 int
	State                       string
	DependencyPackageEcosystem  string
	DependencyPackageName       string
	DependencyManifestPath      string
	DependencyScope             string
	SecurityAdvisoryGHSAID      string
	SecurityAdvisoryCVEID       string
	SecurityAdvisorySummary     string
	SecurityAdvisoryDescription string
	SecurityAdvisorySeverity    string
	SecurityAdvisoryCVSSScore   *float64
	SecurityAdvisoryCVSSVector  string
	SecurityAdvisoryCWEs        []string
	SecurityAdvisoryPublishedAt github.Timestamp
	SecurityAdvisoryUpdatedAt   github.Timestamp
	SecurityAdvisoryWithdrawnAt github.Timestamp
	URL                         string
	HTMLURL                     string
	CreatedAt                   github.Timestamp
	UpdatedAt                   github.Timestamp
	DismissedAt                 github.Timestamp
	DismissedReason             string
	DismissedComment            string
	FixedAt                     github.Timestamp
}

type RepoDeploymentDescription struct {
	steampipemodels.Deployment
	RepoFullName string
}

type RepoEnvironmentDescription struct {
	steampipemodels.Environment
	RepoFullName string
}

type RepoRuleSetDescription struct {
	steampipemodels.Ruleset
	RepoFullName string
}

type RepoSBOMDescription struct {
	RepositoryFullName string
	SPDXID             string
	SPDXVersion        string
	CreationInfo       *github.CreationInfo
	Name               string
	DataLicense        string
	DocumentDescribes  []string
	DocumentNamespace  string
	Packages           []*github.RepoDependencies
}

type RepoVulnerabilityAlertDescription struct {
	RepositoryFullName         string
	Number                     int
	NodeID                     string
	AutoDismissedAt            steampipemodels.NullableTime
	CreatedAt                  steampipemodels.NullableTime
	DependencyScope            githubv4.RepositoryVulnerabilityAlertDependencyScope
	DismissComment             string
	DismissReason              string
	DismissedAt                steampipemodels.NullableTime
	Dismisser                  steampipemodels.BasicUser
	FixedAt                    steampipemodels.NullableTime
	State                      githubv4.RepositoryVulnerabilityAlertState
	SecurityAdvisory           steampipemodels.SecurityAdvisory
	SecurityVulnerability      steampipemodels.SecurityVulnerability
	VulnerableManifestFilename string
	VulnerableManifestPath     string
	VulnerableRequirements     string
	Severity                   githubv4.SecurityAdvisorySeverity
	CvssScore                  float64
}

type SearchCodeDescription struct {
	*github.CodeResult
	RepoFullName string
	Query        string
}

type SearchCommitDescription struct {
	*github.CommitResult
	RepoFullName string
	Query        string
}

type SearchIssueDescription struct {
	IssueDescription
	RepoFullName string
	Query        string
	TextMatches  []steampipemodels.TextMatch
}

type StarDescription struct {
	RepoFullName string
	StarredAt    steampipemodels.NullableTime
	Url          string
}

type StargazerDescription struct {
	RepoFullName string
	StarredAt    steampipemodels.NullableTime
	UserLogin    string
	UserDetail   steampipemodels.BasicUser
}

type TagDescription struct {
	RepositoryFullName string
	Name               string
	TaggerDate         time.Time
	TaggerName         string
	TaggerLogin        string
	Message            string
	Commit             steampipemodels.BaseCommit
}

type ParentTeam struct {
	Id     int
	NodeId string
	Name   string
	Slug   string
}

type TeamDescription struct {
	Organization           string
	Slug                   string
	Name                   string
	ID                     int
	NodeID                 string
	Description            string
	CreatedAt              time.Time
	UpdatedAt              time.Time
	CombinedSlug           string
	ParentTeam             ParentTeam
	Privacy                string
	AncestorsTotalCount    int
	ChildTeamsTotalCount   int
	DiscussionsTotalCount  int
	InvitationsTotalCount  int
	MembersTotalCount      int
	ProjectsV2TotalCount   int
	RepositoriesTotalCount int
	URL                    string
	AvatarURL              string
	DiscussionsURL         string
	EditTeamURL            string
	MembersURL             string
	NewTeamURL             string
	RepositoriesURL        string
	TeamsURL               string
	CanAdminister          bool
	CanSubscribe           bool
	Subscription           string
}

type TeamMembersDescription struct {
	steampipemodels.User
	Organization string
	Slug         string
	Role         githubv4.TeamMemberRole
}

type TrafficViewDailyDescription struct {
	*github.TrafficData
	RepositoryFullName string
}

type TrafficViewWeeklyDescription struct {
	*github.TrafficData
	RepositoryFullName string
}

type TreeDescription struct {
	TreeSHA            string
	RepositoryFullName string
	Recursive          bool
	Truncated          bool
	SHA                string
	Path               string
	Mode               string
	Type               string
	Size               int
	URL                string
}

type UserDescription struct {
	steampipemodels.User
	RepositoriesTotalDiskUsage    int
	FollowersTotalCount           int
	FollowingTotalCount           int
	PublicRepositoriesTotalCount  int
	PrivateRepositoriesTotalCount int
	PublicGistsTotalCount         int
	IssuesTotalCount              int
	OrganizationsTotalCount       int
	PublicKeysTotalCount          int
	OpenPullRequestsTotalCount    int
	MergedPullRequestsTotalCount  int
	ClosedPullRequestsTotalCount  int
	PackagesTotalCount            int
	PinnedItemsTotalCount         int
	SponsoringTotalCount          int
	SponsorsTotalCount            int
	StarredRepositoriesTotalCount int
	WatchingTotalCount            int
}

type WorkflowDescription struct {
	ID                      *int64
	NodeID                  *string
	Name                    *string
	Path                    *string
	State                   *string
	CreatedAt               *github.Timestamp
	UpdatedAt               *github.Timestamp
	URL                     *string
	HTMLURL                 *string
	BadgeURL                *string
	RepositoryFullName      *string
	WorkFlowFileContent     *string
	WorkFlowFileContentJson *github.RepositoryContent
	Pipeline                *goPipeline.Pipeline
}

type CodeOwnerDescription struct {
	RepositoryFullName string
	LineNumber         int64
	Pattern            string
	Users              []string
	Teams              []string
	PreComments        []string
	LineComment        string
}

type OwnerLogin struct {
	Login string `json:"login"`
}

type Package struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	PackageType string     `json:"package_type"`
	Visibility  string     `json:"visibility"`
	HTMLURL     string     `json:"html_url"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
	Owner       OwnerLogin `json:"owner"`
}

type ContainerMetadata struct {
	Container struct {
		Tags []string
	}
}

type ContainerPackageDescription struct {
	ID         int
	Digest     string
	CreatedAt  string
	UpdatedAt  string
	PackageURL string
	Name       string
	MediaType  string
	TotalSize  int64
	Metadata   ContainerMetadata
	Manifest   interface{}

	// -- Add these new fields/methods: --

	// The GitHub "version.Name" or tag,
	// if you want to store that separately from the real Docker digest.

	// When deduplicating, any subsequent tags for the same (ID,digest)
	// can be appended here.
	AdditionalPackageURIs []string `json:"additional_package_uris,omitempty"`
}

// Provide a helper so that code calling `cpd.ActualDigest()` compiles:
func (c ContainerPackageDescription) ActualDigest() string {
	// If you want the *Docker* digest, we simply return `c.Digest`
	// because your code is already overwriting .Digest with the real Docker digest
	return c.Digest
}

type PackageVersion struct {
	ID             int               `json:"id"`
	Name           string            `json:"name"`
	PackageHTMLURL string            `json:"package_html_url"`
	CreatedAt      string            `json:"created_at"`
	UpdatedAt      string            `json:"updated_at"`
	HTMLURL        string            `json:"html_url"`
	Metadata       ContainerMetadata `json:"metadata"`
}

type OwnerDetail struct {
	Login        string `json:"login"`
	ID           int    `json:"id,omitempty"`
	NodeID       string `json:"node_id,omitempty"`
	HTMLURL      string `json:"html_url,omitempty"`
	Type         string `json:"type,omitempty"`
	UserViewType string `json:"user_view_type,omitempty"`
	SiteAdmin    bool   `json:"site_admin,omitempty"`
}

type RepoOwnerDetail struct {
	Login     string `json:"login"`
	ID        int    `json:"id,omitempty"`
	NodeID    string `json:"node_id,omitempty"`
	HTMLURL   string `json:"html_url,omitempty"`
	Type      string `json:"type,omitempty"`
	SiteAdmin bool   `json:"site_admin,omitempty"`
}

type Repository struct {
	ID          int             `json:"id"`
	NodeID      string          `json:"node_id"`
	Name        string          `json:"name"`
	FullName    string          `json:"full_name"`
	Private     bool            `json:"private"`
	Owner       RepoOwnerDetail `json:"owner"`
	HTMLURL     string          `json:"html_url"`
	Description string          `json:"description"`
	Fork        bool            `json:"fork"`
	URL         string          `json:"url"`
}

type PackageListItem struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	PackageType string `json:"package_type"`
	Visibility  string `json:"visibility"`
	HTMLURL     string `json:"html_url"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Owner       struct {
		Login string `json:"login"`
	} `json:"owner"`
	URL string `json:"url"`
}

type PackageDetailDescription struct {
	ID           int         `json:"id"`
	Name         string      `json:"name"`
	PackageType  string      `json:"package_type"`
	Owner        OwnerDetail `json:"owner"`
	VersionCount int         `json:"version_count"`
	Visibility   string      `json:"visibility"`
	URL          string      `json:"url"`
	CreatedAt    string      `json:"created_at"`
	UpdatedAt    string      `json:"updated_at"`
	Repository   Repository  `json:"repository"`
	HTMLURL      string      `json:"html_url"`
}

type PackageDescription struct {
	ID         string
	RegistryID string
	Name       string
	URL        string
	CreatedAt  github.Timestamp
	UpdatedAt  github.Timestamp
}

type PackageVersionDescription struct {
	ID          int
	Name        string
	PackageName string
	VersionURI  string
	Digest      *string
	CreatedAt   github.Timestamp
	UpdatedAt   github.Timestamp
}

type CodeSearchResult struct {
	TotalCount        int             `json:"total_count"`
	IncompleteResults bool            `json:"incomplete_results"`
	Items             []CodeSearchHit `json:"items"`
}

type CodeSearchHit struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	Sha        string `json:"sha"`
	URL        string `json:"url"`
	GitURL     string `json:"git_url"`
	HTMLURL    string `json:"html_url"`
	Repository struct {
		ID       int    `json:"id"`
		NodeID   string `json:"node_id"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		Private  bool   `json:"private"`
		Owner    struct {
			Login   string `json:"login"`
			ID      int    `json:"id"`
			NodeID  string `json:"node_id"`
			URL     string `json:"url"`
			HTMLURL string `json:"html_url"`
			Type    string `json:"type"`
		} `json:"owner"`
		HTMLURL     string `json:"html_url"`
		Description string `json:"description"`
		Fork        bool   `json:"fork"`
	} `json:"repository"`
	Score float64 `json:"score"`
}

type ContentResponse struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Sha      string `json:"sha"`
	Size     int    `json:"size"`
	URL      string `json:"url"`
	HTMLURL  string `json:"html_url"`
	GitURL   string `json:"git_url"`
	Type     string `json:"type"`
	Content  string `json:"content"` // base64
	Encoding string `json:"encoding"`
}

type CommitResponse struct {
	Commit struct {
		Author struct {
			Date string `json:"date"`
		} `json:"author"`
		Committer struct {
			Date string `json:"date"`
		} `json:"committer"`
	} `json:"commit"`
}

type ArtifactDockerFileDescription struct {
	Sha  *string
	Name *string
	//Path                    *string
	LastUpdatedAt *string
	//GitURL                  *string
	HTMLURL *string
	//URI                     *string // Unique identifier
	DockerfileContent       string
	DockerfileContentBase64 *string
	Repository              map[string]interface{}
	Images                  []string // New field to store extracted base images
}
