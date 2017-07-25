package gitlab

import (
	"regexp"
	"strconv"

	"github.com/gogits/gogs/models"
)

// Label represent label
type Label struct {
	ID                    int64  `json:"id"`
	Color                 string `json:"color"`
	Name                  string `json:"name"`
	Description           string `json:"description"`
	OpenCardCount         int    `json:"open_card_count"`
	ClosedCardCount       int    `json:"closed_card_count"`
	OpenMergeRequestCount int    `json:"open_merge_requests_count"`
	Subscribed            bool   `json:"subscribed"`
	Priority              int    `json:"priority"`
}

// Stage represent board stage
type Stage struct {
	Name     string
	Position int
}

// LabelRequest represent request for update label
type LabelRequest struct {
	Name    string `json:"name"`
	Color   string `json:"color"`
	NewName string `json:"new_name"`
}

var (
	stageReg = regexp.MustCompile(`KB\[stage\]\[(\d+)\]\[(.*)\]`)
)

// ParseLabelToStage transform label to stage
func ParseLabelToStage(l string) *Stage {
	m := stageReg.MatchString(l)

	var s Stage
	if m {
		an := stageReg.FindStringSubmatch(l)
		s.Position, _ = strconv.Atoi(an[1])
		s.Name = an[2]
	}

	return &s
}

func MapLabelFromGogs(l *models.Label) *Label {
	return &Label{
		ID:                    l.ID,
		Color:                 l.Color,
		Name:                  l.Name,
		Description:           "",
		OpenCardCount:         0,
		ClosedCardCount:       0,
		OpenMergeRequestCount: 0,
		Subscribed:            false,
		Priority:              0,
	}
}

func MapLabelsFromGogs(labels []*models.Label) *[]string {

	result := make([]string, len(labels))

	for i := range labels {
		result[i] = labels[i].Name
	}

	return &result
}
