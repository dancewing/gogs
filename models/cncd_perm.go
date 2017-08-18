package models

// PermStore persists repository permissions information to storage.
type PermStore interface {
	PermFind(user *User, repo *Repository) (*Perm, error)
	PermUpsert(perm *Perm) error
	PermBatch(perms []*Perm) error
	PermDelete(perm *Perm) error
	PermFlush(user *User, before int64) error
}

// Perm defines a repository permission for an individual user.
type Perm struct {
	UserID int64  `json:"-"      `
	RepoID int64  `json:"-"      `
	Repo   string `json:"-"      xorm:"-"`
	Pull   bool   `json:"pull"   `
	Push   bool   `json:"push"   `
	Admin  bool   `json:"admin"  `
	Synced int64  `json:"synced" `
}


func (t Perm) TableName() string {
	return "cncd_perm"
}
