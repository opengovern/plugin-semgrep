package describers

import (
	"context"
	"fmt"
	"time"

	"github.com/opengovern/og-describer-template/discovery/pkg/models"
	model "github.com/opengovern/og-describer-template/discovery/provider"
	"github.com/shurcooL/githubv4"
	steampipemodels "github.com/turbot/steampipe-plugin-github/github/models"
)

func GetAllBranches(ctx context.Context, githubClient model.GitHubClient, organizationName string, stream *models.StreamSender) ([]models.Resource, error) {
	client := githubClient.RestClient
	owner := organizationName
	repositories, err := getRepositories(ctx, client, owner)
	if err != nil {
		return nil, nil
	}
	var values []models.Resource
	for _, repo := range repositories {
		repoValues, err := GetRepositoryBranches(ctx, githubClient, stream, owner, repo.GetName())
		if err != nil {
			return nil, err
		}
		values = append(values, repoValues...)
	}
	return values, nil
}

func GetRepositoryBranches(ctx context.Context, githubClient model.GitHubClient, stream *models.StreamSender, owner, repo string) ([]models.Resource, error) {
	graphQLClient := githubClient.GraphQLClient
	restClient := githubClient.RestClient
	var query struct {
		RateLimit  steampipemodels.RateLimit
		Repository struct {
			Refs struct {
				TotalCount int
				PageInfo   steampipemodels.PageInfo
				Edges      []struct {
					Node steampipemodels.Branch
				}
			} `graphql:"refs(refPrefix: \"refs/heads/\", first: $pageSize, after: $cursor)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}
	variables := map[string]interface{}{
		"owner":    githubv4.String(owner),
		"repo":     githubv4.String(repo),
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
	}
	appendBranchColumnIncludes(&variables, branchCols())
	repoFullName := formRepositoryFullName(owner, repo)
	var values []models.Resource
	for {
		err := graphQLClient.Query(ctx, &query, variables)
		if err != nil {
			return nil, err
		}
		for _, branch := range query.Repository.Refs.Edges {
			branchInfo, _, err := restClient.Repositories.GetBranch(ctx, owner, repo, branch.Node.Name, true)
			if err != nil {
				return nil, err
			}
			protected := branchInfo.GetProtected()
			id := fmt.Sprintf("%s/%s/%s", owner, repo, branch.Node.Name)

			authoredDate := branch.Node.Target.Commit.AuthoredDate.Format(time.RFC3339)
			authorDate := branch.Node.Target.Commit.Author.Date.Format(time.RFC3339)
			authorUserCreatedAt := branch.Node.Target.Commit.Author.User.CreatedAt.Format(time.RFC3339)
			authorUserUpdatedAt := branch.Node.Target.Commit.Author.User.UpdatedAt.Format(time.RFC3339)
			commiterUserCreatedAt := branch.Node.Target.Commit.Committer.User.CreatedAt.Format(time.RFC3339)
			commiterUserUpdatedAt := branch.Node.Target.Commit.Committer.User.UpdatedAt.Format(time.RFC3339)
			committedDate := branch.Node.Target.Commit.CommittedDate.Format(time.RFC3339)
			commiterDate := branch.Node.Target.Commit.Committer.Date.Format(time.RFC3339)

			authorUser := model.BasicUser{
				Id:        branch.Node.Target.Commit.Author.User.Id,
				NodeId:    &branch.Node.Target.Commit.Author.User.NodeId,
				Name:      &branch.Node.Target.Commit.Author.User.Name,
				Login:     &branch.Node.Target.Commit.Author.User.Login,
				Email:     &branch.Node.Target.Commit.Author.User.Email,
				CreatedAt: &authorUserCreatedAt,
				UpdatedAt: &authorUserUpdatedAt,
				Url:       &branch.Node.Target.Commit.Author.User.Url,
			}

			commiterUser := model.BasicUser{
				Id:        branch.Node.Target.Commit.Committer.User.Id,
				NodeId:    &branch.Node.Target.Commit.Committer.User.NodeId,
				Name:      &branch.Node.Target.Commit.Committer.User.Name,
				Login:     &branch.Node.Target.Commit.Committer.User.Login,
				Email:     &branch.Node.Target.Commit.Committer.User.Email,
				CreatedAt: &commiterUserCreatedAt,
				UpdatedAt: &commiterUserUpdatedAt,
				Url:       &branch.Node.Target.Commit.Committer.User.Url,
			}

			author := model.GitActor{
				AvatarUrl: &branch.Node.Target.Commit.Author.AvatarUrl,
				Date:      &authorDate,
				Email:     &branch.Node.Target.Commit.Author.Email,
				Name:      &branch.Node.Target.Commit.Author.Name,
				User:      authorUser,
			}

			commiter := model.GitActor{
				AvatarUrl: &branch.Node.Target.Commit.Committer.AvatarUrl,
				Date:      &commiterDate,
				Email:     &branch.Node.Target.Commit.Committer.Email,
				Name:      &branch.Node.Target.Commit.Committer.Name,
				User:      commiterUser,
			}

			signature := model.Signature{
				Email:             &branch.Node.Target.Commit.Signature.Email,
				IsValid:           branch.Node.Target.Commit.Signature.IsValid,
				State:             &branch.Node.Target.Commit.Signature.State,
				WasSignedByGitHub: branch.Node.Target.Commit.Signature.WasSignedByGitHub,
				Signer: struct {
					Email *string
					Login *string
				}{Email: &branch.Node.Target.Commit.Signature.Signer.Email, Login: &branch.Node.Target.Commit.Signature.Signer.Login},
			}

			status := model.CommitStatus{
				State: &branch.Node.Target.Commit.Status.State,
			}

			commit := model.BaseCommit{
				Sha:                 &branch.Node.Target.Commit.Sha,
				ShortSha:            &branch.Node.Target.Commit.ShortSha,
				AuthoredDate:        &authoredDate,
				Author:              author,
				CommittedDate:       &committedDate,
				Committer:           commiter,
				Message:             &branch.Node.Target.Commit.Message,
				Url:                 &branch.Node.Target.Commit.Url,
				Additions:           branch.Node.Target.Commit.Additions,
				AuthoredByCommitter: branch.Node.Target.Commit.AuthoredByCommitter,
				ChangedFiles:        branch.Node.Target.Commit.ChangedFiles,
				CommittedViaWeb:     branch.Node.Target.Commit.CommittedViaWeb,
				CommitUrl:           &branch.Node.Target.Commit.CommitUrl,
				Deletions:           branch.Node.Target.Commit.Deletions,
				Signature:           signature,
				TarballUrl:          &branch.Node.Target.Commit.TarballUrl,
				TreeUrl:             &branch.Node.Target.Commit.TreeUrl,
				CanSubscribe:        branch.Node.Target.Commit.CanSubscribe,
				Subscription:        &branch.Node.Target.Commit.Subscription,
				ZipballUrl:          &branch.Node.Target.Commit.ZipballUrl,
				MessageHeadline:     &branch.Node.Target.Commit.MessageHeadline,
				Status:              status,
				NodeId:              &branch.Node.Target.Commit.NodeId,
			}

			branchProtctionRule := model.BranchProtectionRule{
				AllowsDeletions:                branch.Node.BranchProtectionRule.AllowsDeletions,
				AllowsForcePushes:              branch.Node.BranchProtectionRule.AllowsForcePushes,
				BlocksCreations:                branch.Node.BranchProtectionRule.BlocksCreations,
				CreatorLogin:                   &branch.Node.BranchProtectionRule.Creator.Login,
				Id:                             branch.Node.BranchProtectionRule.Id,
				NodeId:                         &branch.Node.BranchProtectionRule.NodeId,
				DismissesStaleReviews:          branch.Node.BranchProtectionRule.DismissesStaleReviews,
				IsAdminEnforced:                branch.Node.BranchProtectionRule.IsAdminEnforced,
				LockAllowsFetchAndMerge:        branch.Node.BranchProtectionRule.LockAllowsFetchAndMerge,
				LockBranch:                     branch.Node.BranchProtectionRule.LockBranch,
				Pattern:                        &branch.Node.BranchProtectionRule.Pattern,
				RequireLastPushApproval:        branch.Node.BranchProtectionRule.RequireLastPushApproval,
				RequiredApprovingReviewCount:   branch.Node.BranchProtectionRule.RequiredApprovingReviewCount,
				RequiredDeploymentEnvironments: branch.Node.BranchProtectionRule.RequiredDeploymentEnvironments,
				RequiredStatusChecks:           branch.Node.BranchProtectionRule.RequiredStatusChecks,
				RequiresApprovingReviews:       branch.Node.BranchProtectionRule.RequiresApprovingReviews,
				RequiresConversationResolution: branch.Node.BranchProtectionRule.RequiresConversationResolution,
				RequiresCodeOwnerReviews:       branch.Node.BranchProtectionRule.RequiresCodeOwnerReviews,
				RequiresCommitSignatures:       branch.Node.BranchProtectionRule.RequiresCommitSignatures,
				RequiresDeployments:            branch.Node.BranchProtectionRule.RequiresDeployments,
				RequiresLinearHistory:          branch.Node.BranchProtectionRule.RequiresLinearHistory,
				RequiresStatusChecks:           branch.Node.BranchProtectionRule.RequiresStatusChecks,
				RequiresStrictStatusChecks:     branch.Node.BranchProtectionRule.RequiresStrictStatusChecks,
				RestrictsPushes:                branch.Node.BranchProtectionRule.RestrictsPushes,
				RestrictsReviewDismissals:      branch.Node.BranchProtectionRule.RestrictsReviewDismissals,
				MatchingBranches:               branch.Node.BranchProtectionRule.MatchingBranches.TotalCount,
			}

			value := models.Resource{
				ID:   id,
				Name: branch.Node.Name,
				Description: model.BranchDescription{
					Name:                 &branch.Node.Name,
					Commit:               commit,
					BranchProtectionRule: branchProtctionRule,
					RepoFullName:         &repoFullName,
					Protected:            protected,
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
		if !query.Repository.Refs.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Repository.Refs.PageInfo.EndCursor)
	}
	return values, nil
}

//TODO: Change the get function api call (The resource type model in Get function is different from model in List function)
//func GetRepositoryBranch(ctx context.Context, githubClient model.GitHubClient, organizationName string, repositoryName string, branchName string, stream *models.StreamSender) (*models.Resource, error) {
//	branchInfo, _, err := githubClient.RestClient.Repositories.GetBranch(ctx, organizationName, repositoryName, branchName, true)
//	if err != nil {
//		return nil, err
//	}
//	repoFullName := formRepositoryFullName(organizationName, repositoryName)
//
//	id := fmt.Sprintf("%s/%s/%s", organizationName, repositoryName, branchInfo.GetName())
//
//	authoredDate := branch.Node.Target.Commit.AuthoredDate.Format(time.RFC3339)
//	authorDate := branch.Node.Target.Commit.Author.Date.Format(time.RFC3339)
//	authorUserCreatedAt := branch.Node.Target.Commit.Author.User.CreatedAt.Format(time.RFC3339)
//	authorUserUpdatedAt := branch.Node.Target.Commit.Author.User.UpdatedAt.Format(time.RFC3339)
//	commiterUserCreatedAt := branch.Node.Target.Commit.Committer.User.CreatedAt.Format(time.RFC3339)
//	commiterUserUpdatedAt := branch.Node.Target.Commit.Committer.User.UpdatedAt.Format(time.RFC3339)
//	committedDate := branch.Node.Target.Commit.CommittedDate.Format(time.RFC3339)
//	commiterDate := branch.Node.Target.Commit.Committer.Date.Format(time.RFC3339)
//
//	authorUser := model.BasicUser{
//		Id:        branch.Node.Target.Commit.Author.User.Id,
//		NodeId:    &branch.Node.Target.Commit.Author.User.NodeId,
//		Name:      &branch.Node.Target.Commit.Author.User.Name,
//		Login:     &branch.Node.Target.Commit.Author.User.Login,
//		Email:     &branch.Node.Target.Commit.Author.User.Email,
//		CreatedAt: &authorUserCreatedAt,
//		UpdatedAt: &authorUserUpdatedAt,
//		Url:       &branch.Node.Target.Commit.Author.User.Url,
//	}
//
//	commiterUser := model.BasicUser{
//		Id:        branch.Node.Target.Commit.Committer.User.Id,
//		NodeId:    &branch.Node.Target.Commit.Committer.User.NodeId,
//		Name:      &branch.Node.Target.Commit.Committer.User.Name,
//		Login:     &branch.Node.Target.Commit.Committer.User.Login,
//		Email:     &branch.Node.Target.Commit.Committer.User.Email,
//		CreatedAt: &commiterUserCreatedAt,
//		UpdatedAt: &commiterUserUpdatedAt,
//		Url:       &branch.Node.Target.Commit.Committer.User.Url,
//	}
//
//	author := model.GitActor{
//		AvatarUrl: &branch.Node.Target.Commit.Author.AvatarUrl,
//		Date:      &authorDate,
//		Email:     &branch.Node.Target.Commit.Author.Email,
//		Name:      &branch.Node.Target.Commit.Author.Name,
//		User:      authorUser,
//	}
//
//	commiter := model.GitActor{
//		AvatarUrl: &branch.Node.Target.Commit.Committer.AvatarUrl,
//		Date:      &commiterDate,
//		Email:     &branch.Node.Target.Commit.Committer.Email,
//		Name:      &branch.Node.Target.Commit.Committer.Name,
//		User:      commiterUser,
//	}
//
//	signature := model.Signature{
//		Email:             &branch.Node.Target.Commit.Signature.Email,
//		IsValid:           branch.Node.Target.Commit.Signature.IsValid,
//		State:             &branch.Node.Target.Commit.Signature.State,
//		WasSignedByGitHub: branch.Node.Target.Commit.Signature.WasSignedByGitHub,
//		Signer: struct {
//			Email *string
//			Login *string
//		}{Email: &branch.Node.Target.Commit.Signature.Signer.Email, Login: &branch.Node.Target.Commit.Signature.Signer.Login},
//	}
//
//	status := model.CommitStatus{
//		State: &branch.Node.Target.Commit.Status.State,
//	}
//
//	commit := model.BaseCommit{
//		Sha:                 &branch.Node.Target.Commit.Sha,
//		ShortSha:            &branch.Node.Target.Commit.ShortSha,
//		AuthoredDate:        &authoredDate,
//		Author:              author,
//		CommittedDate:       &committedDate,
//		Committer:           commiter,
//		Message:             &branch.Node.Target.Commit.Message,
//		Url:                 &branch.Node.Target.Commit.Url,
//		Additions:           branch.Node.Target.Commit.Additions,
//		AuthoredByCommitter: branch.Node.Target.Commit.AuthoredByCommitter,
//		ChangedFiles:        branch.Node.Target.Commit.ChangedFiles,
//		CommittedViaWeb:     branch.Node.Target.Commit.CommittedViaWeb,
//		CommitUrl:           &branch.Node.Target.Commit.CommitUrl,
//		Deletions:           branch.Node.Target.Commit.Deletions,
//		Signature:           signature,
//		TarballUrl:          &branch.Node.Target.Commit.TarballUrl,
//		TreeUrl:             &branch.Node.Target.Commit.TreeUrl,
//		CanSubscribe:        branch.Node.Target.Commit.CanSubscribe,
//		Subscription:        &branch.Node.Target.Commit.Subscription,
//		ZipballUrl:          &branch.Node.Target.Commit.ZipballUrl,
//		MessageHeadline:     &branch.Node.Target.Commit.MessageHeadline,
//		Status:              status,
//		NodeId:              &branch.Node.Target.Commit.NodeId,
//	}
//
//	branchProtctionRule := model.BranchProtectionRule{
//		AllowsDeletions:                branch.Node.BranchProtectionRule.AllowsDeletions,
//		AllowsForcePushes:              branch.Node.BranchProtectionRule.AllowsForcePushes,
//		BlocksCreations:                branch.Node.BranchProtectionRule.BlocksCreations,
//		CreatorLogin:                   &branch.Node.BranchProtectionRule.Creator.Login,
//		Id:                             branch.Node.BranchProtectionRule.Id,
//		NodeId:                         &branch.Node.BranchProtectionRule.NodeId,
//		DismissesStaleReviews:          branch.Node.BranchProtectionRule.DismissesStaleReviews,
//		IsAdminEnforced:                branch.Node.BranchProtectionRule.IsAdminEnforced,
//		LockAllowsFetchAndMerge:        branch.Node.BranchProtectionRule.LockAllowsFetchAndMerge,
//		LockBranch:                     branch.Node.BranchProtectionRule.LockBranch,
//		Pattern:                        &branch.Node.BranchProtectionRule.Pattern,
//		RequireLastPushApproval:        branch.Node.BranchProtectionRule.RequireLastPushApproval,
//		RequiredApprovingReviewCount:   branch.Node.BranchProtectionRule.RequiredApprovingReviewCount,
//		RequiredDeploymentEnvironments: branch.Node.BranchProtectionRule.RequiredDeploymentEnvironments,
//		RequiredStatusChecks:           branch.Node.BranchProtectionRule.RequiredStatusChecks,
//		RequiresApprovingReviews:       branch.Node.BranchProtectionRule.RequiresApprovingReviews,
//		RequiresConversationResolution: branch.Node.BranchProtectionRule.RequiresConversationResolution,
//		RequiresCodeOwnerReviews:       branch.Node.BranchProtectionRule.RequiresCodeOwnerReviews,
//		RequiresCommitSignatures:       branch.Node.BranchProtectionRule.RequiresCommitSignatures,
//		RequiresDeployments:            branch.Node.BranchProtectionRule.RequiresDeployments,
//		RequiresLinearHistory:          branch.Node.BranchProtectionRule.RequiresLinearHistory,
//		RequiresStatusChecks:           branch.Node.BranchProtectionRule.RequiresStatusChecks,
//		RequiresStrictStatusChecks:     branch.Node.BranchProtectionRule.RequiresStrictStatusChecks,
//		RestrictsPushes:                branch.Node.BranchProtectionRule.RestrictsPushes,
//		RestrictsReviewDismissals:      branch.Node.BranchProtectionRule.RestrictsReviewDismissals,
//		MatchingBranches:               branch.Node.BranchProtectionRule.MatchingBranches.TotalCount,
//	}
//
//	value := models.Resource{
//		ID:   id,
//		Name: branchInfo.GetName(),
//		Description: JSONAllFieldsMarshaller{
//			Value: model.BranchDescription{
//				Name: branchInfo.Name,
//				Commit: steampipemodels.BaseCommit{
//					BasicCommit: steampipemodels.BasicCommit{
//						Sha: branchInfo.GetCommit().GetSHA(),
//						Url: branchInfo.GetCommit().GetURL(),
//						Author: steampipemodels.GitActor{
//							Name:      branchInfo.GetCommit().GetAuthor().GetName(),
//							Email:     branchInfo.GetCommit().GetAuthor().GetEmail(),
//							AvatarUrl: branchInfo.GetCommit().GetAuthor().GetAvatarURL(),
//							User: steampipemodels.BasicUser{
//								Login: branchInfo.GetCommit().GetAuthor().GetLogin(),
//								Email: branchInfo.GetCommit().GetAuthor().GetEmail(),
//								Url:   branchInfo.GetCommit().GetAuthor().GetURL(),
//							},
//						},
//						Message: branchInfo.GetCommit().GetCommit().GetMessage(),
//						Committer: steampipemodels.GitActor{
//							Name:      branchInfo.GetCommit().GetCommitter().GetName(),
//							Email:     branchInfo.GetCommit().GetCommitter().GetEmail(),
//							AvatarUrl: branchInfo.GetCommit().GetCommitter().GetAvatarURL(),
//							User: steampipemodels.BasicUser{
//								Login: branchInfo.GetCommit().GetCommitter().GetLogin(),
//								Email: branchInfo.GetCommit().GetCommitter().GetEmail(),
//								Url:   branchInfo.GetCommit().GetCommitter().GetURL(),
//							},
//						},
//						ShortSha: branchInfo.GetCommit().GetCommit().GetSHA(),
//					},
//					Status: steampipemodels.CommitStatus{
//						State: branchInfo.GetCommit().GetStats().String(),
//					},
//				},
//				RepoFullName: repoFullName,
//				Protected:    branchInfo.GetProtected(),
//			},
//		},
//	}
//	if stream != nil {
//		if err := (*stream)(value); err != nil {
//			return nil, err
//		}
//	}
//	return &value, nil
//}
