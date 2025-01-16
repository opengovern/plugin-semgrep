package describers

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/google/go-github/v55/github"
	"github.com/opengovern/og-describer-template/discovery/pkg/models"
	model "github.com/opengovern/og-describer-template/discovery/provider"
	resilientbridge "github.com/opengovern/resilient-bridge"
	"github.com/shurcooL/githubv4"
)

// model.GitHubClient custom struct for defining both rest and graphql clients
type GitHubClient struct {
	RestClient    *github.Client
	GraphQLClient *githubv4.Client
	Token         string
}

const (
	maxPagesCount            = 30
	pageSize                 = 30
	pullRequestsPageSize     = 25
	repoPageSize             = 15
	issuePageSize            = 50
	orgPageSize              = 10
	orgCollaboratorsPageSize = 30
	releasePageSize          = 20
	teamMembersPageSize      = 30
	packagePageSize          = 100
)

func appendRepoColumnIncludes(m *map[string]interface{}, cols []string) {
	optionals := map[string]string{
		"allow_update_branch":              "includeAllowUpdateBranch",
		"archived_at":                      "includeArchivedAt",
		"auto_merge_allowed":               "includeAutoMergeAllowed",
		"can_administer":                   "includeCanAdminister",
		"can_create_projects":              "includeCanCreateProjects",
		"can_subscribe":                    "includeCanSubscribe",
		"can_update_topics":                "includeCanUpdateTopics",
		"code_of_conduct":                  "includeCodeOfConduct",
		"contact_links":                    "includeContactLinks",
		"created_at":                       "includeCreatedAt",
		"default_branch_ref":               "includeDefaultBranchRef",
		"delete_branch_on_merge":           "includeDeleteBranchOnMerge",
		"description":                      "includeDescription",
		"disk_usage":                       "includeDiskUsage",
		"fork_count":                       "includeForkCount",
		"forking_allowed":                  "includeForkingAllowed",
		"funding_links":                    "includeFundingLinks",
		"has_discussions_enabled":          "includeHasDiscussionsEnabled",
		"has_issues_enabled":               "includeHasIssuesEnabled",
		"has_projects_enabled":             "includeHasProjectsEnabled",
		"has_starred":                      "includeHasStarred",
		"has_vulnerability_alerts_enabled": "includeHasVulnerabilityAlertsEnabled",
		"has_wiki_enabled":                 "includeHasWikiEnabled",
		"homepage_url":                     "includeHomepageUrl",
		"interaction_ability":              "includeInteractionAbility",
		"is_archived":                      "includeIsArchived",
		"is_blank_issues_enabled":          "includeIsBlankIssuesEnabled",
		"is_disabled":                      "includeIsDisabled",
		"is_empty":                         "includeIsEmpty",
		"is_fork":                          "includeIsFork",
		"is_in_organization":               "includeIsInOrganization",
		"is_locked":                        "includeIsLocked",
		"is_mirror":                        "includeIsMirror",
		"is_private":                       "includeIsPrivate",
		"is_security_policy_enabled":       "includeIsSecurityPolicyEnabled",
		"is_template":                      "includeIsTemplate",
		"is_user_configuration_repository": "includeIsUserConfigurationRepository",
		"issue_templates":                  "includeIssueTemplates",
		"license_info":                     "includeLicenseInfo",
		"lock_reason":                      "includeLockReason",
		"merge_commit_allowed":             "includeMergeCommitAllowed",
		"merge_commit_message":             "includeMergeCommitMessage",
		"merge_commit_title":               "includeMergeCommitTitle",
		"mirror_url":                       "includeMirrorUrl",
		"open_graph_image_url":             "includeOpenGraphImageUrl",
		"open_issues_total_count":          "includeOpenIssues",
		"possible_commit_emails":           "includePossibleCommitEmails",
		"primary_language":                 "includePrimaryLanguage",
		"projects_url":                     "includeProjectsUrl",
		"pull_request_templates":           "includePullRequestTemplates",
		"pushed_at":                        "includePushedAt",
		"rebase_merge_allowed":             "includeRebaseMergeAllowed",
		"repository_topics_total_count":    "includeRepositoryTopics",
		"security_policy_url":              "includeSecurityPolicyUrl",
		"squash_merge_allowed":             "includeSquashMergeAllowed",
		"squash_merge_commit_message":      "includeSquashMergeCommitMessage",
		"squash_merge_commit_title":        "includeSquashMergeCommitTitle",
		"ssh_url":                          "includeSshUrl",
		"stargazer_count":                  "includeStargazerCount",
		"subscription":                     "includeSubscription",
		"updated_at":                       "includeUpdatedAt",
		"url":                              "includeUrl",
		"uses_custom_open_graph_image":     "includeUsesCustomOpenGraphImage",
		"visibility":                       "includeVisibility",
		"watchers_total_count":             "includeWatchers",
		"web_commit_signoff_required":      "includeWebCommitSignoffRequired",
		"your_permission":                  "includeYourPermission",
	}
	for key, value := range optionals {
		(*m)[value] = githubv4.Boolean(slices.Contains(cols, key))
	}
}

func appendBranchColumnIncludes(m *map[string]interface{}, cols []string) {
	protectionIncluded := githubv4.Boolean(slices.Contains(cols, "protected") || slices.Contains(cols, "branch_protection_rule"))
	(*m)["includeBranchProtectionRule"] = protectionIncluded
	(*m)["includeAllowsDeletions"] = protectionIncluded
	(*m)["includeAllowsForcePushes"] = protectionIncluded
	(*m)["includeBlocksCreations"] = protectionIncluded
	(*m)["includeCreator"] = protectionIncluded
	(*m)["includeBranchProtectionRuleId"] = protectionIncluded
	(*m)["includeDismissesStaleReviews"] = protectionIncluded
	(*m)["includeIsAdminEnforced"] = protectionIncluded
	(*m)["includeLockAllowsFetchAndMerge"] = protectionIncluded
	(*m)["includeLockBranch"] = protectionIncluded
	(*m)["includePattern"] = protectionIncluded
	(*m)["includeRequireLastPushApproval"] = protectionIncluded
	(*m)["includeRequiredApprovingReviewCount"] = protectionIncluded
	(*m)["includeRequiredDeploymentEnvironments"] = protectionIncluded
	(*m)["includeRequiredStatusChecks"] = protectionIncluded
	(*m)["includeRequiresApprovingReviews"] = protectionIncluded
	(*m)["includeRequiresConversationResolution"] = protectionIncluded
	(*m)["includeRequiresCodeOwnerReviews"] = protectionIncluded
	(*m)["includeRequiresCommitSignatures"] = protectionIncluded
	(*m)["includeRequiresDeployments"] = protectionIncluded
	(*m)["includeRequiresLinearHistory"] = protectionIncluded
	(*m)["includeRequiresStatusChecks"] = protectionIncluded
	(*m)["includeRequiresStrictStatusChecks"] = protectionIncluded
	(*m)["includeRestrictsPushes"] = protectionIncluded
	(*m)["includeRestrictsReviewDismissals"] = protectionIncluded
	(*m)["includeMatchingBranches"] = protectionIncluded
}

func appendBranchProtectionRuleColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeAllowsDeletions"] = githubv4.Boolean(slices.Contains(cols, "allows_deletions"))
	(*m)["includeAllowsForcePushes"] = githubv4.Boolean(slices.Contains(cols, "allows_force_pushes"))
	(*m)["includeBlocksCreations"] = githubv4.Boolean(slices.Contains(cols, "blocks_creations"))
	(*m)["includeCreator"] = githubv4.Boolean(slices.Contains(cols, "creator") || slices.Contains(cols, "creator_login"))
	(*m)["includeBranchProtectionRuleId"] = githubv4.Boolean(slices.Contains(cols, "id"))
	(*m)["includeDismissesStaleReviews"] = githubv4.Boolean(slices.Contains(cols, "dismisses_stale_reviews"))
	(*m)["includeIsAdminEnforced"] = githubv4.Boolean(slices.Contains(cols, "is_admin_enforced"))
	(*m)["includeLockAllowsFetchAndMerge"] = githubv4.Boolean(slices.Contains(cols, "lock_allows_fetch_and_merge"))
	(*m)["includeLockBranch"] = githubv4.Boolean(slices.Contains(cols, "lock_branch"))
	(*m)["includePattern"] = githubv4.Boolean(slices.Contains(cols, "pattern"))
	(*m)["includeRequireLastPushApproval"] = githubv4.Boolean(slices.Contains(cols, "require_last_push_approval"))
	(*m)["includeRequiredApprovingReviewCount"] = githubv4.Boolean(slices.Contains(cols, "required_approving_review_count"))
	(*m)["includeRequiredDeploymentEnvironments"] = githubv4.Boolean(slices.Contains(cols, "required_deployment_environments"))
	(*m)["includeRequiredStatusChecks"] = githubv4.Boolean(slices.Contains(cols, "required_status_checks"))
	(*m)["includeRequiresApprovingReviews"] = githubv4.Boolean(slices.Contains(cols, "requires_approving_reviews"))
	(*m)["includeRequiresConversationResolution"] = githubv4.Boolean(slices.Contains(cols, "requires_conversation_resolution"))
	(*m)["includeRequiresCodeOwnerReviews"] = githubv4.Boolean(slices.Contains(cols, "requires_code_owner_reviews"))
	(*m)["includeRequiresCommitSignatures"] = githubv4.Boolean(slices.Contains(cols, "requires_commit_signatures"))
	(*m)["includeRequiresDeployments"] = githubv4.Boolean(slices.Contains(cols, "requires_deployments"))
	(*m)["includeRequiresLinearHistory"] = githubv4.Boolean(slices.Contains(cols, "requires_linear_history"))
	(*m)["includeRequiresStatusChecks"] = githubv4.Boolean(slices.Contains(cols, "requires_status_checks"))
	(*m)["includeRequiresStrictStatusChecks"] = githubv4.Boolean(slices.Contains(cols, "requires_strict_status_checks"))
	(*m)["includeRestrictsPushes"] = githubv4.Boolean(slices.Contains(cols, "restricts_pushes"))
	(*m)["includeRestrictsReviewDismissals"] = githubv4.Boolean(slices.Contains(cols, "restricts_review_dismissals"))
	(*m)["includeMatchingBranches"] = githubv4.Boolean(slices.Contains(cols, "matching_branches"))
}

