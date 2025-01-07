FROM golang:1.22-alpine as builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o exporter .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/exporter .

EXPOSE 8090

CMD ["./exporter"]
