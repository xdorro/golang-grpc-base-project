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

wire.install:
	go install github.com/google/wire/cmd/wire@latest

grpc.install:
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

lint.install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest

lint.run:
	golangci-lint run --fast ./...

go.install: grpc.install lint.install wire.install

go.gen: wire.gen

go.lint: lint.run

cert.gen:
	mkcert -install
	mkcert -key-file ./config/key.pem -cert-file ./config/cert.pem localhost 127.0.0.1 ::1

go.get:
	go get -u ./...

go.tidy:
	go mod tidy -compat=1.17

go.test:
	go test ./...

