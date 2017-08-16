package models

type SenderService interface {
	SenderAllowed(*User, *Repository, *Build, *Config) (bool, error)
	SenderCreate(*Repository, *Sender) error
	SenderUpdate(*Repository, *Sender) error
	SenderDelete(*Repository, string) error
	SenderList(*Repository) ([]*Sender, error)
}

type SenderStore interface {
	SenderFind(*Repository, string) (*Sender, error)
	SenderList(*Repository) ([]*Sender, error)
	SenderCreate(*Sender) error
	SenderUpdate(*Sender) error
	SenderDelete(*Sender) error
}

type Sender struct {
	ID     int64    `json:"-"      meddler:"sender_id,pk"`
	RepoID int64    `json:"-"      meddler:"sender_repo_id"`
	Login  string   `json:"login"  meddler:"sender_login"`
	Allow  bool     `json:"allow"  meddler:"sender_allow"`
	Block  bool     `json:"block"  meddler:"sender_block"`
	Branch []string `json:"branch" meddler:"-"`
	Deploy []string `json:"deploy" meddler:"-"`
	Event  []string `json:"event"  meddler:"-"`
}
