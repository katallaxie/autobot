DIR ?= $(shell pwd)
PROTOBUF_API ?= $(DIR)/pkg/proto/v1

.PHONY: generate
generate: proto
generate:
	@go generate ./...

.PHONY: fmt
fmt: ## Run go fmt against code.
	@go run mvdan.cc/gofumpt -w .

.PHONY: vet
vet: ## Run go vet against code.
	@go vet ./...

.PHONY: build
build:
	@goreleaser build --rm-dist --snapshot

.PHONY: test
test:
	@go test -race ./... -count=1 -cover -coverprofile cover.out

.PHONY: run-scylla
run-scylla:
	@docker pull scylladb/scylla
	@docker run --name scylla -p 9042:9042 --privileged --memory 1G --rm -d scylladb/scylla
	@until docker exec scylla cqlsh -e "DESCRIBE SCHEMA"; do sleep 2; done

download:
	@go mod download

install-tools: download
	@cat internal/tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

.PHONY: stop-scylla
stop-scylla:
	@docker stop scylla

.PHONY: install
install:
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

.PHONY: proto
proto: ## generates protobuf
	@$(eval DEPS=`find $(PROTOBUF_API) -type f -name '*.proto'`)
	@docker run --rm -u $(shell id -u) -v$(PROTOBUF_API):$(PROTOBUF_API) \
		-w$(PROTOBUF_API) ghcr.io/katallaxie/protobuf-docker:development \
		--proto_path=$(PROTOBUF_API) \
		--go_out=$(PROTOBUF_API) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(PROTOBUF_API) \
		--go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(PROTOBUF_API) \
		--grpc-gateway_opt=logtostderr=true \
		--grpc-gateway_opt=paths=source_relative \
		--grpc-gateway_opt=generate_unbound_methods=true \
		$(DEPS)
