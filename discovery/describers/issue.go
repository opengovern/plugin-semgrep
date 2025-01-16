package describers

import (
	"context"
	"math"
	"strconv"
	"time"

	"github.com/opengovern/og-describer-template/discovery/pkg/models"
	model "github.com/opengovern/og-describer-template/discovery/provider"
	"github.com/shurcooL/githubv4"
	steampipemodels "github.com/turbot/steampipe-plugin-github/github/models"
)

func GetIssueList(ctx context.Context, githubClient model.GitHubClient, organizationName string, stream *models.StreamSender) ([]models.Resource, error) {
	client := githubClient.GraphQLClient
	var filters githubv4.IssueFilters
	filters.States = &[]githubv4.IssueState{githubv4.IssueStateOpen, githubv4.IssueStateClosed}
	repositories, err := getRepositories(ctx, githubClient.RestClient, organizationName)
	if err != nil {
		return nil, nil
	}
	var query struct {
		RateLimit  steampipemodels.RateLimit
		Repository struct {
			Issues struct {
				TotalCount int
				PageInfo   steampipemodels.PageInfo
				Nodes      []steampipemodels.Issue
			} `graphql:"issues(first: $pageSize, after: $cursor, filterBy: $filters)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}
	var values []models.Resource

	for _, r := range repositories {
		variables := map[string]interface{}{
			"owner":    githubv4.String(organizationName),
			"name":     githubv4.String(r.GetName()),
			"pageSize": githubv4.Int(pageSize),
			"cursor":   (*githubv4.String)(nil),
			"filters":  filters,
		}
		appendIssueColumnIncludes(&variables, issueCols())
		for {
			err := client.Query(ctx, &query, variables)
			if err != nil {
				return nil, err
			}
			for _, issue := range query.Repository.Issues.Nodes {
				labelsSrcLength := int(math.Min(float64(len(issue.Labels.Nodes)), 100.0))
				labelsSrc := issue.Labels.Nodes[:labelsSrcLength]
				var finalLabelsSrc []model.Label
				for _, labelSrc := range labelsSrc {
					finalLabelsSrc = append(finalLabelsSrc, model.Label{
						NodeId:      &labelSrc.NodeId,
						Name:        &labelSrc.Name,
						Description: &labelSrc.Description,
						IsDefault:   labelSrc.IsDefault,
						Color:       &labelSrc.Color,
					})
				}
				labels := make(map[string]model.Label)
				for _, label := range issue.Labels.Nodes {
					labels[label.Name] = model.Label{
						NodeId:      &label.NodeId,
						Name:        &label.Name,
						Description: &label.Description,
						IsDefault:   label.IsDefault,
						Color:       &label.Color,
					}
				}

				author := model.Actor{
					AvatarUrl: &issue.Author.AvatarUrl,
					Login:     &issue.Author.Login,
					Url:       &issue.Author.Url,
				}

				editor := model.Actor{
					AvatarUrl: &issue.Editor.AvatarUrl,
					Login:     &issue.Editor.Login,
					Url:       &issue.Editor.Url,
				}

				milestoneClosedAt := issue.Milestone.ClosedAt.Format(time.RFC3339)
				milestoneCreatedAt := issue.Milestone.CreatedAt.Format(time.RFC3339)
				milestoneDueOn := issue.Milestone.DueOn.Format(time.RFC3339)
				milestoneUpdatedAt := issue.Milestone.UpdatedAt.Format(time.RFC3339)

				milestoneCreator := model.Actor{
					AvatarUrl: &issue.Milestone.Creator.AvatarUrl,
					Login:     &issue.Milestone.Creator.Login,
					Url:       &issue.Milestone.Creator.Url,
				}

				milestone := model.Milestone{
					Closed:             issue.Milestone.Closed,
					ClosedAt:           &milestoneClosedAt,
					CreatedAt:          &milestoneCreatedAt,
					Creator:            milestoneCreator,
					Description:        &issue.Milestone.Description,
					DueOn:              &milestoneDueOn,
					Number:             issue.Milestone.Number,
					ProgressPercentage: issue.Milestone.ProgressPercentage,
					State:              &issue.Milestone.State,
					Title:              &issue.Milestone.Title,
					UpdatedAt:          &milestoneUpdatedAt,
					UserCanClose:       issue.Milestone.UserCanClose,
					UserCanReopen:      issue.Milestone.UserCanReopen,
				}

				var assignees []model.BaseUser
				for _, assignee := range issue.Assignees.Nodes {
					assigneeCreatedAt := assignee.CreatedAt.Format(time.RFC3339)
					assigneeUpdatedAt := assignee.UpdatedAt.Format(time.RFC3339)
					interactionAbilityExpiresAt := assignee.InteractionAbility.ExpiresAt.Format(time.RFC3339)
					sponsorListingCreatedAt := assignee.SponsorsListing.CreatedAt.Format(time.RFC3339)
					nextPayoutDate := assignee.SponsorsListing.NextPayoutDate.Format(time.RFC3339)
					statusCreatedAt := assignee.Status.CreatedAt.Format(time.RFC3339)
					statusUpdatedAt := assignee.Status.UpdatedAt.Format(time.RFC3339)
					statusExpiresAt := assignee.Status.ExpiresAt.Format(time.RFC3339)

					interactionAbility := model.RepositoryInteractionAbility{
						ExpiresAt: &interactionAbilityExpiresAt,
						Limit:     &assignee.InteractionAbility.Limit,
						Origin:    &assignee.InteractionAbility.Origin,
					}

					activeGoal := model.SponsorsGoal{
						Description:     &assignee.SponsorsListing.ActiveGoal.Description,
						PercentComplete: assignee.SponsorsListing.ActiveGoal.PercentComplete,
						TargetValue:     assignee.SponsorsListing.ActiveGoal.TargetValue,
						Title:           &assignee.SponsorsListing.ActiveGoal.Title,
						Kind:            &assignee.SponsorsListing.ActiveGoal.Kind,
					}

					activeStripeConnectAccount := model.StripeConnectAccount{
						AccountId:              &assignee.SponsorsListing.ActiveStripeConnectAccount.AccountId,
						BillingCountryOrRegion: &assignee.SponsorsListing.ActiveStripeConnectAccount.BillingCountryOrRegion,
						CountryOrRegion:        &assignee.SponsorsListing.ActiveStripeConnectAccount.CountryOrRegion,
						IsActive:               assignee.SponsorsListing.ActiveStripeConnectAccount.IsActive,
						StripeDashboardUrl:     &assignee.SponsorsListing.ActiveStripeConnectAccount.StripeDashboardUrl,
					}

					sponsorListing := model.SponsorsListing{
						Id:                         &assignee.SponsorsListing.Id,
						ActiveGoal:                 activeGoal,
						ActiveStripeConnectAccount: activeStripeConnectAccount,
						BillingCountryOrRegion:     &assignee.SponsorsListing.BillingCountryOrRegion,
						ContactEmailAddress:        &assignee.SponsorsListing.ContactEmailAddress,
						CreatedAt:                  &sponsorListingCreatedAt,
						DashboardUrl:               &assignee.SponsorsListing.DashboardUrl,
						FullDescription:            &assignee.SponsorsListing.FullDescription,
						IsPublic:                   assignee.SponsorsListing.IsPublic,
						Name:                       &assignee.SponsorsListing.Name,
						NextPayoutDate:             &nextPayoutDate,
						ResidenceCountryOrRegion:   &assignee.SponsorsListing.ResidenceCountryOrRegion,
						ShortDescription:           &assignee.SponsorsListing.ShortDescription,
						Slug:                       &assignee.SponsorsListing.Slug,
						Url:                        &assignee.SponsorsListing.Url,
					}

					status := model.UserStatus{
						CreatedAt:                    &statusCreatedAt,
						UpdatedAt:                    &statusUpdatedAt,
						ExpiresAt:                    &statusExpiresAt,
						Emoji:                        &assignee.Status.Emoji,
						Message:                      &assignee.Status.Message,
						IndicatesLimitedAvailability: assignee.Status.IndicatesLimitedAvailability,
					}

					assignees = append(assignees, model.BaseUser{
						BasicUser: model.BasicUser{
							Id:        assignee.Id,
							NodeId:    &assignee.NodeId,
							Name:      &assignee.Name,
							Login:     &assignee.Login,
							Email:     &assignee.Email,
							CreatedAt: &assigneeCreatedAt,
							UpdatedAt: &assigneeUpdatedAt,
							Url:       &assignee.Url,
						},
						AnyPinnableItems:                      assignee.AnyPinnableItems,
						AvatarUrl:                             &assignee.AvatarUrl,
						Bio:                                   &assignee.Bio,
						Company:                               &assignee.Company,
						EstimatedNextSponsorsPayoutInCents:    assignee.EstimatedNextSponsorsPayoutInCents,
						HasSponsorsListing:                    assignee.HasSponsorsListing,
						InteractionAbility:                    interactionAbility,
						IsBountyHunter:                        assignee.IsBountyHunter,
						IsCampusExpert:                        assignee.IsCampusExpert,
						IsDeveloperProgramMember:              assignee.IsDeveloperProgramMember,
						IsEmployee:                            assignee.IsEmployee,
						IsFollowingYou:                        assignee.IsFollowingYou,
						IsGitHubStar:                          assignee.IsGitHubStar,
						IsHireable:                            assignee.IsHireable,
						IsSiteAdmin:                           assignee.IsSiteAdmin,
						IsSponsoringYou:                       assignee.IsSponsoringYou,
						IsYou:                                 assignee.IsYou,
						Location:                              &assignee.Location,
						MonthlyEstimatedSponsorsIncomeInCents: assignee.MonthlyEstimatedSponsorsIncomeInCents,
						PinnedItemsRemaining:                  assignee.PinnedItemsRemaining,
						ProjectsUrl:                           &assignee.ProjectsUrl,
						Pronouns:                              &assignee.Pronouns,
						SponsorsListing:                       sponsorListing,
						Status:                                status,
						TwitterUsername:                       &assignee.TwitterUsername,
						CanChangedPinnedItems:                 assignee.CanChangedPinnedItems,
						CanCreateProjects:                     assignee.CanCreateProjects,
						CanFollow:                             assignee.CanFollow,
						CanSponsor:                            assignee.CanSponsor,
						IsFollowing:                           assignee.IsFollowing,
						IsSponsoring:                          assignee.IsSponsoring,
						WebsiteUrl:                            &assignee.WebsiteUrl,
					})
				}

				closedAt := issue.ClosedAt.Format(time.RFC3339)
				createdAt := issue.CreatedAt.Format(time.RFC3339)
				lastEditedAt := issue.LastEditedAt.Format(time.RFC3339)
				publishedAt := issue.PublishedAt.Format(time.RFC3339)
				updatedAt := issue.UpdatedAt.Format(time.RFC3339)

				value := models.Resource{
					ID:   strconv.Itoa(issue.Id),
					Name: issue.Title,
					Description: model.IssueDescription{
						RepositoryFullName:      r.FullName,
						Id:                      issue.Id,
						NodeId:                  &issue.NodeId,
						Number:                  issue.Number,
						ActiveLockReason:        &issue.ActiveLockReason,
						Author:                  author,
						AuthorLogin:             &issue.Author.Login,
						AuthorAssociation:       &issue.AuthorAssociation,
						Body:                    &issue.Body,
						BodyUrl:                 &issue.BodyUrl,
						Closed:                  issue.Closed,
						ClosedAt:                &closedAt,
						CreatedAt:               &createdAt,
						CreatedViaEmail:         issue.CreatedViaEmail,
						Editor:                  editor,
						FullDatabaseId:          &issue.FullDatabaseId,
						IncludesCreatedEdit:     issue.IncludesCreatedEdit,
						IsPinned:                issue.IsPinned,
						IsReadByUser:            issue.IsReadByUser,
						LastEditedAt:            &lastEditedAt,
						Locked:                  issue.Locked,
						Milestone:               milestone,
						PublishedAt:             &publishedAt,
						State:                   &issue.State,
						StateReason:             &issue.StateReason,
						Title:                   &issue.Title,
						UpdatedAt:               &updatedAt,
						Url:                     &issue.Url,
						UserCanClose:            issue.UserCanClose,
						UserCanReact:            issue.UserCanReact,
						UserCanReopen:           issue.UserCanReopen,
						UserCanSubscribe:        issue.UserCanSubscribe,
						UserCanUpdate:           issue.UserCanUpdate,
						UserCannotUpdateReasons: issue.UserCannotUpdateReasons,
						UserDidAuthor:           issue.UserDidAuthor,
						UserSubscription:        &issue.UserSubscription,
						CommentsTotalCount:      issue.Comments.TotalCount,
						LabelsTotalCount:        issue.Labels.TotalCount,
						LabelsSrc:               finalLabelsSrc,
						Labels:                  labels,
						AssigneesTotalCount:     issue.Assignees.TotalCount,
						Assignees:               assignees,
					},
				}
				if stream != nil {
					if err := (*stream)(value); err != nil {
						return nil, err
					}
				} else {
					values = append(values, value)
				}
			}
			if !query.Repository.Issues.PageInfo.HasNextPage {
				break
			}
			variables["cursor"] = githubv4.NewString(query.Repository.Issues.PageInfo.EndCursor)
		}
	}
	return values, nil
}

func GetIssue(ctx context.Context, githubClient model.GitHubClient, organizationName string, repositoryName string, resourceID string, stream *models.StreamSender) (*models.Resource, error) {
	repoFullName := formRepositoryFullName(organizationName, repositoryName)
	issueID, err := strconv.ParseInt(resourceID, 10, 64)
	if err != nil {
		return nil, err
	}
	client := githubClient.GraphQLClient

	var query struct {
		RateLimit  steampipemodels.RateLimit
		Repository struct {
			Issue steampipemodels.Issue `graphql:"issue(number: $issueNumber)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}

	variables := map[string]interface{}{
		"owner":       githubv4.String(organizationName),
		"repo":        githubv4.String(repositoryName),
		"issueNumber": githubv4.Int(issueID),
	}
	appendIssueColumnIncludes(&variables, issueCols())

	err = client.Query(ctx, &query, variables)
	if err != nil {
		return nil, err
	}
	labelsSrcLength := int(math.Min(float64(len(query.Repository.Issue.Labels.Nodes)), 100.0))
	labelsSrc := query.Repository.Issue.Labels.Nodes[:labelsSrcLength]
	var finalLabelsSrc []model.Label
	for _, labelSrc := range labelsSrc {
		finalLabelsSrc = append(finalLabelsSrc, model.Label{
			NodeId:      &labelSrc.NodeId,
			Name:        &labelSrc.Name,
			Description: &labelSrc.Description,
			IsDefault:   labelSrc.IsDefault,
			Color:       &labelSrc.Color,
		})
	}
	labels := make(map[string]model.Label)
	for _, label := range query.Repository.Issue.Labels.Nodes {
		labels[label.Name] = model.Label{
			NodeId:      &label.NodeId,
			Name:        &label.Name,
			Description: &label.Description,
			IsDefault:   label.IsDefault,
			Color:       &label.Color,
		}
	}

	author := model.Actor{
		AvatarUrl: &query.Repository.Issue.Author.AvatarUrl,
		Login:     &query.Repository.Issue.Author.Login,
		Url:       &query.Repository.Issue.Author.Url,
	}

	editor := model.Actor{
		AvatarUrl: &query.Repository.Issue.Editor.AvatarUrl,
		Login:     &query.Repository.Issue.Editor.Login,
		Url:       &query.Repository.Issue.Editor.Url,
	}

	milestoneClosedAt := query.Repository.Issue.Milestone.ClosedAt.Format(time.RFC3339)
	milestoneCreatedAt := query.Repository.Issue.Milestone.CreatedAt.Format(time.RFC3339)
	milestoneDueOn := query.Repository.Issue.Milestone.DueOn.Format(time.RFC3339)
	milestoneUpdatedAt := query.Repository.Issue.Milestone.UpdatedAt.Format(time.RFC3339)

	milestoneCreator := model.Actor{
		AvatarUrl: &query.Repository.Issue.Milestone.Creator.AvatarUrl,
		Login:     &query.Repository.Issue.Milestone.Creator.Login,
		Url:       &query.Repository.Issue.Milestone.Creator.Url,
	}

	milestone := model.Milestone{
		Closed:             query.Repository.Issue.Milestone.Closed,
		ClosedAt:           &milestoneClosedAt,
		CreatedAt:          &milestoneCreatedAt,
		Creator:            milestoneCreator,
		Description:        &query.Repository.Issue.Milestone.Description,
		DueOn:              &milestoneDueOn,
		Number:             query.Repository.Issue.Milestone.Number,
		ProgressPercentage: query.Repository.Issue.Milestone.ProgressPercentage,
		State:              &query.Repository.Issue.Milestone.State,
		Title:              &query.Repository.Issue.Milestone.Title,
		UpdatedAt:          &milestoneUpdatedAt,
		UserCanClose:       query.Repository.Issue.Milestone.UserCanClose,
		UserCanReopen:      query.Repository.Issue.Milestone.UserCanReopen,
	}

	var assignees []model.BaseUser
	for _, assignee := range query.Repository.Issue.Assignees.Nodes {
		assigneeCreatedAt := assignee.CreatedAt.Format(time.RFC3339)
		assigneeUpdatedAt := assignee.UpdatedAt.Format(time.RFC3339)
		interactionAbilityExpiresAt := assignee.InteractionAbility.ExpiresAt.Format(time.RFC3339)
		sponsorListingCreatedAt := assignee.SponsorsListing.CreatedAt.Format(time.RFC3339)
		nextPayoutDate := assignee.SponsorsListing.NextPayoutDate.Format(time.RFC3339)
		statusCreatedAt := assignee.Status.CreatedAt.Format(time.RFC3339)
		statusUpdatedAt := assignee.Status.UpdatedAt.Format(time.RFC3339)
		statusExpiresAt := assignee.Status.ExpiresAt.Format(time.RFC3339)

		interactionAbility := model.RepositoryInteractionAbility{
			ExpiresAt: &interactionAbilityExpiresAt,
			Limit:     &assignee.InteractionAbility.Limit,
			Origin:    &assignee.InteractionAbility.Origin,
		}

		activeGoal := model.SponsorsGoal{
			Description:     &assignee.SponsorsListing.ActiveGoal.Description,
			PercentComplete: assignee.SponsorsListing.ActiveGoal.PercentComplete,
			TargetValue:     assignee.SponsorsListing.ActiveGoal.TargetValue,
			Title:           &assignee.SponsorsListing.ActiveGoal.Title,
			Kind:            &assignee.SponsorsListing.ActiveGoal.Kind,
		}

		activeStripeConnectAccount := model.StripeConnectAccount{
			AccountId:              &assignee.SponsorsListing.ActiveStripeConnectAccount.AccountId,
			BillingCountryOrRegion: &assignee.SponsorsListing.ActiveStripeConnectAccount.BillingCountryOrRegion,
			CountryOrRegion:        &assignee.SponsorsListing.ActiveStripeConnectAccount.CountryOrRegion,
			IsActive:               assignee.SponsorsListing.ActiveStripeConnectAccount.IsActive,
			StripeDashboardUrl:     &assignee.SponsorsListing.ActiveStripeConnectAccount.StripeDashboardUrl,
		}

		sponsorListing := model.SponsorsListing{
			Id:                         &assignee.SponsorsListing.Id,
			ActiveGoal:                 activeGoal,
			ActiveStripeConnectAccount: activeStripeConnectAccount,
			BillingCountryOrRegion:     &assignee.SponsorsListing.BillingCountryOrRegion,
			ContactEmailAddress:        &assignee.SponsorsListing.ContactEmailAddress,
			CreatedAt:                  &sponsorListingCreatedAt,
			DashboardUrl:               &assignee.SponsorsListing.DashboardUrl,
			FullDescription:            &assignee.SponsorsListing.FullDescription,
			IsPublic:                   assignee.SponsorsListing.IsPublic,
			Name:                       &assignee.SponsorsListing.Name,
			NextPayoutDate:             &nextPayoutDate,
			ResidenceCountryOrRegion:   &assignee.SponsorsListing.ResidenceCountryOrRegion,
			ShortDescription:           &assignee.SponsorsListing.ShortDescription,
			Slug:                       &assignee.SponsorsListing.Slug,
			Url:                        &assignee.SponsorsListing.Url,
		}

		status := model.UserStatus{
			CreatedAt:                    &statusCreatedAt,
			UpdatedAt:                    &statusUpdatedAt,
			ExpiresAt:                    &statusExpiresAt,
			Emoji:                        &assignee.Status.Emoji,
			Message:                      &assignee.Status.Message,
			IndicatesLimitedAvailability: assignee.Status.IndicatesLimitedAvailability,
		}

		assignees = append(assignees, model.BaseUser{
			BasicUser: model.BasicUser{
				Id:        assignee.Id,
				NodeId:    &assignee.NodeId,
				Name:      &assignee.Name,
				Login:     &assignee.Login,
				Email:     &assignee.Email,
				CreatedAt: &assigneeCreatedAt,
				UpdatedAt: &assigneeUpdatedAt,
				Url:       &assignee.Url,
			},
			AnyPinnableItems:                      assignee.AnyPinnableItems,
			AvatarUrl:                             &assignee.AvatarUrl,
			Bio:                                   &assignee.Bio,
			Company:                               &assignee.Company,
			EstimatedNextSponsorsPayoutInCents:    assignee.EstimatedNextSponsorsPayoutInCents,
			HasSponsorsListing:                    assignee.HasSponsorsListing,
			InteractionAbility:                    interactionAbility,
			IsBountyHunter:                        assignee.IsBountyHunter,
			IsCampusExpert:                        assignee.IsCampusExpert,
			IsDeveloperProgramMember:              assignee.IsDeveloperProgramMember,
			IsEmployee:                            assignee.IsEmployee,
			IsFollowingYou:                        assignee.IsFollowingYou,
			IsGitHubStar:                          assignee.IsGitHubStar,
			IsHireable:                            assignee.IsHireable,
			IsSiteAdmin:                           assignee.IsSiteAdmin,
			IsSponsoringYou:                       assignee.IsSponsoringYou,
			IsYou:                                 assignee.IsYou,
			Location:                              &assignee.Location,
			MonthlyEstimatedSponsorsIncomeInCents: assignee.MonthlyEstimatedSponsorsIncomeInCents,
			PinnedItemsRemaining:                  assignee.PinnedItemsRemaining,
			ProjectsUrl:                           &assignee.ProjectsUrl,
			Pronouns:                              &assignee.Pronouns,
			SponsorsListing:                       sponsorListing,
			Status:                                status,
			TwitterUsername:                       &assignee.TwitterUsername,
			CanChangedPinnedItems:                 assignee.CanChangedPinnedItems,
			CanCreateProjects:                     assignee.CanCreateProjects,
			CanFollow:                             assignee.CanFollow,
			CanSponsor:                            assignee.CanSponsor,
			IsFollowing:                           assignee.IsFollowing,
			IsSponsoring:                          assignee.IsSponsoring,
			WebsiteUrl:                            &assignee.WebsiteUrl,
		})
	}

	closedAt := query.Repository.Issue.ClosedAt.Format(time.RFC3339)
	createdAt := query.Repository.Issue.CreatedAt.Format(time.RFC3339)
	lastEditedAt := query.Repository.Issue.LastEditedAt.Format(time.RFC3339)
	publishedAt := query.Repository.Issue.PublishedAt.Format(time.RFC3339)
	updatedAt := query.Repository.Issue.UpdatedAt.Format(time.RFC3339)

	value := models.Resource{
		ID:   strconv.Itoa(query.Repository.Issue.Id),
		Name: query.Repository.Issue.Title,
		Description: model.IssueDescription{
			RepositoryFullName:      &repoFullName,
			Id:                      query.Repository.Issue.Id,
			NodeId:                  &query.Repository.Issue.NodeId,
			Number:                  query.Repository.Issue.Number,
			ActiveLockReason:        &query.Repository.Issue.ActiveLockReason,
			Author:                  author,
			AuthorLogin:             &query.Repository.Issue.Author.Login,
			AuthorAssociation:       &query.Repository.Issue.AuthorAssociation,
			Body:                    &query.Repository.Issue.Body,
			BodyUrl:                 &query.Repository.Issue.BodyUrl,
			Closed:                  query.Repository.Issue.Closed,
			ClosedAt:                &closedAt,
			CreatedAt:               &createdAt,
			CreatedViaEmail:         query.Repository.Issue.CreatedViaEmail,
			Editor:                  editor,
			FullDatabaseId:          &query.Repository.Issue.FullDatabaseId,
			IncludesCreatedEdit:     query.Repository.Issue.IncludesCreatedEdit,
			IsPinned:                query.Repository.Issue.IsPinned,
			IsReadByUser:            query.Repository.Issue.IsReadByUser,
			LastEditedAt:            &lastEditedAt,
			Locked:                  query.Repository.Issue.Locked,
			Milestone:               milestone,
			PublishedAt:             &publishedAt,
			State:                   &query.Repository.Issue.State,
			StateReason:             &query.Repository.Issue.StateReason,
			Title:                   &query.Repository.Issue.Title,
			UpdatedAt:               &updatedAt,
			Url:                     &query.Repository.Issue.Url,
			UserCanClose:            query.Repository.Issue.UserCanClose,
			UserCanReact:            query.Repository.Issue.UserCanReact,
			UserCanReopen:           query.Repository.Issue.UserCanReopen,
			UserCanSubscribe:        query.Repository.Issue.UserCanSubscribe,
			UserCanUpdate:           query.Repository.Issue.UserCanUpdate,
			UserCannotUpdateReasons: query.Repository.Issue.UserCannotUpdateReasons,
			UserDidAuthor:           query.Repository.Issue.UserDidAuthor,
			UserSubscription:        &query.Repository.Issue.UserSubscription,
			CommentsTotalCount:      query.Repository.Issue.Comments.TotalCount,
			LabelsTotalCount:        query.Repository.Issue.Labels.TotalCount,
			LabelsSrc:               finalLabelsSrc,
			Labels:                  labels,
			AssigneesTotalCount:     query.Repository.Issue.Assignees.TotalCount,
			Assignees:               assignees,
		},
	}
	if stream != nil {
		if err := (*stream)(value); err != nil {
			return nil, err
		}
	}
	return &value, nil
}
