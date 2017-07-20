package gitlab

import (
	"time"

	"github.com/gogits/gogs/models"
)

// Issue represents a GitLab issue.
//
// GitLab API docs: http://doc.gitlab.com/ce/api/issues.html
type Issue struct {
	Assignee       *User      `json:"assignee"`
	Author         *User      `json:"author"`
	Description    string     `json:"description"`
	Milestone      *Milestone `json:"milestone"`
	Id             int64      `json:"id"`
	Iid            int64      `json:"iid"`
	Labels         *[]string  `json:"labels"`
	ProjectId      int64      `json:"project_id"`
	State          string     `json:"state"`
	Title          string     `json:"title"`
	UserNotesCount int        `json:"user_notes_count"`
	Subscribed     bool       `json:"subscribed"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DueDate        string     `json:"due_date"`
	Confidential   bool       `json:"confidential"`
	WebURL         string     `json:"web_url"`
}

// IssueRequest represents the available CreateIssue() and UpdateIssue() options.
//
// GitLab API docs: http://doc.gitlab.com/ce/api/issues.html#new-issues
type IssueRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	AssigneeId  *int64 `json:"assignee_id"`
	MilestoneId *int64 `json:"milestone_id"`
	Labels      string `json:"labels"`
	StateEvent  string `json:"state_event,omitempty"`
	DueDate     string `json:"due_date,omitempty"`
}

// MoveIssueRequest moved issue to another project
type MoveIssueRequest struct {
	ToProjectID string `json:"to_project_id"`
}

func MapIssueFromGitlab(n *models.Issue) *Issue {
	if n == nil {
		return nil
	}
	return &Issue{
		Id:          n.ID,
		Assignee:    MapUserFromGitlab(n.Assignee),
		Author:      MapUserFromGitlab(n.Poster),
		Description: n.Content,
		Milestone:   MapMilestoneFromGithub(n.Milestone),
		Title:       n.Title,
		//ProjectId:   n.Repo.RepoID,
	}
}
