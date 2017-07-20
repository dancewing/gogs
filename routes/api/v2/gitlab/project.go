package gitlab

import "github.com/gogits/gogs/models"

// Project represents a GitLab project.
//
// GitLab API docs: http://doc.gitlab.com/ce/api/projects.html
type Project struct {
	Id                int64      `json:"id"`
	Name              string     `json:"name"`
	NamespaceWithName string     `json:"name_with_namespace"`
	PathWithNamespace string     `json:"path_with_namespace"`
	Namespace         *Namespace `json:"namespace,nil,omitempty"`
	Description       string     `json:"description"`
	LastModified      string     `json:"last_modified"`
	CreatedAt         string     `json:"created_at"`
	Owner             *User      `json:"owner,nil,omitempty"`
	AvatarUrl         string     `json:"avatar_url,nil,omitempty"`
}

// Namespace represents a GitLab namespace.
//
// GitLab API docs: http://doc.gitlab.com/ce/api/namespaces.html
type Namespace struct {
	Id     int64   `json:"id"`
	Name   string  `json:"name,omitempty"`
	Avatar *Avatar `json:"avatar,nil,omitempty"`
}

// Avatar represents a GitLab avatar.
type Avatar struct {
	Url string `json:"url"`
}

func MapProjectFromGitlab(r *models.Repository) *Project {
	return &Project{
		Id:                r.ID,
		Name:              r.Name,
		NamespaceWithName: getNamespaceWithName(r),
		PathWithNamespace: getPathWithNamespace(r),
		Namespace:         MapNamespaceFromGitlab(r.Owner),
		Description:       r.Description,
		Owner:             MapUserFromGitlab(r.Owner),
		AvatarUrl:         r.HTMLURL(),
	}
}

func getNamespaceWithName(r *models.Repository) string {
	if r.Owner != nil {
		return r.Owner.Name + " / " + r.Name
	}
	return ""
}

func getPathWithNamespace(r *models.Repository) string {
	if r.Owner != nil {
		return r.Owner.Name + "/" + r.Name
	}
	return ""
}
