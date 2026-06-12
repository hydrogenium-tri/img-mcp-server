# 第一阶段：编译
FROM golang:1.26.4-alpine AS builder
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o img-mcp-server

# 第二阶段：运行（用最小镜像）
FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/img-mcp-server /app/img-mcp-server
EXPOSE 8080
ENTRYPOINT ["/app/img-mcp-server"]