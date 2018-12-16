SHELL := /bin/sh

all: schema proto

proto:
	cd proto/ && make users

schema:
	cd internal/storage/postgres/schema && go generate

test:
	go test -v -race `go list ./... | grep -v kit | grep -v proto` 

cover:
	go test `go list ./... | grep -v kit | grep -v proto` -coverprofile cover.out.tmp && \
		cat cover.out.tmp | grep -v "bindata.go" | grep -v "mock.go" | grep -v "main.go" > cover.out && \
		go tool cover -func cover.out && \
		rm cover.out.tmp && \
		rm cover.out
