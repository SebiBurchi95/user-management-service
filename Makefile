PKGS := $(shell go list ./...)

gomod: 
	go mod tidy

build: 
	go build -o user-management-service

run:
	./user-management-service

test: gomod
	mkdir -p out
	go test -v ./... -race -short -coverprofile out/cover.out
	go tool cover -html=out/cover.out -o out/cover.html
	go tool cover -func=out/cover.out

all: gomod build run