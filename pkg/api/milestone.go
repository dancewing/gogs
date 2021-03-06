package api

import "github.com/gogits/gogs/models"

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

type MilestoneRequest struct {
	MilestoneID int64  `json:"milestone_id"`
	ProjectID   int64  `json:"project_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
}

func MapMilestoneFromGogs(m *models.Milestone) *Milestone {

	if m == nil {
		return nil
	}

	return &Milestone{
		ID:  m.ID,
		IID: m.ID,
		//State: m.State().(string),
		Title:       m.Name,
		DueDate:     m.DeadlineString,
		Description: "",
	}
}
