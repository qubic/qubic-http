PB = $(wildcard *.proto)
GO = $(PB:.proto=.pb.go)
PWD = $(pwd)

all: $(GO)

%.pb.go: %.proto
		protoc -I=. --go-grpc_out=paths=source_relative:. \
		--grpc-gateway_out . \
    --grpc-gateway_opt logtostderr=true \
    --grpc-gateway_opt paths=source_relative \
    --grpc-gateway_opt generate_unbound_methods=true \
    --proto_path=/opt/homebrew/Cellar/protobuf/25.3_1/include \
		--go_out=paths=source_relative:. *.proto

openapi: *.proto
		protoc --openapiv2_out=logtostderr=true:. \
			   --proto_path=.:/opt/homebrew/Cellar/protobuf/25.3_1/include \
			   *.proto

clean:
		rm -f *.pb.go
		rm -f *.go