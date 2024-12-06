FROM golang:1.21-alpine as builder
WORKDIR /websocket_push
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -buildvcs=false -ldflags="-w -s"

FROM buildpack-deps:curl
COPY --from=builder /websocket_push/websocket_push /
CMD ["/websocket_push"]