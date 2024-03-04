package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"argocd-gitlab-plugin/gitlab"
	"argocd-gitlab-plugin/kubernetes"
	"argocd-gitlab-plugin/model"

	"gopkg.in/ini.v1"
	"k8s.io/client-go/util/homedir"
)

const (
	config         = "config.ini"
	appTypeVarName = "REPO_TYPE"
)

func getAppType() model.ApplicationType {
	value, ok := os.LookupEnv(appTypeVarName)
	if !ok {
		panic(fmt.Sprintf("Environment variable $%s is undefined.\n", appTypeVarName))
	}
	return model.ApplicationType(value)
}

func main() {
	// Get config
	cfg, err := ini.Load(config)
	if err != nil {
		fmt.Printf("Fail to read config file: %v\n", err)
		os.Exit(1)
	}

	//Get kubeconfig from $HOME/.kube/config if it exists (for local use only)
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// Read of values from config.ini
	for _, sect := range cfg.Sections() {
		if sect.Name() == ini.DefaultSection {
			continue
		}
		name := sect.Name()
		id := sect.Key("id").String()
		fmt.Printf("Section: %s\nID: %s\n", name, id)

		objectID, iErr := strconv.Atoi(id)
		if iErr != nil {
			fmt.Printf("Section: %s\n ID: %s\nFailed to parse config.\n", name, id)
		}

		var projects []model.Project

		objectType := sect.Key("type").String()
		switch model.ObjectType(objectType) {
		case model.GroupObjectType:
			projects = gitlab.GetGroupsProjects(objectID)
		case model.ProjectObjectType:
			project := gitlab.GetProject(objectID)
			projects = []model.Project{project}
		default:
			fmt.Printf("Unsupported type: %s.\n", objectType)
		}

		//Create Kubernetes resources
		for i := range projects {
			fmt.Printf("Creating application: %s\n", projects[i].Name)
			app, err := kubernetes.PrepareApplication(projects[i], getAppType())
			if err != nil {
				fmt.Printf("Failed to prepare k8s application: %s.\n", err.Error())
				os.Exit(1)
			}
			kubernetes.CreateApplication(kubeconfig, app)
		}
	}
	fmt.Println("Finished sync job without errors... Exiting(0)")
}
