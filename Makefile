include server/.env

.PHONY: all build-server build-client up-server up-client clean

all: build-server build-client

build-server:
	go build -o bin/server server/cmd/main.go

build-client:
	go build -o bin/client client/cmd/main.go

up-client: build-client
	PORT=${CLIENT_PORT} ./bin/client

up-server: build-server
	@PORT=${SERVER_PORT} ADDR=${ADDR} PRIVATE=${PRIVATE} PUBLIC=${PUBLIC} FILE_STORAGE=${FILE_STORAGE} ./bin/server

clean:
	rm -f bin/*

deps:
	@echo "Install dependenciy packages ..."
	go get go get github.com/go-kit/kit/log