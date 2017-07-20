package gitlab

import "github.com/gogits/gogs/models"

// User represents a GitLab user.
//
// GitLab API docs: http://doc.gitlab.com/ce/api/users.html
type User struct {
	Id           int64  `json:"id"`
	Name         string `json:"name,omitempty"`
	AvatarUrl    string `json:"avatar_url,nil,omitempty"`
	State        string `json:"state,omitempty"`
	Username     string `json:"username,omitempty"`
	WebUrl       string `json:"web_url,omitempty"`
	PrivateToken string `json:"private_token"`
}

// mapUserFromGitlab mapped data from gitlab user to kanban user
func MapUserFromGitlab(u *models.User) *User {
	if u == nil {
		return nil
	}
	return &User{
		Id:        u.ID,
		Name:      u.Name,
		Username:  u.Name,
		AvatarUrl: u.AvatarLink(),
		State:     "State",
	}
}

func MapNamespaceFromGitlab(n *models.User) *Namespace {
	if n == nil {
		return nil
	}
	return &Namespace{
		Id:     n.ID,
		Name:   n.Name,
		Avatar: MapAvatarFromGitlab(n),
	}
}

// mapAvatarFromGitlab transform gitlab avatar to kanban avatar
func MapAvatarFromGitlab(n *models.User) *Avatar {
	if n == nil {
		return nil
	}
	return &Avatar{
		Url: n.AvatarLink(),
	}
}
