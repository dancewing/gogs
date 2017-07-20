package gitlab

import "github.com/gogits/gogs/models"

// Milestone represents a GitLab milestone.
//
// GitLab API docs: http://doc.gitlab.com/ce/api/branches.html
type Milestone struct {
	ID          int64  `json:"id"`
	IID         int64  `json:"iid"`
	State       string `json:"state,omitempty"`
	Title       string `json:"title,omitempty"`
	DueDate     string `json:"due_date,omitempty"`
	Description string `json:"description"`
	UpdatedAt   string `json:"updated_at"`
	CreatedAt   string `json:"created_at"`
}

// MilestoneRequest represents the available CreateMilestone() and UpdateMilestone() options.
//
// GitLab API docs: http://doc.gitlab.com/ce/api/milestones.html#create-new-milestone
type MilestoneRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
}

func MapMilestoneFromGithub(m *models.Milestone) *Milestone {

	if m == nil {
		return nil
	}

	return &Milestone{
		ID:  m.ID,
		IID: m.RepoID,
		//State: m.State().(string),
		Title: m.Name,
	}
}
