FROM golang:1.16-alpine as build

WORKDIR /build

COPY . .

RUN  go mod tidy
RUN  go mod vendor
RUN  GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=on \
    go build \
    -o /app ./cmd/server

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app /
COPY config.yaml /

ENTRYPOINT ["/app"]