func appendCommitColumnIncludes(m *map[string]interface{}, cols []string) {
	// For BasicCommit struct
	(*m)["includeCommitShortSha"] = githubv4.Boolean(slices.Contains(cols, "short_sha"))
	(*m)["includeCommitAuthoredDate"] = githubv4.Boolean(slices.Contains(cols, "authored_date"))
	(*m)["includeCommitAuthor"] = githubv4.Boolean(slices.Contains(cols, "author") || slices.Contains(cols, "author_login"))
	(*m)["includeCommitCommittedDate"] = githubv4.Boolean(slices.Contains(cols, "committed_date"))
	(*m)["includeCommitCommitter"] = githubv4.Boolean(slices.Contains(cols, "committer") || slices.Contains(cols, "committer_login"))
	(*m)["includeCommitMessage"] = githubv4.Boolean(slices.Contains(cols, "message"))
	(*m)["includeCommitUrl"] = githubv4.Boolean(slices.Contains(cols, "url"))
	// For Commit struct
	(*m)["includeCommitAdditions"] = githubv4.Boolean(slices.Contains(cols, "additions"))
	(*m)["includeCommitAuthoredByCommitter"] = githubv4.Boolean(slices.Contains(cols, "authored_by_committer"))
	(*m)["includeCommitChangedFiles"] = githubv4.Boolean(slices.Contains(cols, "changed_files"))
	(*m)["includeCommitCommittedViaWeb"] = githubv4.Boolean(slices.Contains(cols, "committed_via_web"))
	(*m)["includeCommitCommitUrl"] = githubv4.Boolean(slices.Contains(cols, "commit_url"))
	(*m)["includeCommitDeletions"] = githubv4.Boolean(slices.Contains(cols, "deletions"))
	(*m)["includeCommitSignature"] = githubv4.Boolean(slices.Contains(cols, "signature"))
	(*m)["includeCommitTarballUrl"] = githubv4.Boolean(slices.Contains(cols, "tarball_url"))
	(*m)["includeCommitTreeUrl"] = githubv4.Boolean(slices.Contains(cols, "tree_url"))
	(*m)["includeCommitCanSubscribe"] = githubv4.Boolean(slices.Contains(cols, "can_subscribe"))
	(*m)["includeCommitSubscription"] = githubv4.Boolean(slices.Contains(cols, "subscription"))
	(*m)["includeCommitZipballUrl"] = githubv4.Boolean(slices.Contains(cols, "zipball_url"))
	(*m)["includeCommitMessageHeadline"] = githubv4.Boolean(slices.Contains(cols, "message_headline"))
	(*m)["includeCommitStatus"] = githubv4.Boolean(slices.Contains(cols, "status"))
	(*m)["includeCommitNodeId"] = githubv4.Boolean(slices.Contains(cols, "node_id"))
}

func appendOrganizationColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeAnnouncement"] = githubv4.Boolean(slices.Contains(cols, "announcement"))
	(*m)["includeAnnouncementExpiresAt"] = githubv4.Boolean(slices.Contains(cols, "announcement_expires_at"))
	(*m)["includeAnnouncementUserDismissible"] = githubv4.Boolean(slices.Contains(cols, "announcement_user_dismissible"))
	(*m)["includeAnyPinnableItems"] = githubv4.Boolean(slices.Contains(cols, "any_pinnable_items"))
	(*m)["includeAvatarUrl"] = githubv4.Boolean(slices.Contains(cols, "avatar_url"))
	(*m)["includeEstimatedNextSponsorsPayoutInCents"] = githubv4.Boolean(slices.Contains(cols, "estimated_next_sponsors_payout_in_cents"))
	(*m)["includeHasSponsorsListing"] = githubv4.Boolean(slices.Contains(cols, "has_sponsors_listing"))
	(*m)["includeInteractionAbility"] = githubv4.Boolean(slices.Contains(cols, "interaction_ability"))
	(*m)["includeIsSponsoringYou"] = githubv4.Boolean(slices.Contains(cols, "is_sponsoring_you"))
	(*m)["includeIsVerified"] = githubv4.Boolean(slices.Contains(cols, "is_verified"))
	(*m)["includeLocation"] = githubv4.Boolean(slices.Contains(cols, "location"))
	(*m)["includeMonthlyEstimatedSponsorsIncomeInCents"] = githubv4.Boolean(slices.Contains(cols, "monthly_estimated_sponsors_income_in_cents"))
	(*m)["includeNewTeamUrl"] = githubv4.Boolean(slices.Contains(cols, "new_team_url"))
	(*m)["includePinnedItemsRemaining"] = githubv4.Boolean(slices.Contains(cols, "pinned_items_remaining"))
	(*m)["includeProjectsUrl"] = githubv4.Boolean(slices.Contains(cols, "projects_url"))
	(*m)["includeSamlIdentityProvider"] = githubv4.Boolean(slices.Contains(cols, "saml_identity_provider"))
	(*m)["includeSponsorsListing"] = githubv4.Boolean(slices.Contains(cols, "sponsors_listing"))
	(*m)["includeTeamsUrl"] = githubv4.Boolean(slices.Contains(cols, "teams_url"))
	(*m)["includeTotalSponsorshipAmountAsSponsorInCents"] = githubv4.Boolean(slices.Contains(cols, "total_sponsorship_amount_as_sponsor_in_cents"))
	(*m)["includeTwitterUsername"] = githubv4.Boolean(slices.Contains(cols, "twitter_username"))
	(*m)["includeOrgViewer"] = githubv4.Boolean(slices.Contains(cols, "can_administer") || slices.Contains(cols, "can_changed_pinned_items") || slices.Contains(cols, "can_create_projects") || slices.Contains(cols, "can_create_repositories") || slices.Contains(cols, "can_create_teams") || slices.Contains(cols, "can_sponsor"))
	(*m)["includeIsAMember"] = githubv4.Boolean(slices.Contains(cols, "is_a_member"))
	(*m)["includeIsFollowing"] = githubv4.Boolean(slices.Contains(cols, "is_following"))
	(*m)["includeIsSponsoring"] = githubv4.Boolean(slices.Contains(cols, "is_sponsoring"))
	(*m)["includeWebsiteUrl"] = githubv4.Boolean(slices.Contains(cols, "website_url"))
	(*m)["includeMembersWithRole"] = githubv4.Boolean(slices.Contains(cols, "members_with_role_total_count"))
	(*m)["includePackages"] = githubv4.Boolean(slices.Contains(cols, "packages_total_count"))
	(*m)["includePinnableItems"] = githubv4.Boolean(slices.Contains(cols, "pinnable_items_total_count"))
	(*m)["includePinnedItems"] = githubv4.Boolean(slices.Contains(cols, "pinned_items_total_count"))
	(*m)["includeProjects"] = githubv4.Boolean(slices.Contains(cols, "projects_total_count"))
	(*m)["includeProjectsV2"] = githubv4.Boolean(slices.Contains(cols, "projects_v2_total_count"))
	(*m)["includeSponsoring"] = githubv4.Boolean(slices.Contains(cols, "sponsoring_total_count"))
	(*m)["includeSponsors"] = githubv4.Boolean(slices.Contains(cols, "sponsors_total_count"))
	(*m)["includeTeams"] = githubv4.Boolean(slices.Contains(cols, "teams_total_count"))
	(*m)["includePrivateRepositories"] = githubv4.Boolean(slices.Contains(cols, "private_repositories_total_count"))
	(*m)["includePublicRepositories"] = githubv4.Boolean(slices.Contains(cols, "public_repositories_total_count"))
	(*m)["includeRepositories"] = githubv4.Boolean(slices.Contains(cols, "repositories_total_count"))
	(*m)["includeRepositories"] = githubv4.Boolean(slices.Contains(cols, "repositories_total_disk_usage"))
}

func appendStarColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeStarNode"] = githubv4.Boolean(slices.Contains(cols, "repository_full_name") || slices.Contains(cols, "url"))
	(*m)["includeStarEdges"] = githubv4.Boolean(slices.Contains(cols, "starred_at"))
}

func appendIssueColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeIssueAuthor"] = githubv4.Boolean(slices.Contains(cols, "author") || slices.Contains(cols, "author_login"))
	(*m)["includeIssueBody"] = githubv4.Boolean(slices.Contains(cols, "body"))
	(*m)["includeIssueEditor"] = githubv4.Boolean(slices.Contains(cols, "editor"))
	(*m)["includeIssueMilestone"] = githubv4.Boolean(slices.Contains(cols, "milestone"))
	(*m)["includeIssueViewer"] = githubv4.Boolean(slices.Contains(cols, "user_can_close") ||
		slices.Contains(cols, "user_can_react") ||
		slices.Contains(cols, "user_can_reopen") ||
		slices.Contains(cols, "user_can_subscribe") ||
		slices.Contains(cols, "user_can_update") ||
		slices.Contains(cols, "user_cannot_update_reasons") ||
		slices.Contains(cols, "user_did_author") ||
		slices.Contains(cols, "user_subscription"))
	(*m)["includeIssueAssignees"] = githubv4.Boolean(slices.Contains(cols, "assignees_total_count") || slices.Contains(cols, "assignees"))
	(*m)["includeIssueCommentCount"] = githubv4.Boolean(slices.Contains(cols, "comments_total_count"))
	(*m)["includeIssueLabels"] = githubv4.Boolean(slices.Contains(cols, "labels") ||
		slices.Contains(cols, "labels_src") ||
		slices.Contains(cols, "labels_total_count"))
	(*m)["includeIssueUrl"] = githubv4.Boolean(slices.Contains(cols, "url"))
	(*m)["includeIssueUpdatedAt"] = githubv4.Boolean(slices.Contains(cols, "updated_at"))
	(*m)["includeIssueTitle"] = githubv4.Boolean(slices.Contains(cols, "title"))
	(*m)["includeIssueStateReason"] = githubv4.Boolean(slices.Contains(cols, "state_reason"))
	(*m)["includeIssueState"] = githubv4.Boolean(slices.Contains(cols, "state"))
	(*m)["includeIssuePublishedAt"] = githubv4.Boolean(slices.Contains(cols, "published_at"))
	(*m)["includeIssueLocked"] = githubv4.Boolean(slices.Contains(cols, "locked"))
	(*m)["includeIssueLastEditedAt"] = githubv4.Boolean(slices.Contains(cols, "last_edited_at"))
	(*m)["includeIssueIsPinned"] = githubv4.Boolean(slices.Contains(cols, "is_pinned"))
	(*m)["includeIssueIncludesCreatedEdit"] = githubv4.Boolean(slices.Contains(cols, "includes_created_edit"))
	(*m)["includeIssueFullDatabaseId"] = githubv4.Boolean(slices.Contains(cols, "full_database_id"))
	(*m)["includeIssueCreatedViaEmail"] = githubv4.Boolean(slices.Contains(cols, "created_via_email"))
	(*m)["includeIssueCreatedAt"] = githubv4.Boolean(slices.Contains(cols, "created_at"))
	(*m)["includeIssueClosedAt"] = githubv4.Boolean(slices.Contains(cols, "closed_at"))
	(*m)["includeIssueClosed"] = githubv4.Boolean(slices.Contains(cols, "closed"))
	(*m)["includeIssueBodyUrl"] = githubv4.Boolean(slices.Contains(cols, "body_url"))
	(*m)["includeIssueAuthorAssociation"] = githubv4.Boolean(slices.Contains(cols, "author_association"))
	(*m)["includeIssueActiveLockReason"] = githubv4.Boolean(slices.Contains(cols, "active_lock_reason"))
	(*m)["includeIssueNodeId"] = githubv4.Boolean(slices.Contains(cols, "node_id"))
	(*m)["includeIssueId"] = githubv4.Boolean(slices.Contains(cols, "id"))
	(*m)["includeIssueIsReadByUser"] = githubv4.Boolean(slices.Contains(cols, "is_read_by_user"))
}

func appendIssuePRCommentColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeIssueCommentAuthor"] = githubv4.Boolean(slices.Contains(cols, "author") || slices.Contains(cols, "author_login"))
	(*m)["includeIssueCommentBody"] = githubv4.Boolean(slices.Contains(cols, "body"))
	(*m)["includeIssueCommentEditor"] = githubv4.Boolean(slices.Contains(cols, "editor") || slices.Contains(cols, "editor_login"))
	(*m)["includeIssueCommentViewer"] = githubv4.Boolean(slices.Contains(cols, "can_delete") ||
		slices.Contains(cols, "can_react") ||
		slices.Contains(cols, "can_minimize") ||
		slices.Contains(cols, "can_update") ||
		slices.Contains(cols, "cannot_update_reasons") ||
		slices.Contains(cols, "did_author"))
	(*m)["includeIssueCommentUrl"] = githubv4.Boolean(slices.Contains(cols, "url"))
	(*m)["includeIssueCommentUpdatedAt"] = githubv4.Boolean(slices.Contains(cols, "updated_at"))
	(*m)["includeIssueCommentPublishedAt"] = githubv4.Boolean(slices.Contains(cols, "published_at"))
	(*m)["includeIssueCommentMinimizedReason"] = githubv4.Boolean(slices.Contains(cols, "minimized_reason"))
	(*m)["includeIssueCommentLastEditedAt"] = githubv4.Boolean(slices.Contains(cols, "last_edited_at"))
	(*m)["includeIssueCommentIsMinimized"] = githubv4.Boolean(slices.Contains(cols, "is_minimized"))
	(*m)["includeIssueCommentIncludesCreatedEdit"] = githubv4.Boolean(slices.Contains(cols, "includes_created_edit"))
	(*m)["includeIssueCommentCreatedViaEmail"] = githubv4.Boolean(slices.Contains(cols, "created_via_email"))
	(*m)["includeIssueCommentCreatedAt"] = githubv4.Boolean(slices.Contains(cols, "created_at"))
	(*m)["includeIssueCommentBody"] = githubv4.Boolean(slices.Contains(cols, "body"))
	(*m)["includeIssueCommentBodyText"] = githubv4.Boolean(slices.Contains(cols, "body_text"))
	(*m)["includeIssueCommentAuthorAssociation"] = githubv4.Boolean(slices.Contains(cols, "author_association"))
	(*m)["includeIssueCommentNodeId"] = githubv4.Boolean(slices.Contains(cols, "node_id"))
	(*m)["includeIssueCommentId"] = githubv4.Boolean(slices.Contains(cols, "id"))
}

func appendLicenseColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeLicenseName"] = githubv4.Boolean(slices.Contains(cols, "name"))
	(*m)["includeLicenseSpdxId"] = githubv4.Boolean(slices.Contains(cols, "spdx_id"))
	(*m)["includeLicenseUrl"] = githubv4.Boolean(slices.Contains(cols, "url"))
	(*m)["includeLicenseConditions"] = githubv4.Boolean(slices.Contains(cols, "conditions"))
	(*m)["includeLicenseDescription"] = githubv4.Boolean(slices.Contains(cols, "description"))
	(*m)["includeLicenseFeatured"] = githubv4.Boolean(slices.Contains(cols, "featured"))
	(*m)["includeLicenseHidden"] = githubv4.Boolean(slices.Contains(cols, "hidden"))
	(*m)["includeLicenseImplementation"] = githubv4.Boolean(slices.Contains(cols, "implementation"))
	(*m)["includeLicenseLimitations"] = githubv4.Boolean(slices.Contains(cols, "limitations"))
	(*m)["includeLicensePermissions"] = githubv4.Boolean(slices.Contains(cols, "permissions"))
	(*m)["includeLicenseNickname"] = githubv4.Boolean(slices.Contains(cols, "nickname"))
	(*m)["includeLicensePseudoLicense"] = githubv4.Boolean(slices.Contains(cols, "pseudo_license"))
}

func appendTeamColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeTeamAvatarUrl"] = githubv4.Boolean(slices.Contains(cols, "avatar_url"))
	(*m)["includeTeamCombinedSlug"] = githubv4.Boolean(slices.Contains(cols, "combined_slug"))
	(*m)["includeTeamCreatedAt"] = githubv4.Boolean(slices.Contains(cols, "created_at"))
	(*m)["includeTeamDescription"] = githubv4.Boolean(slices.Contains(cols, "description"))
	(*m)["includeTeamDiscussionsUrl"] = githubv4.Boolean(slices.Contains(cols, "discussions_url"))
	(*m)["includeTeamEditTeamUrl"] = githubv4.Boolean(slices.Contains(cols, "edit_team_url"))
	(*m)["includeTeamMembersUrl"] = githubv4.Boolean(slices.Contains(cols, "members_url"))
	(*m)["includeTeamNewTeamUrl"] = githubv4.Boolean(slices.Contains(cols, "new_team_url"))
	(*m)["includeTeamParentTeam"] = githubv4.Boolean(slices.Contains(cols, "parent_team"))
	(*m)["includeTeamPrivacy"] = githubv4.Boolean(slices.Contains(cols, "privacy"))
	(*m)["includeTeamRepositoriesUrl"] = githubv4.Boolean(slices.Contains(cols, "repositories_url"))
	(*m)["includeTeamTeamsUrl"] = githubv4.Boolean(slices.Contains(cols, "teams_url"))
	(*m)["includeTeamUpdatedAt"] = githubv4.Boolean(slices.Contains(cols, "updated_at"))
	(*m)["includeTeamUrl"] = githubv4.Boolean(slices.Contains(cols, "url"))
	(*m)["includeTeamCanAdminister"] = githubv4.Boolean(slices.Contains(cols, "can_administer"))
	(*m)["includeTeamCanSubscribe"] = githubv4.Boolean(slices.Contains(cols, "can_subscribe"))
	(*m)["includeTeamSubscription"] = githubv4.Boolean(slices.Contains(cols, "subscription"))
	(*m)["includeTeamAncestors"] = githubv4.Boolean(slices.Contains(cols, "ancestors_total_count"))
	(*m)["includeTeamChildTeams"] = githubv4.Boolean(slices.Contains(cols, "child_teams_total_count"))
	(*m)["includeTeamDiscussions"] = githubv4.Boolean(slices.Contains(cols, "discussions_total_count"))
	(*m)["includeTeamInvitations"] = githubv4.Boolean(slices.Contains(cols, "invitations_total_count"))
	(*m)["includeTeamMembers"] = githubv4.Boolean(slices.Contains(cols, "members_total_count"))
	(*m)["includeTeamProjectsV2"] = githubv4.Boolean(slices.Contains(cols, "projects_v2_total_count"))
	(*m)["includeTeamRepositories"] = githubv4.Boolean(slices.Contains(cols, "repositories_total_count"))
}

func appendOrgCollaboratorColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeOCPermission"] = githubv4.Boolean(slices.Contains(cols, "permission"))
	(*m)["includeOCNode"] = githubv4.Boolean(slices.Contains(cols, "user_login"))
}

func appendOrganizationExternalIdentityColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeOrgExternalIdentityGuid"] = githubv4.Boolean(slices.Contains(cols, "guid"))
	(*m)["includeOrgExternalIdentityUser"] = githubv4.Boolean(slices.Contains(cols, "user_detail") || slices.Contains(cols, "user_login"))
	(*m)["includeOrgExternalIdentitySamlIdentity"] = githubv4.Boolean(slices.Contains(cols, "saml_identity"))
	(*m)["includeOrgExternalIdentityScimIdentity"] = githubv4.Boolean(slices.Contains(cols, "scim_identity"))
	(*m)["includeOrgExternalIdentityOrganizationInvitation"] = githubv4.Boolean(slices.Contains(cols, "organization_invitation"))
}

func appendUserColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeUserAnyPinnableItems"] = githubv4.Boolean(slices.Contains(cols, "any_pinnable_items"))
	(*m)["includeUserAvatarUrl"] = githubv4.Boolean(slices.Contains(cols, "avatar_url"))
	(*m)["includeUserBio"] = githubv4.Boolean(slices.Contains(cols, "bio"))
	(*m)["includeUserCompany"] = githubv4.Boolean(slices.Contains(cols, "company"))
	(*m)["includeUserEstimatedNextSponsorsPayoutInCents"] = githubv4.Boolean(slices.Contains(cols, "estimated_next_sponsors_payout_in_cents"))
	(*m)["includeUserHasSponsorsListing"] = githubv4.Boolean(slices.Contains(cols, "has_sponsors_listing"))
	(*m)["includeUserInteractionAbility"] = githubv4.Boolean(slices.Contains(cols, "interaction_ability"))
	(*m)["includeUserIsBountyHunter"] = githubv4.Boolean(slices.Contains(cols, "is_bounty_hunter"))
	(*m)["includeUserIsCampusExpert"] = githubv4.Boolean(slices.Contains(cols, "is_campus_expert"))
	(*m)["includeUserIsDeveloperProgramMember"] = githubv4.Boolean(slices.Contains(cols, "is_developer_program_member"))
	(*m)["includeUserIsEmployee"] = githubv4.Boolean(slices.Contains(cols, "is_employee"))
	(*m)["includeUserIsFollowingYou"] = githubv4.Boolean(slices.Contains(cols, "is_following_you"))
	(*m)["includeUserIsGitHubStar"] = githubv4.Boolean(slices.Contains(cols, "is_github_star"))
	(*m)["includeUserIsHireable"] = githubv4.Boolean(slices.Contains(cols, "is_hireable"))
	(*m)["includeUserIsSiteAdmin"] = githubv4.Boolean(slices.Contains(cols, "is_site_admin"))
	(*m)["includeUserIsSponsoringYou"] = githubv4.Boolean(slices.Contains(cols, "is_sponsoring_you"))
	(*m)["includeUserIsYou"] = githubv4.Boolean(slices.Contains(cols, "is_you"))
	(*m)["includeUserLocation"] = githubv4.Boolean(slices.Contains(cols, "location"))
	(*m)["includeUserMonthlyEstimatedSponsorsIncomeInCents"] = githubv4.Boolean(slices.Contains(cols, "monthly_estimated_sponsors_income_in_cents"))
	(*m)["includeUserPinnedItemsRemaining"] = githubv4.Boolean(slices.Contains(cols, "pinned_items_remaining"))
	(*m)["includeUserProjectsUrl"] = githubv4.Boolean(slices.Contains(cols, "projects_url"))
	(*m)["includeUserPronouns"] = githubv4.Boolean(slices.Contains(cols, "pronouns"))
	(*m)["includeUserSponsorsListing"] = githubv4.Boolean(slices.Contains(cols, "sponsors_listing"))
	(*m)["includeUserStatus"] = githubv4.Boolean(slices.Contains(cols, "status"))
	(*m)["includeUserTwitterUsername"] = githubv4.Boolean(slices.Contains(cols, "twitter_username"))
	(*m)["includeUserCanChangedPinnedItems"] = githubv4.Boolean(slices.Contains(cols, "can_changed_pinned_items"))
	(*m)["includeUserCanCreateProjects"] = githubv4.Boolean(slices.Contains(cols, "can_create_projects"))
	(*m)["includeUserCanFollow"] = githubv4.Boolean(slices.Contains(cols, "can_follow"))
	(*m)["includeUserCanSponsor"] = githubv4.Boolean(slices.Contains(cols, "can_sponsor"))
	(*m)["includeUserIsFollowing"] = githubv4.Boolean(slices.Contains(cols, "is_following"))
	(*m)["includeUserIsSponsoring"] = githubv4.Boolean(slices.Contains(cols, "is_sponsoring"))
	(*m)["includeUserWebsiteUrl"] = githubv4.Boolean(slices.Contains(cols, "website_url"))
}

func appendRepoCollaboratorColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeRCPermission"] = githubv4.Boolean(slices.Contains(cols, "permission"))
	(*m)["includeRCNode"] = githubv4.Boolean(slices.Contains(cols, "user_login"))
}

func appendRepoDeploymentColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeDeploymentId"] = githubv4.Boolean(slices.Contains(cols, "id"))
	(*m)["includeDeploymentNodeId"] = githubv4.Boolean(slices.Contains(cols, "node_id"))
	(*m)["includeDeploymentCommitSha"] = githubv4.Boolean(slices.Contains(cols, "sha"))
	(*m)["includeDeploymentCreatedAt"] = githubv4.Boolean(slices.Contains(cols, "created_at"))
	(*m)["includeDeploymentCreator"] = githubv4.Boolean(slices.Contains(cols, "creator"))
	(*m)["includeDeploymentDescription"] = githubv4.Boolean(slices.Contains(cols, "description"))
	(*m)["includeDeploymentEnvironment"] = githubv4.Boolean(slices.Contains(cols, "environment"))
	(*m)["includeDeploymentLatestEnvironment"] = githubv4.Boolean(slices.Contains(cols, "latest_environment"))
	(*m)["includeDeploymentLatestStatus"] = githubv4.Boolean(slices.Contains(cols, "latest_status"))
	(*m)["includeDeploymentOriginalEnvironment"] = githubv4.Boolean(slices.Contains(cols, "original_environment"))
	(*m)["includeDeploymentPayload"] = githubv4.Boolean(slices.Contains(cols, "payload"))
	(*m)["includeDeploymentRef"] = githubv4.Boolean(slices.Contains(cols, "ref"))
	(*m)["includeDeploymentState"] = githubv4.Boolean(slices.Contains(cols, "state"))
	(*m)["includeDeploymentTask"] = githubv4.Boolean(slices.Contains(cols, "task"))
	(*m)["includeDeploymentUpdatedAt"] = githubv4.Boolean(slices.Contains(cols, "updated_at"))
}

func appendRepoEnvironmentColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeEnvironmentName"] = githubv4.Boolean(slices.Contains(cols, "name"))
	(*m)["includeEnvironmentNodeId"] = githubv4.Boolean(slices.Contains(cols, "node_id"))
	(*m)["includeEnvironmentId"] = githubv4.Boolean(slices.Contains(cols, "id"))
}

func appendRepoVulnerabilityAlertColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeVulnerabilityAlertNodeId"] = githubv4.Boolean(slices.Contains(cols, "node_id"))
	(*m)["includeVulnerabilityAlertNumber"] = githubv4.Boolean(slices.Contains(cols, "number"))
	(*m)["includeVulnerabilityAlertAutoDismissedAt"] = githubv4.Boolean(slices.Contains(cols, "auto_dismissed_at"))
	(*m)["includeVulnerabilityAlertCreatedAt"] = githubv4.Boolean(slices.Contains(cols, "created_at"))
	(*m)["includeVulnerabilityAlertDependencyScope"] = githubv4.Boolean(slices.Contains(cols, "dependency_scope"))
	(*m)["includeVulnerabilityAlertDismissComment"] = githubv4.Boolean(slices.Contains(cols, "dismiss_comment"))
	(*m)["includeVulnerabilityAlertDismissReason"] = githubv4.Boolean(slices.Contains(cols, "dismiss_reason"))
	(*m)["includeVulnerabilityAlertDismissedAt"] = githubv4.Boolean(slices.Contains(cols, "dismissed_at"))
	(*m)["includeVulnerabilityAlertDismisser"] = githubv4.Boolean(slices.Contains(cols, "dismisser"))
	(*m)["includeVulnerabilityAlertFixedAt"] = githubv4.Boolean(slices.Contains(cols, "fixed_at"))
	(*m)["includeVulnerabilityAlertState"] = githubv4.Boolean(slices.Contains(cols, "state"))
	(*m)["includeVulnerabilityAlertSecurityAdvisory"] = githubv4.Boolean(slices.Contains(cols, "security_advisory") || slices.Contains(cols, "cvss_score"))
	(*m)["includeVulnerabilityAlertSecurityVulnerability"] = githubv4.Boolean(slices.Contains(cols, "security_vulnerability") || slices.Contains(cols, "severity"))
	(*m)["includeVulnerabilityAlertVulnerableManifestFilename"] = githubv4.Boolean(slices.Contains(cols, "vulnerable_manifest_filename"))
	(*m)["includeVulnerabilityAlertVulnerableManifestPath"] = githubv4.Boolean(slices.Contains(cols, "vulnerable_manifest_path"))
	(*m)["includeVulnerabilityAlertVulnerableRequirements"] = githubv4.Boolean(slices.Contains(cols, "vulnerable_requirements"))
}

func appendStargazerColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeStargazerStarredAt"] = githubv4.Boolean(slices.Contains(cols, "starred_at"))
	(*m)["includeStargazerNode"] = githubv4.Boolean(slices.Contains(cols, "user_login") || slices.Contains(cols, "user_detail"))
}

func appendTagColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeTagTarget"] = githubv4.Boolean(slices.Contains(cols, "tagger_date") || slices.Contains(cols, "tagger_name") || slices.Contains(cols, "tagger_login") || slices.Contains(cols, "message") || slices.Contains(cols, "commit"))
	(*m)["includeTagName"] = githubv4.Boolean(slices.Contains(cols, "name"))
}

func appendUserWithCountColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeUserAnyPinnableItems"] = githubv4.Boolean(slices.Contains(cols, "any_pinnable_items"))
	(*m)["includeUserAvatarUrl"] = githubv4.Boolean(slices.Contains(cols, "avatar_url"))
	(*m)["includeUserBio"] = githubv4.Boolean(slices.Contains(cols, "bio"))
	(*m)["includeUserCompany"] = githubv4.Boolean(slices.Contains(cols, "company"))
	(*m)["includeUserEstimatedNextSponsorsPayoutInCents"] = githubv4.Boolean(slices.Contains(cols, "estimated_next_sponsors_payout_in_cents"))
	(*m)["includeUserHasSponsorsListing"] = githubv4.Boolean(slices.Contains(cols, "has_sponsors_listing"))
	(*m)["includeUserInteractionAbility"] = githubv4.Boolean(slices.Contains(cols, "interaction_ability"))
	(*m)["includeUserIsBountyHunter"] = githubv4.Boolean(slices.Contains(cols, "is_bounty_hunter"))
	(*m)["includeUserIsCampusExpert"] = githubv4.Boolean(slices.Contains(cols, "is_campus_expert"))
	(*m)["includeUserIsDeveloperProgramMember"] = githubv4.Boolean(slices.Contains(cols, "is_developer_program_member"))
	(*m)["includeUserIsEmployee"] = githubv4.Boolean(slices.Contains(cols, "is_employee"))
	(*m)["includeUserIsFollowingYou"] = githubv4.Boolean(slices.Contains(cols, "is_following_you"))
	(*m)["includeUserIsGitHubStar"] = githubv4.Boolean(slices.Contains(cols, "is_github_star"))
	(*m)["includeUserIsHireable"] = githubv4.Boolean(slices.Contains(cols, "is_hireable"))
	(*m)["includeUserIsSiteAdmin"] = githubv4.Boolean(slices.Contains(cols, "is_site_admin"))
	(*m)["includeUserIsSponsoringYou"] = githubv4.Boolean(slices.Contains(cols, "is_sponsoring_you"))
	(*m)["includeUserIsYou"] = githubv4.Boolean(slices.Contains(cols, "is_you"))
	(*m)["includeUserLocation"] = githubv4.Boolean(slices.Contains(cols, "location"))
	(*m)["includeUserMonthlyEstimatedSponsorsIncomeInCents"] = githubv4.Boolean(slices.Contains(cols, "monthly_estimated_sponsors_income_in_cents"))
	(*m)["includeUserPinnedItemsRemaining"] = githubv4.Boolean(slices.Contains(cols, "pinned_items_remaining"))
	(*m)["includeUserProjectsUrl"] = githubv4.Boolean(slices.Contains(cols, "projects_url"))
	(*m)["includeUserPronouns"] = githubv4.Boolean(slices.Contains(cols, "pronouns"))
	(*m)["includeUserSponsorsListing"] = githubv4.Boolean(slices.Contains(cols, "sponsors_listing"))
	(*m)["includeUserStatus"] = githubv4.Boolean(slices.Contains(cols, "status"))
	(*m)["includeUserTwitterUsername"] = githubv4.Boolean(slices.Contains(cols, "twitter_username"))
	(*m)["includeUserCanChangedPinnedItems"] = githubv4.Boolean(slices.Contains(cols, "can_changed_pinned_items"))
	(*m)["includeUserCanCreateProjects"] = githubv4.Boolean(slices.Contains(cols, "can_create_projects"))
	(*m)["includeUserCanFollow"] = githubv4.Boolean(slices.Contains(cols, "can_follow"))
	(*m)["includeUserCanSponsor"] = githubv4.Boolean(slices.Contains(cols, "can_sponsor"))
	(*m)["includeUserIsFollowing"] = githubv4.Boolean(slices.Contains(cols, "is_following"))
	(*m)["includeUserIsSponsoring"] = githubv4.Boolean(slices.Contains(cols, "is_sponsoring"))
	(*m)["includeUserWebsiteUrl"] = githubv4.Boolean(slices.Contains(cols, "website_url"))

	(*m)["includeUserRepositories"] = githubv4.Boolean(slices.Contains(cols, "repositories_total_disk_usage"))
	(*m)["includeUserFollowers"] = githubv4.Boolean(slices.Contains(cols, "followers_total_count"))
	(*m)["includeUserFollowing"] = githubv4.Boolean(slices.Contains(cols, "following_total_count"))
	(*m)["includeUserPublicRepositories"] = githubv4.Boolean(slices.Contains(cols, "public_repositories_total_count"))
	(*m)["includeUserPrivateRepositories"] = githubv4.Boolean(slices.Contains(cols, "private_repositories_total_count"))
	(*m)["includeUserPublicGists"] = githubv4.Boolean(slices.Contains(cols, "public_gists_total_count"))
	(*m)["includeUserIssues"] = githubv4.Boolean(slices.Contains(cols, "issues_total_count"))
	(*m)["includeUserOrganizations"] = githubv4.Boolean(slices.Contains(cols, "organizations_total_count"))
	(*m)["includeUserPublicKeys"] = githubv4.Boolean(slices.Contains(cols, "public_keys_total_count"))
	(*m)["includeUserOpenPullRequests"] = githubv4.Boolean(slices.Contains(cols, "open_pull_requests_total_count"))
	(*m)["includeUserMergedPullRequests"] = githubv4.Boolean(slices.Contains(cols, "merged_pull_requests_total_count"))
	(*m)["includeUserClosedPullRequests"] = githubv4.Boolean(slices.Contains(cols, "closed_pull_requests_total_count"))
	(*m)["includeUserPackages"] = githubv4.Boolean(slices.Contains(cols, "packages_total_count"))
	(*m)["includeUserPinnedItems"] = githubv4.Boolean(slices.Contains(cols, "pinned_items_total_count"))
	(*m)["includeUserSponsoring"] = githubv4.Boolean(slices.Contains(cols, "sponsoring_total_count"))
	(*m)["includeUserSponsors"] = githubv4.Boolean(slices.Contains(cols, "sponsors_total_count"))
	(*m)["includeUserStarredRepositories"] = githubv4.Boolean(slices.Contains(cols, "starred_repositories_total_count"))
	(*m)["includeUserWatching"] = githubv4.Boolean(slices.Contains(cols, "watching_total_count"))
}

func appendPullRequestColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includePRAuthor"] = githubv4.Boolean(slices.Contains(cols, "author"))
	(*m)["includePRBody"] = githubv4.Boolean(slices.Contains(cols, "body"))
	(*m)["includePREditor"] = githubv4.Boolean(slices.Contains(cols, "editor"))
	(*m)["includePRMergedBy"] = githubv4.Boolean(slices.Contains(cols, "merged_by"))
	(*m)["includePRMilestone"] = githubv4.Boolean(slices.Contains(cols, "milestone"))

	(*m)["includePRBaseRef"] = githubv4.Boolean(slices.Contains(cols, "base_ref"))
	(*m)["includePRHeadRef"] = githubv4.Boolean(slices.Contains(cols, "head_ref"))
	(*m)["includePRMergeCommit"] = githubv4.Boolean(slices.Contains(cols, "merge_commit"))
	(*m)["includePRSuggested"] = githubv4.Boolean(slices.Contains(cols, "suggested_reviewers"))
	(*m)["includePRViewer"] = githubv4.Boolean(slices.Contains(cols, "can_apply_suggestion") ||
		slices.Contains(cols, "can_close") ||
		slices.Contains(cols, "can_delete_head_ref") ||
		slices.Contains(cols, "can_disable_auto_merge") ||
		slices.Contains(cols, "can_edit_files") ||
		slices.Contains(cols, "can_enable_auto_merge") ||
		slices.Contains(cols, "can_react") ||
		slices.Contains(cols, "can_reopen") ||
		slices.Contains(cols, "can_subscribe") ||
		slices.Contains(cols, "can_update") ||
		slices.Contains(cols, "can_update_branch") ||
		slices.Contains(cols, "did_author") ||
		slices.Contains(cols, "cannot_update_reasons") ||
		slices.Contains(cols, "subscription"))
	(*m)["includePRAssignees"] = githubv4.Boolean(slices.Contains(cols, "assignees_total_count") || slices.Contains(cols, "assignees"))
	(*m)["includePRCommitCount"] = githubv4.Boolean(slices.Contains(cols, "commits_total_count"))
	(*m)["includePRReviewRequestCount"] = githubv4.Boolean(slices.Contains(cols, "review_requests_total_count"))
	(*m)["includePRReviewCount"] = githubv4.Boolean(slices.Contains(cols, "reviews_total_count"))
	(*m)["includePRLabels"] = githubv4.Boolean(slices.Contains(cols, "labels") ||
		slices.Contains(cols, "labels_src") ||
		slices.Contains(cols, "labels_total_count"))
	(*m)["includePRId"] = githubv4.Boolean(slices.Contains(cols, "id"))
	(*m)["includePRNodeId"] = githubv4.Boolean(slices.Contains(cols, "node_id"))
	(*m)["includePRAuthorAssociation"] = githubv4.Boolean(slices.Contains(cols, "author_association"))
	(*m)["includePRBaseRefName"] = githubv4.Boolean(slices.Contains(cols, "base_ref_name"))
	(*m)["includePRActiveLockReason"] = githubv4.Boolean(slices.Contains(cols, "active_lock_reason"))
	(*m)["includePRAdditions"] = githubv4.Boolean(slices.Contains(cols, "additions"))
	(*m)["includePRChangedFiles"] = githubv4.Boolean(slices.Contains(cols, "changed_files"))
	(*m)["includePRChecksUrl"] = githubv4.Boolean(slices.Contains(cols, "checks_url"))
	(*m)["includePRClosed"] = githubv4.Boolean(slices.Contains(cols, "closed"))
	(*m)["includePRClosedAt"] = githubv4.Boolean(slices.Contains(cols, "closed_at"))
	(*m)["includePRCreatedAt"] = githubv4.Boolean(slices.Contains(cols, "created_at"))
	(*m)["includePRCreatedViaEmail"] = githubv4.Boolean(slices.Contains(cols, "created_via_email"))
	(*m)["includePRDeletions"] = githubv4.Boolean(slices.Contains(cols, "deletions"))
	(*m)["includePRHeadRefName"] = githubv4.Boolean(slices.Contains(cols, "head_ref_name"))
	(*m)["includePRHeadRefOid"] = githubv4.Boolean(slices.Contains(cols, "head_ref_oid"))
	(*m)["includePRIncludesCreatedEdit"] = githubv4.Boolean(slices.Contains(cols, "includes_created_edit"))
	(*m)["includePRIsCrossRepository"] = githubv4.Boolean(slices.Contains(cols, "is_cross_repository"))
	(*m)["includePRIsDraft"] = githubv4.Boolean(slices.Contains(cols, "is_draft"))
	(*m)["includePRIsReadByUser"] = githubv4.Boolean(slices.Contains(cols, "is_read_by_user"))
	(*m)["includePRLastEditedAt"] = githubv4.Boolean(slices.Contains(cols, "last_edited_at"))
	(*m)["includePRLocked"] = githubv4.Boolean(slices.Contains(cols, "locked"))
	(*m)["includePRMaintainerCanModify"] = githubv4.Boolean(slices.Contains(cols, "maintainer_can_modify"))
	(*m)["includePRMergeable"] = githubv4.Boolean(slices.Contains(cols, "mergeable"))
	(*m)["includePRMerged"] = githubv4.Boolean(slices.Contains(cols, "merged"))
	(*m)["includePRMergedAt"] = githubv4.Boolean(slices.Contains(cols, "merged_at"))
	(*m)["includePRPermalink"] = githubv4.Boolean(slices.Contains(cols, "permalink"))
	(*m)["includePRPublishedAt"] = githubv4.Boolean(slices.Contains(cols, "published_at"))
	(*m)["includePRRevertUrl"] = githubv4.Boolean(slices.Contains(cols, "revert_url"))
	(*m)["includePRReviewDecision"] = githubv4.Boolean(slices.Contains(cols, "review_decision"))
	(*m)["includePRState"] = githubv4.Boolean(slices.Contains(cols, "state"))
	(*m)["includePRTitle"] = githubv4.Boolean(slices.Contains(cols, "title"))
	(*m)["includePRTotalCommentsCount"] = githubv4.Boolean(slices.Contains(cols, "total_comments_count"))
	(*m)["includePRUpdatedAt"] = githubv4.Boolean(slices.Contains(cols, "updated_at"))
	(*m)["includePRUrl"] = githubv4.Boolean(slices.Contains(cols, "url"))
}

func repositoryCols() []string {
	return []string{
		"id",
		"node_id",
		"name",
		"allow_update_branch",
		"archived_at",
		"auto_merge_allowed",
		"code_of_conduct",
		"contact_links",
		"created_at",
		"default_branch_ref",
		"delete_branch_on_merge",
		"description",
		"disk_usage",
		"fork_count",
		"forking_allowed",
		"funding_links",
		"has_discussions_enabled",
		"has_issues_enabled",
		"has_projects_enabled",
		"has_vulnerability_alerts_enabled",
		"has_wiki_enabled",
		"homepage_url",
		"interaction_ability",
		"is_archived",
		"is_blank_issues_enabled",
		"is_disabled",
		"is_empty",
		"is_fork",
		"is_in_organization",
		"is_locked",
		"is_mirror",
		"is_private",
		"is_security_policy_enabled",
		"is_template",
		"is_user_configuration_repository",
		"issue_templates",
		"license_info",
		"lock_reason",
		"merge_commit_allowed",
		"merge_commit_message",
		"merge_commit_title",
		"mirror_url",
		"name_with_owner",
		"open_graph_image_url",
		"owner_login",
		"primary_language",
		"projects_url",
		"pull_request_templates",
		"pushed_at",
		"rebase_merge_allowed",
		"security_policy_url",
		"squash_merge_allowed",
		"squash_merge_commit_message",
		"squash_merge_commit_title",
		"ssh_url",
		"stargazer_count",
		"updated_at",
		"url",
		"uses_custom_open_graph_image",
		"can_administer",
		"can_create_projects",
		"can_subscribe",
		"can_update_topics",
		"has_starred",
		"possible_commit_emails",
		"subscription",
		"visibility",
		"your_permission",
		"web_commit_signoff_required",
		"repository_topics_total_count",
		"open_issues_total_count",
		"watchers_total_count",
		"hooks",
		"topics",
		"subscribers_count",
		"has_downloads",
		"has_pages",
		"network_count",
	}
}

func branchCols() []string {
	return []string{
		"repository_full_name",
		"name",
		"commit",
		"protected",
		"branch_protection_rule",
	}
}

func branchProtectionCols() []string {
	return []string{
		"repository_full_name",
		"id",
		"node_id",
		"matching_branches",
		"is_admin_enforced",
		"allows_deletions",
		"allows_force_pushes",
		"blocks_creations",
		"creator_login",
		"dismisses_stale_reviews",
		"lock_allows_fetch_and_merge",
		"lock_branch",
		"pattern",
		"require_last_push_approval",
		"requires_approving_reviews",
		"required_approving_review_count",
		"requires_conversation_resolution",
		"requires_code_owner_reviews",
		"requires_commit_signatures",
		"requires_deployments",
		"required_deployment_environments",
		"requires_linear_history",
		"requires_status_checks",
		"required_status_checks",
		"requires_strict_status_checks",
		"restricts_review_dismissals",
		"restricts_pushes",
		"push_allowance_apps",
		"push_allowance_teams",
		"push_allowance_users",
		"bypass_force_push_allowance_apps",
		"bypass_force_push_allowance_teams",
		"bypass_force_push_allowance_users",
		"bypass_pull_request_allowance_apps",
		"bypass_pull_request_allowance_teams",
		"bypass_pull_request_allowance_users",
		"repository_full_name",
		"name",
		"commit",
		"protected",
		"branch_protection_rule",
	}
}

