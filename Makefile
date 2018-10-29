include server/.env
include client/.env

.PHONY: all build-server build-client up-server up-client clean

all: build-server build-client
	@echo ${CLIENT_PORT}
	@echo ${SERVER_PORT}

build-server:
	go build -o bin/server server/cmd/main.go

build-client:
	go build -o bin/client client/cmd/main.go

up-client: build-client
	./bin/client

up-server: build-server
	./bin/server

clean:
	rm -f bin/*

deps:
	@echo "Install dependenciy packages ..."
	@go get github.com/crsmithdev/goenv