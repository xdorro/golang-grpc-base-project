APP_NAME=golang-grpc-base-project
APP_VERSION=latest
DOCKER_REGISTRY=registry.gitlab.com/xdorro/registry
MAIN_DIR=./cmd

docker.build:
	docker build -t $(DOCKER_REGISTRY)/$(APP_NAME):$(APP_VERSION) .

docker.push:
	docker push $(DOCKER_REGISTRY)/$(APP_NAME):$(APP_VERSION)

docker.dev: docker.build docker.push

wire.gen:
	wire ./...

lint.run:
	golangci-lint run --fast ./...

go.install:
#	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

	go install github.com/bufbuild/buf/cmd/buf@latest
	go install github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go@latest

	go install github.com/google/wire/cmd/wire@latest

	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest


go.gen: wire.gen

go.lint: lint.run

cert.gen:
	mkcert -install
	mkcert -key-file ./config/cert/key.pem -cert-file ./config/cert/cert.pem localhost 127.0.0.1 ::1

go.get:
	go get -u ./...

go.tidy:
	go mod tidy -compat=1.18

go.test:
	go test ./...

