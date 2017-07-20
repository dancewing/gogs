package gitlab

// File represents uploaded file to gitlab
//
// Gitlab API docs: http://docs.gitlab.com/ee/api/projects.html#upload-a-file
type File struct {
	Alt      string `json:"alt"`
	URL      string `json:"url"`
	Markdown string `json:"markdown"`
}
