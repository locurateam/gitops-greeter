FROM golang:1.16-alpine
COPY . /go
ARG CGO_ENABLED=0
RUN go build main.go
FROM scratch
COPY --from=0 /go/main /greeter
COPY --from=0 /go/assets /assets
WORKDIR /
ENTRYPOINT ["/greeter"]