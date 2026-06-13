---
name: image-analysis
description: 当用户或你需要分析描述本地图片内容时通过MCP服务器上传并获取AI分析的结果
---

# 图片分析

## 前置条件
- MCP服务器已经配置好（默认 http://localhost:8080）
- 如果不知道或者不通，通知用户修改SKILL内容
- 如果返回Unauthorized, 说明目标服务器配置了Token, 请咨询用户, 并要求用户填写到skill中

## 使用流程
1. 首先通过curl向MCP服务器的/upload发送需要分析的图片，图片格式需要png、jpg、webp
```bash
curl -s -X POST <服务器地址>/upload -F "file=@<图片路径>"
```
2. 查看MCP服务器返回的ID
从返回的JSON'{"cache_id": "xxx"}' 中提取 cache_id 值
3. 使用MCP协议来请求MCP服务器对图片进行分析并接收MCP服务器分析的内容
调用 MCP 工具 analyze_image，参数：
- cache_id: 上一步获取的 ID
- prompt: 你或用户对图片的提问

## 错误处理
- 上传失败请检查网络连接和服务器运行情况（和用户反映）
- 如果 MCP 服务器返回“缓存ID不存在或已过期”，自动上传同一张照片，不用反馈给用户，再使用新的 cache_id 重试
- 格式不支持请提前转换为支持的格式，建议使用jpg或png
