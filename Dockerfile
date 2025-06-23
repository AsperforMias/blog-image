FROM golang:1.23-alpine AS builder

WORKDIR /app

# 复制go mod文件
COPY go.mod go.sum* ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o blog-image main.go

# 运行阶段
FROM alpine:latest

# 安装ca-certificates用于HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/blog-image .
COPY --from=builder /app/static ./static

# 创建图片目录
RUN mkdir -p image/pc image/mobile

# 暴露端口
EXPOSE 8000

# 设置环境变量
ENV port=8000

# 运行应用
CMD ["./blog-image"]
