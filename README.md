# 图像分析 MCP 服务器

让不支持多模态的 LLM 通过 MCP 协议实现对图片内容的分析。

## 功能特性

- MCP 协议兼容，暴露 `analyze_image` 工具
- 支持 `base64` 和上传缓存两种图片传入方式
- SHA256 哈希去重，同一图片多次上传不重复占用内存
- 缓存自动过期清理（10 分钟 TTL）
- 支持 JPEG / PNG / GIF / WebP / BMP 格式自动识别
- 兼容 OpenAI API 格式，可接入硅基流动等国产平台
- Docker 容器化部署，支持离线导出
- 附带集成测试脚本和 Agent Skill 文件

## 项目结构

```
.
├── Dockerfile        # Docker 容器化部署
├── main.go           # 入口
├── config.go         # 配置加载
├── server.go         # MCP 服务器启动 + 路由注册
├── tools.go          # MCP 工具定义 + 处理函数
├── model.go          # 图片格式检测 + 视觉模型调用
├── upload.go         # HTTP 上传接口
├── cache.go          # 图片缓存（SHA256 去重 + TTL 过期）
├── test.sh           # 集成测试脚本
├── test.png          # 测试图片
├── .skill/SKILL.md   # Agent Skill 文件
├── go.mod / go.sum
└── config.json       # 配置文件（不提交到 Git）
```

## 环境要求

- Go 1.21+ 或 Docker
- 一个支持多模态模型的 API Key（如硅基流动）

## 配置

在项目根目录创建 `config.json`：

```json
{
  "provider": "siliconflow",
  "base_url": "https://api.siliconflow.cn/v1",
  "api_key": "sk-xxx",
  "model": "nex-agi/Nex-N2-Pro",
  "port": 8080
}
```

| 字段 | 说明 |
|------|------|
| provider | 模型供应商名称 |
| base_url | 供应商 API 地址（兼容 OpenAI 格式） |
| api_key | API 密钥 ⚠️ 不要提交到 Git |
| model | 多模态模型名称 |
| port | 服务器监听端口 |

## Docker 部署（推荐）

```bash
# 构建镜像（国内已内置 goproxy.cn 代理）
docker build -t img-mcp-server .

# 运行容器
docker run -d --name img-mcp \
  --restart always \
  -p 8080:8080 \
  -v /path/to/config.json:/app/config.json:ro \
  img-mcp-server

# 查看日志
docker logs img-mcp
```

### 离线部署

```bash
# 导出镜像为 tar 文件
docker save -o img-mcp-server.tar img-mcp-server

# 传输到目标服务器
scp img-mcp-server.tar user@server:/path/

# 在服务器上加载
docker load -i img-mcp-server.tar
```

### 端口映射

默认容器内监听 8080，通过 `-p` 映射到任意宿主机端口：

```bash
# 映射到 9999 端口
docker run -d --name img-mcp \
  --restart always \
  -p 9999:8080 \
  -v /path/to/config.json:/app/config.json:ro \
  img-mcp-server
```

## 编译与运行

```bash
# 编译
go build

# 启动服务器（默认监听 :8080）
./img-mcp-server

# 运行集成测试
./test.sh
```

## MCP 工具

### analyze_image

分析图片内容并返回文字描述。

| 参数 | 类型 | 必填 | 说明 |
|------|------|:----:|------|
| base64 | string | | 图片的 Base64 编码（与 cache_id 二选一） |
| cache_id | string | | 通过 /upload 接口上传后获取的缓存 ID（与 base64 二选一） |
| prompt | string | ✅ | 对图片的分析要求或提示词 |

---

## HTTP API

### POST /upload

上传图片到服务器缓存，返回唯一 ID。

```bash
curl -X POST http://localhost:8080/upload -F "file=@图片.jpg"
```

返回示例：

```json
{"cache_id": "abc123def456..."}
```

特点：
- 使用 **SHA256 哈希** 作为缓存键，同一图片多次上传返回相同 `cache_id`
- 图片缓存 **10 分钟** 后自动过期

### POST /mcp

MCP 协议端点（Streamable HTTP，无状态模式）。

```bash
# 查询工具列表
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list"}'

# 分析图片（用 cache_id）
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"analyze_image","arguments":{"cache_id":"abc123","prompt":"描述这张图片"}}}'
```

---

## Agent Skill

项目附带 `.skill/SKILL.md`，安装后 AI Agent 可自动完成"上传图片 → 分析 → 返回结果"的完整流程。

**安装方式：** 将 `.skill/SKILL.md` 复制到全局 Skill 目录或直接引用 GitHub 仓库。

---

## 安全建议

部署到公网时建议：
1. 通过 Nginx/Caddy 反向代理添加 HTTPS
2. 将容器绑定到 `127.0.0.1:8080`，只允许反代访问
3. 使用防火墙限制 /upload 接口的访问来源

## 许可证

MIT
