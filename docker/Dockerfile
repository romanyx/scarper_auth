FROM golang:latest as build_env

COPY . $GOPATH/src/github.com/romanyx/scraper_auth
WORKDIR $GOPATH/src/github.com/romanyx/scraper_auth/cmd/auth

ENV GO111MODULE=on
RUN go get -v
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -a -tags netgo -o /go/bin/auth

FROM alpine
ARG BUILD_DATE
ARG VCS_REF
COPY --from=build_env /go/bin/auth /bin/auth

ENTRYPOINT ["/bin/auth"]
