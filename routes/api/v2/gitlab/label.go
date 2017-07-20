package gitlab

import "github.com/gogits/gogs/models"

// Label represents a GitLab label.
//
// GitLab API docs: http://doc.gitlab.com/ce/api/labels.html
type Label struct {
	ID                    int64  `json:"id"`
	Color                 string `json:"color"`
	Name                  string `json:"name"`
	Description           string `json:"description"`
	OpenIssueCount        int    `json:"open_issue_count"`
	ClosedIssueCount      int    `json:"closed_issue_count"`
	OpenMergeRequestCount int    `json:"open_merge_requests_count"`
	Subscribed            bool   `json:"subscribed"`
	Priority              int    `json:"priority"`
}

// LabelRequest represents the available CreateLabel() and UpdateLabel() options.
type LabelRequest struct {
	Color   string `json:"color"`
	Name    string `json:"name"`
	NewName string `json:"new_name,omitempty"`
}

type LabelDeleteOptions struct {
	Name string `url:"name,omitempty"`
}

func MapLabelFromGitlab(l *models.Label) *Label {
	return &Label{
		ID:                    l.ID,
		Color:                 l.Color,
		Name:                  l.Name,
		Description:           "",
		OpenIssueCount:        0,
		ClosedIssueCount:      0,
		OpenMergeRequestCount: 0,
		Subscribed:            false,
		Priority:              0,
	}
}
