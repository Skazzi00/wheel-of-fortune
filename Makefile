all: server ;

client:
	go build -mod=vendor -v -o ./bin/client ./cmd/client

server:
	go build -mod=vendor -v -o ./bin/server ./cmd/server

vendor:
	go mod vendor

clean:
	rm -fv ./bin/read-server
	rm -fv ./bin/client

.PHONY: all bin/* test vendor clean