LOCAL_BIN := $(CURDIR)/bin

PROTO_PATH := $(CURDIR)/api

PKG_PROTO_PATH := $(CURDIR)/pkg

PROTOC := PATH="$(PATH):$(LOCAL_BIN)" protoc

.bin_deps: export GOBIN := $(LOCAL_BIN)
.bin_deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

.protoc_generate:
	mkdir -p $(PKG_PROTO_PATH)
	$(PROTOC) --proto_path=$(CURDIR) --go_out=$(PKG_PROTO_PATH) --go_opt paths=source_relative \
	--go-grpc_out=$(PKG_PROTO_PATH) --go-grpc_opt paths=source_relative $(PROTO_PATH)/subpub/service.proto

.tidy:
	go mod tidy

# generating .pb files
generate: .bin_deps .protoc_generate .tidy

# building project
build:
	mkdir -p $(LOCAL_BIN)
	go build -o $(LOCAL_BIN) ./cmd/client/...
	go build -o $(LOCAL_BIN) ./cmd/server/...


.PHONY: .bin_deps .protoc_generate generate .tidy