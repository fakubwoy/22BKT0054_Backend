FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .

RUN apk add --no-cache git \
    && go mod download \
    && CGO_ENABLED=0 GOOS=linux go build -o fileshare cmd/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/fileshare .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/uploads ./uploads

RUN mkdir -p /app/uploads \
    && chmod -R 755 /app/uploads

EXPOSE 8080
CMD ["./fileshare"]