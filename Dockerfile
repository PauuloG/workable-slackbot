## BUILDER PART
FROM golang:alpine AS builder

COPY . $GOPATH/src/mypackage/myapp/
WORKDIR $GOPATH/src/mypackage/myapp/
RUN ls -la $GOPATH/src/mypackage/myapp/

RUN apk update && apk add --no-cache git

RUN adduser -D -g '' workable

RUN go get -d -v
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w" -o /go/bin/workable-slackbot

## RUNNER PART
FROM scratch

# We copy the user entry from the builder
COPY --from=builder /etc/passwd /etc/passwd

# We also need the ca-certificates for x509
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# We copy the binary from the builder
COPY --from=builder /go/bin/workable-slackbot /go/bin/workable-slackbot

USER workable

# Run the binary.
ENTRYPOINT ["/go/bin/workable-slackbot"]
