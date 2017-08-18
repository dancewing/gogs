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
	ID     int64    `json:"-"      `
	RepoID int64    `json:"-"      `
	Login  string   `json:"login"  `
	Allow  bool     `json:"allow"  `
	Block  bool     `json:"block"  `
	Branch []string `json:"branch" xorm:"-"`
	Deploy []string `json:"deploy" xorm:"-"`
	Event  []string `json:"event"  xorm:"-"`
}


func (t Sender) TableName() string {
	return "cncd_sender"
}
