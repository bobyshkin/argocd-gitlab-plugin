package gitlab

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"argocd-gitlab-plugin/model"

	"github.com/pkg/errors"
)

const (
	defaultBaseURL    = "https://gitlab.com/"
	apiVersionPath    = "api/v4"
	gitlabAccessToken = "GITLAB_TOKEN"
	gitlabAPI         = defaultBaseURL + apiVersionPath
	URI               = gitlabAPI
	StatusOK          = 200
)

func httpClient() *http.Client {
	client := &http.Client{Timeout: 10 * time.Second}
	return client
}

func sendRequest(client *http.Client, method string, request string) ([]byte, error) {
	req, err := http.NewRequest(method, request, nil)
	if err != nil {
		log.Fatalf("Error Occurred. %+v", err)
	}

	response, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request to API endpoint. %+v", err)
	}
	if response.StatusCode != StatusOK {
		return nil, errors.Errorf("failed to make gitlab api call: %s", response.Status)
	}

	// Close the connection to reuse it
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(response.Body)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Couldn't get response body. %+v", err)
	}
	return body, nil
}

func gitlabToken() string {
	value, ok := os.LookupEnv(gitlabAccessToken)
	if !ok {
		panic(fmt.Sprintf("Environment variable $%s is undefined", gitlabAccessToken))
	}
	return value
}

func GetGroupsProjects(groupId int) []model.Project {
	var (
		projects []model.Project // Declare []Project slice
		perPage  = 20            // PerPage projects result quantity
	)

	//Get Projects in every group listed

	var prj []model.Project

	// Iterate over pages
	for pageNumber := 1; ; pageNumber++ {

		//Request parameters
		c := httpClient()
		request := fmt.Sprintf(
			"%s/groups/%d/projects?access_token=%s&per_page=%d&page=%d&include_subgroups=true&simple=true",
			URI,
			groupId,
			gitlabToken(),
			perPage,
			pageNumber,
		)
		response, err := sendRequest(c, http.MethodGet, request)
		if err != nil {
			os.Exit(1)
		}

		//Decode JSON Response to []Project struct
		JSON := json.Unmarshal(response, &prj)
		if JSON != nil {
			os.Exit(1)
		}

		if len(prj) == 0 {
			break
		}

		projects = append(projects, prj...)
	}

	return projects
}

func GetProject(projectID int) model.Project {
	var project model.Project

	//Request parameters
	c := httpClient()
	request := fmt.Sprintf(
		"%s/projects/%d/?access_token=%s&simple=true",
		URI,
		projectID,
		gitlabToken(),
	)
	response, err := sendRequest(c, http.MethodGet, request)
	if err != nil {
		os.Exit(1)
	}

	//Decode JSON Response to []Project struct
	JSON := json.Unmarshal(response, &project)
	if JSON != nil {
		os.Exit(1)
	}

	return project
}
