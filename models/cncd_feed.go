package models

// Feed represents an item in the user's feed or timeline.
//
// swagger:model feed
type Feed struct {
	Owner    string `json:"owner"         `
	Name     string `json:"name"          `
	FullName string `json:"full_name"     `

	Number   int    `json:"number,omitempty"        `
	Event    string `json:"event,omitempty"         `
	Status   string `json:"status,omitempty"        `
	Created  int64  `json:"created_at,omitempty"    `
	Started  int64  `json:"started_at,omitempty"    `
	Finished int64  `json:"finished_at,omitempty"   `
	Commit   string `json:"commit,omitempty"        `
	Branch   string `json:"branch,omitempty"        `
	Ref      string `json:"ref,omitempty"           `
	Refspec  string `json:"refspec,omitempty"       `
	Remote   string `json:"remote,omitempty"        `
	Title    string `json:"title,omitempty"         `
	Message  string `json:"message,omitempty"       `
	Author   string `json:"author,omitempty"        `
	Avatar   string `json:"author_avatar,omitempty" `
	Email    string `json:"author_email,omitempty"  `
}

func (t Feed) TableName() string {
	return "cncd_feed"
}
