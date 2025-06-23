# 构建和运行指南

## 本地开发

```bash
# 运行服务
go run main.go

# 或者构建后运行
go build -o blog-image main.go
./blog-image
```

## Docker 构建和运行

```bash
# 构建镜像
docker build -t blog-image .

# 运行容器
docker run -p 8000:8000 \
  -v $(pwd)/image:/root/image \
  blog-image

# 或者使用docker-compose
docker-compose up -d
```

## 性能优化建议

1. **静态文件服务**: 在生产环境中，建议使用Nginx等反向代理服务器来处理静态文件
2. **缓存策略**: 可以添加Redis缓存来缓存图片文件列表
3. **CDN**: 将图片托管到CDN服务，提高访问速度
4. **负载均衡**: 在高并发场景下，可以部署多个实例进行负载均衡
