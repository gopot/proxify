ARG GOLANG_IMAGE_TAG=1.11.1-alpine3.8
ARG RUNNER_ALPINE_IMAGE_TAG=3.8

FROM golang:${GOLANG_IMAGE_TAG} as builder

RUN apk update && apk upgrade

RUN apk add git

RUN go version

ADD . $GOPATH/src/project

WORKDIR $GOPATH/src/project

RUN ls

RUN mkdir -p /built

RUN go get ./...

## Test project
#RUN go test -v ./...

# Build service
RUN go build -o /built/service .

FROM alpine:${RUNNER_ALPINE_IMAGE_TAG} AS runner

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /built/service ./

ENTRYPOINT [ "./service" ]
