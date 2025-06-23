# blog-image

随机图片API服务，为 [My Blog](https://zerowolf.cn) 提供随机图片支持。

## 功能特性

- 支持PC端和移动端不同尺寸的随机图片
- 自动设备检测，根据User-Agent自动返回适合的图片
- 安全的路径验证，防止路径遍历攻击
- 高性能，支持多种图片格式（JPEG, PNG, WebP, GIF）
- Docker友好，支持环境变量配置
- 详细的日志记录和错误处理

## API端点

| 端点 | 描述 | 示例 |
|------|------|------|
| `GET /` | 主页面，显示可用API | `http://localhost:8000/` |
| `GET /auto` | 自动检测设备类型并返回对应图片 | `http://localhost:8000/auto` |
| `GET /pc` | 返回PC端随机图片 | `http://localhost:8000/pc` |
| `GET /mobile` | 返回移动端随机图片 | `http://localhost:8000/mobile` |
| `GET /ciallo` | 彩蛋页面 | `http://localhost:8000/ciallo` |

## 快速开始

### 1. 准备图片

```bash
# 创建图片目录
mkdir -p image/pc image/mobile

# 将PC端图片放入 image/pc/ 目录
# 将移动端图片放入 image/mobile/ 目录
```

### 2. 运行服务

```bash
# 使用默认端口8000
go run main.go

# 或指定端口
go run main.go -port=3000

# 或使用环境变量
export port=3000
go run main.go
```

### 3. 访问服务

打开浏览器访问 `http://localhost:8000`

## 支持的图片格式

- JPEG (`.jpg`, `.jpeg`)
- PNG (`.png`)
- WebP (`.webp`)
- GIF (`.gif`)

## 环境变量

| 变量名 | 描述 | 默认值 |
|--------|------|--------|
| `port` | 服务端口 | `8000` |

## Docker 部署

```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o blog-image main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/blog-image .
COPY --from=builder /app/static ./static
COPY --from=builder /app/image ./image
EXPOSE 8000
CMD ["./blog-image"]
```

## 开发说明

### 代码改进点

- 使用现代Go API（`os.ReadDir`替代已废弃的`ioutil.ReadDir`）
- 消除重复代码，统一图片处理逻辑
- 添加路径遍历攻击防护
- 改进错误处理和日志记录
- 添加HTTP超时配置
- 支持更多图片格式
- 自动创建必要目录
- 添加缓存控制头

### 文件结构

```
blog-image/
├── main.go           # 主程序
├── go.mod           # Go模块文件
├── README.md        # 项目文档
├── static/          # 静态文件
│   └── index.html   # 主页面
└── image/           # 图片目录
    ├── pc/          # PC端图片
    └── mobile/      # 移动端图片
```

## 许可证

查看 [LICENSE](LICENSE) 文件了解详情。
