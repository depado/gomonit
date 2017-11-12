FROM golang:1.9 AS build
RUN mkdir -p $GOPATH/src/github.com/Depado/gomonit
ADD . $GOPATH/src/github.com/Depado/gomonit
WORKDIR $GOPATH/src/github.com/Depado/gomonit
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o gomonit
RUN cp gomonit /

FROM golang:1.9
COPY --from=build /gomonit /usr/bin/
ENTRYPOINT ["/usr/bin/gomonit"]
EXPOSE 8080