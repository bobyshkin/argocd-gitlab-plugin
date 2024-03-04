# Argo CD plugin

Argo CD plugin - a configurable tool that automates the process 
of creation ArgoCD applications based on Gitlab groups and projects

## Requirements

### Environment
```shell
GITLAB_TOKEN="read_api_access_token"
REPO_TYPE="helm" # Expected values: "git" or "helm"
```

### Volumes
```shell
config.ini #(currently built in image)
```

### Access to Kubernetes cluster
```shell
# In-Cluster
ServiceAccount: create 'applications.argoproj.io' at 'argocd' namespace

# Outside of Cluster
~/.kube/config
```
