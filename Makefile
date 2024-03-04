BINARY=argocd-gitlab-plugin
CGO_ENABLED=0
GOOS=linux
GOARCH=amd64

default: build

build:
	go build -a -installsuffix cgo -o /go/bin/${BINARY}

install: build

e2e: install
	./${BINARY}