func commitCols() []string {
	return []string{
		"repository_full_name",
		"sha",
		"short_sha",
		"message",
		"author_login",
		"authored_date",
		"author",
		"committer_login",
		"committed_date",
		"committer",
		"additions",
		"authored_by_committer",
		"deletions",
		"changed_files",
		"committed_via_web",
		"commit_url",
		"signature",
		"status",
		"tarball_url",
		"zipball_url",
		"tree_url",
		"can_subscribe",
		"subscription",
		"url",
		"node_id",
		"message_headline",
	}
}

func organizationCols() []string {
	return []string{
		"login",
		"id",
		"node_id",
		"name",
		"created_at",
		"updated_at",
		"description",
		"email",
		"url",
		"announcement",
		"announcement_expires_at",
		"announcement_user_dismissible",
		"any_pinnable_items",
		"avatar_url",
		"estimated_next_sponsors_payout_in_cents",
		"has_sponsors_listing",
		"interaction_ability",
		"is_sponsoring_you",
		"is_verified",
		"location",
		"monthly_estimated_sponsors_income_in_cents",
		"new_team_url",
		"pinned_items_remaining",
		"projects_url",
		"saml_identity_provider",
		"sponsors_listing",
		"teams_url",
		"total_sponsorship_amount_as_sponsor_in_cents",
		"twitter_username",
		"can_administer",
		"can_changed_pinned_items",
		"can_create_projects",
		"can_create_repositories",
		"can_create_teams",
		"can_sponsor",
		"is_a_member",
		"is_following",
		"is_sponsoring",
		"website_url",
		"hooks",
		"billing_email",
		"two_factor_requirement_enabled",
		"default_repo_permission",
		"members_allowed_repository_creation_type",
		"members_can_create_internal_repos",
		"members_can_create_pages",
		"members_can_create_private_repos",
		"members_can_create_public_repos",
		"members_can_create_repos",
		"members_can_fork_private_repos",
		"plan_filled_seats",
		"plan_name",
		"plan_private_repos",
		"plan_seats",
		"plan_space",
		"followers",
		"following",
		"collaborators",
		"has_organization_projects",
		"has_repository_projects",
		"web_commit_signoff_required",
	}
}

func starCols() []string {
	return []string{
		"repository_full_name",
		"starred_at",
		"url",
	}
}

func issueCols() []string {
	return []string{
		"number",
		"id",
		"node_id",
		"active_lock_reason",
		"author",
		"author_login",
		"author_association",
		"body",
		"body_url",
		"closed",
		"closed_at",
		"created_at",
		"created_via_email",
		"editor",
		"full_database_id",
		"includes_created_edit",
		"is_pinned",
		"is_read_by_user",
		"last_edited_at",
		"locked",
		"milestone",
		"published_at",
		"state",
		"state_reason",
		"title",
		"updated_at",
		"url",
		"assignees_total_count",
		"comments_total_count",
		"labels_total_count",
		"labels_src",
		"labels",
		"user_can_close",
		"user_can_react",
		"user_can_reopen",
		"user_can_subscribe",
		"user_can_update",
		"user_cannot_update_reasons",
		"user_did_author",
		"user_subscription",
		"assignees",
	}
}

func issueCommentCols() []string {
	return []string{
		"repository_full_name",
		"number",
		"id",
		"node_id",
		"author",
		"author_login",
		"author_association",
		"body",
		"body_text",
		"created_at",
		"created_via_email",
		"editor",
		"editor_login",
		"includes_created_edit",
		"is_minimized",
		"minimized_reason",
		"last_edited_at",
		"published_at",
		"updated_at",
		"url",
		"can_delete",
		"can_minimize",
		"can_react",
		"can_update",
		"cannot_update_reasons",
		"did_author",
	}
}

func licenseCols() []string {
	return []string{
		"spdx_id",
		"name",
		"url",
		"conditions",
		"description",
		"featured",
		"hidden",
		"implementation",
		"key",
		"limitations",
		"permissions",
		"nickname",
		"pseudo_license",
	}
}

func teamCols() []string {
	return []string{
		"organization",
		"slug",
		"name",
		"id",
		"node_id",
		"description",
		"created_at",
		"updated_at",
		"combined_slug",
		"parent_team",
		"privacy",
		"ancestors_total_count",
		"child_teams_total_count",
		"discussions_total_count",
		"invitations_total_count",
		"members_total_count",
		"projects_v2_total_count",
		"repositories_total_count",
		"url",
		"avatar_url",
		"discussions_url",
		"edit_team_url",
		"members_url",
		"new_team_url",
		"repositories_url",
		"teams_url",
		"can_administer",
		"can_subscribe",
		"subscription",
	}
}

func orgCollaboratorsCols() []string {
	return []string{
		"organization",
		"affiliation",
		"repository_name",
		"permission",
		"user_login",
	}
}

func orgExternalIdentitiesCols() []string {
	return []string{
		"organization",
		"guid",
		"user_login",
		"user_detail",
		"saml_identity",
		"scim_identity",
		"organization_invitation",
	}
}

func orgMembersCols() []string {
	return []string{
		"organization",
		"role",
		"has_two_factor_enabled",
		"login",
		"id",
		"name",
		"node_id",
		"email",
		"url",
		"created_at",
		"updated_at",
		"any_pinnable_items",
		"avatar_url",
		"bio",
		"company",
		"estimated_next_sponsors_payout_in_cents",
		"has_sponsors_listing",
		"interaction_ability",
		"is_bounty_hunter",
		"is_campus_expert",
		"is_developer_program_member",
		"is_employee",
		"is_following_you",
		"is_github_star",
		"is_hireable",
		"is_site_admin",
		"is_sponsoring_you",
		"is_you",
		"location",
		"monthly_estimated_sponsors_income_in_cents",
		"pinned_items_remaining",
		"projects_url",
		"pronouns",
		"sponsors_listing",
		"status",
		"twitter_username",
		"can_changed_pinned_items",
		"can_create_projects",
		"can_follow",
		"can_sponsor",
		"is_following",
		"is_sponsoring",
		"website_url",
	}
}

func repositoryCollaboratorsCols() []string {
	return []string{
		"repository_full_name",
		"affiliation",
		"permission",
		"user_login",
	}
}

func repositoryDeploymentsCols() []string {
	return []string{
		"repository_full_name",
		"id",
		"node_id",
		"commit_sha",
		"created_at",
		"creator",
		"description",
		"environment",
		"latest_environment",
		"latest_status",
		"original_environment",
		"payload",
		"ref",
		"state",
		"task",
		"updated_at",
	}
}

func repositoryEnvironmentsCols() []string {
	return []string{
		"repository_full_name",
		"id",
		"node_id",
		"name",
	}
}

func repositoryVulnerabilityAlertCols() []string {
	return []string{
		"repository_full_name",
		"number",
		"node_id",
		"auto_dismissed_at",
		"created_at",
		"dependency_scope",
		"dismiss_comment",
		"dismiss_reason",
		"dismissed_at",
		"dismisser",
		"fixed_at",
		"state",
		"security_advisory",
		"security_vulnerability",
		"vulnerable_manifest_filename",
		"vulnerable_manifest_path",
		"vulnerable_requirements",
		"severity",
		"cvss_score",
	}
}

func teamMembersCols() []string {
	return []string{
		"organization",
		"slug",
		"name",
		"id",
		"node_id",
		"description",
		"created_at",
		"updated_at",
		"combined_slug",
		"parent_team",
		"privacy",
		"ancestors_total_count",
		"child_teams_total_count",
		"discussions_total_count",
		"invitations_total_count",
		"members_total_count",
		"projects_v2_total_count",
		"repositories_total_count",
		"url",
		"avatar_url",
		"discussions_url",
		"edit_team_url",
		"members_url",
		"new_team_url",
		"repositories_url",
		"teams_url",
		"can_administer",
		"can_subscribe",
		"subscription",
	}
}

func teamRepositoriesCols() []string {
	return []string{
		"organization",
		"slug",
		"permission",
		"id",
		"node_id",
		"name",
		"allow_update_branch",
		"archived_at",
		"auto_merge_allowed",
		"code_of_conduct",
		"contact_links",
		"created_at",
		"default_branch_ref",
		"delete_branch_on_merge",
		"description",
		"disk_usage",
		"fork_count",
		"forking_allowed",
		"funding_links",
		"has_discussions_enabled",
		"has_issues_enabled",
		"has_projects_enabled",
		"has_vulnerability_alerts_enabled",
		"has_wiki_enabled",
		"homepage_url",
		"interaction_ability",
		"is_archived",
		"is_blank_issues_enabled",
		"is_disabled",
		"is_empty",
		"is_fork",
		"is_in_organization",
		"is_locked",
		"is_mirror",
		"is_private",
		"is_security_policy_enabled",
		"is_template",
		"is_user_configuration_repository",
		"issue_templates",
		"license_info",
		"lock_reason",
		"merge_commit_allowed",
		"merge_commit_message",
		"merge_commit_title",
		"mirror_url",
		"name_with_owner",
		"open_graph_image_url",
		"owner_login",
		"primary_language",
		"projects_url",
		"pull_request_templates",
		"pushed_at",
		"rebase_merge_allowed",
		"security_policy_url",
		"squash_merge_allowed",
		"squash_merge_commit_message",
		"squash_merge_commit_title",
		"ssh_url",
		"stargazer_count",
		"updated_at",
		"url",
		"uses_custom_open_graph_image",
		"can_administer",
		"can_create_projects",
		"can_subscribe",
		"can_update_topics",
		"has_starred",
		"possible_commit_emails",
		"subscription",
	}
}

