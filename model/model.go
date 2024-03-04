package model

type ObjectType string

const (
	ProjectObjectType ObjectType = "project"
	GroupObjectType   ObjectType = "group"
)

type ApplicationType string

const (
	GitApplicationType  ApplicationType = "git"
	HelmApplicationType ApplicationType = "helm"
)

type Group struct {
	ID          int    `json:"id"`
	URL         string `json:"web_url"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	FullPath    string `json:"full_Path"`
	Description string `json:"description"`
}

type Project struct {
	ID            int    `json:"id"`
	URL           string `json:"web_url"`
	Name          string `json:"name"`
	Path          string `json:"path"`
	Description   string `json:"description"`
	HttpURLToRepo string `json:"http_url_to_repo"`
}
