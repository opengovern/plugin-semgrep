package github

import (
	opengovernance "github.com/opengovern/og-describer-template/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func sharedOrganizationColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name: "login", Type: proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.Login"),
			Description: "The login name of the organization."},
		{
			Name: "id", Type: proto.ColumnType_INT,
			Transform:   transform.FromField("Description.Id"),
			Description: "The ID number of the organization."},
		{
			Name: "node_id", Type: proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.NodeId"),
			Description: "The node ID of the organization."},
		{
			Name: "name", Type: proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.Name"),
			Description: "The display name of the organization."},
		{
			Name: "created_at", Type: proto.ColumnType_TIMESTAMP,
			Transform:   transform.FromField("Description.CreatedAt").NullIfZero().Transform(convertTimestamp),
			Description: "Timestamp when the organization was created."},
		{
			Name: "updated_at", Type: proto.ColumnType_TIMESTAMP,
			Transform:   transform.FromField("Description.UpdatedAt").NullIfZero().Transform(convertTimestamp),
			Description: "Timestamp when the organization was last updated."},
		{
			Name: "description", Type: proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.Description"),
			Description: "The description of the organization."},
		{
			Name: "email", Type: proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.Email"),
			Description: "The email address associated with the organization."},
		{
			Name: "url", Type: proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.Url"),
			Description: "The URL for this organization."},
		{
			Name: "announcement", Type: proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.Announcement"),
			Description: "The text of the announcement."},
		{
			Name: "announcement_expires_at", Type: proto.ColumnType_TIMESTAMP,
			Transform:   transform.FromField("Description.AnnouncementExpiresAt").NullIfZero().Transform(convertTimestamp),
			Description: "The expiration date of the announcement, if any."},
		{
			Name: "announcement_user_dismissible", Type: proto.ColumnType_BOOL,
			Transform:   transform.FromField("Description.AnnouncementUserDismissible"),
			Description: "If true, the announcement can be dismissed by the user."},
		{
			Name: "any_pinnable_items", Type: proto.ColumnType_BOOL,
			Transform:   transform.FromField("Description.AnyPinnableItems"),
			Description: "If true, this organization has items that can be pinned to their profile."},
		{
			Name: "avatar_url", Type: proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.AvatarUrl"),
			Description: "URL pointing to the organization's public avatar."},
		{
			Name: "estimated_next_sponsors_payout_in_cents", Type: proto.ColumnType_INT,
			Transform:   transform.FromField("Description.EstimatedNextSponsorsPayoutInCents"),
			Description: "The estimated next GitHub Sponsors payout for this organization in cents (USD)."},
		{
			Name: "has_sponsors_listing", Type: proto.ColumnType_BOOL,
			Transform:   transform.FromField("Description.HasSponsorsListing"),
			Description: "If true, this organization has a GitHub Sponsors listing."},
		{
			Name: "interaction_ability", Type: proto.ColumnType_JSON,
			Transform:   transform.FromField("Description.InteractionAbility"),
			Description: "The interaction ability settings for this organization."},
		{
			Name: "is_sponsoring_you", Type: proto.ColumnType_BOOL,
			Transform:   transform.FromField("Description.IsSponsoringYou"),
			Description: "If true, you are sponsored by this organization."},
		{
			Name: "is_verified", Type: proto.ColumnType_BOOL,
			Transform:   transform.FromField("Description.IsVerified"),
			Description: "If true, the organization has verified its profile email and website."},
		{Name: "location", Type: proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.Location"),
			Description: "The organization's public profile location."},
		{
			Name: "monthly_estimated_sponsors_income_in_cents", Type: proto.ColumnType_INT,
			Transform:   transform.FromField("Description.MonthlyEstimatedSponsorsIncomeInCents"),
			Description: "The estimated monthly GitHub Sponsors income for this organization in cents (USD)."},
		{
			Name: "new_team_url", Type: proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.NewTeamUrl"),
			Description: "URL for creating a new team."},
		{Name: "pinned_items_remaining", Type: proto.ColumnType_INT,
			Transform:   transform.FromField("Description.PinnedItemsRemaining"),
			Description: "Returns how many more items this organization can pin to their profile."},
		{
			Name: "projects_url", Type: proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.ProjectsUrl"),
			Description: "URL listing organization's projects."},
		{
			Name: "saml_identity_provider", Type: proto.ColumnType_JSON,
			Transform:   transform.FromField("Description.SamlIdentityProvider"),
			Description: "The Organization's SAML identity provider. Visible to (1) organization owners, (2) organization owners' personal access tokens (classic) with read:org or admin:org scope, (3) GitHub App with an installation token with read or write access to members, else null."},
		{
			Name: "sponsors_listing", Type: proto.ColumnType_JSON,
			Transform:   transform.FromField("Description.SponsorsListing"),
			Description: "The GitHub sponsors listing for this organization."},
		{
			Name: "teams_url", Type: proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.TeamsUrl"),
			Description: "URL listing organization's teams."},
		{
			Name: "total_sponsorship_amount_as_sponsor_in_cents", Type: proto.ColumnType_INT,
			Transform:   transform.FromField("Description.TotalSponsorshipAmountAsSponsorInCents"),
			Description: "The amount in United States cents (e.g., 500 = $5.00 USD) that this entity has spent on GitHub to fund sponsorships. Only returns a value when viewed by the user themselves or by a user who can manage sponsorships for the requested organization."},
		{
			Name: "twitter_username", Type: proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.TwitterUsername"),
			Description: "The organization's Twitter username."},
		{
			Name: "can_administer", Type: proto.ColumnType_BOOL,
			Transform:   transform.FromField("Description.CanAdminister"),
			Description: "If true, you can administer the organization."},
		{
			Name: "can_changed_pinned_items", Type: proto.ColumnType_BOOL,
			Transform: transform.FromField("Description.CanChangedPinnedItems"),

			Description: "If true, you can change the pinned items on the organization's profile."},
		{
			Name: "can_create_projects", Type: proto.ColumnType_BOOL,
			Transform: transform.FromField("Description.CanCreateProjects"),

			Description: "If true, you can create projects for the organization."},
		{
			Name: "can_create_repositories", Type: proto.ColumnType_BOOL,
			Transform: transform.FromField("Description.CanCreateRepositories"),

			Description: "If true, you can create repositories for the organization."},
		{
			Name: "can_create_teams", Type: proto.ColumnType_BOOL,
			Transform: transform.FromField("Description.CanCreateTeams"),

			Description: "If true, you can create teams within the organization."},
		{
			Name: "can_sponsor", Type: proto.ColumnType_BOOL,
			Transform:   transform.FromField("Description.CanSponsor"),
			Description: "If true, you can sponsor this organization."},
		{
			Name: "is_a_member", Type: proto.ColumnType_BOOL,
			Transform: transform.FromField("Description.IsAMember"),

			Description: "If true, you are an active member of the organization."},
		{
			Name: "is_following", Type: proto.ColumnType_BOOL,
			Transform: transform.FromField("Description.IsFollowing"),

			Description: "If true, you are following the organization."},
		{
			Name: "is_sponsoring", Type: proto.ColumnType_BOOL,
			Transform: transform.FromField("Description.IsSponsoring"),

			Description: "If true, you are sponsoring the organization."},
		{
			Name: "website_url", Type: proto.ColumnType_STRING,
			Transform: transform.FromField("Description.WebsiteUrl"),

			Description: "URL for the organization's public website."},
		// Columns from v3 api - hydrates
		{
			Name: "hooks", Type: proto.ColumnType_JSON,
			Description: "The Hooks of the organization.",
			Transform:   transform.FromField("Description.Hooks"),
		},
		{
			Name: "billing_email", Type: proto.ColumnType_STRING,
			Description: "The email address for billing.",
			Transform:   transform.FromField("Description.BillingEmail")},
		{
			Name:        "two_factor_requirement_enabled",
			Type:        proto.ColumnType_BOOL,
			Description: "If true, all members in the organization must have two factor authentication enabled.",
			Transform:   transform.FromField("Description.TwoFactorRequirementEnabled")},
		{
			Name:        "default_repo_permission",
			Type:        proto.ColumnType_STRING,
			Description: "The default repository permissions for the organization.",
			Transform:   transform.FromField("Description.DefaultRepoPermission")},
		{
			Name:        "members_allowed_repository_creation_type",
			Type:        proto.ColumnType_STRING,
			Description: "Specifies which types of repositories non-admin organization members can create",
			Transform:   transform.FromField("Description.MembersAllowedRepositoryCreationType")},
		{
			Name:        "members_can_create_internal_repos",
			Type:        proto.ColumnType_BOOL,
			Description: "If true, members can create internal repositories.",
			Transform:   transform.FromField("Description.MembersCanCreateInternalRepos")},
		{
			Name:        "members_can_create_pages",
			Type:        proto.ColumnType_BOOL,
			Description: "If true, members can create pages.",
			Transform:   transform.FromField("Description.MembersCanCreatePages")},
		{
			Name:        "members_can_create_private_repos",
			Type:        proto.ColumnType_BOOL,
			Description: "If true, members can create private repositories.",
			Transform:   transform.FromField("Description.MembersCanCreatePrivateRepos")},
		{
			Name:        "members_can_create_public_repos",
			Type:        proto.ColumnType_BOOL,
			Description: "If true, members can create public repositories.",
			Transform:   transform.FromField("Description.MembersCanCreatePublicRepos")},
		{
			Name:        "members_can_create_repos",
			Type:        proto.ColumnType_BOOL,
			Description: "If true, members can create repositories.",
			Transform:   transform.FromField("Description.MembersCanCreateRepos")},
		{
			Name:        "members_can_fork_private_repos",
			Type:        proto.ColumnType_BOOL,
			Description: "If true, members can fork private organization repositories.",
			Transform:   transform.FromField("Description.MembersCanForkPrivateRepos")},
		{
			Name:        "plan_filled_seats",
			Type:        proto.ColumnType_INT,
			Description: "The number of used seats for the plan.",
			Transform:   transform.FromField("Description.PlanFilledSeats")},

		{
			Name:        "plan_name",
			Type:        proto.ColumnType_STRING,
			Description: "The name of the GitHub plan.",
			Transform:   transform.FromField("Description.PlanName")},

		{
			Name:        "plan_private_repos",
			Type:        proto.ColumnType_INT,
			Description: "The number of private repositories for the plan.",
			Transform:   transform.FromField("Description.PlanPrivateRepos")},

		{
			Name:        "plan_seats",
			Type:        proto.ColumnType_INT,
			Description: "The number of available seats for the plan",
			Transform:   transform.FromField("Description.PlanSeats")},

		{
			Name:        "plan_space",
			Type:        proto.ColumnType_INT,
			Description: "The total space allocated for the plan.",
			Transform:   transform.FromField("Description.PlanSpace")},

		{
			Name:      "followers",
			Type:      proto.ColumnType_INT,
			Transform: transform.FromField("Description.Followers")},

		{
			Name:      "following",
			Type:      proto.ColumnType_INT,
			Transform: transform.FromField("Description.Following")},

		{
			Name:      "collaborators",
			Type:      proto.ColumnType_INT,
			Transform: transform.FromField("Description.Collaborators")},

		{
			Name:      "has_organization_projects",
			Type:      proto.ColumnType_BOOL,
			Transform: transform.FromField("Description.HasOrganizationProjects")},

		{
			Name:      "has_repository_projects",
			Type:      proto.ColumnType_BOOL,
			Transform: transform.FromField("Description.HasRepositoryProjects")},

		{
			Name:      "web_commit_signoff_required",
			Type:      proto.ColumnType_BOOL,
			Transform: transform.FromField("Description.WebCommitSignoffRequired")},
	}
}

func sharedOrganizationCountColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "members_with_role_total_count", Type: proto.ColumnType_INT,
			Transform:   transform.FromField("Description.MembersWithRoleTotalCount"),
			Description: "Count of members with a role within the organization."},
		{Name: "packages_total_count", Type: proto.ColumnType_INT,
			Transform:   transform.FromField("Description.PackagesTotalCount"),
			Description: "Count of packages within the organization."},
		{Name: "pinnable_items_total_count", Type: proto.ColumnType_INT,
			Transform:   transform.FromField("Description.PinnableItemsTotalCount"),
			Description: "Count of pinnable items within the organization."},
		{Name: "pinned_items_total_count", Type: proto.ColumnType_INT,
			Transform:   transform.FromField("Description.PinnedItemsTotalCount"),
			Description: "Count of itesm pinned to the organization's profile."},
		{Name: "projects_total_count", Type: proto.ColumnType_INT,
			Transform:   transform.FromField("Description.ProjectsTotalCount"),
			Description: "Count of projects within the organization."},
		{Name: "projects_v2_total_count", Type: proto.ColumnType_INT,
			Transform:   transform.FromField("Description.ProjectsV2TotalCount"),
			Description: "Count of V2 projects within the organization."},
		{Name: "repositories_total_count", Type: proto.ColumnType_INT,
			Transform:   transform.FromField("Description.RepositoriesTotalCount"),
			Description: "Count of all repositories within the organization."},
		{Name: "sponsoring_total_count", Type: proto.ColumnType_INT,
			Transform:   transform.FromField("Description.SponsoringTotalCount"),
			Description: "Count of users the organization is sponsoring."},
		{Name: "sponsors_total_count", Type: proto.ColumnType_INT,
			Transform:   transform.FromField("Description.SponsorsTotalCount"),
			Description: "Count of sponsors the organization has."},
		{Name: "teams_total_count", Type: proto.ColumnType_INT,
			Transform:   transform.FromField("Description.TeamsTotalCount"),
			Description: "Count of teams within the organization."},
		{Name: "repositories_total_disk_usage", Type: proto.ColumnType_INT,
			Transform:   transform.FromField("Description.RepositoriesTotalDiskUsage"),
			Description: "Total disk usage for all repositories within the organization."},
		{Name: "private_repositories_total_count", Type: proto.ColumnType_INT,
			Transform:   transform.FromField("Description.PrivateRepositoriesTotalCount"),
			Description: "Count of private repositories within the organization."},
		{Name: "public_repositories_total_count", Type: proto.ColumnType_INT,
			Transform:   transform.FromField("Description.PublicRepositoriesTotalCount"),
			Description: "Count of public repositories within the organization."},
	}
}

func gitHubOrganizationColumns() []*plugin.Column {
	return append(sharedOrganizationColumns(), sharedOrganizationCountColumns()...)
}

func tableGitHubOrganization() *plugin.Table {
	return &plugin.Table{
		Name:        "github_organization",
		Description: "GitHub Organizations are shared accounts where businesses and open-source projects can collaborate across many projects at once.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListOrganization,
		},
		Columns: commonColumns(gitHubOrganizationColumns()),
	}
}
