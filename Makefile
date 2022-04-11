ifndef PORT
	PORT=12345
endif

start-server: build-server
	@./bin/server -port=$(PORT)

run-client: build-client
	@./bin/client $(OPTS)

build: build-server build-client

build-server:
	@mkdir -p ./bin
	@go build cmd/server/server.go
	@mv server ./bin/

build-client:
	@mkdir -p ./bin
	@go build cmd/client/client.go
	@mv client ./bin/
