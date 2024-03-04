# syntax=docker/dockerfile:1

# golang:1.17.6-alpine3.15 linux/amd64 07.01.22
FROM golang@sha256:f28579af8a31c28fc180fb2e26c415bf6211a21fb9f3ed5e81bcdbf062c52893 AS build-env

# All these steps will be cached
WORKDIR /app
COPY go.mod .
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN apk add --no-cache ca-certificates
RUN go mod download
# COPY the source code as the last step
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/argocd-gitlab-plugin

FROM scratch
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build-env /go/bin/argocd-gitlab-plugin /go/bin/argocd-gitlab-plugin
COPY config.ini .
ENTRYPOINT ["/go/bin/argocd-gitlab-plugin"]

