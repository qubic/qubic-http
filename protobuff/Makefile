# Check for macOS arm64 - dynamically detect protobuf version
UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)
PROTOBUF_VERSION := $(shell ls /opt/homebrew/Cellar/protobuf/ 2>/dev/null | head -1)

ifeq ($(UNAME_S),Darwin)
    ifeq ($(UNAME_M),arm64)
        ifneq ($(PROTOBUF_VERSION),)
            OPT_ARGS := --proto_path=/opt/homebrew/Cellar/protobuf/$(PROTOBUF_VERSION)/include
        else
            OPT_ARGS :=
        endif
    else
        OPT_ARGS :=
    endif
else
    OPT_ARGS :=
endif

PB = $(wildcard *.proto)
GO = $(PB:.proto=.pb.go)
PWD = $(pwd)

all: $(GO) openapi-v3

%.pb.go: %.proto
		protoc -I=. --go-grpc_out=paths=source_relative:. \
		--grpc-gateway_out . \
    --grpc-gateway_opt logtostderr=true \
    --grpc-gateway_opt paths=source_relative \
    --grpc-gateway_opt generate_unbound_methods=true $(OPT_ARGS)  \
		--go_out=paths=source_relative:. *.proto

openapi-v3:
		protoc -I=. --openapi_out=output_mode=source_relative:. *.proto

clean:
		rm -f *.pb.go
		rm -f *.go