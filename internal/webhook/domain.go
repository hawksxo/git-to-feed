package webhook

type RepositoryInfo struct {
	Action string `json:"action"`
	Name   string `json:"name"`
	User   string `json:"user"`
	Description string `json:"description"`
	URL    string `json:"url"`
}
