package github

import (
	opengovernance "github.com/opengovern/og-describer-template/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func gitHubIssueColumns() []*plugin.Column {
	tableCols := []*plugin.Column{
		{
			Name:        "repository_full_name",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.RepositoryFullName"),
			Description: "The full name of the repository (login/repo-name).",
		},
	}

	return append(tableCols, sharedIssueColumns()...)
}

func sharedIssueColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "number",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromField("Description.Number"),
			Description: "The issue number."},
		{Name: "id",
			Type:      proto.ColumnType_INT,
			Transform: transform.FromField("Description.Id"),

			Description: "The ID of the issue."},
		{
			Name:        "node_id",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.NodeId"),
			Description: "The node ID of the issue."},
		{
			Name:        "active_lock_reason",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.ActiveLockReason"),
			Description: "Reason that the conversation was locked."},
		{
			Name:        "author",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("Description.Author"),
			Description: "The actor who authored the issue."},
		{
			Name:        "author_login",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.AuthorLogin"),
			Description: "The login of the issue author."},
		{
			Name:        "author_association",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.AuthorAssociation"),
			Description: "Author's association with the subject of the issue."},
		{
			Name:        "body",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.Body"),
			Description: "Identifies the body of the issue."},
		{
			Name:        "body_url",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.BodyUrl"),
			Description: "URL for this issue body."},
		{
			Name:        "closed",
			Type:        proto.ColumnType_BOOL,
			Transform:   transform.FromField("Description.Closed"),
			Description: "If true, issue is closed."},
		{
			Name:        "closed_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Transform:   transform.FromField("Description.ClosedAt").NullIfZero().Transform(convertTimestamp),
			Description: "Timestamp when issue was closed."},
		{
			Name:        "created_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Transform:   transform.FromField("Description.CreatedAt").NullIfZero().Transform(convertTimestamp),
			Description: "Timestamp when issue was created."},
		{
			Name:        "created_via_email",
			Type:        proto.ColumnType_BOOL,
			Transform:   transform.FromField("Description.CreatedViaEmail"),
			Description: "If true, issue was created via email."},
		{
			Name:        "editor",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("Description.Editor"),
			Description: "The actor who edited the issue."},
		{
			Name:        "full_database_id",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromField("Description.FullDatabaseId"),
			Description: "Identifies the primary key from the database as a BigInt."},
		{
			Name:        "includes_created_edit",
			Type:        proto.ColumnType_BOOL,
			Transform:   transform.FromField("Description.IncludesCreatedEdit"),
			Description: "If true, issue was edited and includes an edit with the creation data."},
		{
			Name:        "is_pinned",
			Type:        proto.ColumnType_BOOL,
			Transform:   transform.FromField("Description.IsPinned"),
			Description: "if true, this issue is currently pinned to the repository issues list."},
		{
			Name:        "is_read_by_user",
			Type:        proto.ColumnType_BOOL,
			Transform:   transform.FromField("Description.IsReadByUser"),
			Description: "if true, this issue has been read by the user."},
		{
			Name:        "last_edited_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Transform:   transform.FromField("Description.LastEditedAt").NullIfZero().Transform(convertTimestamp),
			Description: "Timestamp when issue was last edited."},
		{
			Name:        "locked",
			Type:        proto.ColumnType_BOOL,
			Transform:   transform.FromField("Description.Locked"),
			Description: "If true, issue is locked."},
		{
			Name:        "milestone",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("Description.Milestone"),
			Description: "The milestone associated with the issue."},
		{
			Name:        "published_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Transform:   transform.FromField("Description.PublishedAt").NullIfZero().Transform(convertTimestamp),
			Description: "Timestamp when issue was published."},
		{
			Name:        "state",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.State"),
			Description: "The state of the issue."},
		{
			Name:        "state_reason",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.StateReason"),
			Description: "The reason for the issue state."},
		{
			Name:        "title",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.Title"),
			Description: "The title of the issue."},
		{
			Name:        "updated_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Transform:   transform.FromField("Description.UpdatedAt").NullIfZero().Transform(convertTimestamp),
			Description: "Timestamp when issue was last updated."},
		{
			Name:        "url",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Description.Url"),
			Description: "URL for the issue."},
		{
			Name:        "assignees_total_count",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromField("Description.AssigneesTotalCount"),
			Description: "Count of assignees on the issue."},
		{
			Name:        "comments_total_count",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromField("Description.CommentsTotalCount"),
			Description: "Count of comments on the issue."},
		{
			Name:        "labels_total_count",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromField("Description.LabelsTotalCount"),
			Description: "Count of labels on the issue."},
		{
			Name:        "labels_src",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("Description.LabelsSrc"),
			Description: "The first 100 labels associated to the issue."},
		{
			Name: "labels",
			Type: proto.ColumnType_JSON, Description: "A map of labels for the issue.",
			Transform: transform.FromField("Description.Labels"),
		},
		{Name: "user_can_close",
			Type:        proto.ColumnType_BOOL,
			Hydrate:     issueHydrateUserCanClose,
			Transform:   transform.FromField("Description.UserCanClose"),
			Description: "If true, user can close the issue."},
		{Name: "user_can_react",
			Type:      proto.ColumnType_BOOL,
			Transform: transform.FromField("Description.UserCanReact"),

			Description: "If true, user can react on the issue."},
		{Name: "user_can_reopen",
			Type:      proto.ColumnType_BOOL,
			Transform: transform.FromField("Description.UserCanReopen"),

			Description: "If true, user can reopen the issue."},
		{Name: "user_can_subscribe",
			Type:      proto.ColumnType_BOOL,
			Transform: transform.FromField("Description.UserCanSubscribe"),

			Description: "If true, user can subscribe to the issue."},
		{Name: "user_can_update",
			Type:      proto.ColumnType_BOOL,
			Transform: transform.FromField("Description.UserCanUpdate"),

			Description: "If true, user can update the issue,"},
		{Name: "user_cannot_update_reasons",
			Type:      proto.ColumnType_JSON,
			Transform: transform.FromField("Description.UserCannotUpdateReasons"),

			Description: "A list of reason why user cannot update the issue."},
		{Name: "user_did_author",
			Type:      proto.ColumnType_BOOL,
			Transform: transform.FromField("Description.UserDidAuthor"),

			Description: "If true, user authored the issue."},
		{Name: "user_subscription",
			Type:      proto.ColumnType_STRING,
			Transform: transform.FromField("Description.UserSubscription"),

			Description: "Subscription state of the user to the issue."},
		{Name: "assignees",
			Type:      proto.ColumnType_JSON,
			Transform: transform.FromField("Description.Assignees"),

			Description: "A list of Users assigned to the issue."},
	}
}

func tableGitHubIssue() *plugin.Table {
	return &plugin.Table{
		Name:        "github_issue",
		Description: "GitHub Issues are used to track ideas, enhancements, tasks, or bugs for work on GitHub.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListIssue,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"repository_full_name", "number"}),
			Hydrate:    opengovernance.GetIssue,
		},
		Columns: commonColumns(gitHubIssueColumns()),
	}
}
