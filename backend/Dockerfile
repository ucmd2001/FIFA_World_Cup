# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# 安裝 gcc 與 musl-dev 以編譯 CGO (SQLite 需要)
RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 編譯 CGO 開啟，並產出二進位檔
RUN CGO_ENABLED=1 GOOS=linux go build -o main ./cmd/server/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# 安裝 SQLite 與基本的庫
RUN apk add --no-cache ca-certificates sqlite-libs

# 從 builder 階段複製二進位檔與環境變數範例
COPY --from=builder /app/main .
COPY --from=builder /app/.env.example ./.env

EXPOSE 8080

CMD ["./main"]
