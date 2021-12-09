# Go targets are intentionally phony; rely on go build cache.

.PHONY: all
all: server client example

.PHONY: server
server:
	go build -o server ./server

.PHONY: client
client:
	go build -o client ./client

.PHONY: example
example:
	go build -o worker/example ./worker/example

.PHONY: test
test:
	go test -a -race ./...

.PHONY: lint
lint:
	golangci-lint run

protos: protos/jobworker.pb.go protos/jobworker_grpc.pb.go

protos/jobworker.pb.go protos/jobworker_grpc.pb.go: protos/jobworker.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative protos/jobworker.proto

.PHONY: certs
certs:
	cd certs; ./create.sh
