package github

import (
	"context"

	opengovernance "github.com/opengovern/og-describer-template/discovery/pkg/es"

	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubRepositoryRuleset() *plugin.Table {
	return &plugin.Table{
		Name:        "github_repository_ruleset",
		Description: "Retrieve the rulesets of a specified GitHub repository.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListRepoRuleSet,
		},
		Columns: commonColumns(gitHubRulesetColumns()),
	}
}

func gitHubRulesetColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "repository_full_name",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.RepoFullName"),
			Description: "Full name of the repository that contains the ruleset."},
		{Name: "name",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.Name"),
			Description: "The name of the ruleset."},
		{Name: "id",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.ID"),
			Description: "The ID of the ruleset."},
		{Name: "created_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Transform:   transform.FromField("Description.CreatedAt").NullIfZero().Transform(convertRulesetTimestamp),
			Description: "The date and time when the ruleset was created."},
		{Name: "database_id",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromField("Description.DatabaseID"),
			Description: "The database ID of the ruleset."},
		{Name: "enforcement",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.Enforcement"),
			Description: "The enforcement level of the ruleset."},
		{Name: "rules",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("Description.Rules"),
			Description: "The list of rules in the ruleset."},
		{Name: "bypass_actors",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("Description.BypassActors"),
			Description: "The list of actors who can bypass the ruleset."},
		{Name: "conditions",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("Description.Conditions"),
			Description: "The conditions under which the ruleset applies."},
	}
}

//// TRANSFORM FUNCTION

// The timestamp value we are receiving has the layout "2024-06-11 13:18:48 +0000 UTC".
// Our generic timestamp function does not support converting this specific layout to the desired format.
// Additionally, it is not feasible to create a generic function that handles all possible timestamp layouts.
// Therefore, we have opted to implement a specific timestamp conversion function for this table only.
func convertRulesetTimestamp(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}
	t := d.Value.(string)

	// Parse the timestamp into a time.Time object
	parsedTime, err := time.Parse("2006-01-02 15:04:05 -0700 MST", t)
	if err != nil {
		plugin.Logger(ctx).Error("Error parsing time:", err)
		return nil, err
	}
	// Format the time.Time object to RFC 3339 format
	rfc3339Time := parsedTime.Format(time.RFC3339)

	return rfc3339Time, nil
}
