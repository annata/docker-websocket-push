FROM golang:1.19-alpine as builder
WORKDIR /websocket_push
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -buildvcs=false

FROM scratch
COPY --from=builder /websocket_push/websocket_push /
CMD ["/websocket_push"]