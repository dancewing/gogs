package gitlab

import "time"

// Comment represents a GitLab note.
//
// GitLab API docs: http://doc.gitlab.com/ce/api/notes.html
type Comment struct {
	Id        int64     `json:"id"`
	Author    *User     `json:"author"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	System    bool      `json:"system"`
	Upvote    bool      `json:"upvote"`
	Downvote  bool      `json:"downvote"`
}

// commentSlice represents list comments for usage sort.Interface
type commentSlice []*Comment

// CommentRequest represents the available CreateComment() and UpdateComment()
// options.
//
// GitLab API docs:
// http://doc.gitlab.com/ce/api/notes.html#create-new-issue-note
type CommentRequest struct {
	Body string `json:"body"`
}
