
GOPATH := $(shell pwd)
export GOPATH
PATH := ${PATH}:$(shell pwd)/bin
export PATH

build:
	go install github.com/bmeg/arachne

PROTO_INC= -I ./ -I ../googleapis

proto:
	cd src/github.com/bmeg/arachne/aql && protoc \
	$(PROTO_INC) \
	--go_out=\
	Mgoogle/protobuf/struct.proto=github.com/golang/protobuf/ptypes/struct,\
	plugins=grpc:./ \
	--grpc-gateway_out=logtostderr=true:. \
	aql.proto

download:
	go get -d github.com/bmeg/arachne
	go get github.com/golang/protobuf/protoc-gen-go
	go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway/

test:
	go test -v github.com/bmeg/arachne/

tidy:
	@find ./src/github.com/bmeg/arachne -type f | grep -v ".pb." | grep -E '.*\.go$$' | xargs gofmt -w
