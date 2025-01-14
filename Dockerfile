# Build Stage
FROM golang:1.22-alpine as builder

# Set architecture environment variable
ARG ARCH=amd64

WORKDIR /app

COPY . .

RUN go mod tidy

# Build for dynamic architecture (default: amd64)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$ARCH go build -ldflags="-w -s" -o exporter .

# Final image stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/exporter .

EXPOSE 8090

ENTRYPOINT ["./exporter"]
