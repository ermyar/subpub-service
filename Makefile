LOCAL_BIN := $(CURDIR)/bin

PROTO_PATH := $(CURDIR)/api

PKG_PROTO_PATH := $(CURDIR)/pkg

PROTOC := PATH="$(PATH):$(LOCAL_BIN)" protoc

.bin-deps: export GOBIN := $(LOCAL_BIN)
.bin-deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

.protoc_generate:
	mkdir -p $(PKG_PROTO_PATH)
	$(PROTOC) --proto_path=$(CURDIR) --go_out=$(PKG_PROTO_PATH) --go_opt paths=source_relative \
	--go-grpc_out=$(PKG_PROTO_PATH) --go-grpc_opt paths=source_relative $(PROTO_PATH)/subpub/service.proto

.tidy:
	go mod tidy

generate: .bin-deps .protoc_generate .tidy