func stargazerCols() []string {
	return []string{
		"repository_full_name",
		"starred_at",
		"user_login",
		"user_detail",
	}
}

func tagCols() []string {
	return []string{
		"repository_full_name",
		"name",
		"tagger_date",
		"tagger_name",
		"tagger_login",
		"message",
		"commit",
	}
}

func userCols() []string {
	return []string{
		"login",
		"id",
		"name",
		"node_id",
		"email",
		"url",
		"created_at",
		"updated_at",
		"any_pinnable_items",
		"avatar_url",
		"bio",
		"company",
		"estimated_next_sponsors_payout_in_cents",
		"has_sponsors_listing",
		"interaction_ability",
		"is_bounty_hunter",
		"is_campus_expert",
		"is_developer_program_member",
		"is_employee",
		"is_following_you",
		"is_github_star",
		"is_hireable",
		"is_site_admin",
		"is_sponsoring_you",
		"is_you",
		"location",
		"monthly_estimated_sponsors_income_in_cents",
		"pinned_items_remaining",
		"projects_url",
		"pronouns",
		"sponsors_listing",
		"status",
		"twitter_username",
		"can_changed_pinned_items",
		"can_create_projects",
		"can_follow",
		"can_sponsor",
		"is_following",
		"is_sponsoring",
		"website_url",
		"repositories_total_disk_usage",
		"followers_total_count",
		"following_total_count",
		"public_repositories_total_count",
		"private_repositories_total_count",
		"public_gists_total_count",
		"issues_total_count",
		"organizations_total_count",
		"public_keys_total_count",
		"open_pull_requests_total_count",
		"merged_pull_requests_total_count",
		"closed_pull_requests_total_count",
		"packages_total_count",
		"pinned_items_total_count",
		"sponsoring_total_count",
		"sponsors_total_count",
		"starred_repositories_total_count",
		"watching_total_count",
	}
}

func pullRequestCols() []string {
	return []string{
		"repository_full_name",
		"number",
		"id",
		"node_id",
		"active_lock_reason",
		"additions",
		"author",
		"author_association",
		"base_ref_name",
		"body",
		"changed_files",
		"checks_url",
		"closed",
		"closed_at",
		"created_at",
		"created_via_email",
		"deletions",
		"editor",
		"head_ref_name",
		"head_ref_oid",
		"includes_created_edit",
		"is_cross_repository",
		"is_draft",
		"is_read_by_user",
		"last_edited_at",
		"locked",
		"maintainer_can_modify",
		"mergeable",
		"merged",
		"merged_at",
		"merged_by",
		"milestone",
		"permalink",
		"published_at",
		"revert_url",
		"review_decision",
		"state",
		"title",
		"total_comments_count",
		"updated_at",
		"url",
		"assignees",
		"base_ref",
		"head_ref",
		"merge_commit",
		"suggested_reviewers",
		"can_apply_suggestion",
		"can_close",
		"can_delete_head_ref",
		"can_disable_auto_merge",
		"can_edit_files",
		"can_enable_auto_merge",
		"can_merge_as_admin",
		"can_react",
		"can_reopen",
		"can_subscribe",
		"can_update",
		"can_update_branch",
		"did_author",
		"cannot_update_reasons",
		"subscription",
		"labels_src",
		"labels",
		"assignees_total_count",
		"labels_total_count",
		"commits_total_count",
		"review_requests_total_count",
		"reviews_total_count",
	}
}

func getRepositories(ctx context.Context, client *github.Client, owner string) ([]*github.Repository, error) {
	var repositories []*github.Repository
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: maxPagesCount},
	}
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, owner, opt)
		if err != nil {
			return nil, err
		}

		for _, r := range repos {
			if r.GetArchived() || r.GetDisabled() {
				continue
			}
			repositories = append(repositories, r)
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return repositories, nil
}

func getIssues(ctx context.Context, orgName string, client *github.Client) ([]*github.Issue, error) {
	opt := &github.IssueListOptions{
		Filter:      "assigned",
		State:       "all",
		ListOptions: github.ListOptions{PerPage: issuePageSize},
	}
	for {
		issues, resp, err := client.Issues.ListByOrg(ctx, orgName, opt)
		if err != nil {
			return nil, err
		}
		if resp.NextPage == 0 {
			return issues, nil
		}
		opt.Page = resp.NextPage
	}
}

func getOrganizations(ctx context.Context, client *github.Client) ([]*github.Organization, error) {
	opts := &github.ListOptions{PerPage: orgCollaboratorsPageSize}
	for {
		organizations, resp, err := client.Organizations.List(ctx, "", opts)
		if err != nil {
			return nil, err
		}
		if resp.NextPage == 0 {
			return organizations, nil
		}
		opts.Page = resp.NextPage
	}
}

func getTeams(ctx context.Context, client *github.Client) ([]*github.Team, error) {
	var allTeams []*github.Team
	opt := &github.ListOptions{PerPage: 10}
	for {
		teams, resp, err := client.Teams.ListUserTeams(ctx, opt)
		if err != nil {
			return nil, err
		}
		allTeams = append(allTeams, teams...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allTeams, nil
}

func getFileSHAs(client *github.Client, owner, repo string) ([]string, error) {
	fileContent, directoryContent, _, err := client.Repositories.GetContents(context.Background(), owner, repo, "", nil)
	if err != nil {
		return nil, err
	}
	var fileSHAs []string
	if fileContent != nil {
		fileSHAs = append(fileSHAs, fileContent.GetSHA())
	} else {
		for _, content := range directoryContent {
			fileSHAs = append(fileSHAs, content.GetSHA())
		}
	}
	return fileSHAs, nil
}

func getPackages(ctx context.Context, githubClient model.GitHubClient, organizationName string) ([]*github.Package, error) {
	client := githubClient.RestClient
	var packages []*github.Package
	packageTypes := []string{"container", "maven", "npm", "rubygems", "nuget"}
	for _, packageType := range packageTypes {
		page := 1
		for {
			var opts = &github.PackageListOptions{
				PackageType: &packageType,
				ListOptions: github.ListOptions{
					Page:    page,
					PerPage: packagePageSize,
				},
			}
			respPackages, resp, err := client.Organizations.ListPackages(ctx, organizationName, opts)
			if err != nil {
				return nil, err
			}
			packages = append(packages, respPackages...)
			if resp.After == "" {
				break
			}
			opts.ListOptions.Page += 1
		}
	}
	return packages, nil
}

func fetchAllPackages(sdk *resilientbridge.ResilientBridge, org, packageType string) ([]model.PackageListItem, error) {
	var allPackages []model.PackageListItem
	page := 1
	perPage := 100

	for {
		endpoint := fmt.Sprintf("/orgs/%s/packages?package_type=%s&page=%d&per_page=%d", org, packageType, page, perPage)
		listReq := &resilientbridge.NormalizedRequest{
			Method:   "GET",
			Endpoint: endpoint,
			Headers:  map[string]string{"Accept": "application/vnd.github+json"},
		}

		listResp, err := sdk.Request("github", listReq)
		if err != nil {
			return nil, err
		}
		if listResp.StatusCode >= 400 {
			return nil, err
		}

		var packages []model.PackageListItem
		if err := json.Unmarshal(listResp.Data, &packages); err != nil {
			return nil, err
		}

		if len(packages) == 0 {
			break
		}

		allPackages = append(allPackages, packages...)
		page++
	}
	return allPackages, nil
}

func fetchPackageDetails(sdk *resilientbridge.ResilientBridge, org, packageType, packageName string, stream *models.StreamSender) ([]models.Resource, error) {
	var pd model.PackageDetailDescription
	endpoint := fmt.Sprintf("/orgs/%s/packages/%s/%s", org, packageType, packageName)
	req := &resilientbridge.NormalizedRequest{
		Method:   "GET",
		Endpoint: endpoint,
		Headers:  map[string]string{"Accept": "application/vnd.github+json"},
	}

	resp, err := sdk.Request("github", req)
	if err != nil {
		return nil, fmt.Errorf("error fetching package details: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(resp.Data))
	}

	if err := json.Unmarshal(resp.Data, &pd); err != nil {
		return nil, fmt.Errorf("error parsing package details: %w", err)
	}

	var values []models.Resource

	value := models.Resource{
		ID:          strconv.Itoa(pd.ID),
		Name:        pd.Name,
		Description: pd,
	}
	if stream != nil {
		if err := (*stream)(value); err != nil {
			return nil, err
		}
	} else {
		values = append(values, value)
	}

	return values, nil
}

func formRepositoryFullName(owner, repo string) string {
	return fmt.Sprintf("%s/%s", owner, repo)
}

func parseRepoFullName(fullName string) (string, string) {
	owner := ""
	repo := ""
	s := strings.Split(fullName, "/")
	owner = s[0]
	if len(s) > 1 {
		repo = s[1]
	}
	return owner, repo
}
