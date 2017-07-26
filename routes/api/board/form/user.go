package form

import (
	"github.com/gogits/gogs/models"
	"github.com/gogits/gogs/pkg/context"
)

type User struct {
	Id        int64       `json:"id"`
	Name      string      `json:"name"`
	IsLogged  bool        `json:"isLogged"`
	IsActive  bool        `json:"isActive"`
	IsAdmin   bool        `json:"isAdmin"`
	AvatarUrl string      `json:"avatarUrl"`
	State     string      `json:"state"`
	Username  string      `json:"username"`
	Repo      *Repository `json:"repo"`
}

type Repository struct {
	AccessMode models.AccessMode `json:"access_mode"`
	Owner      *User             `json:"owner"`
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
		IsAdmin:   u.IsAdmin,
		IsActive:  u.IsActive,
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

func MapRepoFromGogs(r *context.Repository) *Repository {
	if r == nil {
		return nil
	}
	return &Repository{
		AccessMode: r.AccessMode,
		Owner:      MapUserFromGogs(r.Owner),
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
