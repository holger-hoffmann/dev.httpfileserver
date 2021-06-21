FROM golang:latest as builder

# Build the application
COPY . /go/src/github.com/holger-hoffmann/dev.httpfileserver
WORKDIR /go/src/github.com/holger-hoffmann/dev.httpfileserver
RUN CGO_ENABLED=0 go build -a -o main ./server/http/

# Create the production container with only the minimum necessary
FROM scratch
MAINTAINER Holger Hoffmann <holger.h.hoffmann@gmail.com>

COPY --from=builder /go/src/github.com/holger-hoffmann/dev.httpfileserver/main /

ENTRYPOINT ["/main"]
