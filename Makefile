APP_NAME=golang-grpc-base-project
APP_VERSION=0.0.0
DOCKER_LOCAL=registry.gitlab.com/xdorro/registry
BUILD_DIR=./build
MAIN_DIR=./cmd

docker.build:
	docker build  -f $(BUILD_DIR)/Dockerfile -t $(DOCKER_LOCAL)/$(APP_NAME):$(APP_VERSION) .

docker.push:
	docker push $(DOCKER_LOCAL)/$(APP_NAME):$(APP_VERSION)

docker.dev: docker.build docker.push

grpc.install:
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

buf.gen:
	buf generate

buf.update:
	buf mod update

buf.install:
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install github.com/bufbuild/buf/cmd/protoc-gen-buf-breaking@latest
	go install github.com/bufbuild/buf/cmd/protoc-gen-buf-lint@latest

lint.install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest

lint.run:
	golangci-lint run ./...

ent.install:
	go install entgo.io/ent/cmd/ent@latest
	go install entgo.io/contrib/entproto/cmd/protoc-gen-entgrpc@latest

ent.init:
	go run entgo.io/ent/cmd/ent init --target pkg/ent/schema User

go.gen:
	go generate ./...

go.get:
	go get -u ./...

go.tidy:
	go mod tidy

go.test:
	go test ./...

go.lint: lint.run

go.install: grpc.install buf.install ent.install lint.install
