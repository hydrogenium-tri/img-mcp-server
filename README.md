# 图像分析 MCP 服务器

## 简介
这是一个使用go语言实现的简单的图像MCP服务器，可以让非多模态模型通过MCP协议实现对图像的分析。

## 功能特性
- 暴露了 MCP 工具的图像分析能力
- 传入需要分析的图片，使用配置好的多模态模型对图片进行描述并返回描述词

## 环境要求
- Go 1.21+
- 一个支持多模态模型的API Key

## 配置
配置文件应放置在项目根目录中，命名为config.json
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
| base_url | 供应商 API 地址 |
| api_key | API 密钥 ⚠️ 不要提交到 Git |
| model | 多模态模型名称 |
| port | 服务器监听端口 |

## 编译与运行
```bash
go build
./img-mcp-server
```

## MCP 工具

### analyze_image
分析图片内容并返回文字描述。

**参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|:----:|------|
| base64 | string | ✅ | 图片的 Base64 编码 |
| prompt | string | ✅ | 对图片的分析要求或提示词 |

## 许可证
MIT