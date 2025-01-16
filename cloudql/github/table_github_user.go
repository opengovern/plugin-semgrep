package github

import (
	opengovernance "github.com/opengovern/og-describer-template/discovery/pkg/es"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubUser() *plugin.Table {
	return &plugin.Table{
		Name:        "github_user",
		Description: "GitHub Users are user accounts in GitHub.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListUser,
		},
		Columns: commonColumns(tableGitHubUserColumns()),
	}
}

func tableGitHubUserColumns() []*plugin.Column {
	cols := sharedUserColumns()

	counts := []*plugin.Column{
		{Name: "repositories_total_disk_usage", Type: proto.ColumnType_INT, Description: "Total disk spaced used by the users repositories.",
			Transform: transform.FromField("Description.RepositoriesTotalDiskUsage")},
		{Name: "followers_total_count", Type: proto.ColumnType_INT, Description: "Count of how many users this user follows.",
			Transform: transform.FromField("Description.FollowersTotalCount")},
		{Name: "following_total_count", Type: proto.ColumnType_INT, Description: "Count of how many users follow this user.",
			Transform: transform.FromField("Description.FollowingTotalCount")},
		{Name: "public_repositories_total_count", Type: proto.ColumnType_INT, Description: "Count of public repositories for the user.",
			Transform: transform.FromField("Description.PublicRepositoriesTotalCount")},
		{Name: "private_repositories_total_count", Type: proto.ColumnType_INT, Description: "Count of private repositories for the user.",
			Transform: transform.FromField("Description.PrivateRepositoriesTotalCount")},
		{Name: "public_gists_total_count", Type: proto.ColumnType_INT, Description: "Count of public gists for the user.",
			Transform: transform.FromField("Description.PublicGistsTotalCount")},
		{Name: "issues_total_count", Type: proto.ColumnType_INT, Description: "Count of issues associated with the user.",
			Transform: transform.FromField("Description.IssuesTotalCount")},
		{Name: "organizations_total_count", Type: proto.ColumnType_INT, Description: "Count of organizations the user belongs to.",
			Transform: transform.FromField("Description.OrganizationsTotalCount")},
		{Name: "public_keys_total_count", Type: proto.ColumnType_INT, Description: "Count of public keys associated with the user.",
			Transform: transform.FromField("Description.PublicKeysTotalCount")},
		{Name: "open_pull_requests_total_count", Type: proto.ColumnType_INT, Description: "Count of open pull requests associated with the user.",
			Transform: transform.FromField("Description.OpenPullRequestsTotalCount")},
		{Name: "merged_pull_requests_total_count", Type: proto.ColumnType_INT, Description: "Count of merged pull requests associated with the user.",
			Transform: transform.FromField("Description.MergedPullRequestsTotalCount")},
		{Name: "closed_pull_requests_total_count", Type: proto.ColumnType_INT, Description: "Count of closed pull requests associated with the user.",
			Transform: transform.FromField("Description.ClosedPullRequestsTotalCount")},
		{Name: "packages_total_count", Type: proto.ColumnType_INT, Description: "Count of packages hosted by the user.",
			Transform: transform.FromField("Description.PackagesTotalCount")},
		{Name: "pinned_items_total_count", Type: proto.ColumnType_INT, Description: "Count of items pinned on the users profile.",
			Transform: transform.FromField("Description.PinnedItemsTotalCount")},
		{Name: "sponsoring_total_count", Type: proto.ColumnType_INT, Description: "Count of users that this user is sponsoring.",
			Transform: transform.FromField("Description.SponsoringTotalCount")},
		{Name: "sponsors_total_count", Type: proto.ColumnType_INT, Description: "Count of users sponsoring this user.",
			Transform: transform.FromField("Description.SponsorsTotalCount")},
		{Name: "starred_repositories_total_count", Type: proto.ColumnType_INT, Description: "Count of repositories the user has starred.",
			Transform: transform.FromField("Description.StarredRepositoriesTotalCount")},
		{Name: "watching_total_count", Type: proto.ColumnType_INT, Description: "Count of repositories being watched by the user.",
			Transform: transform.FromField("Description.WatchingTotalCount")},
	}

	cols = append(cols, counts...)

	return cols
}

func sharedUserColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "login", Type: proto.ColumnType_STRING, Description: "The login name of the user.",
			Transform: transform.FromField("Description.Login")},
		{Name: "id", Type: proto.ColumnType_INT, Description: "The ID of the user.",
			Transform: transform.FromField("Description.Id")},
		{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the user.",
			Transform: transform.FromField("Description.Name")},
		{Name: "node_id", Type: proto.ColumnType_STRING, Description: "The node ID of the user.",
			Transform: transform.FromField("Description.NodeId")},
		{Name: "email", Type: proto.ColumnType_STRING, Description: "The email of the user.",
			Transform: transform.FromField("Description.Email")},
		{Name: "url", Type: proto.ColumnType_STRING, Description: "The URL of the user's GitHub page.",
			Transform: transform.FromField("Description.Url")},
		{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when user was created.",
			Transform: transform.FromField("Description.CreatedAt").NullIfZero().Transform(convertTimestamp)},
		{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when user was last updated.",
			Transform: transform.FromField("Description.UpdatedAt").NullIfZero().Transform(convertTimestamp)},
		{Name: "any_pinnable_items", Type: proto.ColumnType_BOOL, Description: "If true, user has pinnable items.",
			Transform: transform.FromField("Description.AnyPinnableItems")},
		{Name: "avatar_url", Type: proto.ColumnType_STRING, Description: "The URL of the user's avatar.",
			Transform: transform.FromField("Description.AvatarUrl")},
		{Name: "bio", Type: proto.ColumnType_STRING, Description: "The biography of the user.",
			Transform: transform.FromField("Description.Bio")},
		{Name: "company", Type: proto.ColumnType_STRING, Description: "The company on the users profile.",
			Transform: transform.FromField("Description.Company")},
		{Name: "estimated_next_sponsors_payout_in_cents", Type: proto.ColumnType_INT, Description: "The estimated next GitHub sponsors payout for this user in cents (USD).",
			Transform: transform.FromField("Description.EstimatedNextSponsorsPayoutInCents")},
		{Name: "has_sponsors_listing", Type: proto.ColumnType_BOOL, Description: "If true, user has a GitHub sponsors listing.",
			Transform: transform.FromField("Description.HasSponsorsListing")},
		{Name: "interaction_ability", Type: proto.ColumnType_JSON, Description: "The interaction ability settings for this user.",
			Transform: transform.FromField("Description.InteractionAbility")},
		{Name: "is_bounty_hunter", Type: proto.ColumnType_BOOL, Description: "If true, user is a participant in the GitHub security bug bounty.",
			Transform: transform.FromField("Description.IsBountyHunter")},
		{Name: "is_campus_expert", Type: proto.ColumnType_BOOL, Description: "If true, user is a participant in the GitHub campus experts program.",
			Transform: transform.FromField("Description.IsCampusExpert")},
		{Name: "is_developer_program_member", Type: proto.ColumnType_BOOL, Description: "If true, user is a GitHub developer program member.",
			Transform: transform.FromField("Description.IsDeveloperProgramMember")},
		{Name: "is_employee", Type: proto.ColumnType_BOOL, Description: "If true, user is a GitHub employee.",
			Transform: transform.FromField("Description.IsEmployee")},
		{Name: "is_following_you", Type: proto.ColumnType_BOOL, Description: "If true, user follows you.",
			Transform: transform.FromField("Description.IsFollowingYou")},
		{Name: "is_github_star", Type: proto.ColumnType_BOOL, Description: "If true, user is a member of the GitHub Stars Program.",
			Transform: transform.FromField("Description.IsGitHubStar")},
		{Name: "is_hireable", Type: proto.ColumnType_BOOL, Description: "If true, user has marked themselves as for hire.",
			Transform: transform.FromField("Description.IsHireable")},
		{Name: "is_site_admin", Type: proto.ColumnType_BOOL, Description: "If true, user is a site administrator.",
			Transform: transform.FromField("Description.IsSiteAdmin")},
		{Name: "is_sponsoring_you", Type: proto.ColumnType_BOOL, Description: "If true, this user is sponsoring you.",
			Transform: transform.FromField("Description.IsSponsoringYou")},
		{Name: "is_you", Type: proto.ColumnType_BOOL, Description: "If true, user is you.",
			Transform: transform.FromField("Description.IsYou")},
		{Name: "location", Type: proto.ColumnType_STRING, Description: "The location of the user.",
			Transform: transform.FromField("Description.Location")},
		{Name: "monthly_estimated_sponsors_income_in_cents", Type: proto.ColumnType_INT, Description: "The estimated monthly GitHub sponsors income for this user in cents (USD).",
			Transform: transform.FromField("Description.MonthlyEstimatedSponsorsIncomeInCents")},
		{Name: "pinned_items_remaining", Type: proto.ColumnType_INT, Description: "How many more items this user can pin to their profile.",
			Transform: transform.FromField("Description.PinnedItemsRemaining")},
		{Name: "projects_url", Type: proto.ColumnType_STRING, Description: "The URL listing user's projects.",
			Transform: transform.FromField("Description.ProjectsUrl")},
		{Name: "pronouns", Type: proto.ColumnType_STRING, Description: "The user's pronouns.",
			Transform: transform.FromField("Description.Pronouns")},
		{Name: "sponsors_listing", Type: proto.ColumnType_JSON, Description: "The GitHub sponsors listing for this user.",
			Transform: transform.FromField("Description.SponsorsListing")},
		{Name: "status", Type: proto.ColumnType_JSON, Description: "The user's status.",
			Transform: transform.FromField("Description.Status")},
		{Name: "twitter_username", Type: proto.ColumnType_STRING, Description: "Twitter username of the user.",
			Transform: transform.FromField("Description.TwitterUsername")},
		{Name: "can_changed_pinned_items", Type: proto.ColumnType_BOOL, Description: "If true, you can change the pinned items for this user.",
			Transform: transform.FromField("Description.CanChangedPinnedItems")},
		{Name: "can_create_projects", Type: proto.ColumnType_BOOL, Description: "If true, you can create projects for this user.",
			Transform: transform.FromField("Description.CanCreateProjects")},
		{Name: "can_follow", Type: proto.ColumnType_BOOL, Description: "If true, you can follow this user.",
			Transform: transform.FromField("Description.CanFollow")},
		{Name: "can_sponsor", Type: proto.ColumnType_BOOL, Description: "If true, you can sponsor this user.",
			Transform: transform.FromField("Description.CanSponsor")},
		{Name: "is_following", Type: proto.ColumnType_BOOL, Description: "If true, you are following this user.",
			Transform: transform.FromField("Description.IsFollowing")},
		{Name: "is_sponsoring", Type: proto.ColumnType_BOOL, Description: "If true, you are sponsoring this user.",
			Transform: transform.FromField("Description.IsSponsoring")},
		{Name: "website_url", Type: proto.ColumnType_STRING, Description: "The URL pointing to the user's public website/blog.",
			Transform: transform.FromField("Description.WebsiteUrl")},
	}
}
