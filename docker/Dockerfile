FROM golang:1.16-alpine
COPY . /app
WORKDIR /app
ARG CGO_ENABLED=0
RUN go build cmd/greeter/main.go
FROM scratch
COPY --from=0 /app/main /greeter
WORKDIR /
ENTRYPOINT ["/greeter"]