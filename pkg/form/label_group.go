package form

type BatchUpdateLabel struct {
	Labels []*CreateLabel `json:"labels"`
}
