package gitlab

import "github.com/gogits/gogs/models"

// Board represents a kanban board.
type Board struct {
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

type BoardRequest struct {
	BoardId string `json:"project_id"`
}

// Board represents a namespace kanban board.
type Namespace struct {
	Id     int64   `json:"id"`
	Name   string  `json:"name,omitempty"`
	Avatar *Avatar `json:"avatar,nil,omitempty"`
}

// Avatar represent a Avatar url.
type Avatar struct {
	Url string `json:"url"`
}

func MapBoardFromGogs(r *models.Repository) *Board {
	return &Board{
		Id:                r.ID,
		Name:              r.Name,
		NamespaceWithName: getNamespaceWithName(r),
		PathWithNamespace: getPathWithNamespace(r),
		Namespace:         MapNamespaceFromGogs(r.Owner),
		Description:       r.Description,
		Owner:             MapUserFromGogs(r.Owner),
		AvatarUrl:         MapAvatarFromGogs(r.Owner).Url,
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
