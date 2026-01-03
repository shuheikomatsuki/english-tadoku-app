# --- ビルドステージ ---
FROM golang:1.24.3-alpine AS builder
RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/api

# --- 実行ステージ ---
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/server .

# Renderが環境変数を注入できるように、ENTRYPOINTを使用
ENTRYPOINT ["./server"]

EXPOSE 8080