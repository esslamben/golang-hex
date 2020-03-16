FROM golang:1.13.3-alpine as builder

# Needed for fetching go deps
RUN apk update && \
    apk add git && \
    apk add make

ENV GO111MODULE=auto
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH

# Copy everything in current working directory over to the equivalent go path
COPY . /golang-hex

# Change current working directory to be the devices service path on the container
WORKDIR  /golang-hex/cmd/golang-hex

# Run a go get to install any required imports run a go build
RUN go get ./... && \
    GOOS=linux GOARCH=amd64 go build -o golang-hex.linux.x86_64 .

FROM alpine:3.5

RUN apk add --no-cache ca-certificates

COPY --from=builder /golang-hex/cmd/golang-hex/golang-hex.linux.x86_64 /usr/bin/golang-hex

ENTRYPOINT ["/usr/bin/golang-hex"]
EXPOSE 5000