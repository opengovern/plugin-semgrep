package describers

import (
	"context"

	"github.com/opengovern/og-describer-template/discovery/pkg/models"
	model "github.com/opengovern/og-describer-template/discovery/provider"
	"github.com/shurcooL/githubv4"
	steampipemodels "github.com/turbot/steampipe-plugin-github/github/models"
)

func GetAllRepositoriesRuleSets(ctx context.Context, githubClient model.GitHubClient, organizationName string, stream *models.StreamSender) ([]models.Resource, error) {
	client := githubClient.RestClient

	var repositoryName string
	if value := ctx.Value(paramKeyRepoName); value != nil {
		repositoryName = value.(string)
	}

	if repositoryName != "" {
		repoValues, err := GetRepositoryRuleSets(ctx, githubClient, stream, organizationName, repositoryName)
		if err != nil {
			return nil, err
		}
		return repoValues, nil
	}

	repositories, err := getRepositories(ctx, client, organizationName)
	if err != nil {
		return nil, nil
	}
	var values []models.Resource
	for _, repo := range repositories {
		repoValues, err := GetRepositoryRuleSets(ctx, githubClient, stream, organizationName, repo.GetName())
		if err != nil {
			return nil, err
		}
		values = append(values, repoValues...)
	}
	return values, nil
}

