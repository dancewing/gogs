package form

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/gogits/gogs/models"
)

// Card represents an card in kanban board
type Card struct {
	Id                int64        `json:"id"`
	Iid               int64        `json:"iid"`
	Assignee          *User        `json:"assignee"`
	Milestone         *Milestone   `json:"milestone"`
	Author            *User        `json:"author"`
	Description       string       `json:"description"`
	Labels            *[]string    `json:"labels"`
	ProjectId         int64        `json:"project_id"`
	BoardID           string       `json:"board_id"`
	PathWithNamespace string       `json:"path_with_namespace"`
	Properties        *Properties  `json:"properties"`
	State             string       `json:"state"`
	Title             string       `json:"title"`
	Todo              []*Todo      `json:"todo"`
	TodoMetrics       *TodoMetrics `json:"todo_metrics"`
	UserCommentsCount int          `json:"user_comments_count"`
	Subscribed        bool         `json:"subscribed"`
	CreatedAt         int64        `json:"created_at"`
	UpdatedAt         int64        `json:"updated_at"`
	DueDate           string       `json:"due_date"`
	Confidential      bool         `json:"confidential"`
	WebURL            string       `json:"web_url"`
}

// TodoMetrics represents metrics by todo for card
type TodoMetrics struct {
	Checked  int `json:"checked"`
	Quantity int `json:"quantity"`
}

// Properties represents a card properties
type Properties struct {
	Andon string `json:"andon"`
}

// Todo represents an todo an card
type Todo struct {
	Body    string `json:"body"`
	Checked bool   `json:"checked"`
}

// CardRequest represents a card request for create, update, delete card on kanban
type CardRequest struct {
	CardId       int64             `json:"issue_id"`
	ProjectId    int64             `json:"project_id"`
	Title        string            `json:"title"`
	Description  string            `json:"description"`
	AssigneeId   int64             `json:"assignee_id"`
	MilestoneId  int64             `json:"milestone_id"`
	Labels       string            `json:"labels"`
	Properties   *Properties       `json:"properties"`
	Stage        map[string]string `json:"stage"`
	Todo         []*Todo           `json:"todo"`
	DueDate      string            `json:"due_date"`
	Confidential bool              `json:"confidential"`
}

func (c *Card) RoutingKey() string {
	return fmt.Sprintf("kanban.%d", c.ProjectId)
}

// mapCardFromGitlab mapped gitlab issue to kanban card
func MapCardFromGogs(c *models.Issue) *Card {
	card := Card{
		Id:          c.ID,
		Iid:         c.Index,
		Title:       c.Title,
		Assignee:    MapUserFromGogs(c.Assignee),
		Author:      MapUserFromGogs(c.Poster),
		Description: mapCardDescriptionFromGitlab(c.Content),
		Milestone:   MapMilestoneFromGogs(c.Milestone),
		Labels:      MapLabelsFromGogs(c.Labels),
		ProjectId:   c.Repo.ID,
		Properties:  mapCardPropertiesFromGitlab(c.Content),
		BoardID:     fmt.Sprintf("%d", c.Repo.ID),
		///Todo:              mapCardTodoFromGitlab(c.Content),
		PathWithNamespace: getPathWithNamespace(c.Repo),
		UserCommentsCount: c.NumComments,
		//	Subscribed:        c.Priority,
		CreatedAt: c.Created.Unix(),
		UpdatedAt: c.Updated.Unix(),
		//	DueDate:           c.DueDate,
		//	Confidential:      c.Confidential,
		//	WebURL:            c.WebURL,
	}

	if c.IsClosed {
		card.State = "closed"
	} else {
		card.State = "open"
	}

	//card.TodoMetrics = mapCardTodoMetrics(card.Todo)

	return &card

}

// removeDuplicates removed duplicates
func removeDuplicates(xs *[]string) *[]string {
	found := make(map[string]bool)
	j := 0
	for i, x := range *xs {
		if !found[x] {
			found[x] = true
			(*xs)[j] = (*xs)[i]
			j++
		}
	}
	*xs = (*xs)[:j]

	return xs
}

// mapCardDescriptionFromGitlab clears gitlab description to card description
func mapCardDescriptionFromGitlab(d string) string {
	var r string
	r = regTodo.ReplaceAllString(d, "")
	r = regProp.ReplaceAllString(r, "")
	return strings.TrimSpace(r)
}

// mapCardPropertiesFromGitlab transforms gitlab description to card properties
func mapCardPropertiesFromGitlab(d string) *Properties {
	m := regProp.MatchString(d)
	var p Properties

	if m {
		an := regProp.FindStringSubmatch(d)
		json.Unmarshal([]byte(an[1]), &p)
	}

	return &p
}

// mapCardTodoMetrics maps card todo metrics from gitlab
//func mapCardTodoMetrics(to []*models.Todo) *TodoMetrics {
//	m := models.TodoMetrics{
//		Checked:  0,
//		Quantity: 0,
//	}
//	for _, t := range to {
//		if t.Checked {
//			m.Checked++
//		}
//		m.Quantity++
//	}
//
//	return &m
//}

var (
	regTodo = regexp.MustCompile(`[-\*]{1}\s(?P<checked>\[.\])(?P<body>.*)`)
	regProp = regexp.MustCompile(`<!--\s@KB:(.*?)\s-->`)
)
