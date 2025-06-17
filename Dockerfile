ARG GOLANG_VERSION=1.23.8
FROM golang:${GOLANG_VERSION}-bullseye as builder

# 设置工作目录
WORKDIR /app

# 设置 GOPROXY 以加速依赖下载
ARG GOPROXY=https://goproxy.cn

# 首先复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o youauth ./main.go

# 使用更小的基础镜像
FROM debian:bullseye-slim

# 安装必要的 CA 证书和 curl（用于健康检查）
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates curl && \
    rm -rf /var/lib/apt/lists/*

# 创建非 root 用户
RUN useradd -m -u 1000 appuser

# 设置工作目录
WORKDIR /app

# 复制必要的文件
COPY --from=builder /app/youauth /usr/local/bin/youauth
COPY --from=builder /app/static /app/static
COPY --from=builder /app/templates /app/templates
COPY --from=builder /app/dist /app/dist

# 设置权限
RUN chown -R appuser:appuser /app

# 切换到非 root 用户
USER appuser

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# 设置入口点
ENTRYPOINT ["/usr/local/bin/youauth", "run"]