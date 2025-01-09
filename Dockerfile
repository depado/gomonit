# Build Step
FROM golang:1.23.4-alpine@sha256:04ec5618ca64098b8325e064aa1de2d3efbbd022a3ac5554d49d5ece99d41ad5 as builder

# Dependencies
RUN apk update && apk add --no-cache make git

# Source
WORKDIR $GOPATH/src/github.com/depado/gomonit
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify
COPY . .

# Build
RUN make tmp

# Final Step
FROM gcr.io/distroless/static@sha256:5c7e2b465ac6a2a4e5f4f7f722ce43b147dabe87cb21ac6c4007ae5178a1fa58
COPY --from=builder /tmp/gomonit /go/bin/gomonit
COPY assets assets/
COPY templates templates/
ENTRYPOINT ["/go/bin/gomonit"]
