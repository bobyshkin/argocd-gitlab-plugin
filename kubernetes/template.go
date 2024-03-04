package kubernetes

import (
	"context"
	"fmt"
	"os"

	"argocd-gitlab-plugin/model"

	argocd "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	argocdclient "github.com/argoproj/argo-cd/v2/pkg/client/clientset/versioned"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func PrepareApplication(project model.Project, appType model.ApplicationType) (*argocd.Application, error) {
	var application *argocd.Application
	switch appType {
	case model.GitApplicationType:
		application = &argocd.Application{
			ObjectMeta: metav1.ObjectMeta{
				Name:       project.Path,
				Finalizers: []string{"resources-finalizer.argocd.argoproj.io"},
			},
			Spec: argocd.ApplicationSpec{
				Destination: argocd.ApplicationDestination{
					Name: "in-cluster",
					//Server:    "https://kubernetes.default.svc",
					Namespace: "microservices", //todo: dynamically change namespace(?)
				},
				Project: "microservices", //todo: dynamically change project(?)
				Source: argocd.ApplicationSource{
					Path:           "chart",
					RepoURL:        project.HttpURLToRepo,
					TargetRevision: "develop",
				},
			},
		}
	case model.HelmApplicationType:
		application = &argocd.Application{ // todo helm type
			ObjectMeta: metav1.ObjectMeta{
				Name:       project.Path,
				Finalizers: []string{"resources-finalizer.argocd.argoproj.io"},
			},
			Spec: argocd.ApplicationSpec{
				Destination: argocd.ApplicationDestination{
					Name: "in-cluster",
					//Server:    "https://kubernetes.default.svc",
					Namespace: "microservices", //todo: dynamically change namespace(?)
				},
				Project: "microservices", //todo: dynamically change project(?)
				Source: argocd.ApplicationSource{
					Chart: project.Path,
					RepoURL: fmt.Sprintf(
						"https://gitlab.com/api/v4/projects/%d/packages/helm/stable",
						project.ID,
					),
					TargetRevision: "x",
				},
			},
		}
	default:
		return &argocd.Application{}, errors.Errorf("Unsupported repo type %s.\n", appType)
	}
	return application, nil
}

func CreateApplication(kubeconfig *string, app *argocd.Application) {

	var config *rest.Config
	var err error

	host, port := os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT")
	// If env variable KUBERNETES_SERVICE_HOST or KUBERNETES_SERVICE_PORT does not set then
	if len(host) == 0 || len(port) == 0 {
		fmt.Println("Using .kube/config")
		// ...uses .kube/config
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	} else {
		// ...else uses ServiceAccount
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		panic(err.Error())
	}

	// Create the clientset
	clientset, err := argocdclient.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	applicationClient := clientset.ArgoprojV1alpha1().Applications("argocd")

	// Create Application
	result, err := applicationClient.Create(context.TODO(), app, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("WARNING: %s (skipping).\n", err.Error())
	} else {
		fmt.Printf("Created application %q.\n", result.GetObjectMeta().GetName())
	}
}
