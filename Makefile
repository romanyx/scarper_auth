SHELL := /bin/sh

all: schema proto

proto:
	cd proto/ && make users

mocks:
	cd internal/verify && go generate
	cd internal/reg && go generate
	cd internal/change && go generate
	cd internal/reset && go generate

schema:
	cd internal/storage/postgres/schema && go generate

test:
	go test -v -race `go list ./... | grep -v kit | grep -v proto` 

cover: mocks
	go test `go list ./... | grep -v kit | grep -v proto` -coverprofile cover.out.tmp && \
		cat cover.out.tmp | grep -v "bindata.go" | grep -v "mock.go" | grep -v "main.go" > cover.out && \
		go tool cover -func cover.out && \
		rm cover.out.tmp && \
		rm cover.out

build:
	cd "$$GOPATH/src/github.com/romanyx/scraper_auth"
	docker build \
		-t scraper/auth:0.0.1 \
		-f docker/Dockerfile \
		--build-arg VCS_REF=`git rev-parse HEAD` \
		--build-arg BUILD_DATE=`date -u +”%Y-%m-%dT%H:%M:%SZ”` \
		.


