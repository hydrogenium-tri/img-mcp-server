# 图像分析 MCP 服务器

让不支持多模态的 LLM 通过 MCP 协议实现对图片内容的分析。

## 功能特性

- MCP 协议兼容，暴露 `analyze_image` 工具
- 支持 `base64` 和上传缓存两种图片传入方式
- SHA256 哈希去重，同一图片多次上传不重复占用内存
- 缓存自动过期清理（10 分钟 TTL）
- 支持 JPEG / PNG / GIF / WebP / BMP 格式自动识别
- 兼容 OpenAI API 格式，可接入硅基流动等国产平台
- Bearer Token 认证，保护服务不被未授权访问
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
├── auth.go           # Bearer Token 认证中间件
├── test.sh           # 集成测试脚本
├── test.png          # 测试图片
├── .skill/SKILL.md   # Agent Skill 文件
├── go.mod / go.sum
└── config.json       # 配置文件（不提交到 Git）
```

## 快速开始

### 1. 创建配置文件

在项目根目录创建 `config.json`：

```json
{
  "provider": "siliconflow",
  "base_url": "https://api.siliconflow.cn/v1",
  "api_key": "sk-xxx",
  "model": "nex-agi/Nex-N2-Pro",
  "port": 8080,
  "auth_token": "your-secret-token"
}
```

| 字段 | 说明 |
|------|------|
| provider | 模型供应商名称 |
| base_url | 供应商 API 地址（兼容 OpenAI 格式） |
| api_key | API 密钥 ⚠️ 不要提交到 Git |
| model | 多模态模型名称（需支持视觉识别） |
| port | 服务器监听端口 |
| auth_token | 认证 Token（不设置则不启用认证） |

> **热更新**：修改 `config.json` 后重启服务即可生效，无需重新编译。

### 2. 编译运行

```bash
go build
./img-mcp-server
```

### 3. 验证服务

```bash
# 查询工具列表（未启用认证时）
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list"}'

# 启用认证后需携带 Token
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-token" \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list"}'
```

### 4. 运行集成测试

```bash
./test.sh
```

## Docker 部署（推荐）

```bash
# 构建镜像（国内已内置 goproxy.cn 代理）
docker build -t img-mcp-server .

# 运行容器（挂载配置文件 + 端口映射）
docker run -d --name img-mcp \
  --restart always \
  -p 9999:8080 \
  -v /path/to/config.json:/app/config.json:ro \
  img-mcp-server

# 查看日志
docker logs img-mcp

# 修改配置后重启
docker restart img-mcp
```

### 离线部署

```bash
# 导出镜像
docker save -o img-mcp-server.tar img-mcp-server

# 传输到目标服务器并加载
scp img-mcp-server.tar user@server:/path/
ssh user@server "docker load -i /path/img-mcp-server.tar"
```

---

## Agent Skill 配置

项目附带 `.skill/SKILL.md`，安装到全局 Skill 目录后，AI Agent 可自动完成"上传图片 → 获取缓存 ID → 分析 → 返回结果"的完整流程。

### 安装 Skill

1. 从 GitHub 下载 `.skill/SKILL.md`
2. 放入全局 Skill 目录 `~/.agents/skills/image-analysis/SKILL.md`
3. 如果服务启用了认证，修改 Skill 中的 curl 命令，添加 Token 头

### 配置服务器地址和 Token

Skill 默认连接 `http://localhost:8080`。如果部署到远程服务器且启用了认证，需要将以下信息告知 Agent：

- **服务器地址**：如 `https://www.megrez.space:9999`
- **认证 Token**：如 `mcp-link`

Agent 会在首次调用失败时询问这些信息，或用户可直接修改 Skill 文件。

---

## Zed / MCP 客户端配置

在 Zed 或其他 MCP 客户端的配置文件中添加：

```json
{
  "Vision_MCP": {
    "url": "https://your-server.com:9999/mcp",
    "headers": {
      "Authorization": "Bearer your-token"
    }
  }
}
```

配置文件路径：
- Zed：`~/.config/zed/mcp.json`
- 其他客户端：参考对应文档

---

## MCP 工具

### analyze_image

分析图片内容并返回文字描述。

| 参数 | 类型 | 必填 | 说明 |
|------|------|:----:|------|
| base64 | string | | 图片的 Base64 编码（与 cache_id 二选一） |
| cache_id | string | | 通过 /upload 接口上传后获取的缓存 ID（与 base64 二选一） |
| prompt | string | ✅ | 对图片的分析要求或提示词 |

---

## HTTP API 参考

所有接口均需认证（除非未配置 `auth_token`）。

### POST /upload

上传图片到服务器缓存，返回唯一 ID。

```bash
curl -X POST http://localhost:8080/upload \
  -H "Authorization: Bearer your-token" \
  -F "file=@图片.jpg"
```

返回：`{"cache_id": "abc123..."}`

- 使用 **SHA256 哈希** 作为缓存键，同一图片多次上传返回相同 ID
- 缓存 **10 分钟** 后自动过期

### POST /mcp

MCP 协议端点（Streamable HTTP，无状态模式）。

```bash
# 查询工具列表
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-token" \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list"}'

# 分析图片（cache_id 方式）
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-token" \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"analyze_image","arguments":{"cache_id":"abc123","prompt":"描述这张图片"}}}'

# 分析图片（base64 方式）
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-token" \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"analyze_image","arguments":{"base64":"iVBORw...","prompt":"描述"}}}'
```

---

## 安全建议

部署到公网时：
1. 务必设置 `auth_token`，防止未授权访问
2. 通过 Nginx/Caddy 反向代理添加 HTTPS
3. 容器绑定到 `127.0.0.1:8080`，仅允许反代访问
4. 使用防火墙限制访问来源 IP

## 许可证

MIT