func GetRepositoryRuleSets(ctx context.Context, githubClient model.GitHubClient, stream *models.StreamSender, owner, repo string) ([]models.Resource, error) {
	client := githubClient.GraphQLClient
	rulesetPageSize := pageSize
	rulePageSize := pageSize
	bypassActorPageSize := pageSize
	var query struct {
		RateLimit  steampipemodels.RateLimit
		Repository struct {
			Rulesets struct {
				PageInfo struct {
					HasNextPage bool
					EndCursor   githubv4.String
				}
				Edges []struct {
					Node struct {
						CreatedAt   githubv4.DateTime
						DatabaseID  int
						Enforcement string
						Name        string
						ID          string
						Rules       struct {
							PageInfo steampipemodels.PageInfo
							Edges    []struct {
								Node steampipemodels.Rule
							}
						} `graphql:"rules(first: $rulePageSize, after: $ruleCursor)"`
						BypassActors struct {
							PageInfo steampipemodels.PageInfo
							Edges    []struct {
								Node steampipemodels.BypassActor
							}
						} `graphql:"bypassActors(first: $bypassActorPageSize, after: $bypassActorCursor)"`
						Conditions steampipemodels.Conditions
					}
				}
			} `graphql:"rulesets(first: $rulesetPageSize, after: $rulesetCursor)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}
	variables := map[string]interface{}{
		"owner":               githubv4.String(owner),
		"name":                githubv4.String(repo),
		"rulesetPageSize":     githubv4.Int(rulesetPageSize),
		"rulesetCursor":       (*githubv4.String)(nil),
		"rulePageSize":        githubv4.Int(rulePageSize),
		"ruleCursor":          (*githubv4.String)(nil),
		"bypassActorPageSize": githubv4.Int(bypassActorPageSize),
		"bypassActorCursor":   (*githubv4.String)(nil),
	}
	var ruleSets []steampipemodels.Ruleset
	repoFullName := formRepositoryFullName(owner, repo)
	var values []models.Resource
	for {
		err := client.Query(ctx, &query, variables)
		if err != nil {
			return nil, err
		}
		for _, edge := range query.Repository.Rulesets.Edges {
			var rules []steampipemodels.Rule
			for _, rule := range edge.Node.Rules.Edges {
				rules = append(rules, rule.Node)
			}
			if edge.Node.Rules.PageInfo.HasNextPage {
				additionalRules := getAdditionalRules(ctx, client, edge.Node.DatabaseID, owner, repo, "")
				rules = append(rules, additionalRules...)
			}
			var bypassActors []steampipemodels.BypassActor
			for _, actor := range edge.Node.BypassActors.Edges {
				bypassActors = append(bypassActors, actor.Node)
			}
			if edge.Node.BypassActors.PageInfo.HasNextPage {
				additionalBypassActors := getAdditionalBypassActors(ctx, client, owner, repo, edge.Node.DatabaseID, "")
				bypassActors = append(bypassActors, additionalBypassActors...)
			}
			ruleset := steampipemodels.Ruleset{
				CreatedAt:    edge.Node.CreatedAt.String(),
				DatabaseID:   edge.Node.DatabaseID,
				Enforcement:  edge.Node.Enforcement,
				Name:         edge.Node.Name,
				ID:           edge.Node.ID,
				Rules:        rules,
				BypassActors: bypassActors,
				Conditions:   edge.Node.Conditions,
			}
			ruleSets = append(ruleSets, ruleset)
		}
		for _, ruleset := range ruleSets {
			value := models.Resource{
				ID:   ruleset.ID,
				Name: ruleset.Name,
				Description: model.RepoRuleSetDescription{
					Ruleset:      ruleset,
					RepoFullName: repoFullName,
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
		if !query.Repository.Rulesets.PageInfo.HasNextPage {
			break
		}
		variables["rulesetCursor"] = githubv4.NewString(query.Repository.Rulesets.PageInfo.EndCursor)
	}
	return values, nil
}

func getAdditionalRules(ctx context.Context, client *githubv4.Client, databaseID int, owner string, repo string, initialCursor githubv4.String) []steampipemodels.Rule {
	var query struct {
		RateLimit  steampipemodels.RateLimit
		Repository struct {
			Ruleset struct {
				Rules struct {
					PageInfo struct {
						HasNextPage bool
						EndCursor   githubv4.String
					}
					Edges []struct {
						Node steampipemodels.Rule
					}
				} `graphql:"rules(first: $pageSize, after: $cursor)"`
			} `graphql:"ruleset(databaseId: $databaseID)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}
	variables := map[string]interface{}{
		"pageSize":   githubv4.Int(100),
		"cursor":     githubv4.NewString(initialCursor),
		"databaseID": githubv4.Int(databaseID),
		"owner":      githubv4.String(owner),
		"name":       githubv4.String(repo),
	}
	var rules []steampipemodels.Rule
	for {
		err := client.Query(ctx, &query, variables)
		if err != nil {
			return nil
		}
		for _, edge := range query.Repository.Ruleset.Rules.Edges {
			rules = append(rules, edge.Node)
		}
		if !query.Repository.Ruleset.Rules.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Repository.Ruleset.Rules.PageInfo.EndCursor)
	}
	return rules
}

func getAdditionalBypassActors(ctx context.Context, client *githubv4.Client, owner string, repo string, databaseID int, initialCursor githubv4.String) []steampipemodels.BypassActor {
	var query struct {
		RateLimit  steampipemodels.RateLimit
		Repository struct {
			Ruleset struct {
				BypassActors struct {
					PageInfo struct {
						HasNextPage bool
						EndCursor   githubv4.String
					}
					Edges []struct {
						Node steampipemodels.BypassActor
					}
				} `graphql:"bypassActors(first: $pageSize, after: $cursor)"`
			} `graphql:"ruleset(databaseId: $databaseID)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}
	variables := map[string]interface{}{
		"owner":      githubv4.String(owner),
		"name":       githubv4.String(repo),
		"pageSize":   githubv4.Int(100),
		"cursor":     githubv4.NewString(initialCursor),
		"databaseID": githubv4.Int(databaseID),
	}
	var bypassActors []steampipemodels.BypassActor
	for {
		err := client.Query(ctx, &query, variables)
		if err != nil {
			return nil
		}
		for _, edge := range query.Repository.Ruleset.BypassActors.Edges {
			bypassActors = append(bypassActors, edge.Node)
		}
		if !query.Repository.Ruleset.BypassActors.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Repository.Ruleset.BypassActors.PageInfo.EndCursor)
	}
	return bypassActors
}
