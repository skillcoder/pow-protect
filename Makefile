BLDCOMMIT = $(shell git rev-parse HEAD)
BLDBRANCH = $(shell git show-ref | grep $(BLDCOMMIT) | sed 's|.*/\(.*\)|\1|' | sort -u | grep -v HEAD | grep -v "^v" | head -n 1)
BLDVER    = $(shell git name-rev --tags --name-only $(BLDBRANCH) | cut -d "^" -f 1)
ECR=skillc0der
ECR_DIR=$(ECR)/
#IMG=$(APPNAME):$(BLDVER)

export DOCKER_BUILDKIT=1

all: test build

deps:
	@echo "+ $@"
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
	# go install github.com/golang/protobuf/protoc-gen-go@v1.5.2
	# sudo apt install protobuf-compiler
	export PATH="$PATH:$(go env GOPATH)/bin"

generate: gen.sdk
	@echo "+ $@"


gen.sdk:
	@echo "+ $@"
	protoc --proto_path=api \
		--go_out=./internal/api/sdk --go_opt=paths=source_relative \
        --go-grpc_out=./internal/api/sdk --go-grpc_opt=paths=source_relative \
        pow/pow.proto wisdom/wisdom.proto


test:
	@echo "+ $@"


build: build.server build.client

build.server:
	@echo "+ $@"
	GOGC=off CGO_ENABLED=0 go build -v -o ./bin/server ./cmd/server


build.client:
	@echo "+ $@"
	GOGC=off CGO_ENABLED=0 go build -v -o ./bin/client ./cmd/client


run: build.server run.server

run.server:
	@echo "+ $@"
	bin/server

run.client: build.client
	@echo "+ $@"
	bin/client

.PHONY: all test build build.server build.client \
	run run.server run.client