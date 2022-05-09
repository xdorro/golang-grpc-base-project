# golang-grpc-base-project

## Getting started

Setup GO PRIVATE

Windows:
```
go env -w GOPRIVATE=github.com/*
```

Linux:
```
export GOPRIVATE=github.com/*
```

Get all modules
```
env GIT_TERMINAL_PROMPT=1 go mod tidy
```