package gitlab

import "github.com/gogits/gogs/models"

type User struct {
	Id        int64
	Name      string
	IsAdmin   bool
	AvatarUrl string
	State     string
	Username  string
	Passwd    string
	Salt      string
	Email     string
}

// mapUserFromGitlab mapped data from gitlab user to kanban user
func MapUserFromGogs(u *models.User) *User {
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

func MapNamespaceFromGogs(n *models.User) *Namespace {
	if n == nil {
		return nil
	}
	return &Namespace{
		Id:     n.ID,
		Name:   n.Name,
		Avatar: MapAvatarFromGogs(n),
	}
}

// mapAvatarFromGitlab transform gitlab avatar to kanban avatar
func MapAvatarFromGogs(n *models.User) *Avatar {
	if n == nil {
		return nil
	}
	return &Avatar{
		Url: n.AvatarLink(),
	}
}
