package describers

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v55/github"
	"github.com/opengovern/og-describer-template/discovery/pkg/models"
	model "github.com/opengovern/og-describer-template/discovery/provider"
	"github.com/shurcooL/githubv4"
	steampipemodels "github.com/turbot/steampipe-plugin-github/github/models"
)

func GetOrganizationList(ctx context.Context, githubClient model.GitHubClient, organizationName string, stream *models.StreamSender) ([]models.Resource, error) {
	client := githubClient.GraphQLClient
	var query struct {
		RateLimit steampipemodels.RateLimit
		Viewer    struct {
			Organization steampipemodels.OrganizationWithCounts `graphql:"organization(login: $name)"`
		}
	}
	variables := map[string]interface{}{
		"name": githubv4.String(organizationName),
	}
	appendOrganizationColumnIncludes(&variables, organizationCols())
	var values []models.Resource
	err := client.Query(ctx, &query, variables)
	if err != nil {
		return nil, err
	}
	org := query.Viewer.Organization
	hooks, err := GetOrganizationHooks(ctx, githubClient.RestClient, org)
	if err != nil {
		return nil, err
	}
	additionalOrgInfo, err := GetOrganizationAdditionalData(ctx, githubClient.RestClient, org)
	if err != nil {
		return nil, err
	}
	value := models.Resource{
		ID:   strconv.Itoa(org.Id),
		Name: org.Name,
		Description: model.OrganizationDescription{
			Id:                                     org.Id,
			NodeId:                                 org.NodeId,
			Name:                                   org.Name,
			Login:                                  org.Login,
			CreatedAt:                              org.CreatedAt.Format(time.RFC3339),
			UpdatedAt:                              org.UpdatedAt.Format(time.RFC3339),
			Description:                            org.Description,
			Email:                                  org.Email,
			Url:                                    org.Url,
			Announcement:                           org.Announcement,
			AnnouncementExpiresAt:                  org.AnnouncementExpiresAt.Format(time.RFC3339),
			AnnouncementUserDismissible:            org.AnnouncementUserDismissible,
			AnyPinnableItems:                       org.AnyPinnableItems,
			AvatarUrl:                              org.AvatarUrl,
			EstimatedNextSponsorsPayoutInCents:     org.EstimatedNextSponsorsPayoutInCents,
			HasSponsorsListing:                     org.HasSponsorsListing,
			InteractionAbility:                     org.InteractionAbility,
			IsSponsoringYou:                        org.IsSponsoringYou,
			IsVerified:                             org.IsVerified,
			Location:                               org.Location,
			MonthlyEstimatedSponsorsIncomeInCents:  org.MonthlyEstimatedSponsorsIncomeInCents,
			NewTeamUrl:                             org.NewTeamUrl,
			PinnedItemsRemaining:                   org.PinnedItemsRemaining,
			ProjectsUrl:                            org.ProjectsUrl,
			SamlIdentityProvider:                   org.SamlIdentityProvider,
			SponsorsListing:                        org.SponsorsListing,
			TeamsUrl:                               org.TeamsUrl,
			TotalSponsorshipAmountAsSponsorInCents: org.TotalSponsorshipAmountAsSponsorInCents,
			TwitterUsername:                        org.TwitterUsername,
			CanAdminister:                          org.CanAdminister,
			CanChangedPinnedItems:                  org.CanChangedPinnedItems,
			CanCreateProjects:                      org.CanCreateProjects,
			CanCreateRepositories:                  org.CanCreateRepositories,
			CanCreateTeams:                         org.CanCreateTeams,
			CanSponsor:                             org.CanSponsor,
			IsAMember:                              org.IsAMember,
			IsFollowing:                            org.IsFollowing,
			IsSponsoring:                           org.IsSponsoring,
			WebsiteUrl:                             org.WebsiteUrl,
			Hooks:                                  hooks,
			BillingEmail:                           additionalOrgInfo.GetBillingEmail(),
			TwoFactorRequirementEnabled:            additionalOrgInfo.GetTwoFactorRequirementEnabled(),
			DefaultRepoPermission:                  additionalOrgInfo.GetDefaultRepoPermission(),
			MembersAllowedRepositoryCreationType:   additionalOrgInfo.GetMembersAllowedRepositoryCreationType(),
			MembersCanCreateInternalRepos:          additionalOrgInfo.GetMembersCanCreateInternalRepos(),
			MembersCanCreatePages:                  additionalOrgInfo.GetMembersCanCreatePages(),
			MembersCanCreatePrivateRepos:           additionalOrgInfo.GetMembersCanCreatePrivateRepos(),
			MembersCanCreatePublicRepos:            additionalOrgInfo.GetMembersCanCreatePublicRepos(),
			MembersCanCreateRepos:                  additionalOrgInfo.GetMembersCanCreateRepos(),
			MembersCanForkPrivateRepos:             additionalOrgInfo.GetMembersCanForkPrivateRepos(),
			PlanFilledSeats:                        additionalOrgInfo.GetPlan().GetFilledSeats(),
			PlanName:                               additionalOrgInfo.GetPlan().GetName(),
			PlanPrivateRepos:                       int(additionalOrgInfo.GetPlan().GetPrivateRepos()),
			PlanSeats:                              additionalOrgInfo.GetPlan().GetSeats(),
			PlanSpace:                              additionalOrgInfo.GetPlan().GetSpace(),
			Followers:                              additionalOrgInfo.GetFollowers(),
			Following:                              additionalOrgInfo.GetFollowing(),
			Collaborators:                          additionalOrgInfo.GetCollaborators(),
			HasOrganizationProjects:                additionalOrgInfo.GetHasOrganizationProjects(),
			HasRepositoryProjects:                  additionalOrgInfo.GetHasRepositoryProjects(),
			WebCommitSignoffRequired:               additionalOrgInfo.GetWebCommitSignoffRequired(),
			MembersWithRoleTotalCount:              org.MembersWithRole.TotalCount,
			PackagesTotalCount:                     org.Packages.TotalCount,
			PinnableItemsTotalCount:                org.PinnableItems.TotalCount,
			PinnedItemsTotalCount:                  org.PinnedItems.TotalCount,
			ProjectsTotalCount:                     org.Projects.TotalCount,
			ProjectsV2TotalCount:                   org.ProjectsV2.TotalCount,
			SponsoringTotalCount:                   org.Sponsoring.TotalCount,
			SponsorsTotalCount:                     org.Sponsors.TotalCount,
			TeamsTotalCount:                        org.Teams.TotalCount,
			PrivateRepositoriesTotalCount:          org.PrivateRepositories.TotalCount,
			PublicRepositoriesTotalCount:           org.PublicRepositories.TotalCount,
			RepositoriesTotalCount:                 org.Repositories.TotalCount,
			RepositoriesTotalDiskUsage:             org.Repositories.TotalDiskUsage,
		},
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

func GetOrganizationHooks(ctx context.Context, client *github.Client, org steampipemodels.OrganizationWithCounts) ([]*github.Hook, error) {
	login := org.Login
	var orgHooks []*github.Hook
	opt := &github.ListOptions{PerPage: pageSize}
	for {
		hooks, resp, err := client.Organizations.ListHooks(ctx, login, opt)
		if err != nil && strings.Contains(err.Error(), "Not Found") {
			return nil, nil
		} else if err != nil {
			return nil, err
		}
		orgHooks = append(orgHooks, hooks...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return orgHooks, nil
}

func GetOrganizationAdditionalData(ctx context.Context, client *github.Client, org steampipemodels.OrganizationWithCounts) (*github.Organization, error) {
	login := org.Login
	organization, _, err := client.Organizations.Get(ctx, login)
	if err != nil {
		return nil, err
	}
	return organization, nil
}
