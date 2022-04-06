ifndef PORT
	PORT=12345
endif

start-server: build-server
	@./server -port=$(PORT)

run-client: build-client
	@./client $(OPTS)

build: build-server build-client

build-server:
	@go build cmd/server/server.go

build-client:
	@go build cmd/client/client.go
