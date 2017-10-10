package api

import "github.com/gogits/gogs/models"

type User struct {
	Id          int64    `json:"id"`
	Name        string   `json:"login"`
	IsActive    bool     `json:"activated"`
	AvatarUrl   string   `json:"avatarUrl"`
	FirstName   string   `json:"firstName"`
	LastName    string   `json:"lastName"`
	Email       string   `json:"email"`
	LangKey     string   `json:"langKey"`
	Authorities []string `json:"authorities"`
	IsOrg       bool     `json:"isOrg"`
}

type Account struct {
	Id          int64    `json:"id"`
	Name        string   `json:"login"`
	IsActive    bool     `json:"activated"`
	AvatarUrl   string   `json:"avatarUrl"`
	FirstName   string   `json:"firstName"`
	LastName    string   `json:"lastName"`
	Email       string   `json:"email"`
	LangKey     string   `json:"langKey"`
	Authorities []string `json:"authorities"`
	IsOrg       bool     `json:"isOrg"`
}

// mapUserFromGitlab mapped data from gitlab user to kanban user
func ConvertUserToAPI(u *models.User) *User {
	if u == nil {
		return nil
	}
	user := &User{
		Id:        u.ID,
		Name:      u.Name,
		FirstName: u.FullName,
		LastName:  u.FullName,
		AvatarUrl: u.AvatarLink(),
		IsActive:  u.IsActive,
		Email:     u.Email,
		LangKey:   u.LangKey,
		IsOrg:     u.IsOrganization(),
	}

	var authorities []string
	if u.IsAdmin {
		authorities = append(authorities, "ROLE_ADMIN")
	} else {
		authorities = append(authorities, "ROLE_USER")
	}

	user.Authorities = authorities

	return user
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
